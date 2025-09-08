# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

Project overview

This repository is a small Go client library for the Tavily Web Search API, with an example CLI under cmd that demonstrates usage via environment configuration and flags.

Toolchain: Go 1.25 (per go.mod)
Common commands

Build (library)

```sh
# Build all packages
go build ./...
```

Lint (standard tooling)

```sh
# Static analysis
go vet ./...
```

Tests

```sh
# Run all tests across the module
go test ./...

# Verbose output
go test -v ./...

# Run a single top-level test by name (PowerShell quoting shown)
# Example: only run TestSearch (and its subtests)
go test -run '^TestSearch$' -v .

# Run a specific subtest (as defined in tavily_test.go)
go test -run '^TestSearch/handles successful search$' -v .

# Code coverage (summary)
go test -cover ./...

# Coverage profile and reports
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
# Opens HTML report in browser
go tool cover -html=coverage.out
```

Example CLI (cmd)

The example program in cmd reads TAVILY_API_KEY from a local .env (via godotenv) or your environment, and accepts -query and -timeout flags.

```sh
# Prepare local env file (PowerShell)
Copy-Item .env.example .env
# Edit .env to set TAVILY_API_KEY, then run the CLI

# Run with defaults
go run ./cmd

# Provide custom query and timeout
go run ./cmd -query 'current weather in Knoxville, TN' -timeout 15s
```

High-level architecture

- Core package (tavily)
  - Interface: WebSearchClient defines Search(ctx, \*SearchRequest) (SearchResponse, error) for dependency inversion and testability.
  - Implementation: tavilyClient is constructed via `New(apiKey string, httpClient *http.Client, logger *slog.Logger) (WebSearchClient, error)`.
    - Validates apiKey, defaults httpClient when nil.
    - Optional structured logging via slog; pass nil to disable logs.
  - HTTP behavior: Search marshals SearchRequest to JSON and POSTs to <https://api.tavily.com/search> with Authorization: Bearer apiKey.
    - Non-2xx responses include status code and body in the returned error.
    - Successful responses decode into SearchResponse.
- Data contracts (types.go)
  - SearchRequest maps directly to Tavily API parameters (auto parameters, answer inclusion, images, etc.).
  - SearchResponse includes the answer, results, images, auto_parameters, and request_id.
- Testing strategy
  - Unit tests use httptest.NewServer and replace the internal base URL (via type assertion to \*tavilyClient) to simulate API responses.
  - MockWebSearchClient provides a simple stub for consumer tests, enabling injection of custom Search behavior.
- Example application (cmd)
  - Loads .env via github.com/joho/godotenv, parses flags, builds an http.Client with timeout, constructs the library client, performs a search, and prints the response with github.com/davecgh/go-spew/spew.

Notes for contributors

- Logging is opt-in: pass a \*slog.Logger to New to enable structured logs; otherwise, logging is suppressed.
- To add additional Tavily endpoints, follow the existing pattern: define request/response types and implement a method on tavilyClient that constructs the HTTP request and decodes the response.
