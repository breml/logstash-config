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
			name: "init",
			pos:  position{line: 7, col: 1, offset: 171},
			expr: &actionExpr{
				pos: position{line: 8, col: 5, offset: 182},
				run: (*parser).calloninit1,
				expr: &seqExpr{
					pos: position{line: 8, col: 5, offset: 182},
					exprs: []interface{}{
						&andCodeExpr{
							pos: position{line: 8, col: 5, offset: 182},
							run: (*parser).calloninit3,
						},
						&labeledExpr{
							pos:   position{line: 10, col: 7, offset: 220},
							label: "conf",
							expr: &choiceExpr{
								pos: position{line: 11, col: 9, offset: 235},
								alternatives: []interface{}{
									&actionExpr{
										pos: position{line: 11, col: 9, offset: 235},
										run: (*parser).calloninit6,
										expr: &seqExpr{
											pos: position{line: 11, col: 9, offset: 235},
											exprs: []interface{}{
												&labeledExpr{
													pos:   position{line: 11, col: 9, offset: 235},
													label: "conf",
													expr: &ruleRefExpr{
														pos:  position{line: 11, col: 14, offset: 240},
														name: "config",
													},
												},
												&ruleRefExpr{
													pos:  position{line: 11, col: 21, offset: 247},
													name: "EOF",
												},
											},
										},
									},
									&actionExpr{
										pos: position{line: 13, col: 13, offset: 294},
										run: (*parser).calloninit11,
										expr: &seqExpr{
											pos: position{line: 13, col: 13, offset: 294},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 13, col: 13, offset: 294},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 13, col: 15, offset: 296},
													name: "EOF",
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
			name: "config",
			pos:  position{line: 24, col: 1, offset: 509},
			expr: &actionExpr{
				pos: position{line: 25, col: 5, offset: 522},
				run: (*parser).callonconfig1,
				expr: &seqExpr{
					pos: position{line: 25, col: 5, offset: 522},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 25, col: 5, offset: 522},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 25, col: 7, offset: 524},
							label: "ps",
							expr: &ruleRefExpr{
								pos:  position{line: 25, col: 10, offset: 527},
								name: "pluginSection",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 25, col: 24, offset: 541},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 25, col: 26, offset: 543},
							label: "pss",
							expr: &zeroOrMoreExpr{
								pos: position{line: 25, col: 30, offset: 547},
								expr: &actionExpr{
									pos: position{line: 26, col: 9, offset: 557},
									run: (*parser).callonconfig9,
									expr: &seqExpr{
										pos: position{line: 26, col: 9, offset: 557},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 26, col: 9, offset: 557},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 26, col: 11, offset: 559},
												label: "ps",
												expr: &ruleRefExpr{
													pos:  position{line: 26, col: 14, offset: 562},
													name: "pluginSection",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 29, col: 8, offset: 622},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 29, col: 10, offset: 624},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "comment",
			pos:  position{line: 37, col: 1, offset: 772},
			expr: &oneOrMoreExpr{
				pos: position{line: 38, col: 5, offset: 786},
				expr: &seqExpr{
					pos: position{line: 38, col: 6, offset: 787},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 38, col: 6, offset: 787},
							expr: &ruleRefExpr{
								pos:  position{line: 38, col: 6, offset: 787},
								name: "whitespace",
							},
						},
						&litMatcher{
							pos:        position{line: 38, col: 18, offset: 799},
							val:        "#",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 38, col: 22, offset: 803},
							expr: &charClassMatcher{
								pos:        position{line: 38, col: 22, offset: 803},
								val:        "[^\\r\\n]",
								chars:      []rune{'\r', '\n'},
								ignoreCase: false,
								inverted:   true,
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 38, col: 31, offset: 812},
							expr: &litMatcher{
								pos:        position{line: 38, col: 31, offset: 812},
								val:        "\r",
								ignoreCase: false,
							},
						},
						&choiceExpr{
							pos: position{line: 38, col: 38, offset: 819},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 38, col: 38, offset: 819},
									val:        "\n",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 38, col: 45, offset: 826},
									name: "EOF",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 44, col: 1, offset: 920},
			expr: &zeroOrMoreExpr{
				pos: position{line: 45, col: 5, offset: 928},
				expr: &choiceExpr{
					pos: position{line: 45, col: 6, offset: 929},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 45, col: 6, offset: 929},
							name: "comment",
						},
						&ruleRefExpr{
							pos:  position{line: 45, col: 16, offset: 939},
							name: "whitespace",
						},
					},
				},
			},
		},
		{
			name: "whitespace",
			pos:  position{line: 51, col: 1, offset: 1035},
			expr: &oneOrMoreExpr{
				pos: position{line: 52, col: 5, offset: 1052},
				expr: &charClassMatcher{
					pos:        position{line: 52, col: 5, offset: 1052},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "pluginSection",
			pos:  position{line: 61, col: 1, offset: 1209},
			expr: &actionExpr{
				pos: position{line: 62, col: 5, offset: 1229},
				run: (*parser).callonpluginSection1,
				expr: &seqExpr{
					pos: position{line: 62, col: 5, offset: 1229},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 62, col: 5, offset: 1229},
							label: "pt",
							expr: &ruleRefExpr{
								pos:  position{line: 62, col: 8, offset: 1232},
								name: "pluginType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 62, col: 19, offset: 1243},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 62, col: 21, offset: 1245},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 62, col: 25, offset: 1249},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 62, col: 27, offset: 1251},
							label: "bops",
							expr: &zeroOrMoreExpr{
								pos: position{line: 62, col: 32, offset: 1256},
								expr: &actionExpr{
									pos: position{line: 63, col: 9, offset: 1266},
									run: (*parser).callonpluginSection10,
									expr: &seqExpr{
										pos: position{line: 63, col: 9, offset: 1266},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 63, col: 9, offset: 1266},
												label: "bop",
												expr: &ruleRefExpr{
													pos:  position{line: 63, col: 13, offset: 1270},
													name: "branchOrPlugin",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 63, col: 28, offset: 1285},
												name: "_",
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 67, col: 9, offset: 1344},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 67, col: 9, offset: 1344},
									val:        "}",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 67, col: 15, offset: 1350},
									run: (*parser).callonpluginSection17,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "branchOrPlugin",
			pos:  position{line: 78, col: 1, offset: 1541},
			expr: &choiceExpr{
				pos: position{line: 79, col: 5, offset: 1562},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 79, col: 5, offset: 1562},
						name: "branch",
					},
					&ruleRefExpr{
						pos:  position{line: 79, col: 14, offset: 1571},
						name: "plugin",
					},
				},
			},
		},
		{
			name: "pluginType",
			pos:  position{line: 85, col: 1, offset: 1650},
			expr: &choiceExpr{
				pos: position{line: 86, col: 5, offset: 1667},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 86, col: 5, offset: 1667},
						run: (*parser).callonpluginType2,
						expr: &litMatcher{
							pos:        position{line: 86, col: 5, offset: 1667},
							val:        "input",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 88, col: 9, offset: 1715},
						run: (*parser).callonpluginType4,
						expr: &litMatcher{
							pos:        position{line: 88, col: 9, offset: 1715},
							val:        "filter",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 90, col: 9, offset: 1765},
						run: (*parser).callonpluginType6,
						expr: &litMatcher{
							pos:        position{line: 90, col: 9, offset: 1765},
							val:        "output",
							ignoreCase: false,
						},
					},
					&andCodeExpr{
						pos: position{line: 92, col: 9, offset: 1815},
						run: (*parser).callonpluginType8,
					},
				},
			},
		},
		{
			name: "plugin",
			pos:  position{line: 123, col: 1, offset: 2454},
			expr: &actionExpr{
				pos: position{line: 124, col: 5, offset: 2467},
				run: (*parser).callonplugin1,
				expr: &seqExpr{
					pos: position{line: 124, col: 5, offset: 2467},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 124, col: 5, offset: 2467},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 124, col: 10, offset: 2472},
								name: "name",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 124, col: 15, offset: 2477},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 124, col: 17, offset: 2479},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 124, col: 21, offset: 2483},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 124, col: 23, offset: 2485},
							label: "attributes",
							expr: &zeroOrOneExpr{
								pos: position{line: 124, col: 34, offset: 2496},
								expr: &actionExpr{
									pos: position{line: 125, col: 9, offset: 2506},
									run: (*parser).callonplugin10,
									expr: &seqExpr{
										pos: position{line: 125, col: 9, offset: 2506},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 125, col: 9, offset: 2506},
												label: "attribute",
												expr: &ruleRefExpr{
													pos:  position{line: 125, col: 19, offset: 2516},
													name: "attribute",
												},
											},
											&labeledExpr{
												pos:   position{line: 125, col: 29, offset: 2526},
												label: "attrs",
												expr: &zeroOrMoreExpr{
													pos: position{line: 125, col: 35, offset: 2532},
													expr: &actionExpr{
														pos: position{line: 126, col: 13, offset: 2546},
														run: (*parser).callonplugin16,
														expr: &seqExpr{
															pos: position{line: 126, col: 13, offset: 2546},
															exprs: []interface{}{
																&ruleRefExpr{
																	pos:  position{line: 126, col: 13, offset: 2546},
																	name: "whitespace",
																},
																&ruleRefExpr{
																	pos:  position{line: 126, col: 24, offset: 2557},
																	name: "_",
																},
																&labeledExpr{
																	pos:   position{line: 126, col: 26, offset: 2559},
																	label: "attribute",
																	expr: &ruleRefExpr{
																		pos:  position{line: 126, col: 36, offset: 2569},
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
							pos:  position{line: 132, col: 8, offset: 2711},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 133, col: 9, offset: 2723},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 133, col: 9, offset: 2723},
									val:        "}",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 133, col: 15, offset: 2729},
									run: (*parser).callonplugin25,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "name",
			pos:  position{line: 147, col: 1, offset: 2979},
			expr: &choiceExpr{
				pos: position{line: 148, col: 7, offset: 2992},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 148, col: 7, offset: 2992},
						run: (*parser).callonname2,
						expr: &oneOrMoreExpr{
							pos: position{line: 148, col: 8, offset: 2993},
							expr: &charClassMatcher{
								pos:        position{line: 148, col: 8, offset: 2993},
								val:        "[A-Za-z0-9_-]",
								chars:      []rune{'_', '-'},
								ranges:     []rune{'A', 'Z', 'a', 'z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
					&actionExpr{
						pos: position{line: 150, col: 9, offset: 3054},
						run: (*parser).callonname5,
						expr: &labeledExpr{
							pos:   position{line: 150, col: 9, offset: 3054},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 150, col: 15, offset: 3060},
								name: "stringValue",
							},
						},
					},
				},
			},
		},
		{
			name: "attribute",
			pos:  position{line: 159, col: 1, offset: 3208},
			expr: &actionExpr{
				pos: position{line: 160, col: 5, offset: 3224},
				run: (*parser).callonattribute1,
				expr: &seqExpr{
					pos: position{line: 160, col: 5, offset: 3224},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 160, col: 5, offset: 3224},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 160, col: 10, offset: 3229},
								name: "name",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 160, col: 15, offset: 3234},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 160, col: 17, offset: 3236},
							val:        "=>",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 160, col: 22, offset: 3241},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 160, col: 24, offset: 3243},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 160, col: 30, offset: 3249},
								name: "value",
							},
						},
					},
				},
			},
		},
		{
			name: "value",
			pos:  position{line: 168, col: 1, offset: 3386},
			expr: &choiceExpr{
				pos: position{line: 169, col: 5, offset: 3398},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 169, col: 5, offset: 3398},
						name: "plugin",
					},
					&ruleRefExpr{
						pos:  position{line: 169, col: 14, offset: 3407},
						name: "bareword",
					},
					&ruleRefExpr{
						pos:  position{line: 169, col: 25, offset: 3418},
						name: "stringValue",
					},
					&ruleRefExpr{
						pos:  position{line: 169, col: 39, offset: 3432},
						name: "number",
					},
					&ruleRefExpr{
						pos:  position{line: 169, col: 48, offset: 3441},
						name: "array",
					},
					&ruleRefExpr{
						pos:  position{line: 169, col: 56, offset: 3449},
						name: "hash",
					},
					&andCodeExpr{
						pos: position{line: 169, col: 63, offset: 3456},
						run: (*parser).callonvalue8,
					},
				},
			},
		},
		{
			name: "arrayValue",
			pos:  position{line: 177, col: 1, offset: 3593},
			expr: &choiceExpr{
				pos: position{line: 178, col: 5, offset: 3610},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 178, col: 5, offset: 3610},
						name: "bareword",
					},
					&ruleRefExpr{
						pos:  position{line: 178, col: 16, offset: 3621},
						name: "stringValue",
					},
					&ruleRefExpr{
						pos:  position{line: 178, col: 30, offset: 3635},
						name: "number",
					},
					&ruleRefExpr{
						pos:  position{line: 178, col: 39, offset: 3644},
						name: "array",
					},
					&ruleRefExpr{
						pos:  position{line: 178, col: 47, offset: 3652},
						name: "hash",
					},
					&andCodeExpr{
						pos: position{line: 178, col: 54, offset: 3659},
						run: (*parser).callonarrayValue7,
					},
				},
			},
		},
		{
			name: "bareword",
			pos:  position{line: 187, col: 1, offset: 3822},
			expr: &actionExpr{
				pos: position{line: 188, col: 5, offset: 3837},
				run: (*parser).callonbareword1,
				expr: &seqExpr{
					pos: position{line: 188, col: 5, offset: 3837},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 188, col: 5, offset: 3837},
							val:        "[A-Za-z_]",
							chars:      []rune{'_'},
							ranges:     []rune{'A', 'Z', 'a', 'z'},
							ignoreCase: false,
							inverted:   false,
						},
						&oneOrMoreExpr{
							pos: position{line: 188, col: 15, offset: 3847},
							expr: &charClassMatcher{
								pos:        position{line: 188, col: 15, offset: 3847},
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
			name: "doubleQuotedString",
			pos:  position{line: 196, col: 1, offset: 4057},
			expr: &actionExpr{
				pos: position{line: 197, col: 5, offset: 4082},
				run: (*parser).callondoubleQuotedString1,
				expr: &seqExpr{
					pos: position{line: 197, col: 7, offset: 4084},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 197, col: 7, offset: 4084},
							val:        "\"",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 197, col: 11, offset: 4088},
							expr: &choiceExpr{
								pos: position{line: 197, col: 13, offset: 4090},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 197, col: 13, offset: 4090},
										val:        "\\\"",
										ignoreCase: false,
									},
									&seqExpr{
										pos: position{line: 197, col: 20, offset: 4097},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 197, col: 20, offset: 4097},
												expr: &litMatcher{
													pos:        position{line: 197, col: 21, offset: 4098},
													val:        "\"",
													ignoreCase: false,
												},
											},
											&anyMatcher{
												line: 197, col: 25, offset: 4102,
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 198, col: 9, offset: 4117},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 198, col: 9, offset: 4117},
									val:        "\"",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 198, col: 15, offset: 4123},
									run: (*parser).callondoubleQuotedString13,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "singleQuotedString",
			pos:  position{line: 209, col: 1, offset: 4366},
			expr: &actionExpr{
				pos: position{line: 210, col: 5, offset: 4391},
				run: (*parser).callonsingleQuotedString1,
				expr: &seqExpr{
					pos: position{line: 210, col: 7, offset: 4393},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 210, col: 7, offset: 4393},
							val:        "'",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 210, col: 11, offset: 4397},
							expr: &choiceExpr{
								pos: position{line: 210, col: 13, offset: 4399},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 210, col: 13, offset: 4399},
										val:        "\\'",
										ignoreCase: false,
									},
									&seqExpr{
										pos: position{line: 210, col: 20, offset: 4406},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 210, col: 20, offset: 4406},
												expr: &litMatcher{
													pos:        position{line: 210, col: 21, offset: 4407},
													val:        "'",
													ignoreCase: false,
												},
											},
											&anyMatcher{
												line: 210, col: 25, offset: 4411,
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 211, col: 9, offset: 4426},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 211, col: 9, offset: 4426},
									val:        "'",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 211, col: 15, offset: 4432},
									run: (*parser).callonsingleQuotedString13,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "stringValue",
			pos:  position{line: 222, col: 1, offset: 4640},
			expr: &actionExpr{
				pos: position{line: 223, col: 5, offset: 4658},
				run: (*parser).callonstringValue1,
				expr: &labeledExpr{
					pos:   position{line: 223, col: 5, offset: 4658},
					label: "str",
					expr: &choiceExpr{
						pos: position{line: 223, col: 11, offset: 4664},
						alternatives: []interface{}{
							&actionExpr{
								pos: position{line: 223, col: 11, offset: 4664},
								run: (*parser).callonstringValue4,
								expr: &labeledExpr{
									pos:   position{line: 223, col: 11, offset: 4664},
									label: "str",
									expr: &ruleRefExpr{
										pos:  position{line: 223, col: 15, offset: 4668},
										name: "doubleQuotedString",
									},
								},
							},
							&actionExpr{
								pos: position{line: 225, col: 9, offset: 4776},
								run: (*parser).callonstringValue7,
								expr: &labeledExpr{
									pos:   position{line: 225, col: 9, offset: 4776},
									label: "str",
									expr: &ruleRefExpr{
										pos:  position{line: 225, col: 13, offset: 4780},
										name: "singleQuotedString",
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
			pos:  position{line: 235, col: 1, offset: 5018},
			expr: &actionExpr{
				pos: position{line: 236, col: 5, offset: 5031},
				run: (*parser).callonregexp1,
				expr: &seqExpr{
					pos: position{line: 236, col: 7, offset: 5033},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 236, col: 7, offset: 5033},
							val:        "/",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 236, col: 11, offset: 5037},
							expr: &choiceExpr{
								pos: position{line: 236, col: 13, offset: 5039},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 236, col: 13, offset: 5039},
										val:        "\\/",
										ignoreCase: false,
									},
									&seqExpr{
										pos: position{line: 236, col: 20, offset: 5046},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 236, col: 20, offset: 5046},
												expr: &litMatcher{
													pos:        position{line: 236, col: 21, offset: 5047},
													val:        "/",
													ignoreCase: false,
												},
											},
											&anyMatcher{
												line: 236, col: 25, offset: 5051,
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 237, col: 9, offset: 5066},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 237, col: 9, offset: 5066},
									val:        "/",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 237, col: 15, offset: 5072},
									run: (*parser).callonregexp13,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "number",
			pos:  position{line: 249, col: 1, offset: 5298},
			expr: &actionExpr{
				pos: position{line: 250, col: 5, offset: 5311},
				run: (*parser).callonnumber1,
				expr: &seqExpr{
					pos: position{line: 250, col: 5, offset: 5311},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 250, col: 5, offset: 5311},
							expr: &litMatcher{
								pos:        position{line: 250, col: 5, offset: 5311},
								val:        "-",
								ignoreCase: false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 250, col: 10, offset: 5316},
							expr: &charClassMatcher{
								pos:        position{line: 250, col: 10, offset: 5316},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 250, col: 17, offset: 5323},
							expr: &seqExpr{
								pos: position{line: 250, col: 18, offset: 5324},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 250, col: 18, offset: 5324},
										val:        ".",
										ignoreCase: false,
									},
									&zeroOrMoreExpr{
										pos: position{line: 250, col: 22, offset: 5328},
										expr: &charClassMatcher{
											pos:        position{line: 250, col: 22, offset: 5328},
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
			pos:  position{line: 266, col: 1, offset: 5645},
			expr: &actionExpr{
				pos: position{line: 267, col: 5, offset: 5657},
				run: (*parser).callonarray1,
				expr: &seqExpr{
					pos: position{line: 267, col: 5, offset: 5657},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 267, col: 5, offset: 5657},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 267, col: 9, offset: 5661},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 267, col: 11, offset: 5663},
							label: "value",
							expr: &zeroOrOneExpr{
								pos: position{line: 267, col: 17, offset: 5669},
								expr: &actionExpr{
									pos: position{line: 268, col: 9, offset: 5679},
									run: (*parser).callonarray7,
									expr: &seqExpr{
										pos: position{line: 268, col: 9, offset: 5679},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 268, col: 9, offset: 5679},
												label: "value",
												expr: &ruleRefExpr{
													pos:  position{line: 268, col: 15, offset: 5685},
													name: "value",
												},
											},
											&labeledExpr{
												pos:   position{line: 268, col: 21, offset: 5691},
												label: "values",
												expr: &zeroOrMoreExpr{
													pos: position{line: 268, col: 28, offset: 5698},
													expr: &actionExpr{
														pos: position{line: 269, col: 13, offset: 5712},
														run: (*parser).callonarray13,
														expr: &seqExpr{
															pos: position{line: 269, col: 13, offset: 5712},
															exprs: []interface{}{
																&ruleRefExpr{
																	pos:  position{line: 269, col: 13, offset: 5712},
																	name: "_",
																},
																&litMatcher{
																	pos:        position{line: 269, col: 15, offset: 5714},
																	val:        ",",
																	ignoreCase: false,
																},
																&ruleRefExpr{
																	pos:  position{line: 269, col: 19, offset: 5718},
																	name: "_",
																},
																&labeledExpr{
																	pos:   position{line: 269, col: 21, offset: 5720},
																	label: "value",
																	expr: &ruleRefExpr{
																		pos:  position{line: 269, col: 27, offset: 5726},
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
							pos:  position{line: 275, col: 8, offset: 5857},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 276, col: 9, offset: 5869},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 276, col: 9, offset: 5869},
									val:        "]",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 276, col: 15, offset: 5875},
									run: (*parser).callonarray23,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "hash",
			pos:  position{line: 292, col: 1, offset: 6125},
			expr: &actionExpr{
				pos: position{line: 293, col: 5, offset: 6136},
				run: (*parser).callonhash1,
				expr: &seqExpr{
					pos: position{line: 293, col: 5, offset: 6136},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 293, col: 5, offset: 6136},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 293, col: 9, offset: 6140},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 293, col: 11, offset: 6142},
							label: "entries",
							expr: &zeroOrOneExpr{
								pos: position{line: 293, col: 19, offset: 6150},
								expr: &ruleRefExpr{
									pos:  position{line: 293, col: 19, offset: 6150},
									name: "hashentries",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 293, col: 32, offset: 6163},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 294, col: 9, offset: 6175},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 294, col: 9, offset: 6175},
									val:        "}",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 294, col: 15, offset: 6181},
									run: (*parser).callonhash11,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "hashentries",
			pos:  position{line: 306, col: 1, offset: 6420},
			expr: &actionExpr{
				pos: position{line: 307, col: 5, offset: 6438},
				run: (*parser).callonhashentries1,
				expr: &seqExpr{
					pos: position{line: 307, col: 5, offset: 6438},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 307, col: 5, offset: 6438},
							label: "hashentry",
							expr: &ruleRefExpr{
								pos:  position{line: 307, col: 15, offset: 6448},
								name: "hashentry",
							},
						},
						&labeledExpr{
							pos:   position{line: 307, col: 25, offset: 6458},
							label: "hashentries1",
							expr: &zeroOrMoreExpr{
								pos: position{line: 307, col: 38, offset: 6471},
								expr: &actionExpr{
									pos: position{line: 308, col: 9, offset: 6481},
									run: (*parser).callonhashentries7,
									expr: &seqExpr{
										pos: position{line: 308, col: 9, offset: 6481},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 308, col: 9, offset: 6481},
												name: "whitespace",
											},
											&labeledExpr{
												pos:   position{line: 308, col: 20, offset: 6492},
												label: "hashentry",
												expr: &ruleRefExpr{
													pos:  position{line: 308, col: 30, offset: 6502},
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
			pos:  position{line: 320, col: 1, offset: 6754},
			expr: &actionExpr{
				pos: position{line: 321, col: 5, offset: 6770},
				run: (*parser).callonhashentry1,
				expr: &seqExpr{
					pos: position{line: 321, col: 5, offset: 6770},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 321, col: 5, offset: 6770},
							label: "name",
							expr: &choiceExpr{
								pos: position{line: 321, col: 11, offset: 6776},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 321, col: 11, offset: 6776},
										name: "number",
									},
									&ruleRefExpr{
										pos:  position{line: 321, col: 20, offset: 6785},
										name: "bareword",
									},
									&ruleRefExpr{
										pos:  position{line: 321, col: 31, offset: 6796},
										name: "stringValue",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 321, col: 44, offset: 6809},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 321, col: 46, offset: 6811},
							val:        "=>",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 321, col: 51, offset: 6816},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 321, col: 53, offset: 6818},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 321, col: 59, offset: 6824},
								name: "value",
							},
						},
					},
				},
			},
		},
		{
			name: "branch",
			pos:  position{line: 332, col: 1, offset: 6991},
			expr: &actionExpr{
				pos: position{line: 333, col: 5, offset: 7004},
				run: (*parser).callonbranch1,
				expr: &seqExpr{
					pos: position{line: 333, col: 5, offset: 7004},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 333, col: 5, offset: 7004},
							label: "ifBlock",
							expr: &ruleRefExpr{
								pos:  position{line: 333, col: 13, offset: 7012},
								name: "ifCond",
							},
						},
						&labeledExpr{
							pos:   position{line: 333, col: 20, offset: 7019},
							label: "elseIfBlocks",
							expr: &zeroOrMoreExpr{
								pos: position{line: 333, col: 33, offset: 7032},
								expr: &actionExpr{
									pos: position{line: 334, col: 9, offset: 7042},
									run: (*parser).callonbranch7,
									expr: &seqExpr{
										pos: position{line: 334, col: 9, offset: 7042},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 334, col: 9, offset: 7042},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 334, col: 11, offset: 7044},
												label: "eib",
												expr: &ruleRefExpr{
													pos:  position{line: 334, col: 15, offset: 7048},
													name: "elseIf",
												},
											},
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 337, col: 12, offset: 7114},
							label: "elseBlock",
							expr: &zeroOrOneExpr{
								pos: position{line: 337, col: 22, offset: 7124},
								expr: &actionExpr{
									pos: position{line: 338, col: 13, offset: 7138},
									run: (*parser).callonbranch14,
									expr: &seqExpr{
										pos: position{line: 338, col: 13, offset: 7138},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 338, col: 13, offset: 7138},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 338, col: 15, offset: 7140},
												label: "eb",
												expr: &ruleRefExpr{
													pos:  position{line: 338, col: 18, offset: 7143},
													name: "elseCond",
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
			name: "ifCond",
			pos:  position{line: 350, col: 1, offset: 7391},
			expr: &actionExpr{
				pos: position{line: 351, col: 5, offset: 7404},
				run: (*parser).callonifCond1,
				expr: &seqExpr{
					pos: position{line: 351, col: 5, offset: 7404},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 351, col: 5, offset: 7404},
							val:        "if",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 351, col: 10, offset: 7409},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 351, col: 12, offset: 7411},
							label: "cond",
							expr: &ruleRefExpr{
								pos:  position{line: 351, col: 17, offset: 7416},
								name: "condition",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 351, col: 27, offset: 7426},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 351, col: 29, offset: 7428},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 351, col: 33, offset: 7432},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 351, col: 35, offset: 7434},
							label: "bops",
							expr: &zeroOrMoreExpr{
								pos: position{line: 351, col: 40, offset: 7439},
								expr: &actionExpr{
									pos: position{line: 352, col: 13, offset: 7453},
									run: (*parser).callonifCond12,
									expr: &seqExpr{
										pos: position{line: 352, col: 13, offset: 7453},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 352, col: 13, offset: 7453},
												label: "bop",
												expr: &ruleRefExpr{
													pos:  position{line: 352, col: 17, offset: 7457},
													name: "branchOrPlugin",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 352, col: 32, offset: 7472},
												name: "_",
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 356, col: 13, offset: 7547},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 356, col: 13, offset: 7547},
									val:        "}",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 356, col: 19, offset: 7553},
									run: (*parser).callonifCond19,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "elseIf",
			pos:  position{line: 368, col: 1, offset: 7827},
			expr: &actionExpr{
				pos: position{line: 369, col: 5, offset: 7840},
				run: (*parser).callonelseIf1,
				expr: &seqExpr{
					pos: position{line: 369, col: 5, offset: 7840},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 369, col: 5, offset: 7840},
							val:        "else",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 369, col: 12, offset: 7847},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 369, col: 14, offset: 7849},
							val:        "if",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 369, col: 19, offset: 7854},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 369, col: 21, offset: 7856},
							label: "cond",
							expr: &ruleRefExpr{
								pos:  position{line: 369, col: 26, offset: 7861},
								name: "condition",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 369, col: 36, offset: 7871},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 369, col: 38, offset: 7873},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 369, col: 42, offset: 7877},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 369, col: 44, offset: 7879},
							label: "bops",
							expr: &zeroOrMoreExpr{
								pos: position{line: 369, col: 49, offset: 7884},
								expr: &actionExpr{
									pos: position{line: 370, col: 9, offset: 7894},
									run: (*parser).callonelseIf14,
									expr: &seqExpr{
										pos: position{line: 370, col: 9, offset: 7894},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 370, col: 9, offset: 7894},
												label: "bop",
												expr: &ruleRefExpr{
													pos:  position{line: 370, col: 13, offset: 7898},
													name: "branchOrPlugin",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 370, col: 28, offset: 7913},
												name: "_",
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 374, col: 9, offset: 7972},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 374, col: 9, offset: 7972},
									val:        "}",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 374, col: 15, offset: 7978},
									run: (*parser).callonelseIf21,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "elseCond",
			pos:  position{line: 386, col: 1, offset: 8220},
			expr: &actionExpr{
				pos: position{line: 387, col: 5, offset: 8235},
				run: (*parser).callonelseCond1,
				expr: &seqExpr{
					pos: position{line: 387, col: 5, offset: 8235},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 387, col: 5, offset: 8235},
							val:        "else",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 387, col: 12, offset: 8242},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 387, col: 14, offset: 8244},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 387, col: 18, offset: 8248},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 387, col: 20, offset: 8250},
							label: "bops",
							expr: &zeroOrMoreExpr{
								pos: position{line: 387, col: 25, offset: 8255},
								expr: &actionExpr{
									pos: position{line: 388, col: 9, offset: 8265},
									run: (*parser).callonelseCond9,
									expr: &seqExpr{
										pos: position{line: 388, col: 9, offset: 8265},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 388, col: 9, offset: 8265},
												label: "bop",
												expr: &ruleRefExpr{
													pos:  position{line: 388, col: 13, offset: 8269},
													name: "branchOrPlugin",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 388, col: 28, offset: 8284},
												name: "_",
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 392, col: 9, offset: 8343},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 392, col: 9, offset: 8343},
									val:        "}",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 392, col: 15, offset: 8349},
									run: (*parser).callonelseCond16,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "condition",
			pos:  position{line: 404, col: 1, offset: 8598},
			expr: &actionExpr{
				pos: position{line: 405, col: 5, offset: 8614},
				run: (*parser).calloncondition1,
				expr: &seqExpr{
					pos: position{line: 405, col: 5, offset: 8614},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 405, col: 5, offset: 8614},
							label: "cond",
							expr: &ruleRefExpr{
								pos:  position{line: 405, col: 10, offset: 8619},
								name: "expression",
							},
						},
						&labeledExpr{
							pos:   position{line: 405, col: 21, offset: 8630},
							label: "conds",
							expr: &zeroOrMoreExpr{
								pos: position{line: 405, col: 27, offset: 8636},
								expr: &actionExpr{
									pos: position{line: 406, col: 9, offset: 8646},
									run: (*parser).calloncondition7,
									expr: &seqExpr{
										pos: position{line: 406, col: 9, offset: 8646},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 406, col: 9, offset: 8646},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 406, col: 11, offset: 8648},
												label: "bo",
												expr: &ruleRefExpr{
													pos:  position{line: 406, col: 14, offset: 8651},
													name: "booleanOperator",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 406, col: 30, offset: 8667},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 406, col: 32, offset: 8669},
												label: "cond",
												expr: &ruleRefExpr{
													pos:  position{line: 406, col: 37, offset: 8674},
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
			pos:  position{line: 425, col: 1, offset: 9073},
			expr: &choiceExpr{
				pos: position{line: 427, col: 9, offset: 9100},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 428, col: 13, offset: 9114},
						run: (*parser).callonexpression2,
						expr: &seqExpr{
							pos: position{line: 428, col: 13, offset: 9114},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 428, col: 13, offset: 9114},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 428, col: 17, offset: 9118},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 428, col: 19, offset: 9120},
									label: "cond",
									expr: &ruleRefExpr{
										pos:  position{line: 428, col: 24, offset: 9125},
										name: "condition",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 428, col: 34, offset: 9135},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 428, col: 36, offset: 9137},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 432, col: 9, offset: 9224},
						name: "negativeExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 433, col: 9, offset: 9251},
						name: "inExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 434, col: 9, offset: 9272},
						name: "notInExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 435, col: 9, offset: 9296},
						name: "compareExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 436, col: 9, offset: 9322},
						name: "regexpExpression",
					},
					&actionExpr{
						pos: position{line: 437, col: 9, offset: 9347},
						run: (*parser).callonexpression15,
						expr: &labeledExpr{
							pos:   position{line: 437, col: 9, offset: 9347},
							label: "rv",
							expr: &ruleRefExpr{
								pos:  position{line: 437, col: 12, offset: 9350},
								name: "rvalue",
							},
						},
					},
				},
			},
		},
		{
			name: "negativeExpression",
			pos:  position{line: 450, col: 1, offset: 9634},
			expr: &choiceExpr{
				pos: position{line: 452, col: 9, offset: 9669},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 453, col: 13, offset: 9683},
						run: (*parser).callonnegativeExpression2,
						expr: &seqExpr{
							pos: position{line: 453, col: 13, offset: 9683},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 453, col: 13, offset: 9683},
									val:        "!",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 453, col: 17, offset: 9687},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 453, col: 19, offset: 9689},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 453, col: 23, offset: 9693},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 453, col: 25, offset: 9695},
									label: "cond",
									expr: &ruleRefExpr{
										pos:  position{line: 453, col: 30, offset: 9700},
										name: "condition",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 453, col: 40, offset: 9710},
									name: "_",
								},
								&choiceExpr{
									pos: position{line: 454, col: 17, offset: 9730},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 454, col: 17, offset: 9730},
											val:        ")",
											ignoreCase: false,
										},
										&andCodeExpr{
											pos: position{line: 454, col: 23, offset: 9736},
											run: (*parser).callonnegativeExpression13,
										},
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 462, col: 11, offset: 9936},
						run: (*parser).callonnegativeExpression14,
						expr: &seqExpr{
							pos: position{line: 462, col: 11, offset: 9936},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 462, col: 11, offset: 9936},
									val:        "!",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 462, col: 15, offset: 9940},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 462, col: 17, offset: 9942},
									label: "sel",
									expr: &ruleRefExpr{
										pos:  position{line: 462, col: 21, offset: 9946},
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
			name: "inExpression",
			pos:  position{line: 473, col: 1, offset: 10145},
			expr: &actionExpr{
				pos: position{line: 474, col: 5, offset: 10164},
				run: (*parser).calloninExpression1,
				expr: &seqExpr{
					pos: position{line: 474, col: 5, offset: 10164},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 474, col: 5, offset: 10164},
							label: "lv",
							expr: &ruleRefExpr{
								pos:  position{line: 474, col: 8, offset: 10167},
								name: "rvalue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 474, col: 15, offset: 10174},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 474, col: 17, offset: 10176},
							name: "inOperator",
						},
						&ruleRefExpr{
							pos:  position{line: 474, col: 28, offset: 10187},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 474, col: 30, offset: 10189},
							label: "rv",
							expr: &ruleRefExpr{
								pos:  position{line: 474, col: 33, offset: 10192},
								name: "rvalue",
							},
						},
					},
				},
			},
		},
		{
			name: "notInExpression",
			pos:  position{line: 483, col: 1, offset: 10371},
			expr: &actionExpr{
				pos: position{line: 484, col: 5, offset: 10393},
				run: (*parser).callonnotInExpression1,
				expr: &seqExpr{
					pos: position{line: 484, col: 5, offset: 10393},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 484, col: 5, offset: 10393},
							label: "lv",
							expr: &ruleRefExpr{
								pos:  position{line: 484, col: 8, offset: 10396},
								name: "rvalue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 484, col: 15, offset: 10403},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 484, col: 17, offset: 10405},
							name: "notInOperator",
						},
						&ruleRefExpr{
							pos:  position{line: 484, col: 31, offset: 10419},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 484, col: 33, offset: 10421},
							label: "rv",
							expr: &ruleRefExpr{
								pos:  position{line: 484, col: 36, offset: 10424},
								name: "rvalue",
							},
						},
					},
				},
			},
		},
		{
			name: "inOperator",
			pos:  position{line: 492, col: 1, offset: 10523},
			expr: &choiceExpr{
				pos: position{line: 493, col: 5, offset: 10540},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 493, col: 5, offset: 10540},
						val:        "in",
						ignoreCase: false,
					},
					&andCodeExpr{
						pos: position{line: 493, col: 12, offset: 10547},
						run: (*parser).calloninOperator3,
					},
				},
			},
		},
		{
			name: "notInOperator",
			pos:  position{line: 501, col: 1, offset: 10669},
			expr: &choiceExpr{
				pos: position{line: 502, col: 5, offset: 10689},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 502, col: 5, offset: 10689},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 502, col: 5, offset: 10689},
								val:        "not ",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 502, col: 12, offset: 10696},
								name: "_",
							},
							&litMatcher{
								pos:        position{line: 502, col: 14, offset: 10698},
								val:        "in",
								ignoreCase: false,
							},
						},
					},
					&andCodeExpr{
						pos: position{line: 502, col: 21, offset: 10705},
						run: (*parser).callonnotInOperator6,
					},
				},
			},
		},
		{
			name: "rvalue",
			pos:  position{line: 513, col: 1, offset: 11031},
			expr: &choiceExpr{
				pos: position{line: 514, col: 5, offset: 11044},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 514, col: 5, offset: 11044},
						name: "stringValue",
					},
					&ruleRefExpr{
						pos:  position{line: 514, col: 19, offset: 11058},
						name: "number",
					},
					&ruleRefExpr{
						pos:  position{line: 514, col: 28, offset: 11067},
						name: "selector",
					},
					&ruleRefExpr{
						pos:  position{line: 514, col: 39, offset: 11078},
						name: "array",
					},
					&ruleRefExpr{
						pos:  position{line: 514, col: 47, offset: 11086},
						name: "regexp",
					},
					&andCodeExpr{
						pos: position{line: 514, col: 56, offset: 11095},
						run: (*parser).callonrvalue7,
					},
				},
			},
		},
		{
			name: "compareExpression",
			pos:  position{line: 548, col: 1, offset: 11833},
			expr: &actionExpr{
				pos: position{line: 549, col: 5, offset: 11857},
				run: (*parser).calloncompareExpression1,
				expr: &seqExpr{
					pos: position{line: 549, col: 5, offset: 11857},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 549, col: 5, offset: 11857},
							label: "lv",
							expr: &ruleRefExpr{
								pos:  position{line: 549, col: 8, offset: 11860},
								name: "rvalue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 549, col: 15, offset: 11867},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 549, col: 17, offset: 11869},
							label: "co",
							expr: &ruleRefExpr{
								pos:  position{line: 549, col: 20, offset: 11872},
								name: "compareOperator",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 549, col: 36, offset: 11888},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 549, col: 38, offset: 11890},
							label: "rv",
							expr: &ruleRefExpr{
								pos:  position{line: 549, col: 41, offset: 11893},
								name: "rvalue",
							},
						},
					},
				},
			},
		},
		{
			name: "compareOperator",
			pos:  position{line: 558, col: 1, offset: 12089},
			expr: &actionExpr{
				pos: position{line: 559, col: 5, offset: 12111},
				run: (*parser).calloncompareOperator1,
				expr: &choiceExpr{
					pos: position{line: 559, col: 6, offset: 12112},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 559, col: 6, offset: 12112},
							val:        "==",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 559, col: 13, offset: 12119},
							val:        "!=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 559, col: 20, offset: 12126},
							val:        "<=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 559, col: 27, offset: 12133},
							val:        ">=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 559, col: 34, offset: 12140},
							val:        "<",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 559, col: 40, offset: 12146},
							val:        ">",
							ignoreCase: false,
						},
						&andCodeExpr{
							pos: position{line: 559, col: 46, offset: 12152},
							run: (*parser).calloncompareOperator9,
						},
					},
				},
			},
		},
		{
			name: "regexpExpression",
			pos:  position{line: 570, col: 1, offset: 12436},
			expr: &actionExpr{
				pos: position{line: 571, col: 5, offset: 12459},
				run: (*parser).callonregexpExpression1,
				expr: &seqExpr{
					pos: position{line: 571, col: 5, offset: 12459},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 571, col: 5, offset: 12459},
							label: "lv",
							expr: &ruleRefExpr{
								pos:  position{line: 571, col: 8, offset: 12462},
								name: "rvalue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 571, col: 15, offset: 12469},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 571, col: 18, offset: 12472},
							label: "ro",
							expr: &ruleRefExpr{
								pos:  position{line: 571, col: 21, offset: 12475},
								name: "regexpOperator",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 571, col: 36, offset: 12490},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 571, col: 38, offset: 12492},
							label: "rv",
							expr: &choiceExpr{
								pos: position{line: 571, col: 42, offset: 12496},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 571, col: 42, offset: 12496},
										name: "stringValue",
									},
									&ruleRefExpr{
										pos:  position{line: 571, col: 56, offset: 12510},
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
			name: "regexpOperator",
			pos:  position{line: 579, col: 1, offset: 12668},
			expr: &actionExpr{
				pos: position{line: 580, col: 5, offset: 12689},
				run: (*parser).callonregexpOperator1,
				expr: &choiceExpr{
					pos: position{line: 580, col: 6, offset: 12690},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 580, col: 6, offset: 12690},
							val:        "=~",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 580, col: 13, offset: 12697},
							val:        "!~",
							ignoreCase: false,
						},
						&andCodeExpr{
							pos: position{line: 580, col: 20, offset: 12704},
							run: (*parser).callonregexpOperator5,
						},
					},
				},
			},
		},
		{
			name: "booleanOperator",
			pos:  position{line: 591, col: 1, offset: 12967},
			expr: &actionExpr{
				pos: position{line: 592, col: 5, offset: 12989},
				run: (*parser).callonbooleanOperator1,
				expr: &choiceExpr{
					pos: position{line: 592, col: 6, offset: 12990},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 592, col: 6, offset: 12990},
							val:        "and",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 592, col: 14, offset: 12998},
							val:        "or",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 592, col: 21, offset: 13005},
							val:        "xor",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 592, col: 29, offset: 13013},
							val:        "nand",
							ignoreCase: false,
						},
						&andCodeExpr{
							pos: position{line: 592, col: 38, offset: 13022},
							run: (*parser).callonbooleanOperator7,
						},
					},
				},
			},
		},
		{
			name: "selector",
			pos:  position{line: 603, col: 1, offset: 13259},
			expr: &actionExpr{
				pos: position{line: 604, col: 5, offset: 13274},
				run: (*parser).callonselector1,
				expr: &labeledExpr{
					pos:   position{line: 604, col: 5, offset: 13274},
					label: "ses",
					expr: &oneOrMoreExpr{
						pos: position{line: 604, col: 9, offset: 13278},
						expr: &ruleRefExpr{
							pos:  position{line: 604, col: 9, offset: 13278},
							name: "selectorElement",
						},
					},
				},
			},
		},
		{
			name: "selectorElement",
			pos:  position{line: 613, col: 1, offset: 13441},
			expr: &actionExpr{
				pos: position{line: 614, col: 5, offset: 13463},
				run: (*parser).callonselectorElement1,
				expr: &seqExpr{
					pos: position{line: 614, col: 5, offset: 13463},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 614, col: 5, offset: 13463},
							val:        "[",
							ignoreCase: false,
						},
						&oneOrMoreExpr{
							pos: position{line: 614, col: 9, offset: 13467},
							expr: &charClassMatcher{
								pos:        position{line: 614, col: 9, offset: 13467},
								val:        "[^\\],]",
								chars:      []rune{']', ','},
								ignoreCase: false,
								inverted:   true,
							},
						},
						&choiceExpr{
							pos: position{line: 615, col: 9, offset: 13485},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 615, col: 9, offset: 13485},
									val:        "]",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 615, col: 15, offset: 13491},
									run: (*parser).callonselectorElement8,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 622, col: 1, offset: 13631},
			expr: &notExpr{
				pos: position{line: 622, col: 7, offset: 13637},
				expr: &anyMatcher{
					line: 622, col: 8, offset: 13638,
				},
			},
		},
	},
}

