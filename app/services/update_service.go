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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	NilBundle  = db.Bundle{}
	cacheStore = cache.New(10*time.Minute, 15*time.Minute)
)

type UpdateService struct {
}

func (svc *UpdateService) GetLatest(ctx context.Context, query GetLatestQuery) (db.Bundle, error) {
	if !query.IsValid() {
		slog.Info("GetLatestQuery is invalid", "query", query)
		return NilBundle, ErrGetLatestQueryInvalid
	}

	doFind := func() (db.Bundle, error) {
		var release db.Release
		err := db.Collections().Releases().FindOne(ctx, bson.M{
			"platform":     query.Platform,
			"app_id":       query.AppID,
			"version_name": query.VersionName,
			"version_code": query.VersionCode,
		}).Decode(&release)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return NilBundle, ErrBundleNotFound
			}
			return NilBundle, err
		}

		bundleID := release.BuiltinBundleID
		if release.ActiveBundleID != nil && !release.ActiveBundleID.IsZero() {
			bundleID = *release.ActiveBundleID
		}
		if bundleID.IsZero() {
			return NilBundle, ErrInvalidBundleForRelease
		}

		var bundle db.Bundle
		err = db.Collections().Bundles().FindOne(ctx, bson.M{"_id": bundleID}).Decode(&bundle)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return NilBundle, ErrBundleNotFound
			}
			return NilBundle, err
		}

		return bundle, nil
	}

	val, found := cacheStore.Get(query.cacheKey())
	if found {
		switch v := val.(type) {
		case db.Bundle:
			return v, nil
		case error:
			return NilBundle, v
		default:
			return NilBundle, ErrCacheInvalid
		}
	}

	bundle, err := doFind()
	if err != nil {
		cacheStore.Set(query.cacheKey(), err, 5*time.Minute)
		return NilBundle, err
	}
	cacheStore.Set(query.cacheKey(), bundle, cache.DefaultExpiration)
	return bundle, nil
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
