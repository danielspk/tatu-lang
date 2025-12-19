// Package parser converts tokens into an Abstract Syntax Tree with syntax validation and sugar expansion.
package parser

import (
	"fmt"

	"github.com/danielspk/tatu-lang/pkg/ast"
	"github.com/danielspk/tatu-lang/pkg/debug"
	"github.com/danielspk/tatu-lang/pkg/location"
	"github.com/danielspk/tatu-lang/pkg/token"
)

// Parser is responsible for analyzing tokens and generating the abstract syntax tree (AST).
type Parser struct {
	tokens  []token.Token
	current int

	sugar    SyntaxSugar
	analyzer SyntaxAnalyzer
}

// NewParser builds a new Parser.
func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		tokens:   tokens,
		sugar:    SyntaxSugar{},
		analyzer: SyntaxAnalyzer{},
	}
}

// Parse parses the tokens and generates a resulting AST.
func (p *Parser) Parse() (*ast.AST, error) {
	if len(p.tokens) == 0 {
		return nil, fmt.Errorf("no tokens found")
	}

	prog, err := p.parseProgram()
	if err != nil {
		return nil, err
	}

	return &ast.AST{
		Program: prog,
	}, nil
}

// advance returns the current token and advances one position.
func (p *Parser) advance() token.Token {
	curToken := p.tokens[p.current]
	p.current++

	return curToken
}

// peek returns the current token.
func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

// previous returns the previous token.
func (p *Parser) previous() token.Token {
	if p.current == 0 {
		return p.tokens[len(p.tokens)-1] // it's assumed to be EOF
	}

	return p.tokens[p.current-1]
}

// match checks if the current token is of a type. If so, advances one position.
func (p *Parser) match(tokenType token.Type) bool {
	if p.isEOF() || p.peek().Type != tokenType {
		return false
	}

	_ = p.advance()

	return true
}

// isEOF checks if the current token is EOF.
func (p *Parser) isEOF() bool {
	return p.peek().Type == token.EOF
}

// parseAtom parses an atom.
//
// <atom> ::= <number> | <string> | <boolean> | <symbol> | "nil"
func (p *Parser) parseAtom() (ast.SExpr, error) {
	atom := p.peek()
	_ = p.advance()

	loc := location.NewLocation(
		atom.File,
		location.NewPosition(atom.Start.Line, atom.Start.Column, atom.Start.Offset),
		location.NewPosition(atom.End.Line, atom.End.Column, atom.Start.Offset),
	)

	switch atom.Type {
	case token.Number:
		return ast.NewNumberExpr(atom.Literal.(float64), loc), nil
	case token.String:
		return ast.NewStringExpr(atom.Literal.(string), loc), nil
	case token.Bool:
		return ast.NewBoolExpr(atom.Literal.(bool), loc), nil
	case token.Symbol:
		return ast.NewSymbolExpr(atom.Literal.(string), loc), nil
	case token.Nil:
		return ast.NewNilExpr(loc), nil
	default:
		return nil, p.error(fmt.Sprintf("unknown atom %d", atom.Type), atom.Location)
	}
}

// parseList parses a list.
//
// <list> ::= "(" <expr>* ")"
func (p *Parser) parseList() (ast.SExpr, error) {
	if !p.match(token.LeftParen) {
		return nil, p.error("expected list", p.peek().Location)
	}

	var exprs []ast.SExpr
	var startLoc, endLoc location.Location

	startLoc = p.previous().Location

	closingParen := false

	for !p.isEOF() {
		if p.match(token.RightParen) {
			endLoc = p.previous().Location
			closingParen = true

			break
		}

		exp, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		exprs = append(exprs, exp)
	}

	if !closingParen {
		return nil, p.error("unclosed parenthesis", startLoc)
	}

	var listExpr ast.SExpr = ast.NewListExpr(exprs, location.NewLocation(
		startLoc.File,
		location.NewPosition(startLoc.Start.Line, startLoc.Start.Column, startLoc.Start.Offset),
		location.NewPosition(endLoc.End.Line, endLoc.End.Column, endLoc.End.Offset),
	))

	if err := p.sugar.Transform(&listExpr); err != nil {
		return nil, err
	}

	if err := p.analyzer.Validate(listExpr); err != nil {
		return nil, err
	}

	return listExpr, nil
}

// parseExpression parses an expression.
//
// <expr> ::= <atom> | <list>
func (p *Parser) parseExpression() (ast.SExpr, error) {
	expr := p.peek()

	if expr.Type == token.Number || expr.Type == token.String || expr.Type == token.Bool || expr.Type == token.Symbol || expr.Type == token.Nil {
		return p.parseAtom()
	} else if expr.Type == token.LeftParen {
		return p.parseList()
	}

	return nil, p.error("expected expression", expr.Location)
}

// parseProgram parses a program.
//
// <program> ::= (<exp>)*
func (p *Parser) parseProgram() ([]ast.SExpr, error) {
	var prog []ast.SExpr

	for !p.isEOF() {
		exp, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		prog = append(prog, exp)
	}

	return prog, nil
}

// error makes an error.
func (p *Parser) error(msg string, loc location.Location) *debug.Error {
	return &debug.Error{
		Msg:    msg,
		Line:   loc.End.Line,
		Column: loc.End.Column,
		File:   loc.File,
	}
}
