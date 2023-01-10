package ast

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Node interface {
	Pos() Pos
}

var (
	_ Node = Config{}
	_ Node = PluginSection{}
	_ Node = Plugin{}
	_ Node = Attribute(nil)
	_ Node = Branch{}
)

type Pos struct {
	Line   int
	Column int
	Offset int
}

func (p Pos) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.Line, p.Column, p.Offset)
}

var InvalidPos = Pos{Offset: -1}

func (c Config) Pos() Pos           { return Pos{Line: 1, Column: 1, Offset: 0} }
func (ps PluginSection) Pos() Pos   { return ps.Start }
func (p Plugin) Pos() Pos           { return p.Start }
func (pa PluginAttribute) Pos() Pos { return pa.Start }
func (sa StringAttribute) Pos() Pos { return sa.Start }
func (na NumberAttribute) Pos() Pos { return na.Start }
func (aa ArrayAttribute) Pos() Pos  { return aa.Start }
func (ha HashAttribute) Pos() Pos   { return ha.Start }
func (he HashEntry) Pos() Pos       { return he.Key.Pos() }
func (b Branch) Pos() Pos           { return b.IfBlock.Start }
func (ib IfBlock) Pos() Pos         { return ib.Start }
func (eib ElseIfBlock) Pos() Pos    { return eib.Start }
func (eb ElseBlock) Pos() Pos       { return eb.Start }
func (c Condition) Pos() Pos {
	if len(c.Expression) == 0 {
		return InvalidPos
	}
	return c.Expression[0].Pos()
}
func (be BoolExpression) Pos() Pos              { return be.Start }
func (ce ConditionExpression) Pos() Pos         { return ce.Start }
func (nc NegativeConditionExpression) Pos() Pos { return nc.Start }
func (ns NegativeSelectorExpression) Pos() Pos  { return ns.Start }
func (ie InExpression) Pos() Pos                { return ie.Start }
func (nie NotInExpression) Pos() Pos            { return nie.Start }
func (re RvalueExpression) Pos() Pos            { return re.Start }
func (ce CompareExpression) Pos() Pos           { return ce.Start }
func (co CompareOperator) Pos() Pos             { return co.Start }
func (re RegexpExpression) Pos() Pos            { return re.Start }
func (r Regexp) Pos() Pos                       { return r.Start }
func (ro RegexpOperator) Pos() Pos              { return ro.Start }
func (bo BooleanOperator) Pos() Pos             { return bo.Start }
func (s Selector) Pos() Pos                     { return s.Start }
func (se SelectorElement) Pos() Pos             { return se.Start }

// A Config node represents the root node of a Logstash configuration.
type Config struct {
	Input         []PluginSection
	Filter        []PluginSection
	Output        []PluginSection
	FooterComment CommentBlock
	Warnings      []string
}

// NewConfig creates a new Logstash config.
func NewConfig(input, filter, output []PluginSection) Config {
	return Config{
		Input:  input,
		Filter: filter,
		Output: output,
	}
}

// String returns a string representation of a Logstash configuration.
func (c Config) String() string {
	var s bytes.Buffer

	s.WriteString(pluginSectionString("input", c.Input))
	s.WriteString(pluginSectionString("filter", c.Filter))
	s.WriteString(pluginSectionString("output", c.Output))
	for _, c := range c.FooterComment {
		s.WriteString(c.String())
	}

	return s.String()
}

type Whitespace struct{}

type CommentBlock []Comment

func NewCommentBlock(comments ...Comment) CommentBlock {
	return comments
}

func (cb CommentBlock) String() string {
	var s bytes.Buffer
	for _, c := range cb {
		s.WriteString(c.String())
	}
	return s.String()
}

type Comment struct {
	comment     string
	SpaceBefore bool
	SpaceAfter  bool
}

func NewComment(comment string, space bool) Comment {
	return Comment{
		comment:     comment,
		SpaceBefore: space,
	}
}

func (c Comment) String() string {
	var s bytes.Buffer
	if c.SpaceBefore {
		s.WriteString("\n")
	}
	s.WriteString(fmt.Sprintf("# %s\n", c.comment))
	if c.SpaceAfter {
		s.WriteString("\n")
	}
	return s.String()
}

func pluginSectionString(pluginType string, ps []PluginSection) string {
	if len(ps) == 0 {
		return ""
	}

	var s bytes.Buffer
	for i, p := range ps {
		if len(ps) > 0 {
			s.WriteString(ps[0].CommentBlock.String())
		}
		s.WriteString(fmt.Sprint(pluginType + " {"))
		s.WriteString(prefix(p.String(), false))
		s.WriteString(prefix(p.FooterComment.String(), false))
		s.WriteString(fmt.Sprintln("}"))
		if i < len(ps)-1 {
			s.WriteString("\n")
		}
	}
	return s.String()
}

