package s3ext

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	Client *s3.Client
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	Client = s3.NewFromConfig(cfg)
}

func NewUploader() *manager.Uploader {
	return manager.NewUploader(Client)
}
