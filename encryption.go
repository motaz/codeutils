package codeutils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"strings"
)

func padding16(text, fill string) (out string) {

	for len(text) < 16 {
		text += fill
	}
	out = text
	return
}

func EncryptText(akey, plaintext string) (crypted string, err error) {

	crypted = ""
	for i := 0; i < 50000; i++ {
		var part string
		if len(plaintext) > 16 {
			part = plaintext[:16]

		} else {
			part = plaintext
		}
		plaintext = plaintext[len(part):]
		part = padding16(part, "`")
		var enc string
		enc, err = EncryptAES(akey, part)
		crypted += enc
		if plaintext == "" {
			break
		}
	}
	return
}

func DecryptText(akey, crypted string) (plaintext string, err error) {

	plaintext = ""
	for i := 0; i < 50000; i++ {
		var part string
		if len(crypted) > 32 {
			part = crypted[:32]

		} else {
			part = crypted
		}
		crypted = crypted[len(part):]
		var plain string
		plain, err = DecryptAES(akey, part)
		if crypted == "" {
			for len(plain) > 0 && strings.HasSuffix(plain, "`") {
				plain = strings.TrimSuffix(plain, "`")
			}
		}
		plaintext += plain
		if crypted == "" {
			break
		}

	}
	return
}

func EncryptAES(akey string, plaintext string) (crypted string, err error) {

	// create cipher
	key := []byte(akey)
	var c cipher.Block
	c, err = aes.NewCipher(key)
	if err == nil {

		// allocate space for ciphered data
		out := make([]byte, len(plaintext))

		// encrypt
		c.Encrypt(out, []byte(plaintext))
		// return hex string
		crypted = hex.EncodeToString(out)
	}
	return
}

func DecryptAES(akey string, ct string) (plaintext string, err error) {

	key := []byte(akey)
	var ciphertext []byte
	ciphertext, err = hex.DecodeString(ct)
	var c cipher.Block
	c, err = aes.NewCipher(key)
	if err == nil {

		pt := make([]byte, len(ciphertext))
		c.Decrypt(pt, ciphertext)

		plaintext = string(pt[:])
	}
	return
}
