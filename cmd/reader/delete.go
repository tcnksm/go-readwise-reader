package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
)

type deleteCmd struct {
	baseCommand
}

func (*deleteCmd) Name() string { return "delete" }
func (*deleteCmd) Synopsis() string {
	return "Delete a document"
}
func (*deleteCmd) Usage() string {
	return `delete <document-id>:
  Delete the document with the specified ID.
  Returns success message or error.
`
}
func (*deleteCmd) SetFlags(f *flag.FlagSet) {}

func (c *deleteCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Parse document ID from args
	args := f.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", c.Usage())
		return subcommands.ExitUsageError
	}
	documentID := args[0]

	// Initialize client
	if err := c.initClient(ctx); err != nil {
		printError(err)
		return subcommands.ExitFailure
	}

	// Call DeleteDocument API
	err := c.client.DeleteDocument(ctx, documentID)
	if err != nil {
		printError(fmt.Errorf("failed to delete document: %w", err))
		return subcommands.ExitFailure
	}

	// Output success confirmation
	fmt.Printf("Document %s deleted successfully\n", documentID)

	return subcommands.ExitSuccess
}
