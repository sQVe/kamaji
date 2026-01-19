package output

import (
	"strings"
	"testing"

	"github.com/sqve/kamaji/internal/config"
)

func TestPrefix_PlainMode(t *testing.T) {
	config.SetPlain(true)
	defer config.ResetPlain()

	tests := []struct {
		msgType  MessageType
		expected string
	}{
		{Success, "[ok] "},
		{Error, "Error: "},
		{Info, "-> "},
		{Warning, "Warning: "},
		{Debug, "[DEBUG] "},
	}

	for _, tt := range tests {
		got := Prefix(tt.msgType)
		if got != tt.expected {
			t.Errorf("Prefix(%d) in plain mode: got %q, want %q", tt.msgType, got, tt.expected)
		}
	}
}

func TestPrefix_StyledMode(t *testing.T) {
	config.SetPlain(false)
	defer config.ResetPlain()

	tests := []struct {
		msgType      MessageType
		containsText string
	}{
		{Success, "[ok]"},
		{Error, "Error:"},
		{Info, "->"},
		{Warning, "Warning:"},
		{Debug, "[DEBUG]"},
	}

	for _, tt := range tests {
		got := Prefix(tt.msgType)
		if !strings.Contains(got, tt.containsText) {
			t.Errorf("Prefix(%d) in styled mode should contain %q, got %q", tt.msgType, tt.containsText, got)
		}
	}
}

func TestStyle_PlainMode(t *testing.T) {
	config.SetPlain(true)
	defer config.ResetPlain()

	tests := []struct {
		msgType  MessageType
		msg      string
		expected string
	}{
		{Success, "done", "[ok] done"},
		{Error, "failed", "Error: failed"},
		{Info, "processing", "-> processing"},
		{Warning, "check this", "Warning: check this"},
		{Debug, "value=42", "[DEBUG] value=42"},
	}

	for _, tt := range tests {
		got := Style(tt.msgType, tt.msg)
		if got != tt.expected {
			t.Errorf("Style(%d, %q) in plain mode: got %q, want %q", tt.msgType, tt.msg, got, tt.expected)
		}
	}
}

func TestStyle_StyledMode(t *testing.T) {
	config.SetPlain(false)
	defer config.ResetPlain()

	tests := []struct {
		msgType      MessageType
		msg          string
		containsText string
	}{
		{Success, "done", "done"},
		{Error, "failed", "failed"},
		{Info, "processing", "processing"},
		{Warning, "check this", "check this"},
		{Debug, "value=42", "value=42"},
	}

	for _, tt := range tests {
		got := Style(tt.msgType, tt.msg)
		if !strings.Contains(got, tt.containsText) {
			t.Errorf("Style(%d, %q) should contain message %q, got %q", tt.msgType, tt.msg, tt.containsText, got)
		}
	}
}

func TestMsgFunctions_PlainMode(t *testing.T) {
	config.SetPlain(true)
	defer config.ResetPlain()

	tests := []struct {
		fn       func(string) string
		msg      string
		expected string
	}{
		{SuccessMsg, "done", "[ok] done"},
		{ErrorMsg, "failed", "Error: failed"},
		{InfoMsg, "processing", "-> processing"},
		{WarningMsg, "check this", "Warning: check this"},
		{DebugMsg, "value=42", "[DEBUG] value=42"},
	}

	for _, tt := range tests {
		got := tt.fn(tt.msg)
		if got != tt.expected {
			t.Errorf("got %q, want %q", got, tt.expected)
		}
	}
}
