package config

//go:generate pigeon -nolint -optimize-grammar -o logstash_config.go logstash_config.peg

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/breml/logstash-config/ast"
)

type exceptionalCommentsWarnings []string

func (w exceptionalCommentsWarnings) Clone() interface{} {
	clone := make(exceptionalCommentsWarnings, 0, len(w))
	clone = append(clone, w...)
	return clone
}

func initState(c *current) error {
	if _, ok := c.state[exceptionalCommentWarning]; !ok {
		warnings := make(exceptionalCommentsWarnings, 0)
		c.state[exceptionalCommentWarning] = warnings
	}
	return nil
}

func initParser() (bool, error) {
	farthestFailure = []errPos{}
	return true, nil
}

func ret(el interface{}) (interface{}, error) {
	return el, nil
}

func retConfig(c *current, el interface{}) (interface{}, error) {
	conf, _ := el.(ast.Config)

	warnings, _ := c.state[exceptionalCommentWarning].(exceptionalCommentsWarnings)

	conf.Warnings = warnings
	return conf, nil
}

func commentBlock(comment1 interface{}, spaceBefore bool, spaceAfter bool) ast.CommentBlock {
	comments1 := toIfaceSlice(comment1)
	var comments []ast.Comment
	for _, icb1 := range comments1 {
		switch t := icb1.(type) {
		case []ast.Comment:
			// If spaceBefore is enabled and it is the first comment and there are
			// actualy comment in t, enable space before for the first comment.
			if spaceBefore && len(comments) == 0 && len(t) > 0 {
				t[0].SpaceBefore = true
			}
			comments = append(comments, t...)
		case ast.Whitespace:
			// If spaceAfter is preserved.
			if len(comments) > 0 && spaceAfter {
				comments[len(comments)-1].SpaceAfter = true
				continue
			}
		}
	}
	return comments
}

func configSection(ps1, psComment1 interface{}) (ast.PluginSection, error) {
	ps, ok := ps1.(ast.PluginSection)
	if !ok {
		return ast.PluginSection{}, fmt.Errorf("Value is not a PluginSection: %#v", ps1)
	}

	psComments := commentBlock(psComment1, false, false)

	ps.CommentBlock = psComments

	return ps, nil
}

func config(ps1, pss1, psComment1, footerComment1 interface{}) (ast.Config, error) {
	var (
		input  []ast.PluginSection
		filter []ast.PluginSection
		output []ast.PluginSection
	)

	ips := toIfaceSlice(ps1)
	ips = append(ips, toIfaceSlice(pss1)...)

	psComment := commentBlock(psComment1, false, true)
	footerComments := commentBlock(footerComment1, true, false)

	first := true
	for _, ips1 := range ips {
		if ps, ok := ips1.(ast.PluginSection); ok {
			if first {
				ps.CommentBlock = psComment
				first = false
			}
			switch ps.PluginType {
			case ast.Input:
				input = append(input, ps)
			case ast.Filter:
				filter = append(filter, ps)
			case ast.Output:
				output = append(output, ps)
			default:
				return ast.Config{}, fmt.Errorf("PluginType is not supported: %#v", ps)
			}
		} else {
			return ast.Config{}, fmt.Errorf("Value is not a PluginSection: %#v", ips1)
		}
	}

	return ast.Config{
		Input:         input,
		Filter:        filter,
		Output:        output,
		FooterComment: footerComments,
	}, nil
}

func warnComment(c *current) error {
	ok, _ := c.globalStore[exceptionalCommentWarning].(bool)
	if ok && bytes.Contains(c.text, []byte("#")) {
		warnings, _ := c.state[exceptionalCommentWarning].(exceptionalCommentsWarnings)
		warnings = append(warnings, fmt.Sprintf("exceptional comment at %s", c.pos))
		c.state[exceptionalCommentWarning] = warnings
	}
	return nil
}

func whitespace() (ast.Whitespace, error) {
	return ast.Whitespace{}, nil
}

func comment(c *current) ([]ast.Comment, error) {
	if ignoreComment, ok := c.globalStore[ignoreComment].(bool); ok && ignoreComment {
		return nil, nil
	}

	lines := strings.Split(strings.Trim(strings.ReplaceAll(string(c.text), "\r", ""), "\n"), "\n")

	var commentLines []ast.Comment
	var space bool
	for i, line := range lines {
		idx := strings.Index(line, "#")
		if idx == -1 {
			if i > 0 {
				space = true
			}
			continue
		}
		var c string
		if len(line) > idx+1 {
			c = string(line[idx+1:])
		}
		c = strings.Trim(c, " ")
		commentLines = append(commentLines, ast.NewComment(c, space))
		space = false
	}

	return commentLines, nil
}

