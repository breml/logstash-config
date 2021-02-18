package astutil

import (
	"fmt"

	"github.com/breml/logstash-config/ast"
)

// An ApplyPluginsFunc is invoked by ApplyPlugins for each plugin p, even if p
// is nil, using a Cursor describing the current plugin and providing operations
// on it.
type ApplyPluginsFunc func(cursor *Cursor)

// ApplyPlugins traverses an AST recursively, starting with root, and calling
// applyPluginsFunc for each plugin. Apply returns the AST, possibly modified.
func ApplyPlugins(root []ast.BranchOrPlugin, applyPluginsFunc ApplyPluginsFunc) (result []ast.BranchOrPlugin) {
	c := Cursor{
		parent: root,
		iter: iterator{
			index: 0,
			step:  1,
		},
	}

	for {
		if c.iter.index >= len(c.parent) {
			break
		}

		c.iter.step = 1

		switch block := c.parent[c.iter.index].(type) {
		case ast.Branch:
			block.IfBlock.Block = ApplyPlugins(block.IfBlock.Block, applyPluginsFunc)

			for i := range block.ElseIfBlock {
				block.ElseIfBlock[i].Block = ApplyPlugins(block.ElseIfBlock[i].Block, applyPluginsFunc)
			}

			block.ElseBlock.Block = ApplyPlugins(block.ElseBlock.Block, applyPluginsFunc)

			c.parent[c.iter.index] = block

		case ast.Plugin:
			applyPluginsFunc(&c)

		case nil:
			applyPluginsFunc(&c)

		default:
			panic(fmt.Sprintf("type %T for block in ApplyPlugins not supported", c.parent[c.iter.index]))
		}

		c.iter.index += c.iter.step
	}

	return c.parent
}

// An iterator controls iteration over a slice of nodes.
type iterator struct {
	index, step int
}

// A Cursor describes a plugin encountered during ApplyPlugins.
// Information about the node and its parent is available from the Plugin,
// Parent and Index methods.
//
// The methods Replace and Delete can be used to change the AST without
// disrupting ApplyPlugins.
type Cursor struct {
	parent []ast.BranchOrPlugin
	iter   iterator // valid if non-nil
}

// Plugin returns the current Plugin.
func (c *Cursor) Plugin() *ast.Plugin {
	p, ok := c.parent[c.iter.index].(ast.Plugin)
	if !ok {
		var p *ast.Plugin
		return p
	}
	return &p
}

// Parent returns the slice of BranchOrPlugin of the current Plugin.
func (c *Cursor) Parent() []ast.BranchOrPlugin { return c.parent }

// Index reports the index of the current Plugin in the slice of BranchOrPlugin
// that contains it.
func (c *Cursor) Index() int {
	return c.iter.index
}

// Delete deletes the current Plugin from its containing slice.
// If the current Plugin is not part of a slice, Delete panics.
func (c *Cursor) Delete() {
	i := c.Index()
	if i < 0 {
		panic("Delete plugin not contained in slice")
	}
	c.parent = append(c.parent[:i], c.parent[i+1:]...)
	c.iter.step--
}

// Replace replaces the current Plugin with p.
// The replacement is not walked by Apply.
// If the current Plugin is not part of a slice, Replace panics.
func (c *Cursor) Replace(p ast.BranchOrPlugin) {
	i := c.Index()
	if i < 0 {
		panic("Replaced plugin not contained in slice")
	}
	c.parent[i] = p
}

// InsertBefore inserts p before the current Plugin in its containing slice.
// If the current Node is not part of a slice, InsertBefore panics.
// Apply will not walk p.
func (c *Cursor) InsertBefore(p ast.BranchOrPlugin) {
	i := c.Index()
	if i < 0 {
		panic("InsertBefore plugin not contained in slice")
	}
	c.parent = append(c.parent[:i], append([]ast.BranchOrPlugin{p}, c.parent[i:]...)...)
	c.iter.step++
}

// InsertAfter inserts p after the current Plugin in its containing slice.
// If the current Node is not part of a slice, InsertAfter panics.
// Apply will not walk p.
func (c *Cursor) InsertAfter(p ast.BranchOrPlugin) {
	i := c.Index()
	if i < 0 {
		panic("InsertAfter plugin not contained in slice")
	}
	c.parent = append(c.parent[:i+1], append([]ast.BranchOrPlugin{p}, c.parent[i+1:]...)...)
	c.iter.step++
}