// Undefined is a placeholder for all undefined values in all available types in this package.
const Undefined = 0

const (
	// Undefined is already defined

	// Input type plugin
	Input = iota + 1
	// Filter type plugin
	Filter
	// Output type plugin
	Output
)

// PluginType defines the type of a Logstash plugin, which is one of:
// Input, Filter or Output.
type PluginType int

// String returns a string representation of a plugin type.
func (pt PluginType) String() string {
	switch pt {
	case Input:
		return "input"
	case Filter:
		return "filter"
	case Output:
		return "output"
	default:
		return "undefined plugin type"
	}
}

// A PluginSection node defines the configuration section with branches or plugins.
type PluginSection struct {
	Start           Pos
	PluginType      PluginType
	BranchOrPlugins []BranchOrPlugin
	CommentBlock    CommentBlock
	FooterComment   CommentBlock
}

// NewPluginSection creates a new plugin section.
func NewPluginSection(pt PluginType, bop ...BranchOrPlugin) PluginSection {
	return PluginSection{
		PluginType:      pt,
		BranchOrPlugins: bop,
	}
}

// NewPluginSections creates an array of plugin sections.
func NewPluginSections(pt PluginType, bop ...BranchOrPlugin) []PluginSection {
	return []PluginSection{
		NewPluginSection(pt, bop...),
	}
}

// String returns a string representation of a plugin section.
func (ps PluginSection) String() string {
	var s bytes.Buffer
	for _, bop := range ps.BranchOrPlugins {
		if bop == nil {
			continue
		}
		if s.Len() > 0 {
			s.WriteString("\n")
		}
		s.WriteString(fmt.Sprintf("%v", bop))
		if s.Len() > 0 {
			s.WriteString("\n")
		}
	}
	return s.String()
}

// BranchOrPlugin interface combines Logstash configuration conditional branches and plugins.
type BranchOrPlugin interface {
	Pos() Pos
	branchOrPlugin()
}

// branchOrPlugin() ensures that only BranchOrPlugin/type nodes can be
// assigned to an BranchOrPlugin.
func (Plugin) branchOrPlugin() {}
func (Branch) branchOrPlugin() {}

// A Plugin node represents a Logstash plugin.
type Plugin struct {
	Start         Pos
	name          string
	Attributes    []Attribute
	Comment       CommentBlock
	FooterComment CommentBlock
}

// NewPlugin creates a new plugin.
func NewPlugin(name string, attributes ...Attribute) Plugin {
	return Plugin{
		name:       name,
		Attributes: attributes,
	}
}

// Name returns the name of the attribute.
func (p Plugin) Name() string {
	return p.name
}

// String returns a string representation of a plugin.
func (p Plugin) String() string {
	var s bytes.Buffer

	s.WriteString(p.Comment.String())

	s.WriteString(fmt.Sprint(p.Name(), " {"))

	var ss bytes.Buffer
	for _, attr := range p.Attributes {
		if attr == nil {
			continue
		}
		ss.WriteString(attr.CommentBlock())
		ss.WriteString(attr.String())
		ss.WriteString("\n")
	}
	if ss.Len() > 0 {
		ss.WriteString("\n")
	}
	ss.WriteString(p.FooterComment.String())
	s.WriteString(prefix(ss.String(), false))

	s.WriteString("}")
	return s.String()
}

// ID returns the id of a Logstash plugin.
// The id attribute is one of the common options, that is optionally available
// on every Logstash plugin. In generall, it is highly recommended for a Logstash
// plugin to have an id.
// If the ID attribute is not present, an error is returned, who implements
// the NotFounder interface.
func (p Plugin) ID() (string, error) {
	for _, attr := range p.Attributes {
		if attr != nil && attr.Name() == "id" {
			switch stringAttr := attr.(type) {
			case StringAttribute:
				return stringAttr.value, nil
			default:
				return "", errors.New("attribut id is not of type string attribute")
			}
		}
	}
	return "", NotFoundErrorf("plugin %s does not contain an id attribute", p.name)
}

// Attribute interface combines Logstash plugin attribute types.
type Attribute interface {
	Name() string
	String() string
	ValueString() string
	CommentBlock() string
	Pos() Pos
	attributeNode()
}