func pluginSection(pt1, bops1, footerComment1 interface{}) (ast.PluginSection, error) {
	pt := ast.PluginType(pt1.(int))
	ibops := toIfaceSlice(bops1)
	var bops []ast.BranchOrPlugin

	for _, bop := range ibops {
		bop := bop.(ast.BranchOrPlugin)
		bops = append(bops, bop)
	}

	return ast.PluginSection{
		PluginType:      pt,
		BranchOrPlugins: bops,
		FooterComment:   commentBlock(footerComment1, false, false),
	}, nil
}

func plugin(name, attributes1, comment1, footerComment1 interface{}) (ast.Plugin, error) {
	var attributes []ast.Attribute
	if attributes1 != nil {
		attributes = attributes1.([]ast.Attribute)
	}

	p := ast.NewPlugin(name.(string), attributes...)

	p.Comment = commentBlock(comment1, false, false)
	p.FooterComment = commentBlock(footerComment1, false, false)

	return p, nil
}

func attributes(attribute1, attributes1, comment1 interface{}) ([]ast.Attribute, error) {
	attribute, _ := attributeComment(attribute1, comment1, false)
	iattributes := toIfaceSlice(attribute)
	iattributes = append(iattributes, toIfaceSlice(attributes1)...)

	var attributes []ast.Attribute

	for _, attr := range iattributes {
		if attr, ok := attr.(ast.Attribute); ok {
			attributes = append(attributes, attr)
		} else {
			return nil, fmt.Errorf("Argument is not an attribute: %#v", attr)
		}
	}

	return attributes, nil
}

func attributeComment(attribute1, comment1 interface{}, spaceBefore bool) (ast.Attribute, error) {
	var attribute ast.Attribute
	switch attr := attribute1.(type) {
	case ast.StringAttribute:
		attr.Comment = commentBlock(comment1, spaceBefore, false)
		attribute = attr
	case ast.NumberAttribute:
		attr.Comment = commentBlock(comment1, spaceBefore, false)
		attribute = attr
	case ast.ArrayAttribute:
		attr.Comment = commentBlock(comment1, spaceBefore, false)
		attribute = attr
	case ast.HashAttribute:
		attr.Comment = commentBlock(comment1, spaceBefore, false)
		attribute = attr
	case ast.PluginAttribute:
		attr.Comment = commentBlock(comment1, spaceBefore, false)
		attribute = attr
	default:
		return nil, fmt.Errorf("Unsupported attribute type %#v", attribute1)
	}
	return attribute, nil
}

func attribute(name, value interface{}) (ast.Attribute, error) {
	var key ast.StringAttribute

	switch name := name.(type) {
	case ast.StringAttribute:
		key = name
	case string:
		key = ast.NewStringAttribute("", name, ast.Bareword)
	default:
		return nil, fmt.Errorf("Type for attribute name is not supported: %#v", name)
	}

	switch value := value.(type) {
	case ast.StringAttribute:
		return ast.NewStringAttribute(key.ValueString(), value.Value(), value.StringAttributeType()), nil
	case ast.NumberAttribute:
		return ast.NewNumberAttribute(key.ValueString(), value.Value()), nil
	case ast.ArrayAttribute:
		aa := ast.NewArrayAttribute(key.ValueString(), value.Value()...)
		aa.FooterComment = value.FooterComment
		return aa, nil
	case ast.HashAttribute:
		ha := ast.NewHashAttribute(key.ValueString(), value.Value()...)
		ha.FooterComment = value.FooterComment
		return ha, nil
	case ast.Plugin:
		return ast.NewPluginAttribute(key.ValueString(), value), nil
	default:
		return nil, fmt.Errorf("Type of value %#v with name %s is not supported", value, key.ValueString())
	}
}

func regexp(c *current) (ast.Regexp, error) {
	val, _ := enclosedValue(c)
	return ast.NewRegexp(val), nil
}

func number(value string) (ast.NumberAttribute, error) {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		// TODO: is this possible to happen? are all values, which are valid floats in Logstash/Ruby also valid floats in Go?
		return ast.NumberAttribute{}, err
	}
	return ast.NewNumberAttribute("", f), nil
}

