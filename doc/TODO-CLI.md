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

## Phase 1: Project Setup & Framework ✅ COMPLETED

### Directory Structure ✅ IMPLEMENTED
```
cmd/reader/
├── README.md           # CLI usage documentation ✅
├── go.mod              # Separate module for CLI dependencies ✅
├── go.sum              # Dependency checksums ✅
├── main.go             # Entry point with subcommands setup ✅
├── base.go             # Base command types, shared functionality, config & output ✅
├── list.go             # List subcommand implementation (Phase 2)
├── create.go           # Create subcommand implementation (Phase 2)
└── delete.go           # Delete subcommand implementation (Phase 2)
```

**Notes**: 
- Combined config.go and output.go functionality into base.go for simplicity
- **Separate Go Module**: CLI has its own go.mod to avoid mixing CLI dependencies with the library package

### Tasks ✅ ALL COMPLETED

1. **Create cmd/reader directory** ✅
   - Initialize with README.md containing basic usage ✅

2. **Set up main.go** ✅
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

3. **Create base.go** ✅
   - Common command structure ✅
   - Shared client initialization ✅
   - Token retrieval from environment ✅
   - JSON formatting utilities ✅

4. **~~Create config.go~~** ✅ (Combined into base.go)
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

5. **~~Create output.go~~** ✅ (Combined into base.go)
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

## Phase 2: Command Implementations ✅ COMPLETED

### List Command (`list.go`) ✅ IMPLEMENTED

**Core Implementation:**
- ✅ Lists documents with flexible filtering options
- ✅ Outputs pretty-printed JSON array of documents
- ✅ Uses baseCommand pattern for client initialization
- ✅ Handles errors with proper exit codes
- ✅ TODO added for future pagination implementation

**Filtering Flags:** ✅ IMPLEMENTED
- ✅ `--location` flag (new, later, archive, feed) - defaults to "new"
- ✅ `--category` flag (article, email, rss, pdf, epub, tweet, video, highlight)
- ✅ `--tag` flag for filtering by single tag name
- ✅ `--since` flag for time-based filtering using Go duration format
- ✅ Flag validation with helpful error messages
- ✅ Support for combining multiple filters

**Time-based Filtering (`--since`):** ✅ IMPLEMENTED
- Uses standard Go `time.ParseDuration` for consistent behavior
- Supports: seconds (`30s`), minutes (`30m`), hours (`24h`), complex (`1h30m`)
- Calculates time "duration ago" from current time for `UpdatedAfter` API
- Proper error handling for invalid duration formats

**Usage Examples:**
```bash
./reader list                           # Default: new location
./reader list -location later           # Later location only
./reader list -category article         # Articles only  
./reader list -location new -category rss  # Combined filters
./reader list -tag "ai"                 # Tag filtering
./reader list -since 1h                 # Documents updated in last hour
./reader list -since 24h                # Documents updated in last day
./reader list -since 1h30m              # Documents updated in last 1.5 hours
./reader list -since 24h -category article  # Recent articles
```

**TODO for future phases:**
- Add `--limit` flag for pagination control
- Add pagination support (currently single page only)

### Create Command (`create.go`) ✅ IMPLEMENTED

- ✅ Accepts URL as command line argument
- ✅ Creates document using CreateDocument API with minimal options
- ✅ Outputs created document as pretty-printed JSON
- ✅ Handles errors with proper exit codes and usage messages
- ✅ Uses baseCommand pattern for client initialization

**TODO for future phases:**
- Add `--title` flag
- Add `--category` flag
- Add `--tags` flag
- Support reading URLs from stdin for batch operations

### Delete Command (`delete.go`) ✅ IMPLEMENTED

- ✅ Accepts document ID as command line argument  
- ✅ Deletes document using DeleteDocument API
- ✅ Outputs success confirmation message
- ✅ Handles errors with proper exit codes and usage messages
- ✅ Uses baseCommand pattern for client initialization

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

### Update Command (`update.go`) ✅ IMPLEMENTED

**Core Implementation:**
- ✅ Updates existing document properties via PATCH API
- ✅ Requires document ID as positional argument  
- ✅ Requires at least one field flag to be provided
- ✅ Returns updated document as pretty-printed JSON
- ✅ Uses baseCommand pattern for client initialization

**Update Flags:** ✅ ALL FIELDS SUPPORTED
- ✅ `--title` - Update document title
- ✅ `--author` - Update document author
- ✅ `--summary` - Update document summary
- ✅ `--location` - Update location (new, later, archive, feed)
- ✅ `--category` - Update category (article, email, rss, pdf, epub, tweet, video, highlight)
- ✅ `--image-url` - Update document image URL
- ✅ `--published-date` - Update published date (RFC3339 format)

**Validation:**
- ✅ Requires at least one update field (fails if no flags provided)
- ✅ Validates location and category enum values
- ✅ Validates published date RFC3339 format
- ✅ Provides helpful error messages for invalid values

**Usage Examples:**
```bash
# Single field updates
./reader update --title "New Title" <document-id>
./reader update --location later <document-id>
./reader update --category article <document-id>

# Multiple field updates
./reader update --title "New Title" --author "John Doe" <document-id>
./reader update --location archive --category article <document-id>

# Date update (RFC3339 format)
./reader update --published-date "2023-01-01T00:00:00Z" <document-id>

# Error cases (properly handled)
./reader update <document-id>                    # Error: at least one field required
./reader update --location invalid <document-id> # Error: invalid location value
```

### Phase 4: Additional Commands
- [x] `reader update <id>` - Update document properties ✅ IMPLEMENTED

## Dependencies

**CLI Module** (`cmd/reader/go.mod`):
- `github.com/google/subcommands` - CLI framework
- `github.com/tcnksm/go-readwise-reader` - API client library

**Library Module** (`go.mod`):
- No CLI-specific dependencies (keeps library clean)

## Notes

- Keep initial version minimal and focused
- Gather user feedback before adding complexity
- Maintain consistency with the library's design patterns
- Follow Go best practices and idioms