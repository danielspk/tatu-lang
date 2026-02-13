package parser

import (
	"fmt"

	"github.com/danielspk/tatu-lang/pkg/ast"
	"github.com/danielspk/tatu-lang/pkg/debug"
	"github.com/danielspk/tatu-lang/pkg/location"
)

// SyntaxAnalyzer is responsible for validates the syntax of S-expressions.
type SyntaxAnalyzer struct {
}

// Validate validates an S-expression structure.
func (sa *SyntaxAnalyzer) Validate(expr ast.SExpr) error {
	listExpr, ok := expr.(*ast.ListExpr)
	if !ok || len(listExpr.List) == 0 {
		return nil
	}

	if listExpr.List[0].Kind() != ast.SymbolKind && listExpr.List[0].Kind() != ast.ListKind {
		return sa.error("expected symbol or list", listExpr.List[0].Location())
	}

	symbolExpr, ok := listExpr.List[0].(*ast.SymbolExpr)
	if !ok {
		// if the first element is not a symbol, it can still be a valid case. Example: ((lambda (x) (+ x 1)) 2)
		return nil
	}

	switch symbolExpr.Symbol {
	case "+":
		return sa.validatePlus(listExpr)
	case "-", "*", "/", "%":
		return sa.validateArithmetic(listExpr)
	case "=", "<", "<=", ">", ">=":
		return sa.validateComparison(listExpr)
	case "and", "or":
		return sa.validateLogical(listExpr)
	case "not":
		return sa.validateNot(listExpr)
	case "include":
		return sa.validateInclude(listExpr)
	case "begin":
		return sa.validateBegin(listExpr)
	case "var":
		return sa.validateVar(listExpr)
	case "set":
		return sa.validateSet(listExpr)
	case "if":
		return sa.validateIf(listExpr)
	case "while":
		return sa.validateWhile(listExpr)
	case "lambda":
		return sa.validateLambda(listExpr)
	case "vector":
		return sa.validateVector(listExpr)
	case "map":
		return sa.validateMap(listExpr)
	case "print":
		return sa.validatePrint(listExpr)
	}

	return nil
}

// validatePlus validates the plus native function.
// Format: (+ <expr>+)
func (sa *SyntaxAnalyzer) validatePlus(expr *ast.ListExpr) error {
	if len(expr.List) < 3 {
		return sa.error("invalid `+` format: expected at least two operands", expr.Location())
	}

	return nil
}

// validateArithmetic validates arithmetic native functions (-, *, /).
// Format: (op <expr>+) or (- <expr>) for unary negation
func (sa *SyntaxAnalyzer) validateArithmetic(expr *ast.ListExpr) error {
	operator := expr.List[0].(*ast.SymbolExpr).Symbol

	if len(expr.List) < 2 {
		return sa.error(fmt.Sprintf("invalid `%s` format: expected at least one operand", operator), expr.Location())
	}

	// special case: unary minus is allowed
	if operator == "-" && len(expr.List) == 2 {
		return nil
	}

	// special case: modulo
	if operator == "%" && len(expr.List) != 3 {
		return sa.error("invalid `%` format: expected exactly two operands", expr.Location())
	}

	if operator != "-" && len(expr.List) < 3 {
		return sa.error(fmt.Sprintf("invalid `%s` format: expected at least two operands", operator), expr.Location())
	}

	return nil
}

// validateComparison validates comparison native functions (=, <, <=, >, >=).
// Format: (op <expr> <expr>)
func (sa *SyntaxAnalyzer) validateComparison(expr *ast.ListExpr) error {
	operator := expr.List[0].(*ast.SymbolExpr).Symbol

	if len(expr.List) != 3 {
		return sa.error(fmt.Sprintf("invalid `%s` format: expected exactly two operands", operator), expr.Location())
	}

	return nil
}

// validateLogical validates logical special forms (and, or).
// Format: (op <expr>+)
func (sa *SyntaxAnalyzer) validateLogical(expr *ast.ListExpr) error {
	operator := expr.List[0].(*ast.SymbolExpr).Symbol

	if len(expr.List) < 3 {
		return sa.error(fmt.Sprintf("invalid `%s` format: expected at least two operands", operator), expr.Location())
	}

	return nil
}

// validateNot validates the `not` native function.
// Format: (not <expr>)
func (sa *SyntaxAnalyzer) validateNot(expr *ast.ListExpr) error {
	if len(expr.List) != 2 {
		return sa.error("invalid `not` format: expected exactly one operand", expr.Location())
	}

	return nil
}