func (c *current) oninit3() (bool, error) {
	return initParser()

}

func (p *parser) calloninit3() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.oninit3()
}

func (c *current) oninit6(conf interface{}) (interface{}, error) {
	return ret(conf)

}

func (p *parser) calloninit6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.oninit6(stack["conf"])
}

func (c *current) oninit11() (interface{}, error) {
	return ast.NewConfig(nil, nil, nil), nil

}

func (p *parser) calloninit11() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.oninit11()
}

func (c *current) oninit1(conf interface{}) (interface{}, error) {
	return ret(conf)

}

func (p *parser) calloninit1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.oninit1(stack["conf"])
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

func (c *current) onpluginSection10(bop interface{}) (interface{}, error) {
	return ret(bop)

}

func (p *parser) callonpluginSection10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onpluginSection10(stack["bop"])
}

func (c *current) onpluginSection17() (bool, error) {
	return pushError("expect closing curly bracket", c)

}

func (p *parser) callonpluginSection17() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onpluginSection17()
}

func (c *current) onpluginSection1(pt, bops interface{}) (interface{}, error) {
	return pluginSection(pt, bops)

}

func (p *parser) callonpluginSection1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onpluginSection1(stack["pt"], stack["bops"])
}

func (c *current) onpluginType2() (interface{}, error) {
	return ast.Input, nil

}

