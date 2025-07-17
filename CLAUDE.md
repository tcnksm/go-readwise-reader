# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go client library for the Readwise Reader API. The package follows Go best practices with interface-based design for testability and clear documentation.

## Development Commands

### Build and Test
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -race -cover ./...

# Run tests in verbose mode
go test -v ./...

# Run a specific test
go test -run TestFunctionName ./...

# Generate test coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Code Quality
```bash
# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Run staticcheck (if installed)
staticcheck ./...

# Check for inefficient assignments
ineffassign ./...

# Run all checks before commit
go fmt ./... && go vet ./... && go test ./...
```

### Documentation
```bash
# View documentation locally
go doc -all

# Start local godoc server (Go 1.19+)
go install golang.org/x/pkgsite/cmd/pkgsite@latest
pkgsite -open .
```

## Architecture Guidelines

### File Structure
Each API endpoint gets its own file with related types and methods:
```
/
├── client.go           # Client interface and constructor
├── client_impl.go      # clientImpl struct implementation
├── list.go            # ListDocuments API - types and implementation
├── list_test.go       # ListDocuments tests
├── list_example_test.go # ListDocuments examples
├── create.go          # CreateDocument API - types and implementation
├── create_test.go     # CreateDocument tests
├── get.go             # GetDocument API - types and implementation
├── get_test.go        # GetDocument tests
├── update.go          # UpdateDocument API - types and implementation
├── update_test.go     # UpdateDocument tests
├── delete.go          # DeleteDocument API - types and implementation
├── delete_test.go     # DeleteDocument tests
├── errors.go          # Common error types
├── doc.go             # Package documentation
├── go.mod
├── go.sum
├── README.md
├── CLAUDE.md
└── _example/
    └── main.go        # Runnable example
```

### Client Design Pattern
```go
// client.go - Client interface for mockability
type Client interface {
    ListDocuments(ctx context.Context, opts *ListDocumentsOptions) (*ListDocumentsResponse, error)
    // Other methods added as implemented
}

// client_impl.go - Internal implementation
type clientImpl struct {
    baseURL    string
    token      string
    httpClient *http.Client
}
```

### API File Organization
Each API endpoint file should contain:

```go
// list.go - Example for ListDocuments endpoint

// Request types
type ListDocumentsOptions struct {
    UpdatedAfter  *time.Time `url:"updated__after,omitempty"`
    Location      string     `url:"location,omitempty"`
    Category      string     `url:"category,omitempty"`
    // ... other fields
}

// Response types
type ListDocumentsResponse struct {
    Count    int         `json:"count"`
    Next     *string     `json:"next"`
    Previous *string     `json:"previous"`
    Results  []Document  `json:"results"`
}

// Implementation method
func (c *clientImpl) ListDocuments(ctx context.Context, opts *ListDocumentsOptions) (*ListDocumentsResponse, error) {
    // Implementation
}
```

### API Implementation Process

1. **Research the API endpoint**
   - Visit https://readwise.io/reader_api
   - Study the specific endpoint documentation
   - Note request parameters, response structure, and authentication

2. **Implement one endpoint at a time**
   - Create new file for the endpoint (e.g., `list.go`)
   - Add all related types in the same file
   - Add method to Client interface in `client.go`
   - Implement method in the endpoint file
   - Write comprehensive tests in `*_test.go`
   - Add godoc examples in `*_example_test.go`
   - Ensure all tests pass before moving to next endpoint

3. **Testing Requirements**
   - Unit tests for each method
   - Table-driven tests for multiple scenarios
   - Mock HTTP responses for API testing
   - Test error cases and edge conditions
   - Achieve >80% code coverage

4. **Documentation Standards**
   - Every exported type, function, and method must have godoc comments
   - Include usage examples in `*_example_test.go` files
   - Document all fields in structs
   - Add package-level documentation in `doc.go`

### Code Standards

1. **Error Handling**
   ```go
   // Wrap errors with context
   if err != nil {
       return nil, fmt.Errorf("failed to list documents: %w", err)
   }
   ```

2. **HTTP Client Best Practices**
   - Always use context.Context for cancellation
   - Set appropriate timeouts
   - Handle rate limiting gracefully
   - Parse API errors properly

3. **Testing Pattern**
   ```go
   func TestClientImpl_ListDocuments(t *testing.T) {
       tests := []struct {
           name    string
           opts    *ListDocumentsOptions
           want    *ListDocumentsResponse
           wantErr bool
       }{
           // Test cases
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               // Test implementation with HTTP mocking
           })
       }
   }
   ```

## Development Workflow

1. **Before implementing a new API endpoint:**
   - Read the API documentation thoroughly
   - Create a new file for the endpoint
   - Design the request/response types
   - Write the test cases first (TDD approach)

2. **Implementation steps:**
   - Create endpoint file (e.g., `list.go`)
   - Define request/response types in the same file
   - Update Client interface in `client.go`
   - Implement the method in the endpoint file
   - Write unit tests in `*_test.go`
   - Add examples in `*_example_test.go`
   - Run all tests and ensure they pass

3. **Before committing:**
   ```bash
   # Must all pass:
   go fmt ./...
   go vet ./...
   go test ./...
   ```

## API Reference

Always refer to https://readwise.io/reader_api for:
- Authentication methods
- Available endpoints
- Request/response formats
- Rate limiting information
- Error responses

## Testing Guidelines

1. **Mock HTTP for unit tests** - Don't make real API calls
2. **Test all error scenarios** - Network errors, API errors, invalid inputs
3. **Use table-driven tests** - Makes it easy to add test cases
4. **Test with examples** - Ensure examples in documentation actually work
5. **Integration tests** - Keep separate, only run with specific flag/tag

## Example Implementation Order

1. Start with `client.go` and `client_impl.go` - basic client structure
2. Implement `list.go` - ListDocuments API (already used in example)
3. Add `get.go` - GetDocument API
4. Add `create.go` - CreateDocument API
5. Add `update.go` - UpdateDocument API
6. Add `delete.go` - DeleteDocument API
7. Continue with other endpoints as needed