// attributeNode() ensures that only attribute/type nodes can be
// assigned to an Attribute.
func (PluginAttribute) attributeNode() {}
func (StringAttribute) attributeNode() {}
func (NumberAttribute) attributeNode() {}
func (ArrayAttribute) attributeNode()  {}
func (HashAttribute) attributeNode()   {}

// A PluginAttribute node represents a plugin attribute of type plugin.
type PluginAttribute struct {
	Start   Pos
	name    string
	value   Plugin
	Comment CommentBlock
}

// NewPluginAttribute creates a new plugin attribute.
func NewPluginAttribute(name string, value Plugin) PluginAttribute {
	return PluginAttribute{
		name:  name,
		value: value,
	}
}

// Name returns the name of the attribute.
func (pa PluginAttribute) Name() string {
	return pa.name
}

// String returns a string representation of a plugin attribute.
func (pa PluginAttribute) String() string {
	return fmt.Sprintf("%s => %s", pa.Name(), pa.ValueString())
}

// ValueString returns the value of the node as a string representation.
func (pa PluginAttribute) ValueString() string {
	return pa.value.String()
}

// CommentBlock returns the comment of the node.
func (pa PluginAttribute) CommentBlock() string {
	return pa.Comment.String()
}

func (pa PluginAttribute) Value() Plugin {
	return pa.value
}

const (
	// Undefined is already defined

	// DoubleQuoted string attribute type
	DoubleQuoted = iota + 1

	// SingleQuoted string attribute type
	SingleQuoted

	// Bareword string attribute type
	Bareword
)

// StringAttributeType defines the string format type of a string attribute.
type StringAttributeType int

// String returns a string representation of a string attribute type.
func (sat StringAttributeType) String() string {
	switch sat {
	case DoubleQuoted:
		return `"`
	case SingleQuoted:
		return `'`
	case Bareword:
		return ``
	default:
		return "undefined string attribute type"
	}
}

// StringAttribute is a plugin attribute of type string.
type StringAttribute struct {
	Start   Pos
	name    string
	value   string
	sat     StringAttributeType
	Comment CommentBlock
}

// NewStringAttribute creates a new plugin attribute of type string.
func NewStringAttribute(name, value string, sat StringAttributeType) StringAttribute {
	return StringAttribute{
		name:  name,
		value: value,
		sat:   sat,
	}
}

// Name returns the name of the attribute.
func (sa StringAttribute) Name() string {
	return sa.name
}

// String returns a string representation of a string attribute.
func (sa StringAttribute) String() string {
	return fmt.Sprintf("%s => %s", sa.Name(), sa.ValueString())
}

// ValueString returns the value of the node as a string representation.
func (sa StringAttribute) ValueString() string {
	return fmt.Sprintf("%s%s%s", sa.StringAttributeType(), sa.Value(), sa.StringAttributeType())
}

// CommentBlock returns the comment of the node.
func (sa StringAttribute) CommentBlock() string {
	return sa.Comment.String()
}

// Value returns the value of the node.
func (sa StringAttribute) Value() string {
	return sa.value
}

// StringAttributeType returns the string attribute type.
func (sa StringAttribute) StringAttributeType() StringAttributeType {
	return sa.sat
}

// A NumberAttribute node represents a plugin attribute of type number.
type NumberAttribute struct {
	Start   Pos
	name    string
	value   float64
	Comment CommentBlock
}

// NewNumberAttribute creates a new number attribute.
func NewNumberAttribute(name string, value float64) NumberAttribute {
	return NumberAttribute{
		name:  name,
		value: value,
	}
}

// Name returns the name of the attribute.
func (na NumberAttribute) Name() string {
	return na.name
}

// String returns a string representation of a number attribute.
func (na NumberAttribute) String() string {
	return fmt.Sprintf("%s => %s", na.Name(), na.ValueString())
}

// ValueString returns the value of the node as a string representation.
func (na NumberAttribute) ValueString() string {
	value := fmt.Sprintf("%g", na.Value())
	if strings.Contains(value, "e") {
		if float64(int64(na.Value())+0) == na.Value() {
			return fmt.Sprintf("%d", int64(na.Value()))
		}
		return strings.TrimRight(fmt.Sprintf("%.10f", na.Value()), "0")
	}
	return value
}

// CommentBlock returns the comment of the node.
func (na NumberAttribute) CommentBlock() string {
	return na.Comment.String()
}

// Value returns the value of the node.
func (na NumberAttribute) Value() float64 {
	return na.value
}

// A ArrayAttribute node represents a plugin attribute of type array.
type ArrayAttribute struct {
	Start         Pos
	name          string
	Attributes    []Attribute
	Comment       CommentBlock
	FooterComment CommentBlock
}

