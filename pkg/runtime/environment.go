package runtime

import "fmt"

// Environment is the symbol table that manages variable scoping.
type Environment struct {
	record map[string]Value
	parent *Environment
}

// NewEnvironment builds a new Environment.
func NewEnvironment(record map[string]Value, parent *Environment) *Environment {
	if record == nil {
		record = make(map[string]Value)
	}

	return &Environment{
		record: record,
		parent: parent,
	}
}

// Define defines a new variable in the current scope.
func (env *Environment) Define(name string, value Value) (Value, error) {
	if _, ok := env.record[name]; ok {
		return nil, fmt.Errorf("symbol `%s` already defined", name)
	}

	env.record[name] = value

	return value, nil
}

// Assign assigns a value in the current or parent scope.
func (env *Environment) Assign(name string, value Value) bool {
	if _, ok := env.record[name]; ok {
		env.record[name] = value
		return true
	}

	if env.parent != nil {
		return env.parent.Assign(name, value)
	}

	return false
}

// Lookup looks up a variable in the current or parent scope.
func (env *Environment) Lookup(name string) (Value, bool) {
	result, ok := env.record[name]
	if ok {
		return result, ok
	}

	if env.parent != nil {
		return env.parent.Lookup(name)
	}

	return nil, false
}
