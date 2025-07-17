package reader

import "context"

// ListDocumentsOptions holds options for listing documents
type ListDocumentsOptions struct {
	// TODO: Add actual fields like Limit, Category, UpdatedAfter
}

// ListDocumentsResponse represents the response from ListDocuments
type ListDocumentsResponse struct {
	// TODO: Add actual fields like Count, Next, Previous, Results
}

// Document represents a document in Readwise Reader
type Document struct {
	// TODO: Add actual fields like ID, Title, URL, etc.
}

// ListDocuments retrieves documents from Readwise Reader
func (c *client) ListDocuments(ctx context.Context, opts *ListDocumentsOptions) (*ListDocumentsResponse, error) {
	// TODO: Implement actual API call
	return nil, &ClientError{
		Type:    "not_implemented",
		Message: "ListDocuments not implemented yet",
	}
}