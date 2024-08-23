package mgmt

import (
	"mime/multipart"
	"time"

	"github.com/tanapoln/capgo-server/app/db"
)

type UploadBundleRequest struct {
	Bundle      multipart.File `form:"bundle"`
	AppID       string         `form:"appId"`
	Platform    string         `form:"platform"`
	VersionName string         `form:"versionName"`
	VersionCode string         `form:"versionCode"`
}

type BundleResponse struct {
	ID                string    `json:"id"`
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
	Platform        string     `json:"platform"`
	VersionName     string     `json:"version_name"`
	VersionCode     string     `json:"version_code"`
	ReleaseDate     *time.Time `json:"release_date"`
	BuiltinBundleID string     `json:"builtin_bundle_id"`
	ActiveBundleID  *string    `json:"active_bundle_id"`
	CreatedAt       time.Time  `json:"created_at"`
}

type ListAllReleasesResponse struct {
	Data []ReleaseResponse `json:"data"`
}

type SetReleaseActiveBundleRequest struct {
	ReleaseID string `json:"release_id"`
	BundleID  string `json:"bundle_id"`
}

func (s *SetReleaseActiveBundleRequest) IsValid() bool {
	return s.ReleaseID != "" && s.BundleID != ""
}

type CreateReleaseRequest struct {
	Platform        string `json:"platform"`
	AppID           string `json:"app_id"`
	VersionName     string `json:"version_name"`
	VersionCode     string `json:"version_code"`
	BuiltinBundleID string `json:"builtin_bundle_id"`
}

func (req *CreateReleaseRequest) IsValid() bool {
	if _, err := db.ParsePlatform(req.Platform); err != nil {
		return false
	}
	return req.AppID != "" && req.VersionName != "" && req.VersionCode != "" && req.BuiltinBundleID != ""
}

func (req *CreateReleaseRequest) GetPlatform() db.Platform {
	p, _ := db.ParsePlatform(req.Platform)
	return p
}

type UpdateReleaseRequest struct {
	ReleaseID   string     `json:"release_id"`
	ReleaseDate *time.Time `json:"release_date"`
}

func (req *UpdateReleaseRequest) IsValid() bool {
	if req.ReleaseID == "" {
		return false
	}
	if req.ReleaseDate != nil {
		if req.ReleaseDate.IsZero() {
			return false
		}
	}
	return true
}
