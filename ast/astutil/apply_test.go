package astutil_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/breml/logstash-config/ast"
	"github.com/breml/logstash-config/ast/astutil"
)

func TestApplyPlugins_Walk(t *testing.T) {
	cases := []struct {
		name  string
		input []ast.BranchOrPlugin

		wantCount int
	}{
		{
			name: "nil",
		},
		{
			name: "nil BranchOrPlugin",
			input: []ast.BranchOrPlugin{
				nil,
			},

			wantCount: 1,
		},
		{
			name: "plugin",
			input: []ast.BranchOrPlugin{
				ast.NewPlugin("plugin"),
			},

			wantCount: 1,
		},
		{
			name: "if with nil, else if with nil else with nil",
			input: []ast.BranchOrPlugin{
				ast.NewBranch(
					ast.NewIfBlock(ast.NewCondition(), nil),
					ast.NewElseBlock(nil),
					ast.NewElseIfBlock(ast.NewCondition(), nil),
				),
			},

			wantCount: 3,
		},
		{
			name: "if with nil, else if with nil else with nil",
			input: []ast.BranchOrPlugin{
				ast.NewBranch(
					ast.NewIfBlock(ast.NewCondition(), ast.NewPlugin("plugin")),
					ast.NewElseBlock(ast.NewPlugin("plugin")),
					ast.NewElseIfBlock(ast.NewCondition(), ast.NewPlugin("plugin"))),
			},

			wantCount: 3,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var count int

			applyFunc := func(c *astutil.Cursor) {
				count++
				if strings.Compare(fmt.Sprint(c.Parent()[c.Index()]), fmt.Sprint(c.Plugin())) != 0 {
					t.Fatalf("expect element at index in parent: %v to be equal to plugin: %v", c.Parent()[c.Index()], c.Plugin())
				}
			}
			astutil.ApplyPlugins(test.input, applyFunc)

			if test.wantCount != count {
				t.Fatalf("Expected walkFn to be called %d, but got: %d", test.wantCount, count)
			}
		})
	}
}

