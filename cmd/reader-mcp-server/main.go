package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	reader "github.com/tcnksm/go-readwise-reader"
)

func main() {
	token := os.Getenv("READWISE_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("READWISE_ACCESS_TOKEN not set")
	}

	readerClient, err := reader.NewClient(token)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	mcpServer := server.NewMCPServer(
		"readwise-reader",
		"0.1.0",
		server.WithLogging(),
		server.WithToolCapabilities(false), // TODO:What is this?
	)
	mcpServer.AddTool(toolSave(readerClient))
	mcpServer.AddTool(toolList(readerClient))

	log.Println("Starting Stdio server")
	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
