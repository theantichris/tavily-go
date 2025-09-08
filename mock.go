package tavily

import "context"

// MockWebSearchClient is a mock implementation of the WebSearchClient interface for testing.
type MockWebSearchClient struct {
	SearchFunc func(ctx context.Context, searchRequest *SearchRequest) (SearchResponse, error)
}

// Search calls the mock Search function if defined, otherwise returns zero values.
func (client *MockWebSearchClient) Search(ctx context.Context, searchRequest *SearchRequest) (SearchResponse, error) {
	if client.SearchFunc != nil {
		return client.SearchFunc(ctx, searchRequest)
	}

	return SearchResponse{}, nil
}
