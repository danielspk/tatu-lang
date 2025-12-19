package parser

import (
	"errors"

	"github.com/danielspk/tatu-lang/pkg/ast"
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
		return errors.New("invalid `def` expression")
	}

	symbolExpr, ok := listExpr.List[0].(*ast.SymbolExpr)
	if !ok || symbolExpr.Symbol != "def" {
		return errors.New("invalid `def` symbol")
	}

	for _, e := range listExpr.List[1:] {
		if err := ss.Transform(&e); err != nil {
			return err
		}
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
// Example: (for init condition body) -> (begin init (while condition (begin body increment)))
func (ss *SyntaxSugar) forToWhile(expr *ast.SExpr) error {
	listExpr, ok := (*expr).(*ast.ListExpr)
	if !ok || len(listExpr.List) != 5 {
		return errors.New("invalid `for` expression")
	}

	symbolExpr, ok := listExpr.List[0].(*ast.SymbolExpr)
	if !ok || symbolExpr.Symbol != "for" {
		return errors.New("invalid `for` symbol")
	}

	for _, e := range listExpr.List[1:] {
		if err := ss.Transform(&e); err != nil {
			return err
		}
	}

	// locations for synthetic tokens (begin and while) are derived from the original expression's location
	*expr = ast.NewListExpr(
		[]ast.SExpr{
			ast.NewSymbolExpr("begin", listExpr.List[0].Location()),
			listExpr.List[1], // init
			ast.NewListExpr([]ast.SExpr{
				ast.NewSymbolExpr("while", listExpr.List[1].Location()),
				listExpr.List[2], // condition
				ast.NewListExpr([]ast.SExpr{
					ast.NewSymbolExpr("begin", listExpr.List[2].Location()),
					listExpr.List[4], // body
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
		return errors.New("invalid `switch` expression")
	}

	symbolExpr, ok := listExpr.List[0].(*ast.SymbolExpr)
	if !ok || symbolExpr.Symbol != "switch" {
		return errors.New("invalid `switch` symbol")
	}

	cases := listExpr.List[1:]

	for _, e := range cases {
		if err := ss.Transform(&e); err != nil {
			return err
		}
	}

	// verify default case
	defaultCase, ok := cases[len(cases)-1].(*ast.ListExpr)
	if !ok || len(defaultCase.List) != 2 {
		return errors.New("invalid default case")
	}

	defaultCond, ok := defaultCase.List[0].(*ast.SymbolExpr)
	if !ok || defaultCond.Symbol != "default" {
		return errors.New("invalid default case symbol")
	}

	// start with the default value
	result := defaultCase.List[1]

	for i := len(cases) - 2; i >= 0; i-- {
		caseExpr, ok := cases[i].(*ast.ListExpr)
		if !ok || len(caseExpr.List) != 2 {
			return errors.New("invalid case expression")
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
