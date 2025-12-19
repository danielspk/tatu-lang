// Package stdlib implements standard library core functions.
package stdlib

import (
	"fmt"
	"math"

	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterMath registers mathematical core functions in the environment.
func RegisterMath(env *runtime.Environment) error {
	functions := map[string]runtime.CoreFunction{
		"math:sqrt": runtime.NewCoreFunction(mathSqrt),
		"math:abs":  runtime.NewCoreFunction(mathAbs),
		"math:pow":  runtime.NewCoreFunction(mathPow),
	}

	for name, fn := range functions {
		if _, err := env.Define(name, fn); err != nil {
			return fmt.Errorf("failed to register math function `%s`: %v", name, err)
		}
	}

	return nil
}

// mathSqrt implements the square core function.
// Usage: (math:sqrt 16) => 4
func mathSqrt(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:sqrt"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	if num.Value < 0 {
		return nil, fmt.Errorf("`%s` cannot compute a negative number", name)
	}

	return runtime.NewNumber(math.Sqrt(num.Value)), nil
}

// mathAbs implements the absolute value function.
// Usage: (math:abs -5) => 5
func mathAbs(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:abs"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Abs(num.Value)), nil
}

// mathPow implements the power function.
// Usage: (math:pow 2 3) => 8
func mathPow(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:pow"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	base, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	exponent, err := expectNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Pow(base.Value, exponent.Value)), nil
}
