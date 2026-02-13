package stdlib

import (
	"fmt"

	"github.com/danielspk/tatu-lang/pkg/core"
	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterMap registers map functions.
func RegisterMap(natives map[string]runtime.NativeFunction) {
	natives["map:len"] = runtime.NewNativeFunction(mapLen)
	natives["map:get"] = runtime.NewNativeFunction(mapGet)
	natives["map:get-in"] = runtime.NewNativeFunction(mapGetIn)
	natives["map:set"] = runtime.NewNativeFunction(mapSet)
	natives["map:delete"] = runtime.NewNativeFunction(mapDelete)
	natives["map:keys"] = runtime.NewNativeFunction(mapKeys)
	natives["map:values"] = runtime.NewNativeFunction(mapValues)
	natives["map:merge"] = runtime.NewNativeFunction(mapMerge)
	natives["map:has"] = runtime.NewNativeFunction(mapHas)
}

// mapLen implements the map length function.
// Usage: (map:len my-map) => 3
func mapLen(args ...runtime.Value) (runtime.Value, error) {
	const name = "map:len"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	mapValue, err := core.ExpectMap(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(float64(len(mapValue.Elements))), nil
}

// mapGet implements the map value access function.
// Usage: (map:get my-map "key") => value
func mapGet(args ...runtime.Value) (runtime.Value, error) {
	const name = "map:get"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	mapValue, err := core.ExpectMap(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	key, err := core.ExpectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	value, exists := mapValue.Elements[key.Value]
	if !exists {
		return runtime.NewNil(), nil
	}

	return value, nil
}

// mapGetIn implements nested access to maps and vectors.
// Usage: (map:get-in data (vector "user" "addresses" 0 "street")) => value
func mapGetIn(args ...runtime.Value) (runtime.Value, error) {
	const name = "map:get-in"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	path, err := core.ExpectVector(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	current := args[0]

	for i, key := range path.Elements {
		switch key.Type() {
		case runtime.StringType:
			mapValue, ok := current.(runtime.Map)
			if !ok {
				return runtime.NewNil(), nil
			}

			value, exists := mapValue.Elements[key.(runtime.String).Value]
			if !exists {
				return runtime.NewNil(), nil
			}

			current = value

		case runtime.NumberType:
			vecValue, ok := current.(runtime.Vector)
			if !ok {
				return runtime.NewNil(), nil
			}

			value := key.(runtime.Number).Value
			index := int(value)
			if float64(index) != value {
				return nil, fmt.Errorf("`%s` expects integer index at path position %d, got %f", name, i, value)
			}
			if index < 0 || index >= len(vecValue.Elements) {
				return runtime.NewNil(), nil
			}

			current = vecValue.Elements[index]

		default:
			return nil, fmt.Errorf("`%s` expects STRING or NUMBER in path at position %d, got %s", name, i, key.Type())
		}
	}

	return current, nil
}

// mapSet implements the map value assignment function.
// Usage: (map:set my-map "key" value) => modified-map
func mapSet(args ...runtime.Value) (runtime.Value, error) {
	const name = "map:set"

	if err := core.ExpectArgs(name, 3, args); err != nil {
		return nil, err
	}

	mapValue, err := core.ExpectMap(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	key, err := core.ExpectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	mapValue.Elements[key.Value] = args[2]

	return mapValue, nil
}

// mapDelete implements the map key deletion function.
// Usage: (map:delete my-map "key") => modified-map
func mapDelete(args ...runtime.Value) (runtime.Value, error) {
	const name = "map:delete"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	mapValue, err := core.ExpectMap(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	key, err := core.ExpectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	delete(mapValue.Elements, key.Value)

	return mapValue, nil
}

// mapKeys implements the map keys extraction function.
// Usage: (map:keys my-map) => vector-of-keys
func mapKeys(args ...runtime.Value) (runtime.Value, error) {
	const name = "map:keys"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	mapValue, err := core.ExpectMap(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	keys := make([]runtime.Value, 0, len(mapValue.Elements))

	for key := range mapValue.Elements {
		keys = append(keys, runtime.NewString(key))
	}

	return runtime.NewVector(keys), nil
}

// mapValues implements the map values extraction function.
// Usage: (map:values my-map) => vector-of-values
func mapValues(args ...runtime.Value) (runtime.Value, error) {
	const name = "map:values"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	mapValue, err := core.ExpectMap(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	values := make([]runtime.Value, 0, len(mapValue.Elements))

	for _, value := range mapValue.Elements {
		values = append(values, value)
	}

	return runtime.NewVector(values), nil
}

// mapMerge implements the map merging function.
// Usage: (map:merge my-map other-map) => modified-map
func mapMerge(args ...runtime.Value) (runtime.Value, error) {
	const name = "map:merge"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	mapValue, err := core.ExpectMap(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	otherMap, err := core.ExpectMap(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	for key, value := range otherMap.Elements {
		mapValue.Elements[key] = value
	}

	return mapValue, nil
}

// mapHas implements the map key existence check function.
// Usage: (map:has my-map "key") => boolean
func mapHas(args ...runtime.Value) (runtime.Value, error) {
	const name = "map:has"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	mapValue, err := core.ExpectMap(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	key, err := core.ExpectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	_, exists := mapValue.Elements[key.Value]

	return runtime.NewBool(exists), nil
}
