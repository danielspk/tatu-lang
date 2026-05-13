package stdlib

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/danielspk/tatu-lang/pkg/core"
	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// regexCache memorizes compiled patterns.
var regexCache sync.Map

// RegisterRegex registers regular expression functions.
func RegisterRegex(env *runtime.Environment) {
	env.DefineNative("regex:matches", runtime.NewNativeFunction(regexMatches))
	env.DefineNative("regex:find", runtime.NewNativeFunction(regexFind))
	env.DefineNative("regex:replace", runtime.NewNativeFunction(regexReplace))
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

	re, err := compileCached(name, pattern.Value)
	if err != nil {
		return nil, err
	}

	return runtime.NewBool(re.MatchString(str.Value)), nil
}

// regexFind finds the first substring that matches a regular expression pattern.
// Usage: (regex:find "hello 123 world" "[0-9]+") => "123"
// Usage: (regex:find "hello" "[0-9]+") => nil
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

	re, err := compileCached(name, pattern.Value)
	if err != nil {
		return nil, err
	}

	loc := re.FindStringIndex(str.Value)
	if loc == nil {
		return runtime.NewNil(), nil
	}

	return runtime.NewString(str.Value[loc[0]:loc[1]]), nil
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

	re, err := compileCached(name, pattern.Value)
	if err != nil {
		return nil, err
	}

	result := re.ReplaceAllString(str.Value, replacement.Value)

	return runtime.NewString(result), nil
}

func compileCached(name, pattern string) (*regexp.Regexp, error) {
	if v, ok := regexCache.Load(pattern); ok {
		return v.(*regexp.Regexp), nil
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("`%s` invalid regex pattern: %w", name, err)
	}

	actual, _ := regexCache.LoadOrStore(pattern, re)

	return actual.(*regexp.Regexp), nil
}