// validateInclude validates the `include` special form.
// Format: (include <string>)
func (sa *SyntaxAnalyzer) validateInclude(expr *ast.ListExpr) error {
	if len(expr.List) != 2 {
		return sa.error("invalid `include` format: expected (include <string>)", expr.Location())
	}

	if expr.List[1].Kind() != ast.StringKind {
		return sa.error("invalid `include` argument: expected string", expr.List[1].Location())
	}

	return nil
}

// validateBegin validates the `begin` special form.
// Format: (begin <expr>+)
func (sa *SyntaxAnalyzer) validateBegin(expr *ast.ListExpr) error {
	if len(expr.List) < 2 {
		return sa.error("invalid `begin` format: expected at least one expression", expr.Location())
	}

	return nil
}

// validateVar validates the `var` special form.
// Format: (var <identifier> <expr>)
func (sa *SyntaxAnalyzer) validateVar(expr *ast.ListExpr) error {
	if len(expr.List) != 3 {
		return sa.error("invalid `var` format: expected (var <identifier> <expr>)", expr.Location())
	}

	if expr.List[1].Kind() != ast.SymbolKind {
		return sa.error("invalid `var` name: expected identifier", expr.List[1].Location())
	}

	return nil
}

// validateSet validates the `set` special form.
// Format: (set <identifier> <expr>)
func (sa *SyntaxAnalyzer) validateSet(expr *ast.ListExpr) error {
	if len(expr.List) != 3 {
		return sa.error("invalid `set` format: expected (set <identifier> <expr>)", expr.Location())
	}

	if expr.List[1].Kind() != ast.SymbolKind {
		return sa.error("invalid `set` name: expected identifier", expr.List[1].Location())
	}

	return nil
}

// validateIf validates the `if` special form.
// Format: (if <expr> <expr> [<expr>])
func (sa *SyntaxAnalyzer) validateIf(expr *ast.ListExpr) error {
	if len(expr.List) < 3 || len(expr.List) > 4 {
		return sa.error("invalid `if` format: expected (if <condition> <then> [<else>])", expr.Location())
	}

	return nil
}

// validateWhile validates the `while` special form.
// Format: (while <expr> <expr>)
func (sa *SyntaxAnalyzer) validateWhile(expr *ast.ListExpr) error {
	if len(expr.List) != 3 {
		return sa.error("invalid `while` format: expected (while <condition> <body>)", expr.Location())
	}

	return nil
}

// validateLambda validates the `lambda` special form.
// Format: (lambda (<identifier>*) <expr>)
func (sa *SyntaxAnalyzer) validateLambda(expr *ast.ListExpr) error {
	if len(expr.List) != 3 {
		return sa.error("invalid `lambda` format: expected (lambda (<params>) <body>)", expr.Location())
	}

	if expr.List[1].Kind() != ast.ListKind {
		return sa.error("invalid `lambda` params: expected list", expr.List[1].Location())
	}

	params := expr.List[1].(*ast.ListExpr)
	for _, param := range params.List {
		if param.Kind() != ast.SymbolKind {
			return sa.error("invalid `lambda` param: expected identifier", param.Location())
		}
	}

	return nil
}

// validateVector validates the `vector` special form.
// Format: (vector <expr>*)
func (sa *SyntaxAnalyzer) validateVector(_ *ast.ListExpr) error {
	return nil // no validation required
}

// validateMap validates the `map` special form.
// Format: (map <key-value>*)
func (sa *SyntaxAnalyzer) validateMap(expr *ast.ListExpr) error {
	if len(expr.List)%2 != 1 {
		return sa.error("invalid `map` format: expected (map <key-value>*)", expr.Location())
	}

	return nil
}

// validatePrint validates the `print` native function.
// Format: (print <expr>*)
func (sa *SyntaxAnalyzer) validatePrint(expr *ast.ListExpr) error {
	if len(expr.List) < 2 {
		return sa.error("invalid `print` format: expected at least one expression", expr.Location())
	}

	return nil
}

// error makes an error.
func (sa *SyntaxAnalyzer) error(msg string, loc location.Location) *debug.Error {
	return &debug.Error{
		Msg:    msg,
		Line:   loc.End.Line,
		Column: loc.End.Column,
		File:   loc.File,
	}
}
