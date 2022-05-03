package codeutils

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

func GetMD5(text string) string {

	data := []byte(text)
	return fmt.Sprintf("%x", md5.Sum(data))

}

func GetRandom(r int) int {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(r)
}
