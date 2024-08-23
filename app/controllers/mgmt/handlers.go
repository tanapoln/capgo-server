package mgmt

import (
	"context"
	"fmt"
	"hash/crc32"
	"io"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/tanapoln/capgo-server/app/controllers/utils"
	"github.com/tanapoln/capgo-server/app/db"
	"github.com/tanapoln/capgo-server/app/external/s3ext"
	"github.com/tanapoln/capgo-server/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewCapgoManagementController() *CapgoManagementController {
	return &CapgoManagementController{}
}

type CapgoManagementController struct {
}

func (ctrl *CapgoManagementController) UploadBundle(ctx *gin.Context) {
	utils.Handle(ctx, func() (interface{}, error) {
		file, header, err := ctx.Request.FormFile("bundle")
		if err != nil {
			return nil, fmt.Errorf("file upload error: %v", err)
		}
		defer file.Close()

		if header.Header.Get("Content-Type") != "application/zip" {
			return nil, fmt.Errorf("invalid file type: only zip files are allowed")
		}

		crc, err := calculateCRC(file)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate CRC: %v", err)
		}

		versionName := ctx.PostForm("version_name")
		publicDownloadURL, err := saveFileToS3Public(ctx.Request.Context(), versionName, header)
		if err != nil {
			return nil, fmt.Errorf("failed to save file: %v", err)
		}

		bundle := db.Bundle{
			ID:                primitive.NewObjectID(),
			VersionName:       versionName,
			Description:       ctx.PostForm("description"),
			CRC:               crc,
			PublicDownloadURL: publicDownloadURL,
			CreatedAt:         time.Now(),
		}

		err = saveBundleToDatabase(ctx.Request.Context(), bundle)
		if err != nil {
			return nil, fmt.Errorf("failed to save bundle to database: %v", err)
		}

		return gin.H{
			"message": "Bundle uploaded successfully",
			"bundle":  mapBundleToResponse(bundle),
		}, nil
	})
}

func (ctrl *CapgoManagementController) CreateRelease(ctx *gin.Context) {
	utils.Handle(ctx, func() (interface{}, error) {
		var reqBody CreateReleaseRequest
		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			return nil, fmt.Errorf("invalid request body: %v", err)
		}

		if !reqBody.IsValid() {
			return nil, fmt.Errorf("invalid request data")
		}

		var bundle db.Bundle
		err := db.Collections().Bundles().FindOne(ctx.Request.Context(), bson.M{"_id": reqBody.BuiltinBundleID}).Decode(&bundle)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch bundle id: %v", reqBody.BuiltinBundleID)
		}

		release := db.Release{
			ID:              primitive.NewObjectID(),
			Platform:        reqBody.GetPlatform(),
			AppID:           reqBody.AppID,
			VersionName:     reqBody.VersionName,
			VersionCode:     reqBody.VersionCode,
			BuiltinBundleID: bundle.ID,
			CreatedAt:       time.Now(),
		}

		_, err = db.Collections().Releases().InsertOne(ctx.Request.Context(), release)
		if err != nil {
			return nil, fmt.Errorf("failed to create release: %v", err)
		}

		return gin.H{
			"message": "Release created successfully",
			"release": mapReleaseToResponse(release),
		}, nil
	})
}

func (ctrl *CapgoManagementController) UpdateRelease(ctx *gin.Context) {
	utils.Handle(ctx, func() (interface{}, error) {
		var req UpdateReleaseRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			return nil, fmt.Errorf("failed to bind request: %v", err)
		}
		if !req.IsValid() {
			return nil, fmt.Errorf("invalid request body")
		}

		var release db.Release
		err := db.Collections().Releases().FindOne(ctx.Request.Context(), bson.M{"_id": req.ReleaseID}).Decode(&release)
		if err != nil {
			return nil, fmt.Errorf("failed to find release id: %v", req.ReleaseID)
		}

		if req.ReleaseDate != nil {
			release.ReleasedDate = req.ReleaseDate
		}

		result, err := db.Collections().Releases().UpdateOne(ctx.Request.Context(), bson.M{"_id": release.ID}, bson.M{"$set": release})
		if err != nil {
			return nil, fmt.Errorf("failed to update release: %v", err)
		}
		if result.ModifiedCount == 0 {
			return nil, fmt.Errorf("failed to update release, no affected. release id: %v", release.ID.Hex())
		}

		return gin.H{
			"message": "Release updated successfully",
			"release": mapReleaseToResponse(release),
		}, nil
	})
}

