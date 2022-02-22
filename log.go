// Write to log function
// Updated Feb 2022

package codeutils

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const WEEKLYLOG = 1
const MONTHLYLOG = 0

var logdaytype byte = MONTHLYLOG

func SetLogType(alogtype byte) {
	logdaytype = alogtype
}

// GetCurrentAppDir returns path from application running directory
func GetCurrentAppDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	if strings.HasPrefix(dir, "/tmp/") {
		dir, _ = os.Getwd()

	}
	return dir
}

// WriteToLog write to log file
func WriteToLog(event string, logFileName string) {

	t := time.Now()
	var today string
	if logdaytype == MONTHLYLOG {
		today = strconv.Itoa(t.Day())
	} else if logdaytype == WEEKLYLOG {
		today = t.Format("Mon")
	}
	old := false
	var dir string
	var logname string

	if strings.Contains(logFileName, string(os.PathSeparator)) {
		dir = path.Dir(logFileName)
	} else {
		dir = GetCurrentAppDir() + string(os.PathSeparator) + "log"
		logFileName = dir + string(os.PathSeparator) + logFileName
	}
	_, err := os.Stat(dir)
	if (err != nil) && (os.IsNotExist(err)) {
		os.Mkdir(dir, 0777)
	}

	// Check current log date, if it is old, overwrite it
	logname = logFileName + "-" + today + ".log"
	logstat, err := os.Stat(logname)
	if err == nil {
		if t.Month() != logstat.ModTime().Month() || t.Day() != logstat.ModTime().Day() {
			old = true
		}
	}
	var f *os.File

	if old {
		os.Remove(logname)
		f, err = os.OpenFile(logname, os.O_CREATE|os.O_RDWR, 0666)

	} else {
		f, err = os.OpenFile(logname, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	}
	if err == nil {
		defer f.Close()

		_, er := f.WriteString(t.String()[1:23] + ": " + event + "\n")
		if er != nil {
			println("Error in writing log: ", er.Error())
		}
	}

}
