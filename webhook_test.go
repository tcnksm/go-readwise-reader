package reader

import (
	"strings"
	"testing"
	"time"
)

func TestDecodeDocumentWebhookPayload(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		wantErr  bool
		validate func(*testing.T, *DocumentWebhookPayload)
	}{
		{
			name: "valid webhook payload with all fields",
			json: `{
				"event_type": "reader.any_document.created",
				"secret": "abc123def456ghi789jkl012mno345pq",
				"id": "doc123",
				"url": "https://example.com/article",
				"title": "Test Article",
				"author": "John Doe",
				"source": "RSS Feed",
				"category": "article",
				"location": "new",
				"tags": {"tech": true, "programming": true},
				"site_name": "Example Site",
				"word_count": 1500,
				"reading_time": "7 min",
				"created_at": "2024-01-15T10:00:00Z",
				"updated_at": "2024-01-16T12:00:00Z",
				"published_date": "2024-01-14",
				"summary": "A test article summary",
				"image_url": "https://example.com/image.jpg",
				"content": "Full article content here",
				"source_url": "https://source.example.com",
				"notes": "My notes",
				"parent_id": "parent123",
				"reading_progress": 45.5,
				"first_opened_at": "2024-01-15T11:00:00Z",
				"last_opened_at": "2024-01-16T14:00:00Z",
				"saved_at": "2024-01-15T09:00:00Z",
				"last_moved_at": "2024-01-15T10:30:00Z"
			}`,
			wantErr: false,
			validate: func(t *testing.T, p *DocumentWebhookPayload) {
				if p.EventType != EventAnyDocumentCreated {
					t.Errorf("EventType = %v, want %v", p.EventType, EventAnyDocumentCreated)
				}
				if p.Secret != "abc123def456ghi789jkl012mno345pq" {
					t.Errorf("Secret = %v, want %v", p.Secret, "abc123def456ghi789jkl012mno345pq")
				}
				if p.ID != "doc123" {
					t.Errorf("ID = %v, want %v", p.ID, "doc123")
				}
				if p.Title != "Test Article" {
					t.Errorf("Title = %v, want %v", p.Title, "Test Article")
				}
				if p.Category != CategoryArticle {
					t.Errorf("Category = %v, want %v", p.Category, CategoryArticle)
				}
				if p.Location != LocationNew {
					t.Errorf("Location = %v, want %v", p.Location, LocationNew)
				}
				if p.WordCount != 1500 {
					t.Errorf("WordCount = %v, want %v", p.WordCount, 1500)
				}
				if p.ReadingProgress != 45.5 {
					t.Errorf("ReadingProgress = %v, want %v", p.ReadingProgress, 45.5)
				}
			},
		},
		{
			name: "webhook payload with null optional fields",
			json: `{
				"event_type": "reader.document.archived",
				"secret": "testsecret123456789012345678901",
				"id": "doc456",
				"url": "https://example.com",
				"title": "Minimal Doc",
				"author": "",
				"source": null,
				"category": "article",
				"location": "archive",
				"tags": {},
				"site_name": "Site",
				"word_count": 100,
				"reading_time": "1 min",
				"created_at": "2024-01-15T10:00:00Z",
				"updated_at": "2024-01-15T10:00:00Z",
				"published_date": "",
				"summary": "",
				"image_url": null,
				"content": null,
				"source_url": "",
				"notes": "",
				"parent_id": null,
				"reading_progress": 0,
				"first_opened_at": null,
				"last_opened_at": null,
				"saved_at": "2024-01-15T10:00:00Z",
				"last_moved_at": "2024-01-15T10:00:00Z"
			}`,
			wantErr: false,
			validate: func(t *testing.T, p *DocumentWebhookPayload) {
				if p.EventType != EventDocumentArchived {
					t.Errorf("EventType = %v, want %v", p.EventType, EventDocumentArchived)
				}
				if p.Source != nil {
					t.Errorf("Source = %v, want nil", p.Source)
				}
				if p.ImageURL != nil {
					t.Errorf("ImageURL = %v, want nil", p.ImageURL)
				}
				if p.Content != nil {
					t.Errorf("Content = %v, want nil", p.Content)
				}
				if p.ParentID != nil {
					t.Errorf("ParentID = %v, want nil", p.ParentID)
				}
				if p.FirstOpenedAt != nil {
					t.Errorf("FirstOpenedAt = %v, want nil", p.FirstOpenedAt)
				}
				if p.LastOpenedAt != nil {
					t.Errorf("LastOpenedAt = %v, want nil", p.LastOpenedAt)
				}
			},
		},
		{
			name: "document moved to later event",
			json: `{
				"event_type": "reader.document.moved_to_later",
				"secret": "secret12345678901234567890123456",
				"id": "doc789",
				"url": "https://example.com/moved",
				"title": "Moved Document",
				"author": "Jane Smith",
				"source": "Manual",
				"category": "article",
				"location": "later",
				"tags": {},
				"site_name": "Example",
				"word_count": 500,
				"reading_time": "3 min",
				"created_at": "2024-01-10T10:00:00Z",
				"updated_at": "2024-01-15T10:00:00Z",
				"published_date": "2024-01-10",
				"summary": "Test",
				"image_url": null,
				"content": null,
				"source_url": "https://example.com/moved",
				"notes": "",
				"parent_id": null,
				"reading_progress": 10.5,
				"first_opened_at": "2024-01-10T11:00:00Z",
				"last_opened_at": "2024-01-14T09:00:00Z",
				"saved_at": "2024-01-10T10:00:00Z",
				"last_moved_at": "2024-01-15T10:00:00Z"
			}`,
			wantErr: false,
			validate: func(t *testing.T, p *DocumentWebhookPayload) {
				if p.EventType != EventDocumentMovedToLater {
					t.Errorf("EventType = %v, want %v", p.EventType, EventDocumentMovedToLater)
				}
				if p.Location != LocationLater {
					t.Errorf("Location = %v, want %v", p.Location, LocationLater)
				}
				expectedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
				if p.LastMovedAt == nil || !p.LastMovedAt.Equal(expectedTime) {
					t.Errorf("LastMovedAt = %v, want %v", p.LastMovedAt, expectedTime)
				}
			},
		},
		{
			name:    "invalid JSON",
			json:    `{"event_type": "reader.any_document.created", invalid}`,
			wantErr: true,
		},
		{
			name:    "empty JSON object",
			json:    `{}`,
			wantErr: false,
			validate: func(t *testing.T, p *DocumentWebhookPayload) {
				if p.EventType != "" {
					t.Errorf("EventType = %v, want empty string", p.EventType)
				}
				if p.Secret != "" {
					t.Errorf("Secret = %v, want empty string", p.Secret)
				}
			},
		},
		{
			name:    "malformed JSON",
			json:    `not json at all`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.json)
			payload, err := DecodeDocumentWebhookPayload(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeDocumentWebhookPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if payload == nil {
				t.Fatal("DecodeDocumentWebhookPayload() returned nil payload")
			}

			if tt.validate != nil {
				tt.validate(t, payload)
			}
		})
	}
}
