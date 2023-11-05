package codeutils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func FormatCommas(num int32) (formatedNum string) {

	formatedNum = FormatFloatCommas(float64(num), 0)
	return
}

func FormatFloatCommas(num float64, digits int) (formatedNum string) {

	digitsStr := strconv.Itoa(digits)

	formatedNum = fmt.Sprintf("%0."+digitsStr+"f", num)
	re := regexp.MustCompile("(\\d+)(\\d{3})")
	for n := ""; n != formatedNum; {
		n = formatedNum
		formatedNum = re.ReplaceAllString(formatedNum, "$1,$2")
	}
	if digits > 0 {
		if strings.Contains(formatedNum, ".") {
			formatedNum = formatedNum[:strings.Index(formatedNum, ".")]
		}

		precesion := fmt.Sprintf("%0."+digitsStr+"f", num)
		formatedNum += precesion[strings.Index(precesion, "."):]
	}
	return
}

func StrToTime(timeStr string) (timeResult time.Time, err error) {

	timeResult, err = time.Parse(time.DateTime, timeStr)
	return
}

func StrToDate(dateStr string) (dateResult time.Time, err error) {

	dateResult, err = time.Parse(time.DateOnly, dateStr)
	return
}

func TimeToStr(atime time.Time) (result string) {

	result = atime.Format(time.DateTime)
	return
}

func DateToStr(adate time.Time) (result string) {

	result = adate.Format(time.DateOnly)
	return
}
