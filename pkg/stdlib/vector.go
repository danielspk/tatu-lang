package stdlib

import (
	"fmt"

	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterVector registers vector core functions in the environment.
func RegisterVector(env *runtime.Environment) error {
	functions := map[string]runtime.CoreFunction{
		"vec:len":      runtime.NewCoreFunction(vectorLen),
		"vec:get":      runtime.NewCoreFunction(vectorGet),
		"vec:set":      runtime.NewCoreFunction(vectorSet),
		"vec:delete":   runtime.NewCoreFunction(vectorDelete),
		"vec:push":     runtime.NewCoreFunction(vectorPush),
		"vec:pop":      runtime.NewCoreFunction(vectorPop),
		"vec:slice":    runtime.NewCoreFunction(vectorSlice),
		"vec:concat":   runtime.NewCoreFunction(vectorConcat),
		"vec:contains": runtime.NewCoreFunction(vectorContains),
		"vec:find":     runtime.NewCoreFunction(vectorFind),
		"vec:reverse":  runtime.NewCoreFunction(vectorReverse),
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

// vectorGet implements the vector element access core function.
// Usage: (vec:get my-vector index) => element
func vectorGet(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:get"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	vector, index, err := validateVectorIndex(name, args)
	if err != nil {
		return nil, err
	}

	return vector.Elements[index], nil
}

// vectorSet implements the vector element update core function.
// Usage: (vec:set my-vector index value) => modified-vector
func vectorSet(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:set"

	if err := expectArgs(name, 3, args); err != nil {
		return nil, err
	}

	vector, index, err := validateVectorIndex(name, args)
	if err != nil {
		return nil, err
	}

	vector.Elements[index] = args[2]

	return vector, nil
}

// vectorDelete implements the vector element deletion core function.
// Usage: (vec:delete my-vector index) => modified-vector
func vectorDelete(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:delete"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	vector, index, err := validateVectorIndex(name, args)
	if err != nil {
		return nil, err
	}

	vector.Elements = append(vector.Elements[:index], vector.Elements[index+1:]...)

	return vector, nil
}

// vectorPush implements the vector element append core function.
// Usage: (vec:push my-vector value) => modified-vector
func vectorPush(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:push"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	vector, err := expectVector(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	vector.Elements = append(vector.Elements, args[1])

	return vector, nil
}

// vectorPop implements the vector element removal core function.
// Usage: (vec:pop my-vector) => modified-vector
func vectorPop(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:pop"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	vector, err := expectVector(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	if len(vector.Elements) == 0 {
		return nil, fmt.Errorf("`%s` cannot pop from empty vector", name)
	}

	vector.Elements = vector.Elements[:len(vector.Elements)-1]

	return vector, nil
}

// vectorSlice implements the vector slice extraction core function.
// Usage: (vec:slice my-vector start end) => new-vector
func vectorSlice(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:slice"

	if err := expectArgs(name, 3, args); err != nil {
		return nil, err
	}

	vector, err := expectVector(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	startNum, err := expectIntegerNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	endNum, err := expectIntegerNumber(name, 2, args[2])
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

// vectorConcat implements the vector concatenation core function.
// Usage: (vec:concat my-vector other-vector) => modified-vector
func vectorConcat(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:concat"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	vector, err := expectVector(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	otherVector, err := expectVector(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	vector.Elements = append(vector.Elements, otherVector.Elements...)

	return vector, nil
}

// vectorContains implements the vector element search core function.
// Usage: (vec:contains my-vector value) => boolean
func vectorContains(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:contains"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	vector, err := expectVector(name, 0, args[0])
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

// vectorFind implements the vector element index search core function.
// Usage: (vec:find my-vector value) => index or nil
func vectorFind(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:find"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	vector, err := expectVector(name, 0, args[0])
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

// vectorReverse implements the vector reversal core function.
// Usage: (vec:reverse my-vector) => modified-vector
func vectorReverse(args ...runtime.Value) (runtime.Value, error) {
	const name = "vec:reverse"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	vector, err := expectVector(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	for i, j := 0, len(vector.Elements)-1; i < j; i, j = i+1, j-1 {
		vector.Elements[i], vector.Elements[j] = vector.Elements[j], vector.Elements[i]
	}

	return vector, nil
}

func validateVectorIndex(name string, args []runtime.Value) (runtime.Vector, int, error) {
	vector, err := expectVector(name, 0, args[0])
	if err != nil {
		return runtime.Vector{}, 0, err
	}

	number, err := expectIntegerNumber(name, 1, args[1])
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
