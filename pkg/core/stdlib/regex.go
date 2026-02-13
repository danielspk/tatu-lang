package stdlib

import (
	"fmt"
	"regexp"

	"github.com/danielspk/tatu-lang/pkg/core"
	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterRegex registers regular expression functions.
func RegisterRegex(natives map[string]runtime.NativeFunction) {
	natives["regex:matches"] = runtime.NewNativeFunction(regexMatches)
	natives["regex:find"] = runtime.NewNativeFunction(regexFind)
	natives["regex:replace"] = runtime.NewNativeFunction(regexReplace)
}

// regexMatches checks if a string matches a regular expression pattern.
// Usage: (regex:matches "hello123" "^[a-z]+[0-9]+$") => true
func regexMatches(args ...runtime.Value) (runtime.Value, error) {
	const name = "regex:matches"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	str, err := core.ExpectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	pattern, err := core.ExpectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	matched, err := regexp.MatchString(pattern.Value, str.Value)
	if err != nil {
		return nil, fmt.Errorf("`%s` invalid regex pattern: %w", name, err)
	}

	return runtime.NewBool(matched), nil
}

// regexFind finds the first substring that matches a regular expression pattern.
// Usage: (regex:find "hello 123 world" "[0-9]+") => "123"
func regexFind(args ...runtime.Value) (runtime.Value, error) {
	const name = "regex:find"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	str, err := core.ExpectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	pattern, err := core.ExpectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(pattern.Value)
	if err != nil {
		return nil, fmt.Errorf("`%s` invalid regex pattern: %w", name, err)
	}

	match := re.FindString(str.Value)

	return runtime.NewString(match), nil
}

// regexReplace replaces all substrings that match a regular expression pattern with a replacement string.
// Usage: (regex:replace "hello 123 world 456" "[0-9]+" "NUM") => "hello NUM world NUM"
func regexReplace(args ...runtime.Value) (runtime.Value, error) {
	const name = "regex:replace"

	if err := core.ExpectArgs(name, 3, args); err != nil {
		return nil, err
	}

	str, err := core.ExpectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	pattern, err := core.ExpectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	replacement, err := core.ExpectString(name, 2, args[2])
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(pattern.Value)
	if err != nil {
		return nil, fmt.Errorf("`%s` invalid regex pattern: %w", name, err)
	}

	result := re.ReplaceAllString(str.Value, replacement.Value)

	return runtime.NewString(result), nil
}