func (p *parser) callonpluginType2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onpluginType2()
}

func (c *current) onpluginType4() (interface{}, error) {
	return ast.Filter, nil

}

func (p *parser) callonpluginType4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onpluginType4()
}

func (c *current) onpluginType6() (interface{}, error) {
	return ast.Output, nil

}

func (p *parser) callonpluginType6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onpluginType6()
}

func (c *current) onpluginType8() (bool, error) {
	return pushError("expect plugin type (input, filter, output)", c)

}

func (p *parser) callonpluginType8() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onpluginType8()
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

func (c *current) onplugin25() (bool, error) {
	return fatalError("expect closing curly bracket", c)

}

func (p *parser) callonplugin25() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onplugin25()
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

func (c *current) onvalue8() (bool, error) {
	return fatalError("invalid value", c)

}

func (p *parser) callonvalue8() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onvalue8()
}

func (c *current) onarrayValue7() (bool, error) {
	return fatalError("invalid array value", c)

}

func (p *parser) callonarrayValue7() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onarrayValue7()
}

func (c *current) onbareword1() (interface{}, error) {
	return ast.NewStringAttribute("", string(c.text), ast.Bareword), nil

}

func (p *parser) callonbareword1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onbareword1()
}

func (c *current) ondoubleQuotedString13() (bool, error) {
	return fatalError("expect closing double quotes (\")", c)

}

