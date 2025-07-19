package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
)

type createCmd struct {
	baseCommand
}

func (*createCmd) Name() string { return "create" }
func (*createCmd) Synopsis() string {
	return "Create a new document"
}
func (*createCmd) Usage() string {
	return `create <url>:
  Create a new document from the specified URL.
  Returns the created document as pretty-printed JSON.
`
}
func (*createCmd) SetFlags(f *flag.FlagSet) {}

func (c *createCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Parse URL from args
	args := f.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", c.Usage())
		return subcommands.ExitUsageError
	}
	url := args[0]

	// Initialize client
	if err := c.initClient(ctx); err != nil {
		printError(err)
		return subcommands.ExitFailure
	}

	// Call CreateDocument API with URL only (Phase 2 - minimal implementation)
	response, err := c.client.CreateDocument(ctx, url, nil)
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