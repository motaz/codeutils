// codeutils package
// Configuration file reader .ini file
// By Motaz Abdel Azeem  code.sd
// June 2017
// Updated March 2022

package codeutils

import (
	"os"
	"strings"
)

func load(filename string, dest map[string]string) error {

	fi, err := os.Stat(filename)
	if err != nil {
		return err
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	} else {
		defer f.Close()
	}
	buff := make([]byte, fi.Size())
	f.Read(buff)

	str := string(buff)
	if !strings.HasSuffix(str, "\n") {
		str = str + "\n"
	}
	s2 := strings.Split(str, "\n")

	for i := 0; i < len(s2); {

		if strings.HasPrefix(s2[i], "#") {
			i++
		} else if strings.Contains(s2[i], "=") {
			key := strings.Trim(s2[i][0:strings.Index(s2[i], "=")], " ")
			val := strings.Trim(s2[i][strings.Index(s2[i], "=")+1:len(s2[i])], " ")

			i++
			val = strings.Replace(val, "\r", "", -1)

			dest[key] = val
		} else {
			i++
		}
	}
	return nil
}

func getConfigValue(configFile, name string) (value string, err error) {

	mymap := make(map[string]string)
	value = ""

	configFile = getFullPathConfigFile(configFile)

	err = load(configFile, mymap)
	if err != nil {

		return
	} else if mymap[name] != "" {
		value = mymap[name]
		return
	} else {
		return
	}
}

func GetConfigValue(configFile, name string) (value string) {
	value, _ = getConfigValue(configFile, name)
	return
}

func GetConfigValueWithError(configFile, name string) (value string, err error) {
	value, err = getConfigValue(configFile, name)
	return
}

func SetConfigValue(configFile, name string, value string) (success bool) {

	mymap := make(map[string]string)

	if !strings.Contains(configFile, string(os.PathSeparator)) {
		configFile = GetCurrentAppDir() + string(os.PathSeparator) + configFile
	}
	load(configFile, mymap)

	mymap[name] = value

	file, err := os.OpenFile(configFile, os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0666)

	if err != nil {
		return false
	} else {
		defer file.Close()
	}
	for key := range mymap {

		file.WriteString(key + "=" + mymap[key] + "\n")
	}

	success = true
	return

}

func GetConfigWithDefault(configFile, name, defaultValue string) (value string) {

	value = GetConfigValue(configFile, name)
	if value == "" {
		value = defaultValue
	}
	return
}

func IsFileExists(fileName string) (exists bool) {

	if !strings.Contains(fileName, string(os.PathSeparator)) {
		fileName = GetCurrentAppDir() + string(os.PathSeparator) + fileName
	}

	if _, err := os.Stat(fileName); err == nil {
		exists = true

	} else {
		exists = false
	}
	return
}
