// Package macro implements macro expansion with AST rewriting.
package macro

import (
	"fmt"

	"github.com/danielspk/tatu-lang/pkg/ast"
	"github.com/danielspk/tatu-lang/pkg/debug"
	"github.com/danielspk/tatu-lang/pkg/location"
)

const maxDepth = 64

var reservedNames = map[string]bool{
	"and": true, "or": true, "block": true, "var": true, "set": true,
	"if": true, "while": true, "lambda": true, "recur": true,
	"vector": true, "map": true, "include": true, "def": true,
	"for": true, "switch": true, "macro": true, "...": true,
}

// rule represents a macro rule with params and a template.
type rule struct {
	params     []string
	template   ast.SExpr
	isVariadic bool
}

// bindings maps variable names to captured AST subtrees.
type bindings map[string][]ast.SExpr

// Expander represents a macro expander.
type Expander struct {
	macros map[string][]rule
}

// NewExpander builds a new Expander.
func NewExpander() *Expander {
	return &Expander{macros: make(map[string][]rule)}
}

// isMacroDef reports whether expr is a macro definition form.
func isMacroDef(expr ast.SExpr) bool {
	list, ok := expr.(*ast.ListExpr)
	if !ok || len(list.List) == 0 {
		return false
	}

	sym, ok := list.List[0].(*ast.SymbolExpr)
	return ok && sym.Symbol == "macro"
}

// Expand registers and expands declared macros in the program.
func (e *Expander) Expand(program *ast.AST) (*ast.AST, error) {
	output := make([]ast.SExpr, 0, len(program.Program))

	for _, expr := range program.Program {
		if err := e.tryRegisterMacro(expr); err != nil {
			return nil, err
		}

		if isMacroDef(expr) {
			continue
		}

		rewritten, err := e.expandExpr(expr, 0)
		if err != nil {
			return nil, err
		}

		output = append(output, rewritten)
	}

	program.Program = output

	return program, nil
}

// tryRegisterMacro registers expr as a macro if it is one.
func (e *Expander) tryRegisterMacro(expr ast.SExpr) error {
	if !isMacroDef(expr) {
		return nil
	}

	form := expr.(*ast.ListExpr)

	if len(form.List) < 2 {
		return e.error("expected (macro <name> (<params>) <body>)", form.Location())
	}

	nameExpr, ok := form.List[1].(*ast.SymbolExpr)
	if !ok {
		return e.error("expected an identifier", form.List[1].Location())
	}

	name := nameExpr.Symbol

	if reservedNames[name] {
		return e.error(fmt.Sprintf("cannot define `%s` as macro", name), form.Location())
	}

	rest := form.List[2:]
	if len(rest) == 0 {
		return e.error("expected (macro <name> (<params>) <body>)", form.Location())
	}

	firstList, ok := rest[0].(*ast.ListExpr)
	if ok && len(firstList.List) > 0 && firstList.List[0].Kind() == ast.ListKind {
		return e.registerMultiRules(name, rest, form.Location())
	}

	return e.registerSingleRule(name, rest, form.Location())
}

// registerSingleRule registers a macro with a single rule.
func (e *Expander) registerSingleRule(name string, parts []ast.SExpr, loc location.Location) error {
	if len(parts) < 2 {
		return e.error("expected (macro <name> (<params>) <body>)", loc)
	}

	paramsList, ok := parts[0].(*ast.ListExpr)
	if !ok {
		return e.error("expected parameter list", parts[0].Location())
	}

	params, isVariadic, err := e.parseParams(paramsList)
	if err != nil {
		return err
	}

	e.macros[name] = append(e.macros[name], rule{
		params:     params,
		template:   parts[1],
		isVariadic: isVariadic,
	})

	return nil
}

// registerMultiRules registers a macro with multiple rules.
func (e *Expander) registerMultiRules(name string, ruleExprs []ast.SExpr, loc location.Location) error {
	for _, ruleExpr := range ruleExprs {
		ruleList, ok := ruleExpr.(*ast.ListExpr)
		if !ok || len(ruleList.List) != 2 {
			return e.error("expected ((<params>) <body>)", ruleExpr.Location())
		}

		paramsList, ok := ruleList.List[0].(*ast.ListExpr)
		if !ok {
			return e.error("expected parameter list", ruleList.List[0].Location())
		}

		params, isVariadic, err := e.parseParams(paramsList)
		if err != nil {
			return err
		}

		e.macros[name] = append(e.macros[name], rule{
			params:     params,
			template:   ruleList.List[1],
			isVariadic: isVariadic,
		})
	}

	return nil
}

// parseParams parses a parameter list and returns params and whether it is variadic.
func (e *Expander) parseParams(paramsList *ast.ListExpr) ([]string, bool, error) {
	params := make([]string, 0, len(paramsList.List))
	isVariadic := false

	for i, param := range paramsList.List {
		sym, ok := param.(*ast.SymbolExpr)
		if !ok {
			return nil, false, e.error("expected identifier", param.Location())
		}

		if sym.Symbol == "..." {
			if i != len(paramsList.List)-1 {
				return nil, false, e.error("`...` must be last in parameter list", param.Location())
			}

			isVariadic = true
			continue
		}

		params = append(params, sym.Symbol)
	}

	return params, isVariadic, nil
}