func array(attributes1, footerComment1 interface{}) (ast.ArrayAttribute, error) {
	var attributes []ast.Attribute
	if attributes1 != nil {
		attributes = attributes1.([]ast.Attribute)
	}

	a := ast.NewArrayAttribute("", attributes...)

	a.FooterComment = commentBlock(footerComment1, false, false)

	return a, nil
}

func hash(attributes1, footerComment1 interface{}) (ast.HashAttribute, error) {
	var hashentries []ast.HashEntry
	if attributes1 != nil {
		hashentries = attributes1.([]ast.HashEntry)
	}

	a := ast.NewHashAttribute("", hashentries...)

	a.FooterComment = commentBlock(footerComment1, false, false)

	return a, nil
}

func hashentries(attribute, attributes1 interface{}) ([]ast.HashEntry, error) {
	entry := attribute.(ast.HashEntry)
	if len(entry.Comment) > 0 {
		entry.Comment[0].SpaceBefore = false
	}

	iattributes := toIfaceSlice(entry)
	iattributes = append(iattributes, toIfaceSlice(attributes1)...)

	var attributes []ast.HashEntry

	for _, attr := range iattributes {
		if attr, ok := attr.(ast.HashEntry); ok {
			attributes = append(attributes, attr)
		} else {
			return nil, fmt.Errorf("Argument is not an attribute")
		}
	}

	return attributes, nil
}

func hashentry(name, value, comment interface{}) (ast.HashEntry, error) {
	key := name.(ast.HashEntryKey)

	he := ast.NewHashEntry(key, value.(ast.Attribute))
	he.Comment = commentBlock(comment, true, false)

	return he, nil
}

func elseIfComment(eib1, eibComment1 interface{}) (ast.ElseIfBlock, error) {
	eib := eib1.(ast.ElseIfBlock)
	eib.Comment = commentBlock(eibComment1, false, false)
	return eib, nil
}

func elseComment(eb1, ebComment1 interface{}) (ast.ElseBlock, error) {
	eb := eb1.(ast.ElseBlock)
	eb.Comment = commentBlock(ebComment1, false, false)
	return eb, nil
}

func branch(ifBlock1, elseIfBlocks1, elseBlock1, ifComment interface{}) (ast.Branch, error) {
	ielseIfBlocks := toIfaceSlice(elseIfBlocks1)

	var elseIfBlocks []ast.ElseIfBlock
	for _, elseIfBlock := range ielseIfBlocks {
		if elseIfBlock, ok := elseIfBlock.(ast.ElseIfBlock); ok {
			elseIfBlocks = append(elseIfBlocks, elseIfBlock)
		} else {
			return ast.Branch{}, fmt.Errorf("Argument is not an elseIfBlock: %#v", elseIfBlock)
		}
	}

	var elseBlock ast.ElseBlock
	if elseBlock1 != nil {
		elseBlock = elseBlock1.(ast.ElseBlock)
	}

	ifBlock := ifBlock1.(ast.IfBlock)
	ifBlock.Comment = commentBlock(ifComment, false, false)

	return ast.NewBranch(ifBlock, elseBlock, elseIfBlocks...), nil
}

func branchOrPluginComment(bop1, comment1 interface{}) (ast.BranchOrPlugin, error) {
	var eop ast.BranchOrPlugin
	switch t := bop1.(type) {
	case ast.Plugin:
		t.Comment = commentBlock(comment1, false, false)
		eop = t
	case ast.Branch:
		t.IfBlock.Comment = commentBlock(comment1, false, false)
		eop = t
	default:
		return nil, fmt.Errorf("invalid value for if block")
	}

	return eop, nil
}

func ifBlock(cond, bops, comment1 interface{}) (ast.IfBlock, error) {
	ib := ast.NewIfBlock(cond.(ast.Condition), branchOrPlugins(bops)...)
	ib.FooterComment = commentBlock(comment1, false, false)
	return ib, nil
}

func elseIfBlock(cond, bops, comment1 interface{}) (ast.ElseIfBlock, error) {
	eib := ast.NewElseIfBlock(cond.(ast.Condition), branchOrPlugins(bops)...)
	eib.FooterComment = commentBlock(comment1, false, false)
	return eib, nil
}

func elseBlock(bops, comment1 interface{}) (ast.ElseBlock, error) {
	eb := ast.NewElseBlock(branchOrPlugins(bops)...)
	eb.FooterComment = commentBlock(comment1, false, false)
	return eb, nil
}

