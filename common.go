package codeutils

import (
	"crypto/md5"
	"fmt"
)

func GetMD5(text string) string {

	data := []byte(text)
	return fmt.Sprintf("%x", md5.Sum(data))

}
