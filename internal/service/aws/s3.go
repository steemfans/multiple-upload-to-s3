package aws

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gogf/gf/v2/os/gctx"
)

type S3 struct {
}

func (s *S3) Download(ctx gctx.Ctx, bucketName, fileKey, dst string) (err error) {
	cfg, err := LoadConfig()
	if err != nil {
		return
	}
	client := s3.NewFromConfig(cfg)
	f, err := os.Create(dst)
	if err != nil {
		return
	}
	downloader := manager.NewDownloader(client)
	_, err = downloader.Download(context.TODO(), f, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
	})
	return
}

func (s *S3) Upload(ctx gctx.Ctx, bucketName, fileKey, src string) (output *manager.UploadOutput, err error) {
	cfg, err := LoadConfig()
	if err != nil {
		return
	}
	client := s3.NewFromConfig(cfg)
	f, err := os.Open(src)
	if err != nil {
		return
	}
	uploader := manager.NewUploader(client)
	output, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
		Body:   f,
	})
	return
}

var SS3 = S3{}
