package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"github.com/steemfans/multiple-upload-to-s3/internal/consts"
	"github.com/steemfans/multiple-upload-to-s3/internal/controller"
	"github.com/steemfans/multiple-upload-to-s3/internal/logic"
	"github.com/steemfans/multiple-upload-to-s3/internal/service/sqlite"
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
	CreateDB = &gcmd.Command{
		Name:        "create-db",
		Brief:       "create an empty db",
		Description: "Because GoFrame gets dao structure from DB. So we should create an empty db before running `gf gen dao`",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			err = sqlite.GenerateNewDb(consts.EMPTY_DB_PATH, consts.EMPTY_DB_NAME)
			if err != nil {
				panic(err)
			}
			return
		},
	}
	S3Put = &gcmd.Command{
		Name:        "s3put",
		Brief:       "upload file to s3",
		Description: "upload file to s3",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			err = logic.S3Put(ctx, parser)
			if err != nil {
				panic(err)
			}
			return
		},
	}
)
