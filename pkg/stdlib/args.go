package stdlib

import (
	"fmt"

	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// expectArgs validates the number of arguments.
func expectArgs(name string, expected int, args []runtime.Value) error {
	if len(args) != expected {
		return fmt.Errorf("`%s` expects %d argument(s), got %d", name, expected, len(args))
	}

	return nil
}

// expectNumber validates that an argument is NUMBER and returns it.
func expectNumber(name string, argIndex int, arg runtime.Value) (runtime.Number, error) {
	if arg.Type() != runtime.NumberType {
		return runtime.Number{}, fmt.Errorf("`%s` expects NUMBER at argument %d, got %s", name, argIndex+1, arg.Type())
	}

	return arg.(runtime.Number), nil
}

// expectString validates that an argument is STRING and returns it.
func expectString(name string, argIndex int, arg runtime.Value) (runtime.String, error) {
	if arg.Type() != runtime.StringType {
		return runtime.String{}, fmt.Errorf("`%s` expects STRING at argument %d, got %s", name, argIndex+1, arg.Type())
	}

	return arg.(runtime.String), nil
}

// expectBool validates that an argument is BOOL and returns it.
func expectBool(name string, argIndex int, arg runtime.Value) (runtime.Bool, error) {
	if arg.Type() != runtime.BoolType {
		return runtime.Bool{}, fmt.Errorf("`%s` expects BOOL at argument %d, got %s", name, argIndex+1, arg.Type())
	}

	return arg.(runtime.Bool), nil
}

// expectVector validates that an argument is VECTOR and returns it.
func expectVector(name string, argIndex int, arg runtime.Value) (runtime.Vector, error) {
	if arg.Type() != runtime.VectorType {
		return runtime.Vector{}, fmt.Errorf("`%s` expects VECTOR at argument %d, got %s", name, argIndex+1, arg.Type())
	}

	return arg.(runtime.Vector), nil
}

// expectMap validates that an argument is MAP and returns it.
func expectMap(name string, argIndex int, arg runtime.Value) (runtime.Map, error) {
	if arg.Type() != runtime.MapType {
		return runtime.Map{}, fmt.Errorf("`%s` expects MAP at argument %d, got %s", name, argIndex+1, arg.Type())
	}

	return arg.(runtime.Map), nil
}
