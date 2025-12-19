# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## About the Tatu Language

Tatu is an educational programming language (toy language) developed for general scripting or embedding within Go applications. It is a functional language with the following features:

- S-expression based syntax
- Everything is an expression (no statements)
- Last evaluated expression always returns its value
- First-class functions with lexical closures
- Lambda functions (anonymous functions)
- Not a Lisp dialect (no metaprogramming support)
- Zero external dependencies (pure Go implementation)

## Development Commands

### Building
- `go build cmd/*.go` - Build the main executable for current platform
- `make` - Build for all platforms (Linux, macOS, Windows for amd64 and arm64) and clean dist directory
- `make build-linux` - Build for Linux (amd64 and arm64)
- `make build-mac` - Build for macOS (amd64 and arm64)
- `make build-windows` - Build for Windows (amd64 and arm64)
- `make clean` - Remove dist directory

### Testing
- `go test ./test` - Run all tests (evaluates .tatu files in test directory)
- Test files use special comment format: `; Expect: <expected_result>`
- Tests are organized in subdirectories: addition, assignments, atom, conditionals, functions, lambdas, loops, recursion, scopes, variables, etc.

### Running
- `go run cmd/*.go` - Execute a Tatu source file or start interactive mode
- `go run cmd/*.go <file.tatu>` - Execute a Tatu source file
- `go run cmd/*.go -printTokens <file.tatu>` - Print generated tokens
- `go run cmd/*.go -printAST <file.tatu>` - Print generated AST
- `go run cmd/*.go -printBytecode <file.tatu>` - Print bytecode (WIP feature)

## Architecture

### Main Components

1. **Scanner** (`pkg/scanner/`) - Lexical analysis
   - Converts source code into tokens
   - Handles UTF-8 strings, numbers, symbols, booleans, nil, comments
   - Supports escape sequences in strings: `\\`, `\"`, `\n`, `\r`, `\t`

2. **Parser** (`pkg/parser/`) - Syntactic analysis
   - Converts tokens into Abstract Syntax Tree (AST)
   - Handles atoms and lists (S-expressions)
   - Includes sugar syntax support and semantic analysis

3. **AST** (`pkg/ast/`) - Abstract Syntax Tree representation
   - `SExpr` interface with different implementations:
     - `NumberExpr`, `StringExpr`, `BoolExpr`, `NilExpr`, `SymbolExpr` (atoms)
     - `ListExpr` (lists/S-expressions)
   - Simple unified expression type for all language constructs
   - Each expression includes location information for error reporting

4. **Token** (`pkg/token/`) - Token definitions and types
   - Defines token types: Number, String, Bool, Nil, Symbol, LParen, RParen, EOF, Illegal

5. **Location** (`pkg/location/`) - Position and location tracking
   - Tracks file, line, column, and offset for tokens and AST nodes
   - Used for accurate error reporting

6. **Builder** (`pkg/builder/`) - Program construction
   - `ProgramBuilder` - Orchestrates scanning and parsing
   - Handles file inclusion and module resolution
   - Provides abstraction over Scanner and Parser interfaces

7. **Runtime** (`pkg/runtime/`) - Runtime value system (shared between interpreter and VM)
   - **Value** - Runtime value representation (Number, String, Bool, Nil, Function, CoreFunction, RecurBindings)
   - **Environment** - Variable scoping and closure support with lexical scoping
   - Shared types used by both tree-walking interpreter and VM

8. **Standard Library** (`pkg/stdlib/`) - Built-in core functions
   - **math.go** - Mathematical functions (sqrt, abs, pow, etc.)
   - **helpers.go** - Argument validation helpers for stdlib functions
   - Pure functions that work with runtime.Value types
   - Shared between interpreter and VM implementations

9. **Interpreter** (`pkg/interpreter/`) - Tree-walking interpreter (currently active)
   - **Interpreter** - Evaluates AST directly using runtime values
   - Implements tail-call optimization with `recur` special form
   - Registers stdlib functions in global environment on initialization

