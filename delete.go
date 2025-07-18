package reader

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// DeleteDocument deletes a document from Readwise Reader.
// It permanently removes the specified document from the user's account.
func (c *client) DeleteDocument(ctx context.Context, documentID string) error {
	if documentID == "" {
		return &ClientError{
			Type:    "invalid_parameter",
			Message: "document ID cannot be empty",
		}
	}

	url := fmt.Sprintf("%s/delete/%s/", c.baseURL, documentID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Token "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			apiErr.StatusCode = resp.StatusCode
			apiErr.Message = fmt.Sprintf("unexpected status code: %d", resp.StatusCode)
		} else {
			apiErr.StatusCode = resp.StatusCode
		}
		return &apiErr
	}

	return nil
}
