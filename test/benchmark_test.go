package test

import (
	"testing"

	"github.com/danielspk/tatu-lang/pkg/builder"
	"github.com/danielspk/tatu-lang/pkg/interpreter"
)

func BenchmarkSum(b *testing.B) {
	source := `
(def sum (n)
  (if (= n 0)
    0
    (+ n (sum (- n 1)))))

(sum 1000)
`
	runTestCode(b, source)
}

func BenchmarkSumWithTCO(b *testing.B) {
	source := `
(def sum (n acc)
  (if (= n 0)
    acc
    (recur (- n 1) (+ acc n))))

(sum 1000 0)
`
	runTestCode(b, source)
}

func runTestCode(b *testing.B, source string) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		progBuilder := builder.NewProgramBuilder(builder.NewDefaultScanner(), builder.NewDefaultParser())
		_, ast, err := progBuilder.BuildFromSource([]byte(source), "")
		if err != nil {
			b.Fatalf("building source: %v", err)
		}

		inter, err := interpreter.NewInterpreter()
		if err != nil {
			b.Fatalf("creating interpreter: %v", err)
		}

		for _, expr := range ast.Program {
			_, err = inter.Eval(expr, nil)
			if err != nil {
				b.Fatalf("evaluating program: %v", err)
			}
		}
	}
}
