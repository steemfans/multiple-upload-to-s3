package utility

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
)

func GetFileType(src string) (fileType string, err error) {
	file, err := os.Open(src)
	if err != nil {
		fmt.Printf("err opening file: %s", err)
		return
	}
	defer file.Close()
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	fileType = http.DetectContentType(buf)
	return
}

func Md5(content *([]byte)) (result string) {
	hash := md5.Sum(*content)
	return base64.StdEncoding.EncodeToString(hash[:])
}
