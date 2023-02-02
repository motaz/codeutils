package codeutils

import (
	"fmt"
	"testing"
)

type Name struct {
	Name string
}

func TestTitle(t *testing.T) {

	num := FormatFloatCommas(34512340001.12309, 2)
	fmt.Println(num)
	num = FormatCommas(1922000)
	fmt.Println(num)
}
