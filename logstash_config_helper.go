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

func (c *current) initState() error {
	if _, ok := c.state[exceptionalCommentWarning]; !ok {
		warnings := make(exceptionalCommentsWarnings, 0)
		c.state[exceptionalCommentWarning] = warnings
	}
	return nil
}

func (c *current) initParser() (bool, error) {
	farthestFailure = []errPos{}
	return true, nil
}

func (c *current) ret(el interface{}) (interface{}, error) {
	return el, nil
}

func (c *current) retConfig(el interface{}) (interface{}, error) {
	conf, _ := el.(ast.Config)

	warnings, _ := c.state[exceptionalCommentWarning].(exceptionalCommentsWarnings)

	conf.Warnings = warnings
	return conf, nil
}

func (c *current) commentBlock(comment1 interface{}, spaceBefore bool, spaceAfter bool) ast.CommentBlock {
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

func (c *current) configSection(ps1, psComment1 interface{}) (ast.PluginSection, error) {
	ps, ok := ps1.(ast.PluginSection)
	if !ok {
		return ast.PluginSection{}, fmt.Errorf("Value is not a PluginSection: %#v", ps1)
	}

	psComments := c.commentBlock(psComment1, false, false)

	ps.CommentBlock = psComments

	return ps, nil
}

func (c *current) config(ps1, pss1, psComment1, footerComment1 interface{}) (ast.Config, error) {
	var (
		input  []ast.PluginSection
		filter []ast.PluginSection
		output []ast.PluginSection
	)

	ips := toIfaceSlice(ps1)
	ips = append(ips, toIfaceSlice(pss1)...)

	psComment := c.commentBlock(psComment1, false, true)
	footerComments := c.commentBlock(footerComment1, true, false)

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

func (c *current) warnComment() error {
	ok, _ := c.globalStore[exceptionalCommentWarning].(bool)
	if ok && bytes.Contains(c.text, []byte("#")) {
		warnings, _ := c.state[exceptionalCommentWarning].(exceptionalCommentsWarnings)
		warnings = append(warnings, fmt.Sprintf("exceptional comment at %s", c.pos))
		c.state[exceptionalCommentWarning] = warnings
	}
	return nil
}

func (c *current) whitespace() (ast.Whitespace, error) {
	return ast.Whitespace{}, nil
}

func (c *current) comment() ([]ast.Comment, error) {
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

func (c *current) pluginSection(pt1, bops1, footerComment1 interface{}) (ast.PluginSection, error) {
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
		FooterComment:   c.commentBlock(footerComment1, false, false),
	}, nil
}

func (c *current) plugin(name, attributes1, comment1, footerComment1 interface{}) (ast.Plugin, error) {
	var attributes []ast.Attribute
	if attributes1 != nil {
		attributes = attributes1.([]ast.Attribute)
	}

	p := ast.NewPlugin(name.(string), attributes...)

	p.Comment = c.commentBlock(comment1, false, false)
	p.FooterComment = c.commentBlock(footerComment1, false, false)

	return p, nil
}

func (c *current) attributes(attribute1, attributes1, comment1 interface{}) ([]ast.Attribute, error) {
	attribute, _ := c.attributeComment(attribute1, comment1, false)
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

func (c *current) attributeComment(attribute1, comment1 interface{}, spaceBefore bool) (ast.Attribute, error) {
	var attribute ast.Attribute
	switch attr := attribute1.(type) {
	case ast.StringAttribute:
		attr.Comment = c.commentBlock(comment1, spaceBefore, false)
		attribute = attr
	case ast.NumberAttribute:
		attr.Comment = c.commentBlock(comment1, spaceBefore, false)
		attribute = attr
	case ast.ArrayAttribute:
		attr.Comment = c.commentBlock(comment1, spaceBefore, false)
		attribute = attr
	case ast.HashAttribute:
		attr.Comment = c.commentBlock(comment1, spaceBefore, false)
		attribute = attr
	case ast.PluginAttribute:
		attr.Comment = c.commentBlock(comment1, spaceBefore, false)
		attribute = attr
	default:
		return nil, fmt.Errorf("Unsupported attribute type %#v", attribute1)
	}
	return attribute, nil
}

func (c *current) attribute(name, value interface{}) (ast.Attribute, error) {
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

func (c *current) regexp() (ast.Regexp, error) {
	val, _ := c.enclosedValue()
	return ast.NewRegexp(val), nil
}

func (c *current) number() (ast.NumberAttribute, error) {
	f, err := strconv.ParseFloat(string(c.text), 64)
	if err != nil {
		// TODO: is this possible to happen? are all values, which are valid floats in Logstash/Ruby also valid floats in Go?
		return ast.NumberAttribute{}, err
	}
	return ast.NewNumberAttribute("", f), nil
}

func (c *current) array(attributes1, footerComment1 interface{}) (ast.ArrayAttribute, error) {
	var attributes []ast.Attribute
	if attributes1 != nil {
		attributes = attributes1.([]ast.Attribute)
	}

	a := ast.NewArrayAttribute("", attributes...)

	a.FooterComment = c.commentBlock(footerComment1, false, false)

	return a, nil
}

func (c *current) hash(attributes1, footerComment1 interface{}) (ast.HashAttribute, error) {
	var hashentries []ast.HashEntry
	if attributes1 != nil {
		hashentries = attributes1.([]ast.HashEntry)
	}

	a := ast.NewHashAttribute("", hashentries...)

	a.FooterComment = c.commentBlock(footerComment1, false, false)

	return a, nil
}

func (c *current) hashentries(attribute, attributes1 interface{}) ([]ast.HashEntry, error) {
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

func (c *current) hashentry(name, value, comment interface{}) (ast.HashEntry, error) {
	key := name.(ast.HashEntryKey)

	he := ast.NewHashEntry(key, value.(ast.Attribute))
	he.Comment = c.commentBlock(comment, true, false)

	return he, nil
}

func (c *current) elseIfComment(eib1, eibComment1 interface{}) (ast.ElseIfBlock, error) {
	eib := eib1.(ast.ElseIfBlock)
	eib.Comment = c.commentBlock(eibComment1, false, false)
	return eib, nil
}

func (c *current) elseComment(eb1, ebComment1 interface{}) (ast.ElseBlock, error) {
	eb := eb1.(ast.ElseBlock)
	eb.Comment = c.commentBlock(ebComment1, false, false)
	return eb, nil
}

func (c *current) branch(ifBlock1, elseIfBlocks1, elseBlock1, ifComment interface{}) (ast.Branch, error) {
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
	ifBlock.Comment = c.commentBlock(ifComment, false, false)

	return ast.NewBranch(ifBlock, elseBlock, elseIfBlocks...), nil
}

func (c *current) branchOrPluginComment(bop1, comment1 interface{}) (ast.BranchOrPlugin, error) {
	var eop ast.BranchOrPlugin
	switch t := bop1.(type) {
	case ast.Plugin:
		t.Comment = c.commentBlock(comment1, false, false)
		eop = t
	case ast.Branch:
		t.IfBlock.Comment = c.commentBlock(comment1, false, false)
		eop = t
	default:
		return nil, fmt.Errorf("invalid value for if block")
	}

	return eop, nil
}

func (c *current) ifBlock(cond, bops, comment1 interface{}) (ast.IfBlock, error) {
	ib := ast.NewIfBlock(cond.(ast.Condition), c.branchOrPlugins(bops)...)
	ib.FooterComment = c.commentBlock(comment1, false, false)
	return ib, nil
}

func (c *current) elseIfBlock(cond, bops, comment1 interface{}) (ast.ElseIfBlock, error) {
	eib := ast.NewElseIfBlock(cond.(ast.Condition), c.branchOrPlugins(bops)...)
	eib.FooterComment = c.commentBlock(comment1, false, false)
	return eib, nil
}

func (c *current) elseBlock(bops, comment1 interface{}) (ast.ElseBlock, error) {
	eb := ast.NewElseBlock(c.branchOrPlugins(bops)...)
	eb.FooterComment = c.commentBlock(comment1, false, false)
	return eb, nil
}

func (c *current) branchOrPlugins(bops1 interface{}) []ast.BranchOrPlugin {
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

func (c *current) condition(expr, exprs interface{}) (ast.Condition, error) {
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

func (c *current) expression(bo, expr1 interface{}) (ast.Expression, error) {
	expr := expr1.(ast.Expression)
	expr.SetBoolOperator(bo.(ast.BooleanOperator))
	return expr, nil
}

func (c *current) conditionExpression(cond interface{}) (ast.ConditionExpression, error) {
	return ast.NewConditionExpression(ast.NoOperator, cond.(ast.Condition)), nil
}

func (c *current) negativeExpression(cond interface{}) (ast.NegativeConditionExpression, error) {
	return ast.NewNegativeConditionExpression(ast.NoOperator, cond.(ast.Condition)), nil
}

func (c *current) negativeSelector(sel interface{}) (ast.NegativeSelectorExpression, error) {
	return ast.NewNegativeSelectorExpression(ast.NoOperator, sel.(ast.Selector)), nil
}

func (c *current) inExpression(lv, rv interface{}) (ast.InExpression, error) {
	return ast.NewInExpression(ast.NoOperator, lv.(ast.Rvalue), rv.(ast.Rvalue)), nil
}

func (c *current) notInExpression(lv, rv interface{}) (ast.NotInExpression, error) {
	return ast.NewNotInExpression(ast.NoOperator, lv.(ast.Rvalue), rv.(ast.Rvalue)), nil
}

func (c *current) compareExpression(lv, co, rv interface{}) (ast.CompareExpression, error) {
	return ast.NewCompareExpression(ast.NoOperator, lv.(ast.Rvalue), co.(ast.CompareOperator), rv.(ast.Rvalue)), nil
}

func (c *current) regexpExpression(lv, ro, rv interface{}) (ast.RegexpExpression, error) {
	return ast.NewRegexpExpression(ast.NoOperator, lv.(ast.Rvalue), ro.(ast.RegexpOperator), rv.(ast.StringOrRegexp)), nil
}

func (c *current) rvalue(rv interface{}) (ast.RvalueExpression, error) {
	return ast.NewRvalueExpression(ast.NoOperator, rv.(ast.Rvalue)), nil
}

func (c *current) compareOperator() (ast.CompareOperator, error) {
	switch string(c.text) {
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

func (c *current) regexpOperator() (ast.RegexpOperator, error) {
	switch string(c.text) {
	case "=~":
		return ast.RegexpMatch, nil
	case "!~":
		return ast.RegexpNotMatch, nil
	}
	return ast.Undefined, nil
}

func (c *current) booleanOperator() (ast.BooleanOperator, error) {
	switch string(c.text) {
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

func (c *current) selector(ses1 interface{}) (ast.Selector, error) {
	ises := toIfaceSlice(ses1)

	var ses []ast.SelectorElement
	for _, se := range ises {
		ses = append(ses, se.(ast.SelectorElement))
	}
	return ast.NewSelector(ses), nil
}

func (c *current) selectorElement() (ast.SelectorElement, error) {
	value := string(c.text)
	return ast.NewSelectorElement(value[1 : len(value)-1]), nil
}

func (c *current) enclosedValue() (string, error) {
	return string(c.text[1 : len(c.text)-1]), nil
}

func (c *current) string() (string, error) {
	return string(c.text), nil
}

func (c *current) astPos() ast.Pos {
	p := ast.Pos{
		Line:   c.pos.line,
		Column: c.pos.col,
		Offset: c.pos.offset,
	}
	return p
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
