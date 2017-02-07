package ast

import "testing"

func TestAst(t *testing.T) {
	cases := []struct {
		config   Config
		expected string
	}{
		{
			config:   Config{},
			expected: ``,
		},
		{
			config: NewConfig(
				NewPluginSections(
					Input, NewPlugin("stdin",
						NewArrayAttribute(
							"tags", StringAttribute{value: "tag1", sat: DoubleQuoted}, StringAttribute{value: "tag2", sat: SingleQuoted}, StringAttribute{value: "tag3", sat: Bareword},
						),
						NewHashAttribute(
							"add_field",
							NewHashEntry("fieldname", NewStringAttribute("", "fieldvalue", DoubleQuoted)),
							NewHashEntry("number", NewNumberAttribute("", 3.1415)),
						),
						NewNumberAttribute("pi", 3.1415),
						NewPluginAttribute("codec", NewPlugin("rubydebug", NewStringAttribute("string", "a value", DoubleQuoted))),
					),
				),
				nil,
				NewPluginSections(
					Output, NewPlugin(
						"stdout", NewStringAttribute("id", "ABC", DoubleQuoted),
					),
				),
			),
			expected: `input {
  stdin {
    tags => [ "tag1", 'tag2', tag3 ]
    add_field => {
      fieldname => "fieldvalue"
      number => 3.1415
    }
    pi => 3.1415
    codec => rubydebug {
      string => "a value"
    }
  }
}
output {
  stdout {
    id => "ABC"
  }
}
`,
		},
	}

	for _, test := range cases {
		got := test.config.String()
		if got != test.expected {
			t.Errorf("Expected:\n%s\n\nGot:\n%s\n\n", test.expected, got)
		}
	}
}

func TestPluginType(t *testing.T) {
	cases := []struct {
		input    PluginType
		expected string
	}{
		{
			input:    Input,
			expected: "input",
		},
		{
			input:    Filter,
			expected: "filter",
		},
		{
			input:    Output,
			expected: "output",
		},
		{
			input:    0,
			expected: "undefined",
		},
		{
			input:    4,
			expected: "undefined",
		},
	}

	for _, test := range cases {
		if test.input.String() != test.expected {
			t.Errorf("Expected: %s, Got: %s", test.expected, test.input)
		}
	}
}
