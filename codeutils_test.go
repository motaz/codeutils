package codeutils

import (
	"fmt"
	"testing"
)

type Name struct {
	Name string
}

func TestTitle(t *testing.T) {

	num := FormatFloatCommas(12340001.12, 1)
	fmt.Println(num)
}
