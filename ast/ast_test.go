package ast_test

import (
	"testing"

	"github.com/breml/logstash-config/ast"
	. "github.com/breml/logstash-config/ast"
)

func TestAst(t *testing.T) {
	cases := []struct {
		name     string
		config   Config
		expected string
	}{
		{
			name:     "empty config",
			config:   Config{},
			expected: ``,
		},
		{
			name: "Plugins with attributes of various types",
			config: NewConfig(
				NewPluginSections(
					Input, NewPlugin("stdin",
						NewArrayAttribute(
							"tags", NewStringAttribute("", "tag1", DoubleQuoted), NewStringAttribute("", "tag2", SingleQuoted), NewStringAttribute("", "tag3", Bareword),
						),
						NewHashAttribute(
							"add_field",
							NewHashEntry(NewStringAttribute("", "bareword", Bareword), NewStringAttribute("", "bareword", Bareword)),
							NewHashEntry(NewStringAttribute("", "single quoted", SingleQuoted), NewStringAttribute("", "single quoted", SingleQuoted)),
							NewHashEntry(NewStringAttribute("", "double quoted", DoubleQuoted), NewStringAttribute("", "double quoted", DoubleQuoted)),
							NewHashEntry(NewNumberAttribute("", 1), NewNumberAttribute("", 3.1415)),
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
    tags => [
      "tag1",
      'tag2',
      tag3
    ]
    add_field => {
      bareword => bareword
      'single quoted' => 'single quoted'
      "double quoted" => "double quoted"
      1 => 3.1415
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
		{
			name: "Simple if (without else) branch",
			config: NewConfig(
				nil,
				NewPluginSections(
					Filter, NewBranch(
						NewIfBlock(
							NewCondition(
								NewCompareExpression(
									ast.BooleanOperator{Op: NoOperator}, NewStringAttribute("", "true", DoubleQuoted), CompareOperator{Op: Equal}, NewStringAttribute("", "true", DoubleQuoted),
								),
							),
							NewPlugin("if-plugin"),
						),
						NewElseBlock(),
					),
				),
				nil,
			),
			expected: `filter {
  if "true" == "true" {
    if-plugin {}
  }
}
`,
		},
		{
			name: "Simple if-else branch",
			config: NewConfig(
				nil,
				NewPluginSections(
					Filter, NewBranch(
						NewIfBlock(
							NewCondition(
								NewCompareExpression(
									ast.BooleanOperator{Op: NoOperator}, NewStringAttribute("", "true", DoubleQuoted), CompareOperator{Op: Equal}, NewStringAttribute("", "true", DoubleQuoted),
								),
							),
							NewPlugin("if-plugin"),
						),
						NewElseBlock(
							NewPlugin("else-plugin"),
						),
					),
				),
				nil,
			),
			expected: `filter {
  if "true" == "true" {
    if-plugin {}
  } else {
    else-plugin {}
  }
}
`,
		},
		{
			name: "if with multiple else-if and a final else branch, test for different condition types",
			config: NewConfig(
				nil,
				NewPluginSections(
					Filter, NewBranch(
						NewIfBlock(
							NewCondition(
								NewCompareExpression(
									ast.BooleanOperator{Op: NoOperator}, NewStringAttribute("", "true", DoubleQuoted), CompareOperator{Op: Equal}, NewStringAttribute("", "true", DoubleQuoted),
								),
							),
							NewPlugin("if-plugin"),
						),
						NewElseBlock(
							NewPlugin("else-plugin"),
						),
						NewElseIfBlock(
							NewCondition(
								NewCompareExpression(
									ast.BooleanOperator{Op: NoOperator}, NewStringAttribute("", "true", DoubleQuoted), CompareOperator{Op: NotEqual}, NewStringAttribute("", "true", DoubleQuoted),
								),
							),
							NewPlugin("else-if-plugin-1"),
						),
						NewElseIfBlock(
							NewCondition(
								NewCompareExpression(
									ast.BooleanOperator{Op: NoOperator}, NewNumberAttribute("", 10), CompareOperator{Op: GreaterThan}, NewNumberAttribute("", 2),
								),
							),
							NewPlugin("else-if-plugin-2"),
						),
						NewElseIfBlock(
							NewCondition(
								NewCompareExpression(
									ast.BooleanOperator{Op: NoOperator}, NewNumberAttribute("", 2), CompareOperator{Op: LessThan}, NewNumberAttribute("", 10),
								),
							),
							NewPlugin("else-if-plugin-3"),
						),
						NewElseIfBlock(
							NewCondition(
								NewCompareExpression(
									ast.BooleanOperator{Op: NoOperator}, NewNumberAttribute("", 10), CompareOperator{Op: GreaterOrEqual}, NewNumberAttribute("", 2),
								),
							),
							NewPlugin("else-if-plugin-4"),
						),
						NewElseIfBlock(
							NewCondition(
								NewCompareExpression(
									ast.BooleanOperator{Op: NoOperator}, NewNumberAttribute("", 2), CompareOperator{Op: LessOrEqual}, NewNumberAttribute("", 10),
								),
							),
							NewPlugin("else-if-plugin-5"),
						),
					),
				),
				nil,
			),
			expected: `filter {
  if "true" == "true" {
    if-plugin {}
  } else if "true" != "true" {
    else-if-plugin-1 {}
  } else if 10 > 2 {
    else-if-plugin-2 {}
  } else if 2 < 10 {
    else-if-plugin-3 {}
  } else if 10 >= 2 {
    else-if-plugin-4 {}
  } else if 2 <= 10 {
    else-if-plugin-5 {}
  } else {
    else-plugin {}
  }
}
`,
		},
		{
			name: "if with multiple compare operators",
			config: NewConfig(
				nil,
				NewPluginSections(
					Filter, NewBranch(
						NewIfBlock(
							NewCondition(
								NewCompareExpression(
									ast.BooleanOperator{Op: NoOperator}, NewStringAttribute("", "true", DoubleQuoted), CompareOperator{Op: Equal}, NewStringAttribute("", "true", DoubleQuoted),
								),
								NewCompareExpression(
									ast.BooleanOperator{Op: And}, NewStringAttribute("", "true", DoubleQuoted), CompareOperator{Op: Equal}, NewStringAttribute("", "true", DoubleQuoted),
								),
								NewCompareExpression(
									ast.BooleanOperator{Op: Or}, NewStringAttribute("", "true", DoubleQuoted), CompareOperator{Op: Equal}, NewStringAttribute("", "true", DoubleQuoted),
								),
								NewCompareExpression(
									ast.BooleanOperator{Op: Nand}, NewStringAttribute("", "true", DoubleQuoted), CompareOperator{Op: Equal}, NewStringAttribute("", "true", DoubleQuoted),
								),
								NewCompareExpression(
									ast.BooleanOperator{Op: Xor}, NewStringAttribute("", "true", DoubleQuoted), CompareOperator{Op: Equal}, NewStringAttribute("", "true", DoubleQuoted),
								),
							),
							NewPlugin("plugin"),
						),
						NewElseBlock(),
					),
				),
				nil,
			),
			expected: `filter {
  if "true" == "true" and "true" == "true" or "true" == "true" nand "true" == "true" xor "true" == "true" {
    plugin {}
  }
}
`,
		},
		{
			name: "Condition in parentheses",
			config: NewConfig(
				nil,
				NewPluginSections(
					Filter, NewBranch(
						NewIfBlock(
							NewCondition(
								NewConditionExpression(
									ast.BooleanOperator{Op: NoOperator},
									NewCondition(
										NewInExpression(
											ast.BooleanOperator{Op: NoOperator}, NewStringAttribute("", "tag", DoubleQuoted), NewSelectorFromNames("tags"),
										),
									),
								),
							),
							NewPlugin("plugin"),
						),
						NewElseBlock(),
					),
				),
				nil,
			),
			expected: `filter {
  if ("tag" in [tags]) {
    plugin {}
  }
}
`,
		},
		{
			name: "Multiple conditions in parentheses",
			config: NewConfig(
				nil,
				NewPluginSections(
					Filter, NewBranch(
						NewIfBlock(
							NewCondition(
								NewConditionExpression(
									ast.BooleanOperator{Op: NoOperator},
									NewCondition(
										NewInExpression(
											ast.BooleanOperator{Op: NoOperator}, NewStringAttribute("", "tag", DoubleQuoted), NewSelectorFromNames("tags"),
										),
										NewConditionExpression(
											ast.BooleanOperator{Op: Or},
											NewCondition(
												NewCompareExpression(
													ast.BooleanOperator{Op: NoOperator},
													NewStringAttribute("", "true", DoubleQuoted),
													CompareOperator{Op: Equal},
													NewStringAttribute("", "true", DoubleQuoted),
												),
												NewCompareExpression(
													ast.BooleanOperator{Op: And},
													NewNumberAttribute("", 1),
													CompareOperator{Op: Equal},
													NewNumberAttribute("", 1),
												),
											),
										),
									),
								),
							),
							NewPlugin("plugin"),
						),
						NewElseBlock(),
					),
				),
				nil,
			),
			expected: `filter {
  if ("tag" in [tags] or ("true" == "true" and 1 == 1)) {
    plugin {}
  }
}
`,
		},
		{
			name: "Negative Condition Expression",
			config: NewConfig(
				nil,
				NewPluginSections(
					Filter, NewBranch(
						NewIfBlock(
							NewCondition(
								NewNegativeConditionExpression(
									ast.BooleanOperator{Op: NoOperator},
									NewCondition(
										NewCompareExpression(
											ast.BooleanOperator{Op: NoOperator},
											NewStringAttribute("", "true", DoubleQuoted),
											CompareOperator{Op: Equal},
											NewStringAttribute("", "true", DoubleQuoted),
										),
									),
								),
							),
							NewPlugin("plugin"),
						),
						NewElseBlock(),
					),
				),
				nil,
			),
			expected: `filter {
  if !("true" == "true") {
    plugin {}
  }
}
`,
		},
		{
			name: "Negative Selector Expression for value in subfield",
			config: NewConfig(
				nil,
				NewPluginSections(
					Filter, NewBranch(
						NewIfBlock(
							NewCondition(
								NewNegativeSelectorExpression(
									ast.BooleanOperator{Op: NoOperator},
									NewSelectorFromNames("field", "subfield"),
								),
							),
							NewPlugin("plugin"),
						),
						NewElseBlock(),
					),
				),
				nil,
			),
			expected: `filter {
  if ![field][subfield] {
    plugin {}
  }
}
`,
		},
		{
			name: "InExpression",
			config: NewConfig(
				nil,
				NewPluginSections(
					Filter, NewBranch(
						NewIfBlock(
							NewCondition(
								NewInExpression(
									ast.BooleanOperator{Op: NoOperator}, NewStringAttribute("", "tag", DoubleQuoted), NewSelectorFromNames("tags"),
								),
							),
							NewPlugin("plugin"),
						),
						NewElseBlock(),
					),
				),
				nil,
			),
			expected: `filter {
  if "tag" in [tags] {
    plugin {}
  }
}
`,
		},
		{
			name: "NotInExpression",
			config: NewConfig(
				nil,
				NewPluginSections(
					Filter, NewBranch(
						NewIfBlock(
							NewCondition(
								NewNotInExpression(
									ast.BooleanOperator{Op: NoOperator}, NewStringAttribute("", "tag", DoubleQuoted), NewSelectorFromNames("field", "subfield"),
								),
							),
							NewPlugin("plugin"),
						),
						NewElseBlock(),
					),
				),
				nil,
			),
			expected: `filter {
  if "tag" not in [field][subfield] {
    plugin {}
  }
}
`,
		},
		{
			name: "RegexpExpression (Match)",
			config: NewConfig(
				nil,
				NewPluginSections(
					Filter, NewBranch(
						NewIfBlock(
							NewCondition(
								NewRegexpExpression(
									ast.BooleanOperator{Op: NoOperator}, NewSelectorFromNames("field"), RegexpOperator{Op: RegexpMatch}, NewRegexp(".*"),
								),
							),
							NewPlugin("plugin"),
						),
						NewElseBlock(),
					),
				),
				nil,
			),
			expected: `filter {
  if [field] =~ /.*/ {
    plugin {}
  }
}
`,
		},
		{
			name: "RegexpExpression (Not Match)",
			config: NewConfig(
				nil,
				NewPluginSections(
					Filter, NewBranch(
						NewIfBlock(
							NewCondition(
								NewRegexpExpression(
									ast.BooleanOperator{Op: NoOperator}, NewSelectorFromNames("field"), RegexpOperator{Op: RegexpNotMatch}, NewRegexp(".*"),
								),
							),
							NewPlugin("plugin"),
						),
						NewElseBlock(),
					),
				),
				nil,
			),
			expected: `filter {
  if [field] !~ /.*/ {
    plugin {}
  }
}
`,
		},
		{
			name: "Rvalue",
			config: NewConfig(
				nil,
				NewPluginSections(
					Filter, NewBranch(
						NewIfBlock(
							NewCondition(
								NewRvalueExpression(
									ast.BooleanOperator{Op: NoOperator}, NewStringAttribute("", "string", DoubleQuoted),
								),
								NewRvalueExpression(
									ast.BooleanOperator{Op: Or}, NewNumberAttribute("", 10),
								),
								NewRvalueExpression(
									ast.BooleanOperator{Op: Or}, NewSelectorFromNames("field", "subfield"),
								),
								NewRvalueExpression(
									ast.BooleanOperator{Op: Or}, NewRegexp(".*"),
								),
							),
							NewPlugin("plugin"),
						),
						NewElseBlock(),
					),
				),
				nil,
			),
			expected: `filter {
  if "string" or 10 or [field][subfield] or /.*/ {
    plugin {}
  }
}
`,
		},
		{
			name: "nil values",
			config: NewConfig(
				nil,
				NewPluginSections(
					Filter,
					nil,
					NewPlugin("mutate", nil),
					nil,
					NewPlugin("alter",
						nil,
						NewStringAttribute("foo", "bar", Bareword),
						NewArrayAttribute("nil", nil),
						NewHashAttribute("nilHash", NewHashEntry(NewStringAttribute("", "nilEntry", Bareword), nil)),
						nil,
					),
					NewBranch(
						NewIfBlock(
							NewCondition(
								NewCompareExpression(
									ast.BooleanOperator{Op: NoOperator}, NewStringAttribute("", "true", DoubleQuoted), CompareOperator{Op: Equal}, NewStringAttribute("", "true", DoubleQuoted),
								),
							),
							nil,
						),
						NewElseBlock(nil),
						NewElseIfBlock(
							NewCondition(
								NewCompareExpression(
									ast.BooleanOperator{Op: NoOperator}, NewStringAttribute("", "false", DoubleQuoted), CompareOperator{Op: Equal}, nil,
								),
								nil,
							),
							nil,
						),
					),
					nil,
				),
				nil,
			),
			expected: `filter {
  mutate {}

  alter {
    foo => bar
    nil => []
    nilHash => {
      nilEntry => 
    }
  }

  if "true" == "true" {
  } else if "false" ==  {
  } else {
  }
}
`,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got := test.config.String()
			if got != test.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s\n\n", test.expected, got)
			}
		})
	}
}

func TestPluginID(t *testing.T) {
	cases := []struct {
		name   string
		plugin Plugin

		wantID  string
		wantErr bool
	}{
		{
			name: "double quoted id",
			plugin: NewPlugin("stdin",
				NewStringAttribute("id", "123", DoubleQuoted),
			),

			wantID: "123",
		},
		{
			name: "bareword id",
			plugin: NewPlugin("stdin",
				NewStringAttribute("id", "123", Bareword),
			),

			wantID: "123",
		},
		{
			name: "bareword id",
			plugin: NewPlugin("stdin",
				NewNumberAttribute("id", 123),
			),

			wantID:  "",
			wantErr: true,
		},
		{
			name: "multiple attributes with id",
			plugin: NewPlugin("stdin",
				NewStringAttribute("name", "fobar", Bareword),
				NewStringAttribute("id", "123", Bareword),
				NewStringAttribute("description", "the description", DoubleQuoted),
			),

			wantID: "123",
		},
		{
			name: "multiple attributes without id",
			plugin: NewPlugin("stdin",
				NewStringAttribute("name", "fobar", Bareword),
				NewStringAttribute("alternative", "baz", Bareword),
				NewStringAttribute("description", "the description", DoubleQuoted),
			),

			wantID:  "",
			wantErr: true,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			id, err := test.plugin.ID()

			if test.wantErr && err == nil {
				t.Error("Expected an error, but go none.")
			}
			if !test.wantErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if test.wantID != id {
				t.Errorf("Expected id to be %q, but got %q", test.wantID, id)
			}
		})
	}
}

func TestPluginType(t *testing.T) {
	cases := []struct {
		name     string
		input    PluginType
		expected string
	}{
		{
			name:     "input",
			input:    Input,
			expected: "input",
		},
		{
			name:     "filter",
			input:    Filter,
			expected: "filter",
		},
		{
			name:     "output",
			input:    Output,
			expected: "output",
		},
		{
			name:     "undefined plugin type 0",
			input:    0,
			expected: "undefined plugin type",
		},
		{
			name:     "undefined plugin type 4",
			input:    4,
			expected: "undefined plugin type",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if test.input.String() != test.expected {
				t.Errorf("Expected: %s, Got: %s", test.expected, test.input)
			}
		})
	}
}

func TestStringAttributeType(t *testing.T) {
	cases := []struct {
		name     string
		input    StringAttributeType
		expected string
	}{
		{
			name:     "double quote",
			input:    DoubleQuoted,
			expected: `"`,
		},
		{
			name:     "single quote",
			input:    SingleQuoted,
			expected: `'`,
		},
		{
			name:     "bareword",
			input:    Bareword,
			expected: ``,
		},
		{
			name:     "undefined string attribute type 0",
			input:    0,
			expected: "undefined string attribute type",
		},
		{
			name:     "undefined string attribute type 4",
			input:    4,
			expected: "undefined string attribute type",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if test.input.String() != test.expected {
				t.Errorf("Expected: %s, Got: %s", test.expected, test.input)
			}
		})
	}
}

func TestCompareOperator(t *testing.T) {
	cases := []struct {
		name     string
		input    CompareOperator
		expected string
	}{
		{
			name:     "equal",
			input:    CompareOperator{Op: Equal},
			expected: `==`,
		},
		{
			name:     "not equal",
			input:    CompareOperator{Op: NotEqual},
			expected: `!=`,
		},
		{
			name:     "less or equal",
			input:    CompareOperator{Op: LessOrEqual},
			expected: `<=`,
		},
		{
			name:     "greater or equal",
			input:    CompareOperator{Op: GreaterOrEqual},
			expected: ">=",
		},
		{
			name:     "less than",
			input:    CompareOperator{Op: LessThan},
			expected: "<",
		},
		{
			name:     "greater than",
			input:    CompareOperator{Op: GreaterThan},
			expected: ">",
		},
		{
			name:     "undefined compare operator 0",
			input:    CompareOperator{Op: 0},
			expected: "undefined compare operator",
		},
		{
			name:     "undefined compare operator 7",
			input:    CompareOperator{Op: 7},
			expected: "undefined compare operator",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if test.input.String() != test.expected {
				t.Errorf("Expected: %s, Got: %s", test.expected, test.input)
			}
		})
	}
}

