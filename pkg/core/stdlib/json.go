package stdlib

import (
	"encoding/json"
	"fmt"

	"github.com/danielspk/tatu-lang/pkg/core"
	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterJSON registers JSON functions.
func RegisterJSON(natives map[string]runtime.NativeFunction) {
	natives["json:encode"] = runtime.NewNativeFunction(jsonEncode)
	natives["json:decode"] = runtime.NewNativeFunction(jsonDecode)
}

// jsonEncode implements the JSON encoding function.
// Usage: (json:encode (map "name" "John" "age" 30)) => "{\"age\":30,\"name\":\"John\"}"
func jsonEncode(args ...runtime.Value) (runtime.Value, error) {
	const name = "json:encode"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	data, err := tatuToJSON(args[0])
	if err != nil {
		return nil, fmt.Errorf("`%s` unsupported type: %w", name, err)
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("`%s` failed to encode: %w", name, err)
	}

	return runtime.NewString(string(jsonBytes)), nil
}

// jsonDecode implements the JSON decoding function.
// Usage: (json:decode "{\"name\":\"John\",\"age\":30}") => (map "name" "John" "age" 30)
func jsonDecode(args ...runtime.Value) (runtime.Value, error) {
	const name = "json:decode"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	str, err := core.ExpectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	var data any
	if err = json.Unmarshal([]byte(str.Value), &data); err != nil {
		return nil, fmt.Errorf("`%s` failed to decode: %w", name, err)
	}

	result, err := jsonToTatu(data)
	if err != nil {
		return nil, fmt.Errorf("`%s` failed to convert: %w", name, err)
	}

	return result, nil
}

// tatuToJSON converts a Tatu runtime.Value to a Go any for json.Marshal.
func tatuToJSON(value runtime.Value) (any, error) {
	switch value.Type() {
	case runtime.NilType:
		return nil, nil
	case runtime.BoolType:
		return value.(runtime.Bool).Value, nil
	case runtime.NumberType:
		return value.(runtime.Number).Value, nil
	case runtime.StringType:
		return value.(runtime.String).Value, nil
	case runtime.VectorType:
		vec := value.(runtime.Vector)
		result := make([]any, len(vec.Elements))

		for i, elem := range vec.Elements {
			val, err := tatuToJSON(elem)
			if err != nil {
				return nil, err
			}
			result[i] = val
		}
		return result, nil
	case runtime.MapType:
		m := value.(runtime.Map)
		result := make(map[string]any)

		for key, val := range m.Elements {
			jsonVal, err := tatuToJSON(val)
			if err != nil {
				return nil, err
			}
			result[key] = jsonVal
		}
		return result, nil
	default:
		return nil, fmt.Errorf("cannot convert %s to JSON", value.Type())
	}
}

// jsonToTatu converts a Go any from json.Unmarshal to a Tatu runtime.Value.
func jsonToTatu(data any) (runtime.Value, error) {
	if data == nil {
		return runtime.NewNil(), nil
	}

	switch v := data.(type) {
	case bool:
		return runtime.NewBool(v), nil
	case float64:
		return runtime.NewNumber(v), nil
	case string:
		return runtime.NewString(v), nil
	case []any:
		elements := make([]runtime.Value, len(v))

		for i, item := range v {
			val, err := jsonToTatu(item)
			if err != nil {
				return nil, err
			}
			elements[i] = val
		}
		return runtime.NewVector(elements), nil
	case map[string]any:
		elements := make(map[string]runtime.Value)

		for key, value := range v {
			val, err := jsonToTatu(value)
			if err != nil {
				return nil, err
			}
			elements[key] = val
		}
		return runtime.NewMap(elements), nil
	default:
		return nil, fmt.Errorf("unsupported JSON type: %T", v)
	}
}
