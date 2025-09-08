package tavily

import (
	"context"
	"testing"
)

func TestMockSearch(t *testing.T) {
	t.Run("calls mock Search function", func(t *testing.T) {
		t.Parallel()

		mockResponse := SearchResponse{Query: "mock query", Answer: "mock answer"}

		mockClient := &MockWebSearchClient{
			SearchFunc: func(ctx context.Context, searchRequest *SearchRequest) (SearchResponse, error) {
				return mockResponse, nil
			},
		}

		ctx := context.Background()
		searchRequest := &SearchRequest{Query: "test"}
		got, err := mockClient.Search(ctx, searchRequest)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got.Query != mockResponse.Query {
			t.Errorf("got Query %q, want %q", got.Query, mockResponse.Query)
		}

		if got.Answer != mockResponse.Answer {
			t.Errorf("got Answer %q, want %q", got.Answer, mockResponse.Answer)
		}
	})

	t.Run("returns zero values when SearchFunc is nil", func(t *testing.T) {
		t.Parallel()

		mockClient := &MockWebSearchClient{}

		ctx := context.Background()
		searchRequest := &SearchRequest{Query: "test"}
		got, err := mockClient.Search(ctx, searchRequest)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got.Query != "" {
			t.Errorf("got Query %q, want empty", got.Query)
		}

		if got.Answer != "" {
			t.Errorf("got Answer %q, want empty", got.Answer)
		}
	})
}
