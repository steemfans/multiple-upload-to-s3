package controller

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	v1 "github.com/steemfans/multiple-upload-to-s3/api/v1"
)

var (
	Hello = cHello{}
)

type cHello struct{}

func (c *cHello) Hello(ctx context.Context, req *v1.HelloReq) (res *v1.HelloRes, err error) {
	g.RequestFromCtx(ctx).Response.Writeln("Hello World!")
	return
}
