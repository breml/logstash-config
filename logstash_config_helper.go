package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/breml/logstash-config/ast"
)

func ret(el interface{}) (interface{}, error) {
	return el, nil
}

func config(ps1, pss1 interface{}) (ast.Config, error) {
	var (
		input  []ast.PluginSection
		filter []ast.PluginSection
		output []ast.PluginSection
	)

	ips := toIfaceSlice(ps1)
	ips = append(ips, toIfaceSlice(pss1)...)

	for _, ips1 := range ips {
		if ps, ok := ips1.(ast.PluginSection); ok {
			switch ps.PluginType {
			case ast.Input:
				input = append(input, ps)
			case ast.Filter:
				filter = append(filter, ps)
			case ast.Output:
				output = append(output, ps)
			default:
				// TODO: Return error
				fmt.Printf("why are we in default? %#v\n", ps)
			}
		} else {
			// TODO: Return error
			fmt.Printf("why dont we get an PluginSection: %#v\n", ips1)
		}
	}

	return ast.Config{
		Input:  input,
		Filter: filter,
		Output: output,
	}, nil
}

func pluginSection(pt1, bops1 ast.BranchOrPlugin) (ast.PluginSection, error) {
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
	}, nil
}

// func plugins(p1 interface{}, ps1 interface{}) ([]ast.Plugin, error) {
// 	var plugins []ast.Plugin

// 	ips := toIfaceSlice(p1)
// 	ips = append(ips, toIfaceSlice(ps1)...)

// 	for _, p := range ips {
// 		if p, ok := p.(ast.Plugin); ok {
// 			plugins = append(plugins, p)
// 		} else {
// 			return nil, fmt.Errorf("Argument is not a plugin")
// 		}
// 	}

// 	return plugins, nil
// }

func plugin(name, attributes interface{}) (ast.Plugin, error) {
	if attributes != nil {
		return ast.NewPlugin(name.(string), attributes.([]ast.Attribute)...), nil
	}
	return ast.NewPlugin(name.(string), nil), nil
}

func attributes(attribute, attributes1 interface{}) ([]ast.Attribute, error) {
	iattributes := toIfaceSlice(attribute)
	iattributes = append(iattributes, toIfaceSlice(attributes1)...)

	var attributes []ast.Attribute

	for _, attr := range iattributes {
		if attr, ok := attr.(ast.Attribute); ok {
			attributes = append(attributes, attr)
		} else {
			return nil, fmt.Errorf("Argument is not an attribute")
		}

	}

	return attributes, nil
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
		return ast.NewArrayAttribute(key.ValueString(), value.Value()...), nil
	case ast.HashAttribute:
		return ast.NewHashAttribute(key.ValueString(), value.Value()...), nil
	case ast.Plugin:
		return ast.NewPluginAttribute(key.ValueString(), value), nil
	default:
		return nil, fmt.Errorf("Type of value %#v with name %s is not supported", value, key.ValueString())
	}
}

func number(value string) (ast.NumberAttribute, error) {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		// TODO: is this possible to happen? are all values, which are valid floats in Logstash/Ruby also valid floats in Go?
		return ast.NumberAttribute{}, err
	}
	return ast.NewNumberAttribute("", f), nil
}

func array(attributes interface{}) (ast.ArrayAttribute, error) {
	if attributes != nil {
		return ast.NewArrayAttribute("", attributes.([]ast.Attribute)...), nil
	}
	// TODO: Is this an error?
	return ast.NewArrayAttribute("", nil), nil
}

func hash(attributes interface{}) (ast.HashAttribute, error) {
	if attributes != nil {
		return ast.NewHashAttribute("", attributes.([]ast.HashEntry)...), nil
	}
	// TODO: Is this an error?
	return ast.HashAttribute{}, nil
}

func hashentries(attribute, attributes1 interface{}) ([]ast.HashEntry, error) {
	// TODO: is this function generalizable with attributes?
	iattributes := toIfaceSlice(attribute)
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

func hashentry(name, value interface{}) (ast.HashEntry, error) {
	var key ast.StringAttribute

	switch name := name.(type) {
	case ast.StringAttribute:
		key = name
	}

	return ast.NewHashEntry(key.ValueString(), value.(ast.Attribute)), nil
}

func quotedvalue(c *current, quotechar string) (interface{}, error) {
	return strings.Trim(string(c.text), quotechar), nil
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
