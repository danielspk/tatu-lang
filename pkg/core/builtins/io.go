package builtins

import (
	"fmt"
	"io"

	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterIO registers I/O functions.
func RegisterIO(env *runtime.Environment, out io.Writer) {
	env.DefineNative("print", runtime.NewNativeFunction(printFn(out)))
	env.DefineNative("println", runtime.NewNativeFunction(printlnFn(out)))
}

// printFn implements the print function.
// Usage: (print "hello" " " "world") => nil | prints: hello world
func printFn(out io.Writer) func(args ...runtime.Value) (runtime.Value, error) {
	return func(args ...runtime.Value) (runtime.Value, error) {
		for _, arg := range args {
			_, _ = fmt.Fprint(out, arg)
		}

		return runtime.NewNil(), nil
	}
}

// printlnFn implements the print new line function.
// Usage: (println "hello" " " "world") => nil | prints: hello\nworld\n
func printlnFn(out io.Writer) func(args ...runtime.Value) (runtime.Value, error) {
	return func(args ...runtime.Value) (runtime.Value, error) {
		for _, arg := range args {
			_, _ = fmt.Fprintln(out, arg)
		}

		return runtime.NewNil(), nil
	}
}
