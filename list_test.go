package reader

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListDocuments(t *testing.T) {
	tests := []struct {
		name            string
		opts            *ListDocumentsOptions
		responseStatus  int
		responseBody    interface{}
		wantErr         bool
		wantCount       int
		wantQueryParams map[string]string
	}{
		{
			name:           "successful list with no options",
			opts:           nil,
			responseStatus: http.StatusOK,
			responseBody: ListDocumentsResponse{
				Count:          2,
				NextPageCursor: stringPtr("next-cursor"),
				Results: []Document{
					{
						ID:       "doc1",
						URL:      "https://example.com/article1",
						Title:    "Test Article 1",
						Category: CategoryArticle,
						Location: LocationNew,
					},
					{
						ID:       "doc2",
						URL:      "https://example.com/article2",
						Title:    "Test Article 2",
						Category: CategoryArticle,
						Location: LocationLater,
					},
				},
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name: "list with filters",
			opts: &ListDocumentsOptions{
				Location:        LocationNew,
				Category:        CategoryArticle,
				Tag:             "programming",
				PageCursor:      "cursor123",
				WithHTMLContent: true,
			},
			responseStatus: http.StatusOK,
			responseBody: ListDocumentsResponse{
				Count: 1,
				Results: []Document{
					{
						ID:       "doc1",
						Title:    "Filtered Article",
						Category: CategoryArticle,
						Location: LocationNew,
					},
				},
			},
			wantErr:   false,
			wantCount: 1,
			wantQueryParams: map[string]string{
				"location":        "new",
				"category":        "article",
				"tag":             "programming",
				"pageCursor":      "cursor123",
				"withHtmlContent": "true",
			},
		},
		{
			name: "list with updated after filter",
			opts: &ListDocumentsOptions{
				UpdatedAfter: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			responseStatus: http.StatusOK,
			responseBody: ListDocumentsResponse{
				Count:   0,
				Results: []Document{},
			},
			wantErr:   false,
			wantCount: 0,
			wantQueryParams: map[string]string{
				"updatedAfter": "2024-01-01T00:00:00Z",
			},
		},
		{
			name:           "API error",
			opts:           nil,
			responseStatus: http.StatusInternalServerError,
			responseBody: map[string]interface{}{
				"error": "Internal server error",
			},
			wantErr: true,
		},
		{
			name:           "unauthorized error",
			opts:           nil,
			responseStatus: http.StatusUnauthorized,
			responseBody: map[string]interface{}{
				"detail": "Invalid token",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify method
				if r.Method != "GET" {
					t.Errorf("expected GET method, got %s", r.Method)
				}

				// Verify path
				if r.URL.Path != "/list/" {
					t.Errorf("expected /list/ path, got %s", r.URL.Path)
				}

				// Verify authorization header
				authHeader := r.Header.Get("Authorization")
				if authHeader != "Token test-token" {
					t.Errorf("expected 'Token test-token' authorization, got %s", authHeader)
				}

				// Verify query parameters
				if tt.wantQueryParams != nil {
					for key, expectedValue := range tt.wantQueryParams {
						actualValue := r.URL.Query().Get(key)
						if actualValue != expectedValue {
							t.Errorf("query param %s: expected %s, got %s", key, expectedValue, actualValue)
						}
					}
				}

				// Send response
				w.WriteHeader(tt.responseStatus)
				json.NewEncoder(w).Encode(tt.responseBody)
			}))
			defer server.Close()

			// Create client with test server URL
			client := &client{
				baseURL:    server.URL,
				token:      "test-token",
				httpClient: &http.Client{},
			}

			// Call ListDocuments
			ctx := context.Background()
			resp, err := client.ListDocuments(ctx, tt.opts)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("ListDocuments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expect an error, we're done
			if tt.wantErr {
				return
			}

			// Verify response
			if resp.Count != tt.wantCount {
				t.Errorf("expected count %d, got %d", tt.wantCount, resp.Count)
			}
		})
	}
}

func TestListDocuments_ContextCancellation(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListDocumentsResponse{})
	}))
	defer server.Close()

	client := &client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: &http.Client{},
	}

	// Create context that cancels immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Call should fail due to cancelled context
	_, err := client.ListDocuments(ctx, nil)
	if err == nil {
		t.Error("expected error due to cancelled context, got nil")
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
