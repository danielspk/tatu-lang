# AGENTS.md

Tatu is an educational scripting language with S-expression syntax, first-class functions, lexical closures, tail-call optimization via `recur`, and macros. Pure Go, zero dependencies.

## Commands

- `go build cmd/*.go` — build
- `go test ./test` — run all tests
- `go run cmd/*.go <file.tatu>` — execute a file
- `go run cmd/*.go -printAST <file.tatu>` — dump AST

## Architecture

Pipeline: Source → Scanner → Parser → Expander → Analyzer → Interpreter → Result

| Package | Purpose |
|---|---|
| `pkg/scanner/` | Lexical analysis |
| `pkg/parser/` | Syntactic analysis + sugar expansion |
| `pkg/ast/` | AST node definitions |
| `pkg/builder/` | File inclusion, program construction |
| `pkg/macro/` | Macro expander (single file, `macro.go`) |
| `pkg/interpreter/` | Tree-walking interpreter |
| `pkg/runtime/` | Values, Environment |
| `pkg/core/builtins/` | Arithmetic, comparison, I/O, type ops |
| `pkg/core/stdlib/` | math, string, vec, map, time, json, fs, regex |
| `pkg/debug/` | Error reporting with location |

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

No subdirectories inside test categories (except `stdlib/` and `builtins/`).
