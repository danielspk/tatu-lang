// Package runtime implements the runtime values and scopes system.
package runtime

import (
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/danielspk/tatu-lang/pkg/ast"
)

// ValueType represents the type of value.
type ValueType string

// Value types.
const (
	NumberType     ValueType = "NUMBER"
	StringType     ValueType = "STRING"
	BoolType       ValueType = "BOOL"
	NilType        ValueType = "NIL"
	VectorType     ValueType = "VECTOR"
	MapType        ValueType = "MAP"
	FuncType       ValueType = "FUNC"
	NativeFuncType ValueType = "NATIVE_FUNC"
	RecurType      ValueType = "RECUR"
)

// Value represents a value interface.
type Value interface {
	Type() ValueType
	String() string
	Equal(other Value) bool
}

// Number represents a value of a number type.
type Number struct {
	Value float64
}

// NewNumber builds a new Number.
func NewNumber(value float64) Number {
	return Number{value}
}

// Type returns the type of the number value.
func (n Number) Type() ValueType {
	return NumberType
}

// String returns the string representation of the number value.
func (n Number) String() string {
	// normalize -0 to 0 because is an IEEE 754 valid number
	if n.Value == 0 {
		return "0"
	}

	if n.Value == math.Trunc(n.Value) {
		return fmt.Sprintf("%.0f", n.Value)
	}

	return strconv.FormatFloat(n.Value, 'g', -1, 64)
}

// Equal compares the number value to another.
func (n Number) Equal(other Value) bool {
	if other.Type() != NumberType {
		return false
	}

	return n.Value == other.(Number).Value
}

// String represents a value of a string type.
type String struct {
	Value string
}

// NewString builds a new String.
func NewString(value string) String {
	return String{value}
}

// Type returns the type of the string value.
func (s String) Type() ValueType {
	return StringType
}

// String returns the string representation of the string value.
func (s String) String() string {
	return s.Value
}

// Equal compares the string value to another.
func (s String) Equal(other Value) bool {
	if other.Type() != StringType {
		return false
	}

	return s.Value == other.(String).Value
}

// Bool represents a value of a boolean type.
type Bool struct {
	Value bool
}

// NewBool builds a new Bool.
func NewBool(value bool) Bool {
	return Bool{value}
}

// Type returns the type of the boolean value.
func (b Bool) Type() ValueType {
	return BoolType
}

// String returns the string representation of the boolean value.
func (b Bool) String() string {
	return fmt.Sprintf("%v", b.Value)
}

// Equal compares the boolean value to another.
func (b Bool) Equal(other Value) bool {
	if other.Type() != BoolType {
		return false
	}

	return b.Value == other.(Bool).Value
}

// Nil represents a value of nil type.
type Nil struct{}

// NewNil builds a new Nil.
func NewNil() Nil {
	return Nil{}
}

// Type returns the type of the nil value.
func (n Nil) Type() ValueType {
	return NilType
}

// String returns the string representation of the nil value.
func (n Nil) String() string {
	return "<nil>"
}

// Equal compares the nil value to another.
func (n Nil) Equal(other Value) bool {
	return other.Type() == NilType
}

// Vector represents a value of a vector type.
type Vector struct {
	Elements []Value
}

// NewVector builds a new Vector.
func NewVector(elements []Value) *Vector {
	return &Vector{elements}
}

// Type returns the type of the vector value.
func (v *Vector) Type() ValueType {
	return VectorType
}

// String returns the string representation of the vector value.
func (v *Vector) String() string {
	out := "("

	for i, e := range v.Elements {
		if i > 0 {
			out += " "
		}
		out += e.String()
	}

	out += ")"

	return out
}

// Equal compares the vector value to another.
func (v *Vector) Equal(other Value) bool {
	if other.Type() != VectorType {
		return false
	}

	o := other.(*Vector)

	if len(v.Elements) != len(o.Elements) {
		return false
	}

	for i, e := range v.Elements {
		if !e.Equal(o.Elements[i]) {
			return false
		}
	}

	return true
}

// Map represents a value of map type.
type Map struct {
	Elements map[string]Value
}

// NewMap builds a new Map.
func NewMap(elements map[string]Value) Map {
	return Map{elements}
}

// Type returns the type of the map value.
func (m Map) Type() ValueType {
	return MapType
}

// String returns the string representation of the map value with sorted keys.
func (m Map) String() string {
	keys := make([]string, 0, len(m.Elements))

	for k := range m.Elements {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	out := "["

	for i, k := range keys {
		if i > 0 {
			out += " "
		}

		out += fmt.Sprintf("%s %s", k, m.Elements[k].String())
	}

	out += "]"

	return out
}

// Equal compares the map value to another.
func (m Map) Equal(other Value) bool {
	if other.Type() != MapType {
		return false
	}

	o := other.(Map)

	if len(m.Elements) != len(o.Elements) {
		return false
	}

	for k, mv := range m.Elements {
		ov, ok := o.Elements[k]
		if !ok || !mv.Equal(ov) {
			return false
		}
	}

	return true
}

// Function represents a user-defined function value (lambda/closure).
// Note: this type is only valid for the interpreted version of the language.
type Function struct {
	Env    *Environment
	Params ast.SExpr
	Body   ast.SExpr
}

// NewFunction builds a new Function.
func NewFunction(env *Environment, params ast.SExpr, body ast.SExpr) Function {
	return Function{env, params, body}
}

// Type returns the type of the function value.
func (uf Function) Type() ValueType {
	return FuncType
}

// String returns the string representation of the function value.
func (uf Function) String() string {
	return "Function()"
}

// Equal compares the function value to another.
func (uf Function) Equal(_ Value) bool {
	return false
}

// NativeFunction represents a native function value.
type NativeFunction struct {
	Value func(args ...Value) (Value, error)
}

// NewNativeFunction builds a new NativeFunction.
func NewNativeFunction(value func(args ...Value) (Value, error)) NativeFunction {
	return NativeFunction{value}
}

// Type returns the type of the native function value.
func (f NativeFunction) Type() ValueType {
	return NativeFuncType
}

// String returns the string representation of the native function value.
func (f NativeFunction) String() string {
	return "NativeFunction()"
}

// Equal compares the native function value to another.
func (f NativeFunction) Equal(_ Value) bool {
	return false
}

// RecurBindings represents a tail-call marker value for TCO.
type RecurBindings struct {
	Args []Value
}

// NewRecurBindings builds a new recurBindings.
func NewRecurBindings(args []Value) RecurBindings {
	return RecurBindings{
		Args: args,
	}
}

// Type returns the type of the recur binding value.
func (r RecurBindings) Type() ValueType {
	return RecurType
}

// String returns the string representation of the recur binding value.
func (r RecurBindings) String() string {
	return "__recur__"
}

// Equal compares the recur binding value to another.
func (r RecurBindings) Equal(_ Value) bool {
	return false
}