func (p *parser) callondoubleQuotedString13() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.ondoubleQuotedString13()
}

func (c *current) ondoubleQuotedString1() (interface{}, error) {
	return enclosedValue(c)

}

func (p *parser) callondoubleQuotedString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.ondoubleQuotedString1()
}

func (c *current) onsingleQuotedString13() (bool, error) {
	return fatalError("expect closing single quote (')", c)

}

func (p *parser) callonsingleQuotedString13() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onsingleQuotedString13()
}

func (c *current) onsingleQuotedString1() (interface{}, error) {
	return enclosedValue(c)

}

func (p *parser) callonsingleQuotedString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onsingleQuotedString1()
}

func (c *current) onstringValue4(str interface{}) (interface{}, error) {
	return ast.NewStringAttribute("", str.(string), ast.DoubleQuoted), nil

}

func (p *parser) callonstringValue4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onstringValue4(stack["str"])
}

func (c *current) onstringValue7(str interface{}) (interface{}, error) {
	return ast.NewStringAttribute("", str.(string), ast.SingleQuoted), nil

}

func (p *parser) callonstringValue7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onstringValue7(stack["str"])
}

func (c *current) onstringValue1(str interface{}) (interface{}, error) {
	return ret(str)

}

func (p *parser) callonstringValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onstringValue1(stack["str"])
}

