package config

import (
	"fmt"
	"strconv"

	"github.com/breml/logstash-config/ast"
)

func initParser() (bool, error) {
	farthestFailure = []errPos{}
	return true, nil
}

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
				return ast.Config{}, fmt.Errorf("PluginType is not supported: %#v", ps)
			}
		} else {
			return ast.Config{}, fmt.Errorf("Value is not a PluginSection: %#v", ips1)
		}
	}

	return ast.Config{
		Input:  input,
		Filter: filter,
		Output: output,
	}, nil
}

func pluginSection(pt1, bops1 interface{}) (ast.PluginSection, error) {
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

func branch(ifBlock, elseIfBlocks1, elseBlock1 interface{}) (ast.Branch, error) {
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

	return ast.NewBranch(ifBlock.(ast.IfBlock), elseBlock, elseIfBlocks...), nil
}

func ifBlock(cond, bops interface{}) (ast.IfBlock, error) {
	return ast.NewIfBlock(cond.(ast.Condition), branchOrPlugins(bops)...), nil
}

func elseIfBlock(cond, bops interface{}) (ast.ElseIfBlock, error) {
	return ast.NewElseIfBlock(cond.(ast.Condition), branchOrPlugins(bops)...), nil
}

func elseBlock(bops interface{}) (ast.ElseBlock, error) {
	return ast.NewElseBlock(branchOrPlugins(bops)...), nil
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
