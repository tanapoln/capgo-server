package capgo

import (
	"github.com/tanapoln/capgo-server/app/db"
)

type UpdateRequest struct {
	Platform string `json:"platform"`
	DeviceID string `json:"device_id"`
	AppID    string `json:"app_id"`
	CustomID string `json:"custom_id"`

	// VersionBuild from application. For Android, it is same as <manifest versionName="...">
	VersionBuild string `json:"version_build"`

	// VersionCode from application. For Android, it is same as <manifest versionCode="...">
	VersionCode string `json:"version_code"`

	VersionOS      string `json:"version_os"`
	VersionName    string `json:"version_name"`
	PluginVersion  string `json:"plugin_version"`
	IsEmulator     bool   `json:"is_emulator"`
	IsProd         bool   `json:"is_prod"`
	DefaultChannel string `json:"defaultChannel"`
}

func (c *UpdateRequest) IsValid() bool {
	_, err := db.ParsePlatform(c.Platform)
	if err != nil {
		return false
	}
	return c.Platform != "" && c.DeviceID != "" && c.AppID != "" && c.VersionBuild != "" && c.VersionCode != ""
}

func (c *UpdateRequest) GetPlatform() db.Platform {
	p, _ := db.ParsePlatform(c.Platform)
	return p
}

type UpdateWithNewMinorVersionResponse struct {
	// Version is a new version string. Capgo will download from URL if this version string doesn't equal to current version
	Version string `json:"version"`
	// URL is a zipped bundle download url
	URL string `json:"url"`
	// SessionKey is Base64 IV + Cipher AES key. Use for decrypt the bundle (encrypted with private key embedded in the app). Can be empty if not use
	SessionKey string `json:"sessionKey"`
	//Checksum is CRC checksum of the bundle
	Checksum string `json:"checksum"`
	//Signature is a signature of the bundle, signed with SHA512 RSA public key that configured in the app. Can be empty if not use
	Signature string `json:"signature"`
}

type CapgoErrorResponse struct {
	Error string `json:"error"`
}

type CapgoIncorrectWithMessageResponse struct {
	Message string `json:"message"`
}

type UpdateBreakingChangeVersionResponse struct {
	// Message is a message to show to user about new major breaking change version.
	Message string `json:"message"`
	// Major is true if this is a major breaking change version. In most case, it is true.
	Major bool `json:"major"`
	//Version is a new major breaking change version.
	Version string `json:"version"`
}
