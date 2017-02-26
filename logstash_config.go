package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/breml/logstash-config/ast"
)

var g = &grammar{
	rules: []*rule{
		{
			name: "config",
			pos:  position{line: 10, col: 1, offset: 252},
			expr: &actionExpr{
				pos: position{line: 11, col: 5, offset: 265},
				run: (*parser).callonconfig1,
				expr: &seqExpr{
					pos: position{line: 11, col: 5, offset: 265},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 11, col: 5, offset: 265},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 11, col: 7, offset: 267},
							label: "ps",
							expr: &ruleRefExpr{
								pos:  position{line: 11, col: 10, offset: 270},
								name: "plugin_section",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 11, col: 25, offset: 285},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 11, col: 27, offset: 287},
							label: "pss",
							expr: &zeroOrMoreExpr{
								pos: position{line: 11, col: 31, offset: 291},
								expr: &actionExpr{
									pos: position{line: 12, col: 9, offset: 301},
									run: (*parser).callonconfig9,
									expr: &seqExpr{
										pos: position{line: 12, col: 9, offset: 301},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 12, col: 9, offset: 301},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 12, col: 11, offset: 303},
												label: "ps",
												expr: &ruleRefExpr{
													pos:  position{line: 12, col: 14, offset: 306},
													name: "plugin_section",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 15, col: 8, offset: 369},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 15, col: 10, offset: 371},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "comment",
			pos:  position{line: 23, col: 1, offset: 521},
			expr: &oneOrMoreExpr{
				pos: position{line: 24, col: 5, offset: 535},
				expr: &seqExpr{
					pos: position{line: 24, col: 6, offset: 536},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 24, col: 6, offset: 536},
							expr: &ruleRefExpr{
								pos:  position{line: 24, col: 6, offset: 536},
								name: "whitespace",
							},
						},
						&litMatcher{
							pos:        position{line: 24, col: 18, offset: 548},
							val:        "#",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 24, col: 22, offset: 552},
							expr: &charClassMatcher{
								pos:        position{line: 24, col: 22, offset: 552},
								val:        "[^\\r\\n]",
								chars:      []rune{'\r', '\n'},
								ignoreCase: false,
								inverted:   true,
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 24, col: 31, offset: 561},
							expr: &litMatcher{
								pos:        position{line: 24, col: 31, offset: 561},
								val:        "\r",
								ignoreCase: false,
							},
						},
						&litMatcher{
							pos:        position{line: 24, col: 37, offset: 567},
							val:        "\n",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 30, col: 1, offset: 661},
			expr: &zeroOrMoreExpr{
				pos: position{line: 31, col: 5, offset: 669},
				expr: &choiceExpr{
					pos: position{line: 31, col: 6, offset: 670},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 31, col: 6, offset: 670},
							name: "comment",
						},
						&ruleRefExpr{
							pos:  position{line: 31, col: 16, offset: 680},
							name: "whitespace",
						},
					},
				},
			},
		},
		{
			name: "whitespace",
			pos:  position{line: 37, col: 1, offset: 776},
			expr: &oneOrMoreExpr{
				pos: position{line: 38, col: 5, offset: 793},
				expr: &charClassMatcher{
					pos:        position{line: 38, col: 5, offset: 793},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "plugin_section",
			pos:  position{line: 47, col: 1, offset: 950},
			expr: &actionExpr{
				pos: position{line: 48, col: 5, offset: 971},
				run: (*parser).callonplugin_section1,
				expr: &seqExpr{
					pos: position{line: 48, col: 5, offset: 971},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 48, col: 5, offset: 971},
							label: "pt",
							expr: &ruleRefExpr{
								pos:  position{line: 48, col: 8, offset: 974},
								name: "plugin_type",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 48, col: 20, offset: 986},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 48, col: 22, offset: 988},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 48, col: 26, offset: 992},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 48, col: 28, offset: 994},
							label: "bops",
							expr: &zeroOrMoreExpr{
								pos: position{line: 48, col: 33, offset: 999},
								expr: &actionExpr{
									pos: position{line: 49, col: 9, offset: 1009},
									run: (*parser).callonplugin_section10,
									expr: &seqExpr{
										pos: position{line: 49, col: 9, offset: 1009},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 49, col: 9, offset: 1009},
												label: "bop",
												expr: &ruleRefExpr{
													pos:  position{line: 49, col: 13, offset: 1013},
													name: "branch_or_plugin",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 49, col: 30, offset: 1030},
												name: "_",
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 52, col: 8, offset: 1082},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "branch_or_plugin",
			pos:  position{line: 60, col: 1, offset: 1196},
			expr: &choiceExpr{
				pos: position{line: 61, col: 5, offset: 1219},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 61, col: 5, offset: 1219},
						name: "branch",
					},
					&ruleRefExpr{
						pos:  position{line: 61, col: 14, offset: 1228},
						name: "plugin",
					},
				},
			},
		},
		{
			name: "plugin_type",
			pos:  position{line: 67, col: 1, offset: 1307},
			expr: &choiceExpr{
				pos: position{line: 68, col: 5, offset: 1325},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 68, col: 5, offset: 1325},
						run: (*parser).callonplugin_type2,
						expr: &litMatcher{
							pos:        position{line: 68, col: 5, offset: 1325},
							val:        "input",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 70, col: 9, offset: 1373},
						run: (*parser).callonplugin_type4,
						expr: &litMatcher{
							pos:        position{line: 70, col: 9, offset: 1373},
							val:        "filter",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 72, col: 9, offset: 1423},
						run: (*parser).callonplugin_type6,
						expr: &litMatcher{
							pos:        position{line: 72, col: 9, offset: 1423},
							val:        "output",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "plugin",
			pos:  position{line: 103, col: 1, offset: 2027},
			expr: &actionExpr{
				pos: position{line: 104, col: 5, offset: 2040},
				run: (*parser).callonplugin1,
				expr: &seqExpr{
					pos: position{line: 104, col: 5, offset: 2040},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 104, col: 5, offset: 2040},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 104, col: 10, offset: 2045},
								name: "name",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 104, col: 15, offset: 2050},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 104, col: 17, offset: 2052},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 104, col: 21, offset: 2056},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 104, col: 23, offset: 2058},
							label: "attributes",
							expr: &zeroOrOneExpr{
								pos: position{line: 104, col: 34, offset: 2069},
								expr: &actionExpr{
									pos: position{line: 105, col: 9, offset: 2080},
									run: (*parser).callonplugin10,
									expr: &seqExpr{
										pos: position{line: 105, col: 9, offset: 2080},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 105, col: 9, offset: 2080},
												label: "attribute",
												expr: &ruleRefExpr{
													pos:  position{line: 105, col: 19, offset: 2090},
													name: "attribute",
												},
											},
											&labeledExpr{
												pos:   position{line: 105, col: 29, offset: 2100},
												label: "attrs",
												expr: &zeroOrMoreExpr{
													pos: position{line: 105, col: 35, offset: 2106},
													expr: &actionExpr{
														pos: position{line: 106, col: 13, offset: 2120},
														run: (*parser).callonplugin16,
														expr: &seqExpr{
															pos: position{line: 106, col: 13, offset: 2120},
															exprs: []interface{}{
																&ruleRefExpr{
																	pos:  position{line: 106, col: 13, offset: 2120},
																	name: "whitespace",
																},
																&ruleRefExpr{
																	pos:  position{line: 106, col: 24, offset: 2131},
																	name: "_",
																},
																&labeledExpr{
																	pos:   position{line: 106, col: 26, offset: 2133},
																	label: "attribute",
																	expr: &ruleRefExpr{
																		pos:  position{line: 106, col: 36, offset: 2143},
																		name: "attribute",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 112, col: 8, offset: 2285},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 112, col: 10, offset: 2287},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "name",
			pos:  position{line: 123, col: 1, offset: 2457},
			expr: &choiceExpr{
				pos: position{line: 124, col: 7, offset: 2470},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 124, col: 7, offset: 2470},
						run: (*parser).callonname2,
						expr: &oneOrMoreExpr{
							pos: position{line: 124, col: 8, offset: 2471},
							expr: &charClassMatcher{
								pos:        position{line: 124, col: 8, offset: 2471},
								val:        "[A-Za-z0-9_-]",
								chars:      []rune{'_', '-'},
								ranges:     []rune{'A', 'Z', 'a', 'z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
					&actionExpr{
						pos: position{line: 126, col: 9, offset: 2532},
						run: (*parser).callonname5,
						expr: &labeledExpr{
							pos:   position{line: 126, col: 9, offset: 2532},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 126, col: 15, offset: 2538},
								name: "string_value",
							},
						},
					},
				},
			},
		},
		{
			name: "attribute",
			pos:  position{line: 135, col: 1, offset: 2687},
			expr: &actionExpr{
				pos: position{line: 136, col: 5, offset: 2703},
				run: (*parser).callonattribute1,
				expr: &seqExpr{
					pos: position{line: 136, col: 5, offset: 2703},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 136, col: 5, offset: 2703},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 136, col: 10, offset: 2708},
								name: "name",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 136, col: 15, offset: 2713},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 136, col: 17, offset: 2715},
							val:        "=>",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 136, col: 22, offset: 2720},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 136, col: 24, offset: 2722},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 136, col: 30, offset: 2728},
								name: "value",
							},
						},
					},
				},
			},
		},
		{
			name: "value",
			pos:  position{line: 144, col: 1, offset: 2865},
			expr: &choiceExpr{
				pos: position{line: 145, col: 5, offset: 2877},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 145, col: 5, offset: 2877},
						name: "plugin",
					},
					&ruleRefExpr{
						pos:  position{line: 145, col: 14, offset: 2886},
						name: "bareword",
					},
					&ruleRefExpr{
						pos:  position{line: 145, col: 25, offset: 2897},
						name: "string_value",
					},
					&ruleRefExpr{
						pos:  position{line: 145, col: 40, offset: 2912},
						name: "number",
					},
					&ruleRefExpr{
						pos:  position{line: 145, col: 49, offset: 2921},
						name: "array",
					},
					&ruleRefExpr{
						pos:  position{line: 145, col: 57, offset: 2929},
						name: "hash",
					},
				},
			},
		},
		{
			name: "array_value",
			pos:  position{line: 151, col: 1, offset: 3016},
			expr: &choiceExpr{
				pos: position{line: 152, col: 5, offset: 3034},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 152, col: 5, offset: 3034},
						name: "bareword",
					},
					&ruleRefExpr{
						pos:  position{line: 152, col: 16, offset: 3045},
						name: "string_value",
					},
					&ruleRefExpr{
						pos:  position{line: 152, col: 31, offset: 3060},
						name: "number",
					},
					&ruleRefExpr{
						pos:  position{line: 152, col: 40, offset: 3069},
						name: "array",
					},
					&ruleRefExpr{
						pos:  position{line: 152, col: 48, offset: 3077},
						name: "hash",
					},
				},
			},
		},
		{
			name: "bareword",
			pos:  position{line: 159, col: 1, offset: 3184},
			expr: &actionExpr{
				pos: position{line: 160, col: 5, offset: 3199},
				run: (*parser).callonbareword1,
				expr: &seqExpr{
					pos: position{line: 160, col: 5, offset: 3199},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 160, col: 5, offset: 3199},
							val:        "[A-Za-z_]",
							chars:      []rune{'_'},
							ranges:     []rune{'A', 'Z', 'a', 'z'},
							ignoreCase: false,
							inverted:   false,
						},
						&oneOrMoreExpr{
							pos: position{line: 160, col: 15, offset: 3209},
							expr: &charClassMatcher{
								pos:        position{line: 160, col: 15, offset: 3209},
								val:        "[A-Za-z0-9_]",
								chars:      []rune{'_'},
								ranges:     []rune{'A', 'Z', 'a', 'z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "double_quoted_string",
			pos:  position{line: 168, col: 1, offset: 3419},
			expr: &actionExpr{
				pos: position{line: 169, col: 5, offset: 3446},
				run: (*parser).callondouble_quoted_string1,
				expr: &seqExpr{
					pos: position{line: 169, col: 7, offset: 3448},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 169, col: 7, offset: 3448},
							val:        "\"",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 169, col: 11, offset: 3452},
							expr: &choiceExpr{
								pos: position{line: 169, col: 13, offset: 3454},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 169, col: 13, offset: 3454},
										val:        "\\\"",
										ignoreCase: false,
									},
									&seqExpr{
										pos: position{line: 169, col: 20, offset: 3461},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 169, col: 20, offset: 3461},
												expr: &litMatcher{
													pos:        position{line: 169, col: 21, offset: 3462},
													val:        "\"",
													ignoreCase: false,
												},
											},
											&anyMatcher{
												line: 169, col: 25, offset: 3466,
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 169, col: 30, offset: 3471},
							val:        "\"",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "single_quoted_string",
			pos:  position{line: 177, col: 1, offset: 3629},
			expr: &actionExpr{
				pos: position{line: 178, col: 5, offset: 3656},
				run: (*parser).callonsingle_quoted_string1,
				expr: &seqExpr{
					pos: position{line: 178, col: 7, offset: 3658},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 178, col: 7, offset: 3658},
							val:        "'",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 178, col: 11, offset: 3662},
							expr: &choiceExpr{
								pos: position{line: 178, col: 13, offset: 3664},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 178, col: 13, offset: 3664},
										val:        "\\'",
										ignoreCase: false,
									},
									&seqExpr{
										pos: position{line: 178, col: 20, offset: 3671},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 178, col: 20, offset: 3671},
												expr: &litMatcher{
													pos:        position{line: 178, col: 21, offset: 3672},
													val:        "'",
													ignoreCase: false,
												},
											},
											&anyMatcher{
												line: 178, col: 25, offset: 3676,
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 178, col: 30, offset: 3681},
							val:        "'",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "string_value",
			pos:  position{line: 186, col: 1, offset: 3806},
			expr: &actionExpr{
				pos: position{line: 187, col: 5, offset: 3825},
				run: (*parser).callonstring_value1,
				expr: &labeledExpr{
					pos:   position{line: 187, col: 5, offset: 3825},
					label: "str",
					expr: &choiceExpr{
						pos: position{line: 187, col: 11, offset: 3831},
						alternatives: []interface{}{
							&actionExpr{
								pos: position{line: 187, col: 11, offset: 3831},
								run: (*parser).callonstring_value4,
								expr: &labeledExpr{
									pos:   position{line: 187, col: 11, offset: 3831},
									label: "str",
									expr: &ruleRefExpr{
										pos:  position{line: 187, col: 15, offset: 3835},
										name: "double_quoted_string",
									},
								},
							},
							&actionExpr{
								pos: position{line: 189, col: 9, offset: 3945},
								run: (*parser).callonstring_value7,
								expr: &labeledExpr{
									pos:   position{line: 189, col: 9, offset: 3945},
									label: "str",
									expr: &ruleRefExpr{
										pos:  position{line: 189, col: 13, offset: 3949},
										name: "single_quoted_string",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "regexp",
			pos:  position{line: 199, col: 1, offset: 4189},
			expr: &actionExpr{
				pos: position{line: 200, col: 5, offset: 4202},
				run: (*parser).callonregexp1,
				expr: &seqExpr{
					pos: position{line: 200, col: 7, offset: 4204},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 200, col: 7, offset: 4204},
							val:        "/",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 200, col: 11, offset: 4208},
							expr: &choiceExpr{
								pos: position{line: 200, col: 13, offset: 4210},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 200, col: 13, offset: 4210},
										val:        "\\/",
										ignoreCase: false,
									},
									&seqExpr{
										pos: position{line: 200, col: 20, offset: 4217},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 200, col: 20, offset: 4217},
												expr: &litMatcher{
													pos:        position{line: 200, col: 21, offset: 4218},
													val:        "/",
													ignoreCase: false,
												},
											},
											&anyMatcher{
												line: 200, col: 25, offset: 4222,
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 200, col: 30, offset: 4227},
							val:        "/",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "number",
			pos:  position{line: 209, col: 1, offset: 4365},
			expr: &actionExpr{
				pos: position{line: 210, col: 5, offset: 4378},
				run: (*parser).callonnumber1,
				expr: &seqExpr{
					pos: position{line: 210, col: 5, offset: 4378},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 210, col: 5, offset: 4378},
							expr: &litMatcher{
								pos:        position{line: 210, col: 5, offset: 4378},
								val:        "-",
								ignoreCase: false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 210, col: 10, offset: 4383},
							expr: &charClassMatcher{
								pos:        position{line: 210, col: 10, offset: 4383},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 210, col: 17, offset: 4390},
							expr: &seqExpr{
								pos: position{line: 210, col: 18, offset: 4391},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 210, col: 18, offset: 4391},
										val:        ".",
										ignoreCase: false,
									},
									&zeroOrMoreExpr{
										pos: position{line: 210, col: 22, offset: 4395},
										expr: &charClassMatcher{
											pos:        position{line: 210, col: 22, offset: 4395},
											val:        "[0-9]",
											ranges:     []rune{'0', '9'},
											ignoreCase: false,
											inverted:   false,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "array",
			pos:  position{line: 226, col: 1, offset: 4712},
			expr: &actionExpr{
				pos: position{line: 227, col: 5, offset: 4724},
				run: (*parser).callonarray1,
				expr: &seqExpr{
					pos: position{line: 227, col: 5, offset: 4724},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 227, col: 5, offset: 4724},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 227, col: 9, offset: 4728},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 227, col: 11, offset: 4730},
							label: "value",
							expr: &zeroOrOneExpr{
								pos: position{line: 227, col: 17, offset: 4736},
								expr: &actionExpr{
									pos: position{line: 228, col: 9, offset: 4747},
									run: (*parser).callonarray7,
									expr: &seqExpr{
										pos: position{line: 228, col: 9, offset: 4747},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 228, col: 9, offset: 4747},
												label: "value",
												expr: &ruleRefExpr{
													pos:  position{line: 228, col: 15, offset: 4753},
													name: "value",
												},
											},
											&labeledExpr{
												pos:   position{line: 228, col: 21, offset: 4759},
												label: "values",
												expr: &zeroOrMoreExpr{
													pos: position{line: 228, col: 28, offset: 4766},
													expr: &actionExpr{
														pos: position{line: 229, col: 13, offset: 4780},
														run: (*parser).callonarray13,
														expr: &seqExpr{
															pos: position{line: 229, col: 13, offset: 4780},
															exprs: []interface{}{
																&ruleRefExpr{
																	pos:  position{line: 229, col: 13, offset: 4780},
																	name: "_",
																},
																&litMatcher{
																	pos:        position{line: 229, col: 15, offset: 4782},
																	val:        ",",
																	ignoreCase: false,
																},
																&ruleRefExpr{
																	pos:  position{line: 229, col: 19, offset: 4786},
																	name: "_",
																},
																&labeledExpr{
																	pos:   position{line: 229, col: 21, offset: 4788},
																	label: "value",
																	expr: &ruleRefExpr{
																		pos:  position{line: 229, col: 27, offset: 4794},
																		name: "value",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 235, col: 8, offset: 4925},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 235, col: 10, offset: 4927},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "hash",
			pos:  position{line: 248, col: 1, offset: 5096},
			expr: &actionExpr{
				pos: position{line: 249, col: 5, offset: 5107},
				run: (*parser).callonhash1,
				expr: &seqExpr{
					pos: position{line: 249, col: 5, offset: 5107},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 249, col: 5, offset: 5107},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 249, col: 9, offset: 5111},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 249, col: 11, offset: 5113},
							label: "entries",
							expr: &zeroOrOneExpr{
								pos: position{line: 249, col: 19, offset: 5121},
								expr: &ruleRefExpr{
									pos:  position{line: 249, col: 19, offset: 5121},
									name: "hashentries",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 249, col: 32, offset: 5134},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 249, col: 34, offset: 5136},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "hashentries",
			pos:  position{line: 258, col: 1, offset: 5295},
			expr: &actionExpr{
				pos: position{line: 259, col: 5, offset: 5313},
				run: (*parser).callonhashentries1,
				expr: &seqExpr{
					pos: position{line: 259, col: 5, offset: 5313},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 259, col: 5, offset: 5313},
							label: "hashentry",
							expr: &ruleRefExpr{
								pos:  position{line: 259, col: 15, offset: 5323},
								name: "hashentry",
							},
						},
						&labeledExpr{
							pos:   position{line: 259, col: 25, offset: 5333},
							label: "hashentries1",
							expr: &zeroOrMoreExpr{
								pos: position{line: 259, col: 38, offset: 5346},
								expr: &actionExpr{
									pos: position{line: 260, col: 9, offset: 5356},
									run: (*parser).callonhashentries7,
									expr: &seqExpr{
										pos: position{line: 260, col: 9, offset: 5356},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 260, col: 9, offset: 5356},
												name: "whitespace",
											},
											&labeledExpr{
												pos:   position{line: 260, col: 20, offset: 5367},
												label: "hashentry",
												expr: &ruleRefExpr{
													pos:  position{line: 260, col: 30, offset: 5377},
													name: "hashentry",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "hashentry",
			pos:  position{line: 272, col: 1, offset: 5629},
			expr: &actionExpr{
				pos: position{line: 273, col: 5, offset: 5645},
				run: (*parser).callonhashentry1,
				expr: &seqExpr{
					pos: position{line: 273, col: 5, offset: 5645},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 273, col: 5, offset: 5645},
							label: "name",
							expr: &choiceExpr{
								pos: position{line: 273, col: 11, offset: 5651},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 273, col: 11, offset: 5651},
										name: "number",
									},
									&ruleRefExpr{
										pos:  position{line: 273, col: 20, offset: 5660},
										name: "bareword",
									},
									&ruleRefExpr{
										pos:  position{line: 273, col: 31, offset: 5671},
										name: "string_value",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 273, col: 45, offset: 5685},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 273, col: 47, offset: 5687},
							val:        "=>",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 273, col: 52, offset: 5692},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 273, col: 54, offset: 5694},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 273, col: 60, offset: 5700},
								name: "value",
							},
						},
					},
				},
			},
		},
		{
			name: "branch",
			pos:  position{line: 284, col: 1, offset: 5867},
			expr: &actionExpr{
				pos: position{line: 285, col: 5, offset: 5880},
				run: (*parser).callonbranch1,
				expr: &seqExpr{
					pos: position{line: 285, col: 5, offset: 5880},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 285, col: 5, offset: 5880},
							label: "ifBlock",
							expr: &ruleRefExpr{
								pos:  position{line: 285, col: 13, offset: 5888},
								name: "if_cond",
							},
						},
						&labeledExpr{
							pos:   position{line: 285, col: 21, offset: 5896},
							label: "elseIfBlocks",
							expr: &zeroOrMoreExpr{
								pos: position{line: 285, col: 34, offset: 5909},
								expr: &actionExpr{
									pos: position{line: 286, col: 9, offset: 5919},
									run: (*parser).callonbranch7,
									expr: &seqExpr{
										pos: position{line: 286, col: 9, offset: 5919},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 286, col: 9, offset: 5919},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 286, col: 11, offset: 5921},
												label: "eib",
												expr: &ruleRefExpr{
													pos:  position{line: 286, col: 15, offset: 5925},
													name: "else_if",
												},
											},
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 289, col: 12, offset: 5992},
							label: "elseBlock",
							expr: &zeroOrOneExpr{
								pos: position{line: 289, col: 22, offset: 6002},
								expr: &actionExpr{
									pos: position{line: 290, col: 13, offset: 6016},
									run: (*parser).callonbranch14,
									expr: &seqExpr{
										pos: position{line: 290, col: 13, offset: 6016},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 290, col: 13, offset: 6016},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 290, col: 15, offset: 6018},
												label: "eb",
												expr: &ruleRefExpr{
													pos:  position{line: 290, col: 18, offset: 6021},
													name: "else_cond",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "if_cond",
			pos:  position{line: 302, col: 1, offset: 6270},
			expr: &actionExpr{
				pos: position{line: 303, col: 5, offset: 6284},
				run: (*parser).callonif_cond1,
				expr: &seqExpr{
					pos: position{line: 303, col: 5, offset: 6284},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 303, col: 5, offset: 6284},
							val:        "if",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 303, col: 10, offset: 6289},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 303, col: 12, offset: 6291},
							label: "cond",
							expr: &ruleRefExpr{
								pos:  position{line: 303, col: 17, offset: 6296},
								name: "condition",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 303, col: 27, offset: 6306},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 303, col: 29, offset: 6308},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 303, col: 33, offset: 6312},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 303, col: 35, offset: 6314},
							label: "bops",
							expr: &zeroOrMoreExpr{
								pos: position{line: 303, col: 40, offset: 6319},
								expr: &actionExpr{
									pos: position{line: 304, col: 13, offset: 6333},
									run: (*parser).callonif_cond12,
									expr: &seqExpr{
										pos: position{line: 304, col: 13, offset: 6333},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 304, col: 13, offset: 6333},
												label: "bop",
												expr: &ruleRefExpr{
													pos:  position{line: 304, col: 17, offset: 6337},
													name: "branch_or_plugin",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 304, col: 34, offset: 6354},
												name: "_",
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 307, col: 12, offset: 6415},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "else_if",
			pos:  position{line: 316, col: 1, offset: 6597},
			expr: &actionExpr{
				pos: position{line: 317, col: 5, offset: 6611},
				run: (*parser).callonelse_if1,
				expr: &seqExpr{
					pos: position{line: 317, col: 5, offset: 6611},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 317, col: 5, offset: 6611},
							val:        "else",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 317, col: 12, offset: 6618},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 317, col: 14, offset: 6620},
							val:        "if",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 317, col: 19, offset: 6625},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 317, col: 21, offset: 6627},
							label: "cond",
							expr: &ruleRefExpr{
								pos:  position{line: 317, col: 26, offset: 6632},
								name: "condition",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 317, col: 36, offset: 6642},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 317, col: 38, offset: 6644},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 317, col: 42, offset: 6648},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 317, col: 44, offset: 6650},
							label: "bops",
							expr: &zeroOrMoreExpr{
								pos: position{line: 317, col: 49, offset: 6655},
								expr: &actionExpr{
									pos: position{line: 318, col: 9, offset: 6665},
									run: (*parser).callonelse_if14,
									expr: &seqExpr{
										pos: position{line: 318, col: 9, offset: 6665},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 318, col: 9, offset: 6665},
												label: "bop",
												expr: &ruleRefExpr{
													pos:  position{line: 318, col: 13, offset: 6669},
													name: "branch_or_plugin",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 318, col: 30, offset: 6686},
												name: "_",
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 321, col: 8, offset: 6735},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "else_cond",
			pos:  position{line: 330, col: 1, offset: 6897},
			expr: &actionExpr{
				pos: position{line: 331, col: 5, offset: 6913},
				run: (*parser).callonelse_cond1,
				expr: &seqExpr{
					pos: position{line: 331, col: 5, offset: 6913},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 331, col: 5, offset: 6913},
							val:        "else",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 331, col: 12, offset: 6920},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 331, col: 14, offset: 6922},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 331, col: 18, offset: 6926},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 331, col: 20, offset: 6928},
							label: "bops",
							expr: &zeroOrMoreExpr{
								pos: position{line: 331, col: 25, offset: 6933},
								expr: &actionExpr{
									pos: position{line: 332, col: 9, offset: 6943},
									run: (*parser).callonelse_cond9,
									expr: &seqExpr{
										pos: position{line: 332, col: 9, offset: 6943},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 332, col: 9, offset: 6943},
												label: "bop",
												expr: &ruleRefExpr{
													pos:  position{line: 332, col: 13, offset: 6947},
													name: "branch_or_plugin",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 332, col: 30, offset: 6964},
												name: "_",
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 335, col: 8, offset: 7013},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "condition",
			pos:  position{line: 344, col: 1, offset: 7182},
			expr: &actionExpr{
				pos: position{line: 345, col: 5, offset: 7198},
				run: (*parser).calloncondition1,
				expr: &seqExpr{
					pos: position{line: 345, col: 5, offset: 7198},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 345, col: 5, offset: 7198},
							label: "cond",
							expr: &ruleRefExpr{
								pos:  position{line: 345, col: 10, offset: 7203},
								name: "expression",
							},
						},
						&labeledExpr{
							pos:   position{line: 345, col: 21, offset: 7214},
							label: "conds",
							expr: &zeroOrMoreExpr{
								pos: position{line: 345, col: 27, offset: 7220},
								expr: &actionExpr{
									pos: position{line: 346, col: 9, offset: 7230},
									run: (*parser).calloncondition7,
									expr: &seqExpr{
										pos: position{line: 346, col: 9, offset: 7230},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 346, col: 9, offset: 7230},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 346, col: 11, offset: 7232},
												label: "bo",
												expr: &ruleRefExpr{
													pos:  position{line: 346, col: 14, offset: 7235},
													name: "boolean_operator",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 346, col: 31, offset: 7252},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 346, col: 33, offset: 7254},
												label: "cond",
												expr: &ruleRefExpr{
													pos:  position{line: 346, col: 38, offset: 7259},
													name: "expression",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "expression",
			pos:  position{line: 365, col: 1, offset: 7658},
			expr: &choiceExpr{
				pos: position{line: 367, col: 9, offset: 7686},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 368, col: 13, offset: 7700},
						run: (*parser).callonexpression2,
						expr: &seqExpr{
							pos: position{line: 368, col: 13, offset: 7700},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 368, col: 13, offset: 7700},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 368, col: 17, offset: 7704},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 368, col: 19, offset: 7706},
									label: "cond",
									expr: &ruleRefExpr{
										pos:  position{line: 368, col: 24, offset: 7711},
										name: "condition",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 368, col: 34, offset: 7721},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 368, col: 36, offset: 7723},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 372, col: 9, offset: 7811},
						name: "negative_expression",
					},
					&ruleRefExpr{
						pos:  position{line: 373, col: 9, offset: 7839},
						name: "in_expression",
					},
					&ruleRefExpr{
						pos:  position{line: 374, col: 9, offset: 7861},
						name: "not_in_expression",
					},
					&ruleRefExpr{
						pos:  position{line: 375, col: 9, offset: 7887},
						name: "compare_expression",
					},
					&ruleRefExpr{
						pos:  position{line: 376, col: 9, offset: 7914},
						name: "regexp_expression",
					},
					&actionExpr{
						pos: position{line: 377, col: 9, offset: 7940},
						run: (*parser).callonexpression15,
						expr: &labeledExpr{
							pos:   position{line: 377, col: 9, offset: 7940},
							label: "rv",
							expr: &ruleRefExpr{
								pos:  position{line: 377, col: 12, offset: 7943},
								name: "rvalue",
							},
						},
					},
				},
			},
		},
		{
			name: "negative_expression",
			pos:  position{line: 389, col: 1, offset: 8166},
			expr: &choiceExpr{
				pos: position{line: 391, col: 9, offset: 8203},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 392, col: 13, offset: 8217},
						run: (*parser).callonnegative_expression2,
						expr: &seqExpr{
							pos: position{line: 392, col: 13, offset: 8217},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 392, col: 13, offset: 8217},
									val:        "!",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 392, col: 17, offset: 8221},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 392, col: 19, offset: 8223},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 392, col: 23, offset: 8227},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 392, col: 25, offset: 8229},
									label: "cond",
									expr: &ruleRefExpr{
										pos:  position{line: 392, col: 30, offset: 8234},
										name: "condition",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 392, col: 40, offset: 8244},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 392, col: 42, offset: 8246},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 397, col: 11, offset: 8345},
						run: (*parser).callonnegative_expression12,
						expr: &seqExpr{
							pos: position{line: 397, col: 11, offset: 8345},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 397, col: 11, offset: 8345},
									val:        "!",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 397, col: 15, offset: 8349},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 397, col: 17, offset: 8351},
									label: "sel",
									expr: &ruleRefExpr{
										pos:  position{line: 397, col: 21, offset: 8355},
										name: "selector",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "in_expression",
			pos:  position{line: 408, col: 1, offset: 8555},
			expr: &actionExpr{
				pos: position{line: 409, col: 5, offset: 8575},
				run: (*parser).callonin_expression1,
				expr: &seqExpr{
					pos: position{line: 409, col: 5, offset: 8575},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 409, col: 5, offset: 8575},
							label: "lv",
							expr: &ruleRefExpr{
								pos:  position{line: 409, col: 8, offset: 8578},
								name: "rvalue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 409, col: 15, offset: 8585},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 409, col: 17, offset: 8587},
							name: "in_operator",
						},
						&ruleRefExpr{
							pos:  position{line: 409, col: 29, offset: 8599},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 409, col: 31, offset: 8601},
							label: "rv",
							expr: &ruleRefExpr{
								pos:  position{line: 409, col: 34, offset: 8604},
								name: "rvalue",
							},
						},
					},
				},
			},
		},
		{
			name: "not_in_expression",
			pos:  position{line: 418, col: 1, offset: 8784},
			expr: &actionExpr{
				pos: position{line: 419, col: 5, offset: 8808},
				run: (*parser).callonnot_in_expression1,
				expr: &seqExpr{
					pos: position{line: 419, col: 5, offset: 8808},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 419, col: 5, offset: 8808},
							label: "lv",
							expr: &ruleRefExpr{
								pos:  position{line: 419, col: 8, offset: 8811},
								name: "rvalue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 419, col: 15, offset: 8818},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 419, col: 17, offset: 8820},
							name: "not_in_operator",
						},
						&ruleRefExpr{
							pos:  position{line: 419, col: 33, offset: 8836},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 419, col: 35, offset: 8838},
							label: "rv",
							expr: &ruleRefExpr{
								pos:  position{line: 419, col: 38, offset: 8841},
								name: "rvalue",
							},
						},
					},
				},
			},
		},
		{
			name: "in_operator",
			pos:  position{line: 427, col: 1, offset: 8942},
			expr: &litMatcher{
				pos:        position{line: 428, col: 5, offset: 8960},
				val:        "in",
				ignoreCase: false,
			},
		},
		{
			name: "not_in_operator",
			pos:  position{line: 434, col: 1, offset: 9023},
			expr: &seqExpr{
				pos: position{line: 435, col: 5, offset: 9045},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 435, col: 5, offset: 9045},
						val:        "not ",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 435, col: 12, offset: 9052},
						name: "_",
					},
					&litMatcher{
						pos:        position{line: 435, col: 14, offset: 9054},
						val:        "in",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "rvalue",
			pos:  position{line: 444, col: 1, offset: 9313},
			expr: &choiceExpr{
				pos: position{line: 445, col: 5, offset: 9326},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 445, col: 5, offset: 9326},
						name: "string_value",
					},
					&ruleRefExpr{
						pos:  position{line: 445, col: 20, offset: 9341},
						name: "number",
					},
					&ruleRefExpr{
						pos:  position{line: 445, col: 29, offset: 9350},
						name: "selector",
					},
					&ruleRefExpr{
						pos:  position{line: 445, col: 40, offset: 9361},
						name: "array",
					},
					&ruleRefExpr{
						pos:  position{line: 445, col: 48, offset: 9369},
						name: "regexp",
					},
				},
			},
		},
		{
			name: "compare_expression",
			pos:  position{line: 477, col: 1, offset: 10045},
			expr: &actionExpr{
				pos: position{line: 478, col: 5, offset: 10070},
				run: (*parser).calloncompare_expression1,
				expr: &seqExpr{
					pos: position{line: 478, col: 5, offset: 10070},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 478, col: 5, offset: 10070},
							label: "lv",
							expr: &ruleRefExpr{
								pos:  position{line: 478, col: 8, offset: 10073},
								name: "rvalue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 478, col: 15, offset: 10080},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 478, col: 17, offset: 10082},
							label: "co",
							expr: &ruleRefExpr{
								pos:  position{line: 478, col: 20, offset: 10085},
								name: "compare_operator",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 478, col: 37, offset: 10102},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 478, col: 39, offset: 10104},
							label: "rv",
							expr: &ruleRefExpr{
								pos:  position{line: 478, col: 42, offset: 10107},
								name: "rvalue",
							},
						},
					},
				},
			},
		},
		{
			name: "compare_operator",
			pos:  position{line: 487, col: 1, offset: 10306},
			expr: &actionExpr{
				pos: position{line: 488, col: 5, offset: 10329},
				run: (*parser).calloncompare_operator1,
				expr: &choiceExpr{
					pos: position{line: 488, col: 6, offset: 10330},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 488, col: 6, offset: 10330},
							val:        "==",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 488, col: 13, offset: 10337},
							val:        "!=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 488, col: 20, offset: 10344},
							val:        "<=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 488, col: 27, offset: 10351},
							val:        ">=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 488, col: 34, offset: 10358},
							val:        "<",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 488, col: 40, offset: 10364},
							val:        ">",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "regexp_expression",
			pos:  position{line: 497, col: 1, offset: 10566},
			expr: &actionExpr{
				pos: position{line: 498, col: 5, offset: 10590},
				run: (*parser).callonregexp_expression1,
				expr: &seqExpr{
					pos: position{line: 498, col: 5, offset: 10590},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 498, col: 5, offset: 10590},
							label: "lv",
							expr: &ruleRefExpr{
								pos:  position{line: 498, col: 8, offset: 10593},
								name: "rvalue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 498, col: 15, offset: 10600},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 498, col: 18, offset: 10603},
							label: "ro",
							expr: &ruleRefExpr{
								pos:  position{line: 498, col: 21, offset: 10606},
								name: "regexp_operator",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 498, col: 37, offset: 10622},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 498, col: 39, offset: 10624},
							label: "rv",
							expr: &choiceExpr{
								pos: position{line: 498, col: 43, offset: 10628},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 498, col: 43, offset: 10628},
										name: "string_value",
									},
									&ruleRefExpr{
										pos:  position{line: 498, col: 58, offset: 10643},
										name: "regexp",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "regexp_operator",
			pos:  position{line: 506, col: 1, offset: 10802},
			expr: &actionExpr{
				pos: position{line: 507, col: 5, offset: 10824},
				run: (*parser).callonregexp_operator1,
				expr: &choiceExpr{
					pos: position{line: 507, col: 6, offset: 10825},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 507, col: 6, offset: 10825},
							val:        "=~",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 507, col: 13, offset: 10832},
							val:        "!~",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "boolean_operator",
			pos:  position{line: 516, col: 1, offset: 11018},
			expr: &actionExpr{
				pos: position{line: 517, col: 5, offset: 11041},
				run: (*parser).callonboolean_operator1,
				expr: &choiceExpr{
					pos: position{line: 517, col: 6, offset: 11042},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 517, col: 6, offset: 11042},
							val:        "and",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 517, col: 14, offset: 11050},
							val:        "or",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 517, col: 21, offset: 11057},
							val:        "xor",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 517, col: 29, offset: 11065},
							val:        "nand",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "selector",
			pos:  position{line: 526, col: 1, offset: 11225},
			expr: &actionExpr{
				pos: position{line: 527, col: 5, offset: 11240},
				run: (*parser).callonselector1,
				expr: &labeledExpr{
					pos:   position{line: 527, col: 5, offset: 11240},
					label: "ses",
					expr: &oneOrMoreExpr{
						pos: position{line: 527, col: 9, offset: 11244},
						expr: &ruleRefExpr{
							pos:  position{line: 527, col: 9, offset: 11244},
							name: "selector_element",
						},
					},
				},
			},
		},
		{
			name: "selector_element",
			pos:  position{line: 536, col: 1, offset: 11408},
			expr: &actionExpr{
				pos: position{line: 537, col: 5, offset: 11431},
				run: (*parser).callonselector_element1,
				expr: &seqExpr{
					pos: position{line: 537, col: 5, offset: 11431},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 537, col: 5, offset: 11431},
							val:        "[",
							ignoreCase: false,
						},
						&oneOrMoreExpr{
							pos: position{line: 537, col: 9, offset: 11435},
							expr: &charClassMatcher{
								pos:        position{line: 537, col: 9, offset: 11435},
								val:        "[^\\],]",
								chars:      []rune{']', ','},
								ignoreCase: false,
								inverted:   true,
							},
						},
						&litMatcher{
							pos:        position{line: 537, col: 17, offset: 11443},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 541, col: 1, offset: 11504},
			expr: &notExpr{
				pos: position{line: 541, col: 7, offset: 11510},
				expr: &anyMatcher{
					line: 541, col: 8, offset: 11511,
				},
			},
		},
	},
}

