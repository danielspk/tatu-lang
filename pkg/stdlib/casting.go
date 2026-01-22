package stdlib

import (
	"fmt"
	"strconv"

	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterCasting registers casting core functions in the environment.
func RegisterCasting(env *runtime.Environment) error {
	functions := map[string]runtime.CoreFunction{
		"to-string": runtime.NewCoreFunction(toString),
		"to-number": runtime.NewCoreFunction(toNumber),
		"to-bool":   runtime.NewCoreFunction(toBool),
	}

	for name, fn := range functions {
		if _, err := env.Define(name, fn); err != nil {
			return fmt.Errorf("failed to register casting function `%s`: %w", name, err)
		}
	}

	return nil
}

// toString implements the to-string conversion core function.
// Usage: (to-string 42) => "42"
func toString(args ...runtime.Value) (runtime.Value, error) {
	const name = "to-string"

	if err := expectArgs(name, 1, args); err != nil {
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

// toNumber implements the to-number conversion core function.
// Usage: (to-number "42") => 42
func toNumber(args ...runtime.Value) (runtime.Value, error) {
	const name = "to-number"

	if err := expectArgs(name, 1, args); err != nil {
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

// toBool implements the to-bool conversion core function.
// Usage: (to-bool 0) => false
func toBool(args ...runtime.Value) (runtime.Value, error) {
	const name = "to-bool"

	if err := expectArgs(name, 1, args); err != nil {
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
