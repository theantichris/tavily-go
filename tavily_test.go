package tavily

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	t.Run("handles successful search", func(t *testing.T) {
		t.Parallel()

		searchRequest := SearchRequest{Query: "test query"}

		want := SearchResponse{
			Query:  "test query",
			Answer: "test answer",
			Results: []SiteResult{
				{Title: "test title", URL: "http://example.com"},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(want)
		}))
		defer server.Close()

		client := New("fake-api-key", server.URL, server.Client(), slog.Default())
		ctx := context.Background()

		got, err := client.Search(ctx, &searchRequest)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got.Query != want.Query {
			t.Errorf("got Query %q, want %q", got.Query, want.Query)
		}

		if got.Answer != want.Answer {
			t.Errorf("got Answer %q, want %q", got.Answer, want.Answer)
		}

		if len(got.Results) != len(want.Results) {
			t.Fatalf("got %d Results, want %d", len(got.Results), len(want.Results))
		}

		if got.Results[0].Title != want.Results[0].Title {
			t.Errorf("got Results[0].Title %q, want %q", got.Results[0].Title, want.Results[0].Title)
		}
	})

	t.Run("handles non-200 status code", func(t *testing.T) {
		t.Parallel()

		searchRequest := SearchRequest{Query: "test query"}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "fail", http.StatusInternalServerError)
		}))
		defer server.Close()

		client := New("fake-api-key", server.URL, server.Client(), slog.Default())
		ctx := context.Background()

		_, err := client.Search(ctx, &searchRequest)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "non-2xx response from Tavily API") {
			t.Errorf("got error %q, want it to contain %q", err.Error(), "non-2xx response from Tavily API")
		}
	})

	t.Run("handles invalid JSON response", func(t *testing.T) {
		t.Parallel()

		searchRequest := SearchRequest{Query: "test query"}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("{invalid json"))
		}))
		defer server.Close()

		client := New("fake-api-key", server.URL, server.Client(), slog.Default())
		ctx := context.Background()

		_, err := client.Search(ctx, &searchRequest)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "failed to decode response JSON") {
			t.Errorf("got error %q, want it to contain %q", err.Error(), "failed to decode response JSON")
		}
	})
}
