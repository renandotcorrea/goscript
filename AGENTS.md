# AGENTS.md

Guide for AI coding agents working on the goscript project — a collection of pure Go utility packages.

## Project Overview

**goscript** is a utility library providing generic, reusable Go packages for common scripting patterns:
- `slice` - Generic slice utilities with functional methods (Filter, Map, Reduce, etc.)
- `http` - Fluent HTTP client builder with chainable API
- `try` - Error handling utilities

Target audience: Go developers building scripts and small tools.

## Package Structure

```
go-scripts/
├── slice/          # Generic Slice[T] type with functional operations
│   ├── slice.go
│   └── slice_test.go
├── http/           # Fluent HTTP request builder
│   ├── http.go
│   └── http_test.go
├── try/            # Error handling utilities
│   ├── try.go
│   └── try_test.go
├── cmd/main.go     # Example usage (optional)
└── go.mod          # Module definition
```

## Setup & Development

### Prerequisites
- Go 1.25.4 or later

### Initialize
- Clone or navigate to the repository
- No external dependencies needed (standard library only)

### Build & Run
- `go build ./cmd/main.go` - Build example
- `go run ./cmd/main.go` - Run example

## Testing Instructions

### Run All Tests
```bash
go test ./...
```

### Run Tests for Specific Package
```bash
go test ./slice
go test ./http
go test ./try
```

### Run with Coverage
```bash
go test -v -cover ./...
```

### Run Specific Test
```bash
go test ./slice -run TestSliceFilter
```

**Important:** All tests must pass before committing. The agent should fix any failing tests.

## Code Style & Conventions

### Go Style
- Follow standard Go conventions (gofmt)
- Use `CamelCase` for exported types and functions
- Use `camelCase` for unexported fields and functions
- One-letter receiver names for small types (e.g., `func (s Slice[T])`)

### Generics & Type Parameters
- Use `[T]` for single type parameter (element type)
- Use `[K, V]` for key-value type parameters
- Keep type constraints simple; use `any` when no constraint needed

### Method Receivers
- Slice methods that modify: use value receiver (methods operate on copy)
- Slice methods that read: use value receiver
- Fluent builders (e.g., HttpRequest): use pointer receiver

### Comments & Documentation
- Document exported types and functions with comments
- Use example comments for complex methods (see http.go for reference)
- Keep comments concise and focused

### Composite Literal Method Calls
- When calling methods on composite literals, wrap in parentheses:
  ```go
  result := (Slice[int]{1, 2, 3}).Filter(predicate)
  ```
  This is a Go syntax requirement for method calls on types.

## Common Implementation Patterns

### Slice Functional Operations
- `Filter(predicate)` - Returns new slice, doesn't modify original
- `Map(fn)` - Returns new slice with transformed elements
- `Reduce(fn, initial)` - Accumulates to single value
- `ForEach(fn)` - Executes function for side effects only

### HTTP Request Building
- `http.Get(url)` returns `*HttpRequest`
- Chain methods: `.Headers()`, `.BodyJSON()`, etc.
- Terminal method: `.JSON(dest)` or `.Raw()` executes request
- Type assertions in tests: check `response.StatusCode` and unmarshal `Body`

### Error Handling
- Return errors explicitly; no panic unless truly exceptional
- Provide context in error messages
- Test error paths in unit tests

## Testing Guidelines

- Write tests alongside implementation (e.g., `slice_test.go` next to `slice.go`)
- Use table-driven tests for multiple cases
- Name test functions: `TestFunctionName` or `TestFunctionNameScenario`
- Test both happy path and edge cases (empty slices, nil values, etc.)
- Mock HTTP responses in http tests (use httptest package if needed)

## Commits & Pull Requests

### Commit Message Format
- Clear, descriptive message summarizing the change
- Reference the package affected: e.g., "slice: add Contains method"
- Keep messages under 72 characters for subject line

### Before Pushing
1. Run `go fmt ./...` to format code
2. Run `go test ./...` to ensure all tests pass
3. Run `go vet ./...` to check for common mistakes
4. Verify no unused imports or variables

### PR Guidelines
- One feature or bug fix per PR
- Update tests if behavior changes
- Add documentation for new exported functions
- Ensure all tests pass in CI

## Common Gotchas

### Generics Syntax
- Method calls on composite literals require parentheses: `(Slice[int]{}).IsEmpty()`
- Type parameters must be explicit in some contexts

### Slice Semantics
- Slice methods that "modify" return a new slice; originals are immutable in spirit
- Be careful with slices of pointers vs slices of values

### HTTP Testing
- Remember to close response bodies in real code (the fluent API handles this)
- Status codes aren't errors; check `response.StatusCode` explicitly

## Agent Tips

- **Exploration:** Start by reading the README.md for project goals, then package files
- **Testing first:** Always run tests after making changes
- **Consistency:** Match existing code style in the target file
- **Documentation:** New exported functions need comment documentation
- **Refactoring:** If moving code between files, verify imports still work (`go mod tidy`)
