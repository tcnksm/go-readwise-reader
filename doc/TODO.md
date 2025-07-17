# Implementation Plan

This document outlines the step-by-step implementation plan for the go-readwise-reader package.

## Development Principles

1. **Incremental Development**: Start with minimal features, add complexity later
2. **Branch Strategy**: Each API endpoint gets its own feature branch
3. **Interface First**: Design interface in README.md, get approval, then implement
4. **Test Driven**: Write tests first, then implementation
5. **Documentation**: Update docs with each change

## Branch Strategy

- `main`: Stable, tested code only
- `feature/client-setup`: Basic client structure
- `feature/list-documents`: List Documents API
- `feature/create-document`: Create Document API  
- `feature/update-document`: Update Document API
- `feature/delete-document`: Delete Document API

## Phase 0: Initial Setup ✓

**Branch**: `main`
**Status**: Complete

- [x] Create go.mod
- [x] Create example directory
- [x] Create README.md with basic info
- [x] Create CLAUDE.md with development guidelines
- [x] Create doc/TODO.md (this file)

## Phase 1: Basic Client Setup ✓

**Branch**: `feature/client-setup`
**Goal**: Establish the foundation for the API client
**Status**: Complete

### Steps:

1. **Design Interface** (README.md update)
   - [x] Basic client interface design
   - [x] Show minimal example usage
   - [x] Get approval before implementation

2. **Implementation**
   - [x] Create `client.go` with Client interface
   - [x] Create `client_impl.go` with basic implementation
   - [x] Create `errors.go` for common error types
   - [x] Create `doc.go` for package documentation

3. **Testing**
   - [x] Unit tests for client creation
   - [x] Example test for documentation

4. **Documentation**
   - [x] Update README.md with actual interface
   - [x] Ensure godoc is clean

## Phase 2: List Documents API ✓

**Branch**: `feature/list-documents`
**Goal**: Implement the first API endpoint
**Status**: Complete

### Steps:

1. **Research & Design**
   - [x] Study https://readwise.io/reader_api for list endpoint
   - [x] Design interface in README.md
   - [x] Define ListDocumentsOptions with comprehensive fields
   - [x] Create typed constants for Location and Category

2. **Implementation**
   - [x] Create `list.go` with types and implementation
   - [x] Implement comprehensive parameters:
     - `ID` (specific document lookup)
     - `UpdatedAfter` (time-based filter)
     - `Location` (typed enum: new, later, archive, feed)
     - `Category` (typed enum: article, email, rss, pdf, epub, tweet, video, highlight)
     - `Tag` (tag-based filter)
     - `PageCursor` (pagination)
     - `WithHTMLContent` (content inclusion)
   - [x] Handle pagination with NextPageCursor
   - [x] Proper HTTP client implementation with context support
   - [x] Error handling for API responses

3. **Testing**
   - [x] Create `list_test.go` with comprehensive table-driven tests
   - [x] Mock HTTP responses for various scenarios
   - [x] Test error cases and edge conditions
   - [x] Test context cancellation
   - [x] Create `list_example_test.go` for godoc examples

4. **Integration**
   - [x] Update example in `_example/main.go`
   - [x] Update README.md with proper constant usage
   - [x] All tests passing
   - [x] Code formatted and linted

## Phase 3: Create Document API

**Branch**: `feature/create-document`
**Goal**: Add ability to create documents

### Steps:

1. **Research & Design**
   - [ ] Study create endpoint in API docs
   - [ ] Design interface in README.md
   - [ ] Define minimal CreateDocumentRequest
   - [ ] Get approval

2. **Implementation**
   - [ ] Create `create.go` with types and implementation
   - [ ] Start with required fields only:
     - `URL` (required)
     - `Title` (optional)
     - `Category` (optional)

3. **Testing**
   - [ ] Create `create_test.go`
   - [ ] Create `create_example_test.go`

## Phase 4: Update Document API

**Branch**: `feature/update-document`
**Goal**: Add ability to update existing documents

### Steps:

1. **Research & Design**
   - [ ] Study update endpoint
   - [ ] Design interface in README.md
   - [ ] Define UpdateDocumentRequest
   - [ ] Get approval

2. **Implementation**
   - [ ] Create `update.go`
   - [ ] Handle partial updates
   - [ ] Support common fields:
     - `Title`
     - `Location` (archive, later, etc.)
     - `Tags`

3. **Testing**
   - [ ] Create `update_test.go`
   - [ ] Create `update_example_test.go`

## Phase 5: Delete Document API

**Branch**: `feature/delete-document`
**Goal**: Add ability to delete documents

### Steps:

1. **Research & Design**
   - [ ] Study delete endpoint
   - [ ] Design interface in README.md
   - [ ] Get approval

2. **Implementation**
   - [ ] Create `delete.go`
   - [ ] Handle soft delete if supported

3. **Testing**
   - [ ] Create `delete_test.go`
   - [ ] Create `delete_example_test.go`

## Future Enhancements (Not in initial scope)

After the basic CRUD operations are complete:

1. **Advanced Features**
   - [ ] Get single document details
   - [ ] Bulk operations
   - [ ] Export functionality
   - [ ] Highlights API
   - [ ] Tags management

2. **Client Improvements**
   - [ ] Rate limiting handling
   - [ ] Retry logic
   - [ ] Request/Response logging
   - [ ] Custom HTTP client support
   - [ ] Context cancellation

3. **Testing & Quality**
   - [ ] Integration test suite
   - [ ] Performance benchmarks
   - [ ] CI/CD setup
   - [ ] Code coverage badges

## Progress Tracking

Each phase should result in:
1. A working feature
2. Comprehensive tests (>80% coverage)
3. Updated documentation
4. Clean git history
5. PR merged to main

## Notes

- Each phase builds on the previous one
- We start with read operations before write operations
- Complex features are deferred to keep PRs small
- Interface design always comes before implementation
- User approval required before coding begins