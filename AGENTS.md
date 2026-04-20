# AGENTS.md

Guide for AI coding agents working on the goscript project — a collection of pure Go utility packages for scripting.

## Project Overview

**goscript** provides:
- `slice` - Generic `Slice[T]` with functional operations (Filter, Map, Reduce, FlatMap, Chunk)
- `http` - Fluent HTTP request builder with chainable API (QueryParams, Retry, Timeout)
- `try` - Error handling utilities (Try, Try1, Try2, etc.)
- `env` - Environment variable helpers (MustGet, GetOr with type conversion, LoadFile)
- `file` - JSON read/write utilities (ReadJson, WriteJson)

Target audience: Go developers building scripts and small tools.

**Documentation:** All packages include godoc comments and canonicalized examples on [pkg.go.dev](https://pkg.go.dev/github.com/renandotcorrea/goscript).

## Setup & Development

### Prerequisites
- Go 1.25.4 or later
- No external dependencies (standard library only)

### Quick Start
```bash
go test ./...        # Run all tests
go fmt ./...         # Format code
go vet ./...         # Lint checks
```

## Code Standards

### Style
- Follow Go conventions: `gofmt`, `CamelCase` exports, `camelCase` unexported
- One-letter receiver names for small types: `func (s Slice[T])`
- Comments for all exported functions and types (checked by godoc)

### Documentation
- Package-level comments exist for all packages (first line summarizes purpose)
- Exported functions have doc comments (first sentence is summary)
- Canonical examples exist as `ExampleXxx()` test functions in `*_test.go`
- See existing code in any package for patterns

### Documentation Guardrails (Required)
- When adding a new package or exported symbol, follow [Go doc comment conventions](https://go.dev/doc/comment)
- Package comment must start with `Package <name> ...` and describe package purpose, scope, and key APIs
- Every exported `type`, `func`, `const`, and `var` must have a doc comment starting with the symbol name
- Boolean-returning docs should prefer `reports whether ...` wording
- Document relevant edge cases and error behavior (what can fail and when)
- Keep comments implementation-agnostic; describe caller-visible behavior, not internal algorithm details
- Use Go doc links when helpful, e.g., `[io.Reader]`, `[json.Unmarshal]`, `[Slice.IsEmpty]`
- Prefer short, complete sentences and stable wording suitable for `go doc` and `pkg.go.dev`
- Add `ExampleXxx()` tests for non-obvious behavior and keep examples runnable and realistic

**Doc comment templates:**

```go
// Package mypkg provides ...
//
// It is intended for ...
// The main entry points are [Foo] and [Bar].
package mypkg
```

```go
// Foo does ...
// It returns ...
// It reports an error if ...
func Foo(...) (..., error)
```

```go
// IsReady reports whether ...
func IsReady(...) bool
```

**Definition of done for docs (new package/function):**
1. Package comment added/updated with purpose and scope
2. All exported symbols documented with caller-focused behavior
3. Error and edge-case behavior documented where relevant
4. `ExampleXxx()` tests added/updated for key API usage
5. Verify rendering with `go doc <module/package>` and sanity-check on pkgsite format

### Generics
- Use `[T]` for single type parameter, `[K, V]` for key-value
- Avoid complex constraints; prefer `any`
- Composite literal method calls need parentheses: `(Slice[int]{}).IsEmpty()`

### Method Semantics
- Slice methods use value receivers (see slice.go)
- HTTP builder methods use pointer receivers for chaining (see http.go)
- All methods return new data; originals unchanged (immutable-spirit)

## Testing & Commit Workflow

### Before Pushing
```bash
go test ./...        # All tests pass
gofmt -w ./...       # Format code
go vet ./...         # Check for mistakes
```

### Commit Messages
- Clear, descriptive; reference package: e.g., "slice: add Chunk method"
- Keep subject line under 72 characters

## Common Patterns

**Adding a New Function:**
1. Write function in package `*.go` file
2. Apply all rules in **Documentation Guardrails (Required)**
3. Complete **Definition of done for docs (new package/function)**
4. Add unit tests in `*_test.go` with edge cases
5. Update `README.md` — add the function to the relevant package example block if it meaningfully changes how the package is used
6. Run tests, format, commit

**Adding a New Package:**
1. Create the package directory and implement the package
2. Apply all rules in **Documentation Guardrails (Required)**
3. Complete **Definition of done for docs (new package/function)**
4. Add unit tests and `ExampleXxx()` functions
5. Update `README.md`:
   - Add a row to the packages table (with link to pkg.go.dev)
   - Add a new `###` section with a usage example and "Full API reference →" link
6. Update `.github/prompts/goscript.prompt.md`:
   - Add a `#fetch https://pkg.go.dev/github.com/renandotcorrea/goscript/<package>` line to the fetch list
7. Update `AGENTS.md` → Project Overview package list
8. Run tests, format, commit

**HTTP Fluent Chains:**
```go
Get(url).
  QueryParams(map[string]string{"page": "1"}).
  Headers(map[string]string{"Auth": "token"}).
  Retry(3, time.Second).
  Timeout(10*time.Second).
  JSON(&result)
```

**Slice Functional Patterns:**
```go
s := Slice[int]{1, 2, 3, 4, 5}
evens := s.Filter(func(x int) bool { return x%2 == 0 }).
  Map(func(x int) int { return x * 2 })
```

## Gotchas

- **Composite literals:** `(Slice[int]{}).IsEmpty()` requires parentheses—Go syntax rule
- **HTTP status codes:** Not errors; check `response.StatusCode` explicitly
- **Slice receivers:** Modifications return new slices; originals unchanged (immutable-spirit semantics)
- **File operations:** LoadFile and JSON read/write return errors; never panic for I/O

## Quick Reference

| Task | Command |
|------|---------|
| Run tests | `go test ./...` |
| Format | `gofmt -w ./...` |
| Check syntax | `go vet ./...` |
| Read docs | Open [pkg.go.dev](https://pkg.go.dev/github.com/renandotcorrea/goscript) or `godoc -h localhost:6060` |
