package reader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// UpdateDocumentRequest represents the request for updating a document
type UpdateDocumentRequest struct {
	// Title is the document title
	Title string `json:"title,omitempty"`

	// Author is the document author
	Author string `json:"author,omitempty"`

	// Summary is the document summary
	Summary string `json:"summary,omitempty"`

	// PublishedDate is the document published date
	PublishedDate *time.Time `json:"published_date,omitempty"`

	// ImageURL is the document image URL
	ImageURL string `json:"image_url,omitempty"`

	// Location is where the document should be stored
	Location Location `json:"location,omitempty"`

	// Category is the document type
	Category Category `json:"category,omitempty"`
}

// UpdateDocumentResponse represents the response from updating a document
type UpdateDocumentResponse struct {
	// ID is the unique identifier of the document
	ID string `json:"id"`

	// URL is the Readwise Reader URL for the document
	URL string `json:"url"`
}

// UpdateDocument updates an existing document in Readwise Reader
func (c *client) UpdateDocument(ctx context.Context, documentID string, req *UpdateDocumentRequest) (*UpdateDocumentResponse, error) {
	if documentID == "" {
		return nil, &ClientError{
			Type:    "invalid_parameter",
			Message: "document ID cannot be empty",
		}
	}

	if req == nil {
		return nil, &ClientError{
			Type:    "invalid_parameter",
			Message: "update request cannot be nil",
		}
	}

	// Build URL
	url := fmt.Sprintf("%s/update/%s/", c.baseURL, documentID)

	// Marshal request body
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request
	httpReq, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", "Token "+c.token)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			apiErr.StatusCode = resp.StatusCode
			apiErr.Message = fmt.Sprintf("unexpected status code: %d", resp.StatusCode)
		} else {
			apiErr.StatusCode = resp.StatusCode
		}
		return nil, &apiErr
	}

	// Decode response
	var response UpdateDocumentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
