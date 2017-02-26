package ast

import (
	"bytes"
	"fmt"
)

type Config struct {
	Input  []PluginSection
	Filter []PluginSection
	Output []PluginSection
}

func NewConfig(input, filter, output []PluginSection) Config {
	return Config{
		Input:  input,
		Filter: filter,
		Output: output,
	}
}

func (c Config) String() string {
	var s bytes.Buffer

	s.WriteString(pluginSectionString("input", c.Input))
	s.WriteString(pluginSectionString("filter", c.Filter))
	s.WriteString(pluginSectionString("output", c.Output))

	return s.String()
}

func pluginSectionString(pluginType string, ps []PluginSection) string {
	var s bytes.Buffer
	if len(ps) > 0 {
		s.WriteString(fmt.Sprintln(pluginType + " {"))
		var ss bytes.Buffer
		for _, p := range ps {
			ss.WriteString(fmt.Sprintf("%v", p))
		}
		s.WriteString(prefix(ss.String()))
		s.WriteString(fmt.Sprintln("}"))
	}
	return s.String()
}

const (
	Undefined = iota
	Input
	Filter
	Output
)

type PluginType int

func (pt PluginType) String() string {
	switch pt {
	case Input:
		return "input"
	case Filter:
		return "filter"
	case Output:
		return "output"
	default:
		return "undefined"
	}
}

type PluginSection struct {
	PluginType      PluginType
	BranchOrPlugins []BranchOrPlugin
}

func NewPluginSection(pt PluginType, bop ...BranchOrPlugin) PluginSection {
	return PluginSection{
		PluginType:      pt,
		BranchOrPlugins: bop,
	}
}

func NewPluginSections(pt PluginType, bop ...BranchOrPlugin) []PluginSection {
	return []PluginSection{
		NewPluginSection(pt, bop...),
	}
}

func (ps PluginSection) String() string {
	var s bytes.Buffer
	for _, bop := range ps.BranchOrPlugins {
		s.WriteString(fmt.Sprintf("%v", bop))
	}
	return s.String()
}

type BranchOrPlugin interface{}

type Plugin struct {
	name       string
	attributes []Attribute
}

func NewPlugin(name string, attributes ...Attribute) Plugin {
	return Plugin{
		name:       name,
		attributes: attributes,
	}
}

func (p Plugin) String() string {
	var s bytes.Buffer
	s.WriteString(fmt.Sprintln(p.name, "{"))
	if p.attributes != nil && len(p.attributes) > 0 {
		var ss bytes.Buffer
		for _, attr := range p.attributes {
			if attr != nil {
				ss.WriteString(fmt.Sprintln(attr.String()))
			}
		}
		s.WriteString(prefix(ss.String()))
	}
	s.WriteString(fmt.Sprintln("}"))
	return s.String()
}

type Attribute interface {
	String() string
	ValueString() string
}

type PluginAttribute struct {
	name  string
	value Plugin
}

func NewPluginAttribute(name string, value Plugin) PluginAttribute {
	return PluginAttribute{
		name:  name,
		value: value,
	}
}

func (pa PluginAttribute) String() string {
	return fmt.Sprintf("%s => %s", pa.name, pa.ValueString())
}

func (pa PluginAttribute) ValueString() string {
	return fmt.Sprintf("%s", pa.value.String())
}

const (
	// Undefined is already defined

	DoubleQuoted = iota + 1
	SingleQuoted
	Bareword
)

type StringAttributeType int

func (sat StringAttributeType) String() string {
	switch sat {
	case DoubleQuoted:
		return `"`
	case SingleQuoted:
		return `'`
	case Bareword:
		return ``
	default:
		return "undefined"
	}
}

type StringAttribute struct {
	name  string
	value string
	sat   StringAttributeType
}

func NewStringAttribute(name, value string, sat StringAttributeType) StringAttribute {
	return StringAttribute{
		name:  name,
		value: value,
		sat:   sat,
	}
}

func (sa StringAttribute) String() string {
	return fmt.Sprintf("%s => %s", sa.name, sa.ValueString())
}

func (sa StringAttribute) ValueString() string {
	return fmt.Sprintf("%s%s%s", sa.StringAttributeType(), sa.Value(), sa.StringAttributeType())
}

func (sa StringAttribute) Value() string {
	return sa.value
}

func (sa StringAttribute) StringAttributeType() StringAttributeType {
	return sa.sat
}

type NumberAttribute struct {
	name  string
	value float64
}

