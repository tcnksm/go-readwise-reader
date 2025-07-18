# go-readwise-reader

[![GoDoc](https://pkg.go.dev/badge/github.com/tcnksm/go-readwise-reader.svg)](https://pkg.go.dev/github.com/tcnksm/go-readwise-reader)

`go-readwise-reader` is a Go package to interact with [Readwise Reader API](https://readwise.io/reader_api).

*NOTE*: This is mainly personal usage for own automation.


## Install

```bash
$ go install github.com/tcnksm/go-readwise-reader@latest
```

## Usage 

See [Go doc](https://pkg.go.dev/github.com/tcnksm/go-readwise-reader). 

## Example

The following is an example to list documents:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	reader "github.com/tcnksm/go-readwise-reader"
)

func main() {
	token := os.Getenv("READWISE_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("READWISE_ACCESS_TOKEN is not set")
	}

	readerClient, err := reader.NewClient(token)
	if err != nil {
		log.Fatalf("failed to create reader client: %v", err)
	}
	
	// List documents with filters
	ctx := context.Background()
	yesterday := time.Now().AddDate(0, 0, -1)
	filteredResponse, err := readerClient.ListDocuments(ctx, &reader.ListDocumentsOptions{
		UpdatedAfter: &yesterday,
		Location:     reader.LocationNew,
		Category:     reader.CategoryArticle,
	})
	if err != nil {
		log.Fatalf("failed to list filtered documents: %v", err)
	}

	fmt.Printf("Found %d new articles from yesterday\n", filteredResponse.Count)
}
```


## License

[MIT](/LICENSE)
