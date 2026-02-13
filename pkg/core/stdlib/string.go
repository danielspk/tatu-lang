package stdlib

import (
	"fmt"
	"strings"

	"github.com/danielspk/tatu-lang/pkg/core"
	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterString registers string functions.
func RegisterString(natives map[string]runtime.NativeFunction) {
	natives["str:len"] = runtime.NewNativeFunction(stringLen)
	natives["str:contains"] = runtime.NewNativeFunction(stringContains)
	natives["str:index"] = runtime.NewNativeFunction(stringIndex)
	natives["str:upper"] = runtime.NewNativeFunction(stringUpper)
	natives["str:lower"] = runtime.NewNativeFunction(stringLower)
	natives["str:trim"] = runtime.NewNativeFunction(stringTrim)
	natives["str:slice"] = runtime.NewNativeFunction(stringSlice)
	natives["str:split"] = runtime.NewNativeFunction(stringSplit)
	natives["str:join"] = runtime.NewNativeFunction(stringJoin)
	natives["str:replace"] = runtime.NewNativeFunction(stringReplace)
	natives["str:starts"] = runtime.NewNativeFunction(stringStarts)
	natives["str:ends"] = runtime.NewNativeFunction(stringEnds)
	natives["str:reverse"] = runtime.NewNativeFunction(stringReverse)
	natives["str:repeat"] = runtime.NewNativeFunction(stringRepeat)
	natives["str:concat"] = runtime.NewNativeFunction(stringConcat)
}

// stringLen implements the string length function.
// Usage: (str:len "hello") => 5
func stringLen(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:len"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	str, err := core.ExpectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(float64(len([]rune(str.Value)))), nil
}

// stringContains implements the string contains check function.
// Usage: (str:contains "hello world" "world") => true
func stringContains(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:contains"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	str, err := core.ExpectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	substr, err := core.ExpectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	return runtime.NewBool(strings.Contains(str.Value, substr.Value)), nil
}

// stringIndex implements the string index search function.
// Usage: (str:index "hello" "ll") => 2
func stringIndex(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:index"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	str, err := core.ExpectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	substr, err := core.ExpectString(name, 1, args[1])
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

// stringUpper implements the string uppercase conversion function.
// Usage: (str:upper "hello") => "HELLO"
func stringUpper(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:upper"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	str, err := core.ExpectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewString(strings.ToUpper(str.Value)), nil
}

// stringLower implements the string lowercase conversion function.
// Usage: (str:lower "HELLO") => "hello"
func stringLower(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:lower"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	str, err := core.ExpectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewString(strings.ToLower(str.Value)), nil
}

// stringTrim implements the string whitespace trimming function.
// Usage: (str:trim "  hello  ") => "hello"
func stringTrim(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:trim"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	str, err := core.ExpectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewString(strings.TrimSpace(str.Value)), nil
}

// stringSlice implements the string slice extraction function.
// Usage: (str:slice "hello" 1 4) => "ell"
func stringSlice(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:slice"

	if err := core.ExpectArgs(name, 3, args); err != nil {
		return nil, err
	}

	str, err := core.ExpectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	startNum, err := core.ExpectIntegerNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	endNum, err := core.ExpectIntegerNumber(name, 2, args[2])
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

// stringSplit implements the string split function.
// Usage: (str:split "a,b,c" ",") => ("a" "b" "c")
func stringSplit(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:split"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	str, err := core.ExpectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	sep, err := core.ExpectString(name, 1, args[1])
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

// stringJoin implements the string join function.
// Usage: (str:join (vector "a" "b" "c") ",") => "a,b,c"
func stringJoin(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:join"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	vec, err := core.ExpectVector(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	sep, err := core.ExpectString(name, 1, args[1])
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

// stringReplace implements the string replacement function.
// Usage: (str:replace "hello world" "world" "Go") => "hello Go"
func stringReplace(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:replace"

	if err := core.ExpectArgs(name, 3, args); err != nil {
		return nil, err
	}

	str, err := core.ExpectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	oldStr, err := core.ExpectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	newStr, err := core.ExpectString(name, 2, args[2])
	if err != nil {
		return nil, err
	}

	return runtime.NewString(strings.ReplaceAll(str.Value, oldStr.Value, newStr.Value)), nil
}

// stringStarts implements the string prefix check function.
// Usage: (str:starts "hello" "he") => true
func stringStarts(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:starts"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	str, err := core.ExpectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	prefix, err := core.ExpectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	return runtime.NewBool(strings.HasPrefix(str.Value, prefix.Value)), nil
}

// stringEnds implements the string suffix check function.
// Usage: (str:ends "hello" "lo") => true
func stringEnds(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:ends"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	str, err := core.ExpectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	suffix, err := core.ExpectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	return runtime.NewBool(strings.HasSuffix(str.Value, suffix.Value)), nil
}

// stringReverse implements the string reversal function.
// Usage: (str:reverse "hello") => "olleh"
func stringReverse(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:reverse"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	str, err := core.ExpectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	runes := []rune(str.Value)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return runtime.NewString(string(runes)), nil
}

// stringRepeat implements the string repetition function.
// Usage: (str:repeat "ha" 3) => "hahaha"
func stringRepeat(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:repeat"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	str, err := core.ExpectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	countNum, err := core.ExpectIntegerNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	count := int(countNum.Value)

	if count < 0 {
		return nil, fmt.Errorf("`%s` count cannot be negative: %d", name, count)
	}

	return runtime.NewString(strings.Repeat(str.Value, count)), nil
}

// stringConcat implements the string concatenation function.
// Usage: (str:concat "hello" " " "world") => "hello world"
func stringConcat(args ...runtime.Value) (runtime.Value, error) {
	const name = "str:concat"

	if len(args) == 0 {
		return runtime.NewString(""), nil
	}

	var result strings.Builder

	for i, arg := range args {
		str, err := core.ExpectString(name, i, arg)
		if err != nil {
			return nil, err
		}

		result.WriteString(str.Value)
	}

	return runtime.NewString(result.String()), nil
}
