package reader

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClient_CreateDocument(t *testing.T) {
	tests := []struct {
		name           string
		req            *CreateDocumentRequest
		serverResponse string
		serverStatus   int
		want           *CreateDocumentResponse
		wantErr        bool
		errType        string
	}{
		{
			name: "successful_creation",
			req: &CreateDocumentRequest{
				URL:   "https://example.com/article",
				Title: "Test Article",
				Tags:  []string{"test", "article"},
			},
			serverResponse: `{"id": "doc123", "url": "https://example.com/article"}`,
			serverStatus:   http.StatusCreated,
			want: &CreateDocumentResponse{
				ID:  "doc123",
				URL: "https://example.com/article",
			},
			wantErr: false,
		},
		{
			name: "existing_document",
			req: &CreateDocumentRequest{
				URL: "https://example.com/existing",
			},
			serverResponse: `{"id": "doc456", "url": "https://example.com/existing"}`,
			serverStatus:   http.StatusOK,
			want: &CreateDocumentResponse{
				ID:  "doc456",
				URL: "https://example.com/existing",
			},
			wantErr: false,
		},
		{
			name: "complete_request",
			req: &CreateDocumentRequest{
				URL:           "https://example.com/complete",
				HTML:          "<html><body>Content</body></html>",
				Title:         "Complete Article",
				Author:        "John Doe",
				Summary:       "A complete test article",
				PublishedDate: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				Tags:          []string{"complete", "test"},
				Location:      LocationNew,
				Category:      CategoryArticle,
				ImageURL:      "https://example.com/image.jpg",
			},
			serverResponse: `{"id": "doc789", "url": "https://example.com/complete"}`,
			serverStatus:   http.StatusCreated,
			want: &CreateDocumentResponse{
				ID:  "doc789",
				URL: "https://example.com/complete",
			},
			wantErr: false,
		},
		{
			name:    "nil_request",
			req:     nil,
			wantErr: true,
			errType: "invalid_request",
		},
		{
			name: "empty_url",
			req: &CreateDocumentRequest{
				Title: "No URL",
			},
			wantErr: true,
			errType: "invalid_request",
		},
		{
			name: "server_error",
			req: &CreateDocumentRequest{
				URL: "https://example.com/error",
			},
			serverResponse: `{"error": "internal server error"}`,
			serverStatus:   http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name: "bad_request",
			req: &CreateDocumentRequest{
				URL: "invalid-url",
			},
			serverResponse: `{"error": "invalid URL"}`,
			serverStatus:   http.StatusBadRequest,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request method
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", r.Method)
				}

				// Verify URL path
				if r.URL.Path != "/save/" {
					t.Errorf("Expected path /save/, got %s", r.URL.Path)
				}

				// Verify headers
				if auth := r.Header.Get("Authorization"); !strings.HasPrefix(auth, "Token ") {
					t.Errorf("Expected Authorization header with Token prefix, got %s", auth)
				}

				if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
					t.Errorf("Expected Content-Type application/json, got %s", contentType)
				}

				// Verify request body if not nil request
				if tt.req != nil {
					var reqBody CreateDocumentRequest
					if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
						t.Errorf("Failed to decode request body: %v", err)
					}

					if reqBody.URL != tt.req.URL {
						t.Errorf("Expected URL %s, got %s", tt.req.URL, reqBody.URL)
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				if tt.serverResponse != "" {
					w.Write([]byte(tt.serverResponse))
				}
			}))
			defer server.Close()

			// Create client
			client := &client{
				baseURL:    server.URL,
				token:      "test-token",
				httpClient: &http.Client{},
			}

			// Execute test
			ctx := context.Background()
			got, err := client.CreateDocument(ctx, tt.req)

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got none")
					return
				}

				if tt.errType != "" {
					var clientErr *ClientError
					if errors.As(err, &clientErr) {
						if clientErr.Type != tt.errType {
							t.Errorf("Expected error type %s, got %s", tt.errType, clientErr.Type)
						}
					}
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check response
			if got == nil {
				t.Error("Expected response, got nil")
				return
			}

			if got.ID != tt.want.ID {
				t.Errorf("Expected ID %s, got %s", tt.want.ID, got.ID)
			}

			if got.URL != tt.want.URL {
				t.Errorf("Expected URL %s, got %s", tt.want.URL, got.URL)
			}
		})
	}
}
