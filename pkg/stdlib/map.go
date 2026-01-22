package stdlib

import (
	"fmt"

	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterMap registers map core functions in the environment.
func RegisterMap(env *runtime.Environment) error {
	functions := map[string]runtime.CoreFunction{
		"map:len":    runtime.NewCoreFunction(mapLen),
		"map:get":    runtime.NewCoreFunction(mapGet),
		"map:set":    runtime.NewCoreFunction(mapSet),
		"map:delete": runtime.NewCoreFunction(mapDelete),
		"map:keys":   runtime.NewCoreFunction(mapKeys),
		"map:values": runtime.NewCoreFunction(mapValues),
		"map:merge":  runtime.NewCoreFunction(mapMerge),
		"map:has":    runtime.NewCoreFunction(mapHas),
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

// mapGet implements the map value access core function.
// Usage: (map:get my-map "key") => value
func mapGet(args ...runtime.Value) (runtime.Value, error) {
	const name = "map:get"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	mapValue, err := expectMap(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	key, err := expectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	value, exists := mapValue.Elements[key.Value]
	if !exists {
		return runtime.NewNil(), nil
	}

	return value, nil
}

// mapSet implements the map value assignment core function.
// Usage: (map:set my-map "key" value) => modified-map
func mapSet(args ...runtime.Value) (runtime.Value, error) {
	const name = "map:set"

	if err := expectArgs(name, 3, args); err != nil {
		return nil, err
	}

	mapValue, err := expectMap(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	key, err := expectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	mapValue.Elements[key.Value] = args[2]

	return mapValue, nil
}

// mapDelete implements the map key deletion core function.
// Usage: (map:delete my-map "key") => modified-map
func mapDelete(args ...runtime.Value) (runtime.Value, error) {
	const name = "map:delete"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	mapValue, err := expectMap(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	key, err := expectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	delete(mapValue.Elements, key.Value)

	return mapValue, nil
}

// mapKeys implements the map keys extraction core function.
// Usage: (map:keys my-map) => vector-of-keys
func mapKeys(args ...runtime.Value) (runtime.Value, error) {
	const name = "map:keys"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	mapValue, err := expectMap(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	keys := make([]runtime.Value, 0, len(mapValue.Elements))
	for key := range mapValue.Elements {
		keys = append(keys, runtime.NewString(key))
	}

	return runtime.NewVector(keys), nil
}

// mapValues implements the map values extraction core function.
// Usage: (map:values my-map) => vector-of-values
func mapValues(args ...runtime.Value) (runtime.Value, error) {
	const name = "map:values"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	mapValue, err := expectMap(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	values := make([]runtime.Value, 0, len(mapValue.Elements))
	for _, value := range mapValue.Elements {
		values = append(values, value)
	}

	return runtime.NewVector(values), nil
}

// mapMerge implements the map merging core function.
// Usage: (map:merge my-map other-map) => modified-map
func mapMerge(args ...runtime.Value) (runtime.Value, error) {
	const name = "map:merge"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	mapValue, err := expectMap(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	otherMap, err := expectMap(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	for key, value := range otherMap.Elements {
		mapValue.Elements[key] = value
	}

	return mapValue, nil
}

// mapHas implements the map key existence check core function.
// Usage: (map:has my-map "key") => boolean
func mapHas(args ...runtime.Value) (runtime.Value, error) {
	const name = "map:has"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	mapValue, err := expectMap(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	key, err := expectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	_, exists := mapValue.Elements[key.Value]

	return runtime.NewBool(exists), nil
}
