# CLAUDE.md

This file provides guidance when working with the Tatu programming language codebase.

## About Tatu

Tatu is an educational functional programming language for general scripting or embedding in Go applications.

**Key Features:**
- S-expression syntax (Lisp-like but not a Lisp)
- Everything is an expression
- First-class functions with lexical closures
- Tail-call optimization via `recur`
- Pure Go implementation (zero external dependencies)

## Development Commands

### Building
- `go build cmd/*.go` - Build the main executable for current platform
- `make` - Build for all platforms (Linux, macOS, Windows for amd64 and arm64) and clean dist directory
- `make build-linux` - Build for Linux (amd64 and arm64)
- `make build-mac` - Build for macOS (amd64 and arm64)
- `make build-windows` - Build for Windows (amd64 and arm64)
- `make clean` - Remove dist directory

### Testing
- `go test ./test` - Run all tests (300+ test files)
- **Success tests**: `; Expect: <expected_result>`
- **Error tests**: `; Expect Error: <optional_error_substring>`
- Tests organized by feature in `test/` subdirectories

### Running
- `go run cmd/*.go` - Execute a Tatu source file or start interactive mode
- `go run cmd/*.go <file.tatu>` - Execute a Tatu source file
- `go run cmd/*.go -printTokens <file.tatu>` - Print generated tokens
- `go run cmd/*.go -printAST <file.tatu>` - Print generated AST
- `go run cmd/*.go -printBytecode <file.tatu>` - Print bytecode (WIP feature)

## Architecture

### Core Packages

- **`pkg/scanner/`** - Lexical analysis (source → tokens)
- **`pkg/parser/`** - Syntactic analysis (tokens → AST)
- **`pkg/ast/`** - AST node definitions (NumberExpr, StringExpr, ListExpr, etc.)
- **`pkg/builder/`** - Program construction, handles file inclusion
- **`pkg/runtime/`** - Runtime values (Number, String, Bool, Nil, Function, Vector, Map) and Environment
- **`pkg/interpreter/`** - Tree-walking interpreter with tail-call optimization
- **`pkg/stdlib/`** - Standard library modules (math, string, vector, map, time, json, file_system, casting, args)
- **`pkg/vm/`** - Virtual machine (work in progress)
- **`pkg/pretty/`** - Colored output and formatting
- **`pkg/debug/`** - Error reporting with location tracking

### Execution Pipeline

Source → Scanner → Parser → Builder → Interpreter → Result

### Language Features

#### Data Types
- **Numbers**: Integers and floating-point numbers (e.g., `42`, `3.14`, `-10`)
- **Strings**: UTF-8 strings with escape sequences (e.g., `"hello"`, `"line\nnew"`)
- **Booleans**: `true` and `false`
- **Nil**: Represents absence of value
- **Symbols**: Identifiers and operators
- **Vectors**: Ordered collections (e.g., `(vector 1 2 3)`)
- **Maps**: Key-value pairs (e.g., `(map "name" "John" "age" 30)`)

#### Core Language Constructs
- **Lists**: S-expression syntax `(operator operand1 operand2 ...)`
- **Variables**:
  - `(var name value)` - Define a new variable in current scope
  - `(set name value)` - Assign to an existing variable
- **Functions**:
  - `(lambda (param1 param2 ...) body)` - Anonymous function (closure)
  - User-defined functions are closures with lexical scoping
- **Control Flow**:
  - `(if condition true-expr false-expr)` - Conditional expression
  - `(while condition body)` - Loop while condition is true
  - `(for init condition body increment)` - For loop (syntactic sugar)
  - `(switch (cond1 result1) (cond2 result2) ... (default default-result))` - Pattern matching (syntactic sugar)
- **Blocks**:
  - `(begin expr1 expr2 ...)` - Evaluate expressions sequentially, return last result
  - Creates new lexical scope
- **Recursion**:
  - `(recur arg1 arg2 ...)` - Tail-recursive call (see Tail-Call Optimization section)

#### Built-in Operators
- **Arithmetic**:
  - `+` - Addition (supports numbers and string concatenation)
  - `-` - Subtraction
  - `*` - Multiplication
  - `/` - Division
- **Comparison**: `>`, `>=`, `<`, `<=`, `=`
- **Logical**: `and`, `or`
- **I/O**: `print` - Print values to stdout

#### Standard Library

All stdlib functions follow the pattern `namespace:function-name`.

- **`math:`** - Mathematical operations (sqrt, abs, pow, sin, cos, tan, log, exp, floor, ceil, round, min, max, between, rand, pi, e)
- **`str:`** - String operations (len, concat, split, join, slice, contains, starts, ends, index, upper, lower, trim, replace, repeat, reverse)
- **`vec:`** - Vector operations (len, get, set, push, pop, concat, slice, find, contains, delete, reverse, sort)
- **`map:`** - Map operations (len, get, set, has, delete, keys, values, merge)
- **`time:`** - Time operations (now, unix, year, month, day, hour, minute, second, format, parse, add, sub, diff, is-leap)
- **`json:`** - JSON encoding/decoding (encode, decode)
- **`fs:`** - File system operations (read, write, append, delete, exists, list, mkdir, move, is-dir, basename, size, temp-dir, read-lines)
- **`regex:`** - Regular expressions (matches, find, replace)
- **Type conversion** - to-string, to-number, to-bool
- **Type checking** - is-bool, is-number, is-int, is-string, is-vector, is-map, is-nil, is-function
- **Args** - args (command-line arguments)

#### Module System
- `(include "path/to/file.tatu")` - Include and evaluate external files with circular dependency prevention

#### Syntactic Sugar
- `def` - Function definition: `(def name (params) body)` → `(var name (lambda (params) body))`
- `switch` - Pattern matching (expands to nested `if`)
- `for` - Loop construct: `(for init cond body inc)` → `(begin init (while cond (begin body inc)))`
- Unary negation: `(- 5)` → `-5`
- Variadic operators: `+`, `*`, `and`, `or`

## Tail-Call Optimization

Use `(recur arg1 arg2 ...)` for tail-recursive calls to prevent stack overflow.

**Example:**
```lisp
(def factorial (n acc)
  (if (= n 0)
    acc
    (recur (- n 1) (* acc n))))

(factorial 10000 1)  ; No stack overflow
```

**Implementation:** The interpreter uses a trampoline loop in `evalCallFunction()` that detects `RecurBindings` and reuses the same stack frame instead of creating new Go stack frames.

**Limitations:** Must be used explicitly (not automatic), only works in tail position, cannot optimize mutual recursion between different functions.

## Contributing Guidelines

### Adding Tests
- Place tests in appropriate `test/` subdirectories
- Use `; Expect: <value>` for success tests
- Use `; Expect Error: <substring>` for error tests (substring is optional)
- Test files should be simple and focused on single functionality

### Adding Standard Library Functions
1. Add function to appropriate file in `pkg/stdlib/` (or create new module)
2. Use `expectArgs()` helper for argument validation
3. Register function in `Register<Module>()` function with namespace
4. Register module in `pkg/interpreter/interpreter.go:NewInterpreter()`
5. Add tests in `test/stdlib/<module>/`
6. Update CLAUDE.md stdlib section