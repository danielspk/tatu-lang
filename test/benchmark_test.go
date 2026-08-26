package test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/danielspk/tatu-lang/pkg/builder"
	"github.com/danielspk/tatu-lang/pkg/interpreter"
)

const benchmarkDir = "benchmark"

func BenchmarkSum(b *testing.B) {
	runTestFile(b, "sum")
}

func BenchmarkSumWithTCO(b *testing.B) {
	runTestFile(b, "sum_tco")
}

func BenchmarkStringConcat(b *testing.B) {
	runTestFile(b, "string_concat")
}

func BenchmarkRegexInLoop(b *testing.B) {
	runTestFile(b, "regex_loop")
}

func BenchmarkMacroExpansion(b *testing.B) {
	runTestFile(b, "macro_expansion")
}

func BenchmarkMacroEquivalent(b *testing.B) {
	runTestFile(b, "macro_equivalent")
}

func BenchmarkBigProgram(b *testing.B) {
	runTestFile(b, "big")
}

func BenchmarkMandel(b *testing.B) {
	runTestFile(b, "mandel")
}

func BenchmarkDragon(b *testing.B) {
	runTestFile(b, "dragon")
}

func runTestFile(b *testing.B, name string) {
	source, err := os.ReadFile(filepath.Join(benchmarkDir, name+".tatu"))
	if err != nil {
		b.Fatalf("reading source: %v", err)
	}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		progBuilder := builder.NewProgramBuilderWithDefaults()

		_, ast, err := progBuilder.BuildFromSource(source, "")
		if err != nil {
			b.Fatalf("building source: %v", err)
		}

		inter := interpreter.NewInterpreter(interpreter.WithWriter(io.Discard))

		_, err = inter.EvalProgram(ast, nil)
		if err != nil {
			b.Fatalf("evaluating program: %v", err)
		}
	}
}
