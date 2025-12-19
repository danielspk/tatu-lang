package builder

import (
	"github.com/danielspk/tatu-lang/pkg/ast"
	"github.com/danielspk/tatu-lang/pkg/parser"
	"github.com/danielspk/tatu-lang/pkg/scanner"
	"github.com/danielspk/tatu-lang/pkg/token"
)

// defaultScanner default implementations of scanner.Scanner.
type defaultScanner struct{}

// NewDefaultScanner build a new default Scanner.
func NewDefaultScanner() Scanner {
	return &defaultScanner{}
}

// Scan builds a new scanner and scan the source code.
func (d *defaultScanner) Scan(source []byte, filename string) ([]token.Token, error) {
	return scanner.NewScanner(source, filename).Scan()
}

// defaultParser default implementations of parser.Parser.
type defaultParser struct{}

// NewDefaultParser build a new default Parser.
func NewDefaultParser() Parser {
	return &defaultParser{}
}

// Parse builds a new parser and parse the tokens.
func (d *defaultParser) Parse(tokens []token.Token) (*ast.AST, error) {
	return parser.NewParser(tokens).Parse()
}
