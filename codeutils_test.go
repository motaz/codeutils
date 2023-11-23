package codeutils

import (
	"fmt"
	"testing"
)

type Name struct {
	Name string
}

func TestTitle(t *testing.T) {

	key := GetMD5("simple")
	enc, err := EncryptText(key, "https://shahid.mbc.net/en/widgets/deal-landing/N2ADSPGOF1?code=B1ALcWn5dMn13123456")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(enc)
		key := GetMD5("simple")

		text, err := DecryptText(key, enc)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(text)
		}
	}

}