func (c *current) onconfig9(ps interface{}) (interface{}, error) {

	return ret(ps)

}

func (p *parser) callonconfig9() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onconfig9(stack["ps"])
}

func (c *current) onconfig1(ps, pss interface{}) (interface{}, error) {

	return config(ps, pss)

}

func (p *parser) callonconfig1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onconfig1(stack["ps"], stack["pss"])
}

func (c *current) onplugin_section10(bop interface{}) (interface{}, error) {

	return ret(bop)

}

func (p *parser) callonplugin_section10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onplugin_section10(stack["bop"])
}

func (c *current) onplugin_section1(pt, bops interface{}) (interface{}, error) {

	return pluginSection(pt, bops)

}

func (p *parser) callonplugin_section1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onplugin_section1(stack["pt"], stack["bops"])
}

func (c *current) onplugin_type2() (interface{}, error) {
	return ast.Input, nil

}

func (p *parser) callonplugin_type2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onplugin_type2()
}

func (c *current) onplugin_type4() (interface{}, error) {
	return ast.Filter, nil

}

func (p *parser) callonplugin_type4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onplugin_type4()
}

func (c *current) onplugin_type6() (interface{}, error) {
	return ast.Output, nil

}

func (p *parser) callonplugin_type6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onplugin_type6()
}

