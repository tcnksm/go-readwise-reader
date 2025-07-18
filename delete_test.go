package reader

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_DeleteDocument(t *testing.T) {
	tests := []struct {
		name           string
		documentID     string
		serverResponse string
		serverStatus   int
		wantErr        bool
		errType        string
	}{
		{
			name:         "successful_deletion",
			documentID:   "doc123",
			serverStatus: http.StatusNoContent,
			wantErr:      false,
		},
		{
			name:       "empty_document_id",
			documentID: "",
			wantErr:    true,
			errType:    "invalid_parameter",
		},
		{
			name:           "document_not_found",
			documentID:     "nonexistent",
			serverResponse: `{"error": "document not found"}`,
			serverStatus:   http.StatusNotFound,
			wantErr:        true,
		},
		{
			name:           "unauthorized",
			documentID:     "doc456",
			serverResponse: `{"error": "unauthorized"}`,
			serverStatus:   http.StatusUnauthorized,
			wantErr:        true,
		},
		{
			name:           "server_error",
			documentID:     "doc789",
			serverResponse: `{"error": "internal server error"}`,
			serverStatus:   http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:           "bad_request",
			documentID:     "invalid-id",
			serverResponse: `{"error": "invalid document ID"}`,
			serverStatus:   http.StatusBadRequest,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request method
				if r.Method != http.MethodDelete {
					t.Errorf("Expected DELETE request, got %s", r.Method)
				}

				// Verify URL path
				expectedPath := "/delete/" + tt.documentID + "/"
				if tt.documentID != "" && r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				// Verify headers
				if auth := r.Header.Get("Authorization"); !strings.HasPrefix(auth, "Token ") {
					t.Errorf("Expected Authorization header with Token prefix, got %s", auth)
				}

				if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
					t.Errorf("Expected Content-Type application/json, got %s", contentType)
				}

				// Verify no request body
				if r.ContentLength > 0 {
					t.Error("Expected no request body for DELETE request")
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
			err := client.DeleteDocument(ctx, tt.documentID)

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
		})
	}
}
