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
								name: "plugin_section",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 19, col: 25, offset: 457},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 19, col: 27, offset: 459},
							label: "pss",
							expr: &zeroOrMoreExpr{
								pos: position{line: 19, col: 31, offset: 463},
								expr: &actionExpr{
									pos: position{line: 20, col: 9, offset: 473},
									run: (*parser).callonconfig9,
									expr: &seqExpr{
										pos: position{line: 20, col: 9, offset: 473},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 20, col: 9, offset: 473},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 20, col: 11, offset: 475},
												label: "ps",
												expr: &ruleRefExpr{
													pos:  position{line: 20, col: 14, offset: 478},
													name: "plugin_section",
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 23, col: 8, offset: 541},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 23, col: 10, offset: 543},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "comment",
			pos:  position{line: 31, col: 1, offset: 693},
			expr: &oneOrMoreExpr{
				pos: position{line: 32, col: 5, offset: 707},
				expr: &seqExpr{
					pos: position{line: 32, col: 6, offset: 708},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 32, col: 6, offset: 708},
							expr: &ruleRefExpr{
								pos:  position{line: 32, col: 6, offset: 708},
								name: "whitespace",
							},
						},
						&litMatcher{
							pos:        position{line: 32, col: 18, offset: 720},
							val:        "#",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 32, col: 22, offset: 724},
							expr: &charClassMatcher{
								pos:        position{line: 32, col: 22, offset: 724},
								val:        "[^\\r\\n]",
								chars:      []rune{'\r', '\n'},
								ignoreCase: false,
								inverted:   true,
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 32, col: 31, offset: 733},
							expr: &litMatcher{
								pos:        position{line: 32, col: 31, offset: 733},
								val:        "\r",
								ignoreCase: false,
							},
						},
						&litMatcher{
							pos:        position{line: 32, col: 37, offset: 739},
							val:        "\n",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 38, col: 1, offset: 833},
			expr: &zeroOrMoreExpr{
				pos: position{line: 39, col: 5, offset: 841},
				expr: &choiceExpr{
					pos: position{line: 39, col: 6, offset: 842},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 39, col: 6, offset: 842},
							name: "comment",
						},
						&ruleRefExpr{
							pos:  position{line: 39, col: 16, offset: 852},
							name: "whitespace",
						},
					},
				},
			},
		},
		{
			name: "whitespace",
			pos:  position{line: 45, col: 1, offset: 948},
			expr: &oneOrMoreExpr{
				pos: position{line: 46, col: 5, offset: 965},
				expr: &charClassMatcher{
					pos:        position{line: 46, col: 5, offset: 965},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "plugin_section",
			pos:  position{line: 55, col: 1, offset: 1122},
			expr: &actionExpr{
				pos: position{line: 56, col: 5, offset: 1143},
				run: (*parser).callonplugin_section1,
				expr: &seqExpr{
					pos: position{line: 56, col: 5, offset: 1143},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 56, col: 5, offset: 1143},
							label: "pt",
							expr: &ruleRefExpr{
								pos:  position{line: 56, col: 8, offset: 1146},
								name: "plugin_type",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 56, col: 20, offset: 1158},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 56, col: 22, offset: 1160},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 56, col: 26, offset: 1164},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 56, col: 28, offset: 1166},
							label: "bops",
							expr: &zeroOrMoreExpr{
								pos: position{line: 56, col: 33, offset: 1171},
								expr: &actionExpr{
									pos: position{line: 57, col: 9, offset: 1181},
									run: (*parser).callonplugin_section10,
									expr: &seqExpr{
										pos: position{line: 57, col: 9, offset: 1181},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 57, col: 9, offset: 1181},
												label: "bop",
												expr: &ruleRefExpr{
													pos:  position{line: 57, col: 13, offset: 1185},
													name: "branch_or_plugin",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 57, col: 30, offset: 1202},
												name: "_",
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 61, col: 9, offset: 1264},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 61, col: 9, offset: 1264},
									val:        "}",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 61, col: 15, offset: 1270},
									run: (*parser).callonplugin_section17,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "branch_or_plugin",
			pos:  position{line: 72, col: 1, offset: 1462},
			expr: &choiceExpr{
				pos: position{line: 73, col: 5, offset: 1485},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 73, col: 5, offset: 1485},
						name: "branch",
					},
					&ruleRefExpr{
						pos:  position{line: 73, col: 14, offset: 1494},
						name: "plugin",
					},
				},
			},
		},
		{
			name: "plugin_type",
			pos:  position{line: 79, col: 1, offset: 1573},
			expr: &choiceExpr{
				pos: position{line: 80, col: 5, offset: 1591},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 80, col: 5, offset: 1591},
						run: (*parser).callonplugin_type2,
						expr: &litMatcher{
							pos:        position{line: 80, col: 5, offset: 1591},
							val:        "input",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 82, col: 9, offset: 1639},
						run: (*parser).callonplugin_type4,
						expr: &litMatcher{
							pos:        position{line: 82, col: 9, offset: 1639},
							val:        "filter",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 84, col: 9, offset: 1689},
						run: (*parser).callonplugin_type6,
						expr: &litMatcher{
							pos:        position{line: 84, col: 9, offset: 1689},
							val:        "output",
							ignoreCase: false,
						},
					},
					&andCodeExpr{
						pos: position{line: 86, col: 9, offset: 1739},
						run: (*parser).callonplugin_type8,
					},
				},
			},
		},
		{
			name: "plugin",
			pos:  position{line: 117, col: 1, offset: 2378},
			expr: &actionExpr{
				pos: position{line: 118, col: 5, offset: 2391},
				run: (*parser).callonplugin1,
				expr: &seqExpr{
					pos: position{line: 118, col: 5, offset: 2391},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 118, col: 5, offset: 2391},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 118, col: 10, offset: 2396},
								name: "name",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 118, col: 15, offset: 2401},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 118, col: 17, offset: 2403},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 118, col: 21, offset: 2407},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 118, col: 23, offset: 2409},
							label: "attributes",
							expr: &zeroOrOneExpr{
								pos: position{line: 118, col: 34, offset: 2420},
								expr: &actionExpr{
									pos: position{line: 119, col: 9, offset: 2431},
									run: (*parser).callonplugin10,
									expr: &seqExpr{
										pos: position{line: 119, col: 9, offset: 2431},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 119, col: 9, offset: 2431},
												label: "attribute",
												expr: &ruleRefExpr{
													pos:  position{line: 119, col: 19, offset: 2441},
													name: "attribute",
												},
											},
											&labeledExpr{
												pos:   position{line: 119, col: 29, offset: 2451},
												label: "attrs",
												expr: &zeroOrMoreExpr{
													pos: position{line: 119, col: 35, offset: 2457},
													expr: &actionExpr{
														pos: position{line: 120, col: 13, offset: 2471},
														run: (*parser).callonplugin16,
														expr: &seqExpr{
															pos: position{line: 120, col: 13, offset: 2471},
															exprs: []interface{}{
																&ruleRefExpr{
																	pos:  position{line: 120, col: 13, offset: 2471},
																	name: "whitespace",
																},
																&ruleRefExpr{
																	pos:  position{line: 120, col: 24, offset: 2482},
																	name: "_",
																},
																&labeledExpr{
																	pos:   position{line: 120, col: 26, offset: 2484},
																	label: "attribute",
																	expr: &ruleRefExpr{
																		pos:  position{line: 120, col: 36, offset: 2494},
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
							pos:  position{line: 126, col: 8, offset: 2636},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 127, col: 9, offset: 2648},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 127, col: 9, offset: 2648},
									val:        "}",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 127, col: 15, offset: 2654},
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
			pos:  position{line: 141, col: 1, offset: 2904},
			expr: &choiceExpr{
				pos: position{line: 142, col: 7, offset: 2917},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 142, col: 7, offset: 2917},
						run: (*parser).callonname2,
						expr: &oneOrMoreExpr{
							pos: position{line: 142, col: 8, offset: 2918},
							expr: &charClassMatcher{
								pos:        position{line: 142, col: 8, offset: 2918},
								val:        "[A-Za-z0-9_-]",
								chars:      []rune{'_', '-'},
								ranges:     []rune{'A', 'Z', 'a', 'z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
					&actionExpr{
						pos: position{line: 144, col: 9, offset: 2979},
						run: (*parser).callonname5,
						expr: &labeledExpr{
							pos:   position{line: 144, col: 9, offset: 2979},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 144, col: 15, offset: 2985},
								name: "string_value",
							},
						},
					},
				},
			},
		},
		{
			name: "attribute",
			pos:  position{line: 153, col: 1, offset: 3134},
			expr: &actionExpr{
				pos: position{line: 154, col: 5, offset: 3150},
				run: (*parser).callonattribute1,
				expr: &seqExpr{
					pos: position{line: 154, col: 5, offset: 3150},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 154, col: 5, offset: 3150},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 154, col: 10, offset: 3155},
								name: "name",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 154, col: 15, offset: 3160},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 154, col: 17, offset: 3162},
							val:        "=>",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 154, col: 22, offset: 3167},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 154, col: 24, offset: 3169},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 154, col: 30, offset: 3175},
								name: "value",
							},
						},
					},
				},
			},
		},
		{
			name: "value",
			pos:  position{line: 162, col: 1, offset: 3312},
			expr: &choiceExpr{
				pos: position{line: 163, col: 5, offset: 3324},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 163, col: 5, offset: 3324},
						name: "plugin",
					},
					&ruleRefExpr{
						pos:  position{line: 163, col: 14, offset: 3333},
						name: "bareword",
					},
					&ruleRefExpr{
						pos:  position{line: 163, col: 25, offset: 3344},
						name: "string_value",
					},
					&ruleRefExpr{
						pos:  position{line: 163, col: 40, offset: 3359},
						name: "number",
					},
					&ruleRefExpr{
						pos:  position{line: 163, col: 49, offset: 3368},
						name: "array",
					},
					&ruleRefExpr{
						pos:  position{line: 163, col: 57, offset: 3376},
						name: "hash",
					},
					&andCodeExpr{
						pos: position{line: 163, col: 64, offset: 3383},
						run: (*parser).callonvalue8,
					},
				},
			},
		},
		{
			name: "array_value",
			pos:  position{line: 171, col: 1, offset: 3520},
			expr: &choiceExpr{
				pos: position{line: 172, col: 5, offset: 3538},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 172, col: 5, offset: 3538},
						name: "bareword",
					},
					&ruleRefExpr{
						pos:  position{line: 172, col: 16, offset: 3549},
						name: "string_value",
					},
					&ruleRefExpr{
						pos:  position{line: 172, col: 31, offset: 3564},
						name: "number",
					},
					&ruleRefExpr{
						pos:  position{line: 172, col: 40, offset: 3573},
						name: "array",
					},
					&ruleRefExpr{
						pos:  position{line: 172, col: 48, offset: 3581},
						name: "hash",
					},
					&andCodeExpr{
						pos: position{line: 172, col: 55, offset: 3588},
						run: (*parser).callonarray_value7,
					},
				},
			},
		},
		{
			name: "bareword",
			pos:  position{line: 181, col: 1, offset: 3751},
			expr: &actionExpr{
				pos: position{line: 182, col: 5, offset: 3766},
				run: (*parser).callonbareword1,
				expr: &seqExpr{
					pos: position{line: 182, col: 5, offset: 3766},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 182, col: 5, offset: 3766},
							val:        "[A-Za-z_]",
							chars:      []rune{'_'},
							ranges:     []rune{'A', 'Z', 'a', 'z'},
							ignoreCase: false,
							inverted:   false,
						},
						&oneOrMoreExpr{
							pos: position{line: 182, col: 15, offset: 3776},
							expr: &charClassMatcher{
								pos:        position{line: 182, col: 15, offset: 3776},
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
			pos:  position{line: 190, col: 1, offset: 3986},
			expr: &actionExpr{
				pos: position{line: 191, col: 5, offset: 4013},
				run: (*parser).callondouble_quoted_string1,
				expr: &seqExpr{
					pos: position{line: 191, col: 7, offset: 4015},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 191, col: 7, offset: 4015},
							val:        "\"",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 191, col: 11, offset: 4019},
							expr: &choiceExpr{
								pos: position{line: 191, col: 13, offset: 4021},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 191, col: 13, offset: 4021},
										val:        "\\\"",
										ignoreCase: false,
									},
									&seqExpr{
										pos: position{line: 191, col: 20, offset: 4028},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 191, col: 20, offset: 4028},
												expr: &litMatcher{
													pos:        position{line: 191, col: 21, offset: 4029},
													val:        "\"",
													ignoreCase: false,
												},
											},
											&anyMatcher{
												line: 191, col: 25, offset: 4033,
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 192, col: 9, offset: 4048},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 192, col: 9, offset: 4048},
									val:        "\"",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 192, col: 15, offset: 4054},
									run: (*parser).callondouble_quoted_string13,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "single_quoted_string",
			pos:  position{line: 203, col: 1, offset: 4297},
			expr: &actionExpr{
				pos: position{line: 204, col: 5, offset: 4324},
				run: (*parser).callonsingle_quoted_string1,
				expr: &seqExpr{
					pos: position{line: 204, col: 7, offset: 4326},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 204, col: 7, offset: 4326},
							val:        "'",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 204, col: 11, offset: 4330},
							expr: &choiceExpr{
								pos: position{line: 204, col: 13, offset: 4332},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 204, col: 13, offset: 4332},
										val:        "\\'",
										ignoreCase: false,
									},
									&seqExpr{
										pos: position{line: 204, col: 20, offset: 4339},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 204, col: 20, offset: 4339},
												expr: &litMatcher{
													pos:        position{line: 204, col: 21, offset: 4340},
													val:        "'",
													ignoreCase: false,
												},
											},
											&anyMatcher{
												line: 204, col: 25, offset: 4344,
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 205, col: 9, offset: 4359},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 205, col: 9, offset: 4359},
									val:        "'",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 205, col: 15, offset: 4365},
									run: (*parser).callonsingle_quoted_string13,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "string_value",
			pos:  position{line: 216, col: 1, offset: 4573},
			expr: &actionExpr{
				pos: position{line: 217, col: 5, offset: 4592},
				run: (*parser).callonstring_value1,
				expr: &labeledExpr{
					pos:   position{line: 217, col: 5, offset: 4592},
					label: "str",
					expr: &choiceExpr{
						pos: position{line: 217, col: 11, offset: 4598},
						alternatives: []interface{}{
							&actionExpr{
								pos: position{line: 217, col: 11, offset: 4598},
								run: (*parser).callonstring_value4,
								expr: &labeledExpr{
									pos:   position{line: 217, col: 11, offset: 4598},
									label: "str",
									expr: &ruleRefExpr{
										pos:  position{line: 217, col: 15, offset: 4602},
										name: "double_quoted_string",
									},
								},
							},
							&actionExpr{
								pos: position{line: 219, col: 9, offset: 4712},
								run: (*parser).callonstring_value7,
								expr: &labeledExpr{
									pos:   position{line: 219, col: 9, offset: 4712},
									label: "str",
									expr: &ruleRefExpr{
										pos:  position{line: 219, col: 13, offset: 4716},
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
			pos:  position{line: 229, col: 1, offset: 4956},
			expr: &actionExpr{
				pos: position{line: 230, col: 5, offset: 4969},
				run: (*parser).callonregexp1,
				expr: &seqExpr{
					pos: position{line: 230, col: 7, offset: 4971},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 230, col: 7, offset: 4971},
							val:        "/",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 230, col: 11, offset: 4975},
							expr: &choiceExpr{
								pos: position{line: 230, col: 13, offset: 4977},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 230, col: 13, offset: 4977},
										val:        "\\/",
										ignoreCase: false,
									},
									&seqExpr{
										pos: position{line: 230, col: 20, offset: 4984},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 230, col: 20, offset: 4984},
												expr: &litMatcher{
													pos:        position{line: 230, col: 21, offset: 4985},
													val:        "/",
													ignoreCase: false,
												},
											},
											&anyMatcher{
												line: 230, col: 25, offset: 4989,
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 231, col: 9, offset: 5004},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 231, col: 9, offset: 5004},
									val:        "/",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 231, col: 15, offset: 5010},
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
			pos:  position{line: 243, col: 1, offset: 5236},
			expr: &actionExpr{
				pos: position{line: 244, col: 5, offset: 5249},
				run: (*parser).callonnumber1,
				expr: &seqExpr{
					pos: position{line: 244, col: 5, offset: 5249},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 244, col: 5, offset: 5249},
							expr: &litMatcher{
								pos:        position{line: 244, col: 5, offset: 5249},
								val:        "-",
								ignoreCase: false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 244, col: 10, offset: 5254},
							expr: &charClassMatcher{
								pos:        position{line: 244, col: 10, offset: 5254},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 244, col: 17, offset: 5261},
							expr: &seqExpr{
								pos: position{line: 244, col: 18, offset: 5262},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 244, col: 18, offset: 5262},
										val:        ".",
										ignoreCase: false,
									},
									&zeroOrMoreExpr{
										pos: position{line: 244, col: 22, offset: 5266},
										expr: &charClassMatcher{
											pos:        position{line: 244, col: 22, offset: 5266},
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
			pos:  position{line: 260, col: 1, offset: 5583},
			expr: &actionExpr{
				pos: position{line: 261, col: 5, offset: 5595},
				run: (*parser).callonarray1,
				expr: &seqExpr{
					pos: position{line: 261, col: 5, offset: 5595},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 261, col: 5, offset: 5595},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 261, col: 9, offset: 5599},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 261, col: 11, offset: 5601},
							label: "value",
							expr: &zeroOrOneExpr{
								pos: position{line: 261, col: 17, offset: 5607},
								expr: &actionExpr{
									pos: position{line: 262, col: 9, offset: 5618},
									run: (*parser).callonarray7,
									expr: &seqExpr{
										pos: position{line: 262, col: 9, offset: 5618},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 262, col: 9, offset: 5618},
												label: "value",
												expr: &ruleRefExpr{
													pos:  position{line: 262, col: 15, offset: 5624},
													name: "value",
												},
											},
											&labeledExpr{
												pos:   position{line: 262, col: 21, offset: 5630},
												label: "values",
												expr: &zeroOrMoreExpr{
													pos: position{line: 262, col: 28, offset: 5637},
													expr: &actionExpr{
														pos: position{line: 263, col: 13, offset: 5651},
														run: (*parser).callonarray13,
														expr: &seqExpr{
															pos: position{line: 263, col: 13, offset: 5651},
															exprs: []interface{}{
																&ruleRefExpr{
																	pos:  position{line: 263, col: 13, offset: 5651},
																	name: "_",
																},
																&litMatcher{
																	pos:        position{line: 263, col: 15, offset: 5653},
																	val:        ",",
																	ignoreCase: false,
																},
																&ruleRefExpr{
																	pos:  position{line: 263, col: 19, offset: 5657},
																	name: "_",
																},
																&labeledExpr{
																	pos:   position{line: 263, col: 21, offset: 5659},
																	label: "value",
																	expr: &ruleRefExpr{
																		pos:  position{line: 263, col: 27, offset: 5665},
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
							pos:  position{line: 269, col: 8, offset: 5796},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 270, col: 9, offset: 5808},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 270, col: 9, offset: 5808},
									val:        "]",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 270, col: 15, offset: 5814},
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
			pos:  position{line: 286, col: 1, offset: 6064},
			expr: &actionExpr{
				pos: position{line: 287, col: 5, offset: 6075},
				run: (*parser).callonhash1,
				expr: &seqExpr{
					pos: position{line: 287, col: 5, offset: 6075},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 287, col: 5, offset: 6075},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 287, col: 9, offset: 6079},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 287, col: 11, offset: 6081},
							label: "entries",
							expr: &zeroOrOneExpr{
								pos: position{line: 287, col: 19, offset: 6089},
								expr: &ruleRefExpr{
									pos:  position{line: 287, col: 19, offset: 6089},
									name: "hashentries",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 287, col: 32, offset: 6102},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 288, col: 9, offset: 6114},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 288, col: 9, offset: 6114},
									val:        "}",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 288, col: 15, offset: 6120},
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
			pos:  position{line: 300, col: 1, offset: 6359},
			expr: &actionExpr{
				pos: position{line: 301, col: 5, offset: 6377},
				run: (*parser).callonhashentries1,
				expr: &seqExpr{
					pos: position{line: 301, col: 5, offset: 6377},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 301, col: 5, offset: 6377},
							label: "hashentry",
							expr: &ruleRefExpr{
								pos:  position{line: 301, col: 15, offset: 6387},
								name: "hashentry",
							},
						},
						&labeledExpr{
							pos:   position{line: 301, col: 25, offset: 6397},
							label: "hashentries1",
							expr: &zeroOrMoreExpr{
								pos: position{line: 301, col: 38, offset: 6410},
								expr: &actionExpr{
									pos: position{line: 302, col: 9, offset: 6420},
									run: (*parser).callonhashentries7,
									expr: &seqExpr{
										pos: position{line: 302, col: 9, offset: 6420},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 302, col: 9, offset: 6420},
												name: "whitespace",
											},
											&labeledExpr{
												pos:   position{line: 302, col: 20, offset: 6431},
												label: "hashentry",
												expr: &ruleRefExpr{
													pos:  position{line: 302, col: 30, offset: 6441},
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
			pos:  position{line: 314, col: 1, offset: 6693},
			expr: &actionExpr{
				pos: position{line: 315, col: 5, offset: 6709},
				run: (*parser).callonhashentry1,
				expr: &seqExpr{
					pos: position{line: 315, col: 5, offset: 6709},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 315, col: 5, offset: 6709},
							label: "name",
							expr: &choiceExpr{
								pos: position{line: 315, col: 11, offset: 6715},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 315, col: 11, offset: 6715},
										name: "number",
									},
									&ruleRefExpr{
										pos:  position{line: 315, col: 20, offset: 6724},
										name: "bareword",
									},
									&ruleRefExpr{
										pos:  position{line: 315, col: 31, offset: 6735},
										name: "string_value",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 315, col: 45, offset: 6749},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 315, col: 47, offset: 6751},
							val:        "=>",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 315, col: 52, offset: 6756},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 315, col: 54, offset: 6758},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 315, col: 60, offset: 6764},
								name: "value",
							},
						},
					},
				},
			},
		},
		{
			name: "branch",
			pos:  position{line: 326, col: 1, offset: 6931},
			expr: &actionExpr{
				pos: position{line: 327, col: 5, offset: 6944},
				run: (*parser).callonbranch1,
				expr: &seqExpr{
					pos: position{line: 327, col: 5, offset: 6944},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 327, col: 5, offset: 6944},
							label: "ifBlock",
							expr: &ruleRefExpr{
								pos:  position{line: 327, col: 13, offset: 6952},
								name: "if_cond",
							},
						},
						&labeledExpr{
							pos:   position{line: 327, col: 21, offset: 6960},
							label: "elseIfBlocks",
							expr: &zeroOrMoreExpr{
								pos: position{line: 327, col: 34, offset: 6973},
								expr: &actionExpr{
									pos: position{line: 328, col: 9, offset: 6983},
									run: (*parser).callonbranch7,
									expr: &seqExpr{
										pos: position{line: 328, col: 9, offset: 6983},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 328, col: 9, offset: 6983},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 328, col: 11, offset: 6985},
												label: "eib",
												expr: &ruleRefExpr{
													pos:  position{line: 328, col: 15, offset: 6989},
													name: "else_if",
												},
											},
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 331, col: 12, offset: 7056},
							label: "elseBlock",
							expr: &zeroOrOneExpr{
								pos: position{line: 331, col: 22, offset: 7066},
								expr: &actionExpr{
									pos: position{line: 332, col: 13, offset: 7080},
									run: (*parser).callonbranch14,
									expr: &seqExpr{
										pos: position{line: 332, col: 13, offset: 7080},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 332, col: 13, offset: 7080},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 332, col: 15, offset: 7082},
												label: "eb",
												expr: &ruleRefExpr{
													pos:  position{line: 332, col: 18, offset: 7085},
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
			pos:  position{line: 344, col: 1, offset: 7334},
			expr: &actionExpr{
				pos: position{line: 345, col: 5, offset: 7348},
				run: (*parser).callonif_cond1,
				expr: &seqExpr{
					pos: position{line: 345, col: 5, offset: 7348},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 345, col: 5, offset: 7348},
							val:        "if",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 345, col: 10, offset: 7353},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 345, col: 12, offset: 7355},
							label: "cond",
							expr: &ruleRefExpr{
								pos:  position{line: 345, col: 17, offset: 7360},
								name: "condition",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 345, col: 27, offset: 7370},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 345, col: 29, offset: 7372},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 345, col: 33, offset: 7376},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 345, col: 35, offset: 7378},
							label: "bops",
							expr: &zeroOrMoreExpr{
								pos: position{line: 345, col: 40, offset: 7383},
								expr: &actionExpr{
									pos: position{line: 346, col: 13, offset: 7397},
									run: (*parser).callonif_cond12,
									expr: &seqExpr{
										pos: position{line: 346, col: 13, offset: 7397},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 346, col: 13, offset: 7397},
												label: "bop",
												expr: &ruleRefExpr{
													pos:  position{line: 346, col: 17, offset: 7401},
													name: "branch_or_plugin",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 346, col: 34, offset: 7418},
												name: "_",
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 350, col: 13, offset: 7493},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 350, col: 13, offset: 7493},
									val:        "}",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 350, col: 19, offset: 7499},
									run: (*parser).callonif_cond19,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "else_if",
			pos:  position{line: 362, col: 1, offset: 7773},
			expr: &actionExpr{
				pos: position{line: 363, col: 5, offset: 7787},
				run: (*parser).callonelse_if1,
				expr: &seqExpr{
					pos: position{line: 363, col: 5, offset: 7787},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 363, col: 5, offset: 7787},
							val:        "else",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 363, col: 12, offset: 7794},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 363, col: 14, offset: 7796},
							val:        "if",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 363, col: 19, offset: 7801},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 363, col: 21, offset: 7803},
							label: "cond",
							expr: &ruleRefExpr{
								pos:  position{line: 363, col: 26, offset: 7808},
								name: "condition",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 363, col: 36, offset: 7818},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 363, col: 38, offset: 7820},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 363, col: 42, offset: 7824},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 363, col: 44, offset: 7826},
							label: "bops",
							expr: &zeroOrMoreExpr{
								pos: position{line: 363, col: 49, offset: 7831},
								expr: &actionExpr{
									pos: position{line: 364, col: 9, offset: 7841},
									run: (*parser).callonelse_if14,
									expr: &seqExpr{
										pos: position{line: 364, col: 9, offset: 7841},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 364, col: 9, offset: 7841},
												label: "bop",
												expr: &ruleRefExpr{
													pos:  position{line: 364, col: 13, offset: 7845},
													name: "branch_or_plugin",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 364, col: 30, offset: 7862},
												name: "_",
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 368, col: 9, offset: 7921},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 368, col: 9, offset: 7921},
									val:        "}",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 368, col: 15, offset: 7927},
									run: (*parser).callonelse_if21,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "else_cond",
			pos:  position{line: 380, col: 1, offset: 8169},
			expr: &actionExpr{
				pos: position{line: 381, col: 5, offset: 8185},
				run: (*parser).callonelse_cond1,
				expr: &seqExpr{
					pos: position{line: 381, col: 5, offset: 8185},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 381, col: 5, offset: 8185},
							val:        "else",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 381, col: 12, offset: 8192},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 381, col: 14, offset: 8194},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 381, col: 18, offset: 8198},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 381, col: 20, offset: 8200},
							label: "bops",
							expr: &zeroOrMoreExpr{
								pos: position{line: 381, col: 25, offset: 8205},
								expr: &actionExpr{
									pos: position{line: 382, col: 9, offset: 8215},
									run: (*parser).callonelse_cond9,
									expr: &seqExpr{
										pos: position{line: 382, col: 9, offset: 8215},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 382, col: 9, offset: 8215},
												label: "bop",
												expr: &ruleRefExpr{
													pos:  position{line: 382, col: 13, offset: 8219},
													name: "branch_or_plugin",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 382, col: 30, offset: 8236},
												name: "_",
											},
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 386, col: 9, offset: 8295},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 386, col: 9, offset: 8295},
									val:        "}",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 386, col: 15, offset: 8301},
									run: (*parser).callonelse_cond16,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "condition",
			pos:  position{line: 398, col: 1, offset: 8550},
			expr: &actionExpr{
				pos: position{line: 399, col: 5, offset: 8566},
				run: (*parser).calloncondition1,
				expr: &seqExpr{
					pos: position{line: 399, col: 5, offset: 8566},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 399, col: 5, offset: 8566},
							label: "cond",
							expr: &ruleRefExpr{
								pos:  position{line: 399, col: 10, offset: 8571},
								name: "expression",
							},
						},
						&labeledExpr{
							pos:   position{line: 399, col: 21, offset: 8582},
							label: "conds",
							expr: &zeroOrMoreExpr{
								pos: position{line: 399, col: 27, offset: 8588},
								expr: &actionExpr{
									pos: position{line: 400, col: 9, offset: 8598},
									run: (*parser).calloncondition7,
									expr: &seqExpr{
										pos: position{line: 400, col: 9, offset: 8598},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 400, col: 9, offset: 8598},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 400, col: 11, offset: 8600},
												label: "bo",
												expr: &ruleRefExpr{
													pos:  position{line: 400, col: 14, offset: 8603},
													name: "boolean_operator",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 400, col: 31, offset: 8620},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 400, col: 33, offset: 8622},
												label: "cond",
												expr: &ruleRefExpr{
													pos:  position{line: 400, col: 38, offset: 8627},
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
			pos:  position{line: 419, col: 1, offset: 9026},
			expr: &choiceExpr{
				pos: position{line: 421, col: 9, offset: 9054},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 422, col: 13, offset: 9068},
						run: (*parser).callonexpression2,
						expr: &seqExpr{
							pos: position{line: 422, col: 13, offset: 9068},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 422, col: 13, offset: 9068},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 422, col: 17, offset: 9072},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 422, col: 19, offset: 9074},
									label: "cond",
									expr: &ruleRefExpr{
										pos:  position{line: 422, col: 24, offset: 9079},
										name: "condition",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 422, col: 34, offset: 9089},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 422, col: 36, offset: 9091},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 426, col: 9, offset: 9179},
						name: "negative_expression",
					},
					&ruleRefExpr{
						pos:  position{line: 427, col: 9, offset: 9207},
						name: "in_expression",
					},
					&ruleRefExpr{
						pos:  position{line: 428, col: 9, offset: 9229},
						name: "not_in_expression",
					},
					&ruleRefExpr{
						pos:  position{line: 429, col: 9, offset: 9255},
						name: "compare_expression",
					},
					&ruleRefExpr{
						pos:  position{line: 430, col: 9, offset: 9282},
						name: "regexp_expression",
					},
					&actionExpr{
						pos: position{line: 431, col: 9, offset: 9308},
						run: (*parser).callonexpression15,
						expr: &labeledExpr{
							pos:   position{line: 431, col: 9, offset: 9308},
							label: "rv",
							expr: &ruleRefExpr{
								pos:  position{line: 431, col: 12, offset: 9311},
								name: "rvalue",
							},
						},
					},
				},
			},
		},
		{
			name: "negative_expression",
			pos:  position{line: 444, col: 1, offset: 9595},
			expr: &choiceExpr{
				pos: position{line: 446, col: 9, offset: 9632},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 447, col: 13, offset: 9646},
						run: (*parser).callonnegative_expression2,
						expr: &seqExpr{
							pos: position{line: 447, col: 13, offset: 9646},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 447, col: 13, offset: 9646},
									val:        "!",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 447, col: 17, offset: 9650},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 447, col: 19, offset: 9652},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 447, col: 23, offset: 9656},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 447, col: 25, offset: 9658},
									label: "cond",
									expr: &ruleRefExpr{
										pos:  position{line: 447, col: 30, offset: 9663},
										name: "condition",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 447, col: 40, offset: 9673},
									name: "_",
								},
								&choiceExpr{
									pos: position{line: 448, col: 17, offset: 9693},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 448, col: 17, offset: 9693},
											val:        ")",
											ignoreCase: false,
										},
										&andCodeExpr{
											pos: position{line: 448, col: 23, offset: 9699},
											run: (*parser).callonnegative_expression13,
										},
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 456, col: 11, offset: 9900},
						run: (*parser).callonnegative_expression14,
						expr: &seqExpr{
							pos: position{line: 456, col: 11, offset: 9900},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 456, col: 11, offset: 9900},
									val:        "!",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 456, col: 15, offset: 9904},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 456, col: 17, offset: 9906},
									label: "sel",
									expr: &ruleRefExpr{
										pos:  position{line: 456, col: 21, offset: 9910},
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
			pos:  position{line: 467, col: 1, offset: 10110},
			expr: &actionExpr{
				pos: position{line: 468, col: 5, offset: 10130},
				run: (*parser).callonin_expression1,
				expr: &seqExpr{
					pos: position{line: 468, col: 5, offset: 10130},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 468, col: 5, offset: 10130},
							label: "lv",
							expr: &ruleRefExpr{
								pos:  position{line: 468, col: 8, offset: 10133},
								name: "rvalue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 468, col: 15, offset: 10140},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 468, col: 17, offset: 10142},
							name: "in_operator",
						},
						&ruleRefExpr{
							pos:  position{line: 468, col: 29, offset: 10154},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 468, col: 31, offset: 10156},
							label: "rv",
							expr: &ruleRefExpr{
								pos:  position{line: 468, col: 34, offset: 10159},
								name: "rvalue",
							},
						},
					},
				},
			},
		},
		{
			name: "not_in_expression",
			pos:  position{line: 477, col: 1, offset: 10339},
			expr: &actionExpr{
				pos: position{line: 478, col: 5, offset: 10363},
				run: (*parser).callonnot_in_expression1,
				expr: &seqExpr{
					pos: position{line: 478, col: 5, offset: 10363},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 478, col: 5, offset: 10363},
							label: "lv",
							expr: &ruleRefExpr{
								pos:  position{line: 478, col: 8, offset: 10366},
								name: "rvalue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 478, col: 15, offset: 10373},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 478, col: 17, offset: 10375},
							name: "not_in_operator",
						},
						&ruleRefExpr{
							pos:  position{line: 478, col: 33, offset: 10391},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 478, col: 35, offset: 10393},
							label: "rv",
							expr: &ruleRefExpr{
								pos:  position{line: 478, col: 38, offset: 10396},
								name: "rvalue",
							},
						},
					},
				},
			},
		},
		{
			name: "in_operator",
			pos:  position{line: 486, col: 1, offset: 10497},
			expr: &choiceExpr{
				pos: position{line: 487, col: 5, offset: 10515},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 487, col: 5, offset: 10515},
						val:        "in",
						ignoreCase: false,
					},
					&andCodeExpr{
						pos: position{line: 487, col: 12, offset: 10522},
						run: (*parser).callonin_operator3,
					},
				},
			},
		},
		{
			name: "not_in_operator",
			pos:  position{line: 495, col: 1, offset: 10644},
			expr: &choiceExpr{
				pos: position{line: 496, col: 5, offset: 10666},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 496, col: 5, offset: 10666},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 496, col: 5, offset: 10666},
								val:        "not ",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 496, col: 12, offset: 10673},
								name: "_",
							},
							&litMatcher{
								pos:        position{line: 496, col: 14, offset: 10675},
								val:        "in",
								ignoreCase: false,
							},
						},
					},
					&andCodeExpr{
						pos: position{line: 496, col: 21, offset: 10682},
						run: (*parser).callonnot_in_operator6,
					},
				},
			},
		},
		{
			name: "rvalue",
			pos:  position{line: 507, col: 1, offset: 11008},
			expr: &choiceExpr{
				pos: position{line: 508, col: 5, offset: 11021},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 508, col: 5, offset: 11021},
						name: "string_value",
					},
					&ruleRefExpr{
						pos:  position{line: 508, col: 20, offset: 11036},
						name: "number",
					},
					&ruleRefExpr{
						pos:  position{line: 508, col: 29, offset: 11045},
						name: "selector",
					},
					&ruleRefExpr{
						pos:  position{line: 508, col: 40, offset: 11056},
						name: "array",
					},
					&ruleRefExpr{
						pos:  position{line: 508, col: 48, offset: 11064},
						name: "regexp",
					},
					&andCodeExpr{
						pos: position{line: 508, col: 57, offset: 11073},
						run: (*parser).callonrvalue7,
					},
				},
			},
		},
		{
			name: "compare_expression",
			pos:  position{line: 542, col: 1, offset: 11811},
			expr: &actionExpr{
				pos: position{line: 543, col: 5, offset: 11836},
				run: (*parser).calloncompare_expression1,
				expr: &seqExpr{
					pos: position{line: 543, col: 5, offset: 11836},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 543, col: 5, offset: 11836},
							label: "lv",
							expr: &ruleRefExpr{
								pos:  position{line: 543, col: 8, offset: 11839},
								name: "rvalue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 543, col: 15, offset: 11846},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 543, col: 17, offset: 11848},
							label: "co",
							expr: &ruleRefExpr{
								pos:  position{line: 543, col: 20, offset: 11851},
								name: "compare_operator",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 543, col: 37, offset: 11868},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 543, col: 39, offset: 11870},
							label: "rv",
							expr: &ruleRefExpr{
								pos:  position{line: 543, col: 42, offset: 11873},
								name: "rvalue",
							},
						},
					},
				},
			},
		},
		{
			name: "compare_operator",
			pos:  position{line: 552, col: 1, offset: 12072},
			expr: &actionExpr{
				pos: position{line: 553, col: 5, offset: 12095},
				run: (*parser).calloncompare_operator1,
				expr: &choiceExpr{
					pos: position{line: 553, col: 6, offset: 12096},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 553, col: 6, offset: 12096},
							val:        "==",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 553, col: 13, offset: 12103},
							val:        "!=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 553, col: 20, offset: 12110},
							val:        "<=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 553, col: 27, offset: 12117},
							val:        ">=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 553, col: 34, offset: 12124},
							val:        "<",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 553, col: 40, offset: 12130},
							val:        ">",
							ignoreCase: false,
						},
						&andCodeExpr{
							pos: position{line: 553, col: 46, offset: 12136},
							run: (*parser).calloncompare_operator9,
						},
					},
				},
			},
		},
		{
			name: "regexp_expression",
			pos:  position{line: 564, col: 1, offset: 12421},
			expr: &actionExpr{
				pos: position{line: 565, col: 5, offset: 12445},
				run: (*parser).callonregexp_expression1,
				expr: &seqExpr{
					pos: position{line: 565, col: 5, offset: 12445},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 565, col: 5, offset: 12445},
							label: "lv",
							expr: &ruleRefExpr{
								pos:  position{line: 565, col: 8, offset: 12448},
								name: "rvalue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 565, col: 15, offset: 12455},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 565, col: 18, offset: 12458},
							label: "ro",
							expr: &ruleRefExpr{
								pos:  position{line: 565, col: 21, offset: 12461},
								name: "regexp_operator",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 565, col: 37, offset: 12477},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 565, col: 39, offset: 12479},
							label: "rv",
							expr: &choiceExpr{
								pos: position{line: 565, col: 43, offset: 12483},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 565, col: 43, offset: 12483},
										name: "string_value",
									},
									&ruleRefExpr{
										pos:  position{line: 565, col: 58, offset: 12498},
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
			pos:  position{line: 573, col: 1, offset: 12657},
			expr: &actionExpr{
				pos: position{line: 574, col: 5, offset: 12679},
				run: (*parser).callonregexp_operator1,
				expr: &choiceExpr{
					pos: position{line: 574, col: 6, offset: 12680},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 574, col: 6, offset: 12680},
							val:        "=~",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 574, col: 13, offset: 12687},
							val:        "!~",
							ignoreCase: false,
						},
						&andCodeExpr{
							pos: position{line: 574, col: 20, offset: 12694},
							run: (*parser).callonregexp_operator5,
						},
					},
				},
			},
		},
		{
			name: "boolean_operator",
			pos:  position{line: 585, col: 1, offset: 12958},
			expr: &actionExpr{
				pos: position{line: 586, col: 5, offset: 12981},
				run: (*parser).callonboolean_operator1,
				expr: &choiceExpr{
					pos: position{line: 586, col: 6, offset: 12982},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 586, col: 6, offset: 12982},
							val:        "and",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 586, col: 14, offset: 12990},
							val:        "or",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 586, col: 21, offset: 12997},
							val:        "xor",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 586, col: 29, offset: 13005},
							val:        "nand",
							ignoreCase: false,
						},
						&andCodeExpr{
							pos: position{line: 586, col: 38, offset: 13014},
							run: (*parser).callonboolean_operator7,
						},
					},
				},
			},
		},
		{
			name: "selector",
			pos:  position{line: 597, col: 1, offset: 13252},
			expr: &actionExpr{
				pos: position{line: 598, col: 5, offset: 13267},
				run: (*parser).callonselector1,
				expr: &labeledExpr{
					pos:   position{line: 598, col: 5, offset: 13267},
					label: "ses",
					expr: &oneOrMoreExpr{
						pos: position{line: 598, col: 9, offset: 13271},
						expr: &ruleRefExpr{
							pos:  position{line: 598, col: 9, offset: 13271},
							name: "selector_element",
						},
					},
				},
			},
		},
		{
			name: "selector_element",
			pos:  position{line: 607, col: 1, offset: 13435},
			expr: &actionExpr{
				pos: position{line: 608, col: 5, offset: 13458},
				run: (*parser).callonselector_element1,
				expr: &seqExpr{
					pos: position{line: 608, col: 5, offset: 13458},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 608, col: 5, offset: 13458},
							val:        "[",
							ignoreCase: false,
						},
						&oneOrMoreExpr{
							pos: position{line: 608, col: 9, offset: 13462},
							expr: &charClassMatcher{
								pos:        position{line: 608, col: 9, offset: 13462},
								val:        "[^\\],]",
								chars:      []rune{']', ','},
								ignoreCase: false,
								inverted:   true,
							},
						},
						&choiceExpr{
							pos: position{line: 609, col: 9, offset: 13480},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 609, col: 9, offset: 13480},
									val:        "]",
									ignoreCase: false,
								},
								&andCodeExpr{
									pos: position{line: 609, col: 15, offset: 13486},
									run: (*parser).callonselector_element8,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 616, col: 1, offset: 13628},
			expr: &notExpr{
				pos: position{line: 616, col: 7, offset: 13634},
				expr: &anyMatcher{
					line: 616, col: 8, offset: 13635,
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

func (c *current) onplugin_section10(bop interface{}) (interface{}, error) {

	return ret(bop)

}

func (p *parser) callonplugin_section10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onplugin_section10(stack["bop"])
}

func (c *current) onplugin_section17() (bool, error) {
	return pushError("expect closing curly bracket", c)

}

func (p *parser) callonplugin_section17() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onplugin_section17()
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

func (c *current) onplugin_type8() (bool, error) {
	return pushError("expect plugin type (input, filter, output)", c)

}

func (p *parser) callonplugin_type8() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onplugin_type8()
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

func (c *current) onarray_value7() (bool, error) {
	return fatalError("invalid array value", c)

}

func (p *parser) callonarray_value7() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onarray_value7()
}

