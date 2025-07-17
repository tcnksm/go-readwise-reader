package reader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateDocumentRequest holds the parameters for creating a document
type CreateDocumentRequest struct {
	// URL is the URL of the document to create (required)
	URL string `json:"url"`
	
	// Title is the optional title for the document
	Title *string `json:"title,omitempty"`
	
	// Category is the optional category for the document
	Category *string `json:"category,omitempty"`
}

// CreateDocumentResponse represents the response from creating a document
type CreateDocumentResponse struct {
	// ID is the unique identifier of the created document
	ID string `json:"id"`
	
	// URL is the URL of the document
	URL string `json:"url"`
	
	// Title is the title of the document
	Title string `json:"title"`
	
	// Category is the category of the document
	Category string `json:"category"`
	
	// CreatedAt is the creation timestamp
	CreatedAt string `json:"created_at"`
	
	// UpdatedAt is the last update timestamp
	UpdatedAt string `json:"updated_at"`
}

// CreateDocument creates a new document in Readwise Reader
func (c *client) CreateDocument(ctx context.Context, req *CreateDocumentRequest) (*CreateDocumentResponse, error) {
	if req == nil {
		return nil, &ClientError{
			Type:    "invalid_request",
			Message: "request cannot be nil",
		}
	}
	
	if req.URL == "" {
		return nil, &ClientError{
			Type:    "invalid_request", 
			Message: "URL is required",
		}
	}
	
	// Marshal the request to JSON
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// Create the HTTP request
	url := fmt.Sprintf("%s/documents/", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	httpReq.Header.Set("Authorization", "Token "+c.token)
	httpReq.Header.Set("Content-Type", "application/json")
	
	// Execute the request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	// Check for API errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil {
			return nil, &APIError{
				StatusCode: resp.StatusCode,
				Message:    fmt.Sprintf("failed to create document"),
				Details:    apiErr,
			}
		}
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("failed to create document: status %d", resp.StatusCode),
		}
	}
	
	// Parse the response
	var createResp CreateDocumentResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &createResp, nil
}