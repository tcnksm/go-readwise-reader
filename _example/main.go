package main

import (
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

	documents, err := readerClient.ListDocuments(&reader.ListDocumentsOptions{
		Limit: 10,
	})

	if err != nil {
		log.Fatalf("failed to get highlights: %v", err)
	}

	for _, document := range documents {
		fmt.Printf("%v\n", document)
	}
}