func NewNumberAttribute(name string, value float64) NumberAttribute {
	return NumberAttribute{
		name:  name,
		value: value,
	}
}

func (na NumberAttribute) String() string {
	return fmt.Sprintf("%s => %s", na.name, na.ValueString())
}

func (na NumberAttribute) ValueString() string {
	return fmt.Sprintf("%v", na.Value())
}

func (na NumberAttribute) Value() float64 {
	return na.value
}

type ArrayAttribute struct {
	name  string
	value []Attribute
}

func NewArrayAttribute(name string, value ...Attribute) ArrayAttribute {
	return ArrayAttribute{
		name:  name,
		value: value,
	}
}

func (aa ArrayAttribute) String() string {
	return fmt.Sprintf("%s => %s", aa.name, aa.ValueString())
}

func (aa ArrayAttribute) ValueString() string {
	var s bytes.Buffer
	s.WriteString(fmt.Sprintf("[ "))

	// TODO: use slice of string and string.Join
	first := true
	for _, a := range aa.Value() {
		if first {
			first = false
		} else {
			s.WriteString(", ")
		}
		s.WriteString(fmt.Sprintf("%s", a.ValueString()))
	}
	s.WriteString(fmt.Sprintf(" ]"))
	return s.String()
}

func (aa ArrayAttribute) Value() []Attribute {
	return aa.value
}

type HashAttribute struct {
	name  string
	value []HashEntry
}

func NewHashAttribute(name string, value ...HashEntry) HashAttribute {
	return HashAttribute{
		name:  name,
		value: value,
	}
}

func (ha HashAttribute) String() string {
	return fmt.Sprintf("%s => %s", ha.name, ha.ValueString())
}

func (ha HashAttribute) ValueString() string {
	var s bytes.Buffer
	s.WriteString(fmt.Sprintln("{"))
	if len(ha.value) > 0 {
		var ss bytes.Buffer
		for _, v := range ha.Value() {
			ss.WriteString(fmt.Sprintln(v.String()))
		}
		s.WriteString(prefix(ss.String()))
	}
	s.WriteString(fmt.Sprint("}"))
	return s.String()
}

func (ha HashAttribute) Value() []HashEntry {
	return ha.value
}

type HashEntry struct {
	name  string
	value Attribute
}

func NewHashEntry(name string, value Attribute) HashEntry {
	return HashEntry{
		name:  name,
		value: value,
	}
}

func (he HashEntry) String() string {
	return fmt.Sprintf("%s => %s", he.name, he.ValueString())
}

func (he HashEntry) ValueString() string {
	return he.value.ValueString()
}

func (he HashEntry) Value() Attribute {
	return he.value
}

type Branch struct {
	ifBlock     IfBlock
	elseIfBlock []ElseIfBlock
	elseBlock   ElseBlock
}

// Arguments for elseBlock and elseIfBlock are in the wrong order from logically point of view.
// This is due to the variadic nature of the elseIfBlock argument.
func NewBranch(ifBlock IfBlock, elseBlock ElseBlock, elseIfBlock ...ElseIfBlock) Branch {
	return Branch{
		ifBlock:     ifBlock,
		elseIfBlock: elseIfBlock,
		elseBlock:   elseBlock,
	}
}

// TODO: Maybe we should add helper functions NewIfBranch and NewIfElseBranch

func (b Branch) String() string {
	var s bytes.Buffer
	s.WriteString(fmt.Sprint(b.ifBlock))
	if b.elseIfBlock != nil && len(b.elseIfBlock) > 0 {
		for _, block := range b.elseIfBlock {
			s.WriteString(fmt.Sprint(block))
		}
	}
	s.WriteString(fmt.Sprintln(b.elseBlock))
	return s.String()
}

type IfBlock struct {
	condition Condition
	block     []BranchOrPlugin
}

func NewIfBlock(condition Condition, block ...BranchOrPlugin) IfBlock {
	return IfBlock{
		condition: condition,
		block:     block,
	}
}

func (ib IfBlock) String() string {
	var s bytes.Buffer
	s.WriteString(fmt.Sprintf("if %v {\n", ib.condition))
	if ib.block != nil && len(ib.block) > 0 {
		var ss bytes.Buffer
		for _, block := range ib.block {
			if block != nil {
				ss.WriteString(fmt.Sprint(block))
			}
		}
		s.WriteString(prefix(ss.String()))
	}
	s.WriteString(fmt.Sprint("}"))
	return s.String()
}

type ElseIfBlock struct {
	condition Condition
	block     []BranchOrPlugin
}

