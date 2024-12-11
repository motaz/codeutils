package codeutils

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"os"
	"strings"
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

func getFullPathConfigFile(filename string) (fullname string) {

	fullname = filename
	if !strings.Contains(fullname, string(os.PathSeparator)) {
		fullname = GetCurrentAppDir() + string(os.PathSeparator) + fullname
	}
	return
}
