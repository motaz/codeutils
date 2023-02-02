package codeutils

import (
	"fmt"
	"regexp"
	"strconv"
)

func FormatCommas(num int32) string {

	str := fmt.Sprintf("%d", num)
	re := regexp.MustCompile("(\\d+)(\\d{3})")
	for n := ""; n != str; {
		n = str
		str = re.ReplaceAllString(str, "$1,$2")
	}
	return str
}

func FormatFloatCommas(num float64, digits int) string {
	digitsStr := strconv.Itoa(digits)
	str := fmt.Sprintf("%0."+digitsStr+"f", num)
	re := regexp.MustCompile("(\\d+)(\\d{3})")
	for n := ""; n != str; {
		n = str
		str = re.ReplaceAllString(str, "$1,$2")
	}
	return str
}
