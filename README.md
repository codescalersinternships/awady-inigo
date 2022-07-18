# INI Parser
Go package provides read and write for INI files.

## Features
- Load from files and strings.
- Change values and add keys and sections.
- Get parsed data as a map.
- Get sections names as a slice.
- Get parsed data as a string.
- Export parsed data to INI file.

## How To Use
```go
parser := Parser{}
```
you can parse from a file:
```go
parser.LoadFromFile("file.ini")
```
or from a string:
```go
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
parser.LoadFromString(iniText)
```