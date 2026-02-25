# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
