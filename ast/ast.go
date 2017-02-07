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
	//Name() string
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

type Branch struct{}

func NewBranch() Branch {
	return Branch{}
}

func (b Branch) String() string {
	return ""
}
