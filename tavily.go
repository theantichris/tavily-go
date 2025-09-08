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

// WebSearchClient defines the interface for a web search client.
type WebSearchClient interface {
	Search(ctx context.Context, searchRequest *SearchRequest) (SearchResponse, error)
}

// tavilyClient implements a web search client for the Tavily API.
type tavilyClient struct {
	apiKey     string
	searchUrl  string
	httpClient *http.Client
	logger     *slog.Logger
}

// New creates a new instance of tavilyClient.
func New(apiKey string, searchURL string, httpClient *http.Client, logger *slog.Logger) WebSearchClient {
	// TODO: Validate apiKey and searchURL are not empty, httpClient and logger are not nil.
	// TODO: Pass nil for no logging
	return &tavilyClient{
		apiKey:     apiKey,
		searchUrl:  searchURL,
		httpClient: httpClient,
		logger:     logger,
	}
}

// Search performs a web search using the Tavily API with the specified query, maximum results, and time range in days.
func (tavilyClient *tavilyClient) Search(ctx context.Context, searchRequest *SearchRequest) (SearchResponse, error) {
	requestBody, err := json.Marshal(searchRequest)
	if err != nil {
		tavilyClient.logger.Error("failed to marshal Tavily request to JSON", "err", err)

		return SearchResponse{}, fmt.Errorf("failed to marshal Tavily request to JSON: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, tavilyClient.searchUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		tavilyClient.logger.Error("failed to create Tavily HTTP request", "err", err)

		return SearchResponse{}, fmt.Errorf("failed to create Tavily HTTP request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+tavilyClient.apiKey)

	tavilyClient.logger.Info("sending request to Tavily API", "url", tavilyClient.searchUrl, "query", searchRequest.Query)

	response, err := tavilyClient.httpClient.Do(request)
	if err != nil {
		tavilyClient.logger.Error("HTTP request to Tavily API failed", "err", err)

		return SearchResponse{}, fmt.Errorf("HTTP request to Tavily API failed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode/100 != 2 {
		bodyBytes, _ := io.ReadAll(response.Body)

		tavilyClient.logger.Error("non-2xx response from Tavily API", "status", response.StatusCode, "body", string(bodyBytes))

		return SearchResponse{}, fmt.Errorf("non-2xx response from Tavily API: %d - %s", response.StatusCode, string(bodyBytes))
	}

	var searchResponse SearchResponse
	if err := json.NewDecoder(response.Body).Decode(&searchResponse); err != nil {
		tavilyClient.logger.Error("failed to decode response JSON", "err", err)

		return SearchResponse{}, fmt.Errorf("failed to decode response JSON: %w", err)
	}

	return searchResponse, nil
}
