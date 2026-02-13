package builtins

import (
	"fmt"

	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterComparison registers comparison and not operator natives.
func RegisterComparison(natives map[string]runtime.NativeFunction) {
	natives["="] = runtime.NewNativeFunction(equal)
	natives[">"] = runtime.NewNativeFunction(greaterThan)
	natives[">="] = runtime.NewNativeFunction(greaterThanOrEqual)
	natives["<"] = runtime.NewNativeFunction(lessThan)
	natives["<="] = runtime.NewNativeFunction(lessThanOrEqual)
	natives["not"] = runtime.NewNativeFunction(not)
}

// equal implements the = operator (equality between same types).
// Usage: (= 1 1) => true
// Usage: (= "a" "b") => false
func equal(args ...runtime.Value) (runtime.Value, error) {
	const name = "="

	if len(args) != 2 {
		return nil, fmt.Errorf("`%s` expected 2 arguments, got %d", name, len(args))
	}

	left, right := args[0], args[1]

	if left.Type() != right.Type() {
		return nil, fmt.Errorf("`%s` cannot compare %s and %s", name, left.Type(), right.Type())
	}

	switch left.Type() {
	case runtime.NumberType:
		return runtime.NewBool(left.(runtime.Number).Value == right.(runtime.Number).Value), nil
	case runtime.StringType:
		return runtime.NewBool(left.(runtime.String).Value == right.(runtime.String).Value), nil
	case runtime.BoolType:
		return runtime.NewBool(left.(runtime.Bool).Value == right.(runtime.Bool).Value), nil
	case runtime.NilType:
		return runtime.NewBool(true), nil
	default:
		return nil, fmt.Errorf("`%s` invalid type %s", name, left.Type())
	}
}

// compareOrdered evaluates an ordering comparison between two values of the same type.
// Returns -1, 0, or 1 for less, equal, or greater respectively.
func compareOrdered(name string, args []runtime.Value) (int, error) {
	if len(args) != 2 {
		return 0, fmt.Errorf("`%s` expected 2 arguments, got %d", name, len(args))
	}

	left, right := args[0], args[1]

	if left.Type() != right.Type() {
		return 0, fmt.Errorf("`%s` cannot compare %s and %s", name, left.Type(), right.Type())
	}

	switch left.Type() {
	case runtime.NumberType:
		l, r := left.(runtime.Number).Value, right.(runtime.Number).Value
		if l < r {
			return -1, nil
		}
		if l > r {
			return 1, nil
		}
		return 0, nil
	case runtime.StringType:
		l, r := left.(runtime.String).Value, right.(runtime.String).Value
		if l < r {
			return -1, nil
		}
		if l > r {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("`%s` invalid type %s", name, left.Type())
	}
}

// greaterThan implements the > operator.
// Usage: (> 3 1) => true
func greaterThan(args ...runtime.Value) (runtime.Value, error) {
	cmp, err := compareOrdered(">", args)
	if err != nil {
		return nil, err
	}

	return runtime.NewBool(cmp > 0), nil
}

// greaterThanOrEqual implements the >= operator.
// Usage: (>= 3 3) => true
func greaterThanOrEqual(args ...runtime.Value) (runtime.Value, error) {
	cmp, err := compareOrdered(">=", args)
	if err != nil {
		return nil, err
	}

	return runtime.NewBool(cmp >= 0), nil
}

// lessThan implements the < operator.
// Usage: (< 1 3) => true
func lessThan(args ...runtime.Value) (runtime.Value, error) {
	cmp, err := compareOrdered("<", args)
	if err != nil {
		return nil, err
	}

	return runtime.NewBool(cmp < 0), nil
}

// lessThanOrEqual implements the <= operator.
// Usage: (<= 3 3) => true
func lessThanOrEqual(args ...runtime.Value) (runtime.Value, error) {
	cmp, err := compareOrdered("<=", args)
	if err != nil {
		return nil, err
	}

	return runtime.NewBool(cmp <= 0), nil
}

// not implements the not operator (boolean negation).
// Usage: (not true) => false
func not(args ...runtime.Value) (runtime.Value, error) {
	const name = "not"

	if len(args) != 1 {
		return nil, fmt.Errorf("`%s` expected 1 argument, got %d", name, len(args))
	}

	if args[0].Type() != runtime.BoolType {
		return nil, fmt.Errorf("`%s` invalid type %s", name, args[0].Type())
	}

	return runtime.NewBool(!args[0].(runtime.Bool).Value), nil
}