func NewElseIfBlock(condition Condition, block ...BranchOrPlugin) ElseIfBlock {
	return ElseIfBlock{
		condition: condition,
		block:     block,
	}
}

func (eib ElseIfBlock) String() string {
	var s bytes.Buffer
	s.WriteString(fmt.Sprintf(" else if %v {\n", eib.condition))
	if eib.block != nil && len(eib.block) > 0 {
		var ss bytes.Buffer
		for _, block := range eib.block {
			if block != nil {
				ss.WriteString(fmt.Sprint(block))
			}
		}
		s.WriteString(prefix(ss.String()))
	}
	s.WriteString(fmt.Sprint("}"))
	return s.String()
}

type ElseBlock struct {
	block []BranchOrPlugin
}

func NewElseBlock(block ...BranchOrPlugin) ElseBlock {
	return ElseBlock{
		block: block,
	}
}

func (eb ElseBlock) String() string {
	if eb.block == nil || len(eb.block) == 0 {
		return ""
	}

	var s bytes.Buffer
	s.WriteString(fmt.Sprintln(" else {"))
	var ss bytes.Buffer
	for _, block := range eb.block {
		if block != nil {
			ss.WriteString(fmt.Sprint(block))
		}
	}
	s.WriteString(prefix(ss.String()))
	s.WriteString(fmt.Sprintln("}"))
	return s.String()
}

type Condition struct {
	expression []Expression
}

func NewCondition(expression ...Expression) Condition {
	return Condition{
		expression: expression,
	}
}

func (c Condition) String() string {
	var s bytes.Buffer
	for _, expression := range c.expression {
		if expression != nil {
			s.WriteString(fmt.Sprint(expression))
		}
	}
	return s.String()
}

type Expression interface {
	BoolOperator() BooleanOperator
	SetBoolOperator(BooleanOperator)
}

type BoolExpression struct {
	boolOperator BooleanOperator
}

func (be *BoolExpression) BoolOperator() BooleanOperator {
	return be.boolOperator
}

func (be *BoolExpression) SetBoolOperator(bo BooleanOperator) {
	be.boolOperator = bo
}

func (be BoolExpression) String() string {
	return be.boolOperator.String()
}

type ConditionExpression struct {
	*BoolExpression
	condition Condition
}

func NewConditionExpression(boolOperator BooleanOperator, condition Condition) ConditionExpression {
	return ConditionExpression{
		BoolExpression: &BoolExpression{
			boolOperator: boolOperator,
		},
		condition: condition,
	}
}

func (ce ConditionExpression) String() string {
	return fmt.Sprintf("%v(%s)", ce.BoolExpression, ce.condition.String())
}

// type NegativeExpression interface {
// 	String() string
// }

type NegativeCondition struct {
	*BoolExpression
	condition Condition
}

func NewNegativeCondition(boolOperator BooleanOperator, condition Condition) NegativeCondition {
	return NegativeCondition{
		BoolExpression: &BoolExpression{
			boolOperator: boolOperator,
		},
		condition: condition,
	}
}

func (nc NegativeCondition) String() string {
	return fmt.Sprintf("%v! (%s)", nc.BoolExpression, nc.condition.String())
}

type NegativeSelector struct {
	*BoolExpression
	selector Selector
}

func NewNegativeSelector(boolOperator BooleanOperator, selector Selector) NegativeSelector {
	return NegativeSelector{
		BoolExpression: &BoolExpression{
			boolOperator: boolOperator,
		},
		selector: selector,
	}
}

func (ns NegativeSelector) String() string {
	return fmt.Sprintf("%v! %s", ns.BoolExpression, ns.selector)
}

type InExpression struct {
	*BoolExpression
	lvalue Rvalue
	rvalue Rvalue
}

func NewInExpression(boolOperator BooleanOperator, lvalue Rvalue, rvalue Rvalue) InExpression {
	return InExpression{
		BoolExpression: &BoolExpression{
			boolOperator: boolOperator,
		},
		lvalue: lvalue,
		rvalue: rvalue,
	}
}

func (ie InExpression) String() string {
	return fmt.Sprintf("%v%v in %v", ie.BoolExpression, ie.lvalue.ValueString(), ie.rvalue.ValueString())
}

type NotInExpression struct {
	*BoolExpression
	rvalue Rvalue
	lvalue Rvalue
}

func NewNotInExpression(boolOperator BooleanOperator, lvalue Rvalue, rvalue Rvalue) NotInExpression {
	return NotInExpression{
		BoolExpression: &BoolExpression{
			boolOperator: boolOperator,
		},
		lvalue: lvalue,
		rvalue: rvalue,
	}
}

