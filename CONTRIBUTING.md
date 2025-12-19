# Contributing to Tatu-Lang

Thank you for your interest in contributing to _Tatu_! This document provides guidelines for contributing to the project.

## How to Contribute

### Reporting Bugs

Before creating a bug report:
- Check the [issue tracker](https://github.com/danielspk/tatu-lang/issues) to avoid duplicates
- Use the latest version from the `master` branch

When reporting a bug, include:
- Tatu version (`tatu --version`)
- Go version (`go version`)
- Operating system and architecture
- Minimal code example that reproduces the issue
- Expected vs actual behavior

### Suggesting Features

Feature suggestions are welcome! Please:
- Check existing issues to avoid duplicates
- Clearly describe the use case and benefits
- Provide examples of the proposed syntax or behavior

### Pull Requests

1. **Fork** the repository and create your branch from `master`
2. **Make your changes** following the code style guidelines
3. **Add tests** for any new functionality
4. **Run tests** to ensure everything passes: `go test ./test`
5. **Update documentation** if needed (README.md, CLAUDE.md, comments)
6. **Commit** using [Conventional Commits](#commit-message-format)
7. **Submit** the pull request with a clear description

## Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/tatu-lang.git
cd tatu-lang

# Build
go build cmd/*.go

# Run tests
go test ./test

# Run a Tatu program
go run cmd/*.go examples/playground.tatu
```

## Code Style

### Go Code
- Follow standard Go conventions ([Effective Go](https://golang.org/doc/effective_go))
- Run `go fmt` before committing
- Keep functions small and focused
- Add comments for exported types and functions
- Maintain consistency with existing code

### Tatu Test Files
- Use the `.tatu` extension
- Include expected output: `; Expect: <result>`
- Organize tests by feature in subdirectories under `test/`

## Commit Message Format

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>: <description>

[optional body]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Maintenance tasks

**Examples:**
```
feat: add tail-call optimization for self-recursion
fix: correct error message for non-callable expressions
docs: update README with installation instructions
test: add tests for logical operators (and, or)
```

## Testing

- All new features must include tests
- Place test files in appropriate subdirectories under `test/`
- Ensure all tests pass before submitting PR: `go test ./test`
- Test coverage is highly valued

## Questions?

Feel free to open an issue for questions or discussions.

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (see [LICENSE](LICENSE)).
