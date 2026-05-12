package parser

import (
	"github.com/danielspk/tatu-lang/pkg/ast"
	"github.com/danielspk/tatu-lang/pkg/debug"
	"github.com/danielspk/tatu-lang/pkg/location"
)

// SyntaxSugar is responsible for transforming syntactic sugar into tatu language constructs.
type SyntaxSugar struct {
}

// Transform applies syntactic sugar transformations to an expression.
func (ss *SyntaxSugar) Transform(expr *ast.SExpr) error {
	listExpr, ok := (*expr).(*ast.ListExpr)
	if !ok || len(listExpr.List) == 0 {
		return nil
	}

	symbolExpr, ok := listExpr.List[0].(*ast.SymbolExpr)
	if !ok {
		return nil
	}

	switch symbolExpr.Symbol {
	case "def":
		return ss.defToVar(expr)
	case "switch":
		return ss.switchToIf(expr)
	case "for":
		return ss.forToWhile(expr)
	}

	return nil
}

// defToVar transforms `def` expression to `var` expression.
// Example: (def name (params) body) -> (var name (lambda (params) body))
func (ss *SyntaxSugar) defToVar(expr *ast.SExpr) error {
	listExpr, ok := (*expr).(*ast.ListExpr)
	if !ok || len(listExpr.List) != 4 {
		return ss.error("invalid `def` expression", (*expr).Location())
	}

	symbolExpr, ok := listExpr.List[0].(*ast.SymbolExpr)
	if !ok || symbolExpr.Symbol != "def" {
		return ss.error("invalid `def` symbol", listExpr.List[0].Location())
	}

	// locations for synthetic tokens (var and lambda) are derived from the original expression's location
	*expr = ast.NewListExpr(
		[]ast.SExpr{
			ast.NewSymbolExpr("var", listExpr.List[0].Location()),
			listExpr.List[1], // name
			ast.NewListExpr([]ast.SExpr{
				ast.NewSymbolExpr("lambda", listExpr.List[2].Location()),
				listExpr.List[2], // params
				listExpr.List[3], // body
			}, listExpr.List[2].Location()),
		},
		listExpr.Location(),
	)

	return nil
}

// forToWhile transforms `for` expression to `while` expression.
// Example: (for init condition increment body) -> (block init (while condition (block (block (var <v> <v>) body) increment)))
func (ss *SyntaxSugar) forToWhile(expr *ast.SExpr) error {
	listExpr, ok := (*expr).(*ast.ListExpr)
	if !ok || len(listExpr.List) != 5 {
		return ss.error("invalid `for` expression", (*expr).Location())
	}

	symbolExpr, ok := listExpr.List[0].(*ast.SymbolExpr)
	if !ok || symbolExpr.Symbol != "for" {
		return ss.error("invalid `for` symbol", listExpr.List[0].Location())
	}

	// extract loop variable name from init for per-iteration shadowing
	initList, ok := listExpr.List[1].(*ast.ListExpr)
	if !ok || len(initList.List) < 2 {
		return ss.error("invalid `for` init clause", listExpr.List[1].Location())
	}

	nameSym, ok := initList.List[1].(*ast.SymbolExpr)
	if !ok {
		return ss.error("invalid `for` loop variable", initList.List[1].Location())
	}

	// wrap body with shadow block: (block (var name name) body)
	shadowedBody := ast.NewListExpr([]ast.SExpr{
		ast.NewSymbolExpr("block", listExpr.List[4].Location()),
		ast.NewListExpr([]ast.SExpr{
			ast.NewSymbolExpr("var", listExpr.List[4].Location()),
			ast.NewSymbolExpr(nameSym.Symbol, listExpr.List[4].Location()),
			ast.NewSymbolExpr(nameSym.Symbol, listExpr.List[4].Location()),
		}, listExpr.List[4].Location()),
		listExpr.List[4],
	}, listExpr.List[4].Location())

	// locations for synthetic tokens (block and while) are derived from the original expression's location
	*expr = ast.NewListExpr(
		[]ast.SExpr{
			ast.NewSymbolExpr("block", listExpr.List[0].Location()),
			listExpr.List[1], // init
			ast.NewListExpr([]ast.SExpr{
				ast.NewSymbolExpr("while", listExpr.List[1].Location()),
				listExpr.List[2], // condition
				ast.NewListExpr([]ast.SExpr{
					ast.NewSymbolExpr("block", listExpr.List[2].Location()),
					shadowedBody,     // body
					listExpr.List[3], // increment
				}, listExpr.List[2].Location()),
			}, listExpr.List[1].Location()),
		},
		listExpr.Location(),
	)

	return nil
}

// switchToIf transforms `switch` expression to `if` expression.
// Example: (switch ((< 10 10) 1) ((> 10 10) 2) (default 3)) -> (if (< 10 10) 1 (if (> 10 10) 2 3))
func (ss *SyntaxSugar) switchToIf(expr *ast.SExpr) error {
	listExpr, ok := (*expr).(*ast.ListExpr)
	if !ok || len(listExpr.List) < 3 {
		return ss.error("invalid `switch` expression", (*expr).Location())
	}

	symbolExpr, ok := listExpr.List[0].(*ast.SymbolExpr)
	if !ok || symbolExpr.Symbol != "switch" {
		return ss.error("invalid `switch` symbol", listExpr.List[0].Location())
	}

	cases := listExpr.List[1:]

	// verify default case
	defaultCase, ok := cases[len(cases)-1].(*ast.ListExpr)
	if !ok || len(defaultCase.List) != 2 {
		return ss.error("invalid default case", cases[len(cases)-1].Location())
	}

	defaultCond, ok := defaultCase.List[0].(*ast.SymbolExpr)
	if !ok || defaultCond.Symbol != "default" {
		return ss.error("invalid default case symbol", defaultCase.List[0].Location())
	}

	// start with the default value
	result := defaultCase.List[1]

	for i := len(cases) - 2; i >= 0; i-- {
		caseExpr, ok := cases[i].(*ast.ListExpr)
		if !ok || len(caseExpr.List) != 2 {
			return ss.error("invalid case expression", cases[i].Location())
		}

		result = ast.NewListExpr(
			[]ast.SExpr{
				ast.NewSymbolExpr("if", caseExpr.List[0].Location()),
				caseExpr.List[0], // condition
				caseExpr.List[1], // true value
				result,           // else value
			},
			listExpr.Location(),
		)
	}

	*expr = result

	return nil
}

// error makes an error.
func (ss *SyntaxSugar) error(msg string, loc location.Location) *debug.Error {
	return &debug.Error{
		Msg:    msg,
		Line:   loc.End.Line,
		Column: loc.End.Column,
		File:   loc.File,
	}
}
