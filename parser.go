package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

var (
	ErrSectionNotFound  = errors.New("section not found")
	ErrKeyNotFound      = errors.New("key not found")
	ErrEmptySectionName = errors.New("section name can't be empty")
	ErrEmptyKey         = errors.New("key can't be empty")
	ErrUnsportedLine    = errors.New("unsupported line")
	ErrGlobalKey        = errors.New("can't parse global keys")
)

//Parser is a struct which parses and holds ini data.
//It can import and export the data and access sections, keys and values from the parsed ini data.
//It can also manipulate the parsed data.
type Parser struct {
	ini map[string]map[string]string
}

// LoadFromString loads and parses ini from iniText.
// It returns any parsing error encountered.
func (p *Parser) LoadFromString(iniText string) (err error) {
	p.ini, err = Parse(iniText)
	return err
}

// LoadFromFile loads and parses ini from iniFile.
// It returns any parsing error encountered.
func (p *Parser) LoadFromFile(iniFile string) (err error) {
	dat, err := os.ReadFile(iniFile)
	if err != nil {
		return err
	}
	iniText := string(dat)
	p.ini, err = Parse(iniText)
	return err
}

// SaveToFile saves ini to outputFile.
// It returns any parsing error encountered.
func (p *Parser) SaveToFile(outputFile string) (err error) {
	file, err := os.Create(outputFile)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.WriteString(p.ToString())
	if err != nil {
		return err
	}
	err = file.Close()
	return err
}

// GetSections gets the ini data as a map in which keys are section names
// and values are maps of keys and values from the ini data.
func (p *Parser) GetSections() map[string]map[string]string {
	return p.ini
}

// GetSectionNames gets the name of sections in the ini data.
// It returns the section names as a slice of strings.
func (p *Parser) GetSectionNames() []string {
	keys := []string{}
	for key := range p.ini {
		keys = append(keys, key)
	}
	return keys
}

// Get gets the value of key in section from the ini data.
// A section and a key in that section must exist in order to get the value or it will return an error
// It returns the value as a string and returns any error encountered
func (p *Parser) Get(section, key string) (string, error) {

	if _, ok := p.ini[section]; !ok {
		return "", ErrSectionNotFound
	}
	if _, ok := p.ini[section][key]; !ok {
		return "", ErrKeyNotFound
	}
	return p.ini[section][key], nil
}

// Set sets the value of a key in section to the ini data.
// If no section exists with input section name a new section will be created and same for keys.
func (p *Parser) Set(section, key, value string) {
	if _, ok := p.ini[section]; !ok {
		p.ini[section] = make(map[string]string)
	}
	p.ini[section][key] = value
}

// ToString converts the ini data to string in the ini format.
// It returns string.
func (p *Parser) ToString() string {
	iniText := ""
	for section, keyValue := range p.ini {
		iniText += "[" + section + "]\n"
		for key, value := range keyValue {
			iniText += key + " = " + value + "\n"
		}
	}
	return iniText
}

// enumeration of line types
const (
	emptyLine int = iota
	sectionLine
	commentLine
	keyValueLine
	unsupportedLine
)

// LineType returns the type of the ini line.
// It could be a section line, key-value line, comment line, empty line or an unsupported line
// It returns int equivilant to which line the input is and any error encountered
func LineType(line string) (int, error) {
	if len(line) == 0 {
		return emptyLine, nil
	}
	if (line[0] == '[') && (strings.Count(line, "]") == 1) &&
		(strings.Count(line, "[") == 1) && (line[len(line)-1] == ']') {
		return sectionLine, nil
	} else if line[0] == ';' {
		return commentLine, nil
	} else if strings.Count(line, "=") > 0 {
		return keyValueLine, nil
	} else if line[0] == '\n' {
		return emptyLine, nil
	} else {
		return unsupportedLine, ErrUnsportedLine
	}

}

// ParseSection parses the section line and returns the name of the section as a string.
// It returns an error if the section name is empty.
func ParseSection(sectionLine string) (string, error) {
	if len(sectionLine) == 2 {
		return "", ErrEmptySectionName
	}
	return sectionLine[1 : len(sectionLine)-1], nil
}

// ParseKeyValue parses the key-value line and returns a key string and a value string.
// It returns an error if the key is empty
func ParseKeyValue(keyValueLine string) (string, string, error) {
	i := strings.Index(keyValueLine, "=")
	key := keyValueLine[0:i]
	key = strings.TrimSpace(key)
	if len(key) == 0 {
		return "", "", ErrEmptyKey
	}
	value := keyValueLine[i+1:]
	value = strings.TrimSpace(value)
	return key, value, nil
}

// Parse parses the iniText and returns it as a map in which keys are section names
// and values are maps of keys and values from the ini data.
// It returns an error if global keys are used.
func Parse(iniText string) (map[string]map[string]string, error) {
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
		case sectionLine:
			currentSection, err = ParseSection(scanner.Text())
			if err != nil {
				return parsedText, err
			}
			parsedText[currentSection] = make(map[string]string)
			continue
		case keyValueLine:
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
			return parsedText, ErrGlobalKey
		}
		parsedText[currentSection][key] = value
	}
	return parsedText, nil
}