// NewArrayAttribute creates a new array attribute.
func NewArrayAttribute(name string, value ...Attribute) ArrayAttribute {
	return ArrayAttribute{
		name:       name,
		Attributes: value,
	}
}

// Name returns the name of the attribute.
func (aa ArrayAttribute) Name() string {
	return aa.name
}

// String returns a string representation of a array attribute.
func (aa ArrayAttribute) String() string {
	return fmt.Sprintf("%s => %s", aa.Name(), aa.ValueString())
}

// ValueString returns the value of the node as a string representation.
func (aa ArrayAttribute) ValueString() string {
	var s bytes.Buffer
	s.WriteString("[")

	var ss bytes.Buffer
	first := true
	for _, a := range aa.Value() {
		if a == nil {
			continue
		}
		if first {
			first = false
		} else {
			ss.WriteString(",\n")
		}
		ss.WriteString(a.CommentBlock())
		ss.WriteString(a.ValueString())
	}
	if ss.Len() > 0 {
		ss.WriteString("\n\n")
	}
	ss.WriteString(aa.FooterComment.String())
	s.WriteString(prefix(ss.String(), false))
	s.WriteString("]")

	return s.String()
}

// CommentBlock returns the comment of the node.
func (aa ArrayAttribute) CommentBlock() string {
	return aa.Comment.String()
}

// Value returns the value of the node.
func (aa ArrayAttribute) Value() []Attribute {
	return aa.Attributes
}

// A HashAttribute node represents a plugin attribute of type hash.
type HashAttribute struct {
	Start         Pos
	name          string
	Entries       []HashEntry
	Comment       CommentBlock
	FooterComment CommentBlock
}

// NewHashAttribute creates a new hash attribute.
func NewHashAttribute(name string, entries ...HashEntry) HashAttribute {
	return HashAttribute{
		name:    name,
		Entries: entries,
	}
}

// Name returns the name of the attribute.
func (ha HashAttribute) Name() string {
	return ha.name
}

// String returns a string representation of a hash attribute.
func (ha HashAttribute) String() string {
	return fmt.Sprintf("%s => %s", ha.Name(), ha.ValueString())
}

// ValueString returns the value of the node as a string representation.
func (ha HashAttribute) ValueString() string {
	var s bytes.Buffer
	s.WriteString("{")

	var ss bytes.Buffer
	for _, v := range ha.Value() {
		ss.WriteString(v.String())
		ss.WriteString("\n")
	}
	if ss.Len() > 0 {
		ss.WriteString("\n")
	}
	ss.WriteString(ha.FooterComment.String())
	s.WriteString(prefix(ss.String(), false))

	s.WriteString("}")
	return s.String()
}

// CommentBlock returns the comment of the node.
func (ha HashAttribute) CommentBlock() string {
	return ha.Comment.String()
}

// Value returns the value of the node.
func (ha HashAttribute) Value() []HashEntry {
	return ha.Entries
}

type HashEntryKey interface {
	Pos() Pos
	ValueString() string
	attributeNode()
	hashEntryKeyAttribute()
}

func (NumberAttribute) hashEntryKeyAttribute() {}
func (StringAttribute) hashEntryKeyAttribute() {}

// A HashEntry node defines a hash entry within a hash attribute.
type HashEntry struct {
	Start   Pos
	Key     HashEntryKey
	Value   Attribute
	Comment CommentBlock
}

// NewHashEntry creates a new hash entry for a hash attribute.
func NewHashEntry(name HashEntryKey, value Attribute) HashEntry {
	return HashEntry{
		Key:   name,
		Value: value,
	}
}

// Name returns the name of the attribute.
func (he HashEntry) Name() string {
	return he.Key.ValueString()
}

// String returns a string representation of a hash entry.
func (he HashEntry) String() string {
	return fmt.Sprintf("%s%s => %s", he.Comment, he.Name(), he.ValueString())
}

// ValueString returns the value of the node as a string representation.
func (he HashEntry) ValueString() string {
	if he.Value == nil {
		return ""
	}
	return he.Value.ValueString()
}

// A Branch node represents a conditional branch within a Logstash configuration.
type Branch struct {
	IfBlock     IfBlock
	ElseIfBlock []ElseIfBlock
	ElseBlock   ElseBlock
}

// NewBranch creates a new branch.
// Arguments for elseBlock and elseIfBlock are in the wrong order from logically point of view.
// This is due to the variadic nature of the elseIfBlock argument.
func NewBranch(ifBlock IfBlock, elseBlock ElseBlock, elseIfBlock ...ElseIfBlock) Branch {
	return Branch{
		IfBlock:     ifBlock,
		ElseIfBlock: elseIfBlock,
		ElseBlock:   elseBlock,
	}
}