func (c *current) onbareword1() (interface{}, error) {
	return ast.NewStringAttribute("", string(c.text), ast.Bareword), nil

}

func (p *parser) callonbareword1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onbareword1()
}

func (c *current) ondouble_quoted_string13() (bool, error) {
	return fatalError("expect closing double quotes (\")", c)

}

func (p *parser) callondouble_quoted_string13() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.ondouble_quoted_string13()
}

func (c *current) ondouble_quoted_string1() (interface{}, error) {
	return enclosedValue(c)

}

func (p *parser) callondouble_quoted_string1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.ondouble_quoted_string1()
}

func (c *current) onsingle_quoted_string13() (bool, error) {
	return fatalError("expect closing single quote (')", c)

}

func (p *parser) callonsingle_quoted_string13() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onsingle_quoted_string13()
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

func (c *current) onif_cond12(bop interface{}) (interface{}, error) {
	return ret(bop)

}

func (p *parser) callonif_cond12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onif_cond12(stack["bop"])
}

func (c *current) onif_cond19() (bool, error) {
	return fatalError("expect closing curly bracket", c)

}

func (p *parser) callonif_cond19() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onif_cond19()
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

func (c *current) onelse_if21() (bool, error) {
	return fatalError("expect closing curly bracket", c)

}

func (p *parser) callonelse_if21() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onelse_if21()
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