func (c *current) onplugin16(attribute interface{}) (interface{}, error) {
	return ret(attribute)

}

func (p *parser) callonplugin16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onplugin16(stack["attribute"])
}

func (c *current) onplugin10(attribute, attrs interface{}) (interface{}, error) {
	return attributes(attribute, attrs)

}

func (p *parser) callonplugin10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onplugin10(stack["attribute"], stack["attrs"])
}

func (c *current) onplugin1(name, attributes interface{}) (interface{}, error) {
	return plugin(name, attributes)

}

func (p *parser) callonplugin1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onplugin1(stack["name"], stack["attributes"])
}

func (c *current) onname2() (interface{}, error) {
	return string(c.text), nil

}

func (p *parser) callonname2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onname2()
}

func (c *current) onname5(value interface{}) (interface{}, error) {
	return ret(value)

}

func (p *parser) callonname5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onname5(stack["value"])
}

func (c *current) onattribute1(name, value interface{}) (interface{}, error) {
	return attribute(name, value)

}

func (p *parser) callonattribute1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onattribute1(stack["name"], stack["value"])
}

func (c *current) onbareword1() (interface{}, error) {
	return ast.NewStringAttribute("", string(c.text), ast.Bareword), nil

}

func (p *parser) callonbareword1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onbareword1()
}

