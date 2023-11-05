package codeutils

import (
	"fmt"
	"testing"
)

type Name struct {
	Name string
}

func TestTitle(t *testing.T) {

	num, _ := StrToDate("2002-10-08")
	fmt.Println(num)

}
