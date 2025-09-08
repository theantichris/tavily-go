package tavily

import "context"

// MockSearchProvider is a mock implementation of the WebSearchProvider interface for testing.
type MockSearchProvider struct {
	Response SearchResponse
	Err      error
}

// Search returns the predefined response or error.
func (provider *MockSearchProvider) Search(ctx context.Context, query string) (SearchResponse, error) {
	return provider.Response, provider.Err
}
