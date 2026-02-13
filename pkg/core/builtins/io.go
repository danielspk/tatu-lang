package builtins

import (
	"fmt"

	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterIO registers I/O functions.
func RegisterIO(natives map[string]runtime.NativeFunction) {
	natives["print"] = runtime.NewNativeFunction(printFn)
}

// printFn implements the print function.
// Usage: (print "hello" " " "world") => nil (prints: hello world)
func printFn(args ...runtime.Value) (runtime.Value, error) {
	for _, arg := range args {
		fmt.Print(arg)
	}

	fmt.Println()

	return runtime.NewNil(), nil
}
