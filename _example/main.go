package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	// Example 1: Create a new document
	fmt.Println("Creating a new document...")
	title := "Example Article"
	category := "article"
	
	createResp, err := readerClient.CreateDocument(ctx, &reader.CreateDocumentRequest{
		URL:      "https://example.com/interesting-article",
		Title:    &title,
		Category: &category,
	})
	if err != nil {
		log.Fatalf("failed to create document: %v", err)
	}
	
	fmt.Printf("Created document: ID=%s, Title=%s, URL=%s\n", 
		createResp.ID, createResp.Title, createResp.URL)

	// Example 2: List documents
	fmt.Println("\nListing documents...")
	documents, err := readerClient.ListDocuments(ctx, &reader.ListDocumentsOptions{})
	if err != nil {
		log.Fatalf("failed to get documents: %v", err)
	}

	for _, document := range documents.Results {
		fmt.Printf("Document: %+v\n", document)
	}
}
