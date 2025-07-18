package reader

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Location represents document location in Readwise Reader
type Location string

// Location constants for document filtering
const (
	LocationNew     Location = "new"
	LocationLater   Location = "later"
	LocationArchive Location = "archive"
	LocationFeed    Location = "feed"
)

// Category represents document category/type in Readwise Reader
type Category string

// Category constants for document filtering
const (
	CategoryArticle   Category = "article"
	CategoryEmail     Category = "email"
	CategoryRSS       Category = "rss"
	CategoryPDF       Category = "pdf"
	CategoryEPUB      Category = "epub"
	CategoryTweet     Category = "tweet"
	CategoryVideo     Category = "video"
	CategoryHighlight Category = "highlight"
)

// ListDocumentsOptions holds options for listing documents
type ListDocumentsOptions struct {
	// ID fetches a specific document by ID
	ID string `url:"id,omitempty"`

	// UpdatedAfter fetches documents updated after this time
	UpdatedAfter *time.Time `url:"updatedAfter,omitempty"`

	// Location filters by document location
	Location Location `url:"location,omitempty"`

	// Category filters by document type
	Category Category `url:"category,omitempty"`

	// Tag filters by document tags
	Tag string `url:"tag,omitempty"`

	// PageCursor retrieves the next page of results
	PageCursor string `url:"pageCursor,omitempty"`

	// WithHTMLContent includes HTML content in the response
	WithHTMLContent bool `url:"withHtmlContent,omitempty"`
}

// ListDocumentsResponse represents the response from ListDocuments
type ListDocumentsResponse struct {
	// Count is the total number of documents
	Count int `json:"count"`

	// NextPageCursor is used for pagination
	NextPageCursor *string `json:"nextPageCursor"`

	// Results contains the list of documents
	Results []Document `json:"results"`
}

// Document represents a document in Readwise Reader
type Document struct {
	// ID is the unique identifier of the document
	ID string `json:"id"`

	// URL is the original URL of the document
	URL string `json:"url"`

	// SourceURL is the URL of the source of the document
	SourceURL string `json:"source_url"`

	// Title is the document title
	Title string `json:"title"`

	// Author is the document author
	Author string `json:"author"`

	// Category is the document type
	Category Category `json:"category"`

	// Location is where the document is stored
	Location Location `json:"location"`

	// CreatedAt is when the document was added
	CreatedAt *time.Time `json:"created_at"`

	// UpdatedAt is when the document was last updated
	UpdatedAt *time.Time `json:"updated_at"`

	// Summary of the document
	Summary string `json:"summary"`

	// ReadingProgressPercent is the reading progress (0-100)
	ReadingProgressPercent float64 `json:"reading_progress_percent"`

	// WordCount is the number of words in the document
	WordCount int `json:"word_count"`
}

// ListDocuments retrieves documents from Readwise Reader
func (c *client) ListDocuments(ctx context.Context, opts *ListDocumentsOptions) (*ListDocumentsResponse, error) {
	if opts == nil {
		opts = &ListDocumentsOptions{}
	}

	// Build URL with query parameters
	u, err := url.Parse(c.baseURL + "/list/")
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Add query parameters
	q := u.Query()
	if opts.ID != "" {
		q.Set("id", opts.ID)
	}
	if opts.UpdatedAfter != nil {
		q.Set("updatedAfter", opts.UpdatedAfter.Format(time.RFC3339))
	}
	if opts.Location != "" {
		q.Set("location", string(opts.Location))
	}
	if opts.Category != "" {
		q.Set("category", string(opts.Category))
	}
	if opts.Tag != "" {
		q.Set("tag", opts.Tag)
	}
	if opts.PageCursor != "" {
		q.Set("pageCursor", opts.PageCursor)
	}
	if opts.WithHTMLContent {
		q.Set("withHtmlContent", "true")
	}
	u.RawQuery = q.Encode()

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Token "+c.token)
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
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
	var response ListDocumentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
