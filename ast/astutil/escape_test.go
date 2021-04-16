package astutil_test

import (
	"testing"

	"github.com/breml/logstash-config/ast/astutil"
)

var tt = []struct {
	name      string
	escaped   string
	unescaped string
}{
	{
		name:      "carriage return (ASCII 13)",
		escaped:   `\r`,
		unescaped: "\r",
	},
	{
		name:      "new line (ASCII 10)",
		escaped:   `\n`,
		unescaped: "\n",
	},
	{
		name:      "tab (ASCII 9)",
		escaped:   `\t`,
		unescaped: "\t",
	},
	{
		name:      "backslash (ASCII 92)",
		escaped:   `\\`,
		unescaped: `\`,
	},
	{
		name:      "double quote (ASCII 34)",
		escaped:   `\"`,
		unescaped: `"`,
	},
	{
		name:      "single quote (ASCII 39)",
		escaped:   `\'`,
		unescaped: `'`,
	},
	{
		name:      "value containing all special characters",
		escaped:   `foo\r\n\t\\\"\'bar`,
		unescaped: "foo\r\n\t\\\"'bar",
	},
}

func TestEscape(t *testing.T) {
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := astutil.Escape(tc.unescaped)
			if tc.escaped != got {
				t.Errorf("Escape(%q) != %q, got: %q", tc.unescaped, tc.escaped, got)
			}
		})
	}
}

func TestUnescape(t *testing.T) {
	tt := append(tt, struct {
		name      string
		escaped   string
		unescaped string
	}{
		name:      "not valid escape",
		escaped:   `\x`,
		unescaped: "\\x",
	})

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := astutil.Unescape(tc.escaped)
			if tc.unescaped != got {
				t.Errorf("Unescape(%q) != %q, got: %q", tc.escaped, tc.unescaped, got)
			}
		})
	}
}
