package codeutils

import (
	"testing"
)

func TestTitle(t *testing.T) {
	logdaytype = WEEKLYLOG
	WriteToLog("Third line", "test")
}
