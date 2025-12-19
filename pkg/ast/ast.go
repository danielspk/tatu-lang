// Package ast provides the Abstract Syntax Tree representation.
package ast

import (
	"github.com/danielspk/tatu-lang/pkg/location"
)

// ExprKind represents the kind of an AST expression node.
type ExprKind uint8

// Expression kinds.
const (
	NumberKind ExprKind = iota + 1
	StringKind
	BoolKind
	SymbolKind
	NilKind
	ListKind
)

// SExpr represents an S-expression interface.
type SExpr interface {
	Kind() ExprKind
	Location() location.Location

	exprNode() // private marker method
}

// Node represents a base node expression.
type node struct {
	kind     ExprKind
	location location.Location
}

// Kind returns the expression kind.
func (n *node) Kind() ExprKind {
	return n.kind
}

// Location returns the node location.
func (n *node) Location() location.Location {
	return n.location
}

// exprNode expression marker method.
func (n *node) exprNode() {}

// NumberExpr represents a number atom expression.
type NumberExpr struct {
	node
	Number float64
}

// NewNumberExpr builds a new NumberExpr.
func NewNumberExpr(value float64, loc location.Location) *NumberExpr {
	return &NumberExpr{
		Number: value,
		node:   node{NumberKind, loc},
	}
}

// StringExpr represents a string atom expression.
type StringExpr struct {
	node
	String string
}

// NewStringExpr builds a new StringExpr.
func NewStringExpr(value string, loc location.Location) *StringExpr {
	return &StringExpr{
		String: value,
		node:   node{StringKind, loc},
	}
}

// BoolExpr represents a boolean atom expression.
type BoolExpr struct {
	node
	Bool bool
}

// NewBoolExpr builds a new BoolExpr.
func NewBoolExpr(value bool, loc location.Location) *BoolExpr {
	return &BoolExpr{
		Bool: value,
		node: node{BoolKind, loc},
	}
}

// SymbolExpr represents a symbol atom expression.
type SymbolExpr struct {
	node
	Symbol string
}

// NewSymbolExpr builds a new SymbolExpr.
func NewSymbolExpr(value string, loc location.Location) *SymbolExpr {
	return &SymbolExpr{
		Symbol: value,
		node:   node{SymbolKind, loc},
	}
}

// NilExpr represents a nil atom expression.
type NilExpr struct {
	node
}

// NewNilExpr builds a new NilExpr.
func NewNilExpr(loc location.Location) *NilExpr {
	return &NilExpr{
		node: node{NilKind, loc},
	}
}

// ListExpr represents a list expression.
type ListExpr struct {
	node
	List []SExpr
}

// NewListExpr builds a new ListExpr.
func NewListExpr(value []SExpr, loc location.Location) *ListExpr {
	return &ListExpr{
		List: value,
		node: node{ListKind, loc},
	}
}

// AST represents an ast for the program.
type AST struct {
	Program []SExpr
}
