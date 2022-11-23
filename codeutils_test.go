package codeutils

import (
	"fmt"

	"testing"
)

type Name struct {
	Name string
}

func TestTitle(t *testing.T) {

	//result := CallURLAsGet("http://localhost", 10)
	req, _ := PrepareURLCall("http://localhost/sources", "GET", nil)
	result := CallURL(req, 10)
	fmt.Printf("%d\n%+s", result.StatusCode, string(result.Content))
}
