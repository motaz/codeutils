// WriteToLog function
// Updated November 2022

package codeutils

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const MONTHLYLOG = 0
const WEEKLYLOG = 1
const HOURLYLOG = 2

var logdaytype byte = MONTHLYLOG

func SetLogType(alogtype byte) {
	logdaytype = alogtype
}

// GetCurrentAppDir returns path from application running directory
func GetCurrentAppDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	if strings.HasPrefix(dir, "/tmp/") || strings.Contains(dir, "go-build") {
		dir, _ = os.Getwd()

	}
	return dir
}

// WriteToLog write to log file
func WriteToLog(event string, logFileName string) (err error) {

	t := time.Now()
	var logID string
	if logdaytype == MONTHLYLOG {
		logID = strconv.Itoa(t.Day())
	} else if logdaytype == WEEKLYLOG {
		logID = t.Format("Mon")
	} else if logdaytype == HOURLYLOG {
		logID = strconv.Itoa(t.Hour())
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
	_, err = os.Stat(dir)
	if (err != nil) && (os.IsNotExist(err)) {
		err = os.Mkdir(dir, 0777)
	}

	// Check current log date, if it is old, overwrite it
	logname = logFileName + "-" + logID + ".log"
	var logstat os.FileInfo
	logstat, err = os.Stat(logname)
	if err == nil {
		if t.Month() != logstat.ModTime().Month() || t.Day() != logstat.ModTime().Day() {
			old = true
		}
		if logdaytype == HOURLYLOG {
			if t.Hour() != logstat.ModTime().Hour() {
				old = true
			}
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

		_, err = f.WriteString(t.String()[1:23] + ": " + event + "\n")
		if err != nil {
			println("Error in writing log: ", err.Error())
		}
	}
	return

}
