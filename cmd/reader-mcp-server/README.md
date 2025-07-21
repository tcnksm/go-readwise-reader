# Readwise Reader MCP Server

The Readwise Reader MCP Server connects AI tools directly to Readwise Reader. This enables AI agent to fetch unread RSS feeds e.g., for creating summary, and manage contents in feed, inbox, and archieve e.g., for moving interesting articles from feed to inbox to read it later, and various automations to make consuming RSS easier.

## Tools

- **save** - Save given URL link to Readwise Reader
  - `url`: URL of the document to save (string, required)
  - `summary`: Brief summary of the document (string, optional)
- **list** - List the documents
  - `location`: Location of the documents. One of new, later, archive, or feed (string, required)
  - `since`: Filter documents updated since duration ago (e.g., 10s, 30m, 24h) (string, optional)
- **move** - Move the documents to different location
  - `id`: ID of the document (given by list tools) (string, required)
  - `location`: Location of the documents. One of new, later, archive, or feed (string, required)


## Installation

```json
{
  "mcp": {
    "servers": {
      "readwise-reader": {
        "command": "/path/to/readwise-reader-mcp-server",
        "env": {
          "READWISE_ACCESS_TOKEN": "<YOUR_TOKEN>"
        }
      }
    }
  }
}
```