func TestRegexpOperator(t *testing.T) {
	cases := []struct {
		name     string
		input    RegexpOperator
		expected string
	}{
		{
			name: "regex match",
			input: ast.RegexpOperator{
				Op: RegexpMatch,
			},
			expected: `=~`,
		},
		{
			name: "regex not match",
			input: ast.RegexpOperator{
				Op: RegexpNotMatch,
			},
			expected: `!~`,
		},
		{
			name: "undefined regexp operator 0",
			input: ast.RegexpOperator{
				Op: 0,
			},
			expected: "undefined regexp operator",
		},
		{
			name: "undefined regexp operator 3",
			input: ast.RegexpOperator{
				Op: 3,
			},
			expected: "undefined regexp operator",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if test.input.String() != test.expected {
				t.Errorf("Expected: %s, Got: %s", test.expected, test.input)
			}
		})
	}
}

func TestBooleanOperator(t *testing.T) {
	cases := []struct {
		name     string
		input    BooleanOperator
		expected string
	}{
		{
			name:     "no operator",
			input:    ast.BooleanOperator{Op: NoOperator},
			expected: ``,
		},
		{
			name:     "and",
			input:    ast.BooleanOperator{Op: And},
			expected: ` and `,
		},
		{
			name:     "or",
			input:    ast.BooleanOperator{Op: Or},
			expected: ` or `,
		},
		{
			name:     "nand",
			input:    ast.BooleanOperator{Op: Nand},
			expected: ` nand `,
		},
		{
			name:     "xor",
			input:    ast.BooleanOperator{Op: Xor},
			expected: ` xor `,
		},
		{
			name:     "undefined boolean operator 0",
			input:    ast.BooleanOperator{Op: 0},
			expected: "undefined boolean operator",
		},
		{
			name:     "undefined boolean operator 6",
			input:    ast.BooleanOperator{Op: 6},
			expected: "undefined boolean operator",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if test.input.String() != test.expected {
				t.Errorf("Expected: %s, Got: %s", test.expected, test.input)
			}
		})
	}
}