func (c *current) ondouble_quoted_string1() (interface{}, error) {
	return enclosedValue(c)

}

func (p *parser) callondouble_quoted_string1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.ondouble_quoted_string1()
}

func (c *current) onsingle_quoted_string1() (interface{}, error) {
	return enclosedValue(c)

}

func (p *parser) callonsingle_quoted_string1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onsingle_quoted_string1()
}

func (c *current) onstring_value4(str interface{}) (interface{}, error) {
	return ast.NewStringAttribute("", str.(string), ast.DoubleQuoted), nil

}

func (p *parser) callonstring_value4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onstring_value4(stack["str"])
}

func (c *current) onstring_value7(str interface{}) (interface{}, error) {
	return ast.NewStringAttribute("", str.(string), ast.SingleQuoted), nil

}

func (p *parser) callonstring_value7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onstring_value7(stack["str"])
}

func (c *current) onstring_value1(str interface{}) (interface{}, error) {
	return ret(str)

}

func (p *parser) callonstring_value1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onstring_value1(stack["str"])
}

func (c *current) onregexp1() (interface{}, error) {
	return regexp(c)

}

func (p *parser) callonregexp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onregexp1()
}

func (c *current) onnumber1() (interface{}, error) {
	return number(string(c.text))

}

func (p *parser) callonnumber1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onnumber1()
}

