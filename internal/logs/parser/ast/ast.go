// Package ast provides the log query language abstract syntax tree.
package ast

import (
	"fmt"
	"strconv"
	"strings"
)

// Op type.
type Op string

// Op types.
const (
	LNOT Op = "not"
	NOT     = "!"
	IN      = "in"
	OR      = "||"
	AND     = "&&"
	NE      = "!="
	EQ      = "="
	GT      = ">"
	LT      = "<"
	GE      = ">="
	LE      = "<="
)

// Node interface.
type Node interface {
	String() string
}

// Root node.
type Root struct {
	Node Node
}

// String implementation.
func (n Root) String() string {
	return fmt.Sprintf(`{ %s }`, n.Node)
}

// Expr node.
type Expr struct {
	Node Node
}

// String implementation.
func (n Expr) String() string {
	return fmt.Sprintf(`(%s)`, n.Node)
}

// Literal node.
type Literal string

// String implementation.
func (n Literal) String() string {
	return fmt.Sprintf(`%s`, string(n))
}

// Tuple node.
type Tuple []Node

// String implementation.
func (n Tuple) String() string {
	return fmt.Sprintf(`%#v`, n)
}

// Contains node.
type Contains struct {
	Node Node
}

// String implementation.
func (n Contains) String() string {
	switch v := n.Node.(type) {
	case String:
		return fmt.Sprintf(`"*%s*"`, string(v))
	default:
		return fmt.Sprintf(`%s`, n.Node)
	}
}

// String node.
type String string

// String implementation.
func (n String) String() string {
	return fmt.Sprintf(`%q`, string(n))
}

// Property node.
type Property string

// String implementation.
func (n Property) String() string {
	return fmt.Sprintf(`$.%s`, string(n))
}

// Field node.
type Field string

// String implementation.
func (n Field) String() string {
	return fmt.Sprintf(`$.fields.%s`, string(n))
}

// Subscript node.
type Subscript struct {
	Left  Node
	Right Node
}

// String implementation.
func (n Subscript) String() string {
	return fmt.Sprintf(`%s[%s]`, n.Left, n.Right)
}

// Member node.
type Member struct {
	Left  Node
	Right Node
}

// String implementation.
func (n Member) String() string {
	return fmt.Sprintf(`%s.%s`, n.Left, n.Right)
}

// Number node.
type Number struct {
	Value float64
	Unit  string
}

// String implementation.
func (n Number) String() string {
	v := n.Value

	switch n.Unit {
	case "kb":
		v *= 1 << 10
	case "mb":
		v *= 1 << 20
	case "gb":
		v *= 1 << 30
	case "s":
		v *= 1000
	}

	return strconv.FormatFloat(v, 'f', -1, 64)
}

// Binary node.
type Binary struct {
	Op    Op
	Left  Node
	Right Node
}

// String implementation.
func (n Binary) String() string {
	switch n.Op {
	case IN:
		var s []string
		for _, v := range n.Right.(Tuple) {
			s = append(s, fmt.Sprintf(`%s %s %s`, n.Left, EQ, v))
		}
		return fmt.Sprintf(`(%s)`, strings.Join(s, " || "))
	default:
		return fmt.Sprintf(`%s %s %s`, n.Left, n.Op, n.Right)
	}
}

// Unary node.
type Unary struct {
	Op    Op
	Right Node
}

// String implementation.
func (n Unary) String() string {
	switch n.Op {
	case LNOT:
		return fmt.Sprintf(`!(%s)`, n.Right)
	default:
		return fmt.Sprintf(`%s%s`, n.Op, n.Right)
	}
}
