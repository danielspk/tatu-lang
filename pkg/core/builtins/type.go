package builtins

import (
	"fmt"
	"math"
	"strconv"

	"github.com/danielspk/tatu-lang/pkg/core"
	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterTypes registers type checking and conversion functions.
func RegisterTypes(natives map[string]runtime.NativeFunction) {
	natives["is-bool"] = runtime.NewNativeFunction(isBool)
	natives["is-number"] = runtime.NewNativeFunction(isNumber)
	natives["is-int"] = runtime.NewNativeFunction(isInt)
	natives["is-string"] = runtime.NewNativeFunction(isString)
	natives["is-vector"] = runtime.NewNativeFunction(isVector)
	natives["is-map"] = runtime.NewNativeFunction(isMap)
	natives["is-nil"] = runtime.NewNativeFunction(isNil)
	natives["is-function"] = runtime.NewNativeFunction(isFunction)
	natives["to-string"] = runtime.NewNativeFunction(toString)
	natives["to-number"] = runtime.NewNativeFunction(toNumber)
	natives["to-bool"] = runtime.NewNativeFunction(toBool)
}

// isBool implements the boolean type checking function.
// Usage: (is-bool true) => true
func isBool(args ...runtime.Value) (runtime.Value, error) {
	const name = "is-bool"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	return runtime.NewBool(args[0].Type() == runtime.BoolType), nil
}

// isNumber implements the number type checking function.
// Usage: (is-number 42) => true
func isNumber(args ...runtime.Value) (runtime.Value, error) {
	const name = "is-number"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	return runtime.NewBool(args[0].Type() == runtime.NumberType), nil
}

// isInt implements the integer type checking function.
// Usage: (is-int 42) => true
func isInt(args ...runtime.Value) (runtime.Value, error) {
	const name = "is-int"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	if args[0].Type() != runtime.NumberType {
		return runtime.NewBool(false), nil
	}

	num := args[0].(runtime.Number)
	return runtime.NewBool(num.Value == math.Trunc(num.Value)), nil
}

// isString implements the string type checking function.
// Usage: (is-string "hello") => true
func isString(args ...runtime.Value) (runtime.Value, error) {
	const name = "is-string"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	return runtime.NewBool(args[0].Type() == runtime.StringType), nil
}

// isVector implements the vector type checking function.
// Usage: (is-vector (vector 1 2 3)) => true
func isVector(args ...runtime.Value) (runtime.Value, error) {
	const name = "is-vector"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	return runtime.NewBool(args[0].Type() == runtime.VectorType), nil
}

// isMap implements the map type checking function.
// Usage: (is-map (map "key" "value")) => true
func isMap(args ...runtime.Value) (runtime.Value, error) {
	const name = "is-map"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	return runtime.NewBool(args[0].Type() == runtime.MapType), nil
}

// isNil implements the nil type checking function.
// Usage: (is-nil nil) => true
func isNil(args ...runtime.Value) (runtime.Value, error) {
	const name = "is-nil"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	return runtime.NewBool(args[0].Type() == runtime.NilType), nil
}

// isFunction implements the function type checking function.
// Usage: (is-function (lambda (x) x)) => true
func isFunction(args ...runtime.Value) (runtime.Value, error) {
	const name = "is-function"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	typ := args[0].Type()
	return runtime.NewBool(typ == runtime.FuncType || typ == runtime.NativeFuncType), nil
}

// toString implements the to-string conversion function.
// Usage: (to-string 42) => "42"
func toString(args ...runtime.Value) (runtime.Value, error) {
	const name = "to-string"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	switch args[0].Type() {
	case runtime.StringType:
		return args[0], nil
	case runtime.NumberType:
		num := args[0].(runtime.Number)
		return runtime.NewString(num.String()), nil
	case runtime.BoolType:
		b := args[0].(runtime.Bool)
		return runtime.NewString(b.String()), nil
	case runtime.NilType:
		n := args[0].(runtime.Nil)
		return runtime.NewString(n.String()), nil
	default:
		return nil, fmt.Errorf("`%s` cannot convert %s to STRING", name, args[0].Type())
	}
}

// toNumber implements the to-number conversion function.
// Usage: (to-number "42") => 42
func toNumber(args ...runtime.Value) (runtime.Value, error) {
	const name = "to-number"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	switch args[0].Type() {
	case runtime.NumberType:
		return args[0], nil
	case runtime.StringType:
		str := args[0].(runtime.String)
		num, err := strconv.ParseFloat(str.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("`%s` cannot parse STRING '%s' to NUMBER: %w", name, str.Value, err)
		}
		return runtime.NewNumber(num), nil
	case runtime.BoolType:
		b := args[0].(runtime.Bool)
		if b.Value {
			return runtime.NewNumber(1), nil
		}
		return runtime.NewNumber(0), nil
	case runtime.NilType:
		return runtime.NewNumber(0), nil
	default:
		return nil, fmt.Errorf("`%s` cannot convert %s to NUMBER", name, args[0].Type())
	}
}

// toBool implements the to-bool conversion function.
// Usage: (to-bool 0) => false
func toBool(args ...runtime.Value) (runtime.Value, error) {
	const name = "to-bool"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	switch args[0].Type() {
	case runtime.BoolType:
		return args[0], nil
	case runtime.NumberType:
		num := args[0].(runtime.Number)
		return runtime.NewBool(num.Value != 0), nil
	case runtime.StringType:
		str := args[0].(runtime.String)
		return runtime.NewBool(str.Value != ""), nil
	case runtime.NilType:
		return runtime.NewBool(false), nil
	default:
		return nil, fmt.Errorf("`%s` cannot convert %s to BOOL", name, args[0].Type())
	}
}
