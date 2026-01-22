// Package interpreter implements the tree-walking interpreter.
package interpreter

import (
	"fmt"
	"strings"

	"github.com/danielspk/tatu-lang/pkg/ast"
	"github.com/danielspk/tatu-lang/pkg/debug"
	"github.com/danielspk/tatu-lang/pkg/location"
	"github.com/danielspk/tatu-lang/pkg/runtime"
	"github.com/danielspk/tatu-lang/pkg/stdlib"
)

// Interpreter represents a tree-walking interpreter.
type Interpreter struct {
	global *runtime.Environment
}

// NewInterpreter builds a new Interpreter.
func NewInterpreter() (*Interpreter, error) {
	env := runtime.NewEnvironment(nil, nil)

	if err := stdlib.RegisterCasting(env); err != nil {
		return nil, err
	}
	if err := stdlib.RegisterFileSystem(env); err != nil {
		return nil, err
	}
	if err := stdlib.RegisterJSON(env); err != nil {
		return nil, err
	}
	if err := stdlib.RegisterMap(env); err != nil {
		return nil, err
	}
	if err := stdlib.RegisterMath(env); err != nil {
		return nil, err
	}
	if err := stdlib.RegisterRegex(env); err != nil {
		return nil, err
	}
	if err := stdlib.RegisterString(env); err != nil {
		return nil, err
	}
	if err := stdlib.RegisterTime(env); err != nil {
		return nil, err
	}
	if err := stdlib.RegisterTypes(env); err != nil {
		return nil, err
	}
	if err := stdlib.RegisterVector(env); err != nil {
		return nil, err
	}

	return &Interpreter{
		global: env,
	}, nil
}

// Eval evaluates an S-expression and returns the resulting value.
// Note: the format of the S-expressions is guaranteed by the syntax analyzer.
func (i *Interpreter) Eval(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	return i.eval(expr, env)
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
	if !found {
		return nil, i.error(fmt.Sprintf("unknown symbol `%s`", exprSymbol.Symbol), exprSymbol.Location())
	}

	return value, nil
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
		case "+":
			return i.evalPlusSymbol(exprList, env)
		case "-", "*", "/":
			return i.evalMathSymbol(exprList, env)
		case "=", ">", ">=", "<", "<=", "and", "or":
			return i.evalLogicalSymbol(exprList, env)
		case "include":
			return nil, i.error("include not resolver", exprList.Location())
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
		case "print":
			return i.evalPrint(exprList, env)
		}
	}

	// call function
	return i.evalCallFunction(exprList, env)
}

// evalPlusSymbol evaluates the plus operator (addition or concatenation).
func (i *Interpreter) evalPlusSymbol(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	operator := expr.(*ast.ListExpr).List[0].(*ast.SymbolExpr).Symbol

	results := make([]runtime.Value, 0, len(expr.(*ast.ListExpr).List)-1)
	hasString := false

	for _, e := range expr.(*ast.ListExpr).List[1:] {
		result, err := i.eval(e, env)
		if err != nil {
			return nil, err
		}

		if result.Type() != runtime.NumberType && result.Type() != runtime.StringType {
			return nil, i.error(fmt.Sprintf("invalid type %s for `%s`", result.Type(), operator), e.Location())
		} else if result.Type() == runtime.StringType {
			hasString = true
		}

		results = append(results, result)
	}

	if hasString {
		var out strings.Builder

		for _, r := range results {
			out.WriteString(fmt.Sprintf("%v", r))
		}

		return runtime.NewString(out.String()), nil
	}

	var total float64

	for _, r := range results {
		total += r.(runtime.Number).Value
	}

	return runtime.NewNumber(total), nil
}

// evalMathSymbol evaluates mathematical operators.
func (i *Interpreter) evalMathSymbol(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	operator := expr.(*ast.ListExpr).List[0].(*ast.SymbolExpr).Symbol

	results := make([]runtime.Value, 0, len(expr.(*ast.ListExpr).List)-1)

	for _, e := range expr.(*ast.ListExpr).List[1:] {
		result, err := i.eval(e, env)
		if err != nil {
			return nil, err
		}

		if result.Type() != runtime.NumberType {
			return nil, i.error(fmt.Sprintf("invalid type %s for `%s`", result.Type(), operator), e.Location())
		}

		results = append(results, result)
	}

	total := results[0].(runtime.Number).Value

	if len(results) == 1 {
		if operator != "-" {
			return nil, i.error("invalid operand length", expr.Location())
		}

		return runtime.NewNumber(-total), nil
	}

	for _, r := range results[1:] {
		value := r.(runtime.Number).Value

		switch operator {
		case "-":
			total -= value
		case "*":
			total *= value
		case "/":
			if value == 0 {
				return nil, i.error("division by zero", expr.Location())
			}

			total /= value
		}
	}

	return runtime.NewNumber(total), nil
}

