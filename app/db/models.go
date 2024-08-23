package db

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bundle struct {
	ID          primitive.ObjectID `bson:"_id"`
	VersionName string             `bson:"version_name"`
	Description string             `bson:"description"`
	CRC         string             `bson:"crc_checksum"`
	//Signature is a signature of the bundle, signed with SHA512 RSA public key that configured in the app. Can be empty if not use
	Signature         string    `bson:"signature"`
	PublicDownloadURL string    `bson:"public_download_url"` //a quick MVP solution for capgo
	CreatedAt         time.Time `bson:"created_at"`
}

type Release struct {
	ID       primitive.ObjectID `bson:"_id"`
	Platform Platform           `bson:"platform"`
	AppID    string             `bson:"app_id"`

	// VersionName is release version, mostly semver is used. Usually, it's shown to the user. For example, 1.5.0.
	// Android is same as <manifest versionName="...">
	// iOS is Bundle.main.infoDictionary['CFBundleShortVersionString']
	VersionName string `bson:"version_name"`

	// VersionCode is usually an increment build number. Mostly, it's used internally to track a newer build.
	// Android is same as <manifest versionCode="...">.
	// iOS is Bundle.main.infoDictionary['CFBundleVersion']
	VersionCode  string     `bson:"version_code"`
	ReleasedDate *time.Time `bson:"released_date"`

	// BuiltinBundleID is a bundle ID that's already embedded into released executable.
	BuiltinBundleID primitive.ObjectID `bson:"builtin_bundle_id"`
	// ActiveBundleID is a bundle ID that's app must be used.
	ActiveBundleID *primitive.ObjectID `bson:"active_bundle_id"`

	UpdatedAt time.Time `bson:"updated_at"`
	CreatedAt time.Time `bson:"created_at"`
}

type Platform string

const (
	PlatformAndroid Platform = "android"
	PlatformIOS     Platform = "ios"
)

func ParsePlatform(val string) (Platform, error) {
	s := strings.TrimSpace(strings.ToLower(val))
	switch s {
	case "android":
		return PlatformAndroid, nil
	case "ios":
		return PlatformIOS, nil
	default:
		return "", errors.New("invalid platform: " + val)
	}
}
