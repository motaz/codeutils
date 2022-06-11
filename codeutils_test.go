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

	result, status, err := CallURLAsGet("http://localhost", 10)

	fmt.Printf("status code: %d\nerror: %s\nContents:\n%s", status, err, result)
}
