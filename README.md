# goscript

Pure Go utility packages for writing scripts with minimal boilerplate. No external dependencies.

## Installation

```bash
go get github.com/renandotcorrea/goscript@latest
```

Requires Go 1.25.4 or higher.

## Quick Start

```go
package main

import (
    "fmt"
    "os"

    "github.com/renandotcorrea/goscript/env"
    "github.com/renandotcorrea/goscript/http"
    "github.com/renandotcorrea/goscript/slice"
    "github.com/renandotcorrea/goscript/try"
)

type Item struct {
    Name  string `json:"name"`
    Price int    `json:"price"`
}

func main() {
    defer try.Handle(func(err error) {
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    })

    token := env.MustGet("API_TOKEN")

    var items []Item
    try.Try(http.Get("https://api.example.com/items").
        Headers(map[string]string{"Authorization": "Bearer " + token}).
        Retry(3, time.Second).
        JSON(&items))

    slice.Slice[Item](items).
        Filter(func(i Item) bool { return i.Price < 100 }).
        ForEach(func(i Item) { fmt.Println(i.Name) })
}
```

## Packages

| Package | Purpose | Docs |
|---------|---------|------|
| `slice` | Generic `Slice[T]` with functional operations (Filter, Map, Reduce, …) | [pkg.go.dev/slice](https://pkg.go.dev/github.com/renandotcorrea/goscript/slice) |
| `http` | Fluent HTTP request builder with chainable API | [pkg.go.dev/http](https://pkg.go.dev/github.com/renandotcorrea/goscript/http) |
| `try` | Panic-based error handling to eliminate `if err != nil` boilerplate | [pkg.go.dev/try](https://pkg.go.dev/github.com/renandotcorrea/goscript/try) |
| `env` | Environment variable helpers with typed defaults | [pkg.go.dev/env](https://pkg.go.dev/github.com/renandotcorrea/goscript/env) |
| `file` | JSON read/write utilities | [pkg.go.dev/file](https://pkg.go.dev/github.com/renandotcorrea/goscript/file) |

### `slice`

[Full API reference →](https://pkg.go.dev/github.com/renandotcorrea/goscript/slice)

```go
s := slice.Slice[int]{1, 2, 3, 4, 5}

s.Filter(func(x int) bool { return x%2 == 0 }) // [2 4]
s.Map(func(x int) int { return x * 2 })         // [2 4 6 8 10]
s.Reduce(func(a, b int) int { return a + b }, 0) // 15
s.Reverse()                                       // [5 4 3 2 1]
s.Unique()                                        // deduplicates
s.Chunk(2)                                        // [[1 2] [3 4] [5]]
s.Contains(3)                                     // true
s.First() / s.Last()                              // *T, nil if empty
```

### `http`

[Full API reference →](https://pkg.go.dev/github.com/renandotcorrea/goscript/http)

```go
// Decode JSON response
var result MyStruct
err := http.Get("https://api.example.com/items").
    QueryParams(map[string]string{"page": "1"}).
    Headers(map[string]string{"Authorization": "Bearer " + token}).
    Retry(3, time.Second).
    Timeout(10 * time.Second).
    JSON(&result)

// POST with JSON body
err = http.Post("https://api.example.com/items").
    BodyJSON(map[string]any{"name": "foo"}).
    JSON(&result)

// Raw response
resp, err := http.Get(url).Do()
// resp.StatusCode, resp.Body ([]byte), resp.Headers
```

### `try`

[Full API reference →](https://pkg.go.dev/github.com/renandotcorrea/goscript/try)

```go
// Centralise error handling at the top of main
defer try.Handle(func(err error) {
    fmt.Fprintln(os.Stderr, "error:", err)
    os.Exit(1)
})

try.Try(os.Remove("file.txt"))                  // error-only calls
data := try.Try1(os.ReadFile("input.json"))     // (value, error) pairs
v1, v2 := try.Try2(someFunc())                  // (v1, v2, error) triples
```

### `env`

[Full API reference →](https://pkg.go.dev/github.com/renandotcorrea/goscript/env)

```go
token   := env.MustGet("API_TOKEN")         // panics if missing or empty
port    := env.GetOr("PORT", "8080")        // string with fallback
workers := env.GetIntOr("WORKERS", 4)       // int with fallback
debug   := env.GetBoolOr("DEBUG", false)    // bool with fallback
env.LoadFile(".env")                         // load a .env file
```

### `file`

[Full API reference →](https://pkg.go.dev/github.com/renandotcorrea/goscript/file)

```go
var cfg Config
try.Try(file.ReadJson("config.json", &cfg))
try.Try(file.WriteJson("output.json", result))
```

## Generating scripts with natural language (GitHub Copilot)

This repository includes a [Copilot prompt](.github/prompts/goscript.prompt.md) that generates ready-to-run Go scripts from a plain-language description, prioritising goscript packages.

**Requirements:** VS Code with the [GitHub Copilot](https://marketplace.visualstudio.com/items?itemName=GitHub.copilot) extension.

**Steps:**

1. Open the Chat panel (`Ctrl+Alt+I` / `Cmd+Alt+I`).
2. Click the **paperclip** icon → **Prompt...** → select **goscript**.
3. Describe what you need. For example:

   > Read a JSON file with a list of products, filter the ones with stock > 0, sort by price, and write the result to a new JSON file.

Copilot will produce a `package main` file that follows Go scripting best practices (DRY, KISS, YAGNI).

## Testing

```bash
go test ./...
```

## License
