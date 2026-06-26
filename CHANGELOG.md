# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.7.0](https://github.com/danielspk/tatu-lang/releases/tag/v0.7.0) - _2026-06-25_

### Added

- User-defined macros with positional arg binding and variadic support.
- New ellipsis symbol `...` in the parser.
- New `builder.NewProgramBuilderWithDefaults()` method.
- Expand test suite.

### Changed

- `NewInterpreter()` no longer returns an error.
- `ProgramBuilder` runs macro expansion and syntax analysis.
- `ProgramBuilder` now takes an `Expander` and an `Analyzer`.
- `Scanner` is reusable: `Scan` takes the source and filename.
- `Parser` is reusable: `Parse` takes the tokens.
- Syntax validation moved out of parsing.

### Removed

- `builder.NewDefaultScanner` and `builder.NewDefaultParser`.

## [v0.6.0](https://github.com/danielspk/tatu-lang/releases/tag/v0.6.0) - _2026-05-13_

### Changed

- `regex:matches`, `regex:find`, `regex:replace` cache compiled patterns.
- `Interpreter.Environment()` renamed to `Globals()`.
- `+` concatenates strings faster.
- `str:index` searches faster.
- `recur` allocates less memory and runs faster.

### Fixed

- `var` could shadow native functions.

## [v0.5.0](https://github.com/danielspk/tatu-lang/releases/tag/v0.5.0) - _2026-05-12_

### Changed

- `for` now creates a per-iteration binding of the loop variable.
- `Number.String()` outputs full IEEE 754 precision, matching `=`.
- `Map.String()` sorts keys alphabetically for stable output.
- `regex:find` returns `nil` on no match _(was `""`)_.
- `=` returns `false` when comparing different types, not an error.
- `=` does structural equality on vectors and maps; functions are unequal.
- `to-string` now accepts vectors, maps and functions.
- `to-bool` now accepts vectors, maps and functions.
- `to-number` no longer accepts `"Inf"` and `"NaN"` strings.
- `math:pow` no longer returns infinity or NaN on edge cases.
- `math:exp` no longer returns infinity on large arguments.
- `debug.Error.Dump` removed; error rendering moved to `pretty`.
- Test harness `; Expect Error:` matches build-time errors too.
- Remove unreachable recursive `Transform` calls in `sugar.go`.

### Fixed

- `sugar.go` errors now include line, column and file location.
- Scanner escape processing corrupting literal `\\n` into a newline.
- Scanner mis-terminating strings ending with escaped backslash `"a\\"`.
- Lambda call with wrong arity panicking with `index out of range`.
- CLI taking the last argument as filename even if it was a flag.
- Parser atom location had `End.Offset` copying `Start.Offset`.
- Scanner line count off-by-one when string ends with a newline.
- `vec:push`, `vec:pop`, `vec:delete`, `vec:concat` now mutate the receiver.
- `vec:sort` errors on mixed or unsortable types _(was indeterminate order)_.
- `map:keys` and `map:values` non-deterministic order.
- `fs:read-lines` left trailing empty line and `\r` on CRLF files.

## [v0.4.0](https://github.com/danielspk/tatu-lang/releases/tag/v0.4.0) - _2026-02-26_

### Changed

- `begin` renamed to `block`.

## [v0.3.1](https://github.com/danielspk/tatu-lang/releases/tag/v0.3.1) - _2026-02-24_

### Changed

- Update BNF.

## [v0.3.0](https://github.com/danielspk/tatu-lang/releases/tag/v0.3.0) - _2026-02-12_

### Added

- `map:get-in` stdlib function.
- `not` operator.

### Changed

- Standard library removed from environment.
- Builtins and stdlib reorganized under `pkg/core/`.
- Arithmetic, comparison, and `not` operators extracted as native functions.
- `time:format` and `time:parse` use date pattern syntax.
- `if` expression `else` branch is now optional.
- `%` operator validated by syntax analyzer.

### Fixed

- `map` key symbol resolution.
- `and`/`or` short-circuit evaluation.
- `for` syntax documentation.

## [v0.2.0](https://github.com/danielspk/tatu-lang/releases/tag/v0.2.0) - _2026-01-28_

### Added

- `Interpreter.Environment()` method to access global environment.

## [v0.1.0](https://github.com/danielspk/tatu-lang/releases/tag/v0.1.0) - _2026-01-27_

### Added

- `Interpreter.EvalProgram()` method to evaluate the AST program.

## [v0.1.0-beta](https://github.com/danielspk/tatu-lang/releases/tag/v0.1.0-beta) - _2026-01-26_

### Added

- Modulo operator `%`.
- Extends the standard language library.

## [v0.1.0-alpha](https://github.com/danielspk/tatu-lang/releases/tag/v0.1.0-alpha) - _2025-12-19_

### Added

- Initial alpha release.
