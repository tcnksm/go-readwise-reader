package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	reader "github.com/tcnksm/go-readwise-reader"
)

func toolSave(client reader.Client) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool(
			"save",
			mcp.WithDescription("Save a given URL link to Readwise Reader"),
			mcp.WithToolAnnotation(
				mcp.ToolAnnotation{
					Title:        "Save a URL link to Readwise Reader",
					ReadOnlyHint: ToBoolPtr(false),
				},
			),
			mcp.WithString(
				"url",
				mcp.Description("The URL of the document to save"),
				mcp.Required(),
			),
			mcp.WithString(
				"summary",
				mcp.Description("A brief summary of the given the link"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			url, err := req.RequireString("url")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			resp, err := client.CreateDocument(ctx, url, &reader.CreateDocumentRequest{
				Summary:  req.GetString("summary", ""),
				Location: reader.LocationNew,
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			r, err := json.Marshal(resp)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal issue: %w", err)
			}
			return mcp.NewToolResultText(string(r)), nil
		}
}

// ToBoolPtr converts a bool to a *bool pointer.
func ToBoolPtr(b bool) *bool {
	return &b
}
