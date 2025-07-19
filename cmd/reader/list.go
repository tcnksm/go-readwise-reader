package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
	reader "github.com/tcnksm/go-readwise-reader"
)

type listCmd struct {
	baseCommand
}

func (*listCmd) Name() string { return "list" }
func (*listCmd) Synopsis() string {
	return "List documents in new location"
}
func (*listCmd) Usage() string {
	return `list:
  List all documents in the new location.
  Output is pretty-printed JSON array.
`
}
func (*listCmd) SetFlags(f *flag.FlagSet) {}

func (c *listCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Initialize client
	if err := c.initClient(ctx); err != nil {
		printError(err)
		return subcommands.ExitFailure
	}

	// Set up options for ListDocuments - Phase 2 lists only new location
	opts := &reader.ListDocumentsOptions{
		Location: reader.LocationNew,
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