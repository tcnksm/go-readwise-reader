package reader_test

import (
	"context"
	"fmt"
	"log"
	"time"

	reader "github.com/tcnksm/go-readwise-reader"
)

func ExampleClient_CreateDocument() {
	client, err := reader.NewClient("your-token-here")
	if err != nil {
		log.Fatal(err)
	}

	url := "https://example.com/article"
	req := &reader.CreateDocumentRequest{
		Title: "Interesting Article",
		Tags:  []string{"golang", "api"},
	}

	ctx := context.Background()
	resp, err := client.CreateDocument(ctx, url, req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created document with ID: %s\n", resp.ID)
	fmt.Printf("Document URL: %s\n", resp.URL)
}

func ExampleClient_CreateDocument_withAllFields() {
	client, err := reader.NewClient("your-token-here")
	if err != nil {
		log.Fatal(err)
	}

	publishedDate := time.Date(2023, 12, 1, 10, 0, 0, 0, time.UTC)

	url := "https://example.com/complete-article"
	req := &reader.CreateDocumentRequest{
		HTML:          "<html><body><h1>Complete Article</h1><p>Content here</p></body></html>",
		Title:         "Complete Article Example",
		Author:        "Jane Doe",
		Summary:       "This is a complete example of creating a document with all fields",
		PublishedDate: &publishedDate,
		Tags:          []string{"example", "complete", "documentation"},
		Location:      reader.LocationNew,
		Category:      reader.CategoryArticle,
		ImageURL:      "https://example.com/image.jpg",
	}

	ctx := context.Background()
	resp, err := client.CreateDocument(ctx, url, req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created document with ID: %s\n", resp.ID)
	fmt.Printf("Document URL: %s\n", resp.URL)
}