// TODO: Maybe we should add helper functions NewIfBranch and NewIfElseBranch

// String returns a string representation of a branch.
func (b Branch) String() string {
	var s bytes.Buffer
	s.WriteString(b.IfBlock.String())
	for _, block := range b.ElseIfBlock {
		s.WriteString(block.String())
	}
	s.WriteString(b.ElseBlock.String())
	return s.String()
}

// A IfBlock node represents an if-block of a Branch.
type IfBlock struct {
	Start         Pos
	Condition     Condition
	Block         []BranchOrPlugin
	Comment       CommentBlock
	FooterComment CommentBlock
}

// NewIfBlock creates a new if-block.
func NewIfBlock(condition Condition, block ...BranchOrPlugin) IfBlock {
	return IfBlock{
		Condition: condition,
		Block:     block,
	}
}

// String returns a string representation of an if-block.
func (ib IfBlock) String() string {
	var s bytes.Buffer
	s.WriteString(ib.Comment.String())
	s.WriteString(fmt.Sprintf("if %v {", ib.Condition))

	var ss bytes.Buffer
	for _, block := range ib.Block {
		if block == nil {
			continue
		}
		if ss.Len() > 0 {
			ss.WriteString("\n")
		}
		ss.WriteString(fmt.Sprintf("%v", block))
		if ss.Len() > 0 {
			ss.WriteString("\n")
		}
	}
	if ss.Len() > 0 {
		ss.WriteString("\n")
	}
	ss.WriteString(ib.FooterComment.String())
	s.WriteString(prefix(ss.String(), true))

	s.WriteString("}")
	return s.String()
}

// A ElseIfBlock node represents an else-if-block of a Branch.
type ElseIfBlock struct {
	Start         Pos
	Condition     Condition
	Block         []BranchOrPlugin
	Comment       CommentBlock
	FooterComment CommentBlock
}

// NewElseIfBlock creates a new else-if-block of a Branch.
func NewElseIfBlock(condition Condition, block ...BranchOrPlugin) ElseIfBlock {
	return ElseIfBlock{
		Condition: condition,
		Block:     block,
	}
}

// String returns a string representation of an else if block.
func (eib ElseIfBlock) String() string {
	var s bytes.Buffer
	if len(eib.Comment) > 0 {
		s.WriteString("\n")
		s.WriteString(eib.Comment.String())
	} else {
		s.WriteString(" ")
	}
	s.WriteString(fmt.Sprintf("else if %v {", eib.Condition))

	var ss bytes.Buffer
	for _, block := range eib.Block {
		if block == nil {
			continue
		}
		if ss.Len() > 0 {
			ss.WriteString("\n")
		}
		ss.WriteString(fmt.Sprint(block))
		if ss.Len() > 0 {
			ss.WriteString("\n")
		}
	}
	if ss.Len() > 0 {
		ss.WriteString("\n")
	}
	ss.WriteString(eib.FooterComment.String())
	s.WriteString(prefix(ss.String(), true))

	s.WriteString("}")
	return s.String()
}

// A ElseBlock node represents a else-block of a Branch.
type ElseBlock struct {
	Start         Pos
	Block         []BranchOrPlugin
	Comment       CommentBlock
	FooterComment CommentBlock
}

// NewElseBlock creates a new else-block
func NewElseBlock(block ...BranchOrPlugin) ElseBlock {
	return ElseBlock{
		Block: block,
	}
}

// String returns a string representation of an else block.
func (eb ElseBlock) String() string {
	if len(eb.Block) == 0 && len(eb.Comment) == 0 && len(eb.FooterComment) == 0 {
		return ""
	}

	var s bytes.Buffer
	if len(eb.Comment) > 0 {
		s.WriteString("\n")
		s.WriteString(eb.Comment.String())
	} else {
		s.WriteString(" ")
	}
	s.WriteString("else {")
	var ss bytes.Buffer
	for _, block := range eb.Block {
		if block == nil {
			continue
		}
		if ss.Len() > 0 {
			ss.WriteString("\n")
		}
		ss.WriteString(fmt.Sprint(block))
		if ss.Len() > 0 {
			ss.WriteString("\n")
		}
	}
	if ss.Len() > 0 {
		ss.WriteString("\n")
	}
	ss.WriteString(eb.FooterComment.String())
	s.WriteString(prefix(ss.String(), true))
	s.WriteString("}")
	return s.String()
}

