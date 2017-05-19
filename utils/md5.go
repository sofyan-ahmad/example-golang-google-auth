package utils

import (
	"crypto/md5"
	"fmt"
	"io"
)

func HashPassword(source string) string {
	h := md5.New()
	io.WriteString(h, source)
	result := fmt.Sprintf("%x", md5.Sum(nil))

	return result
}
