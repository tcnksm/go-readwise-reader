package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	reader "github.com/tcnksm/go-readwise-reader"
)

func toolMove(client reader.Client) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool(
			"move",
			mcp.WithDescription("Move a document to a different location in Readwise Reader"),
			mcp.WithToolAnnotation(
				mcp.ToolAnnotation{
					Title:        "Move a document to a different location",
					ReadOnlyHint: ToBoolPtr(false),
				},
			),
			mcp.WithString(
				"id",
				mcp.Description("The ID of the document to move"),
				mcp.Required(),
			),
			mcp.WithString(
				"location",
				mcp.Description("The target location: new, inbox, later, archive, or feed"),
				mcp.Required(),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Extract parameters
			id, err := req.RequireString("id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			location, err := req.RequireString("location")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			// Validate location
			var loc reader.Location
			switch location {
			case "new", "inbox":
				loc = reader.LocationNew
			case "later":
				loc = reader.LocationLater
			case "archive":
				loc = reader.LocationArchive
			case "feed":
				loc = reader.LocationFeed
			default:
				return mcp.NewToolResultError("invalid location: must be one of new, inbox, later, archive, or feed"), nil
			}

			// Update document location
			updateReq := &reader.UpdateDocumentRequest{
				Location: loc,
			}

			resp, err := client.UpdateDocument(ctx, id, updateReq)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to move document: %v", err)), nil
			}

			// Return JSON response
			jsonData, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to format response: %v", err)), nil
			}

			return mcp.NewToolResultText(string(jsonData)), nil
		}
}