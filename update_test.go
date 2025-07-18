package reader

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUpdateDocument(t *testing.T) {
	publishedDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		documentID     string
		req            *UpdateDocumentRequest
		responseStatus int
		responseBody   interface{}
		wantErr        bool
		wantID         string
		wantURL        string
	}{
		{
			name:       "successful update with location change",
			documentID: "doc123",
			req: &UpdateDocumentRequest{
				Location: LocationNew,
			},
			responseStatus: http.StatusOK,
			responseBody: UpdateDocumentResponse{
				ID:  "doc123",
				URL: "https://read.readwise.io/read/doc123",
			},
			wantErr: false,
			wantID:  "doc123",
			wantURL: "https://read.readwise.io/read/doc123",
		},
		{
			name:       "successful update with multiple fields",
			documentID: "doc456",
			req: &UpdateDocumentRequest{
				Title:    "Updated Title",
				Author:   "Updated Author",
				Summary:  "Updated Summary",
				Location: LocationArchive,
				Category: CategoryArticle,
			},
			responseStatus: http.StatusOK,
			responseBody: UpdateDocumentResponse{
				ID:  "doc456",
				URL: "https://read.readwise.io/read/doc456",
			},
			wantErr: false,
			wantID:  "doc456",
			wantURL: "https://read.readwise.io/read/doc456",
		},
		{
			name:       "successful update with published date",
			documentID: "doc789",
			req: &UpdateDocumentRequest{
				Title:         "Article with Date",
				PublishedDate: &publishedDate,
			},
			responseStatus: http.StatusOK,
			responseBody: UpdateDocumentResponse{
				ID:  "doc789",
				URL: "https://read.readwise.io/read/doc789",
			},
			wantErr: false,
			wantID:  "doc789",
			wantURL: "https://read.readwise.io/read/doc789",
		},
		{
			name:       "empty document ID",
			documentID: "",
			req: &UpdateDocumentRequest{
				Title: "Test",
			},
			wantErr: true,
		},
		{
			name:       "nil request",
			documentID: "doc123",
			req:        nil,
			wantErr:    true,
		},
		{
			name:       "API error - not found",
			documentID: "nonexistent",
			req: &UpdateDocumentRequest{
				Title: "Test",
			},
			responseStatus: http.StatusNotFound,
			responseBody: map[string]interface{}{
				"detail": "Document not found",
			},
			wantErr: true,
		},
		{
			name:       "API error - unauthorized",
			documentID: "doc123",
			req: &UpdateDocumentRequest{
				Title: "Test",
			},
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
				if r.Method != "PATCH" {
					t.Errorf("expected PATCH method, got %s", r.Method)
				}

				// Verify path
				expectedPath := "/update/" + tt.documentID + "/"
				if r.URL.Path != expectedPath {
					t.Errorf("expected %s path, got %s", expectedPath, r.URL.Path)
				}

				// Verify authorization header
				authHeader := r.Header.Get("Authorization")
				if authHeader != "Token test-token" {
					t.Errorf("expected 'Token test-token' authorization, got %s", authHeader)
				}

				// Verify content type
				contentType := r.Header.Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("expected 'application/json' content type, got %s", contentType)
				}

				// Verify request body for valid requests
				if tt.req != nil && tt.documentID != "" {
					var reqBody UpdateDocumentRequest
					if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
						t.Errorf("failed to decode request body: %v", err)
					}
				}

				// Send response
				w.WriteHeader(tt.responseStatus)
				if tt.responseBody != nil {
					json.NewEncoder(w).Encode(tt.responseBody)
				}
			}))
			defer server.Close()

			// Create client with test server URL
			client := &client{
				baseURL:    server.URL,
				token:      "test-token",
				httpClient: &http.Client{},
			}

			// Call UpdateDocument
			ctx := context.Background()
			resp, err := client.UpdateDocument(ctx, tt.documentID, tt.req)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expect an error, we're done
			if tt.wantErr {
				return
			}

			// Verify response
			if resp.ID != tt.wantID {
				t.Errorf("expected ID %s, got %s", tt.wantID, resp.ID)
			}
			if resp.URL != tt.wantURL {
				t.Errorf("expected URL %s, got %s", tt.wantURL, resp.URL)
			}
		})
	}
}

func TestUpdateDocument_ContextCancellation(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(UpdateDocumentResponse{
			ID:  "doc123",
			URL: "https://read.readwise.io/read/doc123",
		})
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
	_, err := client.UpdateDocument(ctx, "doc123", &UpdateDocumentRequest{Title: "Test"})
	if err == nil {
		t.Error("expected error due to cancelled context, got nil")
	}
}