func (c *current) onelse_cond16() (bool, error) {
	return fatalError("expect closing curly bracket", c)

}

func (p *parser) callonelse_cond16() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onelse_cond16()
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

func (c *current) onnegative_expression13() (bool, error) {
	return fatalError("expect closing parenthesis", c)

}

func (p *parser) callonnegative_expression13() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onnegative_expression13()
}

func (c *current) onnegative_expression2(cond interface{}) (interface{}, error) {
	return negative_expression(cond)

}

func (p *parser) callonnegative_expression2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onnegative_expression2(stack["cond"])
}

func (c *current) onnegative_expression14(sel interface{}) (interface{}, error) {
	return negative_selector(sel)

}

func (p *parser) callonnegative_expression14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onnegative_expression14(stack["sel"])
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

func (c *current) onin_operator3() (bool, error) {
	return pushError("expect in operator (in)", c)

}

func (p *parser) callonin_operator3() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onin_operator3()
}

func (c *current) onnot_in_operator6() (bool, error) {
	return pushError("expect not in operator (not in)", c)

}

func (p *parser) callonnot_in_operator6() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onnot_in_operator6()
}

func (c *current) onrvalue7() (bool, error) {
	return pushError("invalid value for expression", c)

}

func (p *parser) callonrvalue7() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onrvalue7()
}

func (c *current) oncompare_expression1(lv, co, rv interface{}) (interface{}, error) {
	return compare_expression(lv, co, rv)

}

func (p *parser) calloncompare_expression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.oncompare_expression1(stack["lv"], stack["co"], stack["rv"])
}

func (c *current) oncompare_operator9() (bool, error) {
	return pushError("expect compare operator (==, !=, <=, >=, <, >)", c)

}

func (p *parser) calloncompare_operator9() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.oncompare_operator9()
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

func (c *current) onregexp_operator5() (bool, error) {
	return pushError("expect regexp comparison operator (=~, !~)", c)

}

func (p *parser) callonregexp_operator5() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onregexp_operator5()
}

func (c *current) onregexp_operator1() (interface{}, error) {
	return regexp_operator(string(c.text))

}

func (p *parser) callonregexp_operator1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onregexp_operator1()
}

func (c *current) onboolean_operator7() (bool, error) {
	return pushError("expect boolean operator (and, or, xor, nand)", c)

}

func (p *parser) callonboolean_operator7() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onboolean_operator7()
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

func (c *current) onselector_element8() (bool, error) {
	return fatalError("expect closing square bracket", c)

}

func (p *parser) callonselector_element8() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onselector_element8()
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
