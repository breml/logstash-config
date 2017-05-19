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
							expr: &ruleRefExpr{
								pos:  position{line: 10, col: 12, offset: 225},
								name: "config",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 10, col: 19, offset: 232},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "config",
			pos:  position{line: 18, col: 1, offset: 424},
			expr: &actionExpr{
				pos: position{line: 19, col: 5, offset: 437},
				run: (*parser).callonconfig1,
				expr: &seqExpr{
					pos: position{line: 19, col: 5, offset: 437},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 19, col: 5, offset: 437},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 19, col: 7, offset: 439},
							label: "ps",
							expr: &ruleRefExpr{
								pos:  position{line: 19, col: 10, offset: 442},
								name: "pluginSection",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 19, col: 24, offset: 456},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 19, col: 26, offset: 458},
							label: "pss",
							expr: &zeroOrMoreExpr{
								pos: position{line: 19, col: 30, offset: 462},
								expr: &actionExpr{
									pos: position{line: 20, col: 9, offset: 472},
									run: (*parser).callonconfig9,
									expr: &seqExpr{
										pos: position{line: 20, col: 9, offset: 472},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 20, col: 9, offset: 472},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 20, col: 11, offset: 474},
												label: "ps",
												expr: &ruleRefExpr{
													pos:  position{line: 20, col: 14, offset: 477},
													name: "pluginSection",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 23, col: 8, offset: 537},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 23, col: 10, offset: 539},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "comment",
			pos:  position{line: 31, col: 1, offset: 687},
			expr: &oneOrMoreExpr{
				pos: position{line: 32, col: 5, offset: 701},
				expr: &seqExpr{
					pos: position{line: 32, col: 6, offset: 702},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 32, col: 6, offset: 702},
							expr: &ruleRefExpr{
								pos:  position{line: 32, col: 6, offset: 702},
								name: "whitespace",
							},
						},
						&litMatcher{
							pos:        position{line: 32, col: 18, offset: 714},
							val:        "#",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 32, col: 22, offset: 718},
							expr: &charClassMatcher{
								pos:        position{line: 32, col: 22, offset: 718},
								val:        "[^\\r\\n]",
								chars:      []rune{'\r', '\n'},
								ignoreCase: false,
								inverted:   true,
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 32, col: 31, offset: 727},
							expr: &litMatcher{
								pos:        position{line: 32, col: 31, offset: 727},
								val:        "\r",
								ignoreCase: false,
							},
						},
						&litMatcher{
							pos:        position{line: 32, col: 37, offset: 733},
							val:        "\n",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 38, col: 1, offset: 827},
			expr: &zeroOrMoreExpr{
				pos: position{line: 39, col: 5, offset: 835},
				expr: &choiceExpr{
					pos: position{line: 39, col: 6, offset: 836},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 39, col: 6, offset: 836},
							name: "comment",
						},
						&ruleRefExpr{
							pos:  position{line: 39, col: 16, offset: 846},
							name: "whitespace",
						},
					},
				},
			},
		},
		{
			name: "whitespace",
			pos:  position{line: 45, col: 1, offset: 942},
			expr: &oneOrMoreExpr{
				pos: position{line: 46, col: 5, offset: 959},
				expr: &charClassMatcher{
					pos:        position{line: 46, col: 5, offset: 959},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "pluginSection",
			pos:  position{line: 55, col: 1, offset: 1116},
			expr: &actionExpr{
				pos: position{line: 56, col: 5, offset: 1136},
				run: (*parser).callonpluginSection1,
				expr: &seqExpr{
					pos: position{line: 56, col: 5, offset: 1136},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 56, col: 5, offset: 1136},
							label: "pt",
							expr: &ruleRefExpr{
								pos:  position{line: 56, col: 8, offset: 1139},
								name: "pluginType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 56, col: 19, offset: 1150},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 56, col: 21, offset: 1152},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 56, col: 25, offset: 1156},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 56, col: 27, offset: 1158},
							label: "bops",
							expr: &zeroOrMoreExpr{
								pos: position{line: 56, col: 32, offset: 1163},
								expr: &actionExpr{
									pos: position{line: 57, col: 9, offset: 1173},
									run: (*parser).callonpluginSection10,
									expr: &seqExpr{
										pos: position{line: 57, col: 9, offset: 1173},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 57, col: 9, offset: 1173},
												label: "bop",
												expr: &ruleRefExpr{
													pos:  position{line: 57, col: 13, offset: 1177},
													name: "branchOrPlugin",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 57, col: 28, offset: 1192},
												name: "_",
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 61, col: 9, offset: 1251},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 61, col: 9, offset: 1251},
									val:        "}",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 61, col: 15, offset: 1257},
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
			pos:  position{line: 72, col: 1, offset: 1448},
			expr: &choiceExpr{
				pos: position{line: 73, col: 5, offset: 1469},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 73, col: 5, offset: 1469},
						name: "branch",
					},
					&ruleRefExpr{
						pos:  position{line: 73, col: 14, offset: 1478},
						name: "plugin",
					},
				},
			},
		},
		{
			name: "pluginType",
			pos:  position{line: 79, col: 1, offset: 1557},
			expr: &choiceExpr{
				pos: position{line: 80, col: 5, offset: 1574},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 80, col: 5, offset: 1574},
						run: (*parser).callonpluginType2,
						expr: &litMatcher{
							pos:        position{line: 80, col: 5, offset: 1574},
							val:        "input",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 82, col: 9, offset: 1622},
						run: (*parser).callonpluginType4,
						expr: &litMatcher{
							pos:        position{line: 82, col: 9, offset: 1622},
							val:        "filter",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 84, col: 9, offset: 1672},
						run: (*parser).callonpluginType6,
						expr: &litMatcher{
							pos:        position{line: 84, col: 9, offset: 1672},
							val:        "output",
							ignoreCase: false,
						},
					},
					&andCodeExpr{
						pos: position{line: 86, col: 9, offset: 1722},
						run: (*parser).callonpluginType8,
					},
				},
			},
		},
		{
			name: "plugin",
			pos:  position{line: 117, col: 1, offset: 2361},
			expr: &actionExpr{
				pos: position{line: 118, col: 5, offset: 2374},
				run: (*parser).callonplugin1,
				expr: &seqExpr{
					pos: position{line: 118, col: 5, offset: 2374},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 118, col: 5, offset: 2374},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 118, col: 10, offset: 2379},
								name: "name",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 118, col: 15, offset: 2384},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 118, col: 17, offset: 2386},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 118, col: 21, offset: 2390},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 118, col: 23, offset: 2392},
							label: "attributes",
							expr: &zeroOrOneExpr{
								pos: position{line: 118, col: 34, offset: 2403},
								expr: &actionExpr{
									pos: position{line: 119, col: 9, offset: 2413},
									run: (*parser).callonplugin10,
									expr: &seqExpr{
										pos: position{line: 119, col: 9, offset: 2413},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 119, col: 9, offset: 2413},
												label: "attribute",
												expr: &ruleRefExpr{
													pos:  position{line: 119, col: 19, offset: 2423},
													name: "attribute",
												},
											},
											&labeledExpr{
												pos:   position{line: 119, col: 29, offset: 2433},
												label: "attrs",
												expr: &zeroOrMoreExpr{
													pos: position{line: 119, col: 35, offset: 2439},
													expr: &actionExpr{
														pos: position{line: 120, col: 13, offset: 2453},
														run: (*parser).callonplugin16,
														expr: &seqExpr{
															pos: position{line: 120, col: 13, offset: 2453},
															exprs: []interface{}{
																&ruleRefExpr{
																	pos:  position{line: 120, col: 13, offset: 2453},
																	name: "whitespace",
																},
																&ruleRefExpr{
																	pos:  position{line: 120, col: 24, offset: 2464},
																	name: "_",
																},
																&labeledExpr{
																	pos:   position{line: 120, col: 26, offset: 2466},
																	label: "attribute",
																	expr: &ruleRefExpr{
																		pos:  position{line: 120, col: 36, offset: 2476},
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
							pos:  position{line: 126, col: 8, offset: 2618},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 127, col: 9, offset: 2630},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 127, col: 9, offset: 2630},
									val:        "}",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 127, col: 15, offset: 2636},
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
			pos:  position{line: 141, col: 1, offset: 2886},
			expr: &choiceExpr{
				pos: position{line: 142, col: 7, offset: 2899},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 142, col: 7, offset: 2899},
						run: (*parser).callonname2,
						expr: &oneOrMoreExpr{
							pos: position{line: 142, col: 8, offset: 2900},
							expr: &charClassMatcher{
								pos:        position{line: 142, col: 8, offset: 2900},
								val:        "[A-Za-z0-9_-]",
								chars:      []rune{'_', '-'},
								ranges:     []rune{'A', 'Z', 'a', 'z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
					&actionExpr{
						pos: position{line: 144, col: 9, offset: 2961},
						run: (*parser).callonname5,
						expr: &labeledExpr{
							pos:   position{line: 144, col: 9, offset: 2961},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 144, col: 15, offset: 2967},
								name: "stringValue",
							},
						},
					},
				},
			},
		},
		{
			name: "attribute",
			pos:  position{line: 153, col: 1, offset: 3115},
			expr: &actionExpr{
				pos: position{line: 154, col: 5, offset: 3131},
				run: (*parser).callonattribute1,
				expr: &seqExpr{
					pos: position{line: 154, col: 5, offset: 3131},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 154, col: 5, offset: 3131},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 154, col: 10, offset: 3136},
								name: "name",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 154, col: 15, offset: 3141},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 154, col: 17, offset: 3143},
							val:        "=>",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 154, col: 22, offset: 3148},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 154, col: 24, offset: 3150},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 154, col: 30, offset: 3156},
								name: "value",
							},
						},
					},
				},
			},
		},
		{
			name: "value",
			pos:  position{line: 162, col: 1, offset: 3293},
			expr: &choiceExpr{
				pos: position{line: 163, col: 5, offset: 3305},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 163, col: 5, offset: 3305},
						name: "plugin",
					},
					&ruleRefExpr{
						pos:  position{line: 163, col: 14, offset: 3314},
						name: "bareword",
					},
					&ruleRefExpr{
						pos:  position{line: 163, col: 25, offset: 3325},
						name: "stringValue",
					},
					&ruleRefExpr{
						pos:  position{line: 163, col: 39, offset: 3339},
						name: "number",
					},
					&ruleRefExpr{
						pos:  position{line: 163, col: 48, offset: 3348},
						name: "array",
					},
					&ruleRefExpr{
						pos:  position{line: 163, col: 56, offset: 3356},
						name: "hash",
					},
					&andCodeExpr{
						pos: position{line: 163, col: 63, offset: 3363},
						run: (*parser).callonvalue8,
					},
				},
			},
		},
		{
			name: "arrayValue",
			pos:  position{line: 171, col: 1, offset: 3500},
			expr: &choiceExpr{
				pos: position{line: 172, col: 5, offset: 3517},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 172, col: 5, offset: 3517},
						name: "bareword",
					},
					&ruleRefExpr{
						pos:  position{line: 172, col: 16, offset: 3528},
						name: "stringValue",
					},
					&ruleRefExpr{
						pos:  position{line: 172, col: 30, offset: 3542},
						name: "number",
					},
					&ruleRefExpr{
						pos:  position{line: 172, col: 39, offset: 3551},
						name: "array",
					},
					&ruleRefExpr{
						pos:  position{line: 172, col: 47, offset: 3559},
						name: "hash",
					},
					&andCodeExpr{
						pos: position{line: 172, col: 54, offset: 3566},
						run: (*parser).callonarrayValue7,
					},
				},
			},
		},
		{
			name: "bareword",
			pos:  position{line: 181, col: 1, offset: 3729},
			expr: &actionExpr{
				pos: position{line: 182, col: 5, offset: 3744},
				run: (*parser).callonbareword1,
				expr: &seqExpr{
					pos: position{line: 182, col: 5, offset: 3744},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 182, col: 5, offset: 3744},
							val:        "[A-Za-z_]",
							chars:      []rune{'_'},
							ranges:     []rune{'A', 'Z', 'a', 'z'},
							ignoreCase: false,
							inverted:   false,
						},
						&oneOrMoreExpr{
							pos: position{line: 182, col: 15, offset: 3754},
							expr: &charClassMatcher{
								pos:        position{line: 182, col: 15, offset: 3754},
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
			pos:  position{line: 190, col: 1, offset: 3964},
			expr: &actionExpr{
				pos: position{line: 191, col: 5, offset: 3989},
				run: (*parser).callondoubleQuotedString1,
				expr: &seqExpr{
					pos: position{line: 191, col: 7, offset: 3991},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 191, col: 7, offset: 3991},
							val:        "\"",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 191, col: 11, offset: 3995},
							expr: &choiceExpr{
								pos: position{line: 191, col: 13, offset: 3997},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 191, col: 13, offset: 3997},
										val:        "\\\"",
										ignoreCase: false,
									},
									&seqExpr{
										pos: position{line: 191, col: 20, offset: 4004},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 191, col: 20, offset: 4004},
												expr: &litMatcher{
													pos:        position{line: 191, col: 21, offset: 4005},
													val:        "\"",
													ignoreCase: false,
												},
											},
											&anyMatcher{
												line: 191, col: 25, offset: 4009,
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 192, col: 9, offset: 4024},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 192, col: 9, offset: 4024},
									val:        "\"",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 192, col: 15, offset: 4030},
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
			pos:  position{line: 203, col: 1, offset: 4273},
			expr: &actionExpr{
				pos: position{line: 204, col: 5, offset: 4298},
				run: (*parser).callonsingleQuotedString1,
				expr: &seqExpr{
					pos: position{line: 204, col: 7, offset: 4300},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 204, col: 7, offset: 4300},
							val:        "'",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 204, col: 11, offset: 4304},
							expr: &choiceExpr{
								pos: position{line: 204, col: 13, offset: 4306},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 204, col: 13, offset: 4306},
										val:        "\\'",
										ignoreCase: false,
									},
									&seqExpr{
										pos: position{line: 204, col: 20, offset: 4313},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 204, col: 20, offset: 4313},
												expr: &litMatcher{
													pos:        position{line: 204, col: 21, offset: 4314},
													val:        "'",
													ignoreCase: false,
												},
											},
											&anyMatcher{
												line: 204, col: 25, offset: 4318,
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 205, col: 9, offset: 4333},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 205, col: 9, offset: 4333},
									val:        "'",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 205, col: 15, offset: 4339},
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
			pos:  position{line: 216, col: 1, offset: 4547},
			expr: &actionExpr{
				pos: position{line: 217, col: 5, offset: 4565},
				run: (*parser).callonstringValue1,
				expr: &labeledExpr{
					pos:   position{line: 217, col: 5, offset: 4565},
					label: "str",
					expr: &choiceExpr{
						pos: position{line: 217, col: 11, offset: 4571},
						alternatives: []interface{}{
							&actionExpr{
								pos: position{line: 217, col: 11, offset: 4571},
								run: (*parser).callonstringValue4,
								expr: &labeledExpr{
									pos:   position{line: 217, col: 11, offset: 4571},
									label: "str",
									expr: &ruleRefExpr{
										pos:  position{line: 217, col: 15, offset: 4575},
										name: "doubleQuotedString",
									},
								},
							},
							&actionExpr{
								pos: position{line: 219, col: 9, offset: 4683},
								run: (*parser).callonstringValue7,
								expr: &labeledExpr{
									pos:   position{line: 219, col: 9, offset: 4683},
									label: "str",
									expr: &ruleRefExpr{
										pos:  position{line: 219, col: 13, offset: 4687},
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
			pos:  position{line: 229, col: 1, offset: 4925},
			expr: &actionExpr{
				pos: position{line: 230, col: 5, offset: 4938},
				run: (*parser).callonregexp1,
				expr: &seqExpr{
					pos: position{line: 230, col: 7, offset: 4940},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 230, col: 7, offset: 4940},
							val:        "/",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 230, col: 11, offset: 4944},
							expr: &choiceExpr{
								pos: position{line: 230, col: 13, offset: 4946},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 230, col: 13, offset: 4946},
										val:        "\\/",
										ignoreCase: false,
									},
									&seqExpr{
										pos: position{line: 230, col: 20, offset: 4953},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 230, col: 20, offset: 4953},
												expr: &litMatcher{
													pos:        position{line: 230, col: 21, offset: 4954},
													val:        "/",
													ignoreCase: false,
												},
											},
											&anyMatcher{
												line: 230, col: 25, offset: 4958,
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 231, col: 9, offset: 4973},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 231, col: 9, offset: 4973},
									val:        "/",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 231, col: 15, offset: 4979},
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
			pos:  position{line: 243, col: 1, offset: 5205},
			expr: &actionExpr{
				pos: position{line: 244, col: 5, offset: 5218},
				run: (*parser).callonnumber1,
				expr: &seqExpr{
					pos: position{line: 244, col: 5, offset: 5218},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 244, col: 5, offset: 5218},
							expr: &litMatcher{
								pos:        position{line: 244, col: 5, offset: 5218},
								val:        "-",
								ignoreCase: false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 244, col: 10, offset: 5223},
							expr: &charClassMatcher{
								pos:        position{line: 244, col: 10, offset: 5223},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 244, col: 17, offset: 5230},
							expr: &seqExpr{
								pos: position{line: 244, col: 18, offset: 5231},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 244, col: 18, offset: 5231},
										val:        ".",
										ignoreCase: false,
									},
									&zeroOrMoreExpr{
										pos: position{line: 244, col: 22, offset: 5235},
										expr: &charClassMatcher{
											pos:        position{line: 244, col: 22, offset: 5235},
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
			pos:  position{line: 260, col: 1, offset: 5552},
			expr: &actionExpr{
				pos: position{line: 261, col: 5, offset: 5564},
				run: (*parser).callonarray1,
				expr: &seqExpr{
					pos: position{line: 261, col: 5, offset: 5564},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 261, col: 5, offset: 5564},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 261, col: 9, offset: 5568},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 261, col: 11, offset: 5570},
							label: "value",
							expr: &zeroOrOneExpr{
								pos: position{line: 261, col: 17, offset: 5576},
								expr: &actionExpr{
									pos: position{line: 262, col: 9, offset: 5586},
									run: (*parser).callonarray7,
									expr: &seqExpr{
										pos: position{line: 262, col: 9, offset: 5586},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 262, col: 9, offset: 5586},
												label: "value",
												expr: &ruleRefExpr{
													pos:  position{line: 262, col: 15, offset: 5592},
													name: "value",
												},
											},
											&labeledExpr{
												pos:   position{line: 262, col: 21, offset: 5598},
												label: "values",
												expr: &zeroOrMoreExpr{
													pos: position{line: 262, col: 28, offset: 5605},
													expr: &actionExpr{
														pos: position{line: 263, col: 13, offset: 5619},
														run: (*parser).callonarray13,
														expr: &seqExpr{
															pos: position{line: 263, col: 13, offset: 5619},
															exprs: []interface{}{
																&ruleRefExpr{
																	pos:  position{line: 263, col: 13, offset: 5619},
																	name: "_",
																},
																&litMatcher{
																	pos:        position{line: 263, col: 15, offset: 5621},
																	val:        ",",
																	ignoreCase: false,
																},
																&ruleRefExpr{
																	pos:  position{line: 263, col: 19, offset: 5625},
																	name: "_",
																},
																&labeledExpr{
																	pos:   position{line: 263, col: 21, offset: 5627},
																	label: "value",
																	expr: &ruleRefExpr{
																		pos:  position{line: 263, col: 27, offset: 5633},
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
							pos:  position{line: 269, col: 8, offset: 5764},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 270, col: 9, offset: 5776},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 270, col: 9, offset: 5776},
									val:        "]",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 270, col: 15, offset: 5782},
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
			pos:  position{line: 286, col: 1, offset: 6032},
			expr: &actionExpr{
				pos: position{line: 287, col: 5, offset: 6043},
				run: (*parser).callonhash1,
				expr: &seqExpr{
					pos: position{line: 287, col: 5, offset: 6043},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 287, col: 5, offset: 6043},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 287, col: 9, offset: 6047},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 287, col: 11, offset: 6049},
							label: "entries",
							expr: &zeroOrOneExpr{
								pos: position{line: 287, col: 19, offset: 6057},
								expr: &ruleRefExpr{
									pos:  position{line: 287, col: 19, offset: 6057},
									name: "hashentries",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 287, col: 32, offset: 6070},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 288, col: 9, offset: 6082},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 288, col: 9, offset: 6082},
									val:        "}",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 288, col: 15, offset: 6088},
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
			pos:  position{line: 300, col: 1, offset: 6327},
			expr: &actionExpr{
				pos: position{line: 301, col: 5, offset: 6345},
				run: (*parser).callonhashentries1,
				expr: &seqExpr{
					pos: position{line: 301, col: 5, offset: 6345},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 301, col: 5, offset: 6345},
							label: "hashentry",
							expr: &ruleRefExpr{
								pos:  position{line: 301, col: 15, offset: 6355},
								name: "hashentry",
							},
						},
						&labeledExpr{
							pos:   position{line: 301, col: 25, offset: 6365},
							label: "hashentries1",
							expr: &zeroOrMoreExpr{
								pos: position{line: 301, col: 38, offset: 6378},
								expr: &actionExpr{
									pos: position{line: 302, col: 9, offset: 6388},
									run: (*parser).callonhashentries7,
									expr: &seqExpr{
										pos: position{line: 302, col: 9, offset: 6388},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 302, col: 9, offset: 6388},
												name: "whitespace",
											},
											&labeledExpr{
												pos:   position{line: 302, col: 20, offset: 6399},
												label: "hashentry",
												expr: &ruleRefExpr{
													pos:  position{line: 302, col: 30, offset: 6409},
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
			pos:  position{line: 314, col: 1, offset: 6661},
			expr: &actionExpr{
				pos: position{line: 315, col: 5, offset: 6677},
				run: (*parser).callonhashentry1,
				expr: &seqExpr{
					pos: position{line: 315, col: 5, offset: 6677},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 315, col: 5, offset: 6677},
							label: "name",
							expr: &choiceExpr{
								pos: position{line: 315, col: 11, offset: 6683},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 315, col: 11, offset: 6683},
										name: "number",
									},
									&ruleRefExpr{
										pos:  position{line: 315, col: 20, offset: 6692},
										name: "bareword",
									},
									&ruleRefExpr{
										pos:  position{line: 315, col: 31, offset: 6703},
										name: "stringValue",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 315, col: 44, offset: 6716},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 315, col: 46, offset: 6718},
							val:        "=>",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 315, col: 51, offset: 6723},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 315, col: 53, offset: 6725},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 315, col: 59, offset: 6731},
								name: "value",
							},
						},
					},
				},
			},
		},
		{
			name: "branch",
			pos:  position{line: 326, col: 1, offset: 6898},
			expr: &actionExpr{
				pos: position{line: 327, col: 5, offset: 6911},
				run: (*parser).callonbranch1,
				expr: &seqExpr{
					pos: position{line: 327, col: 5, offset: 6911},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 327, col: 5, offset: 6911},
							label: "ifBlock",
							expr: &ruleRefExpr{
								pos:  position{line: 327, col: 13, offset: 6919},
								name: "ifCond",
							},
						},
						&labeledExpr{
							pos:   position{line: 327, col: 20, offset: 6926},
							label: "elseIfBlocks",
							expr: &zeroOrMoreExpr{
								pos: position{line: 327, col: 33, offset: 6939},
								expr: &actionExpr{
									pos: position{line: 328, col: 9, offset: 6949},
									run: (*parser).callonbranch7,
									expr: &seqExpr{
										pos: position{line: 328, col: 9, offset: 6949},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 328, col: 9, offset: 6949},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 328, col: 11, offset: 6951},
												label: "eib",
												expr: &ruleRefExpr{
													pos:  position{line: 328, col: 15, offset: 6955},
													name: "elseIf",
												},
											},
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 331, col: 12, offset: 7021},
							label: "elseBlock",
							expr: &zeroOrOneExpr{
								pos: position{line: 331, col: 22, offset: 7031},
								expr: &actionExpr{
									pos: position{line: 332, col: 13, offset: 7045},
									run: (*parser).callonbranch14,
									expr: &seqExpr{
										pos: position{line: 332, col: 13, offset: 7045},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 332, col: 13, offset: 7045},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 332, col: 15, offset: 7047},
												label: "eb",
												expr: &ruleRefExpr{
													pos:  position{line: 332, col: 18, offset: 7050},
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
			pos:  position{line: 344, col: 1, offset: 7298},
			expr: &actionExpr{
				pos: position{line: 345, col: 5, offset: 7311},
				run: (*parser).callonifCond1,
				expr: &seqExpr{
					pos: position{line: 345, col: 5, offset: 7311},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 345, col: 5, offset: 7311},
							val:        "if",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 345, col: 10, offset: 7316},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 345, col: 12, offset: 7318},
							label: "cond",
							expr: &ruleRefExpr{
								pos:  position{line: 345, col: 17, offset: 7323},
								name: "condition",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 345, col: 27, offset: 7333},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 345, col: 29, offset: 7335},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 345, col: 33, offset: 7339},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 345, col: 35, offset: 7341},
							label: "bops",
							expr: &zeroOrMoreExpr{
								pos: position{line: 345, col: 40, offset: 7346},
								expr: &actionExpr{
									pos: position{line: 346, col: 13, offset: 7360},
									run: (*parser).callonifCond12,
									expr: &seqExpr{
										pos: position{line: 346, col: 13, offset: 7360},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 346, col: 13, offset: 7360},
												label: "bop",
												expr: &ruleRefExpr{
													pos:  position{line: 346, col: 17, offset: 7364},
													name: "branchOrPlugin",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 346, col: 32, offset: 7379},
												name: "_",
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 350, col: 13, offset: 7454},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 350, col: 13, offset: 7454},
									val:        "}",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 350, col: 19, offset: 7460},
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
			pos:  position{line: 362, col: 1, offset: 7734},
			expr: &actionExpr{
				pos: position{line: 363, col: 5, offset: 7747},
				run: (*parser).callonelseIf1,
				expr: &seqExpr{
					pos: position{line: 363, col: 5, offset: 7747},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 363, col: 5, offset: 7747},
							val:        "else",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 363, col: 12, offset: 7754},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 363, col: 14, offset: 7756},
							val:        "if",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 363, col: 19, offset: 7761},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 363, col: 21, offset: 7763},
							label: "cond",
							expr: &ruleRefExpr{
								pos:  position{line: 363, col: 26, offset: 7768},
								name: "condition",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 363, col: 36, offset: 7778},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 363, col: 38, offset: 7780},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 363, col: 42, offset: 7784},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 363, col: 44, offset: 7786},
							label: "bops",
							expr: &zeroOrMoreExpr{
								pos: position{line: 363, col: 49, offset: 7791},
								expr: &actionExpr{
									pos: position{line: 364, col: 9, offset: 7801},
									run: (*parser).callonelseIf14,
									expr: &seqExpr{
										pos: position{line: 364, col: 9, offset: 7801},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 364, col: 9, offset: 7801},
												label: "bop",
												expr: &ruleRefExpr{
													pos:  position{line: 364, col: 13, offset: 7805},
													name: "branchOrPlugin",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 364, col: 28, offset: 7820},
												name: "_",
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 368, col: 9, offset: 7879},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 368, col: 9, offset: 7879},
									val:        "}",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 368, col: 15, offset: 7885},
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
			pos:  position{line: 380, col: 1, offset: 8127},
			expr: &actionExpr{
				pos: position{line: 381, col: 5, offset: 8142},
				run: (*parser).callonelseCond1,
				expr: &seqExpr{
					pos: position{line: 381, col: 5, offset: 8142},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 381, col: 5, offset: 8142},
							val:        "else",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 381, col: 12, offset: 8149},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 381, col: 14, offset: 8151},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 381, col: 18, offset: 8155},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 381, col: 20, offset: 8157},
							label: "bops",
							expr: &zeroOrMoreExpr{
								pos: position{line: 381, col: 25, offset: 8162},
								expr: &actionExpr{
									pos: position{line: 382, col: 9, offset: 8172},
									run: (*parser).callonelseCond9,
									expr: &seqExpr{
										pos: position{line: 382, col: 9, offset: 8172},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 382, col: 9, offset: 8172},
												label: "bop",
												expr: &ruleRefExpr{
													pos:  position{line: 382, col: 13, offset: 8176},
													name: "branchOrPlugin",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 382, col: 28, offset: 8191},
												name: "_",
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 386, col: 9, offset: 8250},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 386, col: 9, offset: 8250},
									val:        "}",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 386, col: 15, offset: 8256},
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
			pos:  position{line: 398, col: 1, offset: 8505},
			expr: &actionExpr{
				pos: position{line: 399, col: 5, offset: 8521},
				run: (*parser).calloncondition1,
				expr: &seqExpr{
					pos: position{line: 399, col: 5, offset: 8521},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 399, col: 5, offset: 8521},
							label: "cond",
							expr: &ruleRefExpr{
								pos:  position{line: 399, col: 10, offset: 8526},
								name: "expression",
							},
						},
						&labeledExpr{
							pos:   position{line: 399, col: 21, offset: 8537},
							label: "conds",
							expr: &zeroOrMoreExpr{
								pos: position{line: 399, col: 27, offset: 8543},
								expr: &actionExpr{
									pos: position{line: 400, col: 9, offset: 8553},
									run: (*parser).calloncondition7,
									expr: &seqExpr{
										pos: position{line: 400, col: 9, offset: 8553},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 400, col: 9, offset: 8553},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 400, col: 11, offset: 8555},
												label: "bo",
												expr: &ruleRefExpr{
													pos:  position{line: 400, col: 14, offset: 8558},
													name: "booleanOperator",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 400, col: 30, offset: 8574},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 400, col: 32, offset: 8576},
												label: "cond",
												expr: &ruleRefExpr{
													pos:  position{line: 400, col: 37, offset: 8581},
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
			pos:  position{line: 419, col: 1, offset: 8980},
			expr: &choiceExpr{
				pos: position{line: 421, col: 9, offset: 9007},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 422, col: 13, offset: 9021},
						run: (*parser).callonexpression2,
						expr: &seqExpr{
							pos: position{line: 422, col: 13, offset: 9021},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 422, col: 13, offset: 9021},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 422, col: 17, offset: 9025},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 422, col: 19, offset: 9027},
									label: "cond",
									expr: &ruleRefExpr{
										pos:  position{line: 422, col: 24, offset: 9032},
										name: "condition",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 422, col: 34, offset: 9042},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 422, col: 36, offset: 9044},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 426, col: 9, offset: 9131},
						name: "negativeExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 427, col: 9, offset: 9158},
						name: "inExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 428, col: 9, offset: 9179},
						name: "notInExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 429, col: 9, offset: 9203},
						name: "compareExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 430, col: 9, offset: 9229},
						name: "regexpExpression",
					},
					&actionExpr{
						pos: position{line: 431, col: 9, offset: 9254},
						run: (*parser).callonexpression15,
						expr: &labeledExpr{
							pos:   position{line: 431, col: 9, offset: 9254},
							label: "rv",
							expr: &ruleRefExpr{
								pos:  position{line: 431, col: 12, offset: 9257},
								name: "rvalue",
							},
						},
					},
				},
			},
		},
		{
			name: "negativeExpression",
			pos:  position{line: 444, col: 1, offset: 9541},
			expr: &choiceExpr{
				pos: position{line: 446, col: 9, offset: 9576},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 447, col: 13, offset: 9590},
						run: (*parser).callonnegativeExpression2,
						expr: &seqExpr{
							pos: position{line: 447, col: 13, offset: 9590},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 447, col: 13, offset: 9590},
									val:        "!",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 447, col: 17, offset: 9594},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 447, col: 19, offset: 9596},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 447, col: 23, offset: 9600},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 447, col: 25, offset: 9602},
									label: "cond",
									expr: &ruleRefExpr{
										pos:  position{line: 447, col: 30, offset: 9607},
										name: "condition",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 447, col: 40, offset: 9617},
									name: "_",
								},
								&choiceExpr{
									pos: position{line: 448, col: 17, offset: 9637},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 448, col: 17, offset: 9637},
											val:        ")",
											ignoreCase: false,
										},
										&andCodeExpr{
											pos: position{line: 448, col: 23, offset: 9643},
											run: (*parser).callonnegativeExpression13,
										},
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 456, col: 11, offset: 9843},
						run: (*parser).callonnegativeExpression14,
						expr: &seqExpr{
							pos: position{line: 456, col: 11, offset: 9843},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 456, col: 11, offset: 9843},
									val:        "!",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 456, col: 15, offset: 9847},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 456, col: 17, offset: 9849},
									label: "sel",
									expr: &ruleRefExpr{
										pos:  position{line: 456, col: 21, offset: 9853},
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
			pos:  position{line: 467, col: 1, offset: 10052},
			expr: &actionExpr{
				pos: position{line: 468, col: 5, offset: 10071},
				run: (*parser).calloninExpression1,
				expr: &seqExpr{
					pos: position{line: 468, col: 5, offset: 10071},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 468, col: 5, offset: 10071},
							label: "lv",
							expr: &ruleRefExpr{
								pos:  position{line: 468, col: 8, offset: 10074},
								name: "rvalue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 468, col: 15, offset: 10081},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 468, col: 17, offset: 10083},
							name: "inOperator",
						},
						&ruleRefExpr{
							pos:  position{line: 468, col: 28, offset: 10094},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 468, col: 30, offset: 10096},
							label: "rv",
							expr: &ruleRefExpr{
								pos:  position{line: 468, col: 33, offset: 10099},
								name: "rvalue",
							},
						},
					},
				},
			},
		},
		{
			name: "notInExpression",
			pos:  position{line: 477, col: 1, offset: 10278},
			expr: &actionExpr{
				pos: position{line: 478, col: 5, offset: 10300},
				run: (*parser).callonnotInExpression1,
				expr: &seqExpr{
					pos: position{line: 478, col: 5, offset: 10300},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 478, col: 5, offset: 10300},
							label: "lv",
							expr: &ruleRefExpr{
								pos:  position{line: 478, col: 8, offset: 10303},
								name: "rvalue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 478, col: 15, offset: 10310},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 478, col: 17, offset: 10312},
							name: "notInOperator",
						},
						&ruleRefExpr{
							pos:  position{line: 478, col: 31, offset: 10326},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 478, col: 33, offset: 10328},
							label: "rv",
							expr: &ruleRefExpr{
								pos:  position{line: 478, col: 36, offset: 10331},
								name: "rvalue",
							},
						},
					},
				},
			},
		},
		{
			name: "inOperator",
			pos:  position{line: 486, col: 1, offset: 10430},
			expr: &choiceExpr{
				pos: position{line: 487, col: 5, offset: 10447},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 487, col: 5, offset: 10447},
						val:        "in",
						ignoreCase: false,
					},
					&andCodeExpr{
						pos: position{line: 487, col: 12, offset: 10454},
						run: (*parser).calloninOperator3,
					},
				},
			},
		},
		{
			name: "notInOperator",
			pos:  position{line: 495, col: 1, offset: 10576},
			expr: &choiceExpr{
				pos: position{line: 496, col: 5, offset: 10596},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 496, col: 5, offset: 10596},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 496, col: 5, offset: 10596},
								val:        "not ",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 496, col: 12, offset: 10603},
								name: "_",
							},
							&litMatcher{
								pos:        position{line: 496, col: 14, offset: 10605},
								val:        "in",
								ignoreCase: false,
							},
						},
					},
					&andCodeExpr{
						pos: position{line: 496, col: 21, offset: 10612},
						run: (*parser).callonnotInOperator6,
					},
				},
			},
		},
		{
			name: "rvalue",
			pos:  position{line: 507, col: 1, offset: 10938},
			expr: &choiceExpr{
				pos: position{line: 508, col: 5, offset: 10951},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 508, col: 5, offset: 10951},
						name: "stringValue",
					},
					&ruleRefExpr{
						pos:  position{line: 508, col: 19, offset: 10965},
						name: "number",
					},
					&ruleRefExpr{
						pos:  position{line: 508, col: 28, offset: 10974},
						name: "selector",
					},
					&ruleRefExpr{
						pos:  position{line: 508, col: 39, offset: 10985},
						name: "array",
					},
					&ruleRefExpr{
						pos:  position{line: 508, col: 47, offset: 10993},
						name: "regexp",
					},
					&andCodeExpr{
						pos: position{line: 508, col: 56, offset: 11002},
						run: (*parser).callonrvalue7,
					},
				},
			},
		},
		{
			name: "compareExpression",
			pos:  position{line: 542, col: 1, offset: 11740},
			expr: &actionExpr{
				pos: position{line: 543, col: 5, offset: 11764},
				run: (*parser).calloncompareExpression1,
				expr: &seqExpr{
					pos: position{line: 543, col: 5, offset: 11764},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 543, col: 5, offset: 11764},
							label: "lv",
							expr: &ruleRefExpr{
								pos:  position{line: 543, col: 8, offset: 11767},
								name: "rvalue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 543, col: 15, offset: 11774},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 543, col: 17, offset: 11776},
							label: "co",
							expr: &ruleRefExpr{
								pos:  position{line: 543, col: 20, offset: 11779},
								name: "compareOperator",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 543, col: 36, offset: 11795},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 543, col: 38, offset: 11797},
							label: "rv",
							expr: &ruleRefExpr{
								pos:  position{line: 543, col: 41, offset: 11800},
								name: "rvalue",
							},
						},
					},
				},
			},
		},
		{
			name: "compareOperator",
			pos:  position{line: 552, col: 1, offset: 11996},
			expr: &actionExpr{
				pos: position{line: 553, col: 5, offset: 12018},
				run: (*parser).calloncompareOperator1,
				expr: &choiceExpr{
					pos: position{line: 553, col: 6, offset: 12019},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 553, col: 6, offset: 12019},
							val:        "==",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 553, col: 13, offset: 12026},
							val:        "!=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 553, col: 20, offset: 12033},
							val:        "<=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 553, col: 27, offset: 12040},
							val:        ">=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 553, col: 34, offset: 12047},
							val:        "<",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 553, col: 40, offset: 12053},
							val:        ">",
							ignoreCase: false,
						},
						&andCodeExpr{
							pos: position{line: 553, col: 46, offset: 12059},
							run: (*parser).calloncompareOperator9,
						},
					},
				},
			},
		},
		{
			name: "regexpExpression",
			pos:  position{line: 564, col: 1, offset: 12343},
			expr: &actionExpr{
				pos: position{line: 565, col: 5, offset: 12366},
				run: (*parser).callonregexpExpression1,
				expr: &seqExpr{
					pos: position{line: 565, col: 5, offset: 12366},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 565, col: 5, offset: 12366},
							label: "lv",
							expr: &ruleRefExpr{
								pos:  position{line: 565, col: 8, offset: 12369},
								name: "rvalue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 565, col: 15, offset: 12376},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 565, col: 18, offset: 12379},
							label: "ro",
							expr: &ruleRefExpr{
								pos:  position{line: 565, col: 21, offset: 12382},
								name: "regexpOperator",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 565, col: 36, offset: 12397},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 565, col: 38, offset: 12399},
							label: "rv",
							expr: &choiceExpr{
								pos: position{line: 565, col: 42, offset: 12403},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 565, col: 42, offset: 12403},
										name: "stringValue",
									},
									&ruleRefExpr{
										pos:  position{line: 565, col: 56, offset: 12417},
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
			pos:  position{line: 573, col: 1, offset: 12575},
			expr: &actionExpr{
				pos: position{line: 574, col: 5, offset: 12596},
				run: (*parser).callonregexpOperator1,
				expr: &choiceExpr{
					pos: position{line: 574, col: 6, offset: 12597},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 574, col: 6, offset: 12597},
							val:        "=~",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 574, col: 13, offset: 12604},
							val:        "!~",
							ignoreCase: false,
						},
						&andCodeExpr{
							pos: position{line: 574, col: 20, offset: 12611},
							run: (*parser).callonregexpOperator5,
						},
					},
				},
			},
		},
		{
			name: "booleanOperator",
			pos:  position{line: 585, col: 1, offset: 12874},
			expr: &actionExpr{
				pos: position{line: 586, col: 5, offset: 12896},
				run: (*parser).callonbooleanOperator1,
				expr: &choiceExpr{
					pos: position{line: 586, col: 6, offset: 12897},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 586, col: 6, offset: 12897},
							val:        "and",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 586, col: 14, offset: 12905},
							val:        "or",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 586, col: 21, offset: 12912},
							val:        "xor",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 586, col: 29, offset: 12920},
							val:        "nand",
							ignoreCase: false,
						},
						&andCodeExpr{
							pos: position{line: 586, col: 38, offset: 12929},
							run: (*parser).callonbooleanOperator7,
						},
					},
				},
			},
		},
		{
			name: "selector",
			pos:  position{line: 597, col: 1, offset: 13166},
			expr: &actionExpr{
				pos: position{line: 598, col: 5, offset: 13181},
				run: (*parser).callonselector1,
				expr: &labeledExpr{
					pos:   position{line: 598, col: 5, offset: 13181},
					label: "ses",
					expr: &oneOrMoreExpr{
						pos: position{line: 598, col: 9, offset: 13185},
						expr: &ruleRefExpr{
							pos:  position{line: 598, col: 9, offset: 13185},
							name: "selectorElement",
						},
					},
				},
			},
		},
		{
			name: "selectorElement",
			pos:  position{line: 607, col: 1, offset: 13348},
			expr: &actionExpr{
				pos: position{line: 608, col: 5, offset: 13370},
				run: (*parser).callonselectorElement1,
				expr: &seqExpr{
					pos: position{line: 608, col: 5, offset: 13370},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 608, col: 5, offset: 13370},
							val:        "[",
							ignoreCase: false,
						},
						&oneOrMoreExpr{
							pos: position{line: 608, col: 9, offset: 13374},
							expr: &charClassMatcher{
								pos:        position{line: 608, col: 9, offset: 13374},
								val:        "[^\\],]",
								chars:      []rune{']', ','},
								ignoreCase: false,
								inverted:   true,
							},
						},
						&choiceExpr{
							pos: position{line: 609, col: 9, offset: 13392},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 609, col: 9, offset: 13392},
									val:        "]",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 609, col: 15, offset: 13398},
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
			pos:  position{line: 616, col: 1, offset: 13538},
			expr: &notExpr{
				pos: position{line: 616, col: 7, offset: 13544},
				expr: &anyMatcher{
					line: 616, col: 8, offset: 13545,
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
