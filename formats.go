package codeutils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func FormatCommas(num int32) (formatedNum string) {

	formatedNum = FormatFloatCommas(float64(num), 0)
	return
}

func callFormatFloatCommas(num float64, digits int, trim bool) (formatedNum string) {

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
		precesion = precesion[strings.Index(precesion, "."):]
		if !trim || strings.Trim(precesion, "0") != "." {
			formatedNum += precesion
		}
	}
	return
}

func FormatFloatCommasTrim(num float64, digits int) (formatedNum string) {

	formatedNum = callFormatFloatCommas(num, digits, true)
	return
}

func FormatFloatCommas(num float64, digits int) (formatedNum string) {

	formatedNum = callFormatFloatCommas(num, digits, false)
	return
}

func StrToTime(timeStr string) (timeResult time.Time, err error) {

	timeResult, err = time.Parse("2006-01-02 15:04:05", timeStr)
	return
}

func StrToDate(dateStr string) (dateResult time.Time, err error) {

	dateResult, err = time.Parse("2006-01-02", dateStr)
	return
}

func TimeToStr(atime time.Time) (result string) {

	result = atime.Format("2006-01-02 15:04:05")
	return
}

func DateToStr(adate time.Time) (result string) {

	result = adate.Format("2006-01-02")
	return
}

func FormatBytes(bytes int64) string {
	if bytes < 0 {
		return "Invalid size (negative bytes)"
	}
	if bytes == 0 {
		return "0 B"
	}

	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
		PB = 1024 * TB
		EB = 1024 * PB // Exabyte
	)

	var (
		value float64
		unit  string
	)

	switch {
	case bytes < KB:
		value = float64(bytes)
		unit = "B"
	case bytes < MB:
		value = float64(bytes) / KB
		unit = "KB"
	case bytes < GB:
		value = float64(bytes) / MB
		unit = "MB"
	case bytes < TB:
		value = float64(bytes) / GB
		unit = "GB"
	case bytes < PB:
		value = float64(bytes) / TB
		unit = "TB"
	case bytes < EB:
		value = float64(bytes) / PB
		unit = "PB"
	default:
		// For sizes larger than Exabytes, we can continue the pattern or cap it.
		// Here, we'll cap it at EB for simplicity, but you could extend to ZB, YB.
		value = float64(bytes) / EB
		unit = "EB"
	}

	// Format to one decimal place if it's not a whole number, otherwise no decimal.
	if value == math.Trunc(value) {
		return fmt.Sprintf("%.0f %s", value, unit)
	}
	return fmt.Sprintf("%.1f %s", value, unit)
}

func GetLocalTimeAsUTC() time.Time {

	return GetATimeAsUTC(time.Now())
}

func GetATimeAsUTC(atime time.Time) time.Time {

	return time.Date(
		atime.Year(), atime.Month(), atime.Day(),
		atime.Hour(), atime.Minute(), atime.Second(), atime.Nanosecond(),
		time.UTC,
	)
}
