package main

import (
	_ "github.com/steemfans/multiple-upload-to-s3/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"github.com/steemfans/multiple-upload-to-s3/internal/cmd"
)

func main() {
	err := cmd.Main.AddCommand(
		cmd.Web,
		cmd.S3Put,
	)
	if err != nil {
		panic(err)
	}
	cmd.Main.Run(gctx.New())
}
