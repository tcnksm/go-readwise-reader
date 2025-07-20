package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/google/subcommands"
	reader "github.com/tcnksm/go-readwise-reader"
)

type createCmd struct {
	baseCommand

	// Flag values
	location string
	notes    string
	summary  string
	title    string
	author   string
	html     string
}

func (*createCmd) Name() string { return "create" }
func (*createCmd) Synopsis() string {
	return "Create a new document"
}
func (*createCmd) Usage() string {
	return `create [flags] <url>:
  Create a new document from the specified URL.
  Returns the created document as pretty-printed JSON.

  Flags:
    -location string     Document location (new, later, archive, feed)
    -notes string        Top-level note for the document (use "-" to read from stdin)
    -summary string      Brief summary of the document
    -title string        Document title
    -author string       Document author
    -html string         Document content in valid HTML format (use "-" to read from stdin)
`
}
func (c *createCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.location, "location", "", "Document location (new, later, archive, feed)")
	f.StringVar(&c.notes, "notes", "", "Top-level note for the document (use '-' to read from stdin)")
	f.StringVar(&c.summary, "summary", "", "Brief summary of the document")
	f.StringVar(&c.title, "title", "", "Document title")
	f.StringVar(&c.author, "author", "", "Document author")
	f.StringVar(&c.html, "html", "", "Document content in valid HTML format")
}

func (c *createCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Parse URL from args
	args := f.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", c.Usage())
		return subcommands.ExitUsageError
	}
	url := args[0]

	// Validate that both -notes - and -html - are not used together
	if c.notes == "-" && c.html == "-" {
		fmt.Fprintf(os.Stderr, "Error: cannot use both -notes and -html with stdin input\n")
		return subcommands.ExitUsageError
	}

	// Initialize client
	if err := c.initClient(ctx); err != nil {
		printError(err)
		return subcommands.ExitFailure
	}

	// Build CreateDocumentRequest from flags
	req := &reader.CreateDocumentRequest{}

	// Set fields only if flags were provided
	if c.location != "" {
		loc := reader.Location(c.location)
		req.Location = loc
	}
	if c.notes != "" {
		if c.notes == "-" {
			// Read notes content from stdin
			notesBytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				printError(fmt.Errorf("failed to read notes from stdin: %w", err))
				return subcommands.ExitFailure
			}
			req.Notes = string(notesBytes)
		} else {
			req.Notes = c.notes
		}
	}
	if c.summary != "" {
		req.Summary = c.summary
	}
	if c.title != "" {
		req.Title = c.title
	}
	if c.author != "" {
		req.Author = c.author
	}
	if c.html != "" {
		if c.html == "-" {
			// Read HTML content from stdin
			htmlBytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				printError(fmt.Errorf("failed to read HTML from stdin: %w", err))
				return subcommands.ExitFailure
			}
			req.HTML = string(htmlBytes)
		} else {
			req.HTML = c.html
		}
		// Enable HTML cleaning when HTML is provided
		req.ShouldCleanHTML = true
	}

	// Call CreateDocument API
	response, err := c.client.CreateDocument(ctx, url, req)
	if err != nil {
		printError(fmt.Errorf("failed to create document: %w", err))
		return subcommands.ExitFailure
	}

	// Output created document JSON
	if err := printJSON(response); err != nil {
		printError(fmt.Errorf("failed to output JSON: %w", err))
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
