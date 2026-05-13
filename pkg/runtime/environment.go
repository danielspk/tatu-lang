package runtime

import "fmt"

// Binding represents an entry in the symbol table.
type Binding struct {
	Value  Value
	Native bool
}

// Environment is the symbol table that manages variable scoping.
type Environment struct {
	record map[string]Binding
	parent *Environment
}

// NewEnvironment builds a new Environment.
func NewEnvironment(record map[string]Binding, parent *Environment) *Environment {
	if record == nil {
		record = make(map[string]Binding)
	}

	return &Environment{
		record: record,
		parent: parent,
	}
}

// Define defines a new user binding in the current scope.
func (env *Environment) Define(name string, value Value) (Value, error) {
	if env.hasNative(name) {
		return nil, fmt.Errorf("cannot redefine native `%s`", name)
	}

	if _, ok := env.record[name]; ok {
		return nil, fmt.Errorf("symbol `%s` already defined", name)
	}

	env.record[name] = Binding{Value: value}

	return value, nil
}

// DefineNative defines a runtime-provided binding in the current scope.
func (env *Environment) DefineNative(name string, value Value) {
	env.record[name] = Binding{Value: value, Native: true}
}

// Assign assigns a value in the current or parent scope.
func (env *Environment) Assign(name string, value Value) error {
	if b, ok := env.record[name]; ok {
		if b.Native {
			return fmt.Errorf("cannot assign to native `%s`", name)
		}

		env.record[name] = Binding{Value: value}

		return nil
	}

	if env.parent != nil {
		return env.parent.Assign(name, value)
	}

	return fmt.Errorf("undefined variable `%s`", name)
}

// Lookup looks up a binding in the current or parent scope.
func (env *Environment) Lookup(name string) (Value, bool) {
	if b, ok := env.record[name]; ok {
		return b.Value, true
	}

	if env.parent != nil {
		return env.parent.Lookup(name)
	}

	return nil, false
}

// Variables returns the user-defined variables in this scope.
func (env *Environment) Variables() map[string]Value {
	out := make(map[string]Value, len(env.record))

	for name, b := range env.record {
		if !b.Native {
			out[name] = b.Value
		}
	}

	return out
}

// hasNative checks for a native binding in the current or parent scope.
func (env *Environment) hasNative(name string) bool {
	if b, ok := env.record[name]; ok && b.Native {
		return true
	}

	if env.parent != nil {
		return env.parent.hasNative(name)
	}

	return false
}
