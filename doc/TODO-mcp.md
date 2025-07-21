# TODO: Readwise Reader MCP Server Implementation

## Overview

Implement remaining tools for the Readwise Reader MCP server following the existing pattern established in `save.go`.

## Phase 1: Implement and Test 'list' Tool

### Implementation
1. **Create `cmd/reader-mcp-server/list.go`**
   - Tool definition with parameters:
     - `location` (required): new, later, archive, or feed
     - `since` (optional): duration string (e.g., "24h", "30m") using `time.ParseDuration`
     - `limit` (optional): max number of results, default 10
   - Handler implementation:
     - Validate location parameter
     - Parse duration using standard `time.ParseDuration`
     - Call ListDocuments API with appropriate filters
     - Return JSON array of documents with essential fields

2. **Update `main.go`**
   - Register the 'list' tool

3. **Test with MCP Inspector**
   - Verify tool appears in available tools
   - Test various parameter combinations
   - Validate error handling

## Phase 2: Implement and Test 'move' Tool

### Implementation
1. **Create `cmd/reader-mcp-server/move.go`**
   - Tool definition with parameters:
     - `id` (required): document ID from list
     - `location` (required): target location (new, later, archive, or feed)
   - Handler implementation:
     - Validate location parameter
     - Call UpdateDocument API with location update
     - Return success confirmation with updated document details

2. **Update `main.go`**
   - Register the 'move' tool

3. **Test with MCP Inspector**
   - Verify tool appears in available tools
   - Test moving documents between locations
   - Validate error handling for invalid IDs

## Key Design Decisions

1. **No Pagination**: Use `limit` parameter (default 10) instead of complex pagination
2. **Built-in Functions**: Use `time.ParseDuration` instead of custom parsing
3. **Manual Testing**: Use MCP Inspector instead of unit tests
4. **Incremental Registration**: Register and test each tool before implementing the next

## Success Criteria

- [x] 'list' tool implemented and functional
- [x] 'move' tool implemented and functional
- [ ] Both tools tested with MCP Inspector
- [x] Error handling for invalid parameters
- [x] Consistent JSON response format across all tools