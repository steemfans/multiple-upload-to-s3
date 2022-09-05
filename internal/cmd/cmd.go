package cmd

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/steemfans/multiple-upload-to-s3/internal/controller"
)

var (
	Main = &gcmd.Command{
		Name:        "muts",
		Brief:       "Multiple Upload To S3",
		Description: "Multiple Upload To S3",
	}
	Web = &gcmd.Command{
		Name:        "web",
		Brief:       "start web server",
		Description: "This is for starting api web manager service.",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					controller.Hello,
				)
			})
			s.Run()
			return nil
		},
	}
	S3 = &gcmd.Command{
		Name:        "s3",
		Brief:       "s3 transfer",
		Description: "s3 transfer",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			p, err := gcmd.Parse(g.MapStrBool{
				"src":        true,
				"bucketname": true,
				"objectname": true,
				"endpoint":   true,
			})
			if err != nil {
				panic(err)
			}
			src := p.GetOpt("src")
			bucketname := p.GetOpt("bucketname")
			objectname := p.GetOpt("objectname")

			awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
			awsSecretKey := os.Getenv("AWS_SECRET_KEY")
			endpoint := os.Getenv("ENDPOINT")
			if awsAccessKey == "" || awsSecretKey == "" {
				log.Fatalln("need aws key settings")
			}

			s3Client, err := minio.New(endpoint, &minio.Options{
				Creds:  credentials.NewStaticV4(awsAccessKey, awsSecretKey, ""),
				Secure: true,
			})
			if err != nil {
				log.Fatalln(err)
			}

			mtype, err := mimetype.DetectFile(src.String())
			if err != nil {
				log.Fatalln(err)
			}

			fileinfo, err := os.Stat(src.String())
			if err != nil {
				log.Fatalln(err)
			}
			filesize := fileinfo.Size()

			progress := pb.New64(filesize)
			progress.Start()

			s3Ctx, cancel := context.WithTimeout(context.Background(), 7200*time.Second)
			defer cancel()

			if _, err := s3Client.FPutObject(s3Ctx, bucketname.String(), objectname.String(), src.String(), minio.PutObjectOptions{
				ContentType: mtype.String(),
				Progress:    progress,
			}); err != nil {
				log.Fatalln(err)
			}
			log.Println("Successfully uploaded " + src.String())
			return
		},
	}
)
