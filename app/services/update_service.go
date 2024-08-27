package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/tanapoln/capgo-server/app/db"
	"github.com/tanapoln/capgo-server/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	NilLatestResult = GetLatestResult{}
	cacheStore      = (func() *cache.Cache {
		dur := config.Get().CacheResultDuration
		return cache.New(dur, dur+time.Minute*5)
	})()
)

type UpdateService struct {
}

func (svc *UpdateService) GetLatest(ctx context.Context, query GetLatestQuery) (GetLatestResult, error) {
	if !query.IsValid() {
		slog.Info("GetLatestQuery is invalid", "query", query)
		return NilLatestResult, ErrGetLatestQueryInvalid
	}

	doFind := func() (GetLatestResult, error) {
		var release db.Release
		err := db.Collections().Releases().FindOne(ctx, bson.M{
			"platform":     query.Platform,
			"app_id":       query.AppID,
			"version_name": query.VersionName,
			"version_code": query.VersionCode,
		}).Decode(&release)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return NilLatestResult, ErrBundleNotFound
			}
			return NilLatestResult, err
		}

		result := GetLatestResult{
			Builtin: true,
		}
		bundleID := release.BuiltinBundleID
		if release.ActiveBundleID != nil && !release.ActiveBundleID.IsZero() {
			bundleID = *release.ActiveBundleID
			result.Builtin = false
		}
		if bundleID.IsZero() {
			return NilLatestResult, ErrInvalidBundleForRelease
		}

		var bundle db.Bundle
		err = db.Collections().Bundles().FindOne(ctx, bson.M{"_id": bundleID}).Decode(&bundle)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return NilLatestResult, ErrBundleNotFound
			}
			return NilLatestResult, err
		}
		result.Bundle = bundle

		return result, nil
	}

	val, found := cacheStore.Get(query.cacheKey())
	if found {
		switch v := val.(type) {
		case GetLatestResult:
			return v, nil
		case error:
			return NilLatestResult, v
		default:
			return NilLatestResult, ErrCacheInvalid
		}
	}

	result, err := doFind()
	if err != nil {
		cacheStore.Set(query.cacheKey(), err, 5*time.Minute)
		return NilLatestResult, err
	}
	cacheStore.Set(query.cacheKey(), result, cache.DefaultExpiration)
	return result, nil
}

func (svc *UpdateService) CreateBundleDownloadURL(ctx context.Context, bundleID primitive.ObjectID) (*url.URL, error) {
	//TODO: secure bundle download url will be implemented in the future
	return nil, nil
}

type GetLatestQuery struct {
	AppID       string
	Platform    db.Platform
	VersionName string
	VersionCode string
}

func (c GetLatestQuery) IsValid() bool {
	return c.AppID != "" && c.Platform != "" && c.VersionName != "" && c.VersionCode != ""
}

func (c GetLatestQuery) cacheKey() string {
	return fmt.Sprintf("%s|%s|%s|%s", c.Platform, c.VersionName, c.VersionCode, c.AppID)
}

type GetLatestResult struct {
	Bundle  db.Bundle
	Builtin bool
}

func (r GetLatestResult) VersionName() string {
	if r.Builtin {
		return "builtin"
	} else {
		return r.Bundle.VersionName
	}
}

func (r GetLatestResult) Checksum() string {
	return r.Bundle.CRC
}

func (r GetLatestResult) PublicDownloadURL() string {
	return r.Bundle.PublicDownloadURL
}

func (r GetLatestResult) Signature() string {
	return r.Bundle.Signature
}
