// Package reader provides a Go client library for the Readwise Reader API.
//
// This package allows you to interact with Readwise Reader's API to manage
// your saved documents, articles, and reading list.
//
// Basic usage:
//
//	client, err := reader.NewClient(token)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	documents, err := client.ListDocuments(ctx, &reader.ListDocumentsOptions{
//		Limit: 10,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Authentication:
//
// You need a Readwise Reader API token to use this package. You can get one
// from https://readwise.io/reader_api
//
// Set your token as an environment variable:
//
//	export READWISE_ACCESS_TOKEN="your-token-here"
package reader
