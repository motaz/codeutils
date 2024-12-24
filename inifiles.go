package codeutils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
)

func lineIsNotCommented(line string) (notCommented bool) {

	notCommented = !strings.HasPrefix(line, ";") && !strings.HasPrefix(line, "#")
	return
}

func lineIsSection(line string) (isSection bool) {

	isSection = strings.HasPrefix(line, "[") && strings.Contains(line, "]")
	return
}

func getSectionFromLine(line string) (section string) {

	section = line[1:]
	section = section[:strings.Index(section, "]")]
	section = strings.TrimSpace(section)
	section = strings.ToLower(section)
	return
}

func ReadINIAsInt(filename, section, key string) (value int, err error) {

	var valueStr string
	valueStr, err = ReadINIValue(filename, section, key, "")
	if err == nil {
		if valueStr == "" {
			err = errors.New("Not found")
		} else {
			value, err = strconv.Atoi(valueStr)
		}
	}
	return
}

func ReadINIAsBool(filename, section, key string) (value bool, err error) {

	var valueStr string
	valueStr, err = ReadINIValue(filename, section, key, "")
	if err == nil {
		valueStr = strings.TrimSpace(strings.ToLower(valueStr))
		if valueStr == "" {
			err = errors.New("Not found")
		} else if valueStr == "yes" {
			value = true
		} else if valueStr == "no" {
			value = false
		} else {
			value, err = strconv.ParseBool(valueStr)
		}
	}
	return
}

func ReadINIValue(filename, section, key, defaultValue string) (value string, err error) {

	section = strings.ToLower(section)
	section = strings.TrimSpace(section)
	value = defaultValue
	key = strings.ToLower(key)
	var lines []string
	lines, err = readFileAsLines(filename)
	if err == nil && lines != nil {
		foundSection := ""
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if lineIsNotCommented(line) {
				if lineIsSection(line) {
					foundSection = getSectionFromLine(line)
				}
				if foundSection == section && strings.Index(line, "=") > 1 {
					akey := line[:strings.Index(line, "=")]
					if key == strings.ToLower(strings.TrimSpace(akey)) {
						value = line[strings.Index(line, "=")+1:]
						if strings.Index(value, " ;") > 2 {
							value = value[:strings.Index(value, " ;")-1]
						}
						value = strings.TrimSpace(value)

						break
					}
				}
			}

		}
	}
	return
}

func readFileAsLines(filename string) (lines []string, err error) {

	filename = getFullPathConfigFile(filename)
	var info fs.FileInfo
	info, err = os.Stat(filename)
	if err == nil {
		var file *os.File
		file, err = os.Open(filename)
		if err == nil {

			defer file.Close()

			buff := make([]byte, info.Size())
			file.Read(buff)

			str := string(buff)
			if !strings.HasSuffix(str, "\n") {
				str = str + "\n"
			}

			lines = strings.Split(str, "\n")
		}
	}
	return
}

func ReadINISections(filename string) (sections []string, err error) {

	sections = make([]string, 0)
	var lines []string
	lines, err = readFileAsLines(filename)
	if err == nil && lines != nil {
		foundSection := ""
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if lineIsNotCommented(line) {
				if lineIsSection(line) {
					foundSection = getSectionFromLine(line)
					sections = append(sections, foundSection)
				}

			}

		}
	}
	return
}

func ReadINISectionsKeys(filename, section string) (keys []string, err error) {

	keys = make([]string, 0)
	var lines []string
	lines, err = readFileAsLines(filename)
	if err == nil && lines != nil {
		foundSection := ""
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if lineIsNotCommented(line) && strings.TrimSpace(line) != "" {
				if lineIsSection(line) {
					foundSection = getSectionFromLine(line)
				} else {
					if foundSection == section {
						akey := strings.TrimSpace(line)
						if strings.Contains(akey, "=") {
							akey = akey[:strings.Index(akey, "=")]

							akey = strings.TrimSpace(akey)
						}
						keys = append(keys, akey)

					}
				}

			}

		}
	}
	return
}

func SetIniValue(filename, section, key, value string) error {
	// Open the INI file
	if !IsFileExists(filename) {
		os.Create(filename)
	}
	// Open the INI file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read the contents of the file
	var buffer bytes.Buffer
	scanner := bufio.NewScanner(file)

	sectionFound := false
	keyUpdated := false
	insideTargetSection := false

	for scanner.Scan() {
		line := scanner.Text()

		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "[") && strings.HasSuffix(trimmedLine, "]") {
			// If we are leaving the target section and haven't added the key-value pair, do so now
			if insideTargetSection && !keyUpdated {
				buffer.WriteString(fmt.Sprintf("%s=%s\n", key, value))
				keyUpdated = true
			}

			// Check if this is the target section
			sectionName := strings.TrimSpace(trimmedLine[1 : len(trimmedLine)-1])
			insideTargetSection = sectionName == section
			if insideTargetSection {
				sectionFound = true
			}
		}

		// If inside the target section, update the key-value pair
		if insideTargetSection && strings.HasPrefix(trimmedLine, key+"=") {
			buffer.WriteString(fmt.Sprintf("%s=%s\n", key, value))
			keyUpdated = true
			continue
		}

		// Write the line to the buffer
		buffer.WriteString(line + "\n")
	}

	// If the target section was not found, add it at the end of the file
	if !sectionFound {
		buffer.WriteString(fmt.Sprintf("\n[%s]\n", section))
		buffer.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	} else if insideTargetSection && !keyUpdated {
		// If the section exists but the key was not updated, add the key-value pair
		buffer.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	}

	// Check for scan errors
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Write the updated content back to the file
	err = os.WriteFile(filename, buffer.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
