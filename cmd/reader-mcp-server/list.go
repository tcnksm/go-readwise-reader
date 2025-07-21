package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	reader "github.com/tcnksm/go-readwise-reader"
)

func toolList(client reader.Client) (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool(
			"list",
			mcp.WithDescription("List documents from Readwise Reader"),
			mcp.WithToolAnnotation(
				mcp.ToolAnnotation{
					Title:        "List documents from Readwise Reader",
					ReadOnlyHint: ToBoolPtr(true),
				},
			),
			mcp.WithString(
				"location",
				mcp.Description("Location of documents: new, later, archive, or feed"),
				mcp.Required(),
			),
			mcp.WithString(
				"since",
				mcp.Description("Filter documents updated since duration ago (e.g., 10s, 30m, 24h) (default: 12h)"),
			),
			mcp.WithNumber(
				"limit",
				mcp.Description("Maximum number of documents to return (default: 5)"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Extract parameters
			location, err := req.RequireString("location")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			// Validate location
			var loc reader.Location
			switch location {
			case "new":
				loc = reader.LocationNew
			case "later":
				loc = reader.LocationLater
			case "archive":
				loc = reader.LocationArchive
			case "feed":
				loc = reader.LocationFeed
			default:
				return mcp.NewToolResultError("invalid location: must be one of new, later, archive, or feed"), nil
			}

			// Parse optional parameters
			opts := &reader.ListDocumentsOptions{
				Location: loc,
			}

			// Handle since parameter (default: 12h)
			sinceStr := req.GetString("since", "12h")
			duration, err := time.ParseDuration(sinceStr)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("invalid since duration: %v", err)), nil
			}
			updatedAfter := time.Now().Add(-duration)
			opts.UpdatedAfter = &updatedAfter

			// Handle limit parameter
			limit := int(req.GetFloat("limit", 5))

			if limit <= 0 {
				return mcp.NewToolResultError("limit must be greater than 0"), nil
			}

			// Call Readwise API
			resp, err := client.ListDocuments(ctx, opts)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to list documents: %v", err)), nil
			}

			// Limit results if needed
			if len(resp.Results) > limit {
				resp.Results = resp.Results[:limit]
			}

			// Return JSON response
			jsonData, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to format response: %v", err)), nil
			}

			return mcp.NewToolResultText(string(jsonData)), nil
		}
}