10. **VM** (`pkg/vm/`) - Virtual machine for bytecode execution (WIP)
    - **Compiler** - Bytecode compilation (work in progress)
    - **VM** - Stack-based virtual machine with 512 element stack limit
    - **Opcodes** - Bytecode operation definitions
    - Will share stdlib functions with interpreter

11. **Pretty** (`pkg/pretty/`) - Output formatting
    - Colored output formatting for REPL and error messages
    - AST and token pretty-printing

12. **Debug** (`pkg/debug/`) - Error handling
    - Structured error reporting with location information

### Execution Pipeline

**Current (Interpreter-based):**
1. Source code → Scanner → Tokens
2. Tokens → Parser → AST
3. AST → Builder (resolves includes) → Final AST
4. AST → Interpreter → Result

**Future (VM-based - WIP):**
1. Source code → Scanner → Tokens
2. Tokens → Parser → AST
3. AST → Builder (resolves includes) → Final AST
4. AST → Compiler → Bytecode
5. Bytecode → Virtual Machine → Result

### Language Features

#### Data Types
- **Numbers**: Integers and floating-point numbers (e.g., `42`, `3.14`, `-10`)
- **Strings**: UTF-8 strings with escape sequences (e.g., `"hello"`, `"line\nnew"`)
- **Booleans**: `true` and `false`
- **Nil**: Represents absence of value
- **Symbols**: Identifiers and operators

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
- **Blocks**:
  - `(begin expr1 expr2 ...)` - Evaluate expressions sequentially, return last result
  - Creates new lexical scope

#### Built-in Operators
- **Arithmetic**:
  - `+` - Addition (supports numbers and string concatenation)
  - `-` - Subtraction
  - `*` - Multiplication
  - `/` - Division
- **Comparison**: `>`, `>=`, `<`, `<=`, `=`
- **Logical**: `and`, `or`
- **I/O**: `print` - Print values to stdout
- **Math Library** (namespace: `math:`):
  - `math:sqrt` - Square root (e.g., `(math:sqrt 16)` → `4`)
  - `math:abs` - Absolute value (e.g., `(math:abs -5)` → `5`)
  - `math:pow` - Power function (e.g., `(math:pow 2 3)` → `8`)
- **Module System**: `include` - Include external files (WIP)

#### Special Features
- All expressions return a value (no statements)
- Last expression in a block is the return value
- Functions are first-class values (can be assigned, passed, and returned)
- Lexical closures (functions capture their defining environment)
- UTF-8 native support for strings and identifiers
- Tail recursion support (pending optimization - design in TCO section below)

### Test Structure

Tests are located in the `/test` directory as `.tatu` files with expected results in comments:
- Each test file should end with `; Expect: expected_output`
- Test runner (`test/evaluator_test.go`) reads all `.tatu` files and validates output
- **109+ comprehensive test files** organized by feature:
  - **Operators**: addition, subtraction, multiplication, division, logical (and/or), concatenation
  - **Language constructs**: atom, assignments, conditionals, functions, lambdas, loops, recursion, scopes, variables
  - **Advanced**: closures, IILE (Immediately Invoked Lambda Expressions), tail recursion
  - **Standard library**: math (sqrt, abs, pow), math operations
  - **Special**: list, includes

### Implementation Details

**Pure Go Implementation:**
- Zero external dependencies (only Go standard library)
- No CGO, no external bindings
- Cross-platform compilation for 6 targets: Linux, macOS, Windows (amd64 & arm64)
- Single binary distribution with no runtime dependencies

**Error Handling:**
- Structured error reporting with file, line, and column tracking
- Source code context display with visual markers (`↑` pointing to error location)
- Detailed error messages for debugging

**Developer Tools:**
- Debug flags: `-printTokens`, `-printAST`, `-printBytecode`
- Pretty-printed output with ANSI colors
- Declarative testing format (`.tatu` files with `; Expect:` comments)

**Syntactic Sugar:**
- `def` expands to `var` + `lambda`
- Unary negation: `(- 5)` → -5
- Variadic operators: `+`, `*`, `and`, `or` accept multiple arguments

### Code Quality Notes

