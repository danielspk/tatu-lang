// Package interpreter implements the tree-walking interpreter.
package interpreter

import (
	"fmt"

	"github.com/danielspk/tatu-lang/pkg/ast"
	"github.com/danielspk/tatu-lang/pkg/core/builtins"
	"github.com/danielspk/tatu-lang/pkg/core/stdlib"
	"github.com/danielspk/tatu-lang/pkg/debug"
	"github.com/danielspk/tatu-lang/pkg/location"
	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// Interpreter represents a tree-walking interpreter.
type Interpreter struct {
	natives map[string]runtime.NativeFunction
	global     *runtime.Environment
}

// NewInterpreter builds a new Interpreter.
func NewInterpreter() (*Interpreter, error) {
	natives := make(map[string]runtime.NativeFunction)

	builtins.RegisterArithmetic(natives)
	builtins.RegisterComparison(natives)
	builtins.RegisterIO(natives)
	builtins.RegisterTypes(natives)

	stdlib.RegisterFileSystem(natives)
	stdlib.RegisterJSON(natives)
	stdlib.RegisterMap(natives)
	stdlib.RegisterMath(natives)
	stdlib.RegisterRegex(natives)
	stdlib.RegisterString(natives)
	stdlib.RegisterTime(natives)
	stdlib.RegisterVector(natives)

	return &Interpreter{
		natives: natives,
		global:     runtime.NewEnvironment(nil, nil),
	}, nil
}

// Environment returns the global environment.
func (i *Interpreter) Environment() *runtime.Environment {
	return i.global
}

// Eval evaluates an S-expression and returns the resulting value.
// Note: the format of the S-expressions is guaranteed by the syntax analyzer.
func (i *Interpreter) Eval(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	return i.eval(expr, env)
}

// EvalProgram evaluates an AST and returns the resulting value.
func (i *Interpreter) EvalProgram(ast *ast.AST, env *runtime.Environment) (runtime.Value, error) {
	var lastValue runtime.Value
	var err error

	for _, expr := range ast.Program {
		lastValue, err = i.eval(expr, env)
		if err != nil {
			return nil, err
		}
	}

	return lastValue, nil
}

// eval evaluates an S-expression in non-tail position.
func (i *Interpreter) eval(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	result, err := i.evalInTailPosition(expr, env)
	if err != nil {
		return nil, err
	}

	if result != nil && result.Type() == runtime.RecurType {
		return nil, i.error("recur can only be used in tail position of a function", expr.Location())
	}

	return result, nil
}

// evalInTailPosition evaluates an S-expression in tail position.
func (i *Interpreter) evalInTailPosition(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	if env == nil {
		env = i.global
	}

	switch expr.(type) {
	case *ast.NumberExpr, *ast.StringExpr, *ast.BoolExpr, *ast.NilExpr, *ast.SymbolExpr:
		return i.evalAtom(expr, env)
	case *ast.ListExpr:
		return i.evalList(expr, env)
	default:
		return nil, i.error("unknown expression type", expr.Location())
	}
}

// evalAtom evaluates an atomic expression.
func (i *Interpreter) evalAtom(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	switch expr.(type) {
	case *ast.NumberExpr:
		return runtime.NewNumber(expr.(*ast.NumberExpr).Number), nil
	case *ast.StringExpr:
		return runtime.NewString(expr.(*ast.StringExpr).String), nil
	case *ast.BoolExpr:
		return runtime.NewBool(expr.(*ast.BoolExpr).Bool), nil
	case *ast.NilExpr:
		return runtime.NewNil(), nil
	case *ast.SymbolExpr:
		return i.evalSymbol(expr, env)
	default:
		return nil, i.error("unknown atom", expr.Location())
	}
}

// evalSymbol evaluates a symbol expression.
func (i *Interpreter) evalSymbol(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	if expr.Kind() != ast.SymbolKind {
		return nil, i.error("invalid symbol expression", expr.Location())
	}

	exprSymbol := expr.(*ast.SymbolExpr)

	value, found := env.Lookup(exprSymbol.Symbol)
	if found {
		return value, nil
	}

	if fn, ok := i.natives[exprSymbol.Symbol]; ok {
		return fn, nil
	}

	return nil, i.error(fmt.Sprintf("unknown symbol `%s`", exprSymbol.Symbol), exprSymbol.Location())
}

// evalList evaluates a list expression.
func (i *Interpreter) evalList(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	if expr.Kind() != ast.ListKind {
		return nil, i.error("invalid list expression", expr.Location())
	}

	exprList := expr.(*ast.ListExpr)

	if len(exprList.List) == 0 {
		return runtime.NewNil(), nil
	}

	if exprList.List[0].Kind() == ast.SymbolKind {
		exprSymbol := exprList.List[0].(*ast.SymbolExpr)

		switch exprSymbol.Symbol {
		case "include":
			return nil, i.error("include not resolved", exprList.Location())
		case "and", "or":
			return i.evalLogical(exprList, env)
		case "begin":
			return i.evalBegin(exprList, env)
		case "var":
			return i.evalVar(exprList, env)
		case "set":
			return i.evalSet(exprList, env)
		case "if":
			return i.evalIf(exprList, env)
		case "while":
			return i.evalWhile(exprList, env)
		case "lambda":
			return i.evalLambda(exprList, env)
		case "recur":
			return i.evalRecur(exprList, env)
		case "vector":
			return i.evalVector(exprList, env)
		case "map":
			return i.evalMap(exprList, env)
		}
	}

	// call function
	return i.evalCallFunction(exprList, env)
}

// evalLogical evaluates the `and` and `or` logical special forms.
func (i *Interpreter) evalLogical(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	exprList := expr.(*ast.ListExpr)
	operator := exprList.List[0].(*ast.SymbolExpr).Symbol

	for _, e := range exprList.List[1:] {
		result, err := i.eval(e, env)
		if err != nil {
			return nil, err
		}

		if result.Type() != runtime.BoolType {
			return nil, i.error(fmt.Sprintf("invalid type %s for `%s`", result.Type(), operator), e.Location())
		}

		value := result.(runtime.Bool).Value

		if operator == "and" && !value {
			return runtime.NewBool(false), nil
		}

		if operator == "or" && value {
			return runtime.NewBool(true), nil
		}
	}

	return runtime.NewBool(operator == "and"), nil
}

// evalBegin evaluates a `begin` expression (block of expressions).
func (i *Interpreter) evalBegin(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	exprList := expr.(*ast.ListExpr)

	newEnv := runtime.NewEnvironment(nil, env)

	// eval all expressions except the last (in tail position)
	for _, e := range exprList.List[1 : len(exprList.List)-1] {
		if _, err := i.eval(e, newEnv); err != nil {
			return nil, err
		}
	}

	return i.evalInTailPosition(exprList.List[len(exprList.List)-1], newEnv)
}

// evalVar evaluates a `var` expression.
func (i *Interpreter) evalVar(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	exprList := expr.(*ast.ListExpr)

	value, err := i.eval(exprList.List[2], env)
	if err != nil {
		return nil, err
	}

	return env.Define(exprList.List[1].(*ast.SymbolExpr).Symbol, value)
}

// evalSet evaluates a `set` expression.
func (i *Interpreter) evalSet(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	exprList := expr.(*ast.ListExpr)

	value, err := i.eval(exprList.List[2], env)
	if err != nil {
		return nil, err
	}

	if found := env.Assign(exprList.List[1].(*ast.SymbolExpr).Symbol, value); !found {
		return nil, i.error(fmt.Sprintf("undefined variable `%s`", exprList.List[1].(*ast.SymbolExpr).Symbol), exprList.List[1].Location())
	}

	return value, nil
}

// evalIf evaluates an `if` expression.
func (i *Interpreter) evalIf(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	exprList := expr.(*ast.ListExpr)

	condition := exprList.List[1]
	consequent := exprList.List[2]

	value, err := i.eval(condition, env)
	if err != nil {
		return nil, err
	}

	if value.Type() != runtime.BoolType {
		return nil, i.error(fmt.Sprintf("expected BOOL, found %s", value.Type()), condition.Location())
	}

	if value.(runtime.Bool).Value {
		return i.evalInTailPosition(consequent, env)
	}

	if len(exprList.List) == 4 {
		return i.evalInTailPosition(exprList.List[3], env)
	}

	return runtime.NewNil(), nil
}

// evalWhile evaluates a `while` expression.
func (i *Interpreter) evalWhile(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	exprList := expr.(*ast.ListExpr)

	condition := exprList.List[1]
	body := exprList.List[2]

	var lastValue runtime.Value

	for {
		value, err := i.eval(condition, env)
		if err != nil {
			return nil, err
		}

		if value.Type() != runtime.BoolType {
			return nil, i.error(fmt.Sprintf("expected BOOL, found %s", value.Type()), condition.Location())
		}

		if !value.(runtime.Bool).Value {
			break
		}

		lastValue, err = i.eval(body, env)
		if err != nil {
			return nil, err
		}
	}

	return lastValue, nil
}

// evalLambda evaluates a `lambda` expression.
func (i *Interpreter) evalLambda(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	exprList := expr.(*ast.ListExpr)

	params := exprList.List[1]
	body := exprList.List[2]

	return runtime.NewFunction(env, params, body), nil
}

// evalRecur evaluates a `recur` expression for TCO.
func (i *Interpreter) evalRecur(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	exprList := expr.(*ast.ListExpr)

	args := make([]runtime.Value, 0, len(exprList.List)-1)

	for _, e := range exprList.List[1:] {
		result, err := i.eval(e, env)
		if err != nil {
			return nil, err
		}

		args = append(args, result)
	}

	return runtime.NewRecurBindings(args), nil
}

// evalVector evaluates a `vector` expression.
func (i *Interpreter) evalVector(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	exprList := expr.(*ast.ListExpr)

	elements := make([]runtime.Value, 0, len(exprList.List)-1)

	for _, e := range exprList.List[1:] {
		result, err := i.eval(e, env)
		if err != nil {
			return nil, err
		}

		elements = append(elements, result)
	}

	return runtime.NewVector(elements), nil
}

// evalMap evaluates a `map` expression.
func (i *Interpreter) evalMap(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	exprList := expr.(*ast.ListExpr)

	elements := make(map[string]runtime.Value, (len(exprList.List)-1)/2)

	for idx := 1; idx < len(exprList.List); idx += 2 {
		keyExpr := exprList.List[idx]
		valueExpr := exprList.List[idx+1]

		key, err := i.eval(keyExpr, env)
		if err != nil {
			return nil, err
		}

		if key.Type() != runtime.StringType {
			return nil, i.error(fmt.Sprintf("invalid map key type: expected string, got %s", key.Type()), keyExpr.Location())
		}

		result, err := i.eval(valueExpr, env)
		if err != nil {
			return nil, err
		}

		elements[key.(runtime.String).Value] = result
	}

	return runtime.NewMap(elements), nil
}

// evalCallFunction evaluates a call function expression with tail-call optimization support.
func (i *Interpreter) evalCallFunction(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	exprList := expr.(*ast.ListExpr)

	funcValue, err := i.eval(exprList.List[0], env)
	if err != nil {
		return nil, err
	}

	if funcValue.Type() != runtime.NativeFuncType && funcValue.Type() != runtime.FuncType {
		return nil, i.error("expression is not a function", exprList.List[0].Location())
	}

	valArgs, err := i.evalFunctionArguments(exprList.List[1:], env)
	if err != nil {
		return nil, err
	}

	// native function
	if funcValue.Type() == runtime.NativeFuncType {
		result, err := funcValue.(runtime.NativeFunction).Value(valArgs...)
		if err != nil {
			return nil, i.error(err.Error(), exprList.Location())
		}

		return result, nil
	}

	// lambda function
	fn := funcValue.(runtime.Function)
	currentArgs := valArgs

	for {
		activationRecord := make(map[string]runtime.Value)
		for pidx, p := range fn.Params.(*ast.ListExpr).List {
			activationRecord[p.(*ast.SymbolExpr).Symbol] = currentArgs[pidx]
		}
		activationEnv := runtime.NewEnvironment(activationRecord, fn.Env)

		result, err := i.evalInTailPosition(fn.Body, activationEnv)
		if err != nil {
			return nil, err
		}

		if result.Type() == runtime.RecurType {
			currentArgs = result.(runtime.RecurBindings).Args

			continue
		}

		return result, nil
	}
}

// evalFunctionArguments evaluates a list of argument expressions.
func (i *Interpreter) evalFunctionArguments(exprArgs []ast.SExpr, env *runtime.Environment) ([]runtime.Value, error) {
	results := make([]runtime.Value, 0, len(exprArgs))

	for _, e := range exprArgs {
		result, err := i.eval(e, env)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

// error makes an error.
func (i *Interpreter) error(msg string, loc location.Location) *debug.Error {
	return &debug.Error{
		Msg:    msg,
		Line:   loc.End.Line,
		Column: loc.End.Column,
		File:   loc.File,
	}
}
