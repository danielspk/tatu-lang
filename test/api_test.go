package test

import (
	"testing"

	"github.com/danielspk/tatu-lang/pkg/builder"
	"github.com/danielspk/tatu-lang/pkg/interpreter"
	"github.com/danielspk/tatu-lang/pkg/runtime"
)

func TestGlobals(t *testing.T) {
	const source = `(var result (double answer))`

	_, program, err := builder.NewProgramBuilderWithDefaults().BuildFromSource([]byte(source), "")
	if err != nil {
		t.Fatalf("building source: %s", err)
	}

	inter := interpreter.NewInterpreter()
	inter.SetGlobal("answer", runtime.NewNumber(21))
	inter.SetGlobal("double", runtime.NewNativeFunction(fnDouble))

	if _, err := inter.EvalProgram(program, nil); err != nil {
		t.Fatalf("evaluating program: %s", err)
	}

	out := inter.Globals()

	if answer := out["answer"]; answer == nil || answer.(runtime.Number).Value != 21 {
		t.Errorf("expected seeded answer 21, found %v", answer)
	}
	if result := out["result"]; result == nil || result.(runtime.Number).Value != 42 {
		t.Errorf("expected result 42, found %v", result)
	}
}

func TestCall(t *testing.T) {
	const source = `(def add (a b) (+ a b))`

	_, program, err := builder.NewProgramBuilderWithDefaults().BuildFromSource([]byte(source), "")
	if err != nil {
		t.Fatalf("building source: %s", err)
	}

	inter := interpreter.NewInterpreter()
	if _, err := inter.EvalProgram(program, nil); err != nil {
		t.Fatalf("evaluating program: %s", err)
	}

	result, err := inter.Call("add", runtime.NewNumber(2), runtime.NewNumber(3))
	if err != nil {
		t.Fatalf("calling function: %s", err)
	}
	if result == nil || result.(runtime.Number).Value != 5 {
		t.Errorf("expected 5, found %v", result)
	}
}

// fnDouble is a native function example.
func fnDouble(args ...runtime.Value) (runtime.Value, error) {
	return runtime.NewNumber(args[0].(runtime.Number).Value * 2), nil
}