func (ctrl *CapgoManagementController) ListAllBundles(ctx *gin.Context) {
	utils.Handle(ctx, func() (interface{}, error) {
		cursor, err := db.Collections().Bundles().Find(
			ctx.Request.Context(), bson.M{},
			options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
		if err != nil {
			return nil, fmt.Errorf("failed to fetch bundles: %v", err)
		}
		defer cursor.Close(ctx.Request.Context())

		var bundles []db.Bundle
		if err = cursor.All(ctx.Request.Context(), &bundles); err != nil {
			return nil, fmt.Errorf("failed to decode bundles: %v", err)
		}

		response := make([]BundleResponse, len(bundles))
		for i, bundle := range bundles {
			response[i] = mapBundleToResponse(bundle)
		}

		return ListAllBundlesResponse{
			Data: response,
		}, nil
	})
}

func (ctrl *CapgoManagementController) ListAllReleases(ctx *gin.Context) {
	utils.Handle(ctx, func() (interface{}, error) {
		cursor, err := db.Collections().Releases().Find(
			ctx.Request.Context(), bson.M{},
			options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
		if err != nil {
			return nil, fmt.Errorf("failed to fetch bundles: %v", err)
		}
		defer cursor.Close(ctx.Request.Context())

		var releases []db.Release
		if err = cursor.All(ctx.Request.Context(), &releases); err != nil {
			return nil, fmt.Errorf("failed to decode bundles: %v", err)
		}

		response := make([]ReleaseResponse, len(releases))
		for i, release := range releases {
			response[i] = mapReleaseToResponse(release)
		}

		return ListAllReleasesResponse{
			Data: response,
		}, nil
	})
}

func (ctrl *CapgoManagementController) SetReleaseActiveBundle(ctx *gin.Context) {
	utils.Handle(ctx, func() (interface{}, error) {
		var req SetReleaseActiveBundleRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			return nil, fmt.Errorf("failed to bind request: %v", err)
		}

		if !req.IsValid() {
			return nil, fmt.Errorf("invalid request body")
		}

		var bundle db.Bundle
		err := db.Collections().Bundles().FindOne(ctx.Request.Context(), bson.M{"_id": req.BundleID}).Decode(&bundle)
		if err != nil {
			return nil, fmt.Errorf("failed to find bundle id: %v", req.BundleID)
		}

		var release db.Release
		err = db.Collections().Releases().FindOne(ctx.Request.Context(), bson.M{"_id": req.ReleaseID}).Decode(&release)
		if err != nil {
			return nil, fmt.Errorf("failed to find release id: %v", req.ReleaseID)
		}

		result, err := db.Collections().Releases().UpdateOne(
			ctx.Request.Context(),
			bson.M{"_id": release.ID},
			bson.M{"$set": bson.M{"active_bundle_id": bundle.ID}},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update release: %v", err)
		}
		if result.ModifiedCount == 0 {
			return nil, fmt.Errorf("failed to update release, no affected. release id: %v", release.ID.Hex())
		}

		return gin.H{
			"message": "Release updated successfully",
		}, nil
	})
}

// Placeholder functions (implement these according to your actual storage and database setup)
func saveFileToS3Public(ctx context.Context, versionName string, header *multipart.FileHeader) (string, error) {
	filename := versionName + "_" + xid.New().String() + ".zip"

	file, err := header.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	uploader := s3ext.NewUploader()
	result, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(config.Get().S3Bucket),
		Key:    aws.String(filename),
		Body:   file,
		ACL:    types.ObjectCannedACLPublicRead,
		Metadata: map[string]string{
			"Content-Type": header.Header.Get("Content-Type"),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	// Return the public URL of the uploaded file
	return result.Location, nil
}

func calculateCRC(file io.Reader) (string, error) {
	hash := crc32.NewIEEE()
	_, err := io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	checksum := hash.Sum32()
	crc := fmt.Sprintf("%08x", checksum)

	return crc, nil
}

func saveBundleToDatabase(ctx context.Context, bundle db.Bundle) error {
	_, err := db.Collections().Bundles().InsertOne(ctx, bundle)
	if err != nil {
		return err
	}
	return nil
}

func mapBundleToResponse(bundle db.Bundle) BundleResponse {
	return BundleResponse{
		ID:                bundle.ID.Hex(),
		VersionName:       bundle.VersionName,
		Description:       bundle.Description,
		CRC:               bundle.CRC,
		PublicDownloadURL: bundle.PublicDownloadURL,
		CreatedAt:         bundle.CreatedAt,
	}
}

func mapReleaseToResponse(release db.Release) ReleaseResponse {
	r := ReleaseResponse{
		ID:              release.ID.Hex(),
		VersionName:     release.VersionName,
		Platform:        string(release.Platform),
		VersionCode:     release.VersionCode,
		BuiltinBundleID: release.BuiltinBundleID.Hex(),
		CreatedAt:       release.CreatedAt,
	}
	if release.ReleasedDate != nil {
		r.ReleaseDate = release.ReleasedDate
	}
	if release.ActiveBundleID != nil {
		s := release.ActiveBundleID.Hex()
		r.ActiveBundleID = &s
	}
	return r
}
