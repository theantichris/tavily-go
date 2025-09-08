package tavily

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// WebSearchProvider defines the interface for a web search provider.
type WebSearchProvider interface {
	Search(ctx context.Context, query string) (SearchResponse, error)
}

// tavilyProvider implements a web search provider using the Tavily API.
type tavilyProvider struct {
	apiKey     string
	searchUrl  string
	httpClient *http.Client
	logger     *slog.Logger
}

// New creates a new instance of TavilyProvider.
func New(apiKey string, searchURL string, httpClient *http.Client, logger *slog.Logger) WebSearchProvider {
	return &tavilyProvider{
		apiKey:     apiKey,
		searchUrl:  searchURL,
		httpClient: httpClient,
		logger:     logger,
	}
}

// Search performs a web search using the Tavily API with the specified query, maximum results, and time range in days.
func (provider *tavilyProvider) Search(ctx context.Context, query string) (SearchResponse, error) {
	tavilyRequest := searchRequest{
		Query:                    query,
		AutoParameters:           true,
		IncludeAnswer:            true,
		IncludeRawContent:        true,
		IncludeImages:            true,
		IncludeImageDescriptions: true,
		MaxResults:               20,
	}

	requestBody, err := json.Marshal(tavilyRequest)
	if err != nil {
		provider.logger.Error("failed to marshal Tavily request to JSON", "err", err)

		return SearchResponse{}, fmt.Errorf("failed to marshal Tavily request to JSON: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, provider.searchUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		provider.logger.Error("failed to create Tavily HTTP request", "err", err)

		return SearchResponse{}, fmt.Errorf("failed to create Tavily HTTP request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+provider.apiKey)

	provider.logger.Info("sending request to Tavily API", "url", provider.searchUrl, "query", query)

	response, err := provider.httpClient.Do(request)
	if err != nil {
		provider.logger.Error("HTTP request to Tavily API failed", "err", err)

		return SearchResponse{}, fmt.Errorf("HTTP request to Tavily API failed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode/100 != 2 {
		bodyBytes, _ := io.ReadAll(response.Body)

		provider.logger.Error("non-2xx response from Tavily API", "status", response.StatusCode, "body", string(bodyBytes))

		return SearchResponse{}, fmt.Errorf("non-2xx response from Tavily API: %d - %s", response.StatusCode, string(bodyBytes))
	}

	var searchResponse SearchResponse
	if err := json.NewDecoder(response.Body).Decode(&searchResponse); err != nil {
		provider.logger.Error("failed to decode response JSON", "err", err)

		return SearchResponse{}, fmt.Errorf("failed to decode response JSON: %w", err)
	}

	return searchResponse, nil
}
