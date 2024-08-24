package mgmt

import (
	"archive/zip"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/tanapoln/capgo-server/app/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadBundleRequest struct {
	Bundle      *multipart.FileHeader `form:"bundle"`
	AppID       string                `form:"app_id"`
	VersionName string                `form:"version_name"`
	Description string                `form:"description"`
}

func (req *UploadBundleRequest) IsValid() error {
	b := req.Bundle != nil && req.VersionName != "" && req.AppID != ""
	if !b {
		return fmt.Errorf("invalid request body")
	}

	err := req.validateBundleZip()
	if err != nil {
		return err
	}
	return nil
}

func (req *UploadBundleRequest) validateBundleZip() error {
	f, err := req.Bundle.Open()
	if err != nil {
		return fmt.Errorf("invalid bundle zip file: %v", err)
	}
	defer f.Close()

	_, err = zip.NewReader(f, req.Bundle.Size)
	if err != nil {
		return fmt.Errorf("invalid bundle zip file: %v", err)
	}

	return nil
}

type BundleResponse struct {
	ID                string    `json:"id"`
	AppID             string    `json:"app_id"`
	VersionName       string    `json:"version_name"`
	Description       string    `json:"description"`
	CRC               string    `json:"crc_checksum"`
	PublicDownloadURL string    `json:"public_download_url"`
	CreatedAt         time.Time `json:"created_at"`
}

type ListAllBundlesResponse struct {
	Data []BundleResponse `json:"data"`
}

type ReleaseResponse struct {
	ID              string     `json:"id"`
	AppID           string     `json:"app_id"`
	Platform        string     `json:"platform"`
	VersionName     string     `json:"version_name"`
	VersionCode     string     `json:"version_code"`
	ReleaseDate     *time.Time `json:"release_date"`
	BuiltinBundleID string     `json:"builtin_bundle_id"`
	ActiveBundleID  *string    `json:"active_bundle_id"`
	UpdatedAt       time.Time  `json:"updated_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type ListAllReleasesResponse struct {
	Data []ReleaseResponse `json:"data"`
}

type SetReleaseActiveBundleRequest struct {
	ReleaseID string `json:"release_id"`
	BundleID  string `json:"bundle_id"`
}

func (s *SetReleaseActiveBundleRequest) IsValid() error {
	if s.ReleaseID == "" || s.BundleID == "" {
		return fmt.Errorf("invalid request body")
	}
	_, err := primitive.ObjectIDFromHex(s.ReleaseID)
	if err != nil {
		return fmt.Errorf("invalid release id: %v", err)
	}
	_, err = primitive.ObjectIDFromHex(s.BundleID)
	if err != nil {
		return fmt.Errorf("invalid bundle id: %v", err)
	}
	return nil
}

func (s *SetReleaseActiveBundleRequest) GetReleaseID() primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(s.ReleaseID)
	return id
}

func (s *SetReleaseActiveBundleRequest) GetBundleID() primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(s.BundleID)
	return id
}

type CreateReleaseRequest struct {
	Platform        string `json:"platform"`
	AppID           string `json:"app_id"`
	VersionName     string `json:"version_name"`
	VersionCode     string `json:"version_code"`
	BuiltinBundleID string `json:"builtin_bundle_id"`
}

func (req *CreateReleaseRequest) IsValid() error {
	if _, err := db.ParsePlatform(req.Platform); err != nil {
		return fmt.Errorf("invalid platform: %v", err)
	}
	_, err := primitive.ObjectIDFromHex(req.BuiltinBundleID)
	if err != nil {
		return fmt.Errorf("invalid builtin bundle id: %v", err)
	}
	b := req.AppID != "" && req.VersionName != "" && req.VersionCode != "" && req.BuiltinBundleID != ""
	if !b {
		return fmt.Errorf("invalid request body")
	}
	return nil
}

func (req *CreateReleaseRequest) GetBuiltinBundleID() primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(req.BuiltinBundleID)
	return id
}

func (req *CreateReleaseRequest) GetPlatform() db.Platform {
	p, _ := db.ParsePlatform(req.Platform)
	return p
}

type UpdateReleaseRequest struct {
	ReleaseID   string     `json:"release_id"`
	ReleaseDate *time.Time `json:"release_date"`
}

func (req *UpdateReleaseRequest) IsValid() error {
	if req.ReleaseID == "" {
		return fmt.Errorf("missing release id")
	}
	_, err := primitive.ObjectIDFromHex(req.ReleaseID)
	if err != nil {
		return fmt.Errorf("invalid release id: %v", err)
	}
	if req.ReleaseDate != nil {
		if req.ReleaseDate.IsZero() {
			return fmt.Errorf("release date is empty")
		}
	}
	return nil
}

func (req *UpdateReleaseRequest) GetReleaseID() primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(req.ReleaseID)
	return id
}

type DeleteReleaseRequest struct {
	ReleaseID string `json:"release_id"`
}

func (req *DeleteReleaseRequest) IsValid() error {
	if req.ReleaseID == "" {
		return fmt.Errorf("missing release id")
	}
	_, err := primitive.ObjectIDFromHex(req.ReleaseID)
	if err != nil {
		return fmt.Errorf("invalid release id: %v", err)
	}
	return nil
}

func (req *DeleteReleaseRequest) GetReleaseID() primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(req.ReleaseID)
	return id
}
