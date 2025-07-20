package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/google/subcommands"
	reader "github.com/tcnksm/go-readwise-reader"
)

type listCmd struct {
	baseCommand
	location string
	category string
	tag      string
	since    string
	html     bool
}

func (*listCmd) Name() string { return "list" }
func (*listCmd) Synopsis() string {
	return "List documents with optional filtering"
}
func (*listCmd) Usage() string {
	return `list [flags]:
  List documents with optional filtering.
  Output is pretty-printed JSON array.

Flags:
  -location   Filter by location (new, later, archive, feed). Default: new
  -category   Filter by category (article, email, rss, pdf, epub, tweet, video, highlight)
  -tag        Filter by tag name
  -since      Filter documents updated since duration ago (e.g., 10s, 30m, 24h)
  -html       Include HTML content in the response
`
}
func (c *listCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.location, "location", "new", "Filter by location (new, later, archive, feed)")
	f.StringVar(&c.category, "category", "", "Filter by category (article, email, rss, pdf, epub, tweet, video, highlight)")
	f.StringVar(&c.tag, "tag", "", "Filter by tag name")
	f.StringVar(&c.since, "since", "", "Filter documents updated since duration ago (e.g., 10s, 30m, 24h)")
	f.BoolVar(&c.html, "html", false, "Include HTML content in the response")
}

func (c *listCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Initialize client
	if err := c.initClient(ctx); err != nil {
		printError(err)
		return subcommands.ExitFailure
	}

	// Validate location
	var location reader.Location
	switch c.location {
	case "new":
		location = reader.LocationNew
	case "later":
		location = reader.LocationLater
	case "archive":
		location = reader.LocationArchive
	case "feed":
		location = reader.LocationFeed
	default:
		printError(fmt.Errorf("invalid location: %s. Valid values: new, later, archive, feed", c.location))
		return subcommands.ExitUsageError
	}

	// Validate category
	var category reader.Category
	if c.category != "" {
		switch c.category {
		case "article":
			category = reader.CategoryArticle
		case "email":
			category = reader.CategoryEmail
		case "rss":
			category = reader.CategoryRSS
		case "pdf":
			category = reader.CategoryPDF
		case "epub":
			category = reader.CategoryEPUB
		case "tweet":
			category = reader.CategoryTweet
		case "video":
			category = reader.CategoryVideo
		case "highlight":
			category = reader.CategoryHighlight
		default:
			printError(fmt.Errorf("invalid category: %s. Valid values: article, email, rss, pdf, epub, tweet, video, highlight", c.category))
			return subcommands.ExitUsageError
		}
	}

	// Parse since duration if provided
	var updatedAfter *time.Time
	if c.since != "" {
		duration, err := time.ParseDuration(c.since)
		if err != nil {
			printError(fmt.Errorf("invalid duration format: %s. Use formats like 10s, 30m, 24h", c.since))
			return subcommands.ExitUsageError
		}

		// Calculate time "duration" ago from now
		timeAgo := time.Now().Add(-duration)
		updatedAfter = &timeAgo
	}

	// Set up options for ListDocuments
	opts := &reader.ListDocumentsOptions{
		Location:        location,
		Category:        category,
		Tag:             c.tag,
		UpdatedAfter:    updatedAfter,
		WithHTMLContent: c.html,
	}

	// Call ListDocuments API
	response, err := c.client.ListDocuments(ctx, opts)
	if err != nil {
		printError(fmt.Errorf("failed to list documents: %w", err))
		return subcommands.ExitFailure
	}

	// Output documents as pretty JSON
	// TODO: Handle pagination in future implementation
	if err := printJSON(response.Results); err != nil {
		printError(fmt.Errorf("failed to output JSON: %w", err))
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