func (c *current) onarray13(value interface{}) (interface{}, error) {
	return ret(value)

}

func (p *parser) callonarray13() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onarray13(stack["value"])
}

func (c *current) onarray7(value, values interface{}) (interface{}, error) {
	return attributes(value, values)

}

func (p *parser) callonarray7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onarray7(stack["value"], stack["values"])
}

func (c *current) onarray1(value interface{}) (interface{}, error) {
	return array(value)

}

func (p *parser) callonarray1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onarray1(stack["value"])
}

func (c *current) onhash1(entries interface{}) (interface{}, error) {
	return hash(entries)

}

func (p *parser) callonhash1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onhash1(stack["entries"])
}

func (c *current) onhashentries7(hashentry interface{}) (interface{}, error) {
	return ret(hashentry)

}

func (p *parser) callonhashentries7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onhashentries7(stack["hashentry"])
}

func (c *current) onhashentries1(hashentry, hashentries1 interface{}) (interface{}, error) {
	return hashentries(hashentry, hashentries1)

}

func (p *parser) callonhashentries1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onhashentries1(stack["hashentry"], stack["hashentries1"])
}

func (c *current) onhashentry1(name, value interface{}) (interface{}, error) {
	return hashentry(name, value)

}

func (p *parser) callonhashentry1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onhashentry1(stack["name"], stack["value"])
}

