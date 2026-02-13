package stdlib

import (
	"fmt"
	"sort"

	"github.com/danielspk/tatu-lang/pkg/core"
	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterVector registers vector functions.
func RegisterVector(natives map[string]runtime.NativeFunction) {
	natives["vec:len"] = runtime.NewNativeFunction(vectorLen)
	natives["vec:get"] = runtime.NewNativeFunction(vectorGet)
	natives["vec:set"] = runtime.NewNativeFunction(vectorSet)
	natives["vec:delete"] = runtime.NewNativeFunction(vectorDelete)
	natives["vec:push"] = runtime.NewNativeFunction(vectorPush)
	natives["vec:pop"] = runtime.NewNativeFunction(vectorPop)
	natives["vec:slice"] = runtime.NewNativeFunction(vectorSlice)
	natives["vec:concat"] = runtime.NewNativeFunction(vectorConcat)
	natives["vec:contains"] = runtime.NewNativeFunction(vectorContains)
	natives["vec:find"] = runtime.NewNativeFunction(vectorFind)
	natives["vec:reverse"] = runtime.NewNativeFunction(vectorReverse)
	natives["vec:sort"] = runtime.NewNativeFunction(vectorSort)
}

// vectorLen implements the vector length function.
// Usage: (vec:len my-vector) => 3
func vectorLen(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:len"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	vector, err := core.ExpectVector(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(float64(len(vector.Elements))), nil
}

// vectorGet implements the vector element access function.
// Usage: (vec:get my-vector index) => element
func vectorGet(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:get"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	vector, index, err := validateVectorIndex(name, args)
	if err != nil {
		return nil, err
	}

	return vector.Elements[index], nil
}

// vectorSet implements the vector element update function.
// Usage: (vec:set my-vector index value) => modified-vector
func vectorSet(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:set"

	if err := core.ExpectArgs(name, 3, args); err != nil {
		return nil, err
	}

	vector, index, err := validateVectorIndex(name, args)
	if err != nil {
		return nil, err
	}

	vector.Elements[index] = args[2]

	return vector, nil
}

// vectorDelete implements the vector element deletion function.
// Usage: (vec:delete my-vector index) => modified-vector
func vectorDelete(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:delete"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	vector, index, err := validateVectorIndex(name, args)
	if err != nil {
		return nil, err
	}

	vector.Elements = append(vector.Elements[:index], vector.Elements[index+1:]...)

	return vector, nil
}

// vectorPush implements the vector element append function.
// Usage: (vec:push my-vector value) => modified-vector
func vectorPush(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:push"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	vector, err := core.ExpectVector(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	vector.Elements = append(vector.Elements, args[1])

	return vector, nil
}

// vectorPop implements the vector element removal function.
// Usage: (vec:pop my-vector) => modified-vector
func vectorPop(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:pop"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	vector, err := core.ExpectVector(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	if len(vector.Elements) == 0 {
		return nil, fmt.Errorf("`%s` cannot pop from empty vector", name)
	}

	vector.Elements = vector.Elements[:len(vector.Elements)-1]

	return vector, nil
}

// vectorSlice implements the vector slice extraction function.
// Usage: (vec:slice my-vector start end) => new-vector
func vectorSlice(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:slice"

	if err := core.ExpectArgs(name, 3, args); err != nil {
		return nil, err
	}

	vector, err := core.ExpectVector(name, 0, args[0])
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

	start := int(startNum.Value)
	end := int(endNum.Value)

	if start < 0 || start > len(vector.Elements) {
		return nil, fmt.Errorf("`%s` start index out of bounds: %d (vector length: %d)", name, start, len(vector.Elements))
	}

	if end < 0 || end > len(vector.Elements) {
		return nil, fmt.Errorf("`%s` end index out of bounds: %d (vector length: %d)", name, end, len(vector.Elements))
	}

	if start > end {
		return nil, fmt.Errorf("`%s` start index (%d) cannot be greater than end index (%d)", name, start, end)
	}

	newElements := make([]runtime.Value, end-start)
	copy(newElements, vector.Elements[start:end])

	return runtime.NewVector(newElements), nil
}

// vectorConcat implements the vector concatenation function.
// Usage: (vec:concat my-vector other-vector) => modified-vector
func vectorConcat(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:concat"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	vector, err := core.ExpectVector(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	otherVector, err := core.ExpectVector(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	vector.Elements = append(vector.Elements, otherVector.Elements...)

	return vector, nil
}

// vectorContains implements the vector element search function.
// Usage: (vec:contains my-vector value) => boolean
func vectorContains(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:contains"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	vector, err := core.ExpectVector(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	searchValue := args[1]

	for _, elem := range vector.Elements {
		if valuesEqual(elem, searchValue) {
			return runtime.NewBool(true), nil
		}
	}

	return runtime.NewBool(false), nil
}

// vectorFind implements the vector element index search function.
// Usage: (vec:find my-vector value) => index or nil
func vectorFind(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:find"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	vector, err := core.ExpectVector(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	searchValue := args[1]

	for i, elem := range vector.Elements {
		if valuesEqual(elem, searchValue) {
			return runtime.NewNumber(float64(i)), nil
		}
	}

	return runtime.NewNil(), nil
}

// vectorReverse implements the vector reversal function.
// Usage: (vec:reverse my-vector) => modified-vector
func vectorReverse(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:reverse"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	vector, err := core.ExpectVector(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	for i, j := 0, len(vector.Elements)-1; i < j; i, j = i+1, j-1 {
		vector.Elements[i], vector.Elements[j] = vector.Elements[j], vector.Elements[i]
	}

	return vector, nil
}

// vectorSort sorts a vector in ascending order.
// Usage: (vec:sort (vector 3 1 2)) => (vector 1 2 3)
func vectorSort(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:sort"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	vector, err := core.ExpectVector(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	if len(vector.Elements) == 0 {
		return vector, nil
	}

	firstType := vector.Elements[0].Type()

	sort.Slice(vector.Elements, func(i, j int) bool {
		a := vector.Elements[i]
		b := vector.Elements[j]

		if a.Type() != firstType || b.Type() != firstType {
			return false
		}

		switch firstType {
		case runtime.NumberType:
			return a.(runtime.Number).Value < b.(runtime.Number).Value
		case runtime.StringType:
			return a.(runtime.String).Value < b.(runtime.String).Value
		case runtime.BoolType:
			return !a.(runtime.Bool).Value && b.(runtime.Bool).Value
		default:
			return false
		}
	})

	return vector, nil
}

func validateVectorIndex(name string, args []runtime.Value) (runtime.Vector, int, error) {
	vector, err := core.ExpectVector(name, 0, args[0])
	if err != nil {
		return runtime.Vector{}, 0, err
	}

	number, err := core.ExpectIntegerNumber(name, 1, args[1])
	if err != nil {
		return runtime.Vector{}, 0, err
	}

	index := int(number.Value)

	if index < 0 || index >= len(vector.Elements) {
		return runtime.Vector{}, 0, fmt.Errorf("`%s` index out of bounds: %d (vector length: %d)", name, index, len(vector.Elements))
	}

	return vector, index, nil
}

func valuesEqual(a, b runtime.Value) bool {
	if a.Type() != b.Type() {
		return false
	}

	switch a.Type() {
	case runtime.NumberType:
		return a.(runtime.Number).Value == b.(runtime.Number).Value
	case runtime.StringType:
		return a.(runtime.String).Value == b.(runtime.String).Value
	case runtime.BoolType:
		return a.(runtime.Bool).Value == b.(runtime.Bool).Value
	case runtime.NilType:
		return true
	default:
		return false
	}
}
