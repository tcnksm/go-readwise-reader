package reader_test

import (
	"context"
	"fmt"
	"log"
	"time"

	reader "github.com/tcnksm/go-readwise-reader"
)

func ExampleClient_ListDocuments() {
	// Create client
	client, err := reader.NewClient("your-token-here")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// List all documents
	resp, err := client.ListDocuments(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total documents: %d\n", resp.Count)
	for _, doc := range resp.Results {
		fmt.Printf("- %s (%s)\n", doc.Title, doc.Location)
	}
}

func ExampleClient_ListDocuments_withFilters() {
	// Create client
	client, err := reader.NewClient("your-token-here")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// List documents with filters
	yesterday := time.Now().AddDate(0, 0, -1)
	opts := &reader.ListDocumentsOptions{
		UpdatedAfter: yesterday,
		Location:     reader.LocationNew,
		Category:     reader.CategoryArticle,
		Tag:          "programming",
	}

	resp, err := client.ListDocuments(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d new programming articles from yesterday\n", resp.Count)
}

func ExampleClient_ListDocuments_pagination() {
	// Create client
	client, err := reader.NewClient("your-token-here")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Paginate through all documents
	var allDocuments []reader.Document
	var cursor *string

	for {
		opts := &reader.ListDocumentsOptions{
			PageCursor: func() string {
				if cursor != nil {
					return *cursor
				}
				return ""
			}(),
		}

		resp, err := client.ListDocuments(ctx, opts)
		if err != nil {
			log.Fatal(err)
		}

		allDocuments = append(allDocuments, resp.Results...)

		// Check if there are more pages
		if resp.NextPageCursor == nil {
			break
		}
		cursor = resp.NextPageCursor
	}

	fmt.Printf("Total documents retrieved: %d\n", len(allDocuments))
}