func (c *current) onbranch7(eib interface{}) (interface{}, error) {
	return ret(eib)

}

func (p *parser) callonbranch7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onbranch7(stack["eib"])
}

func (c *current) onbranch14(eb interface{}) (interface{}, error) {
	return ret(eb)

}

func (p *parser) callonbranch14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onbranch14(stack["eb"])
}

func (c *current) onbranch1(ifBlock, elseIfBlocks, elseBlock interface{}) (interface{}, error) {
	return branch(ifBlock, elseIfBlocks, elseBlock)

}

func (p *parser) callonbranch1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onbranch1(stack["ifBlock"], stack["elseIfBlocks"], stack["elseBlock"])
}

func (c *current) onif_cond12(bop interface{}) (interface{}, error) {
	return ret(bop)

}

func (p *parser) callonif_cond12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onif_cond12(stack["bop"])
}

func (c *current) onif_cond1(cond, bops interface{}) (interface{}, error) {
	return ifBlock(cond, bops)

}

func (p *parser) callonif_cond1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onif_cond1(stack["cond"], stack["bops"])
}

func (c *current) onelse_if14(bop interface{}) (interface{}, error) {
	return ret(bop)

}

func (p *parser) callonelse_if14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onelse_if14(stack["bop"])
}

func (c *current) onelse_if1(cond, bops interface{}) (interface{}, error) {
	return elseIfBlock(cond, bops)

}

