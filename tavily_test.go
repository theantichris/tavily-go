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

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("creates client with valid parameters", func(t *testing.T) {
		t.Parallel()

		client, err := New("valid-api-key", http.DefaultClient, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if client == nil {
			t.Fatal("expected client to be non-nil")
		}

		if _, ok := client.(*tavilyClient); !ok {
			t.Fatalf("got type %T, want *tavilyClient", client)
		}

		if client.(*tavilyClient).apiKey != "valid-api-key" {
			t.Errorf("got apiKey %q, want %q", client.(*tavilyClient).apiKey, "valid-api-key")
		}

		if client.(*tavilyClient).searchUrl != searchURL {
			t.Errorf("got searchUrl %q, want %q", client.(*tavilyClient).searchUrl, searchURL)
		}
	})

	t.Run("returns error with empty apiKey", func(t *testing.T) {
		t.Parallel()

		_, err := New("", http.DefaultClient, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "apiKey cannot be empty") {
			t.Errorf("got error %q, want it to contain %q", err.Error(), "apiKey cannot be empty")
		}
	})

	t.Run("uses default http.Client when nil is provided", func(t *testing.T) {
		t.Parallel()

		client, err := New("valid-api-key", nil, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if client == nil {
			t.Fatal("expected client to be non-nil")
		}

		if client.(*tavilyClient).httpClient != http.DefaultClient {
			t.Errorf("got httpClient %v, want %v", client.(*tavilyClient).httpClient, http.DefaultClient)
		}
	})
}

func TestSearch(t *testing.T) {
	t.Parallel()

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

		client, _ := New("fake-api-key", server.Client(), nil)
		if c, ok := client.(*tavilyClient); ok {
			c.searchUrl = server.URL
		}

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

		client, _ := New("fake-api-key", server.Client(), nil)
		if c, ok := client.(*tavilyClient); ok {
			c.searchUrl = server.URL
		}

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

		client, _ := New("fake-api-key", server.Client(), nil)
		if c, ok := client.(*tavilyClient); ok {
			c.searchUrl = server.URL
		}

		ctx := context.Background()

		_, err := client.Search(ctx, &searchRequest)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "failed to decode response JSON") {
			t.Errorf("got error %q, want it to contain %q", err.Error(), "failed to decode response JSON")
		}
	})

	t.Run("handles HTTP request failure", func(t *testing.T) {
		t.Parallel()

		searchRequest := SearchRequest{Query: "test query"}

		client, _ := New("fake-api-key", &http.Client{}, nil)
		if c, ok := client.(*tavilyClient); ok {
			c.searchUrl = "http://invalid-url"
		}

		ctx := context.Background()

		_, err := client.Search(ctx, &searchRequest)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "HTTP request to Tavily API failed") {
			t.Errorf("got error %q, want it to contain %q", err.Error(), "HTTP request to Tavily API failed")
		}
	})

	t.Run("handles HTTP request creation failure", func(t *testing.T) {
		t.Parallel()

		client, _ := New("valid-api-key", http.DefaultClient, nil)
		// Set an invalid URL to force http.NewRequestWithContext to fail
		if c, ok := client.(*tavilyClient); ok {
			c.searchUrl = ":"
		}

		searchRequest := &SearchRequest{Query: "test query"}
		_, err := client.Search(context.Background(), searchRequest)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "failed to create Tavily HTTP request") {
			t.Errorf("got error %q, want it to contain %q", err.Error(), "failed to create Tavily HTTP request")
		}
	})
}

func TestLogError(t *testing.T) {
	t.Parallel()

	client, err := New("valid-api-key", http.DefaultClient, slog.Default())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if c, ok := client.(*tavilyClient); ok {
		c.logError("test error message", "key", "value")
	}
}

func TestLogInfo(t *testing.T) {
	t.Parallel()

	client, err := New("valid-api-key", http.DefaultClient, slog.Default())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if c, ok := client.(*tavilyClient); ok {
		c.logInfo("test info message", "key", "value")
	}
}
