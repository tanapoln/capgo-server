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
		var req UploadBundleRequest
		if err := ctx.Bind(&req); err != nil {
			return nil, fmt.Errorf("failed to bind request: %v", err)
		}
		if err := req.IsValid(); err != nil {
			return nil, err
		}

		crc, err := calculateCRC(req.Bundle)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate CRC: %v", err)
		}

		publicDownloadURL, err := saveFileToS3Public(ctx.Request.Context(), req.VersionName, req.Bundle)
		if err != nil {
			return nil, fmt.Errorf("failed to save file: %v", err)
		}

		bundle := db.Bundle{
			ID:                primitive.NewObjectID(),
			VersionName:       req.VersionName,
			Description:       req.Description,
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
		var req CreateReleaseRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			return nil, fmt.Errorf("invalid request body: %v", err)
		}
		if err := req.IsValid(); err != nil {
			return nil, err
		}

		var bundle db.Bundle
		err := db.Collections().Bundles().FindOne(ctx.Request.Context(), bson.M{"_id": req.GetBuiltinBundleID()}).Decode(&bundle)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch bundle id: %v", req.BuiltinBundleID)
		}

		release := db.Release{
			ID:              primitive.NewObjectID(),
			Platform:        req.GetPlatform(),
			AppID:           req.AppID,
			VersionName:     req.VersionName,
			VersionCode:     req.VersionCode,
			BuiltinBundleID: bundle.ID,
			UpdatedAt:       time.Now(),
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
		if err := req.IsValid(); err != nil {
			return nil, err
		}

		var release db.Release
		err := db.Collections().Releases().FindOne(
			ctx.Request.Context(),
			bson.M{
				"_id": req.GetReleaseID(),
			},
		).Decode(&release)
		if err != nil {
			return nil, fmt.Errorf("failed to find release id: %v", req.ReleaseID)
		}

		if req.ReleaseDate != nil {
			release.ReleasedDate = req.ReleaseDate
		}
		release.UpdatedAt = time.Now()

		result, err := db.Collections().Releases().UpdateOne(
			ctx.Request.Context(),
			bson.M{"_id": release.ID},
			bson.M{"$set": release},
		)
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

		if err := req.IsValid(); err != nil {
			return nil, err
		}

		var bundle db.Bundle
		err := db.Collections().Bundles().FindOne(ctx.Request.Context(), bson.M{"_id": req.GetBundleID()}).Decode(&bundle)
		if err != nil {
			return nil, fmt.Errorf("failed to find bundle id: %v", req.BundleID)
		}

		var release db.Release
		err = db.Collections().Releases().FindOne(ctx.Request.Context(), bson.M{"_id": req.GetReleaseID()}).Decode(&release)
		if err != nil {
			return nil, fmt.Errorf("failed to find release id: %v", req.ReleaseID)
		}

		result, err := db.Collections().Releases().UpdateOne(
			ctx.Request.Context(),
			bson.M{"_id": release.ID},
			bson.M{
				"$set": bson.M{
					"active_bundle_id": bundle.ID,
					"updated_at":       time.Now(),
				},
			},
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
func saveFileToS3Public(ctx context.Context, versionName string, file *multipart.FileHeader) (string, error) {
	filename := versionName + "_" + xid.New().String() + ".zip"

	r, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer r.Close()

	uploader := s3ext.NewUploader()
	result, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(config.Get().S3Bucket),
		Key:    aws.String(filename),
		Body:   r,
		ACL:    types.ObjectCannedACLPublicRead,
		Metadata: map[string]string{
			"Content-Type": file.Header.Get("Content-Type"),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	// Return the public URL of the uploaded file
	return result.Location, nil
}

func calculateCRC(file *multipart.FileHeader) (string, error) {
	r, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer r.Close()

	hash := crc32.NewIEEE()
	_, err = io.Copy(hash, r)
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
		AppID:           release.AppID,
		VersionName:     release.VersionName,
		Platform:        string(release.Platform),
		VersionCode:     release.VersionCode,
		BuiltinBundleID: release.BuiltinBundleID.Hex(),
		UpdatedAt:       release.UpdatedAt,
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