// evalLogicalSymbol evaluates logical and comparison operators.
func (i *Interpreter) evalLogicalSymbol(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	operator := expr.(*ast.ListExpr).List[0].(*ast.SymbolExpr).Symbol

	// = => same type (any) between 2 expressions
	if operator == "=" {
		resultLeft, err := i.eval(expr.(*ast.ListExpr).List[1], env)
		if err != nil {
			return nil, err
		}

		resultRight, err := i.eval(expr.(*ast.ListExpr).List[2], env)
		if err != nil {
			return nil, err
		}

		if resultLeft.Type() != resultRight.Type() {
			return nil, i.error(fmt.Sprintf("cannot apply %s operator for %s and %s expressiones", operator, resultLeft.Type(), resultRight.Type()), expr.Location())
		}

		if resultLeft.Type() == runtime.NumberType {
			return runtime.NewBool(resultLeft.(runtime.Number).Value == resultRight.(runtime.Number).Value), nil
		}

		if resultLeft.Type() == runtime.StringType {
			return runtime.NewBool(resultLeft.(runtime.String).Value == resultRight.(runtime.String).Value), nil
		}

		if resultLeft.Type() == runtime.BoolType {
			return runtime.NewBool(resultLeft.(runtime.Bool).Value == resultRight.(runtime.Bool).Value), nil
		}

		if resultLeft.Type() == runtime.NilType {
			return runtime.NewBool(true), nil
		}

		return nil, i.error(fmt.Sprintf("invalid type %s for `%s`", resultLeft.Type(), operator), expr.Location())
	}

	// > >= < <= => same type (string or number) between 2 expressions
	if operator == "<" || operator == "<=" || operator == ">" || operator == ">=" {
		resultLeft, err := i.eval(expr.(*ast.ListExpr).List[1], env)
		if err != nil {
			return nil, err
		}

		resultRight, err := i.eval(expr.(*ast.ListExpr).List[2], env)
		if err != nil {
			return nil, err
		}

		if resultLeft.Type() != resultRight.Type() {
			return nil, i.error(fmt.Sprintf("cannot apply %s operator for %s and %s expressions", operator, resultLeft.Type(), resultRight.Type()), expr.Location())
		}

		if resultLeft.Type() == runtime.NumberType {
			switch operator {
			case "<":
				return runtime.NewBool(resultLeft.(runtime.Number).Value < resultRight.(runtime.Number).Value), nil
			case "<=":
				return runtime.NewBool(resultLeft.(runtime.Number).Value <= resultRight.(runtime.Number).Value), nil
			case ">":
				return runtime.NewBool(resultLeft.(runtime.Number).Value > resultRight.(runtime.Number).Value), nil
			case ">=":
				return runtime.NewBool(resultLeft.(runtime.Number).Value >= resultRight.(runtime.Number).Value), nil
			}
		}

		if resultLeft.Type() == runtime.StringType {
			switch operator {
			case "<":
				return runtime.NewBool(resultLeft.(runtime.String).Value < resultRight.(runtime.String).Value), nil
			case "<=":
				return runtime.NewBool(resultLeft.(runtime.String).Value <= resultRight.(runtime.String).Value), nil
			case ">":
				return runtime.NewBool(resultLeft.(runtime.String).Value > resultRight.(runtime.String).Value), nil
			case ">=":
				return runtime.NewBool(resultLeft.(runtime.String).Value >= resultRight.(runtime.String).Value), nil
			}
		}

		return nil, i.error(fmt.Sprintf("invalid type %s for `%s`", resultLeft.Type(), operator), expr.Location())
	}

	// and or => only booleans between multiple expressions
	if operator == "and" || operator == "or" {
		results := make([]runtime.Value, 0, len(expr.(*ast.ListExpr).List)-1)

		for _, e := range expr.(*ast.ListExpr).List[1:] {
			result, err := i.eval(e, env)
			if err != nil {
				return nil, err
			}

			if result.Type() != runtime.BoolType {
				return nil, i.error(fmt.Sprintf("invalid type %s for `%s`", result.Type(), operator), e.Location())
			}

			results = append(results, result)
		}

		logical := results[0].(runtime.Bool).Value

		for _, r := range results[1:] {
			value := r.(runtime.Bool).Value

			switch operator {
			case "and":
				logical = logical && value
			case "or":
				logical = logical || value
			}
		}

		return runtime.NewBool(logical), nil
	}

	return nil, i.error(fmt.Sprintf("unknown operator `%s`", operator), expr.Location())
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
	alternate := exprList.List[3]

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

	return i.evalInTailPosition(alternate, env)
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

// evalPrint evaluates a `print` expression.
func (i *Interpreter) evalPrint(expr ast.SExpr, env *runtime.Environment) (runtime.Value, error) {
	exprList := expr.(*ast.ListExpr)

	for _, e := range exprList.List[1:] {
		result, err := i.eval(e, env)
		if err != nil {
			return nil, err
		}

		fmt.Print(result)
	}

	fmt.Println()

	return runtime.NewNil(), nil
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

		var key string

		if keyExpr.Kind() == ast.SymbolKind {
			key = keyExpr.(*ast.SymbolExpr).Symbol
		} else if keyExpr.Kind() == ast.StringKind {
			key = keyExpr.(*ast.StringExpr).String
		}

		result, err := i.eval(valueExpr, env)
		if err != nil {
			return nil, err
		}

		elements[key] = result
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

	if funcValue.Type() != runtime.CoreFuncType && funcValue.Type() != runtime.FuncType {
		return nil, i.error("expression is not a function", exprList.List[0].Location())
	}

	valArgs, err := i.evalFunctionArguments(exprList.List[1:], env)
	if err != nil {
		return nil, err
	}

	// native core function
	if funcValue.Type() == runtime.CoreFuncType {
		return funcValue.(runtime.CoreFunction).Value(valArgs...)
	}

	// lambda function with TCO
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
