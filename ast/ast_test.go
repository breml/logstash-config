package ast_test

import (
	"testing"

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
    tags => [
      "tag1",
      'tag2',
      tag3
    ]
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
		{
			name: "Simple if (without else) branch",
			config: NewConfig(
				nil,
				NewPluginSections(
					Filter, NewBranch(
						NewIfBlock(
							NewCondition(
								NewCompareExpression(
									NoOperator, NewStringAttribute("", "true", DoubleQuoted), Equal, NewStringAttribute("", "true", DoubleQuoted),
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
									NoOperator, NewStringAttribute("", "true", DoubleQuoted), Equal, NewStringAttribute("", "true", DoubleQuoted),
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
									NoOperator, NewStringAttribute("", "true", DoubleQuoted), Equal, NewStringAttribute("", "true", DoubleQuoted),
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
									NoOperator, NewStringAttribute("", "true", DoubleQuoted), NotEqual, NewStringAttribute("", "true", DoubleQuoted),
								),
							),
							NewPlugin("else-if-plugin-1"),
						),
						NewElseIfBlock(
							NewCondition(
								NewCompareExpression(
									NoOperator, NewNumberAttribute("", 10), GreaterThan, NewNumberAttribute("", 2),
								),
							),
							NewPlugin("else-if-plugin-2"),
						),
						NewElseIfBlock(
							NewCondition(
								NewCompareExpression(
									NoOperator, NewNumberAttribute("", 2), LessThan, NewNumberAttribute("", 10),
								),
							),
							NewPlugin("else-if-plugin-3"),
						),
						NewElseIfBlock(
							NewCondition(
								NewCompareExpression(
									NoOperator, NewNumberAttribute("", 10), GreaterOrEqual, NewNumberAttribute("", 2),
								),
							),
							NewPlugin("else-if-plugin-4"),
						),
						NewElseIfBlock(
							NewCondition(
								NewCompareExpression(
									NoOperator, NewNumberAttribute("", 2), LessOrEqual, NewNumberAttribute("", 10),
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
									NoOperator, NewStringAttribute("", "true", DoubleQuoted), Equal, NewStringAttribute("", "true", DoubleQuoted),
								),
								NewCompareExpression(
									And, NewStringAttribute("", "true", DoubleQuoted), Equal, NewStringAttribute("", "true", DoubleQuoted),
								),
								NewCompareExpression(
									Or, NewStringAttribute("", "true", DoubleQuoted), Equal, NewStringAttribute("", "true", DoubleQuoted),
								),
								NewCompareExpression(
									Nand, NewStringAttribute("", "true", DoubleQuoted), Equal, NewStringAttribute("", "true", DoubleQuoted),
								),
								NewCompareExpression(
									Xor, NewStringAttribute("", "true", DoubleQuoted), Equal, NewStringAttribute("", "true", DoubleQuoted),
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
									NoOperator,
									NewCondition(
										NewInExpression(
											NoOperator, NewStringAttribute("", "tag", DoubleQuoted), NewSelectorFromNames("tags"),
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
									NoOperator,
									NewCondition(
										NewInExpression(
											NoOperator, NewStringAttribute("", "tag", DoubleQuoted), NewSelectorFromNames("tags"),
										),
										NewConditionExpression(
											Or,
											NewCondition(
												NewCompareExpression(
													NoOperator,
													NewStringAttribute("", "true", DoubleQuoted),
													Equal,
													NewStringAttribute("", "true", DoubleQuoted),
												),
												NewCompareExpression(
													And,
													NewNumberAttribute("", 1),
													Equal,
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
									NoOperator,
									NewCondition(
										NewCompareExpression(
											NoOperator,
											NewStringAttribute("", "true", DoubleQuoted),
											Equal,
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
  if ! ("true" == "true") {
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
									NoOperator,
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
  if ! [field][subfield] {
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
									NoOperator, NewStringAttribute("", "tag", DoubleQuoted), NewSelectorFromNames("tags"),
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
									NoOperator, NewStringAttribute("", "tag", DoubleQuoted), NewSelectorFromNames("field", "subfield"),
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
									NoOperator, NewSelectorFromNames("field"), RegexpMatch, NewRegexp(".*"),
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
									NoOperator, NewSelectorFromNames("field"), RegexpNotMatch, NewRegexp(".*"),
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
									NoOperator, NewStringAttribute("", "string", DoubleQuoted),
								),
								NewRvalueExpression(
									Or, NewNumberAttribute("", 10),
								),
								NewRvalueExpression(
									Or, NewSelectorFromNames("field", "subfield"),
								),
								NewRvalueExpression(
									Or, NewRegexp(".*"),
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
						NewHashAttribute("nilHash", NewHashEntry("nilEntry", nil)),
						nil,
					),
					NewBranch(
						NewIfBlock(
							NewCondition(
								NewCompareExpression(
									NoOperator, NewStringAttribute("", "true", DoubleQuoted), Equal, NewStringAttribute("", "true", DoubleQuoted),
								),
							),
							nil,
						),
						NewElseBlock(nil),
						NewElseIfBlock(
							NewCondition(
								NewCompareExpression(
									NoOperator, NewStringAttribute("", "false", DoubleQuoted), Equal, nil,
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
			input:    Equal,
			expected: `==`,
		},
		{
			name:     "not equal",
			input:    NotEqual,
			expected: `!=`,
		},
		{
			name:     "less or equal",
			input:    LessOrEqual,
			expected: `<=`,
		},
		{
			name:     "greater or equal",
			input:    GreaterOrEqual,
			expected: ">=",
		},
		{
			name:     "less than",
			input:    LessThan,
			expected: "<",
		},
		{
			name:     "greater than",
			input:    GreaterThan,
			expected: ">",
		},
		{
			name:     "undefined compare operator 0",
			input:    0,
			expected: "undefined compare operator",
		},
		{
			name:     "undefined compare operator 7",
			input:    7,
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
			name:     "regex match",
			input:    RegexpMatch,
			expected: `=~`,
		},
		{
			name:     "regex not match",
			input:    RegexpNotMatch,
			expected: `!~`,
		},
		{
			name:     "undefined regexp operator 0",
			input:    0,
			expected: "undefined regexp operator",
		},
		{
			name:     "undefined regexp operator 3",
			input:    3,
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
			input:    NoOperator,
			expected: ``,
		},
		{
			name:     "and",
			input:    And,
			expected: ` and `,
		},
		{
			name:     "or",
			input:    Or,
			expected: ` or `,
		},
		{
			name:     "nand",
			input:    Nand,
			expected: ` nand `,
		},
		{
			name:     "xor",
			input:    Xor,
			expected: ` xor `,
		},
		{
			name:     "undefined boolean operator 0",
			input:    0,
			expected: "undefined boolean operator",
		},
		{
			name:     "undefined boolean operator 6",
			input:    6,
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
