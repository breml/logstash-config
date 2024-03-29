package astutil_test

import (
	"testing"

	"github.com/breml/logstash-config/ast"
	"github.com/breml/logstash-config/ast/astutil"
)

func TestQuote(t *testing.T) {
	tt := []struct {
		name string
		in   string

		want    []string
		wantErr []bool
	}{
		{
			name: "bareword",
			in:   `bareword`,

			want: []string{
				ast.DoubleQuoted: `"bareword"`,
				ast.SingleQuoted: `'bareword'`,
				ast.Bareword:     `bareword`,
			},
			wantErr: []bool{
				ast.DoubleQuoted: false,
				ast.SingleQuoted: false,
				ast.Bareword:     false,
			},
		},
		{
			name: "multiple words",
			in:   `multiple words`,

			want: []string{
				ast.DoubleQuoted: `"multiple words"`,
				ast.SingleQuoted: `'multiple words'`,
				ast.Bareword:     ``,
			},
			wantErr: []bool{
				ast.DoubleQuoted: false,
				ast.SingleQuoted: false,
				ast.Bareword:     true,
			},
		},
		{
			name: "double quote",
			in:   `value with " (double quote)`,

			want: []string{
				ast.DoubleQuoted: ``,
				ast.SingleQuoted: `'value with " (double quote)'`,
				ast.Bareword:     ``,
			},
			wantErr: []bool{
				ast.DoubleQuoted: true,
				ast.SingleQuoted: false,
				ast.Bareword:     true,
			},
		},
		{
			name: "escaped double quote",
			in:   `value with \" (escaped double quote)`,

			want: []string{
				ast.DoubleQuoted: `"value with \" (escaped double quote)"`,
				ast.SingleQuoted: `'value with \" (escaped double quote)'`,
				ast.Bareword:     ``,
			},
			wantErr: []bool{
				ast.DoubleQuoted: false,
				ast.SingleQuoted: false,
				ast.Bareword:     true,
			},
		},
		{
			name: "single quote",
			in:   `value with ' (single quote)`,

			want: []string{
				ast.DoubleQuoted: `"value with ' (single quote)"`,
				ast.SingleQuoted: ``,
				ast.Bareword:     ``,
			},
			wantErr: []bool{
				ast.DoubleQuoted: false,
				ast.SingleQuoted: true,
				ast.Bareword:     true,
			},
		},
		{
			name: "escaped single quote",
			in:   `value with \' (escaped single quote)`,

			want: []string{
				ast.DoubleQuoted: `"value with \' (escaped single quote)"`,
				ast.SingleQuoted: `'value with \' (escaped single quote)'`,
				ast.Bareword:     ``,
			},
			wantErr: []bool{
				ast.DoubleQuoted: false,
				ast.SingleQuoted: false,
				ast.Bareword:     true,
			},
		},
	}

	quoteTypes := map[string]ast.StringAttributeType{
		"double quote": ast.DoubleQuoted,
		"single quote": ast.SingleQuoted,
		"bareword":     ast.Bareword,
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.want) != 4 && len(tc.wantErr) != 4 {
				t.Error("test case has an invalid number of want or wantErr values")
			}
			for name, quoteType := range quoteTypes {
				t.Run(name, func(t *testing.T) {
					got, err := astutil.Quote(tc.in, quoteType)
					if tc.wantErr[quoteType] != (err != nil) {
						t.Errorf("wantErr %t, err: %v", tc.wantErr[quoteType], err)
					}
					if tc.want[quoteType] != got {
						t.Errorf("want: %q, got: %q", tc.want[quoteType], got)
					}
				})
			}
		})
	}
}

func TestQuoteWithEscape(t *testing.T) {
	tt := []struct {
		name string
		in   string

		wantEscaped []string
	}{
		{
			name: "bareword",
			in:   `bareword`,

			wantEscaped: []string{
				ast.DoubleQuoted: `"bareword"`,
				ast.SingleQuoted: `'bareword'`,
				ast.Bareword:     `bareword`,
			},
		},
		{
			name: "multiple words",
			in:   `multiple words`,

			wantEscaped: []string{
				ast.DoubleQuoted: `"multiple words"`,
				ast.SingleQuoted: `'multiple words'`,
				ast.Bareword:     `multiple_words`,
			},
		},
		{
			name: "double quote",
			in:   `value with " (double quote)`,

			wantEscaped: []string{
				ast.DoubleQuoted: `"value with \" (double quote)"`,
				ast.SingleQuoted: `'value with " (double quote)'`,
				ast.Bareword:     `value_with____double_quote_`,
			},
		},
		{
			name: "escaped double quote",
			in:   `value with \" (escaped double quote)`,

			wantEscaped: []string{
				ast.DoubleQuoted: `"value with \" (escaped double quote)"`,
				ast.SingleQuoted: `'value with \" (escaped double quote)'`,
				ast.Bareword:     `value_with_____escaped_double_quote_`,
			},
		},
		{
			name: "single quote",
			in:   `value with ' (single quote)`,

			wantEscaped: []string{
				ast.DoubleQuoted: `"value with ' (single quote)"`,
				ast.SingleQuoted: `'value with \' (single quote)'`,
				ast.Bareword:     `value_with____single_quote_`,
			},
		},
		{
			name: "escaped single quote",
			in:   `value with \' (escaped single quote)`,

			wantEscaped: []string{
				ast.DoubleQuoted: `"value with \' (escaped single quote)"`,
				ast.SingleQuoted: `'value with \' (escaped single quote)'`,
				ast.Bareword:     `value_with_____escaped_single_quote_`,
			},
		},
	}

	quoteTypes := map[string]ast.StringAttributeType{
		"double quote": ast.DoubleQuoted,
		"single quote": ast.SingleQuoted,
		"bareword":     ast.Bareword,
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.wantEscaped) != 4 {
				t.Error("test case has an invalid number of want values")
			}
			for name, quoteType := range quoteTypes {
				t.Run(name, func(t *testing.T) {
					gotEscaped := astutil.QuoteWithEscape(tc.in, quoteType)
					if tc.wantEscaped[quoteType] != gotEscaped {
						t.Errorf("want: %q, got: %q", tc.wantEscaped[quoteType], gotEscaped)
					}
				})
			}
		})
	}
}
