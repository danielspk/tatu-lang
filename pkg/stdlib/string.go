package stdlib

import (
	"fmt"
	"strings"

	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterString registers string core functions in the environment.
func RegisterString(env *runtime.Environment) error {
	functions := map[string]runtime.CoreFunction{
		"str:len":      runtime.NewCoreFunction(stringLen),
		"str:contains": runtime.NewCoreFunction(stringContains),
		"str:index":    runtime.NewCoreFunction(stringIndex),
		"str:upper":    runtime.NewCoreFunction(stringUpper),
		"str:lower":    runtime.NewCoreFunction(stringLower),
		"str:trim":     runtime.NewCoreFunction(stringTrim),
		"str:slice":    runtime.NewCoreFunction(stringSlice),
		"str:split":    runtime.NewCoreFunction(stringSplit),
		"str:join":     runtime.NewCoreFunction(stringJoin),
		"str:replace":  runtime.NewCoreFunction(stringReplace),
		"str:starts":   runtime.NewCoreFunction(stringStarts),
		"str:ends":     runtime.NewCoreFunction(stringEnds),
		"str:reverse":  runtime.NewCoreFunction(stringReverse),
		"str:repeat":   runtime.NewCoreFunction(stringRepeat),
		"str:concat":   runtime.NewCoreFunction(stringConcat),
	}

	for name, fn := range functions {
		if _, err := env.Define(name, fn); err != nil {
			return fmt.Errorf("failed to register string function `%s`: %w", name, err)
		}
	}

	return nil
}

// stringLen implements the string length core function.
// Usage: (str:len "hello") => 5
func stringLen(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:len"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	str, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(float64(len([]rune(str.Value)))), nil
}

// stringContains implements the string contains check core function.
// Usage: (str:contains "hello world" "world") => true
func stringContains(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:contains"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	str, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	substr, err := expectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	return runtime.NewBool(strings.Contains(str.Value, substr.Value)), nil
}

// stringIndex implements the string index search core function.
// Usage: (str:index "hello" "ll") => 2
func stringIndex(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:index"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	str, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	substr, err := expectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	runes := []rune(str.Value)
	subRunes := []rune(substr.Value)
	index := -1

	for i := 0; i <= len(runes)-len(subRunes); i++ {
		match := true

		for j := 0; j < len(subRunes); j++ {
			if runes[i+j] != subRunes[j] {
				match = false
				break
			}
		}

		if match {
			index = i
			break
		}
	}

	return runtime.NewNumber(float64(index)), nil
}

// stringUpper implements the string uppercase conversion core function.
// Usage: (str:upper "hello") => "HELLO"
func stringUpper(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:upper"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	str, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewString(strings.ToUpper(str.Value)), nil
}

// stringLower implements the string lowercase conversion core function.
// Usage: (str:lower "HELLO") => "hello"
func stringLower(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:lower"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	str, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewString(strings.ToLower(str.Value)), nil
}

// stringTrim implements the string whitespace trimming core function.
// Usage: (str:trim "  hello  ") => "hello"
func stringTrim(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:trim"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	str, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewString(strings.TrimSpace(str.Value)), nil
}

// stringSlice implements the string slice extraction core function.
// Usage: (str:slice "hello" 1 4) => "ell"
func stringSlice(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:slice"

	if err := expectArgs(name, 3, args); err != nil {
		return nil, err
	}

	str, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	startNum, err := expectIntegerNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	endNum, err := expectIntegerNumber(name, 2, args[2])
	if err != nil {
		return nil, err
	}

	runes := []rune(str.Value)
	start := int(startNum.Value)
	end := int(endNum.Value)

	if start < 0 || start > len(runes) {
		return nil, fmt.Errorf("`%s` start index out of bounds: %d (string length: %d)", name, start, len(runes))
	}

	if end < 0 || end > len(runes) {
		return nil, fmt.Errorf("`%s` end index out of bounds: %d (string length: %d)", name, end, len(runes))
	}

	if start > end {
		return nil, fmt.Errorf("`%s` start index (%d) cannot be greater than end index (%d)", name, start, end)
	}

	return runtime.NewString(string(runes[start:end])), nil
}

// stringSplit implements the string split core function.
// Usage: (str:split "a,b,c" ",") => ("a" "b" "c")
func stringSplit(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:split"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	str, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	sep, err := expectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	parts := strings.Split(str.Value, sep.Value)
	elements := make([]runtime.Value, len(parts))

	for i, part := range parts {
		elements[i] = runtime.NewString(part)
	}

	return runtime.NewVector(elements), nil
}

// stringJoin implements the string join core function.
// Usage: (str:join (vector "a" "b" "c") ",") => "a,b,c"
func stringJoin(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:join"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	vec, err := expectVector(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	sep, err := expectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	parts := make([]string, len(vec.Elements))

	for i, elem := range vec.Elements {
		if elem.Type() != runtime.StringType {
			return nil, fmt.Errorf("`%s` expects vector of strings, got %s at index %d", name, elem.Type(), i)
		}
		parts[i] = elem.(runtime.String).Value
	}

	return runtime.NewString(strings.Join(parts, sep.Value)), nil
}

// stringReplace implements the string replacement core function.
// Usage: (str:replace "hello world" "world" "Go") => "hello Go"
func stringReplace(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:replace"

	if err := expectArgs(name, 3, args); err != nil {
		return nil, err
	}

	str, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	oldStr, err := expectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	newStr, err := expectString(name, 2, args[2])
	if err != nil {
		return nil, err
	}

	return runtime.NewString(strings.ReplaceAll(str.Value, oldStr.Value, newStr.Value)), nil
}

// stringStarts implements the string prefix check core function.
// Usage: (str:starts "hello" "he") => true
func stringStarts(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:starts"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	str, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	prefix, err := expectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	return runtime.NewBool(strings.HasPrefix(str.Value, prefix.Value)), nil
}

// stringEnds implements the string suffix check core function.
// Usage: (str:ends "hello" "lo") => true
func stringEnds(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:ends"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	str, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	suffix, err := expectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	return runtime.NewBool(strings.HasSuffix(str.Value, suffix.Value)), nil
}

// stringReverse implements the string reversal core function.
// Usage: (str:reverse "hello") => "olleh"
func stringReverse(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:reverse"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	str, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	runes := []rune(str.Value)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return runtime.NewString(string(runes)), nil
}

// stringRepeat implements the string repetition core function.
// Usage: (str:repeat "ha" 3) => "hahaha"
func stringRepeat(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:repeat"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	str, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	countNum, err := expectIntegerNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	count := int(countNum.Value)

	if count < 0 {
		return nil, fmt.Errorf("`%s` count cannot be negative: %d", name, count)
	}

	return runtime.NewString(strings.Repeat(str.Value, count)), nil
}

// stringConcat implements the string concatenation core function.
// Usage: (str:concat "hello" " " "world") => "hello world"
func stringConcat(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:concat"

	if len(args) == 0 {
		return runtime.NewString(""), nil
	}

	var result strings.Builder

	for i, arg := range args {
		str, err := expectString(name, i, arg)
		if err != nil {
			return nil, err
		}

		result.WriteString(str.Value)
	}

	return runtime.NewString(result.String()), nil
}
