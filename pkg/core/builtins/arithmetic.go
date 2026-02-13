package builtins

import (
	"fmt"
	"math"
	"strings"

	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterArithmetic registers arithmetic operator natives.
func RegisterArithmetic(natives map[string]runtime.NativeFunction) {
	natives["+"] = runtime.NewNativeFunction(add)
	natives["-"] = runtime.NewNativeFunction(subtract)
	natives["*"] = runtime.NewNativeFunction(multiply)
	natives["/"] = runtime.NewNativeFunction(divide)
	natives["%"] = runtime.NewNativeFunction(modulo)
}

// add implements the + operator (addition and string concatenation).
// Usage: (+ 1 2 3) => 6
// Usage: (+ "hello" " " "world") => "hello world"
func add(args ...runtime.Value) (runtime.Value, error) {
	const name = "+"

	if len(args) < 2 {
		return nil, fmt.Errorf("`%s` expected at least 2 arguments, got %d", name, len(args))
	}

	hasString := false

	for _, arg := range args {
		if arg.Type() != runtime.NumberType && arg.Type() != runtime.StringType {
			return nil, fmt.Errorf("`%s` invalid type %s", name, arg.Type())
		}

		if arg.Type() == runtime.StringType {
			hasString = true
		}
	}

	if hasString {
		var out strings.Builder

		for _, arg := range args {
			out.WriteString(fmt.Sprintf("%v", arg))
		}

		return runtime.NewString(out.String()), nil
	}

	var total float64

	for _, arg := range args {
		total += arg.(runtime.Number).Value
	}

	return runtime.NewNumber(total), nil
}

// subtract implements the - operator (subtraction and unary negation).
// Usage: (- 10 3) => 7
// Usage: (- 5) => -5
func subtract(args ...runtime.Value) (runtime.Value, error) {
	const name = "-"

	if len(args) < 1 {
		return nil, fmt.Errorf("`%s` expected at least 1 argument, got 0", name)
	}

	for _, arg := range args {
		if arg.Type() != runtime.NumberType {
			return nil, fmt.Errorf("`%s` invalid type %s", name, arg.Type())
		}
	}

	total := args[0].(runtime.Number).Value

	if len(args) == 1 {
		return runtime.NewNumber(-total), nil
	}

	for _, arg := range args[1:] {
		total -= arg.(runtime.Number).Value
	}

	return runtime.NewNumber(total), nil
}

// multiply implements the * operator.
// Usage: (* 2 3 4) => 24
func multiply(args ...runtime.Value) (runtime.Value, error) {
	const name = "*"

	if len(args) < 2 {
		return nil, fmt.Errorf("`%s` expected at least 2 arguments, got %d", name, len(args))
	}

	for _, arg := range args {
		if arg.Type() != runtime.NumberType {
			return nil, fmt.Errorf("`%s` invalid type %s", name, arg.Type())
		}
	}

	total := args[0].(runtime.Number).Value

	for _, arg := range args[1:] {
		total *= arg.(runtime.Number).Value
	}

	return runtime.NewNumber(total), nil
}

// divide implements the / operator.
// Usage: (/ 10 2) => 5
func divide(args ...runtime.Value) (runtime.Value, error) {
	const name = "/"

	if len(args) < 2 {
		return nil, fmt.Errorf("`%s` expected at least 2 arguments, got %d", name, len(args))
	}

	for _, arg := range args {
		if arg.Type() != runtime.NumberType {
			return nil, fmt.Errorf("`%s` invalid type %s", name, arg.Type())
		}
	}

	total := args[0].(runtime.Number).Value

	for _, arg := range args[1:] {
		value := arg.(runtime.Number).Value

		if value == 0 {
			return nil, fmt.Errorf("`%s` division by zero", name)
		}

		total /= value
	}

	return runtime.NewNumber(total), nil
}

// modulo implements the % operator.
// Usage: (% 10 3) => 1
func modulo(args ...runtime.Value) (runtime.Value, error) {
	const name = "%"

	if len(args) != 2 {
		return nil, fmt.Errorf("`%s` expected exactly 2 arguments, got %d", name, len(args))
	}

	for _, arg := range args {
		if arg.Type() != runtime.NumberType {
			return nil, fmt.Errorf("`%s` invalid type %s", name, arg.Type())
		}
	}

	left := args[0].(runtime.Number).Value
	right := args[1].(runtime.Number).Value

	if right == 0 {
		return nil, fmt.Errorf("`%s` modulo by zero", name)
	}

	return runtime.NewNumber(math.Mod(left, right)), nil
}
