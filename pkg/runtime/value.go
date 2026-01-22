// Package runtime implements the runtime values and scopes system.
package runtime

import (
	"fmt"
	"math"
	"strconv"

	"github.com/danielspk/tatu-lang/pkg/ast"
)

// ValueType represents the type of value.
type ValueType string

// Value types.
const (
	NumberType   ValueType = "NUMBER"
	StringType   ValueType = "STRING"
	BoolType     ValueType = "BOOL"
	NilType      ValueType = "NIL"
	VectorType   ValueType = "VECTOR"
	MapType      ValueType = "MAP"
	FuncType     ValueType = "FUNC"
	CoreFuncType ValueType = "CORE_FUNC"
	RecurType    ValueType = "RECUR"
)

// Value represents a value interface.
type Value interface {
	Type() ValueType
	String() string
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
	if n.Value == 0 {
		return "0"
	}

	if n.Value == math.Trunc(n.Value) {
		return fmt.Sprintf("%.0f", n.Value)
	}

	formatted := fmt.Sprintf("%.10f", n.Value)
	value, _ := strconv.ParseFloat(formatted, 64)

	return fmt.Sprintf("%g", value)
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

// Vector represents a value of a vector type.
type Vector struct {
	Elements []Value
}

// NewVector builds a new Vector.
func NewVector(elements []Value) Vector {
	return Vector{elements}
}

// Type returns the type of the vector value.
func (v Vector) Type() ValueType {
	return VectorType
}

// String returns the string representation of the vector value.
func (v Vector) String() string {
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

// String returns the string representation of the map value.
func (m Map) String() string {
	out := "["

	i := 0
	for k, v := range m.Elements {
		if i > 0 {
			out += " "
		}
		out += fmt.Sprintf("%s %s", k, v.String())
		i++
	}

	out += "]"

	return out
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

// CoreFunction represents a built-in/native function value.
type CoreFunction struct {
	Value func(args ...Value) (Value, error)
}

// NewCoreFunction builds a new CoreFunction.
func NewCoreFunction(value func(args ...Value) (Value, error)) CoreFunction {
	return CoreFunction{value}
}

// Type returns the type of the core function value.
func (f CoreFunction) Type() ValueType {
	return CoreFuncType
}

// String returns the string representation of the core function value.
func (f CoreFunction) String() string {
	return "CoreFunction()"
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
