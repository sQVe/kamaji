package output

import (
	"strings"
	"testing"

	"github.com/sqve/kamaji/internal/config"
	"github.com/sqve/kamaji/internal/mcp"
	"github.com/sqve/kamaji/internal/testutil"
)

func TestFormatSignal_PlainMode(t *testing.T) {
	config.SetPlain(true)
	defer config.ResetPlain()

	tests := []struct {
		name     string
		signal   mcp.Signal
		expected string
	}{
		{
			name: "task_complete pass",
			signal: mcp.Signal{
				Tool:    mcp.SignalToolTaskComplete,
				Status:  "pass",
				Summary: "Created login component",
			},
			expected: "[ok] Task completed: Created login component",
		},
		{
			name: "task_complete fail",
			signal: mcp.Signal{
				Tool:    mcp.SignalToolTaskComplete,
				Status:  "fail",
				Summary: "Validation failed",
			},
			expected: "Error: Task failed: Validation failed",
		},
		{
			name: "note_insight",
			signal: mcp.Signal{
				Tool:    mcp.SignalToolNoteInsight,
				Summary: "Found a bug in existing code",
			},
			expected: "-> Insight: Found a bug in existing code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatSignal(tt.signal)
			if got != tt.expected {
				t.Errorf("FormatSignal() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestFormatSignal_StyledMode(t *testing.T) {
	config.SetPlain(false)
	defer config.ResetPlain()

	tests := []struct {
		name        string
		signal      mcp.Signal
		mustContain string
	}{
		{
			name: "task_complete pass contains success text",
			signal: mcp.Signal{
				Tool:    mcp.SignalToolTaskComplete,
				Status:  "pass",
				Summary: "Created component",
			},
			mustContain: "Task completed: Created component",
		},
		{
			name: "task_complete fail contains error text",
			signal: mcp.Signal{
				Tool:    mcp.SignalToolTaskComplete,
				Status:  "fail",
				Summary: "Build failed",
			},
			mustContain: "Task failed: Build failed",
		},
		{
			name: "note_insight contains insight text",
			signal: mcp.Signal{
				Tool:    mcp.SignalToolNoteInsight,
				Summary: "Discovered issue",
			},
			mustContain: "Insight: Discovered issue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatSignal(tt.signal)
			if got == "" {
				t.Errorf("FormatSignal() returned empty string")
			}
			if !strings.Contains(got, tt.mustContain) {
				t.Errorf("FormatSignal() = %q, missing %q", got, tt.mustContain)
			}
		})
	}
}

func TestFormatSignal_UnknownTool(t *testing.T) {
	config.SetPlain(true)
	defer config.ResetPlain()

	signal := mcp.Signal{
		Tool:    "unknown_tool",
		Summary: "Some message",
	}

	got := FormatSignal(signal)
	expected := "-> Some message"
	if got != expected {
		t.Errorf("FormatSignal() = %q, want %q", got, expected)
	}
}

func TestPrintSignal(t *testing.T) {
	config.SetPlain(true)
	defer config.ResetPlain()

	tests := []struct {
		name        string
		signal      mcp.Signal
		mustContain string
	}{
		{
			name: "pass signal",
			signal: mcp.Signal{
				Tool:    mcp.SignalToolTaskComplete,
				Status:  "pass",
				Summary: "Done",
			},
			mustContain: "Task completed: Done",
		},
		{
			name: "fail signal",
			signal: mcp.Signal{
				Tool:    mcp.SignalToolTaskComplete,
				Status:  "fail",
				Summary: "Failed",
			},
			mustContain: "Task failed: Failed",
		},
		{
			name: "insight signal",
			signal: mcp.Signal{
				Tool:    mcp.SignalToolNoteInsight,
				Summary: "Note",
			},
			mustContain: "Insight: Note",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := testutil.CaptureStdout(t, func() {
				PrintSignal(tt.signal)
			})
			testutil.AssertContains(t, output, tt.mustContain)
		})
	}
}
