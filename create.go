package reader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CreateDocumentRequest represents the request to create a document
type CreateDocumentRequest struct {
	// HTML is the document content in valid HTML format (optional)
	HTML string `json:"html,omitempty"`

	// Title is the document title (optional)
	Title string `json:"title,omitempty"`

	// Author is the document author (optional)
	Author string `json:"author,omitempty"`

	// Summary is a brief summary of the document (optional)
	Summary string `json:"summary,omitempty"`

	// PublishedDate is when the document was published (optional)
	PublishedDate *time.Time `json:"published_date,omitempty"`

	// Tags is a list of tags to associate with the document (optional)
	Tags []string `json:"tags,omitempty"`

	// Location is where the document should be initially stored (optional)
	// Valid values: "new", "later", "archive", "feed"
	Location Location `json:"location,omitempty"`

	// Category is the document type (optional)
	// Valid values: "article", "email", "rss", "pdf", "epub", "tweet", "video", "highlight"
	Category Category `json:"category,omitempty"`

	// ImageURL is the URL of an image associated with the document (optional)
	ImageURL string `json:"image_url,omitempty"`

	// Notes is A top-level note of the document (optional)
	Notes string `json:"notes,omitempty"`

	// ShouldCleanHTML instructs Readwise to clean the provided HTML and parse metadata (optional)
	// Only valid when HTML is provided. Defaults to false.
	ShouldCleanHTML bool `json:"should_clean_html,omitempty"`
}

// CreateDocumentResponse represents the response from creating a document
type CreateDocumentResponse struct {
	// ID is the unique identifier of the created document
	ID string `json:"id"`

	// URL is the URL of the created document
	URL string `json:"url"`
}

// CreateDocument creates a new document in Readwise Reader
func (c *client) CreateDocument(ctx context.Context, url string, req *CreateDocumentRequest) (*CreateDocumentResponse, error) {
	if url == "" {
		return nil, &ClientError{
			Type:    "invalid_request",
			Message: "URL is required",
		}
	}

	if req == nil {
		req = &CreateDocumentRequest{}
	}

	// Create request body with URL and other fields
	reqBody := struct {
		URL string `json:"url"`
		*CreateDocumentRequest
	}{
		URL:                   url,
		CreateDocumentRequest: req,
	}

	// Prepare request body
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/save/", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", "Token "+c.token)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code - both 200 and 201 are successful
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
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
	var response CreateDocumentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