func branchOrPlugins(bops1 interface{}) []ast.BranchOrPlugin {
	bops := toIfaceSlice(bops1)

	var branchOrPlugins []ast.BranchOrPlugin
	for _, bop := range bops {
		if bop, ok := bop.(ast.BranchOrPlugin); ok {
			branchOrPlugins = append(branchOrPlugins, bop)
		} else {
			return []ast.BranchOrPlugin{}
		}
	}

	return branchOrPlugins
}

func condition(expr, exprs interface{}) (ast.Condition, error) {
	iexprs := toIfaceSlice(expr)
	iexprs = append(iexprs, toIfaceSlice(exprs)...)

	var expressions []ast.Expression
	for _, ex := range iexprs {
		if ex, ok := ex.(ast.Expression); ok {
			expressions = append(expressions, ex)
		} else {
			return ast.Condition{}, fmt.Errorf("Argument is not an expression: %#v", ex)
		}
	}

	return ast.NewCondition(expressions...), nil
}

func expression(bo, expr1 interface{}) (ast.Expression, error) {
	expr := expr1.(ast.Expression)
	expr.SetBoolOperator(bo.(ast.BooleanOperator))
	return expr, nil
}

func conditionExpression(cond interface{}) (ast.ConditionExpression, error) {
	return ast.NewConditionExpression(ast.NoOperator, cond.(ast.Condition)), nil
}

func negativeExpression(cond interface{}) (ast.NegativeConditionExpression, error) {
	return ast.NewNegativeConditionExpression(ast.NoOperator, cond.(ast.Condition)), nil
}

func negativeSelector(sel interface{}) (ast.NegativeSelectorExpression, error) {
	return ast.NewNegativeSelectorExpression(ast.NoOperator, sel.(ast.Selector)), nil
}

func inExpression(lv, rv interface{}) (ast.InExpression, error) {
	return ast.NewInExpression(ast.NoOperator, lv.(ast.Rvalue), rv.(ast.Rvalue)), nil
}

func notInExpression(lv, rv interface{}) (ast.NotInExpression, error) {
	return ast.NewNotInExpression(ast.NoOperator, lv.(ast.Rvalue), rv.(ast.Rvalue)), nil
}

func compareExpression(lv, co, rv interface{}) (ast.CompareExpression, error) {
	return ast.NewCompareExpression(ast.NoOperator, lv.(ast.Rvalue), co.(ast.CompareOperator), rv.(ast.Rvalue)), nil
}

func regexpExpression(lv, ro, rv interface{}) (ast.RegexpExpression, error) {
	return ast.NewRegexpExpression(ast.NoOperator, lv.(ast.Rvalue), ro.(ast.RegexpOperator), rv.(ast.StringOrRegexp)), nil
}

func rvalue(rv interface{}) (ast.RvalueExpression, error) {
	return ast.NewRvalueExpression(ast.NoOperator, rv.(ast.Rvalue)), nil
}

func compareOperator(value string) (ast.CompareOperator, error) {
	switch value {
	case "==":
		return ast.Equal, nil
	case "!=":
		return ast.NotEqual, nil
	case "<=":
		return ast.LessOrEqual, nil
	case ">=":
		return ast.GreaterOrEqual, nil
	case "<":
		return ast.LessThan, nil
	case ">":
		return ast.GreaterThan, nil
	}
	return ast.Undefined, nil
}

func regexpOperator(value string) (ast.RegexpOperator, error) {
	switch value {
	case "=~":
		return ast.RegexpMatch, nil
	case "!~":
		return ast.RegexpNotMatch, nil
	}
	return ast.Undefined, nil
}

func booleanOperator(value string) (ast.BooleanOperator, error) {
	switch value {
	case "and":
		return ast.And, nil
	case "or":
		return ast.Or, nil
	case "xor":
		return ast.Xor, nil
	case "nand":
		return ast.Nand, nil
	}
	return ast.Undefined, nil
}

func selector(ses1 interface{}) (ast.Selector, error) {
	ises := toIfaceSlice(ses1)

	var ses []ast.SelectorElement
	for _, se := range ises {
		ses = append(ses, se.(ast.SelectorElement))
	}
	return ast.NewSelector(ses), nil
}

func selectorElement(value string) (ast.SelectorElement, error) {
	return ast.NewSelectorElement(value[1 : len(value)-1]), nil
}

func enclosedValue(c *current) (string, error) {
	return string(c.text[1 : len(c.text)-1]), nil
}

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	switch v1 := v.(type) {
	case []interface{}:
		return v1
	case interface{}:
		return []interface{}{v1}
	default:
		return nil
	}
}
