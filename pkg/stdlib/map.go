package stdlib

import (
	"fmt"

	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterMap registers map core functions in the environment.
func RegisterMap(env *runtime.Environment) error {
	functions := map[string]runtime.CoreFunction{
		"map:len": runtime.NewCoreFunction(mapLen),
	}

	for name, fn := range functions {
		if _, err := env.Define(name, fn); err != nil {
			return fmt.Errorf("failed to register map function `%s`: %v", name, err)
		}
	}

	return nil
}

// mapLen implements the map length core function.
// Usage: (map:len my-map) => 3
func mapLen(args ...runtime.Value) (runtime.Value, error) {
	const name = "map:len"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	mapValue, err := expectMap(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(float64(len(mapValue.Elements))), nil
}
