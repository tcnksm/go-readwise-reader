package reader

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://readwise.io/api/v3"
	defaultTimeout = 30 * time.Second
)

// Client interface for interacting with Readwise Reader API
type Client interface {
	ListDocuments(ctx context.Context, opts *ListDocumentsOptions) (*ListDocumentsResponse, error)
	CreateDocument(ctx context.Context, req *CreateDocumentRequest) (*CreateDocumentResponse, error)
  UpdateDocument(ctx context.Context, documentID string, req *UpdateDocumentRequest) (*UpdateDocumentResponse, error)
}

// client is the implementation of the Client interface
type client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewClient creates a new Readwise Reader client
func NewClient(token string) (Client, error) {
	if token == "" {
		return nil, &ClientError{
			Type:    "invalid_token",
			Message: "token cannot be empty",
		}
	}

	return &client{
		baseURL: defaultBaseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}, nil
}

// ClientError represents an error from the client
type ClientError struct {
	Type    string
	Message string
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// APIError represents an error from the Readwise Reader API
type APIError struct {
	StatusCode int
	Message    string
	Details    map[string]interface{}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}
