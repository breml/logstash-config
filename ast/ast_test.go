package ast_test

import (
	"testing"

	. "github.com/breml/logstash-config/ast"
)

func TestAst(t *testing.T) {
	cases := []struct {
		config   Config
		expected string
	}{
		// The empty config
		{
			config:   Config{},
			expected: ``,
		},

		// Plugins with attributes of various types
		{
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

		// Simple if (without else) branch
		{
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

		// Simple if-else branch
		{
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

		// if with multiple else-if and a final else branch
		// test for different condition types
		{
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

		// if with multiple compare operators
		{
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

		// Condition in parentheses
		{
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

		// Multiple conditions in parentheses
		{
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

		// Negative Condition Expression
		{
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

		// Negative Selector Expression for value in subfield
		{
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

		// InExpression
		{
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

		// NotInExpression
		{
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

		// RegexpExpression (Match)
		{
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

		// RegexpExpression (Not Match)
		{
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

		// Rvalue
		{
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
    nil => [  ]
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
			expected: "undefined plugin type",
		},
		{
			input:    4,
			expected: "undefined plugin type",
		},
	}

	for _, test := range cases {
		if test.input.String() != test.expected {
			t.Errorf("Expected: %s, Got: %s", test.expected, test.input)
		}
	}
}

func TestStringAttributeType(t *testing.T) {
	cases := []struct {
		input    StringAttributeType
		expected string
	}{
		{
			input:    DoubleQuoted,
			expected: `"`,
		},
		{
			input:    SingleQuoted,
			expected: `'`,
		},
		{
			input:    Bareword,
			expected: ``,
		},
		{
			input:    0,
			expected: "undefined string attribute type",
		},
		{
			input:    4,
			expected: "undefined string attribute type",
		},
	}

	for _, test := range cases {
		if test.input.String() != test.expected {
			t.Errorf("Expected: %s, Got: %s", test.expected, test.input)
		}
	}
}

func TestCompareOperator(t *testing.T) {
	cases := []struct {
		input    CompareOperator
		expected string
	}{
		{
			input:    Equal,
			expected: `==`,
		},
		{
			input:    NotEqual,
			expected: `!=`,
		},
		{
			input:    LessOrEqual,
			expected: `<=`,
		},
		{
			input:    GreaterOrEqual,
			expected: ">=",
		},
		{
			input:    LessThan,
			expected: "<",
		},
		{
			input:    GreaterThan,
			expected: ">",
		},
		{
			input:    0,
			expected: "undefined compare operator",
		},
		{
			input:    7,
			expected: "undefined compare operator",
		},
	}

	for _, test := range cases {
		if test.input.String() != test.expected {
			t.Errorf("Expected: %s, Got: %s", test.expected, test.input)
		}
	}
}

func TestRegexpOperator(t *testing.T) {
	cases := []struct {
		input    RegexpOperator
		expected string
	}{
		{
			input:    RegexpMatch,
			expected: `=~`,
		},
		{
			input:    RegexpNotMatch,
			expected: `!~`,
		},
		{
			input:    0,
			expected: "undefined regexp operator",
		},
		{
			input:    3,
			expected: "undefined regexp operator",
		},
	}

	for _, test := range cases {
		if test.input.String() != test.expected {
			t.Errorf("Expected: %s, Got: %s", test.expected, test.input)
		}
	}
}

func TestBooleanOperator(t *testing.T) {
	cases := []struct {
		input    BooleanOperator
		expected string
	}{
		{
			input:    NoOperator,
			expected: ``,
		},
		{
			input:    And,
			expected: ` and `,
		},
		{
			input:    Or,
			expected: ` or `,
		},
		{
			input:    Nand,
			expected: ` nand `,
		},
		{
			input:    Xor,
			expected: ` xor `,
		},
		{
			input:    0,
			expected: "undefined boolean operator",
		},
		{
			input:    6,
			expected: "undefined boolean operator",
		},
	}

	for _, test := range cases {
		if test.input.String() != test.expected {
			t.Errorf("Expected: %s, Got: %s", test.expected, test.input)
		}
	}
}