func (nie NotInExpression) String() string {
	return fmt.Sprintf("%v%v not in %v", nie.BoolExpression, nie.lvalue.ValueString(), nie.rvalue.ValueString())
}

type Rvalue interface {
	String() string
	ValueString() string
}

type RvalueExpression struct {
	*BoolExpression
	rvalue Rvalue
}

func NewRvalueExpression(boolOperator BooleanOperator, rvalue Rvalue) RvalueExpression {
	return RvalueExpression{
		BoolExpression: &BoolExpression{
			boolOperator: boolOperator,
		},
		rvalue: rvalue,
	}
}

func (re RvalueExpression) String() string {
	return fmt.Sprintf("%v%v", re.BoolExpression, re.rvalue.ValueString())
}

type CompareExpression struct {
	*BoolExpression
	lvalue          Rvalue
	compareOperator CompareOperator
	rvalue          Rvalue
}

func NewCompareExpression(boolOperator BooleanOperator, lvalue Rvalue, compareOperator CompareOperator, rvalue Rvalue) CompareExpression {
	return CompareExpression{
		BoolExpression: &BoolExpression{
			boolOperator: boolOperator,
		},
		lvalue:          lvalue,
		compareOperator: compareOperator,
		rvalue:          rvalue,
	}
}

func (ce CompareExpression) String() string {
	return fmt.Sprintf("%v%v %v %v", ce.BoolExpression, ce.lvalue.ValueString(), ce.compareOperator, ce.rvalue.ValueString())
}

type CompareOperator int

const (
	Equal = iota + 1
	NotEqual
	LessOrEqual
	GreaterOrEqual
	LessThan
	GreaterThan
)

func (co CompareOperator) String() string {
	switch co {
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
		return "undefined"
	}
}

type RegexpExpression struct {
	*BoolExpression
	lvalue         Rvalue
	regexpOperator RegexpOperator
	rvalue         StringOrRegexp
}

func NewRegexpExpression(boolOperator BooleanOperator, lvalue Rvalue, regexpOperator RegexpOperator, rvalue StringOrRegexp) RegexpExpression {
	return RegexpExpression{
		BoolExpression: &BoolExpression{
			boolOperator: boolOperator,
		},
		lvalue:         lvalue,
		regexpOperator: regexpOperator,
		rvalue:         rvalue,
	}
}

func (re RegexpExpression) String() string {
	return fmt.Sprintf("%v%v %v %v", re.BoolExpression, re.lvalue.ValueString(), re.regexpOperator, re.rvalue.ValueString())
}

type StringOrRegexp interface {
	String() string
	ValueString() string
}

type RegexpOperator int

const (
	RegexpMatch = iota + 1
	RegexpNotMatch
)

func (ro RegexpOperator) String() string {
	switch ro {
	case RegexpMatch:
		return fmt.Sprint("=~")
	case RegexpNotMatch:
		return fmt.Sprint("!~")
	default:
		return fmt.Sprint(" undefined ")
	}
}

type Regexp struct {
	regexp string
}

func NewRegexp(regexp string) Regexp {
	return Regexp{
		regexp: regexp,
	}
}

func (r Regexp) String() string {
	return fmt.Sprintf("/%s/", r.regexp)
}

func (r Regexp) ValueString() string {
	return r.String()
}

type BooleanOperator int

const (
	NoOperator = iota + 1
	And
	Or
	Xor
	Nand
)

func (be BooleanOperator) String() string {
	switch be {
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
		return " undefined "
	}
}

type SelectorExpression struct {
	*BoolExpression
	selectorElement []SelectorElement
}

type Selector struct {
	elements []SelectorElement
}

func NewSelector(elements []SelectorElement) Selector {
	return Selector{
		elements: elements,
	}
}

func NewSelectorFromNames(names ...string) Selector {
	var elements []SelectorElement
	for _, name := range names {
		elements = append(elements, NewSelectorElement(name))
	}
	return NewSelector(elements)
}

func (s Selector) String() string {
	var bb bytes.Buffer
	for _, element := range s.elements {
		bb.WriteString(element.String())
	}
	return bb.String()
}

func (s Selector) ValueString() string {
	return s.String()
}

type SelectorElement struct {
	name string
}

func NewSelectorElement(name string) SelectorElement {
	return SelectorElement{
		name: name,
	}
}

func (se SelectorElement) String() string {
	return fmt.Sprintf("[%s]", se.name)
}
