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
		want := sectionLine

		assertNoError(t, err)
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	})

	t.Run("keyValue line", func(t *testing.T) {

		line := "name = John Doe"
		got, err := LineType(line)
		want := keyValueLine

		assertNoError(t, err)
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	})

	t.Run("comment line", func(t *testing.T) {

		line := "; last modified 1 April 2001 by John Doe"
		got, err := LineType(line)
		want := commentLine

		assertNoError(t, err)
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	})
	t.Run("empty line", func(t *testing.T) {

		line := "\n"
		got, err := LineType(line)
		want := emptyLine

		assertNoError(t, err)
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	})
	t.Run("unsupported line", func(t *testing.T) {

		line := "[section]]"
		_, err := LineType(line)
		want := ErrUnsportedLine

		assertError(t, err, want)
	})

	t.Run("comment at the end of the line", func(t *testing.T) {

		line := "[section] ;testing"
		_, err := LineType(line)
		want := ErrUnsportedLine

		assertError(t, err, want)
	})
	t.Run("line with more than one equal sign", func(t *testing.T) {

		line := "key = val=ue"
		got, err := LineType(line)
		want := keyValueLine

		assertNoError(t, err)
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

func TestParseSection(t *testing.T) {
	t.Run("section line", func(t *testing.T) {
		line := "[section]"
		got, err := ParseSection(line)
		want := "section"

		assertNoError(t, err)
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	})
	t.Run("empty section line", func(t *testing.T) {
		line := "[]"
		_, err := ParseSection(line)
		want := ErrEmptySectionName
		assertError(t, err, want)

	})
}

func TestParseKeyValue(t *testing.T) {
	t.Run("keyValue Line", func(t *testing.T) {
		line := "name = John Doe"
		gotKey, gotValue, err := ParseKeyValue(line)
		key := "name"
		value := "John Doe"

		assertNoError(t, err)
		if gotKey != key || gotValue != value {
			t.Errorf("got %v want %v", gotKey, key)
		}

	})
	t.Run("empty key in keyValue Line", func(t *testing.T) {
		line := " = John Doe"
		_, _, err := ParseKeyValue(line)
		want := ErrEmptyKey

		assertError(t, err, want)

	})

	t.Run("empty value in keyValue Line", func(t *testing.T) {
		line := "name = "
		gotKey, gotValue, err := ParseKeyValue(line)
		key := "name"
		value := ""

		assertNoError(t, err)
		if gotKey != key || gotValue != value {
			t.Errorf("got %v want %v", gotKey, key)
		}

	})

}

func TestParse(t *testing.T) {
	t.Run("ini text", func(t *testing.T) {
		iniText := `; last modified 1 April 2001 by John Doe
[owner]
name = John Doe
organization = Acme Widgets Inc.

[database]
; use IP address in case network name resolution is not working
server = 192.0.2.62     
port = 143
file = "payroll.dat"`
		got, err := Parse(iniText)
		want := map[string]map[string]string{
			"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
			"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
		}

		assertNoError(t, err)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %#v want\n %#v", got, want)
		}

	})

	t.Run("ini text with global keys", func(t *testing.T) {
		iniText := `; last modified 1 April 2001 by John Doe
name = Test
[owner]
name = John Doe
organization = Acme Widgets Inc.

[database]
; use IP address in case network name resolution is not working
server = 192.0.2.62     
port = 143
file = "payroll.dat"`
		_, err := Parse(iniText)
		want := ErrGlobalKey

		assertError(t, err, want)

	})
}

func TestLoadFromString(t *testing.T) {
	t.Run("passing ini string", func(t *testing.T) {
		iniText := `; last modified 1 April 2001 by John Doe
[owner]
name = John Doe
organization = Acme Widgets Inc.

[database]
; use IP address in case network name resolution is not working
server = 192.0.2.62     
port = 143
file = "payroll.dat"`
		parser := Parser{}
		err := parser.LoadFromString(iniText)
		assertNoError(t, err)

		want := Parser{map[string]map[string]string{
			"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
			"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
		}}

		if !reflect.DeepEqual(parser, want) {
			t.Errorf("got %v want %v", parser, want)
		}

	})
	t.Run("passing ini string with unsupported global keys", func(t *testing.T) {
		iniText := `; last modified 1 April 2001 by John Doe
name = test
[owner]
name = John Doe
organization = Acme Widgets Inc.

[database]
; use IP address in case network name resolution is not working
server = 192.0.2.62     
port = 143
file = "payroll.dat"`
		parser := Parser{}
		err := parser.LoadFromString(iniText)
		want := ErrGlobalKey
		assertError(t, err, want)

	})
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
	t.Run("getting value using key and section", func(t *testing.T) {
		parser := Parser{map[string]map[string]string{
			"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
			"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
		}}
		got, err := parser.Get("owner", "name")
		want := "John Doe"

		assertNoError(t, err)
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	})
	t.Run("getting value using non-existing section", func(t *testing.T) {
		parser := Parser{map[string]map[string]string{
			"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
			"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
		}}
		_, err := parser.Get("test", "name")
		want := ErrSectionNotFound

		assertError(t, err, want)
	})
	t.Run("getting value using non-existing key", func(t *testing.T) {
		parser := Parser{map[string]map[string]string{
			"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
			"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
		}}
		_, err := parser.Get("owner", "data")
		want := ErrKeyNotFound

		assertError(t, err, want)
	})
}
func TestSet(t *testing.T) {
	t.Run("changing value using key in section", func(t *testing.T) {
		parser := Parser{map[string]map[string]string{
			"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
			"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
		}}
		parser.Set("owner", "name", "Abdo")

		want := Parser{map[string]map[string]string{
			"owner":    {"name": "Abdo", "organization": "Acme Widgets Inc."},
			"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
		}}

		if !reflect.DeepEqual(parser, want) {
			t.Errorf("got %v want %v", parser, want)
		}

	})
	t.Run("changing value in non-existing section", func(t *testing.T) {
		parser := Parser{map[string]map[string]string{
			"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
			"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
		}}
		parser.Set("owne", "name", "Abdo")
		want := Parser{map[string]map[string]string{
			"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
			"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
			"owne":     {"name": "Abdo"},
		}}

		if !reflect.DeepEqual(parser, want) {
			t.Errorf("got %v want %v", parser, want)
		}

	})
	t.Run("changing value using non-existing key", func(t *testing.T) {
		parser := Parser{map[string]map[string]string{
			"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
			"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
		}}
		parser.Set("owner", "names", "Abdo")
		want := Parser{map[string]map[string]string{
			"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc.", "names": "Abdo"},
			"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
		}}

		if !reflect.DeepEqual(parser, want) {
			t.Errorf("got %v want %v", parser, want)
		}

	})
}

func TestToString(t *testing.T) {
	parser := Parser{}
	parser.LoadFromFile("file.ini")

	got, err := Parse(parser.ToString())
	want := map[string]map[string]string{
		"owner":    {"name": "John Doe", "organization": "Acme Widgets Inc."},
		"database": {"server": "192.0.2.62", "port": "143", "file": "\"payroll.dat\""},
	}

	assertNoError(t, err)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("got an error:%s but didn't want one", got)
	}
}

func assertError(t testing.TB, got error, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
