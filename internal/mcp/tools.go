package mcp

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
)

type TaskCompleteArgs struct {
	Status  string `json:"status"`
	Summary string `json:"summary"`
}

type TaskCompleteResult struct {
	Status       string `json:"status"`
	Summary      string `json:"summary"`
	Acknowledged bool   `json:"acknowledged"`
}

func HandleTaskComplete(ctx context.Context, req mcp.CallToolRequest, args TaskCompleteArgs) (*mcp.CallToolResult, error) {
	if args.Status != "pass" && args.Status != "fail" {
		return mcp.NewToolResultError("status must be pass or fail"), nil
	}

	if args.Summary == "" {
		return mcp.NewToolResultError("summary is required"), nil
	}

	result := TaskCompleteResult{
		Status:       args.Status,
		Summary:      args.Summary,
		Acknowledged: true,
	}

	data, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

type NoteInsightArgs struct {
	Text string `json:"text"`
}

type NoteInsightResult struct {
	Text     string `json:"text"`
	Recorded bool   `json:"recorded"`
}

func HandleNoteInsight(ctx context.Context, req mcp.CallToolRequest, args NoteInsightArgs) (*mcp.CallToolResult, error) {
	if args.Text == "" {
		return mcp.NewToolResultError("text is required"), nil
	}

	result := NoteInsightResult{
		Text:     args.Text,
		Recorded: true,
	}

	data, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}
