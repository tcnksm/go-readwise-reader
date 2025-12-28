package reader

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// WebhookEventType represents the type of webhook event from Readwise Reader
type WebhookEventType string

// Webhook event type constants
const (
	// EventAnyDocumentCreated is triggered when any new document is added
	EventAnyDocumentCreated WebhookEventType = "reader.any_document.created"

	// EventFeedDocumentCreated is triggered when a new document from RSS or newsletter is added
	EventFeedDocumentCreated WebhookEventType = "reader.feed_document.created"

	// EventNonFeedDocumentCreated is triggered when a manually added document is created
	EventNonFeedDocumentCreated WebhookEventType = "reader.non_feed_document.created"

	// EventDocumentTagsUpdated is triggered when tags are modified on a document
	EventDocumentTagsUpdated WebhookEventType = "reader.document.tags_updated"

	// EventDocumentFinished is triggered when a document is marked as read
	EventDocumentFinished WebhookEventType = "reader.document.finished"

	// EventDocumentArchived is triggered when a document is archived
	EventDocumentArchived WebhookEventType = "reader.document.archived"

	// EventDocumentMovedToLater is triggered when a document is moved to "Later" queue
	EventDocumentMovedToLater WebhookEventType = "reader.document.moved_to_later"

	// EventDocumentMovedToInbox is triggered when a document is moved back to inbox
	EventDocumentMovedToInbox WebhookEventType = "reader.document.moved_to_inbox"

	// EventDocumentShortlisted is triggered when a document is added to favorites
	EventDocumentShortlisted WebhookEventType = "reader.document.shortlisted"
)

// DocumentWebhookPayload represents the webhook payload for document events
type DocumentWebhookPayload struct {
	// EventType specifies which webhook event triggered this payload
	EventType WebhookEventType `json:"event_type"`

	// Secret is a 32-character verification token for authenticating the webhook
	Secret string `json:"secret"`

	// ID is the unique identifier of the document
	ID string `json:"id"`

	// URL is the original URL of the document
	URL string `json:"url"`

	// Title is the document title
	Title string `json:"title"`

	// Author is the document author
	Author string `json:"author"`

	// Source describes how the document was added
	Source *string `json:"source"`

	// Category is the document type
	Category Category `json:"category"`

	// Location is where the document is stored
	Location Location `json:"location"`

	// Tags contains the document tags
	Tags map[string]interface{} `json:"tags"`

	// SiteName is the name of the website
	SiteName string `json:"site_name"`

	// WordCount is the number of words in the document
	WordCount int `json:"word_count"`

	// ReadingTime is the estimated reading time
	ReadingTime string `json:"reading_time"`

	// CreatedAt is when the document was added
	CreatedAt *time.Time `json:"created_at"`

	// UpdatedAt is when the document was last updated
	UpdatedAt *time.Time `json:"updated_at"`

	// PublishedDate is when the document was originally published
	PublishedDate string `json:"published_date"`

	// Summary of the document
	Summary string `json:"summary"`

	// ImageURL is the URL of the document's image
	ImageURL *string `json:"image_url"`

	// Content is the full content of the document
	Content *string `json:"content"`

	// SourceURL is the URL of the source of the document
	SourceURL string `json:"source_url"`

	// Notes is a top-level note of the document
	Notes string `json:"notes"`

	// ParentID is the ID of the parent document if this is a thread
	ParentID *string `json:"parent_id"`

	// ReadingProgress is the reading progress (0-100)
	ReadingProgress float64 `json:"reading_progress"`

	// FirstOpenedAt is when the document was first opened
	FirstOpenedAt *time.Time `json:"first_opened_at"`

	// LastOpenedAt is when the document was last opened
	LastOpenedAt *time.Time `json:"last_opened_at"`

	// SavedAt is when the document was initially saved
	SavedAt *time.Time `json:"saved_at"`

	// LastMovedAt is when the document was last moved between locations
	LastMovedAt *time.Time `json:"last_moved_at"`
}

// DecodeDocumentWebhookPayload decodes a JSON webhook payload into a DocumentWebhookPayload struct.
// This function is typically used in webhook handlers to parse incoming POST request bodies.
func DecodeDocumentWebhookPayload(r io.Reader) (*DocumentWebhookPayload, error) {
	var payload DocumentWebhookPayload
	if err := json.NewDecoder(r).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode webhook payload: %w", err)
	}
	return &payload, nil
}
