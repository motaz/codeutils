package codeutils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func FormatCommas(num int32) (formatedNum string) {

	formatedNum = FormatFloatCommas(float64(num), 0)
	return
}

func FormatFloatCommas(num float64, digits int) (formatedNum string) {

	formatedNum = fmt.Sprintf("%0.0f", num)
	re := regexp.MustCompile("(\\d+)(\\d{3})")
	for n := ""; n != formatedNum; {
		n = formatedNum
		formatedNum = re.ReplaceAllString(formatedNum, "$1,$2")
	}
	if digits > 0 {
		digitsStr := strconv.Itoa(digits)

		precesion := fmt.Sprintf("%0."+digitsStr+"f", num)
		formatedNum += precesion[strings.Index(precesion, "."):]
	}
	return
}
