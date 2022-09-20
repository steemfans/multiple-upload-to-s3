package aws

import (
	"bytes"
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/cheggaaa/pb/v3"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/steemfans/multiple-upload-to-s3/utility"
)

var (
	maxRetries = 10
)

func init() {
	maxRetries = genv.Get("MAX_RETRIES", 10).Int()
}

type S3 struct {
}

func (s *S3) GetInstance() (client *s3.Client, err error) {
	cfg, err := LoadConfig()
	if err != nil {
		return
	}
	client = s3.NewFromConfig(cfg)
	return
}

func (s *S3) Download(ctx gctx.Ctx, bucketName, fileKey, dst string) (err error) {
	client, err := s.GetInstance()
	if err != nil {
		return
	}
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
	client, err := s.GetInstance()
	if err != nil {
		return
	}
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

func (s *S3) CreateMultipartUpload(ctx gctx.Ctx, bucketName, fileKey, src string) (output *s3.CreateMultipartUploadOutput, err error) {
	fileType, err := utility.GetFileType(src)
	if err != nil {
		return
	}
	client, err := s.GetInstance()
	if err != nil {
		return
	}
	input := &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fileKey),
		ContentType: aws.String(fileType),
	}
	output, err = client.CreateMultipartUpload(ctx, input)
	return
}

func (s *S3) UploadPart(ctx gctx.Ctx, resp *s3.CreateMultipartUploadOutput, fileBytes *([]byte), partNumber int32) (out *s3.UploadPartOutput, err error) {
	tryNum := 1

	client, err := s.GetInstance()
	if err != nil {
		return
	}

	for tryNum <= maxRetries {
		// init bar
		ioReaderData := bytes.NewReader(*fileBytes)
		bar := pb.Full.Start(len(*fileBytes))
		barReader := bar.NewProxyReader(ioReaderData)

		// init part input
		partInput := &s3.UploadPartInput{
			Body:          barReader,
			Bucket:        resp.Bucket,
			Key:           resp.Key,
			PartNumber:    partNumber,
			UploadId:      resp.UploadId,
			ContentLength: int64(len(*fileBytes)),
			ContentMD5:    aws.String(utility.Md5(fileBytes)),
		}

		// start upload
		out, err := client.UploadPart(ctx, partInput)

		// finish bar
		bar.Finish()

		// err
		if err != nil {
			if tryNum == maxRetries {
				if aerr, ok := err.(awserr.Error); ok {
					return nil, aerr
				}
				return nil, err
			}
			glog.Warningf(ctx, "Retrying to upload part #%v", partNumber)
			tryNum++
		} else {
			glog.Infof(ctx, "Uploaded part #%v", partNumber)
			return out, err
		}
	}
	return nil, nil
}

func (s *S3) CompleteMultipartUpload(ctx gctx.Ctx, resp *s3.CreateMultipartUploadOutput, completedParts []s3types.CompletedPart) (out *s3.CompleteMultipartUploadOutput, err error) {
	completeInput := &s3.CompleteMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
		MultipartUpload: &s3types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	}
	client, err := s.GetInstance()
	if err != nil {
		return
	}
	return client.CompleteMultipartUpload(ctx, completeInput)
}

var SS3 = S3{}
