package utility_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/steemfans/multiple-upload-to-s3/utility"
)

func TestMd5(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expectStr := "XrY7u+Ae7tCTyyK7j1rNww=="
		buf := []byte("hello world")
		result := utility.Md5(&buf)
		t.Assert(expectStr, result)
	})
}
