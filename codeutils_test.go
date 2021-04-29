package codeutils

import (
	"fmt"
	"testing"
)

func TestTitle(t *testing.T) {
	aname, err := getConfigValue("test.ini", "name")
	if err != nil {
		fmt.Println("Error: " + err.Error())
	} else {
		fmt.Println("aname: ", aname)
	}
	WriteToLog("test", "test")
}
