package s3ext

import (
	"context"
	"log"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tanapoln/capgo-server/config"
)

var (
	Client *s3.Client
)

func init() {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var s3Opts []func(*s3.Options)
	if config.Get().S3BaseEndpoint != nil {
		cfg.BaseEndpoint = config.Get().S3BaseEndpoint
		s3Opts = append(s3Opts, func(o *s3.Options) {
			o.UsePathStyle = true
		})
	}

	// Create an Amazon S3 service client
	Client = s3.NewFromConfig(cfg, s3Opts...)
}

func NewUploader() *manager.Uploader {
	return manager.NewUploader(Client)
}
