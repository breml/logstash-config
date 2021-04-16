package astutil

import (
	"regexp"
)

var unescapeRe = regexp.MustCompile(`\\.`)

// Unescape converts an input containing escape sequences as understood by
// Logstash (https://www.elastic.co/guide/en/logstash/current/configuration-file-structure.html#_escape_sequences)
// to an unescaped string.
//
// Unescaped character sequences:
//
//     \r : carriage return (ASCII 13)
//     \n : new line (ASCII 10)
//     \t : tab (ASCII 9)
//     \\ : backslash (ASCII 92)
//     \" : double quote (ASCII 34)
//     \' : single quote (ASCII 39)
//
// Based on:
// https://github.com/elastic/logstash/blob/e9c9865f4066b54048f8d708612a72d25e2fe5d9/logstash-core/lib/logstash/config/string_escape.rb
func Unescape(value string) string {
	return unescapeRe.ReplaceAllStringFunc(value, func(s string) string {
		switch s[1] {
		case '"', '\'', '\\':
			return string(s[1])
		case 'n':
			return "\n"
		case 'r':
			return "\r"
		case 't':
			return "\t"
		default:
			return s
		}
	})
}

var escapeRe = regexp.MustCompile(`(?s:.)`)

// Escape converts an input to an escaped representation as understood by
// Logstash (https://www.elastic.co/guide/en/logstash/current/configuration-file-structure.html#_escape_sequences).
//
// Escaped characters:
//
//     \r : carriage return (ASCII 13)
//     \n : new line (ASCII 10)
//     \t : tab (ASCII 9)
//     \  : backslash (ASCII 92)
//     "  : double quote (ASCII 34)
//     '  : single quote (ASCII 39)
func Escape(value string) string {
	return escapeRe.ReplaceAllStringFunc(value, func(s string) string {
		switch s {
		case `"`, `'`, `\`:
			return `\` + s
		case "\n":
			return `\n`
		case "\r":
			return `\r`
		case "\t":
			return `\t`
		default:
			return s
		}
	})
}
