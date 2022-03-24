package codeutils

import (
	"testing"
)

func TestTitle(t *testing.T) {
	SetLogType(WEEKLYLOG)
	WriteToLog("Third line", "test")
	println(GetCurrentAppDir())
	println(GetMD5("help21"))
}
