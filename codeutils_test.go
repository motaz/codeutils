package codeutils

import (
	"testing"
)

type Name struct {
	Name string
}

func TestTitle(t *testing.T) {

	SetLogType(HOURLYLOG)
	WriteToLog("second line", "testhours")
}