func (p *parser) callonelse_if1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onelse_if1(stack["cond"], stack["bops"])
}

func (c *current) onelse_cond9(bop interface{}) (interface{}, error) {
	return ret(bop)

}

func (p *parser) callonelse_cond9() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onelse_cond9(stack["bop"])
}

func (c *current) onelse_cond1(bops interface{}) (interface{}, error) {
	return elseBlock(bops)

}

func (p *parser) callonelse_cond1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onelse_cond1(stack["bops"])
}

func (c *current) oncondition7(bo, cond interface{}) (interface{}, error) {
	return expression(bo, cond)

}

func (p *parser) calloncondition7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.oncondition7(stack["bo"], stack["cond"])
}

func (c *current) oncondition1(cond, conds interface{}) (interface{}, error) {
	return condition(cond, conds)

}

func (p *parser) calloncondition1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.oncondition1(stack["cond"], stack["conds"])
}

func (c *current) onexpression2(cond interface{}) (interface{}, error) {
	return condition_expression(cond)

}

func (p *parser) callonexpression2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onexpression2(stack["cond"])
}

func (c *current) onexpression15(rv interface{}) (interface{}, error) {
	return rvalue(rv)

}

func (p *parser) callonexpression15() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onexpression15(stack["rv"])
}

func (c *current) onnegative_expression2(cond interface{}) (interface{}, error) {
	return negative_expression(cond)

}