// A Condition node represents a condition used by if- or else-if-blocks.
type Condition struct {
	Expression []Expression
}

// NewCondition creates a new condition.
func NewCondition(expression ...Expression) Condition {
	return Condition{
		Expression: expression,
	}
}

// String returns a string representation of a condition.
func (c Condition) String() string {
	var s bytes.Buffer
	for _, expression := range c.Expression {
		if expression == nil {
			continue
		}
		s.WriteString(fmt.Sprint(expression))
	}
	return s.String()
}

// An Expression node defines an expression.
// An Expression is chainable with a preceding Expression by
// the the boolean operator.
type Expression interface {
	Pos() Pos
	BoolOperator() BooleanOperator
	SetBoolOperator(BooleanOperator)
	expressionNode()
}

// expressionNode() ensures that only expression/type nodes can be
// assigned to an Expression.
func (ConditionExpression) expressionNode()         {}
func (NegativeConditionExpression) expressionNode() {}
func (NegativeSelectorExpression) expressionNode()  {}
func (InExpression) expressionNode()                {}
func (NotInExpression) expressionNode()             {}
func (CompareExpression) expressionNode()           {}
func (RegexpExpression) expressionNode()            {}
func (RvalueExpression) expressionNode()            {}

// A BoolExpression node represents a boolean operator.
type BoolExpression struct {
	Start        Pos
	boolOperator BooleanOperator
}

// BoolOperator returns the boolean operator of the node.
func (be *BoolExpression) BoolOperator() BooleanOperator {
	return be.boolOperator
}

// SetBoolOperator sets the boolean operator for the node.
func (be *BoolExpression) SetBoolOperator(bo BooleanOperator) {
	be.boolOperator = bo
	be.Start = bo.Start
}

// String returns a string representation of a boolean expression.
func (be BoolExpression) String() string {
	return be.BoolOperator().String()
}

// A ConditionExpression node represents an Expression, which is enclosed in parentheses.
type ConditionExpression struct {
	Start Pos
	*BoolExpression
	Condition Condition
}

// NewConditionExpression creates a new condition expression.
func NewConditionExpression(boolOperator BooleanOperator, condition Condition) ConditionExpression {
	return ConditionExpression{
		BoolExpression: &BoolExpression{
			boolOperator: boolOperator,
		},
		Condition: condition,
	}
}

// String returns a string representation of a condition expression.
func (ce ConditionExpression) String() string {
	return fmt.Sprintf("%v(%s)", ce.BoolExpression, ce.Condition.String())
}

// A NegativeConditionExpression node represents an Expression within parentheses, which is negated.
type NegativeConditionExpression struct {
	Start Pos
	*BoolExpression
	Condition Condition
}

// NewNegativeConditionExpression creates a new negative condition expression.
func NewNegativeConditionExpression(boolOperator BooleanOperator, condition Condition) NegativeConditionExpression {
	return NegativeConditionExpression{
		BoolExpression: &BoolExpression{
			boolOperator: boolOperator,
		},
		Condition: condition,
	}
}

// String returns a string representation of a negative condition expression.
func (nc NegativeConditionExpression) String() string {
	return fmt.Sprintf("%v!(%s)", nc.BoolExpression, nc.Condition.String())
}

// A NegativeSelectorExpression node represents a field selector expression, which is negated.
type NegativeSelectorExpression struct {
	Start Pos
	*BoolExpression
	Selector Selector
}

// NewNegativeSelectorExpression creates a new negative selector expression.
func NewNegativeSelectorExpression(boolOperator BooleanOperator, selector Selector) NegativeSelectorExpression {
	return NegativeSelectorExpression{
		BoolExpression: &BoolExpression{
			boolOperator: boolOperator,
		},
		Selector: selector,
	}
}

// String returns a string representation of a negative selector expression.
func (ns NegativeSelectorExpression) String() string {
	return fmt.Sprintf("%v!%s", ns.BoolExpression, ns.Selector)
}

// An InExpression node represents an in expression.
type InExpression struct {
	Start Pos
	*BoolExpression
	LValue Rvalue
	RValue Rvalue
}

// NewInExpression creates a new in expression.
func NewInExpression(boolOperator BooleanOperator, lvalue Rvalue, rvalue Rvalue) InExpression {
	return InExpression{
		BoolExpression: &BoolExpression{
			boolOperator: boolOperator,
		},
		LValue: lvalue,
		RValue: rvalue,
	}
}

// String returns a string representation of an in expression.
func (ie InExpression) String() string {
	return fmt.Sprintf("%v%v in %v", ie.BoolExpression, ie.LValue.ValueString(), ie.RValue.ValueString())
}

