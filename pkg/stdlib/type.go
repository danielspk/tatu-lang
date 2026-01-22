package stdlib

import (
	"fmt"
	"math"

	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterTypes registers types core functions in the environment.
func RegisterTypes(env *runtime.Environment) error {
	functions := map[string]runtime.CoreFunction{
		"is-bool":     runtime.NewCoreFunction(isBool),
		"is-number":   runtime.NewCoreFunction(isNumber),
		"is-int":      runtime.NewCoreFunction(isInt),
		"is-string":   runtime.NewCoreFunction(isString),
		"is-vector":   runtime.NewCoreFunction(isVector),
		"is-map":      runtime.NewCoreFunction(isMap),
		"is-nil":      runtime.NewCoreFunction(isNil),
		"is-function": runtime.NewCoreFunction(isFunction),
	}

	for name, fn := range functions {
		if _, err := env.Define(name, fn); err != nil {
			return fmt.Errorf("failed to register type function `%s`: %w", name, err)
		}
	}

	return nil
}

// isBool implements the boolean type checking core function.
// Usage: (is-bool true) => true
func isBool(args ...runtime.Value) (runtime.Value, error) {
	const name = "is-bool"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	return runtime.NewBool(args[0].Type() == runtime.BoolType), nil
}

// isNumber implements the number type checking core function.
// Usage: (is-number 42) => true
func isNumber(args ...runtime.Value) (runtime.Value, error) {
	const name = "is-number"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	return runtime.NewBool(args[0].Type() == runtime.NumberType), nil
}

// isInt implements the integer type checking core function.
// Usage: (is-int 42) => true
func isInt(args ...runtime.Value) (runtime.Value, error) {
	const name = "is-int"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	if args[0].Type() != runtime.NumberType {
		return runtime.NewBool(false), nil
	}

	num := args[0].(runtime.Number)
	return runtime.NewBool(num.Value == math.Trunc(num.Value)), nil
}

// isString implements the string type checking core function.
// Usage: (is-string "hello") => true
func isString(args ...runtime.Value) (runtime.Value, error) {
	const name = "is-string"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	return runtime.NewBool(args[0].Type() == runtime.StringType), nil
}

// isVector implements the vector type checking core function.
// Usage: (is-vector (vector 1 2 3)) => true
func isVector(args ...runtime.Value) (runtime.Value, error) {
	const name = "is-vector"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	return runtime.NewBool(args[0].Type() == runtime.VectorType), nil
}

// isMap implements the map type checking core function.
// Usage: (is-map (map "key" "value")) => true
func isMap(args ...runtime.Value) (runtime.Value, error) {
	const name = "is-map"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	return runtime.NewBool(args[0].Type() == runtime.MapType), nil
}

// isNil implements the nil type checking core function.
// Usage: (is-nil nil) => true
func isNil(args ...runtime.Value) (runtime.Value, error) {
	const name = "is-nil"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	return runtime.NewBool(args[0].Type() == runtime.NilType), nil
}

// isFunction implements the function type checking core function.
// Usage: (is-function (lambda (x) x)) => true
func isFunction(args ...runtime.Value) (runtime.Value, error) {
	const name = "is-function"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	typ := args[0].Type()
	return runtime.NewBool(typ == runtime.FuncType || typ == runtime.CoreFuncType), nil
}