From the codebase comments and TODOs:
- AST changes resulted in extensive type casting that needs cleanup
- Error messages in interpreter need improvement
- Semantic analysis phase needs implementation to validate structures before evaluation
- Many eval functions assume correct format without validation
- Include/module system is incomplete
- Virtual machine and compiler are work in progress

### Architecture Decisions

**Runtime/Stdlib Separation:**
- `pkg/runtime/` contains shared value types (Number, String, Bool, Nil, Function, CoreFunction, Environment)
- `pkg/stdlib/` contains pure functions that operate on runtime values
- Both interpreter and VM (future) share these components to avoid code duplication
- CoreFunction signature: `func(args ...runtime.Value) (runtime.Value, error)` - no Environment parameter for purity
- Function type currently contains `ast.SExpr` (interpreter-specific); will need refactoring when VM is complete

## Future Implementation: Self-Tail-Recursion Optimization (TCO)

### Current Problem

The interpreter currently creates a new Go stack frame for every recursive function call in `pkg/interpreter/interpreter.go`:

```go
return i.evalBody(value.(Function).body, activationEnv)
```

This causes **stack overflow** with deep recursions (e.g., `factorial(10000)`).

### Design: Simple Self-Tail-Recursion Optimization

**Goal**: Optimize tail-recursive calls where a function calls itself in tail position (last operation before return).

**Key Insight**: Replace recursive calls with iteration using a trampoline loop - reuse the same Go stack frame instead of creating new ones.

### What is Tail Position?

A function call is in **tail position** when it's the last operation before returning:

✅ **Tail Position (optimizable):**
```lisp
(def factorial-tail (n acc)
  (if (= n 0)
    acc
    (factorial-tail (- n 1) (* acc n))))  ; Last operation - no work after
```

❌ **NOT Tail Position:**
```lisp
(def factorial (n)
  (if (= n 1)
    1
    (* n (factorial (- n 1)))))  ; Multiplication happens AFTER recursive call
```

### Implementation Plan

**Note:** Tail-call optimization is already implemented using the `recur` special form in `pkg/interpreter/interpreter.go`.

**Alternative approach (not yet implemented):** Automatic TCO detection

**Changes needed:**

1. **Add trampoline loop** around the function call logic
2. **Detect self-recursion** in tail position
3. **Reuse parameters** instead of creating new stack frame
4. **Extract tail expression** from function body (handle `begin`, `if`)

### Pseudo-code Implementation

```go
func (i *Interpreter) evalCallFunction(expr ast.SExpr, env *Environment) (Value, error) {
    exprList := expr.(*ast.ListExpr)

    // TRAMPOLINE LOOP: Replace recursion with iteration
    for {
        value, err := i.Eval(exprList.List[0], env)
        if err != nil {
            return nil, err
        }

        // Handle core functions (no change needed)
        if value.Type() == CoreFuncType {
            // ... existing code ...
            return value.(CoreFunction).value(env, results...)
        }

        // Handle lambda/user functions
        if value.Type() == FuncType {
            fn := value.(Function)

            // Evaluate arguments
            var args []Value
            for _, e := range exprList.List[1:] {
                result, err := i.Eval(e, env)
                if err != nil {
                    return nil, err
                }
                args = append(args, result)
            }

            // Create activation environment
            activationRecord := make(map[string]Value)
            params := fn.params.(*ast.ListExpr).List
            for idx, p := range params {
                activationRecord[p.(*ast.SymbolExpr).Symbol] = args[idx]
            }
            activationEnv := NewEnvironment(activationRecord, fn.env)

            // Extract tail expression from body
            tailExpr := extractTailExpression(fn.body, activationEnv)

            // Check if tail expression is self-recursive call
            if isSelfRecursiveCall(tailExpr, exprList.List[0], activationEnv) {
                // OPTIMIZE: Reuse current stack frame
                exprList = tailExpr.(*ast.ListExpr)
                env = activationEnv
                continue  // Loop back instead of recursing
            }

            // Not tail-recursive: evaluate normally
            return i.evalBody(fn.body, activationEnv)
        }

        return nil, i.error("not a function", expr.Location())
    }
}
```

