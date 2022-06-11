package codeutils

import (
	"fmt"

	"testing"
)

type Name struct {
	Name string
}

func TestTitle(t *testing.T) {

	headers := make(map[string]string)
	headers["content-type"] = "application/json"
	var aName Name
	aName.Name = "Motaz"

	result := CallURLAsPost("http://localhost", nil, 10)

	fmt.Printf("Status code: %d\nerror: %s\nContents:\n%s", result.StatusCode, result.Err, result.Content)
}
