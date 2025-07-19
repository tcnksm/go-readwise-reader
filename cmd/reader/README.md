# reader - Readwise Reader CLI

A command-line interface for interacting with the Readwise Reader API.

## Installation

```bash
go install github.com/tcnksm/go-readwise-reader/cmd/reader@latest
```

## Configuration

Set your Readwise access token as an environment variable:

```bash
export READWISE_ACCESS_TOKEN="your-token-here"
```

You can get your access token from: https://readwise.io/access_token

## Usage

### List Documents

List all documents in your "new" location:

```bash
reader list
```

Output is pretty-printed JSON array of documents.

### Create Document

Add a new document by URL:

```bash
reader create https://example.com/blog/interesting-article
```

Returns the created document as pretty-printed JSON.

### Update Document

Update existing document properties:

```bash
reader update --title "New Title" 01k0g64pkqq9w6vh6mz7jtwbvv
reader update --location later --category article 01k0g64pkqq9w6vh6mz7jtwbvv
```

At least one field must be specified. Returns updated document as JSON.

### Delete Document

Remove a document by ID:

```bash
reader delete 01k0g64pkqq9w6vh6mz7jtwbvv
```

Returns success confirmation or error message.