package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	reader "github.com/tcnksm/go-readwise-reader"
)

type baseCommand struct {
	client reader.Client
}

func (c *baseCommand) initClient(ctx context.Context) error {
	token, err := getToken()
	if err != nil {
		return err
	}

	client, err := reader.NewClient(token)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	c.client = client
	return nil
}

func getToken() (string, error) {
	token := os.Getenv("READWISE_ACCESS_TOKEN")
	if token == "" {
		return "", fmt.Errorf("READWISE_ACCESS_TOKEN not set")
	}
	return token, nil
}

func printJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

func printError(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
}
