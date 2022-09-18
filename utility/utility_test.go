package utility_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/steemfans/multiple-upload-to-s3/utility"
)

func TestMd5(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		test := "hello world"
		expectStr := "XrY7u+Ae7tCTyyK7j1rNww=="
		result := utility.Md5([]byte(test))
		t.Assert(expectStr, result)
	})
}
