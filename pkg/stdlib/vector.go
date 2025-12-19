package stdlib

import (
	"fmt"

	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterVector registers vector core functions in the environment.
func RegisterVector(env *runtime.Environment) error {
	functions := map[string]runtime.CoreFunction{
		"vec:len": runtime.NewCoreFunction(vectorLen),
	}

	for name, fn := range functions {
		if _, err := env.Define(name, fn); err != nil {
			return fmt.Errorf("failed to register vector function `%s`: %v", name, err)
		}
	}

	return nil
}

// vectorLen implements the vector length core function.
// Usage: (vec:len my-vector) => 3
func vectorLen(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:len"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	vector, err := expectVector(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(float64(len(vector.Elements))), nil
}
