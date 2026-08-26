# AGENTS.md

Tatu is an educational scripting language with S-expression syntax, first-class functions, lexical closures, tail-call optimization via `recur`, and macros. Pure Go, zero dependencies.

## Commands

- `go build cmd/*.go` — build
- `go test ./test` — run all tests
- `go test -bench=. ./test` — run benchmarks
- `go run cmd/*.go <file.tatu>` — execute a file
- `go run cmd/*.go -printTokens <file.tatu>` — dump tokens
- `go run cmd/*.go -printAST <file.tatu>` — dump AST
- `go run cmd/*.go -printInfo=false <file.tatu>` — omit the header

## Architecture

Pipeline: Source → Scanner → Parser → Expander → Analyzer → Interpreter → Result

| Package | Purpose |
|---|---|
| `pkg/ast/` | AST node definitions |
| `pkg/builder/` | File inclusion, program construction |
| `pkg/core/builtins/` | Arithmetic, comparison, I/O, type ops |
| `pkg/core/stdlib/` | math, string, vec, map, time, json, fs, regex |
| `pkg/debug/` | Error reporting with location |
| `pkg/interpreter/` | Tree-walking interpreter |
| `pkg/location/` | Source positions |
| `pkg/macro/` | Macro expander (single file, `macro.go`) |
| `pkg/parser/` | Syntactic analysis + sugar expansion |
| `pkg/pretty/` | Colored dumps of tokens, AST and errors |
| `pkg/runtime/` | Values, Environment |
| `pkg/scanner/` | Lexical analysis |
| `pkg/token/` | Token definitions |

## Embedding

`interpreter.NewInterpreter(interpreter.WithWriter(w))` redirects `print` and `println`. The host seeds globals with `SetGlobal`, reads them back with `Globals`, and calls script functions with `Call`. See the README for a full example.

## Code Style

- Comments: short, third-person, no paragraphs. `// eval evaluates an S-expression.` not `// This method evaluates an S-expression by walking the AST tree.`
- Error messages: lowercase, backtick names for builtins (`` `%s` invalid type %s``), short for interpreter (`expected BOOL, found %s`)
- Receiver: single letter (`i *Interpreter`, `e *Expander`, `p *Parser`)
- Functions: public before private, private ordered by dependency
- Guard clauses: early returns, no deep nesting

## Tests

Place in `test/<category>/`. One `.tatu` per test. Format:

```lisp
; Test <short description>

(code here)

; Expect: <result>
```

Error tests use `; Expect Error: <substring>`.

No subdirectories inside test categories, except `builtins/`, `stdlib/` and `includes/`.

A sibling `<file>.tatu.expected` is compared byte for byte against what the program printed. Optional: no file, no check.

`test/benchmark/` holds the benchmark corpus. Those programs are also regular fixtures.
