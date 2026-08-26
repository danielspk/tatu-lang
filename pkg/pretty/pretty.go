// Package pretty provides colored output formatting for tokens, AST, and error messages.
package pretty

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/danielspk/tatu-lang/pkg/ast"
	"github.com/danielspk/tatu-lang/pkg/debug"
	"github.com/danielspk/tatu-lang/pkg/token"
)

var headerRule = strings.Repeat("-", 80)

// DumpExecution dumps the running execution.
func DumpExecution(version, filename string) string {
	return fmt.Sprintf("%s>>> Tatu lang %s(%s)\n%s>>> File - %s\n%s>>> Result:%s\n",
		ColorPurple, ColorDarkGray, version, ColorGreen, filename, ColorPink, ColorReset)
}

// DumpError dumps an error.
func DumpError(err error, sources map[string][]byte) string {
	if tatuErr, ok := errors.AsType[*debug.Error](err); ok {
		return dumpError(tatuErr, sources[tatuErr.Loc.File])
	}

	return fmt.Sprintf("%s>>> Error: %s%s\n", ColorRed, err, ColorReset)
}

// DumpTokensHeader dumps the token header.
func DumpTokensHeader() string {
	return fmt.Sprintf("%s%s\n-- TOKENS\n%s\n%s", ColorPurple, headerRule, headerRule, ColorReset)
}

// DumpTokens dumps the token information.
func DumpTokens(tokens []token.Token) string {
	var sb strings.Builder

	for _, tok := range tokens {
		dumpToken(&sb, tok)
	}

	return sb.String()
}

// DumpASTHeader dumps the ast header.
func DumpASTHeader() string {
	return fmt.Sprintf("%s%s\n-- AST\n%s\n%s", ColorGreen, headerRule, headerRule, ColorReset)
}

// dumpToken writes the token representation.
func dumpToken(sb *strings.Builder, tok token.Token) {
	tokenType := ""
	tokenColor := ""

	switch tok.Type {
	case token.LeftParen:
		tokenType = "LEFT_PAREN"
		tokenColor = ColorPurple
	case token.RightParen:
		tokenType = "RIGHT_PAREN"
		tokenColor = ColorPurple
	case token.Number:
		tokenType = "NUMBER"
		tokenColor = ColorGreen
	case token.String:
		tokenType = "STRING"
		tokenColor = ColorOrange
	case token.Bool:
		tokenType = "BOOL"
		tokenColor = ColorPink
	case token.Nil:
		tokenType = "NIL"
		tokenColor = ColorLightGray
	case token.Symbol:
		tokenType = "SYMBOL"
		tokenColor = ColorCyan
	case token.EOF:
		tokenType = "EOF"
		tokenColor = ColorYellow
	}

	lexeme := strings.Replace(tok.Lexeme, "\n", "\\n", -1)

	_, _ = fmt.Fprintf(sb, "%s(%s: `%s`)%s => start(%d:%d) end(%d:%d) file(%s)%s\n",
		tokenColor, tokenType, lexeme, ColorDarkGray, tok.Location.Start.Line,
		tok.Location.Start.Column, tok.Location.End.Line, tok.Location.End.Column,
		relPath(tok.Location.File), ColorReset)
}

// DumpAST dumps the ast information.
func DumpAST(ast *ast.AST) string {
	var sb strings.Builder

	for _, exp := range ast.Program {
		dumpExpression(&sb, exp, 0)
	}

	return sb.String()
}

// dumpError dumps the error message next to reference source code.
func dumpError(e *debug.Error, source []byte) string {
	var sb strings.Builder

	startCol := int(e.Loc.Start.Column)
	line := int(e.Loc.Start.Line)

	_, _ = fmt.Fprintf(&sb, "%s>>> Error on line %d, column %d, file `%s`:\n\n",
		ColorRed, line, startCol, relPath(e.Loc.File))

	lines := strings.Split(string(source), "\n")
	if line == 0 || line > len(lines) {
		_, _ = fmt.Fprintf(&sb, "%s%s\n", e.Msg, ColorReset)

		return sb.String()
	}

	endCol := int(e.Loc.End.Column)
	if e.Loc.Start.Line != e.Loc.End.Line {
		endCol = utf8.RuneCountInString(lines[line-1]) + 1
	}

	width := len(fmt.Sprintf("%d", line))

	for i := max(line-2, 1); i <= line; i++ {
		_, _ = fmt.Fprintf(&sb, "%s%*d │ %s%s\n", ColorDarkGray, width, i, ColorLightGray, lines[i-1])
	}

	span := max(endCol-startCol, 1)
	offset := width + 3 + startCol - 1

	_, _ = fmt.Fprintf(&sb, "%s%s%s\n", ColorYellow, strings.Repeat(" ", offset), strings.Repeat("^", span))
	_, _ = fmt.Fprintf(&sb, "%s└─ %s%s\n", strings.Repeat(" ", offset+span-1), e.Msg, ColorReset)

	return sb.String()
}

func dumpExpression(sb *strings.Builder, expr ast.SExpr, depth int) {
	switch expr.(type) {
	case *ast.NumberExpr:
		sb.WriteString(fmt.Sprintf("%s(Number %v)", ColorGreen, expr.(*ast.NumberExpr).Number))
	case *ast.StringExpr:
		sb.WriteString(fmt.Sprintf("%s(String \"%s\")", ColorOrange, expr.(*ast.StringExpr).String))
	case *ast.BoolExpr:
		sb.WriteString(fmt.Sprintf("%s(Bool %t)", ColorPink, expr.(*ast.BoolExpr).Bool))
	case *ast.SymbolExpr:
		sb.WriteString(fmt.Sprintf("%s(Symbol %s)", ColorCyan, expr.(*ast.SymbolExpr).Symbol))
	case *ast.NilExpr:
		sb.WriteString(fmt.Sprintf("%s(Nil)", ColorLightGray))
	case *ast.ListExpr:
		listExpr := expr.(*ast.ListExpr)

		if len(listExpr.List) == 0 {
			sb.WriteString(fmt.Sprintf("%s(List)", ColorPurple))
			break
		}

		indent := strings.Repeat("    ", depth)

		sb.WriteString(fmt.Sprintf("%s(List\n", ColorPurple))
		for i, e := range listExpr.List {
			connector := "├─ "
			if i == len(listExpr.List)-1 {
				connector = "└─ "
			}

			sb.WriteString(indent + fmt.Sprintf(" %s%s", ColorPurple, connector))
			dumpExpression(sb, e, depth+1)
		}

		sb.WriteString(indent + fmt.Sprintf("%s)", ColorPurple))
	}

	sb.WriteString(fmt.Sprintf("%s => start(%d:%d) end(%d:%d) file(%s)%s\n",
		ColorDarkGray, expr.Location().Start.Line, expr.Location().Start.Column,
		expr.Location().End.Line, expr.Location().End.Column, relPath(expr.Location().File), ColorReset))
}

// relPath returns the path relative to the current working directory.
func relPath(path string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return path
	}

	rel, err := filepath.Rel(cwd, path)
	if err != nil {
		return path
	}

	return rel
}