// expandExpr recursively expands macro uses in an expression.
func (e *Expander) expandExpr(expr ast.SExpr, depth int) (ast.SExpr, error) {
	list, ok := expr.(*ast.ListExpr)
	if !ok || len(list.List) == 0 {
		return expr, nil
	}

	if list.List[0].Kind() == ast.SymbolKind {
		name := list.List[0].(*ast.SymbolExpr).Symbol

		if rules, found := e.macros[name]; found {
			return e.expandMacro(list, rules, name, depth)
		}
	}

	return e.expandChildren(list, depth)
}

// expandMacro expands a macro call against its matching rule.
func (e *Expander) expandMacro(list *ast.ListExpr, rules []rule, name string, depth int) (ast.SExpr, error) {
	if depth >= maxDepth {
		return nil, e.error(fmt.Sprintf("macro `%s` expansion depth limit exceeded (%d)", name, maxDepth), list.Location())
	}

	for _, r := range rules {
		caps := bindings{}

		if e.bindArgs(r, list.List[1:], caps) {
			return e.expandExpr(e.substitute(r.template, caps, list.Location()), depth+1)
		}
	}

	return nil, e.error(fmt.Sprintf("no rule matches macro `%s` (%d arguments)", name, len(list.List)-1), list.Location())
}

// bindArgs binds arguments to parameters.
func (e *Expander) bindArgs(r rule, args []ast.SExpr, caps bindings) bool {
	if r.isVariadic {
		if len(args) < len(r.params) {
			return false
		}

		for i, param := range r.params {
			caps[param] = []ast.SExpr{args[i]}
		}

		caps["..."] = args[len(r.params):]
	} else {
		if len(args) != len(r.params) {
			return false
		}

		for i, param := range r.params {
			caps[param] = []ast.SExpr{args[i]}
		}
	}

	return true
}

// expandChildren expands all children of a list expression.
func (e *Expander) expandChildren(list *ast.ListExpr, depth int) (ast.SExpr, error) {
	changed := false
	children := make([]ast.SExpr, len(list.List))

	for idx, child := range list.List {
		rewritten, err := e.expandExpr(child, depth)
		if err != nil {
			return nil, err
		}

		children[idx] = rewritten

		if rewritten != child {
			changed = true
		}
	}

	if !changed {
		return list, nil
	}

	return ast.NewListExpr(children, list.Location()), nil
}

// substitute builds an AST from a template, splicing captures for pattern variables.
func (e *Expander) substitute(template ast.SExpr, caps bindings, loc location.Location) ast.SExpr {
	if sym, ok := template.(*ast.SymbolExpr); ok {
		if captured, has := caps[sym.Symbol]; has {
			return captured[0]
		}

		return ast.NewSymbolExpr(sym.Symbol, loc)
	}

	list, ok := template.(*ast.ListExpr)
	if !ok {
		switch n := template.(type) {
		case *ast.NumberExpr:
			return ast.NewNumberExpr(n.Number, loc)
		case *ast.StringExpr:
			return ast.NewStringExpr(n.String, loc)
		case *ast.BoolExpr:
			return ast.NewBoolExpr(n.Bool, loc)
		case *ast.NilExpr:
			return ast.NewNilExpr(loc)
		}

		return template
	}

	return ast.NewListExpr(e.substituteList(list.List, caps, loc), loc)
}

// substituteList builds a list from template items, splicing captures.
func (e *Expander) substituteList(items []ast.SExpr, caps bindings, loc location.Location) []ast.SExpr {
	out := make([]ast.SExpr, 0, len(items))

	for _, item := range items {
		sym, ok := item.(*ast.SymbolExpr)
		if !ok {
			out = append(out, e.substitute(item, caps, loc))
			continue
		}

		spliced, ok := spliceCapture(sym.Symbol, caps, out)
		if ok {
			out = spliced
			continue
		}

		out = append(out, e.substitute(item, caps, loc))
	}

	return out
}

// spliceCapture appends captured values for name and returns the updated slice.
// It handles both ellipsis and regular parameter captures.
func spliceCapture(name string, caps bindings, out []ast.SExpr) ([]ast.SExpr, bool) {
	if name == "..." {
		captured, has := caps["..."]
		if !has {
			return out, true
		}

		return append(out, captured...), true
	}

	captured, has := caps[name]
	if !has {
		return out, false
	}

	return append(out, captured...), true
}

// error builds a macro expansion error with location.
func (e *Expander) error(msg string, loc location.Location) *debug.Error {
	return &debug.Error{Msg: msg, Line: loc.End.Line, Column: loc.End.Column, File: loc.File}
}
