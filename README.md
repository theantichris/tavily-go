# tavily-go

[![Go Reference](https://pkg.go.dev/badge/github.com/theantichris/tavily-go.svg)](https://pkg.go.dev/github.com/theantichris/tavily-go) [![CI](https://github.com/theantichris/tavily-go/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/theantichris/tavily-go/actions/workflows/ci.yml)

A Go client library for the Tavily Web Search API.

## Requirements

- Go 1.25+
- A Tavily API key (set as `TAVILY_API_KEY` in your environment or pass directly to the client constructor)

## Installation

```sh
go get github.com/theantichris/tavily-go@latest
```

## Quickstart

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "time"

    tavily "github.com/theantichris/tavily-go"
)

func main() {
    apiKey := os.Getenv("TAVILY_API_KEY")
    if apiKey == "" {
        panic("TAVILY_API_KEY must be set")
    }

    httpClient := &http.Client{Timeout: 10 * time.Second}

    client, err := tavily.New(apiKey, httpClient, nil) // pass a *slog.Logger instead of nil to enable structured logging
    if err != nil {
        panic(err)
    }

    req := &tavily.SearchRequest{
        Query:          "golang news",
        AutoParameters: true,
        IncludeAnswer:  true,
        MaxResults:     5,
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    resp, err := client.Search(ctx, req)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Query: %s\nAnswer: %s\nResults: %d\n", resp.Query, resp.Answer, len(resp.Results))
}
```

## Configuration

- API key: The client requires a Tavily API key. You can source it from your environment:
  - PowerShell: `$env:TAVILY_API_KEY = "<your-key>"`
  - Bash/zsh: `export TAVILY_API_KEY="<your-key>"`
- Logging: Pass a `*slog.Logger` to `New(...)` to enable structured logs, or `nil` to disable logging.
- HTTP client: Pass a custom `*http.Client` (timeouts, transport, proxy, etc.). If `nil` is provided, the default client is used.

## Example CLI

An example program is provided under `cmd`:

```powershell
# PowerShell
Copy-Item .env.example .env
# Edit .env to set TAVILY_API_KEY

go run ./cmd -query 'current weather in Knoxville, TN' -timeout 15s
```

```bash
# Bash
cp .env.example .env
# Edit .env to set TAVILY_API_KEY

go run ./cmd -query 'current weather in Knoxville, TN' -timeout 15s
```

## Testing

```sh
# Run all tests
go test ./...

# Verbose
go test -v ./...

# Run a single test (and its subtests)
go test -run '^TestSearch$' -v .

# Run a specific subtest (example)
go test -run '^TestSearch/handles successful search$' -v .

# Coverage summary
go test -cover ./...

# Coverage profile and reports
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
# HTML report
go tool cover -html=coverage.out
```

## Mocking in your tests

Use the provided mock to stub the `WebSearchClient` interface in your own code:

```go
mock := &tavily.MockWebSearchClient{
    SearchFunc: func(ctx context.Context, r *tavily.SearchRequest) (tavily.SearchResponse, error) {
        return tavily.SearchResponse{
            Query:  r.Query,
            Answer: "mock-answer",
        }, nil
    },
}

// Example usage in your unit tests
resp, err := mock.Search(context.Background(), &tavily.SearchRequest{Query: "test"})
if err != nil {
    t.Fatalf("unexpected error: %v", err)
}
```

## Error semantics

- Non-2xx responses from the Tavily API return an error that includes the HTTP status code and response body.
- JSON decode failures also return a descriptive error.

## API reference

- <https://pkg.go.dev/github.com/theantichris/tavily-go>

## License

This project is licensed under the terms of the license found in the [LICENSE](./LICENSE) file.
