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
		"math:pi":    runtime.NewCoreFunction(mathPi),
		"math:e":     runtime.NewCoreFunction(mathE),
		"math:abs":   runtime.NewCoreFunction(mathAbs),
		"math:floor": runtime.NewCoreFunction(mathFloor),
		"math:ceil":  runtime.NewCoreFunction(mathCeil),
		"math:round": runtime.NewCoreFunction(mathRound),
		"math:sin":   runtime.NewCoreFunction(mathSin),
		"math:cos":   runtime.NewCoreFunction(mathCos),
		"math:tan":   runtime.NewCoreFunction(mathTan),
		"math:min":   runtime.NewCoreFunction(mathMin),
		"math:max":   runtime.NewCoreFunction(mathMax),
		"math:sqrt":  runtime.NewCoreFunction(mathSqrt),
		"math:pow":   runtime.NewCoreFunction(mathPow),
		"math:log":   runtime.NewCoreFunction(mathLog),
		"math:exp":   runtime.NewCoreFunction(mathExp),
	}

	for name, fn := range functions {
		if _, err := env.Define(name, fn); err != nil {
			return fmt.Errorf("failed to register math function `%s`: %v", name, err)
		}
	}

	return nil
}

// mathPi implements the pi constant.
// Usage: (math:pi) => 3.1415926536
func mathPi(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:pi"

	if err := expectArgs(name, 0, args); err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Pi), nil
}

// mathE implements the e constant.
// Usage: (math:e) => 2.7182818285
func mathE(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:e"

	if err := expectArgs(name, 0, args); err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.E), nil
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

// mathFloor implements the floor function.
// Usage: (math:floor 3.7) => 3
func mathFloor(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:floor"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Floor(num.Value)), nil
}

// mathCeil implements the ceiling function.
// Usage: (math:ceil 3.2) => 4
func mathCeil(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:ceil"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Ceil(num.Value)), nil
}

// mathRound implements the rounding function.
// Usage: (math:round 3.5) => 4
func mathRound(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:round"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Round(num.Value)), nil
}

// mathSin implements the sine function.
// Usage: (math:sin 0) => 0
func mathSin(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:sin"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Sin(num.Value)), nil
}

// mathCos implements the cosine function.
// Usage: (math:cos 0) => 1
func mathCos(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:cos"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Cos(num.Value)), nil
}

// mathTan implements the tangent function.
// Usage: (math:tan 0) => 0
func mathTan(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:tan"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Tan(num.Value)), nil
}

// mathMin implements the minimum function.
// Usage: (math:min 3 5) => 3
func mathMin(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:min"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	a, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	b, err := expectNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Min(a.Value, b.Value)), nil
}

// mathMax implements the maximum function.
// Usage: (math:max 3 5) => 5
func mathMax(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:max"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	a, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	b, err := expectNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Max(a.Value, b.Value)), nil
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

// mathLog implements the natural logarithm function.
// Usage: (math:log 2.718281828459045) => 1
func mathLog(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:log"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	if num.Value <= 0 {
		return nil, fmt.Errorf("`%s` requires a positive number", name)
	}

	return runtime.NewNumber(math.Log(num.Value)), nil
}

// mathExp implements the exponential function.
// Usage: (math:exp 1) => 2.718281828459045
func mathExp(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:exp"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Exp(num.Value)), nil
}
