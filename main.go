package main

import (
	_ "multiple-upload-to-s3/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"multiple-upload-to-s3/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.New())
}