func (p *parser) callonnegative_expression2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onnegative_expression2(stack["cond"])
}

func (c *current) onnegative_expression12(sel interface{}) (interface{}, error) {
	return negative_selector(sel)

}

func (p *parser) callonnegative_expression12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onnegative_expression12(stack["sel"])
}

func (c *current) onin_expression1(lv, rv interface{}) (interface{}, error) {
	return in_expression(lv, rv)

}

func (p *parser) callonin_expression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onin_expression1(stack["lv"], stack["rv"])
}

func (c *current) onnot_in_expression1(lv, rv interface{}) (interface{}, error) {
	return not_in_expression(lv, rv)

}

func (p *parser) callonnot_in_expression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onnot_in_expression1(stack["lv"], stack["rv"])
}

func (c *current) oncompare_expression1(lv, co, rv interface{}) (interface{}, error) {
	return compare_expression(lv, co, rv)

}

func (p *parser) calloncompare_expression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.oncompare_expression1(stack["lv"], stack["co"], stack["rv"])
}

func (c *current) oncompare_operator1() (interface{}, error) {
	return compare_operator(string(c.text))

}

func (p *parser) calloncompare_operator1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.oncompare_operator1()
}

func (c *current) onregexp_expression1(lv, ro, rv interface{}) (interface{}, error) {
	return regexp_expression(lv, ro, rv)

}

func (p *parser) callonregexp_expression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onregexp_expression1(stack["lv"], stack["ro"], stack["rv"])
}

func (c *current) onregexp_operator1() (interface{}, error) {
	return regexp_operator(string(c.text))

}

func (p *parser) callonregexp_operator1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onregexp_operator1()
}

func (c *current) onboolean_operator1() (interface{}, error) {
	return boolean_operator(string(c.text))

}

func (p *parser) callonboolean_operator1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onboolean_operator1()
}

func (c *current) onselector1(ses interface{}) (interface{}, error) {
	return selector(ses)

}

func (p *parser) callonselector1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onselector1(stack["ses"])
}

func (c *current) onselector_element1() (interface{}, error) {
	return selector_element(string(c.text))

}

func (p *parser) callonselector_element1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onselector_element1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n == 1 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
