package reader

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_CreateDocument(t *testing.T) {
	tests := []struct {
		name           string
		request        *CreateDocumentRequest
		serverResponse string
		statusCode     int
		wantErr        bool
		wantErrType    string
	}{
		{
			name: "successful creation with all fields",
			request: &CreateDocumentRequest{
				URL:      "https://example.com/article",
				Title:    stringPtr("Example Article"),
				Category: stringPtr("article"),
			},
			serverResponse: `{
				"id": "doc123",
				"url": "https://example.com/article",
				"title": "Example Article",
				"category": "article",
				"created_at": "2023-01-01T00:00:00Z",
				"updated_at": "2023-01-01T00:00:00Z"
			}`,
			statusCode: 201,
			wantErr:    false,
		},
		{
			name: "successful creation with URL only",
			request: &CreateDocumentRequest{
				URL: "https://example.com/minimal",
			},
			serverResponse: `{
				"id": "doc456",
				"url": "https://example.com/minimal",
				"title": "Minimal Article",
				"category": "",
				"created_at": "2023-01-01T00:00:00Z",
				"updated_at": "2023-01-01T00:00:00Z"
			}`,
			statusCode: 201,
			wantErr:    false,
		},
		{
			name:        "nil request",
			request:     nil,
			wantErr:     true,
			wantErrType: "*reader.ClientError",
		},
		{
			name: "empty URL",
			request: &CreateDocumentRequest{
				URL: "",
			},
			wantErr:     true,
			wantErrType: "*reader.ClientError",
		},
		{
			name: "API error response",
			request: &CreateDocumentRequest{
				URL: "https://example.com/invalid",
			},
			serverResponse: `{
				"error": "Invalid URL format"
			}`,
			statusCode:  400,
			wantErr:     true,
			wantErrType: "*reader.APIError",
		},
		{
			name: "server error",
			request: &CreateDocumentRequest{
				URL: "https://example.com/article",
			},
			serverResponse: `Internal Server Error`,
			statusCode:     500,
			wantErr:        true,
			wantErrType:    "*reader.APIError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify the request method and headers
				if r.Method != "POST" {
					t.Errorf("Expected POST request, got %s", r.Method)
				}
				
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
				}
				
				authHeader := r.Header.Get("Authorization")
				if !strings.HasPrefix(authHeader, "Token ") {
					t.Errorf("Expected Authorization header with Token prefix, got %s", authHeader)
				}
				
				// Verify the request body for valid requests
				if tt.request != nil && tt.request.URL != "" {
					var body CreateDocumentRequest
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Errorf("Failed to decode request body: %v", err)
					}
					
					if body.URL != tt.request.URL {
						t.Errorf("Expected URL %s, got %s", tt.request.URL, body.URL)
					}
				}
				
				// Send the test response
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// Create a client with the test server URL
			client := &client{
				baseURL:    server.URL,
				token:      "test-token",
				httpClient: &http.Client{},
			}

			// Execute the method
			ctx := context.Background()
			result, err := client.CreateDocument(ctx, tt.request)

			// Check error expectations
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.wantErrType != "" {
					errType := getErrorType(err)
					if errType != tt.wantErrType {
						t.Errorf("Expected error type %s, got %s", tt.wantErrType, errType)
					}
				}
				return
			}

			// Check success case
			if result == nil {
				t.Error("Expected non-nil result for successful case")
				return
			}

			// Parse expected response for comparison
			var expected CreateDocumentResponse
			if err := json.Unmarshal([]byte(tt.serverResponse), &expected); err != nil {
				t.Fatalf("Failed to parse expected response: %v", err)
			}

			if result.ID != expected.ID {
				t.Errorf("Expected ID %s, got %s", expected.ID, result.ID)
			}
			if result.URL != expected.URL {
				t.Errorf("Expected URL %s, got %s", expected.URL, result.URL)
			}
		})
	}
}

func TestCreateDocumentRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request *CreateDocumentRequest
		wantErr bool
	}{
		{
			name: "valid request with all fields",
			request: &CreateDocumentRequest{
				URL:      "https://example.com",
				Title:    stringPtr("Test"),
				Category: stringPtr("article"),
			},
			wantErr: false,
		},
		{
			name: "valid request with URL only",
			request: &CreateDocumentRequest{
				URL: "https://example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid request with empty URL",
			request: &CreateDocumentRequest{
				URL: "",
			},
			wantErr: true,
		},
	}

	client := &client{
		baseURL:    "https://api.example.com",
		token:      "test-token",
		httpClient: &http.Client{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := client.CreateDocument(ctx, tt.request)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function to get error type as string
func getErrorType(err error) string {
	if err == nil {
		return ""
	}
	switch err.(type) {
	case *ClientError:
		return "*reader.ClientError"
	case *APIError:
		return "*reader.APIError"
	default:
		return "unknown"
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}