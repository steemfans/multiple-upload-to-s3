package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/glog"

	"github.com/steemfans/multiple-upload-to-s3/internal/controller"
	"github.com/steemfans/multiple-upload-to-s3/internal/service/aws"
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
	S3Put = &gcmd.Command{
		Name:        "s3put",
		Brief:       "upload file to s3",
		Description: "upload file to s3",
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
			bucketName := p.GetOpt("bucketname")
			objectName := p.GetOpt("objectname")

			output, err := aws.SS3.Upload(ctx, bucketName.String(), objectName.String(), src.String())
			if err != nil {
				glog.Fatal(ctx, "Upload Failed: ", err)
			}
			glog.Info(ctx, output)
			return
		},
	}
)
