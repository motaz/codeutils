package codeutils

import (
	"fmt"
	"testing"
)

type Name struct {
	Name string
}

func TestTitle(t *testing.T) {

	num := FormatFloatCommas(1236689.799, 5)
	fmt.Println(num)

}
