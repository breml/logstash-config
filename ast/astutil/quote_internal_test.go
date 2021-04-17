package astutil

import "testing"

func TestEscapeQuotes(t *testing.T) {
	tt := []struct {
		name       string
		wantDouble string
		wantSingle string
	}{
		{
			name:       ``,
			wantDouble: ``,
			wantSingle: ``,
		},
		{
			name:       `"`,
			wantDouble: `\"`,
			wantSingle: `"`,
		},
		{
			name:       `"foo"bar"`,
			wantDouble: `\"foo\"bar\"`,
			wantSingle: `"foo"bar"`,
		},
		{
			name:       `'`,
			wantDouble: `'`,
			wantSingle: `\'`,
		},
		{
			name:       `'foo'bar'`,
			wantDouble: `'foo'bar'`,
			wantSingle: `\'foo\'bar\'`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := escapeQuotes(tc.name, '"')
			if tc.wantDouble != got {
				t.Errorf("want: %q, got: %q", tc.wantDouble, got)
			}

			got = escapeQuotes(tc.name, '\'')
			if tc.wantSingle != got {
				t.Errorf("want: %q, got: %q", tc.wantSingle, got)
			}
		})
	}
}

func TestEscapeBareword(t *testing.T) {
	tt := []struct {
		name string
		want string
	}{
		{
			name: "",
			want: "",
		},
		{
			name: "0",
			want: "_",
		},
		{
			name: "bareword",
			want: "bareword",
		},
		{
			name: "BAREWORD",
			want: "BAREWORD",
		},
		{
			name: "_bare_word_",
			want: "_bare_word_",
		},
		{
			name: "0bare1word9",
			want: "_bare1word9",
		},
		{
			name: "-() ",
			want: "____",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := escapeBareword(tc.name)
			if tc.want != got {
				t.Errorf("want: %q, got: %q", tc.want, got)
			}
		})
	}
}
