package main

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

type ParsedText map[string]map[string]string

func (p ParsedText) GetSectionNames() []string {
	keys := []string{}
	for key, _ := range p {
		keys = append(keys, key)
	}
	return keys
}
func (p ParsedText) Get(section string) []string {

	keys := []string{}
	for key, _ := range p[section] {
		keys = append(keys, key)
	}

	return keys
}

func LineType(line string) (string, error) {
	if len(line) == 0 {
		return "emptyLine", nil
	}
	if (line[0] == '[') && (strings.Count(line, "]") == 1) &&
		(strings.Count(line, "[") == 1) && (line[len(line)-1] == ']') {
		return "section", nil
	} else if strings.Count(line, "=") == 1 {
		return "keyValue", nil
	} else if line[0] == ';' {
		return "comment", nil
	} else if line[0] == '\n' {
		return "emptyLine", nil
	} else {
		return "", errors.New("Unkown Line")
	}

}

func ParseSection(sectionLine string) (string, error) {
	if len(sectionLine) == 2 {
		return "", errors.New("section can't be empty")
	}
	return sectionLine[1 : len(sectionLine)-1], nil
}

func ParseKeyValue(keyValueLine string) (string, string, error) {
	i := strings.Index(keyValueLine, "=")
	key := keyValueLine[0:i]
	key = strings.TrimSpace(key)
	if len(key) == 0 {
		return "", "", errors.New("key can't be empty")
	}
	value := keyValueLine[i+1:]
	value = strings.TrimSpace(value)
	return key, value, nil
}

func Parse(iniText string) (ParsedText, error) {
	scanner := bufio.NewScanner(strings.NewReader(iniText))
	parsedText := make(map[string]map[string]string)
	currentSection := ""
	for scanner.Scan() {
		lineType, err := LineType(scanner.Text())
		key := ""
		value := ""
		if err != nil {
			return parsedText, err
		}
		switch lineType {
		case "section":
			currentSection, err = ParseSection(scanner.Text())
			if err != nil {
				return parsedText, err
			}
			parsedText[currentSection] = make(map[string]string)
			continue
		case "keyValue":
			key, value, err = ParseKeyValue(scanner.Text())
			if err != nil {
				return parsedText, err
			}
		default:
			continue
		}
		if err != nil {
			return parsedText, err
		}
		if currentSection == "" {
			return parsedText, errors.New("keys have to be under a section")
		}
		parsedText[currentSection][key] = value
	}
	return parsedText, nil
}
func main() {
	iniText := `; last modified 1 April 2001 by John Doe
[owner]
name = John Doe
organization = Acme Widgets Inc.

[database]
; use IP address in case network name resolution is not working
server = 192.0.2.62     
port = 143
file = "payroll.dat"`
	parsedText, _ := Parse(iniText)
	fmt.Println(parsedText.GetSectionNames())
	fmt.Println(parsedText.Get("database"))
}