func TestApplyPlugins_Delete(t *testing.T) {
	cases := []struct {
		name  string
		input []ast.BranchOrPlugin

		want []ast.BranchOrPlugin
	}{
		{
			name: "plugin",
			input: []ast.BranchOrPlugin{
				ast.NewPlugin("delete_me"),
			},

			want: []ast.BranchOrPlugin{},
		},
		{
			name: "first and last",
			input: []ast.BranchOrPlugin{
				ast.NewPlugin("delete_me"),
				ast.NewPlugin("plugin 1"),
				ast.NewPlugin("plugin 2"),
				ast.NewPlugin("delete_me"),
			},

			want: []ast.BranchOrPlugin{
				ast.NewPlugin("plugin 1"),
				ast.NewPlugin("plugin 2"),
			},
		},
		{
			name: "middle",
			input: []ast.BranchOrPlugin{
				ast.NewPlugin("plugin 1"),
				ast.NewPlugin("delete_me"),
				ast.NewPlugin("plugin 2"),
			},

			want: []ast.BranchOrPlugin{
				ast.NewPlugin("plugin 1"),
				ast.NewPlugin("plugin 2"),
			},
		},
		{
			name: "if with plugin, else if with plugin else with plugin",
			input: []ast.BranchOrPlugin{
				ast.Branch{
					IfBlock: ast.IfBlock{
						Block: []ast.BranchOrPlugin{
							ast.NewPlugin("delete_me"),
						},
					},
					ElseBlock: ast.ElseBlock{
						Block: []ast.BranchOrPlugin{
							ast.NewPlugin("delete_me"),
						},
					},
					ElseIfBlock: []ast.ElseIfBlock{
						{
							Block: []ast.BranchOrPlugin{
								ast.NewPlugin("delete_me"),
							},
						},
					},
				},
			},

			want: []ast.BranchOrPlugin{
				ast.Branch{
					IfBlock: ast.IfBlock{
						Block: []ast.BranchOrPlugin{},
					},
					ElseBlock: ast.ElseBlock{
						Block: []ast.BranchOrPlugin{},
					},
					ElseIfBlock: []ast.ElseIfBlock{
						{
							Block: []ast.BranchOrPlugin{},
						},
					},
				},
			},
		},
		{
			name: "nested if with plugin",
			input: []ast.BranchOrPlugin{
				ast.Branch{
					IfBlock: ast.IfBlock{
						Block: []ast.BranchOrPlugin{
							ast.Branch{
								IfBlock: ast.IfBlock{
									Block: []ast.BranchOrPlugin{
										ast.NewPlugin("delete_me"),
									},
								},
							},
						},
					},
				},
			},

			want: []ast.BranchOrPlugin{
				ast.Branch{
					IfBlock: ast.IfBlock{
						Block: []ast.BranchOrPlugin{
							ast.Branch{
								IfBlock: ast.IfBlock{
									Block: []ast.BranchOrPlugin{},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			applyFunc := func(c *astutil.Cursor) {
				if c.Plugin().Name() != "delete_me" {
					return
				}
				c.Delete()
			}
			got := astutil.ApplyPlugins(test.input, applyFunc)

			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("Expect %#v to be equal to %#v", test.want, got)
			}
		})
	}
}

func TestApplyPlugins_Replace(t *testing.T) {
	cases := []struct {
		name  string
		input []ast.BranchOrPlugin

		want []ast.BranchOrPlugin
	}{
		{
			name: "plugin",
			input: []ast.BranchOrPlugin{
				ast.NewPlugin("plugin"),
			},

			want: []ast.BranchOrPlugin{
				ast.NewPlugin("replacement"),
			},
		},
		{
			name: "if with plugin, else if with plugin else with plugin",
			input: []ast.BranchOrPlugin{
				ast.Branch{
					IfBlock: ast.IfBlock{
						Block: []ast.BranchOrPlugin{
							ast.NewPlugin("plugin"),
						},
					},
					ElseBlock: ast.ElseBlock{
						Block: []ast.BranchOrPlugin{
							ast.NewPlugin("plugin"),
						},
					},
					ElseIfBlock: []ast.ElseIfBlock{
						{
							Block: []ast.BranchOrPlugin{
								ast.NewPlugin("plugin"),
							},
						},
					},
				},
			},

			want: []ast.BranchOrPlugin{
				ast.Branch{
					IfBlock: ast.IfBlock{
						Block: []ast.BranchOrPlugin{
							ast.NewPlugin("replacement"),
						},
					},
					ElseBlock: ast.ElseBlock{
						Block: []ast.BranchOrPlugin{
							ast.NewPlugin("replacement"),
						},
					},
					ElseIfBlock: []ast.ElseIfBlock{
						{
							Block: []ast.BranchOrPlugin{
								ast.NewPlugin("replacement"),
							},
						},
					},
				},
			},
		},
		{
			name: "nested if with plugin",
			input: []ast.BranchOrPlugin{
				ast.Branch{
					IfBlock: ast.IfBlock{
						Block: []ast.BranchOrPlugin{
							ast.Branch{
								IfBlock: ast.IfBlock{
									Block: []ast.BranchOrPlugin{
										ast.NewPlugin("plugin"),
									},
								},
							},
						},
					},
				},
			},

			want: []ast.BranchOrPlugin{
				ast.Branch{
					IfBlock: ast.IfBlock{
						Block: []ast.BranchOrPlugin{
							ast.Branch{
								IfBlock: ast.IfBlock{
									Block: []ast.BranchOrPlugin{
										ast.NewPlugin("replacement"),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			applyFunc := func(c *astutil.Cursor) {
				c.Replace(ast.NewPlugin("replacement"))
			}
			got := astutil.ApplyPlugins(test.input, applyFunc)

			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("Expect %#v to be equal to %#v", test.want, got)
			}
		})
	}
}
