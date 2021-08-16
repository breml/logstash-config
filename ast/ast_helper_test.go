package ast

import "testing"

func TestPrefix(t *testing.T) {
	tt := []struct {
		name         string
		input        string
		emptyNewline bool

		want string
	}{
		{
			name: "empty string",
		},
		{
			name:         `empty string with empty newline`,
			emptyNewline: true,

			want: "\n",
		},
		{
			name:  "simple attribute",
			input: `value => 3.1415`,

			want: `
  value => 3.1415
`,
		},
		{
			name: "simple attribute with newline",
			input: `value => 3.1415
`,

			want: `
  value => 3.1415
`,
		},
		{
			name: "simple attribute with comment",
			input: `// comment
value => 3.1415`,

			want: `
  // comment
  value => 3.1415
`,
		},
		{
			name: "block",
			input: `add_field {
  // comment
  value => 3.1415
}

add_tag => [ "foobar" ]
`,

			want: `
  add_field {
    // comment
    value => 3.1415
  }

  add_tag => [ "foobar" ]
`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := prefix(tc.input, tc.emptyNewline)

			if tc.want != got {
				t.Errorf("want %q, got %q", tc.want, got)
			}
		})
	}
}
