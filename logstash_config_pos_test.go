package config_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	. "github.com/breml/logstash-config"
	"github.com/breml/logstash-config/ast"
)

func TestPos(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		recordForTypes map[string]bool

		want []ast.Pos
	}{
		{
			name: "config",
			input: `input {}
`,
			recordForTypes: map[string]bool{
				T(ast.Config{}):        true,
				T(ast.PluginSection{}): true,
			},

			want: []ast.Pos{
				{
					// Config
					Line:   1,
					Column: 1,
					Offset: 0,
				},
				{
					// PluginSection input
					Line:   1,
					Column: 1,
					Offset: 0,
				},
			},
		},
		{
			name: "plugin sections",
			input: `  input {}
  filter {}

  output {}
		`,
			recordForTypes: map[string]bool{
				T(ast.PluginSection{}): true,
			},

			want: []ast.Pos{
				{
					// PluginSection input
					Line:   1,
					Column: 3,
					Offset: 2,
				},
				{
					// PluginSection filter
					Line:   2,
					Column: 3,
					Offset: 13,
				},
				{
					// PluginSection output
					Line:   4,
					Column: 3,
					Offset: 26,
				},
			},
		},
		{
			name: "plugins",
			input: `input {
  stdin {}
}
filter {
  mutate {}
}
output {
  stdout {}
}
`,
			recordForTypes: map[string]bool{
				T(ast.Plugin{}): true,
			},

			want: []ast.Pos{
				{
					// PluginSection stdin
					Line:   2,
					Column: 3,
					Offset: 10,
				},
				{
					// PluginSection mutate
					Line:   5,
					Column: 3,
					Offset: 32,
				},
				{
					// PluginSection filter
					Line:   8,
					Column: 3,
					Offset: 55,
				},
			},
		},
		{
			name: "plugin attribute",
			input: `input {
  stdin {
    codec => plain {}
  }
}
`,
			recordForTypes: map[string]bool{
				T(ast.PluginAttribute{}): true,
			},

			want: []ast.Pos{
				{
					// plugin
					Line:   3,
					Column: 5,
					Offset: 22,
				},
			},
		},
		{
			name: "string attribute",
			input: `input {
  stdin {
    barestring => value
    singlequotestring => 'value'
    doublequotestrign => "value"
  }
}
`,
			recordForTypes: map[string]bool{
				T(ast.StringAttribute{}): true,
			},

			want: []ast.Pos{
				{
					// barestring
					Line:   3,
					Column: 5,
					Offset: 22,
				},
				{
					// singlequotestring
					Line:   4,
					Column: 5,
					Offset: 46,
				},
				{
					// doublequotestring
					Line:   5,
					Column: 5,
					Offset: 79,
				},
			},
		},
		{
			name: "number attribute",
			input: `input {
  stdin {
    number => 10
  }
}
`,
			recordForTypes: map[string]bool{
				T(ast.NumberAttribute{}): true,
			},

			want: []ast.Pos{
				{
					// number
					Line:   3,
					Column: 5,
					Offset: 22,
				},
			},
		},
		{
			name: "array attribute empty",
			input: `input {
  stdin {
    array => []
  }
}
`,
			recordForTypes: map[string]bool{
				T(ast.ArrayAttribute{}): true,
			},

			want: []ast.Pos{
				{
					// array
					Line:   3,
					Column: 5,
					Offset: 22,
				},
			},
		},
		{
			name: "array attribute with elements",
			input: `input {
  stdin {
    array => [
      string,
      "string",
      'string'
    ]
  }
}
`,
			recordForTypes: map[string]bool{
				T(ast.StringAttribute{}): true,
				T(ast.ArrayAttribute{}):  true,
			},

			want: []ast.Pos{
				{
					// array
					Line:   3,
					Column: 5,
					Offset: 22,
				},
				{
					// barestring
					Line:   4,
					Column: 7,
					Offset: 39,
				},
				{
					// doublequotedstring
					Line:   5,
					Column: 7,
					Offset: 53,
				},
				{
					// singlequotedstring
					Line:   6,
					Column: 7,
					Offset: 69,
				},
			},
		},
		{
			name: "hash attribute empty",
			input: `input {
  stdin {
    hash => {}
  }
}
`,
			recordForTypes: map[string]bool{
				T(ast.HashAttribute{}): true,
			},

			want: []ast.Pos{
				{
					// hash
					Line:   3,
					Column: 5,
					Offset: 22,
				},
			},
		},
		{
			name: "hash attribute with elements",
			input: `input {
  stdin {
    hash => {
      bareword => value
      "doublequoted" => "value"
      'singlequoted' => 'value'
    }
  }
}
`,
			recordForTypes: map[string]bool{
				T(ast.HashAttribute{}):   true,
				T(ast.HashEntry{}):       true,
				T(ast.StringAttribute{}): true,
			},

			want: []ast.Pos{
				{
					// hash
					Line:   3,
					Column: 5,
					Offset: 22,
				},
				{
					// barestring hash entry
					Line:   4,
					Column: 7,
					Offset: 38,
				},
				{
					// barestring hash entry key
					Line:   4,
					Column: 7,
					Offset: 38,
				},
				{
					// barestring hash entry value
					Line:   4,
					Column: 19,
					Offset: 50,
				},
				{
					// doublequoted hash entry
					Line:   5,
					Column: 7,
					Offset: 62,
				},
				{
					// doublequoted hash entry key
					Line:   5,
					Column: 7,
					Offset: 62,
				},
				{
					// doublequoted hash entry value
					Line:   5,
					Column: 25,
					Offset: 80,
				},
				{
					// singlequoted hash entry
					Line:   6,
					Column: 7,
					Offset: 94,
				},
				{
					// singlequoted hash entry key
					Line:   6,
					Column: 7,
					Offset: 94,
				},
				{
					// singlequoted hash entry value
					Line:   6,
					Column: 25,
					Offset: 112,
				},
			},
		},
		{
			name: "if with multiple else-if and a final else branch, test for different condition types",
			input: `filter {
  if "true" == "true" {
    if-plugin {}
  } else if !("true" == "true") {
    else-if-plugin-1 {}
  } else if "true" == "true" and "true" == "true" or "true" == "true" nand "true" == "true" xor "true" == "true" {
    else-if-plugin-2 {}
  } else if ("tag" in [tags]) {
    else-if-plugin-3 {}
  } else if ("tag" in [tags] or ("true" == "true" and 1 == 1)) {
    else-if-plugin-4 {}
  } else if ![field][subfield] {
    else-if-plugin-5 {}
  } else if "tag" in [tags] {
    else-if-plugin-6 {}
  } else if "tag" not in [field][subfield] {
    else-if-plugin-7 {}
  } else if [field] =~ /.*/ {
    else-if-plugin-8 {}
  } else if [field] !~ /.*/ {
    else-if-plugin-9 {}
  } else if "string" or 10 or [field][subfield] or /.*/ {
    else-if-plugin-10 {}
  } else {
    else-plugin {}
  }
}
`,
			recordForTypes: map[string]bool{
				T(ast.Branch{}):                      true,
				T(ast.IfBlock{}):                     true,
				T(ast.ElseIfBlock{}):                 true,
				T(ast.ElseBlock{}):                   true,
				T(ast.Condition{}):                   true,
				T(ast.ConditionExpression{}):         true,
				T(ast.NegativeConditionExpression{}): true,
				T(ast.NegativeSelectorExpression{}):  true,
				T(ast.InExpression{}):                true,
				T(ast.NotInExpression{}):             true,
				T(ast.RvalueExpression{}):            true,
				T(ast.CompareExpression{}):           true,
				T(ast.CompareOperator{}):             true,
				T(ast.RegexpExpression{}):            true,
				T(ast.Regexp{}):                      true,
				T(ast.RegexpOperator{}):              true,
				T(ast.BooleanOperator{}):             true,
				T(ast.Selector{}):                    true,
				T(ast.SelectorElement{}):             true,
			},

			want: []ast.Pos{
				{
					// ast.Branch
					Line:   2,
					Column: 3,
					Offset: 11,
				},
				{
					// ast.IfBlock
					Line:   2,
					Column: 3,
					Offset: 11,
				},
				{
					// ast.Condition
					Line:   2,
					Column: 6,
					Offset: 14,
				},
				{
					// ast.CompareExpression
					Line:   2,
					Column: 6,
					Offset: 14,
				},
				{
					// ast.CompareOperator
					Line:   2,
					Column: 13,
					Offset: 21,
				},
				{
					// ast.ElseIfBlock
					Line:   4,
					Column: 5,
					Offset: 54,
				},
				{
					// ast.Condition
					Line:   4,
					Column: 13,
					Offset: 62,
				},
				{
					// ast.NegativeConditionExpression
					Line:   4,
					Column: 13,
					Offset: 62,
				},
				{
					// ast.Condition
					Line:   4,
					Column: 15,
					Offset: 64,
				},
				{
					// ast.CompareExpression
					Line:   4,
					Column: 15,
					Offset: 64,
				},
				{
					// ast.CompareOperator
					Line:   4,
					Column: 22,
					Offset: 71,
				},
				{
					// ast.ElseIfBlock
					Line:   6,
					Column: 5,
					Offset: 112,
				},
				{
					// ast.Condition
					Line:   6,
					Column: 13,
					Offset: 120,
				},
				{
					// ast.CompareExpression
					Line:   6,
					Column: 13,
					Offset: 120,
				},
				{
					// ast.CompareOperator
					Line:   6,
					Column: 20,
					Offset: 127,
				},
				{
					// ast.BooleanOperator
					Line:   6,
					Column: 30,
					Offset: 137,
				},
				{
					// ast.CompareExpression
					Line:   6,
					Column: 34,
					Offset: 141,
				},
				{
					// ast.CompareOperator
					Line:   6,
					Column: 41,
					Offset: 148,
				},
				{
					// ast.BooleanOperator
					Line:   6,
					Column: 51,
					Offset: 158,
				},
				{
					// ast.CompareExpression
					Line:   6,
					Column: 54,
					Offset: 161,
				},
				{
					// ast.CompareOperator
					Line:   6,
					Column: 61,
					Offset: 168,
				},
				{
					// ast.BooleanOperator
					Line:   6,
					Column: 71,
					Offset: 178,
				},
				{
					// ast.CompareExpression
					Line:   6,
					Column: 76,
					Offset: 183,
				},
				{
					// ast.CompareOperator
					Line:   6,
					Column: 83,
					Offset: 190,
				},
				{
					// ast.BooleanOperator
					Line:   6,
					Column: 93,
					Offset: 200,
				},
				{
					// ast.CompareExpression
					Line:   6,
					Column: 97,
					Offset: 204,
				},
				{
					// ast.CompareOperator
					Line:   6,
					Column: 104,
					Offset: 211,
				},
				{
					// ast.ElseIfBlock
					Line:   8,
					Column: 5,
					Offset: 251,
				},
				{
					// ast.Condition
					Line:   8,
					Column: 13,
					Offset: 259,
				},
				{
					// ast.ConditionExpression
					Line:   8,
					Column: 13,
					Offset: 259,
				},
				{
					// ast.Condition
					Line:   8,
					Column: 14,
					Offset: 260,
				},
				{
					// ast.InExpression
					Line:   8,
					Column: 14,
					Offset: 260,
				},
				{
					// ast.Selector
					Line:   8,
					Column: 23,
					Offset: 269,
				},
				{
					// ast.SelectorElement
					Line:   8,
					Column: 23,
					Offset: 269,
				},
				{
					// ast.ElseIfBlock
					Line:   10,
					Column: 5,
					Offset: 307,
				},
				{
					// ast.Condition
					Line:   10,
					Column: 13,
					Offset: 315,
				},
				{
					// ast.ConditionExpression
					Line:   10,
					Column: 13,
					Offset: 315,
				},
				{
					// ast.Condition
					Line:   10,
					Column: 14,
					Offset: 316,
				},
				{
					// ast.InExpression
					Line:   10,
					Column: 14,
					Offset: 316,
				},
				{
					// ast.Selector
					Line:   10,
					Column: 23,
					Offset: 325,
				},
				{
					// ast.SelectorElement
					Line:   10,
					Column: 23,
					Offset: 325,
				},
				{
					// ast.BooleanOperator
					Line:   10,
					Column: 30,
					Offset: 332,
				},
				{
					// ast.ConditionExpression
					Line:   10,
					Column: 33,
					Offset: 335,
				},
				{
					// ast.Condition
					Line:   10,
					Column: 34,
					Offset: 336,
				},
				{
					// ast.CompareExpression
					Line:   10,
					Column: 34,
					Offset: 336,
				},
				{
					// ast.CompareOperator
					Line:   10,
					Column: 41,
					Offset: 343,
				},
				{
					// ast.BooleanOperator
					Line:   10,
					Column: 51,
					Offset: 353,
				},
				{
					// ast.CompareExpression
					Line:   10,
					Column: 55,
					Offset: 357,
				},
				{
					// ast.CompareOperator
					Line:   10,
					Column: 57,
					Offset: 359,
				},
				{
					// ast.ElseIfBlock
					Line:   12,
					Column: 5,
					Offset: 396,
				},
				{
					// ast.Condition
					Line:   12,
					Column: 13,
					Offset: 404,
				},
				{
					// ast.NegativeSelectorExpression
					Line:   12,
					Column: 13,
					Offset: 404,
				},
				{
					// ast.Selector
					Line:   12,
					Column: 14,
					Offset: 405,
				},
				{
					// ast.SelectorElement
					Line:   12,
					Column: 14,
					Offset: 405,
				},
				{
					// ast.SelectorElement
					Line:   12,
					Column: 21,
					Offset: 412,
				},
				{
					// ast.ElseIfBlock
					Line:   14,
					Column: 5,
					Offset: 453,
				},
				{
					// ast.Condition
					Line:   14,
					Column: 13,
					Offset: 461,
				},
				{
					// ast.InExpression
					Line:   14,
					Column: 13,
					Offset: 461,
				},
				{
					// ast.Selector
					Line:   14,
					Column: 22,
					Offset: 470,
				},
				{
					// ast.SelectorElement
					Line:   14,
					Column: 22,
					Offset: 470,
				},
				{
					// ast.ElseIfBlock
					Line:   16,
					Column: 5,
					Offset: 507,
				},
				{
					// ast.Condition
					Line:   16,
					Column: 13,
					Offset: 515,
				},
				{
					// ast.NotInExpression
					Line:   16,
					Column: 13,
					Offset: 515,
				},
				{
					// ast.Selector
					Line:   16,
					Column: 26,
					Offset: 528,
				},
				{
					// ast.SelectorElement
					Line:   16,
					Column: 26,
					Offset: 528,
				},
				{
					// ast.SelectorElement
					Line:   16,
					Column: 33,
					Offset: 535,
				},
				{
					// ast.ElseIfBlock
					Line:   18,
					Column: 5,
					Offset: 576,
				},
				{
					// ast.Condition
					Line:   18,
					Column: 13,
					Offset: 584,
				},
				{
					// ast.RegexpExpression
					Line:   18,
					Column: 13,
					Offset: 584,
				},
				{
					// ast.Selector
					Line:   18,
					Column: 13,
					Offset: 584,
				},
				{
					// ast.SelectorElement
					Line:   18,
					Column: 13,
					Offset: 584,
				},
				{
					// ast.RegexpOperator
					Line:   18,
					Column: 21,
					Offset: 592,
				},
				{
					// ast.Regexp
					Line:   18,
					Column: 24,
					Offset: 595,
				},
				{
					// ast.ElseIfBlock
					Line:   20,
					Column: 5,
					Offset: 630,
				},
				{
					// ast.Condition
					Line:   20,
					Column: 13,
					Offset: 638,
				},
				{
					// ast.RegexpExpression
					Line:   20,
					Column: 13,
					Offset: 638,
				},
				{
					// ast.Selector
					Line:   20,
					Column: 13,
					Offset: 638,
				},
				{
					// ast.SelectorElement
					Line:   20,
					Column: 13,
					Offset: 638,
				},
				{
					// ast.RegexpOperator
					Line:   20,
					Column: 21,
					Offset: 646,
				},
				{
					// ast.Regexp
					Line:   20,
					Column: 24,
					Offset: 649,
				},
				{
					// ast.ElseIfBlock
					Line:   22,
					Column: 5,
					Offset: 684,
				},
				{
					// ast.Condition
					Line:   22,
					Column: 13,
					Offset: 692,
				},
				{
					// ast.RvalueExpression
					Line:   22,
					Column: 13,
					Offset: 692,
				},
				{
					// ast.BooleanOperator
					Line:   22,
					Column: 22,
					Offset: 701,
				},
				{
					// ast.RvalueExpression
					Line:   22,
					Column: 25,
					Offset: 704,
				},
				{
					// ast.BooleanOperator
					Line:   22,
					Column: 28,
					Offset: 707,
				},
				{
					// ast.RvalueExpression
					Line:   22,
					Column: 31,
					Offset: 710,
				},
				{
					// ast.Selector
					Line:   22,
					Column: 31,
					Offset: 710,
				},
				{
					// ast.SelectorElement
					Line:   22,
					Column: 31,
					Offset: 710,
				},
				{
					// ast.SelectorElement
					Line:   22,
					Column: 38,
					Offset: 717,
				},
				{
					// ast.BooleanOperator
					Line:   22,
					Column: 49,
					Offset: 728,
				},
				{
					// ast.RvalueExpression
					Line:   22,
					Column: 52,
					Offset: 731,
				},
				{
					// ast.Regexp
					Line:   22,
					Column: 52,
					Offset: 731,
				},
				{
					// ast.ElseBlock
					Line:   24,
					Column: 5,
					Offset: 767,
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseReader("test", strings.NewReader(test.input))
			if err != nil {
				t.Fatalf("Expected to parse without error: %s, input:\n%s", err, test.input)
			}
			conf := got.(ast.Config)

			p := positionWalker{
				recordForTypes: test.recordForTypes,
			}
			p.walk(conf)

			if !reflect.DeepEqual(test.want, p.positions) {
				t.Fatalf("Expect %#v to be equal to %#v", test.want, p.positions)
			}
		})
	}
}

