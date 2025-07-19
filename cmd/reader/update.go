package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/google/subcommands"
	reader "github.com/tcnksm/go-readwise-reader"
)

type updateCmd struct {
	baseCommand
	title         string
	author        string
	summary       string
	location      string
	category      string
	imageURL      string
	publishedDate string
}

func (*updateCmd) Name() string { return "update" }
func (*updateCmd) Synopsis() string {
	return "Update document properties"
}
func (*updateCmd) Usage() string {
	return `update <document-id> [flags]:
  Update properties of an existing document.
  At least one flag must be provided to specify what to update.
  Returns the updated document as pretty-printed JSON.

Flags:
  -title          Update document title
  -author         Update document author
  -summary        Update document summary
  -location       Update document location (new, later, archive, feed)
  -category       Update document category (article, email, rss, pdf, epub, tweet, video, highlight)
  -image-url      Update document image URL
  -published-date Update document published date (RFC3339 format, e.g., 2023-01-01T00:00:00Z)
`
}

func (c *updateCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.title, "title", "", "Update document title")
	f.StringVar(&c.author, "author", "", "Update document author")
	f.StringVar(&c.summary, "summary", "", "Update document summary")
	f.StringVar(&c.location, "location", "", "Update document location (new, later, archive, feed)")
	f.StringVar(&c.category, "category", "", "Update document category (article, email, rss, pdf, epub, tweet, video, highlight)")
	f.StringVar(&c.imageURL, "image-url", "", "Update document image URL")
	f.StringVar(&c.publishedDate, "published-date", "", "Update document published date (RFC3339 format)")
}

func (c *updateCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Parse document ID from args
	args := f.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", c.Usage())
		return subcommands.ExitUsageError
	}
	documentID := args[0]

	// Check if at least one flag is provided
	if c.title == "" && c.author == "" && c.summary == "" && c.location == "" &&
		c.category == "" && c.imageURL == "" && c.publishedDate == "" {
		printError(fmt.Errorf("at least one field must be specified to update"))
		return subcommands.ExitUsageError
	}

	// Initialize client
	if err := c.initClient(ctx); err != nil {
		printError(err)
		return subcommands.ExitFailure
	}

	// Build update request
	req := &reader.UpdateDocumentRequest{}

	// Set string fields
	if c.title != "" {
		req.Title = c.title
	}
	if c.author != "" {
		req.Author = c.author
	}
	if c.summary != "" {
		req.Summary = c.summary
	}
	if c.imageURL != "" {
		req.ImageURL = c.imageURL
	}

	// Validate and set location
	if c.location != "" {
		switch c.location {
		case "new":
			req.Location = reader.LocationNew
		case "later":
			req.Location = reader.LocationLater
		case "archive":
			req.Location = reader.LocationArchive
		case "feed":
			req.Location = reader.LocationFeed
		default:
			printError(fmt.Errorf("invalid location: %s. Valid values: new, later, archive, feed", c.location))
			return subcommands.ExitUsageError
		}
	}

	// Validate and set category
	if c.category != "" {
		switch c.category {
		case "article":
			req.Category = reader.CategoryArticle
		case "email":
			req.Category = reader.CategoryEmail
		case "rss":
			req.Category = reader.CategoryRSS
		case "pdf":
			req.Category = reader.CategoryPDF
		case "epub":
			req.Category = reader.CategoryEPUB
		case "tweet":
			req.Category = reader.CategoryTweet
		case "video":
			req.Category = reader.CategoryVideo
		case "highlight":
			req.Category = reader.CategoryHighlight
		default:
			printError(fmt.Errorf("invalid category: %s. Valid values: article, email, rss, pdf, epub, tweet, video, highlight", c.category))
			return subcommands.ExitUsageError
		}
	}

	// Parse and set published date
	if c.publishedDate != "" {
		parsedTime, err := time.Parse(time.RFC3339, c.publishedDate)
		if err != nil {
			printError(fmt.Errorf("invalid published date format: %s. Use RFC3339 format like 2023-01-01T00:00:00Z", c.publishedDate))
			return subcommands.ExitUsageError
		}
		req.PublishedDate = &parsedTime
	}

	// Call UpdateDocument API
	response, err := c.client.UpdateDocument(ctx, documentID, req)
	if err != nil {
		printError(fmt.Errorf("failed to update document: %w", err))
		return subcommands.ExitFailure
	}

	// Output updated document JSON
	if err := printJSON(response); err != nil {
		printError(fmt.Errorf("failed to output JSON: %w", err))
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