// A NotInExpression node defines a not in expression.
type NotInExpression struct {
	Start Pos
	*BoolExpression
	RValue Rvalue
	LValue Rvalue
}

// NewNotInExpression creates a new not in expression.
func NewNotInExpression(boolOperator BooleanOperator, lvalue Rvalue, rvalue Rvalue) NotInExpression {
	return NotInExpression{
		BoolExpression: &BoolExpression{
			boolOperator: boolOperator,
		},
		LValue: lvalue,
		RValue: rvalue,
	}
}

// String returns a string representation of a not in expression.
func (nie NotInExpression) String() string {
	var lvalue string
	var rvalue string
	if nie.LValue != nil {
		lvalue = nie.LValue.ValueString()
	}
	if nie.RValue != nil {
		rvalue = nie.RValue.ValueString()
	}
	return fmt.Sprintf("%v%v not in %v", nie.BoolExpression, lvalue, rvalue)
}

// A Rvalue node represents an right (or in some cases also an left) side value of an expression.
type Rvalue interface {
	Pos() Pos
	String() string
	ValueString() string
	rvalueNode()
}

// rvalueNode() ensures that only rvalue/type nodes can be
// assigned to an Rvalue.
func (StringAttribute) rvalueNode() {}
func (NumberAttribute) rvalueNode() {}
func (Selector) rvalueNode()        {}
func (ArrayAttribute) rvalueNode()  {}
func (Regexp) rvalueNode()          {}

// A RvalueExpression node defines an expression consisting only of a Rvalue.
type RvalueExpression struct {
	Start Pos
	*BoolExpression
	RValue Rvalue
}

// NewRvalueExpression creates a new rvalue expression.
func NewRvalueExpression(boolOperator BooleanOperator, rvalue Rvalue) RvalueExpression {
	return RvalueExpression{
		BoolExpression: &BoolExpression{
			boolOperator: boolOperator,
		},
		RValue: rvalue,
	}
}

// String returns a string representation of a rvalue expression.
func (re RvalueExpression) String() string {
	var rvalue string
	if re.RValue != nil {
		rvalue = re.RValue.ValueString()
	}
	return fmt.Sprintf("%v%v", re.BoolExpression, rvalue)
}

// A CompareExpression node represents a expression, which compares lvalue and rvalue
// based on the comparison operator.
type CompareExpression struct {
	Start Pos
	*BoolExpression
	LValue          Rvalue
	CompareOperator CompareOperator
	RValue          Rvalue
}

// NewCompareExpression creates a new compare expression.
func NewCompareExpression(boolOperator BooleanOperator, lvalue Rvalue, compareOperator CompareOperator, rvalue Rvalue) CompareExpression {
	return CompareExpression{
		BoolExpression: &BoolExpression{
			boolOperator: boolOperator,
		},
		LValue:          lvalue,
		CompareOperator: compareOperator,
		RValue:          rvalue,
	}
}

// String returns a string representation of a compare expression.
func (ce CompareExpression) String() string {
	var lvalue string
	var rvalue string
	if ce.LValue != nil {
		lvalue = ce.LValue.ValueString()
	}
	if ce.RValue != nil {
		rvalue = ce.RValue.ValueString()
	}
	return fmt.Sprintf("%v%v %v %v", ce.BoolExpression, lvalue, ce.CompareOperator, rvalue)
}

// A CompareOperator represents the comparison operator, used to compare two values.
type CompareOperator struct {
	Op    int
	Start Pos
}

const (
	// Undefined is already defined

	// Equal defines the equal operator (==)
	Equal = iota + 1

	// NotEqual defines the not equal operator (!=)
	NotEqual

	// LessOrEqual defines the less or equal operator (<=)
	LessOrEqual

	// GreaterOrEqual defines the greater or equal operator (>=)
	GreaterOrEqual

	// LessThan defines the less than operator (<)
	LessThan

	// GreaterThan defines the greater than operator (>)
	GreaterThan
)

// String returns a string representation of a compare operator.
func (co CompareOperator) String() string {
	switch co.Op {
	case Equal:
		return "=="
	case NotEqual:
		return "!="
	case LessOrEqual:
		return "<="
	case GreaterOrEqual:
		return ">="
	case LessThan:
		return "<"
	case GreaterThan:
		return ">"
	default:
		return "undefined compare operator"
	}
}

// A RegexpExpression node defines a regular expression node.
type RegexpExpression struct {
	Start Pos
	*BoolExpression
	LValue         Rvalue
	RegexpOperator RegexpOperator
	RValue         StringOrRegexp
}

