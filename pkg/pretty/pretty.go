// Package pretty provides colored output formatting for tokens, AST, and error messages.
package pretty

import (
	"errors"
	"fmt"
	"strings"

	"github.com/danielspk/tatu-lang/pkg/ast"
	"github.com/danielspk/tatu-lang/pkg/debug"
	"github.com/danielspk/tatu-lang/pkg/token"
)

// FormatRunningExecution formats running execution.
func FormatRunningExecution(version, filename string) string {
	return fmt.Sprintf("%s>>> Running tatu lang %s(%s)%s - %s%s",
		ColorPurple, ColorDarkGray, version, ColorGreen, filename, ColorReset)
}

// FormatRunningOutput formats running output.
func FormatRunningOutput() string {
	return fmt.Sprintf("%s>>> Result:%s", ColorPink, ColorReset)
}

// FormatError formats error.
func FormatError(err error) string {
	var tatuErr *debug.Error

	if errors.As(err, &tatuErr) {
		return fmt.Sprintf("%s>>> %s%s\n", ColorRed, tatuErr.Dump(), ColorReset)
	}

	return fmt.Sprintf("%s>>> Error: %s%s\n", ColorRed, err, ColorReset)
}

// FormatToken formats token information.
func FormatToken(tok token.Token) string {
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

	return fmt.Sprintf("%s(%s: `%s`)%s => start(%d:%d) end(%d:%d) file(%s)%s",
		tokenColor, tokenType, lexeme, ColorDarkGray, tok.Location.Start.Line,
		tok.Location.Start.Column, tok.Location.End.Line, tok.Location.End.Column,
		tok.Location.File, ColorReset)
}

// FormatAST prints formated ast node.
func FormatAST(ast *ast.AST) string {
	var sb strings.Builder

	for _, exp := range ast.Program {
		prettyExpression(&sb, exp, 0)
	}

	return sb.String()
}

func prettyExpression(sb *strings.Builder, expr ast.SExpr, depth int) {
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
			prettyExpression(sb, e, depth+1)
		}

		sb.WriteString(indent + fmt.Sprintf("%s)", ColorPurple))
	}

	sb.WriteString(fmt.Sprintf("%s => start(%d:%d) end(%d:%d) file(%s)%s\n",
		ColorDarkGray, expr.Location().Start.Line, expr.Location().Start.Column,
		expr.Location().End.Line, expr.Location().End.Column, expr.Location().File, ColorReset))
}
