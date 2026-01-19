package mcp

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
)

// Signal tool name constants.
const (
	SignalToolTaskComplete = "task_complete"
	SignalToolNoteInsight  = "note_insight"
)

// Status constants for task_complete results.
const (
	StatusPass = "pass"
	StatusFail = "fail"
)

// Signal represents a tool call event emitted by the MCP server.
type Signal struct {
	Tool    string // SignalToolTaskComplete or SignalToolNoteInsight
	Status  string // "pass" or "fail" (only for task_complete)
	Summary string // task_complete summary or note_insight text
}

type TaskCompleteArgs struct {
	Status  string `json:"status"`
	Summary string `json:"summary"`
}

type TaskCompleteResult struct {
	Status       string `json:"status"`
	Summary      string `json:"summary"`
	Acknowledged bool   `json:"acknowledged"`
}

//nolint:unparam // error return required by mcp-go TypedToolHandler interface
func HandleTaskComplete(_ context.Context, _ mcp.CallToolRequest, args TaskCompleteArgs) (*mcp.CallToolResult, error) {
	if args.Status != StatusPass && args.Status != StatusFail {
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

//nolint:unparam // error return required by mcp-go TypedToolHandler interface
func HandleNoteInsight(_ context.Context, _ mcp.CallToolRequest, args NoteInsightArgs) (*mcp.CallToolResult, error) {
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
