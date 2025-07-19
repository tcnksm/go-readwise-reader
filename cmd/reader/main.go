package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
)


type createCmd struct{}

func (*createCmd) Name() string             { return "create" }
func (*createCmd) Synopsis() string         { return "Create document (Phase 2)" }
func (*createCmd) Usage() string            { return "create: Create document (Phase 2)\n" }
func (*createCmd) SetFlags(f *flag.FlagSet) {}
func (*createCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	fmt.Fprintln(os.Stderr, "create command not yet implemented (Phase 2)")
	return subcommands.ExitFailure
}

type deleteCmd struct{}

func (*deleteCmd) Name() string             { return "delete" }
func (*deleteCmd) Synopsis() string         { return "Delete document (Phase 2)" }
func (*deleteCmd) Usage() string            { return "delete: Delete document (Phase 2)\n" }
func (*deleteCmd) SetFlags(f *flag.FlagSet) {}
func (*deleteCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	fmt.Fprintln(os.Stderr, "delete command not yet implemented (Phase 2)")
	return subcommands.ExitFailure
}

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	subcommands.Register(&listCmd{}, "")
	subcommands.Register(&createCmd{}, "")
	subcommands.Register(&deleteCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