func T(in interface{}) string {
	return fmt.Sprintf("%T", in)
}

type positionWalker struct {
	positions      []ast.Pos
	recordForTypes map[string]bool
}

func (p *positionWalker) walk(node ast.Node) {
	if n, ok := node.(ast.Expression); ok {
		if n.BoolOperator().Op != ast.NoOperator {
			p.walk(n.BoolOperator())
		}
	}

	if p.recordForTypes[T(node)] {
		p.positions = append(p.positions, node.Pos())
		// fmt.Printf("%v %T\n", node.Pos(), node)
	}

	switch n := node.(type) {
	case ast.Config:
		for _, input := range n.Input {
			p.walk(input)
		}
		for _, filter := range n.Filter {
			p.walk(filter)
		}
		for _, output := range n.Output {
			p.walk(output)
		}

	case ast.PluginSection:
		for _, bop := range n.BranchOrPlugins {
			p.walk(bop)
		}

	case ast.Plugin:
		for _, attr := range n.Attributes {
			p.walk(attr)
		}

	case ast.ArrayAttribute:
		for _, attr := range n.Attributes {
			p.walk(attr)
		}

	case ast.HashAttribute:
		for _, entry := range n.Entries {
			p.walk(entry)
		}

	case ast.HashEntry:
		p.walk(n.Key)
		p.walk(n.Value)

	case ast.Branch:
		p.walk(n.IfBlock)
		for _, block := range n.ElseIfBlock {
			p.walk(block)
		}
		p.walk(n.ElseBlock)

	case ast.IfBlock:
		p.walk(n.Condition)
		for _, block := range n.Block {
			p.walk(block)
		}

	case ast.ElseIfBlock:
		p.walk(n.Condition)
		for _, block := range n.Block {
			p.walk(block)
		}

	case ast.ElseBlock:
		for _, block := range n.Block {
			p.walk(block)
		}

	case ast.Condition:
		for _, expression := range n.Expression {
			p.walk(expression)
		}

	case ast.ConditionExpression:
		p.walk(n.Condition)

	case ast.NegativeConditionExpression:
		p.walk(n.Condition)

	case ast.NegativeSelectorExpression:
		p.walk(n.Selector)

	case ast.InExpression:
		p.walk(n.LValue)
		p.walk(n.RValue)

	case ast.NotInExpression:
		p.walk(n.LValue)
		p.walk(n.RValue)

	case ast.RvalueExpression:
		p.walk(n.RValue)

	case ast.CompareExpression:
		p.walk(n.LValue)
		p.walk(n.CompareOperator)
		p.walk(n.RValue)

	case ast.RegexpExpression:
		p.walk(n.LValue)
		p.walk(n.RegexpOperator)
		p.walk(n.RValue)

	case ast.Selector:
		for _, element := range n.Elements {
			p.walk(element)
		}

	}
}
