// Package stdlib implements standard library functions.
package stdlib

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/danielspk/tatu-lang/pkg/core"
	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterMath registers mathematical functions.
func RegisterMath(natives map[string]runtime.NativeFunction) {
	natives["math:pi"] = runtime.NewNativeFunction(mathPi)
	natives["math:e"] = runtime.NewNativeFunction(mathE)
	natives["math:abs"] = runtime.NewNativeFunction(mathAbs)
	natives["math:floor"] = runtime.NewNativeFunction(mathFloor)
	natives["math:ceil"] = runtime.NewNativeFunction(mathCeil)
	natives["math:round"] = runtime.NewNativeFunction(mathRound)
	natives["math:sin"] = runtime.NewNativeFunction(mathSin)
	natives["math:cos"] = runtime.NewNativeFunction(mathCos)
	natives["math:tan"] = runtime.NewNativeFunction(mathTan)
	natives["math:min"] = runtime.NewNativeFunction(mathMin)
	natives["math:max"] = runtime.NewNativeFunction(mathMax)
	natives["math:sqrt"] = runtime.NewNativeFunction(mathSqrt)
	natives["math:pow"] = runtime.NewNativeFunction(mathPow)
	natives["math:log"] = runtime.NewNativeFunction(mathLog)
	natives["math:exp"] = runtime.NewNativeFunction(mathExp)
	natives["math:between"] = runtime.NewNativeFunction(mathBetween)
	natives["math:rand"] = runtime.NewNativeFunction(mathRand)
}

// mathPi implements the pi constant.
// Usage: (math:pi) => 3.1415926536
func mathPi(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:pi"

	if err := core.ExpectArgs(name, 0, args); err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Pi), nil
}

// mathE implements the e constant.
// Usage: (math:e) => 2.7182818285
func mathE(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:e"

	if err := core.ExpectArgs(name, 0, args); err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.E), nil
}

// mathAbs implements the absolute value function.
// Usage: (math:abs -5) => 5
func mathAbs(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:abs"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Abs(num.Value)), nil
}

// mathFloor implements the floor function.
// Usage: (math:floor 3.7) => 3
func mathFloor(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:floor"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Floor(num.Value)), nil
}

// mathCeil implements the ceiling function.
// Usage: (math:ceil 3.2) => 4
func mathCeil(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:ceil"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Ceil(num.Value)), nil
}

// mathRound implements the rounding function.
// Usage: (math:round 3.5) => 4
func mathRound(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:round"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Round(num.Value)), nil
}

// mathSin implements the sine function.
// Usage: (math:sin 0) => 0
func mathSin(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:sin"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Sin(num.Value)), nil
}

// mathCos implements the cosine function.
// Usage: (math:cos 0) => 1
func mathCos(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:cos"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Cos(num.Value)), nil
}

// mathTan implements the tangent function.
// Usage: (math:tan 0) => 0
func mathTan(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:tan"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Tan(num.Value)), nil
}

// mathMin implements the minimum function.
// Usage: (math:min 3 5) => 3
func mathMin(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:min"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	a, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	b, err := core.ExpectNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Min(a.Value, b.Value)), nil
}

// mathMax implements the maximum function.
// Usage: (math:max 3 5) => 5
func mathMax(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:max"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	a, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	b, err := core.ExpectNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Max(a.Value, b.Value)), nil
}

// mathSqrt implements the square function.
// Usage: (math:sqrt 16) => 4
func mathSqrt(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:sqrt"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := core.ExpectNumber(name, 0, args[0])
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

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	base, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	exponent, err := core.ExpectNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Pow(base.Value, exponent.Value)), nil
}

// mathLog implements the natural logarithm function.
// Usage: (math:log 2.718281828459045) => 1
func mathLog(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:log"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := core.ExpectNumber(name, 0, args[0])
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

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	num, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(math.Exp(num.Value)), nil
}

// mathBetween checks if a value is between min and max (inclusive).
// Usage: (math:between 5 1 10) => true
func mathBetween(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:between"

	if err := core.ExpectArgs(name, 3, args); err != nil {
		return nil, err
	}

	value, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	min, err := core.ExpectNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	max, err := core.ExpectNumber(name, 2, args[2])
	if err != nil {
		return nil, err
	}

	result := value.Value >= min.Value && value.Value <= max.Value

	return runtime.NewBool(result), nil
}

// mathRand generates a random integer between min and max (inclusive).
// Usage: (math:rand 1 10) => 7
func mathRand(args ...runtime.Value) (runtime.Value, error) {
	const name = "math:rand"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	minNum, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	maxNum, err := core.ExpectNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	minInt := int(math.Floor(minNum.Value))
	maxInt := int(math.Floor(maxNum.Value))

	if minInt > maxInt {
		return nil, fmt.Errorf("`%s` min (%d) cannot be greater than max (%d)", name, minInt, maxInt)
	}

	result := minInt + rand.Intn(maxInt-minInt+1)

	return runtime.NewNumber(float64(result)), nil
}
