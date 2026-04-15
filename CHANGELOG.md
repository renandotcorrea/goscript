# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.1.2] - 2026-04-15

### Added
- Complete godoc coverage for all packages with package-level comments.
- Canonical examples as `ExampleXxx()` test functions for:
  - `http.Get()` — fluent request chain with QueryParams and Headers
  - `env.GetIntOr()` — integer environment variable parsing with fallback
  - `env.LoadFile()` — loading environment variables from file
  - `file.ReadJson()` and `file.WriteJson()` — JSON file roundtrip
  - `try.Try1()` and `try.Try()` — panic-on-error patterns
- Documentation for previously undocumented functions:
  - `env.GetIntOr()`, `env.GetBoolOr()`, `env.GetDurationOr()`
  - `slice.IsEmpty()`, `slice.ForEach()`
- Type documentation for `HttpRequest`, `HttpResponse`, and `Slice[T]`.

### Changed
- `AGENTS.md` simplified and aligned with actual package structure and documentation practices.

## [v0.1.1] - 2026-04-13

### Added

- MIT License

## [v0.1.0] - 2026-04-13

### Added
- New `file` package with JSON helpers:
	- `ReadJson(filePath string, dest any) error`
	- `WriteJson(filePath string, src any) error`
- Unit tests for the new `file` package.

### Changed
- `env` package now includes `LoadFile(filePath ...string) error` to load env vars from a file (`.env` by default).
- Extended test coverage for `env` package (`LoadFile` scenarios).

## [0.0.2] - 2026-04-10

### Added

#### `env` package
- `MustGet(key string) string` — returns the value of an environment variable or panics if it is missing or empty.
- `GetOr(key, def string) string` — returns the value of an environment variable or a default value if it is missing or empty.

#### `http` package
- Fluent HTTP request builder with chainable API.
- `Get`, `Post`, `Put`, `Patch`, `Delete` constructors.
- `Headers(map[string]string)` — set request headers.
- `QueryParams(map[string]string)` — set URL query string parameters.
- `BodyJSON(src any)` — set a JSON-encoded body.
- `BodyXML(src any)` — set an XML-encoded body.
- `Body([]byte)` — set a raw body.
- `Retry(n int, backoff time.Duration)` — automatic retry with backoff for 429 and 5xx responses.
- `Timeout(d time.Duration)` — configurable per-request timeout.
- `JSON(dest any)` — execute request and unmarshal JSON response.
- `XML(dest any)` — execute request and unmarshal XML response.
- `Do()` — execute request and return raw `*HttpResponse`.

#### `slice` package
- Generic `Slice[T]` type backed by a plain Go slice.
- `Contains(value T) bool` — check if a value is present.
- `Filter(predicate func(T) bool) Slice[T]` — return elements that satisfy a predicate.
- `Map(fn func(T) T) Slice[T]` — transform each element.
- `FlatMap[T, U](s []T, fn func(T) []U) Slice[U]` — transform and flatten results.
- `Reduce(reducer func(T, T) T, initial T) T` — reduce to a single value.
- `ForEach(action func(T))` — iterate for side effects.
- `First() *T` / `Last() *T` — return a pointer to the first or last element, or nil if empty.
- `Reverse() Slice[T]` — return a new slice in reverse order.
- `Unique() Slice[T]` — return a new slice with duplicates removed.
- `Chunk(n int) []Slice[T]` — split the slice into chunks of size n.
- `IsEmpty() bool` — check if the slice has no elements.
- `Append(values ...T) Slice[T]` — append elements and return a new slice.
- `Len() int` / `Cap() int` — return length and capacity.
- `ToMap[K, V](s []V, keySelector func(V) K) map[K]V` — convert a slice to a map.

#### `try` package
- `Try(err error)` — panic if the error is not nil.
- `Try1[T](value T, err error) T` — return value or panic on error.
- `Try2[T1, T2](v1 T1, v2 T2, err error) (T1, T2)` — return two values or panic on error.

### Changed
- Module renamed to `github.com/renandotcorrea/goscript`.

[0.0.2]: https://github.com/renandotcorrea/goscript/releases/tag/0.0.2
[0.1.0]: https://github.com/renandotcorrea/goscript/releases/tag/v0.1.0
[0.1.1]: https://github.com/renandotcorrea/goscript/releases/tag/v0.1.1
[0.1.2]: https://github.com/renandotcorrea/goscript/releases/tag/v0.1.2