// NewRegexpExpression creates a new regexp (regular expression) expression.
func NewRegexpExpression(boolOperator BooleanOperator, lvalue Rvalue, regexpOperator RegexpOperator, rvalue StringOrRegexp) RegexpExpression {
	return RegexpExpression{
		BoolExpression: &BoolExpression{
			boolOperator: boolOperator,
		},
		LValue:         lvalue,
		RegexpOperator: regexpOperator,
		RValue:         rvalue,
	}
}

// String returns a string representation of a regexp expression.
func (re RegexpExpression) String() string {
	var lvalue string
	var rvalue string
	if re.LValue != nil {
		lvalue = re.LValue.ValueString()
	}
	if re.RValue != nil {
		rvalue = re.RValue.ValueString()
	}
	return fmt.Sprintf("%v%v %v %v", re.BoolExpression, lvalue, re.RegexpOperator, rvalue)
}

// A StringOrRegexp node is a string attribute node or a regexp node.
type StringOrRegexp interface {
	Pos() Pos
	String() string
	ValueString() string
	stringOrRegexp()
}

// stringOrRegexp() ensures that only stringOrRegexp/type nodes can be.
// assigned to an StringOrRegexp.
//
func (StringAttribute) stringOrRegexp() {}
func (Regexp) stringOrRegexp()          {}

// A RegexpOperator is an operator, used to compare a regular expression with an other value.
type RegexpOperator struct {
	Op    int
	Start Pos
}

const (
	// Undefined is already defined

	// RegexpMatch is the regular expression match operator (=~)
	RegexpMatch = iota + 1

	// RegexpNotMatch is the regular expression not match operator (!~)
	RegexpNotMatch
)

// String returns a string representation of a regexp operator.
func (ro RegexpOperator) String() string {
	switch ro.Op {
	case RegexpMatch:
		return "=~"
	case RegexpNotMatch:
		return "!~"
	default:
		return "undefined regexp operator"
	}
}

// A Regexp node represents a regular expression.
type Regexp struct {
	Start  Pos
	Regexp string
}

// NewRegexp creates a new Regexp.
func NewRegexp(regexp string) Regexp {
	return Regexp{
		Regexp: regexp,
	}
}

// String returns a string representation of a regexp.
func (r Regexp) String() string {
	return fmt.Sprintf("/%s/", r.Regexp)
}

// ValueString returns the value of the node as a string representation.
func (r Regexp) ValueString() string {
	return r.String()
}

// A BooleanOperator represents a boolean operator.
type BooleanOperator struct {
	Op    int
	Start Pos
}

const (
	// Undefined is already defined

	// NoOperator is used for the first expression, which is not chained by a boolean operator
	NoOperator = iota + 1

	// And is the and boolean operator
	And

	// Or is the or boolean operator
	Or

	// Xor is the xor boolean operator
	Xor

	// Nand is the nand boolean operator
	Nand
)

// String returns a string representation of a boolean operator.
func (be BooleanOperator) String() string {
	switch be.Op {
	case NoOperator:
		return ""
	case And:
		return " and "
	case Or:
		return " or "
	case Nand:
		return " nand "
	case Xor:
		return " xor "
	default:
		return "undefined boolean operator"
	}
}

// A Selector node represents a field selector.
type Selector struct {
	Start    Pos
	Elements []SelectorElement
}

// NewSelector creates a new Selector.
func NewSelector(elements []SelectorElement) Selector {
	return Selector{
		Elements: elements,
	}
}

// NewSelectorFromNames creates a new Selector form a slice of field names.
func NewSelectorFromNames(names ...string) Selector {
	var elements []SelectorElement
	for _, name := range names {
		elements = append(elements, NewSelectorElement(name))
	}
	return NewSelector(elements)
}

// String returns a string representation of a selector.
func (s Selector) String() string {
	var bb bytes.Buffer
	for _, element := range s.Elements {
		bb.WriteString(element.String())
	}
	return bb.String()
}

// ValueString returns the value of the node as a string representation.
func (s Selector) ValueString() string {
	return s.String()
}

// A SelectorElement node defines a selector element.
type SelectorElement struct {
	Start Pos
	name  string
}

// NewSelectorElement creates a new selector element.
func NewSelectorElement(name string) SelectorElement {
	return SelectorElement{
		name: name,
	}
}

// String returns a string representation of a selector element.
func (se SelectorElement) String() string {
	return fmt.Sprintf("[%s]", se.name)
}

// FIXME: Do I need this interface?

// A Commentable node is an ast node, which accepts comments.
type Commentable interface {
	SetComment(cb CommentBlock)
}