func (c *current) onregexp13() (bool, error) {
	return fatalError("expect closing slash (/) for regexp", c)

}

func (p *parser) callonregexp13() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onregexp13()
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

func (c *current) onarray23() (bool, error) {
	return fatalError("expect closing square bracket", c)

}

func (p *parser) callonarray23() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onarray23()
}

func (c *current) onarray1(value interface{}) (interface{}, error) {
	return array(value)

}

func (p *parser) callonarray1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onarray1(stack["value"])
}

func (c *current) onhash11() (bool, error) {
	return fatalError("expect closing curly bracket", c)

}

func (p *parser) callonhash11() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onhash11()
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

func (c *current) onifCond12(bop interface{}) (interface{}, error) {
	return ret(bop)

}

func (p *parser) callonifCond12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onifCond12(stack["bop"])
}

func (c *current) onifCond19() (bool, error) {
	return fatalError("expect closing curly bracket", c)

}

func (p *parser) callonifCond19() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onifCond19()
}

func (c *current) onifCond1(cond, bops interface{}) (interface{}, error) {
	return ifBlock(cond, bops)

}

func (p *parser) callonifCond1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onifCond1(stack["cond"], stack["bops"])
}

func (c *current) onelseIf14(bop interface{}) (interface{}, error) {
	return ret(bop)

}

