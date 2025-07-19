# CLI Implementation Plan: `reader` Command

This document outlines the implementation plan for the `reader` CLI tool that uses the go-readwise-reader library.

## Overview

Develop a command-line tool named `reader` that provides a simple interface to interact with Readwise Reader API.

### Basic Usage
```bash
$ export READWISE_ACCESS_TOKEN="your token"
$ reader list                                    # List documents (pretty JSON)
$ reader create https://example.com/blog/reader  # Create document (returns JSON)
$ reader delete 01k0g64pkqq9w6vh6mz7jtwbvv      # Delete document
```

## Development Principles

1. **Minimal Initial Version**: Start with required fields only, no options
2. **JSON Output**: Pretty-printed JSON for human readability
3. **Error Handling**: Plain text errors to stderr with proper exit codes
4. **Simplicity First**: Additional features in later phases
5. **No Testing Initially**: Manual execution for validation

## Implementation Phases

## Phase 1: Project Setup & Framework

### Directory Structure
```
cmd/reader/
├── README.md           # CLI usage documentation
├── main.go             # Entry point with subcommands setup
├── base.go             # Base command types and shared functionality
├── list.go             # List subcommand implementation
├── create.go           # Create subcommand implementation
├── delete.go           # Delete subcommand implementation
├── config.go           # Environment variable handling
└── output.go           # JSON formatting utilities
```

### Tasks

1. **Create cmd/reader directory**
   - Initialize with README.md containing basic usage

2. **Set up main.go**
   ```go
   // Basic structure using github.com/google/subcommands
   package main
   
   import (
       "context"
       "flag"
       "os"
       
       "github.com/google/subcommands"
   )
   
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
   ```

3. **Create base.go**
   - Common command structure
   - Shared client initialization
   - Token retrieval from environment

4. **Create config.go**
   ```go
   // Handle READWISE_ACCESS_TOKEN environment variable
   func getToken() (string, error) {
       token := os.Getenv("READWISE_ACCESS_TOKEN")
       if token == "" {
           return "", fmt.Errorf("READWISE_ACCESS_TOKEN not set")
       }
       return token, nil
   }
   ```

5. **Create output.go**
   ```go
   // Pretty print JSON responses
   func printJSON(v interface{}) error {
       encoder := json.NewEncoder(os.Stdout)
       encoder.SetIndent("", "  ")
       return encoder.Encode(v)
   }
   
   // Print error to stderr
   func printError(err error) {
       fmt.Fprintln(os.Stderr, "Error:", err)
   }
   ```

## Phase 2: Command Implementations

### List Command (`list.go`)

```go
type listCmd struct{}

func (*listCmd) Name() string     { return "list" }
func (*listCmd) Synopsis() string { return "List documents in new location" }
func (*listCmd) Usage() string {
    return `list:
    List all documents in the new location.
    Output is pretty-printed JSON array.
`
}

func (c *listCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
    // 1. Get token from environment
    // 2. Create client
    // 3. Call ListDocuments with Location: LocationNew
    // 4. Handle pagination (fetch all pages)
    // 5. Output pretty JSON
    // 6. Handle errors appropriately
}
```

**TODO for future phases:**
- Add `--location` flag (new, later, archive, feed)
- Add `--category` flag (article, email, rss, pdf, epub, tweet, video, highlight)
- Add `--tag` flag for filtering by tags
- Add `--limit` flag for pagination control
- Add `--updated-after` flag for time-based filtering

### Create Command (`create.go`)

```go
type createCmd struct{}

func (*createCmd) Name() string     { return "create" }
func (*createCmd) Synopsis() string { return "Create a new document" }
func (*createCmd) Usage() string {
    return `create <url>:
    Create a new document from the specified URL.
    Returns the created document as pretty-printed JSON.
`
}

func (c *createCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
    // 1. Parse URL from args
    // 2. Get token and create client
    // 3. Call CreateDocument with URL only
    // 4. Output created document JSON
    // 5. Handle errors
}
```

**TODO for future phases:**
- Add `--title` flag
- Add `--category` flag
- Add `--tags` flag
- Support reading URLs from stdin for batch operations

### Delete Command (`delete.go`)

```go
type deleteCmd struct{}

func (*deleteCmd) Name() string     { return "delete" }
func (*deleteCmd) Synopsis() string { return "Delete a document" }
func (*deleteCmd) Usage() string {
    return `delete <document-id>:
    Delete the document with the specified ID.
    Returns success message or error.
`
}

func (c *deleteCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
    // 1. Parse document ID from args
    // 2. Get token and create client
    // 3. Call DeleteDocument
    // 4. Output success confirmation
    // 5. Handle errors
}
```

**TODO for future phases:**
- Add `--confirm` flag for safety
- Support multiple IDs for batch deletion

## Phase 3: Documentation & Polish

### README.md for cmd/reader

Include:
1. Installation instructions
   ```bash
   go install github.com/tcnksm/go-readwise-reader/cmd/reader@latest
   ```

2. Environment setup
   ```bash
   export READWISE_ACCESS_TOKEN="your-token-here"
   ```

3. Usage examples for each command

4. Error handling explanation

5. JSON output format examples

### Exit Codes

- `0`: Success
- `1`: General error (API errors, network issues)
- `2`: Usage error (missing arguments, invalid flags)

### Error Output Format

All errors printed to stderr in plain text:
```
Error: READWISE_ACCESS_TOKEN not set
Error: failed to create document: API error (status 400): Invalid URL
Error: document not found
```

## Future Enhancements

### Phase 4: Additional Commands
- [ ] `reader update <id>` - Update document properties

## Dependencies

- `github.com/google/subcommands` - CLI framework
- `github.com/tcnksm/go-readwise-reader` - API client library

## Notes

- Keep initial version minimal and focused
- Gather user feedback before adding complexity
- Maintain consistency with the library's design patterns
- Follow Go best practices and idioms