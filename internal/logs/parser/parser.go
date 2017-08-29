//go:generate peg -inline -switch grammar.peg

// Package parser provides a parser for Up's
// log query language, abstracting away provider
// specifics.
package parser

import (
	"strconv"

	"github.com/apex/up/internal/logs/parser/ast"
)

// Parse query string.
func Parse(s string) (ast.Node, error) {
	p := &parser{Buffer: s}
	p.Init()

	if err := p.Parse(); err != nil {
		return nil, err
	}

	p.Execute()
	n := ast.Root{Node: p.stack[0]}
	return n, nil
}

// push node.
func (p *parser) push(n ast.Node) {
	p.stack = append(p.stack, n)
}

// pop node.
func (p *parser) pop() ast.Node {
	if len(p.stack) == 0 {
		panic("pop: no nodes")
	}

	n := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return n
}

// AddLevel node.
func (p *parser) AddLevel(s string) {
	p.AddField("level")
	p.AddString(s)
	p.AddBinary(ast.EQ)
	p.AddExpr()
}

// AddExpr node.
func (p *parser) AddExpr() {
	p.push(ast.Expr{
		Node: p.pop(),
	})
}

// AddField node.
func (p *parser) AddField(s string) {
	switch s {
	case "level", "message", "timestamp":
		p.push(ast.Property(s))
	default:
		p.push(ast.Field(s))
	}
}

// AddString node.
func (p *parser) AddString(s string) {
	p.push(ast.String(s))
}

// AddSubscript node.
func (p *parser) AddSubscript(s string) {
	p.push(ast.Subscript{
		Left:  p.pop(),
		Right: ast.Literal(s),
	})
}

// AddMember node.
func (p *parser) AddMember(s string) {
	p.push(ast.Member{
		Left:  p.pop(),
		Right: ast.Literal(s),
	})
}

// SetNumber text.
func (p *parser) SetNumber(s string) {
	p.number = s
}

// AddNumber node.
func (p *parser) AddNumber(unit string) {
	f, _ := strconv.ParseFloat(p.number, 64)
	p.push(ast.Number{
		Value: f,
		Unit:  unit,
	})
}

// AddTuple node.
func (p *parser) AddTuple() {
	p.push(ast.Tuple{})
}

// AddTupleValue node.
func (p *parser) AddTupleValue() {
	v := p.pop()
	t := p.pop().(ast.Tuple)
	t = append(t, v)
	p.push(t)
}

// AddBinary node.
func (p *parser) AddBinary(op ast.Op) {
	p.push(ast.Binary{
		Op:    op,
		Right: p.pop(),
		Left:  p.pop(),
	})
}

// AddBinaryContains node.
func (p *parser) AddBinaryContains() {
	p.push(ast.Binary{
		Op:    ast.EQ,
		Right: ast.Contains{Node: p.pop()},
		Left:  p.pop(),
	})
}

// AddUnary node.
func (p *parser) AddUnary(op ast.Op) {
	p.push(ast.Unary{
		Op:    op,
		Right: p.pop(),
	})
}
