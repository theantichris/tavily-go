package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"github.com/theantichris/tavily-go"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	query := flag.String("query", "current weather in Knoxville, TN", "the search query to execute")
	timeout := flag.Duration("timeout", 30*time.Second, "the timeout for the search request")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		logger.Error("error loading .env file", "err", err)
		os.Exit(1)
	}

	apiKey := os.Getenv("TAVILY_API_KEY")
	searchURL := os.Getenv("TAVILY_SEARCH_URL")
	if apiKey == "" || searchURL == "" {
		logger.Error("TAVILY_API_KEY and TAVILY_SEARCH_URL environment variables must be set")
		os.Exit(1)
	}

	httpClient := &http.Client{Timeout: *timeout}
	provider, err := tavily.New(apiKey, httpClient, logger)
	if err != nil {
		logger.Error("error creating Tavily client", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	searchRequest := &tavily.SearchRequest{
		Query:                    *query,
		AutoParameters:           true,
		IncludeAnswer:            true,
		IncludeRawContent:        true,
		IncludeImages:            true,
		IncludeImageDescriptions: true,
		MaxResults:               20,
	}

	response, err := provider.Search(ctx, searchRequest)
	if err != nil {
		logger.Error("error performing search", "err", err)
		os.Exit(1)
	}

	spew.Dump(response)
}