func (p *parser) callonelseIf14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onelseIf14(stack["bop"])
}

func (c *current) onelseIf21() (bool, error) {
	return fatalError("expect closing curly bracket", c)

}

func (p *parser) callonelseIf21() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onelseIf21()
}

func (c *current) onelseIf1(cond, bops interface{}) (interface{}, error) {
	return elseIfBlock(cond, bops)

}

func (p *parser) callonelseIf1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onelseIf1(stack["cond"], stack["bops"])
}

func (c *current) onelseCond9(bop interface{}) (interface{}, error) {
	return ret(bop)

}

func (p *parser) callonelseCond9() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onelseCond9(stack["bop"])
}

func (c *current) onelseCond16() (bool, error) {
	return fatalError("expect closing curly bracket", c)

}

func (p *parser) callonelseCond16() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onelseCond16()
}

func (c *current) onelseCond1(bops interface{}) (interface{}, error) {
	return elseBlock(bops)

}

func (p *parser) callonelseCond1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onelseCond1(stack["bops"])
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
	return conditionExpression(cond)

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

func (c *current) onnegativeExpression13() (bool, error) {
	return fatalError("expect closing parenthesis", c)

}

func (p *parser) callonnegativeExpression13() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onnegativeExpression13()
}

func (c *current) onnegativeExpression2(cond interface{}) (interface{}, error) {
	return negativeExpression(cond)

}

