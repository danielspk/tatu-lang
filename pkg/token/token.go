// Package token defines token types for the lexical analyzer.
package token

import (
	"github.com/danielspk/tatu-lang/pkg/location"
)

// Type represents a token type.
type Type uint8

// Token types.
const (
	LeftParen  Type = iota + 1 // (
	RightParen                 // )
	Number                     // 0-9.
	String                     // "..."
	Bool                       // "true" | "false"
	Nil                        // "nil"
	Symbol                     // alphanumeric | operators
	EOF
)

// Token represents a token extracted from the Lexer tokenization.
type Token struct {
	Type    Type
	Lexeme  string
	Literal any
	location.Location
}

// NewToken builds a new Token.
func NewToken(tokenType Type, lexeme string, literal any, location location.Location) Token {
	return Token{
		Type:     tokenType,
		Lexeme:   lexeme,
		Literal:  literal,
		Location: location,
	}
}
