# goscript

A collection of useful pure Go utility packages for working with slices, HTTP requests, and error handling.

This package allow you to build scripts with Go using few lines of code.

## Packages

### `slice` - Generic Slice Utilities

Provides a generic `Slice[T]` type with a comprehensive set of functional methods for working with slices.

**Key methods:**
- `Contains(value T)` - Check if slice contains a value
- `Filter(predicate func(T) bool)` - Filter elements by predicate
- `Map(fn func(T) U)` - Transform elements
- `FlatMap(slice Slice[T], fn func(T) []U)` - Transform and flatten results
- `Reduce(fn func(acc, x T) T, initial T)` - Reduce to a single value
- `ForEach(fn func(T))` - Iterate over elements
- `First()` / `Last()` - Get first or last element
- `Reverse()` - Reverse slice order
- `Unique()` - Get unique elements
- `Chunk(n int)` - Split slice into smaller chunks
- `IsEmpty()` - Check if slice is empty
- `Append(elements...T)` - Append elements
- `Len()` / `Cap()` - Get length and capacity
- `ToMap(slice Slice[T], keyFn func(T) K)` - Convert slice to map

**Example:**
```go
import "github.com/renandotcorrea/goscript/slice"

s := slice.Slice[int]{1, 2, 3, 4, 5}
evens := s.Filter(func(x int) bool { return x%2 == 0 })
doubled := evens.Map(func(x int) int { return x * 2 })
```

### http - Fluent HTTP Client

A fluent HTTP request builder for making HTTP requests with a clean, chainable API.

**Methods:**
- `Get(url)` - Create a GET request
- `Post(url)` - Create a POST request
- `Put(url)` - Create a PUT request
- `QueryParams(map)` - Set URL query parameters
- `Retry(n, backoff)` - Retry transient failures (429/5xx and transient network errors)
- `Timeout(d)` - Set per-request timeout
- `Headers(map)` - Set request headers
- `BodyJSON(data)` - Set JSON body
- `BodyXML(data)` - Set XML body
- `BodyRaw(data)` - Set raw body
- `JSON(dest)` - Execute and unmarshal JSON response
- `XML(dest)` - Execute and unmarshal XML response
- `Raw()` - Execute and get raw response

**Example:**
```go
import "github.com/renandotcorrea/goscript/http"

var result map[string]interface{}
err := http.Get("http://api.example.com/data").
    QueryParams(map[string]string{"page": "1", "page_size": "100"}).
    Retry(3, time.Second).
    Timeout(10 * time.Second).
    Headers(map[string]string{"Authorization": "Bearer token"}).
    JSON(&result)
```

### env - Environment Variable Helpers

Helpers to make environment variable reads simpler in scripts and CLIs.

**Functions:**
- `MustGet(key)` - Return value or panic if missing/empty
- `GetOr(key, def)` - Return value or fallback default

**Example:**
```go
import "github.com/renandotcorrea/goscript/env"

token := env.MustGet("API_TOKEN")
region := env.GetOr("REGION", "us-east-1")
_ = token
_ = region
```

### try - Error Handling Utilities

Utility functions for error handling that panic when errors occur. Useful for simplifying code when error handling isn't critical.

**Functions:**
- `Try(err)` - Panic if error is not nil
- `Try1[T](value T, err)` - Return value or panic on error
- `Try2[T1, T2](v1 T1, v2 T2, err)` - Return two values or panic on error

**Example:**
```go
import "github.com/renandotcorrea/goscript/try"

data := try.Try1(ioutil.ReadFile("file.txt"))
```

## Installation

```bash
go get github.com/renandotcorrea/goscript
```

## Testing

Run tests for all packages:

```bash
go test ./...
```

Run tests for a specific package:

```bash
go test ./slice
go test ./http
go test ./try
```

## Go Version

Requires Go 1.25.4 or higher

## License