func (p *parser) callonnegativeExpression2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onnegativeExpression2(stack["cond"])
}

func (c *current) onnegativeExpression14(sel interface{}) (interface{}, error) {
	return negativeSelector(sel)

}

func (p *parser) callonnegativeExpression14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onnegativeExpression14(stack["sel"])
}

func (c *current) oninExpression1(lv, rv interface{}) (interface{}, error) {
	return inExpression(lv, rv)

}

func (p *parser) calloninExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.oninExpression1(stack["lv"], stack["rv"])
}

func (c *current) onnotInExpression1(lv, rv interface{}) (interface{}, error) {
	return notInExpression(lv, rv)

}

func (p *parser) callonnotInExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onnotInExpression1(stack["lv"], stack["rv"])
}

func (c *current) oninOperator3() (bool, error) {
	return pushError("expect in operator (in)", c)

}

func (p *parser) calloninOperator3() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.oninOperator3()
}

func (c *current) onnotInOperator6() (bool, error) {
	return pushError("expect not in operator (not in)", c)

}

func (p *parser) callonnotInOperator6() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onnotInOperator6()
}

func (c *current) onrvalue7() (bool, error) {
	return pushError("invalid value for expression", c)

}

func (p *parser) callonrvalue7() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onrvalue7()
}