### Helper Functions Needed

```go
// extractTailExpression extracts the expression in tail position
func extractTailExpression(body ast.SExpr, env *Environment) ast.SExpr {
    if body.Kind() != ast.ListKind {
        return body
    }

    list := body.(*ast.ListExpr)
    if len(list.List) == 0 {
        return body
    }

    if list.List[0].Kind() != ast.SymbolKind {
        return body
    }

    symbol := list.List[0].(*ast.SymbolExpr).Symbol

    // Handle (begin ...) - tail expression is the last one
    if symbol == "begin" {
        if len(list.List) > 1 {
            // Evaluate all expressions except last
            for _, e := range list.List[1:len(list.List)-1] {
                i.Eval(e, env)
            }
            // Recurse on last expression
            return extractTailExpression(list.List[len(list.List)-1], env)
        }
    }

    // Handle (if cond then else) - tail expression is in then/else branch
    if symbol == "if" && len(list.List) == 4 {
        condValue, _ := i.Eval(list.List[1], env)
        if condValue.Type() == BoolType {
            if condValue.(Bool).value {
                return extractTailExpression(list.List[2], env) // then branch
            } else {
                return extractTailExpression(list.List[3], env) // else branch
            }
        }
    }

    return body
}

// isSelfRecursiveCall checks if tailExpr is a call to the same function
func isSelfRecursiveCall(tailExpr ast.SExpr, originalFunc ast.SExpr, env *Environment) bool {
    if tailExpr.Kind() != ast.ListKind {
        return false
    }

    tailList := tailExpr.(*ast.ListExpr)
    if len(tailList.List) == 0 {
        return false
    }

    // Both must be symbols referring to the same function
    if tailList.List[0].Kind() != ast.SymbolKind || originalFunc.Kind() != ast.SymbolKind {
        return false
    }

    tailSym := tailList.List[0].(*ast.SymbolExpr).Symbol
    origSym := originalFunc.(*ast.SymbolExpr).Symbol

    // Check if symbols are the same
    return tailSym == origSym
}
```

### Benefits

✅ **No new types needed**: Works with existing `Value` types
✅ **Localized changes**: Only modify `evalCallFunction` and add 2 helper functions
✅ **Covers 95% of cases**: `factorial`, `fibonacci`, `sum`, etc.
✅ **No stack overflow**: Can handle recursions of millions of levels
✅ **Same performance as loops**: Tail recursion as fast as `while`
✅ **Simple to test and debug**: Straightforward logic

### Limitations

⚠️ **Only self-recursion**: Doesn't optimize mutual recursion (function A calls B calls A)
⚠️ **Manual tail form**: User must write functions in tail-recursive style
⚠️ **Evaluates conditions twice**: `if` condition evaluated in `extractTailExpression` and main eval

### Testing Strategy

Create tests in `test/recursion/`:

1. **`factorial_tail.tatu`** - Tail-recursive factorial with large numbers
```lisp
(def factorial-tail (n acc)
  (if (= n 0) acc (factorial-tail (- n 1) (* acc n))))
(factorial-tail 1000 1)
```

2. **`fibonacci_tail.tatu`** - Tail-recursive fibonacci
```lisp
(def fib-tail (n a b)
  (if (= n 0) a (fib-tail (- n 1) b (+ a b))))
(fib-tail 50 0 1)
```

3. **`sum_tail.tatu`** - Tail-recursive sum
```lisp
(def sum-tail (n acc)
  (if (= n 0) acc (sum-tail (- n 1) (+ acc n))))
(sum-tail 10000 0)
```

4. **`deep_recursion.tatu`** - Test with very deep recursion
```lisp
(def countdown (n)
  (if (= n 0) n (countdown (- n 1))))
(countdown 100000)
```

### Estimated Effort

- **Implementation**: ~2-3 hours
- **Testing**: ~1 hour
- **Total**: Half day of work

### Future Extensions

Once this works, can extend to:
1. **Mutual recursion** (function A → B → A)
2. **Tail calls between different functions**
3. **More sophisticated tail position detection**