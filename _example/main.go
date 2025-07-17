package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	reader "github.com/tcnksm/go-readwise-reader"
)

func main() {
	token := os.Getenv("READWISE_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("READWISE_ACCESS_TOKEN is not set")
	}

	readerClient, err := reader.NewClient(token)
	if err != nil {
		log.Fatalf("failed to create reader client: %v", err)
	}

	ctx := context.Background()

	// List all documents
	response, err := readerClient.ListDocuments(ctx, nil)
	if err != nil {
		log.Fatalf("failed to list documents: %v", err)
	}

	fmt.Printf("Found %d documents\n", response.Count)
	for _, document := range response.Results {
		fmt.Printf("- %s: %s\n", document.Title, document.URL)
	}

	// List documents with filters
	yesterday := time.Now().AddDate(0, 0, -1)
	filteredResponse, err := readerClient.ListDocuments(ctx, &reader.ListDocumentsOptions{
		UpdatedAfter: &yesterday,
		Location:     reader.LocationNew,
		Category:     reader.CategoryArticle,
	})
	if err != nil {
		log.Fatalf("failed to list filtered documents: %v", err)
	}

	fmt.Printf("Found %d new articles from yesterday\n", filteredResponse.Count)
}
