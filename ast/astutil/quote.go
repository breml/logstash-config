package astutil

import (
	"bytes"
	"errors"
	"regexp"

	"github.com/breml/logstash-config/ast"
)

var barewordRe = regexp.MustCompile("(?s:^[A-Za-z_][A-Za-z0-9_]+$)")

// Quote returns a string with quotes for Logstash. Supported quote types
// are ast.DoubleQuoted, ast.SingleQuoted and ast.Bareword.
// If the result is not a valid quoted value, an error is returned.
func Quote(value string, quoteType ast.StringAttributeType) (string, error) {
	switch quoteType {
	case ast.DoubleQuoted, ast.SingleQuoted:
		var hasQuote bool
		quote := quoteType.String()[0]
		for i := 0; i < len(value); i++ {
			if value[i] == quote && i > 1 && value[i-1] != '\\' {
				hasQuote = true
				break
			}
		}

		if hasQuote {
			return "", errors.New("value %q contains unescaped quotes and can not be quoted without escaping")
		}
		return quoteType.String() + value + quoteType.String(), nil

	case ast.Bareword:
		if !barewordRe.MatchString(value) {
			return "", errors.New("value %q contains non bareword characters and can not be quoted as bareword without escaping")
		}
		return value, nil

	default:
		panic("quote type not supported")
	}
}

// QuoteWithEscape returns a string with quotes for Logstash. Supported quote
// types are ast.DoubleQuoted, ast.SingleQuoted and ast.Bareword.
// The value will be escaped if necessary such, that the returned value is a
// valid quoted Logstash string.
// For ast.DoubleQuoted, all double quotes (`"`) are escaped to `\"`.
// For ast.SingleQuoted, all single quotes (`'`) are escaped to `\'`.
// For ast.Bareword, all characters not matching "[A-Za-z_][A-Za-z0-9_]+" are
// replaced with `_`.
func QuoteWithEscape(value string, quoteType ast.StringAttributeType) string {
	switch quoteType {
	case ast.DoubleQuoted, ast.SingleQuoted:
		quote := quoteType.String()[0]
		return quoteType.String() + escapeQuotes(value, quote) + quoteType.String()

	case ast.Bareword:
		return escapeBareword(value)

	default:
		panic("quote type not supported")
	}
}

func escapeQuotes(value string, quote byte) string {
	b := []byte(value)

	for i := 0; i < len(b); i++ {
		if b[i] == quote && (i == 0 || i > 1 && b[i-1] != '\\') {
			b = append(b[:i], append([]byte{'\\'}, b[i:]...)...)
		}
	}

	return string(b)
}

func escapeBareword(value string) string {
	if len(value) == 0 {
		return ""
	}
	b := []byte(value)
	if b[0] >= '0' && b[0] <= '9' {
		b[0] = '_'
	}
	barewordMap := func(r rune) rune {
		switch {
		case r >= '0' && r <= '9':
			return r
		case r >= 'A' && r <= 'Z':
			return r
		case r >= 'a' && r <= 'z':
			return r
		default:
			return '_'
		}
	}
	b = bytes.Map(barewordMap, b)

	return string(b)
}
