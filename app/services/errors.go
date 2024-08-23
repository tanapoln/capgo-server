package services

import "errors"

var ErrReleaseNotFound = errors.New("release is found")
var ErrInvalidBundleForRelease = errors.New("release contains invalid bundle id")
var ErrBundleNotFound = errors.New("bundle is not found")
var ErrGetLatestQueryInvalid = errors.New("GetLatest information is incompleted, cannot find the latest release")
var ErrCacheInvalid = errors.New("unexpected error (cache issue)")
