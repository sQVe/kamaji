package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestHandleTaskComplete_Pass(t *testing.T) {
	args := TaskCompleteArgs{
		Status:  "pass",
		Summary: "Done",
	}

	result, err := HandleTaskComplete(context.Background(), mcp.CallToolRequest{}, args)
	if err != nil {
		t.Fatalf("HandleTaskComplete() error = %v", err)
	}

	if result.IsError {
		t.Fatal("HandleTaskComplete() returned error result")
	}

	var content TaskCompleteResult
	if err := json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &content); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if content.Status != "pass" {
		t.Errorf("result.Status = %q, want %q", content.Status, "pass")
	}
	if content.Summary != "Done" {
		t.Errorf("result.Summary = %q, want %q", content.Summary, "Done")
	}
	if !content.Acknowledged {
		t.Error("result.Acknowledged = false, want true")
	}
}

func TestHandleTaskComplete_Fail(t *testing.T) {
	args := TaskCompleteArgs{
		Status:  "fail",
		Summary: "Error occurred",
	}

	result, err := HandleTaskComplete(context.Background(), mcp.CallToolRequest{}, args)
	if err != nil {
		t.Fatalf("HandleTaskComplete() error = %v", err)
	}

	if result.IsError {
		t.Fatal("HandleTaskComplete() returned error result")
	}

	var content TaskCompleteResult
	if err := json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &content); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if content.Status != "fail" {
		t.Errorf("result.Status = %q, want %q", content.Status, "fail")
	}
	if content.Summary != "Error occurred" {
		t.Errorf("result.Summary = %q, want %q", content.Summary, "Error occurred")
	}
	if !content.Acknowledged {
		t.Error("result.Acknowledged = false, want true")
	}
}

func TestHandleTaskComplete_InvalidStatus(t *testing.T) {
	args := TaskCompleteArgs{
		Status:  "invalid",
		Summary: "Some summary",
	}

	result, err := HandleTaskComplete(context.Background(), mcp.CallToolRequest{}, args)
	if err != nil {
		t.Fatalf("HandleTaskComplete() error = %v", err)
	}

	if !result.IsError {
		t.Error("HandleTaskComplete() IsError = false, want true for invalid status")
	}
}

func TestHandleTaskComplete_EmptySummary(t *testing.T) {
	args := TaskCompleteArgs{
		Status:  "pass",
		Summary: "",
	}

	result, err := HandleTaskComplete(context.Background(), mcp.CallToolRequest{}, args)
	if err != nil {
		t.Fatalf("HandleTaskComplete() error = %v", err)
	}

	if !result.IsError {
		t.Error("HandleTaskComplete() IsError = false, want true for empty summary")
	}
}

func TestHandleNoteInsight_Valid(t *testing.T) {
	args := NoteInsightArgs{
		Text: "Found pattern X",
	}

	result, err := HandleNoteInsight(context.Background(), mcp.CallToolRequest{}, args)
	if err != nil {
		t.Fatalf("HandleNoteInsight() error = %v", err)
	}

	if result.IsError {
		t.Fatal("HandleNoteInsight() returned error result")
	}

	var content NoteInsightResult
	if err := json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &content); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if content.Text != "Found pattern X" {
		t.Errorf("result.Text = %q, want %q", content.Text, "Found pattern X")
	}
	if !content.Recorded {
		t.Error("result.Recorded = false, want true")
	}
}

func TestHandleNoteInsight_Empty(t *testing.T) {
	args := NoteInsightArgs{
		Text: "",
	}

	result, err := HandleNoteInsight(context.Background(), mcp.CallToolRequest{}, args)
	if err != nil {
		t.Fatalf("HandleNoteInsight() error = %v", err)
	}

	if !result.IsError {
		t.Error("HandleNoteInsight() IsError = false, want true for empty text")
	}
}
