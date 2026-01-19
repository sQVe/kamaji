package output

import (
	"fmt"
	"os"

	"github.com/sqve/kamaji/internal/mcp"
)

// FormatSignal returns a styled string for an MCP signal.
func FormatSignal(sig mcp.Signal) string {
	switch sig.Tool {
	case mcp.SignalToolTaskComplete:
		return formatTaskComplete(sig)
	case mcp.SignalToolNoteInsight:
		return formatNoteInsight(sig)
	default:
		return Style(Info, sig.Summary)
	}
}

// PrintSignal outputs an MCP signal with appropriate styling.
func PrintSignal(sig mcp.Signal) {
	_, _ = fmt.Fprintln(os.Stdout, FormatSignal(sig))
}

func formatTaskComplete(sig mcp.Signal) string {
	switch sig.Status {
	case "pass":
		return Style(Success, "Task completed: "+sig.Summary)
	case "fail":
		return Style(Error, "Task failed: "+sig.Summary)
	default:
		return Style(Info, "Task: "+sig.Summary)
	}
}

func formatNoteInsight(sig mcp.Signal) string {
	return Style(Info, "Insight: "+sig.Summary)
}
