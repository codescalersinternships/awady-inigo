package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestLineType(t *testing.T) {
	t.Run("section line", func(t *testing.T) {

		line := "[owner]"
		got, err := LineType(line)
		want := "section"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
		if err != nil {
			t.Fatalf("got an error but should get none")
		}
	})

	t.Run("keyValue line", func(t *testing.T) {

		line := "name = John Doe"
		got, err := LineType(line)
		want := "keyValue"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
		if err != nil {
			t.Fatalf("got an error but should get none")
		}
	})

	t.Run("comment line", func(t *testing.T) {

		line := "; last modified 1 April 2001 by John Doe"
		got, _ := LineType(line)
		want := "comment"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	})
	t.Run("empty line", func(t *testing.T) {

		line := "\n"
		got, _ := LineType(line)
		want := "emptyLine"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	})
	t.Run("unknown line", func(t *testing.T) {

		line := "[section]]"
		got, err := LineType(line)
		want := ""

		if err == nil {
			t.Fatalf("should get an Unkown error")
		}

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

func TestParseSection(t *testing.T) {
	t.Run("section line", func(t *testing.T) {
		line := "[section]"
		got, _ := ParseSection(line)
		want := "section"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	})
	t.Run("empty section line", func(t *testing.T) {
		line := "[]"
		_, err := ParseSection(line)

		if err == nil {
			t.Fatalf("should get an error: section can't be empty")
		}

	})
}

func TestParseKeyValue(t *testing.T) {
	t.Run("keyValue Line", func(t *testing.T) {
		line := "name = John Doe"
		gotKey, gotValue, err := ParseKeyValue(line)
		key := "name"
		value := "John Doe"

		if err != nil {
			t.Errorf("got an error but should get none")
		}

		if gotKey != key || gotValue != value {
			t.Errorf("got %v want %v", gotKey, key)
		}

	})
	t.Run("empty key in keyValue Line", func(t *testing.T) {
		line := " = John Doe"
		_, _, err := ParseKeyValue(line)

		if err == nil {
			t.Errorf("should get an error: key can't be empty")
		}

	})

}

func TestParse(t *testing.T) {
	iniText := `; last modified 1 April 2001 by John Doe
[owner]
name = John Doe
organization = Acme Widgets Inc.

[database]
; use IP address in case network name resolution is not working
server = 192.0.2.62     
port = 143
file = "payroll.dat"`
	got, _ := Parse(iniText)
	want := map[string]map[string]string{
		"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
		"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v want\n %#v", got, want)
	}
}

func TestGetSections(t *testing.T) {
	parser := Parser{}
	parser.LoadFromFile("file.ini")
	got := parser.GetSections()
	want := map[string]map[string]string{
		"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
		"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v want %#v", got, want)
	}
}

func TestGetSectionNames(t *testing.T) {
	parser := Parser{map[string]map[string]string{
		"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
		"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
	}}
	got := parser.GetSectionNames()
	want := []string{"owner", "database"}
	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
func TestGet(t *testing.T) {
	parser := Parser{map[string]map[string]string{
		"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
		"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
	}}
	got := parser.Get("owner")
	want := []string{"name", "organization"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
func TestSet(t *testing.T) {
	t.Run("changed \"name\" key in \"owner\" section", func(t *testing.T) {
		parser := Parser{map[string]map[string]string{
			"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
			"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
		}}
		err := parser.Set("owner", "name", "Abdo")

		want := Parser{map[string]map[string]string{
			"owner":    {"name": "Abdo", "organization": "Acme Widgets Inc."},
			"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
		}}
		if err != nil {
			t.Fatalf("should not get an error but got: %q", err.Error())
		}

		if !reflect.DeepEqual(parser, want) {
			t.Errorf("got %v want %v", parser, want)
		}

	})
	t.Run("changed \"name\" key in non-existing section", func(t *testing.T) {
		parser := Parser{map[string]map[string]string{
			"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
			"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
		}}
		err := parser.Set("owne", "name", "Abdo")

		if err == nil {
			t.Fatalf("should get an error section not found but got none")
		}

	})
	t.Run("changed non-existing key in \"owner\" section", func(t *testing.T) {
		parser := Parser{map[string]map[string]string{
			"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
			"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
		}}
		err := parser.Set("owner", "names", "Abdo")

		if err == nil {
			t.Fatalf("should get an error key not found but got none")
		}

	})
}
