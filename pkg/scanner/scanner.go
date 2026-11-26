// Package scanner performs lexical analysis, converting source code into tokens.
package scanner

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/danielspk/tatu-lang/pkg/debug"
	"github.com/danielspk/tatu-lang/pkg/location"
	"github.com/danielspk/tatu-lang/pkg/token"
)

// escapeReplacer string escape replacer.
var escapeReplacer = strings.NewReplacer(
	`\\`, `\`,
	`\n`, "\n",
	`\t`, "\t",
	`\r`, "\r",
	`\"`, `"`,
)

type cursor struct {
	offset uint
	line   uint
	column uint
}

// Scanner is the lexical analyzer that converts a script of rules into tokens.
type Scanner struct {
	source   string
	filename string
	start    cursor
	current  cursor
	tokens   []token.Token
}

// NewScanner builds a new Scanner.
func NewScanner() *Scanner {
	return &Scanner{}
}

// Scan tokenizes the source code to generate a slice of tokens.
func (s *Scanner) Scan(source []byte, filename string) ([]token.Token, error) {
	s.source = string(source)
	s.filename = filename
	s.start = cursor{offset: 0, line: 1, column: 1}
	s.current = cursor{offset: 0, line: 1, column: 1}
	s.tokens = make([]token.Token, 0)

	for !s.isAtEnd() {
		if err := s.scanToken(); err != nil {
			return nil, err
		}
	}

	_ = s.addToken(token.EOF)

	return s.tokens, nil
}

// scanToken scans the next token.
func (s *Scanner) scanToken() error {
	s.start = s.current

	chr := s.advance()

	switch chr {
	case ' ', '\t', '\r':
		return nil

	case '\n':
		s.current.line++
		s.current.column = 1
		return nil

	case ';':
		s.readComment()
		return nil

	case '(':
		return s.addToken(token.LeftParen)

	case ')':
		return s.addToken(token.RightParen)

	case '"':
		if err := s.readString(); err != nil {
			return err
		}
		return s.addToken(token.String)

	case '.':
		if s.peek() == '.' && s.lookAhead() == '.' {
			_ = s.advance()
			_ = s.advance()
			return s.addToken(token.Symbol)
		}
		return s.error("incomplete ellipsis")

	default:
		if s.isDigit(chr) || (chr == '-' && s.isDigit(s.peek())) {
			if err := s.readNumber(); err != nil {
				return err
			}
			return s.addToken(token.Number)
		}

		if s.isSymbol(chr) {
			s.readSymbol()
			return s.addToken(s.symbolToken())
		}

		return s.error(fmt.Sprintf("unexpected character `%c`", chr))
	}
}

// advance gets the current UTF8 character and advances to the next position.
func (s *Scanner) advance() rune {
	char, size := utf8.DecodeRuneInString(s.source[s.current.offset:])

	s.current.offset += uint(size)
	s.current.column++

	return char
}

// peek gets the current character.
func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}

	char, _ := utf8.DecodeRuneInString(s.source[s.current.offset:])

	return char
}

// lookAhead gets the next character.
func (s *Scanner) lookAhead() rune {
	if s.current.offset+1 >= uint(len(s.source)) {
		return 0
	}

	char, _ := utf8.DecodeRuneInString(s.source[s.current.offset+1:])

	return char
}

// isAtEnd checks if it is the end of a file.
func (s *Scanner) isAtEnd() bool {
	return s.current.offset >= uint(len(s.source))
}

// markNewline updates position on a newline character.
func (s *Scanner) markNewline(chr rune) {
	if chr == '\n' {
		s.current.line++
		s.current.column = 1
	}
}

// currentLexeme gets the current lexeme.
func (s *Scanner) currentLexeme() string {
	return s.source[s.start.offset:s.current.offset]
}

// currentLiteral gets the current literal value.
func (s *Scanner) currentLiteral(tokenType token.Type) (any, error) {
	lexeme := s.currentLexeme()

	switch tokenType {
	case token.Number:
		literal, err := strconv.ParseFloat(lexeme, 64)
		if err != nil {
			return nil, s.error(fmt.Sprintf("invalid value `%s` for a number: %s", lexeme, err.Error()))
		}

		return literal, nil

	case token.String:
		str := lexeme[1 : len(lexeme)-1]
		return s.processEscapes(str), nil

	case token.Symbol:
		return lexeme, nil

	case token.Bool:
		if lexeme == "true" {
			return true, nil
		}

		if lexeme == "false" {
			return false, nil
		}

		return nil, s.error(fmt.Sprintf("invalid value `%s` for a boolean", lexeme))

	default:
		return nil, nil
	}
}

// addToken adds a token to the list based on the current position.
func (s *Scanner) addToken(tokenType token.Type) error {
	lexeme := s.currentLexeme()

	if tokenType == token.EOF {
		lexeme = ""
		s.start = s.current
	}

	literal, err := s.currentLiteral(tokenType)
	if err != nil {
		return err
	}

	s.tokens = append(s.tokens, token.NewToken(
		tokenType,
		lexeme,
		literal,
		location.NewLocation(s.filename,
			location.NewPosition(
				s.start.line,
				s.start.column,
				s.start.offset),
			location.NewPosition(
				s.current.line,
				s.current.column,
				s.current.offset),
		),
	))

	return nil
}

// readComment advances positions until you finish reading a comment.
func (s *Scanner) readComment() {
	for !s.isAtEnd() && s.peek() != '\n' {
		_ = s.advance()
	}
}

// readNumber advances positions until you finish reading a number.
func (s *Scanner) readNumber() error {
	for s.isDigit(s.peek()) {
		_ = s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.lookAhead()) {
		_ = s.advance() // consume the dot

		for s.isDigit(s.peek()) {
			_ = s.advance()
		}
	}

	if s.isSymbol(s.peek()) {
		s.readSymbol()
		return s.error(fmt.Sprintf("invalid number `%s`", s.currentLexeme()))
	}

	return nil
}

// readString advances positions until you finish reading a string.
func (s *Scanner) readString() error {
	for !s.isAtEnd() && s.peek() != '"' {
		chr := s.advance()
		s.markNewline(chr)

		if chr == '\\' && !s.isAtEnd() {
			s.markNewline(s.advance())
		}
	}

	if s.isAtEnd() {
		return s.error("unterminated string")
	}

	_ = s.advance() // consume the end quote

	return nil
}

// readSymbol advances positions until you finish reading a symbol.
func (s *Scanner) readSymbol() {
	for s.isSymbol(s.peek()) {
		_ = s.advance()
	}
}

// symbolToken retrieves the token corresponding to a reserved word or symbol.
func (s *Scanner) symbolToken() token.Type {
	word := s.source[s.start.offset:s.current.offset]

	switch word {
	case "true", "false":
		return token.Bool
	case "nil":
		return token.Nil
	default:
		return token.Symbol
	}
}

// processEscapes processes escape sequences.
func (s *Scanner) processEscapes(str string) string {
	return escapeReplacer.Replace(str)
}

// isDigit checks if it is a digit.
func (s *Scanner) isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// isAlpha checks if it is a letter.
func (s *Scanner) isAlpha(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

// isAlphanumeric checks if it is an alphanumeric.
func (s *Scanner) isAlphanumeric(r rune) bool {
	return s.isAlpha(r) || s.isDigit(r)
}

// isOperator checks if it is an operator.
func (s *Scanner) isOperator(r rune) bool {
	return strings.ContainsRune("+-*/%=><", r)
}

// isIdentifier checks if it is an identifier.
func (s *Scanner) isIdentifier(r rune) bool {
	return s.isAlphanumeric(r) || strings.ContainsRune("-_?:", r)
}

// isSymbol checks if it is a symbol.
func (s *Scanner) isSymbol(r rune) bool {
	return s.isIdentifier(r) || s.isOperator(r)
}

// error makes an error.
func (s *Scanner) error(msg string) *debug.Error {
	return debug.NewError(msg, location.NewLocation(
		s.filename,
		location.NewPosition(s.start.line, s.start.column, s.start.offset),
		location.NewPosition(s.current.line, s.current.column, s.current.offset),
	))
}