func (c *current) oncompareExpression1(lv, co, rv interface{}) (interface{}, error) {
	return compareExpression(lv, co, rv)

}

func (p *parser) calloncompareExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.oncompareExpression1(stack["lv"], stack["co"], stack["rv"])
}

func (c *current) oncompareOperator9() (bool, error) {
	return pushError("expect compare operator (==, !=, <=, >=, <, >)", c)

}

func (p *parser) calloncompareOperator9() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.oncompareOperator9()
}

func (c *current) oncompareOperator1() (interface{}, error) {
	return compareOperator(string(c.text))

}

func (p *parser) calloncompareOperator1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.oncompareOperator1()
}

func (c *current) onregexpExpression1(lv, ro, rv interface{}) (interface{}, error) {
	return regexpExpression(lv, ro, rv)

}

func (p *parser) callonregexpExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onregexpExpression1(stack["lv"], stack["ro"], stack["rv"])
}

func (c *current) onregexpOperator5() (bool, error) {
	return pushError("expect regexp comparison operator (=~, !~)", c)

}

func (p *parser) callonregexpOperator5() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onregexpOperator5()
}

func (c *current) onregexpOperator1() (interface{}, error) {
	return regexpOperator(string(c.text))

}

func (p *parser) callonregexpOperator1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onregexpOperator1()
}

func (c *current) onbooleanOperator7() (bool, error) {
	return pushError("expect boolean operator (and, or, xor, nand)", c)

}

func (p *parser) callonbooleanOperator7() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onbooleanOperator7()
}

func (c *current) onbooleanOperator1() (interface{}, error) {
	return booleanOperator(string(c.text))

}

func (p *parser) callonbooleanOperator1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onbooleanOperator1()
}

func (c *current) onselector1(ses interface{}) (interface{}, error) {
	return selector(ses)

}

func (p *parser) callonselector1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onselector1(stack["ses"])
}

func (c *current) onselectorElement8() (bool, error) {
	return pushError("expect closing square bracket", c)

}

func (p *parser) callonselectorElement8() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onselectorElement8()
}

func (c *current) onselectorElement1() (interface{}, error) {
	return selectorElement(string(c.text))

}

func (p *parser) callonselectorElement1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onselectorElement1()
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
