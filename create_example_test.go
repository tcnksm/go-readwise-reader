package reader_test

import (
	"context"
	"log"
	"os"

	reader "github.com/tcnksm/go-readwise-reader"
)

func ExampleClient_CreateDocument() {
	// Get the API token from environment variable
	token := os.Getenv("READWISE_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("READWISE_ACCESS_TOKEN environment variable is required")
	}

	// Create a new client
	client, err := reader.NewClient(token)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create a new document with minimal information (URL only)
	ctx := context.Background()
	response, err := client.CreateDocument(ctx, &reader.CreateDocumentRequest{
		URL: "https://example.com/interesting-article",
	})
	if err != nil {
		log.Fatalf("Failed to create document: %v", err)
	}

	log.Printf("Created document with ID: %s", response.ID)
}

func ExampleClient_CreateDocument_withTitle() {
	token := os.Getenv("READWISE_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("READWISE_ACCESS_TOKEN environment variable is required")
	}

	client, err := reader.NewClient(token)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create a document with custom title
	ctx := context.Background()
	title := "My Custom Title"
	response, err := client.CreateDocument(ctx, &reader.CreateDocumentRequest{
		URL:   "https://example.com/article",
		Title: &title,
	})
	if err != nil {
		log.Fatalf("Failed to create document: %v", err)
	}

	log.Printf("Created document '%s' with ID: %s", response.Title, response.ID)
}

func ExampleClient_CreateDocument_withCategory() {
	token := os.Getenv("READWISE_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("READWISE_ACCESS_TOKEN environment variable is required")
	}

	client, err := reader.NewClient(token)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create a document with title and category
	ctx := context.Background()
	title := "Important Research Paper"
	category := "research"
	
	response, err := client.CreateDocument(ctx, &reader.CreateDocumentRequest{
		URL:      "https://example.com/research-paper",
		Title:    &title,
		Category: &category,
	})
	if err != nil {
		log.Fatalf("Failed to create document: %v", err)
	}

	log.Printf("Created %s document '%s' with ID: %s", 
		response.Category, response.Title, response.ID)
}