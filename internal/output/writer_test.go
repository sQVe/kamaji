package output

import (
	"bytes"
	"testing"

	"github.com/sqve/kamaji/internal/config"
)

func TestWriter_SingleLine(t *testing.T) {
	config.SetPlain(true)
	defer config.ResetPlain()

	var buf bytes.Buffer
	w := NewWriter(&buf, Info)

	n, err := w.Write([]byte("hello world\n"))
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if n != 12 {
		t.Errorf("Write returned %d, want 12", n)
	}

	expected := "-> hello world\n"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}

func TestWriter_MultiLine(t *testing.T) {
	config.SetPlain(true)
	defer config.ResetPlain()

	var buf bytes.Buffer
	w := NewWriter(&buf, Error)

	_, err := w.Write([]byte("line one\nline two\nline three\n"))
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}

	expected := "Error: line one\nError: line two\nError: line three\n"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}

func TestWriter_PartialLineBuffering(t *testing.T) {
	config.SetPlain(true)
	defer config.ResetPlain()

	var buf bytes.Buffer
	w := NewWriter(&buf, Success)

	_, err := w.Write([]byte("part"))
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("partial write should not output, got %q", buf.String())
	}

	_, err = w.Write([]byte("ial line\n"))
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}

	expected := "[ok] partial line\n"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}

func TestWriter_Flush(t *testing.T) {
	config.SetPlain(true)
	defer config.ResetPlain()

	var buf bytes.Buffer
	w := NewWriter(&buf, Warning)

	_, err := w.Write([]byte("incomplete"))
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("partial write should not output, got %q", buf.String())
	}

	if err := w.Flush(); err != nil {
		t.Fatalf("Flush error: %v", err)
	}

	expected := "Warning: incomplete"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}

func TestWriter_DifferentTypes(t *testing.T) {
	config.SetPlain(true)
	defer config.ResetPlain()

	tests := []struct {
		msgType  MessageType
		expected string
	}{
		{Success, "[ok] test\n"},
		{Error, "Error: test\n"},
		{Info, "-> test\n"},
		{Warning, "Warning: test\n"},
		{Debug, "[DEBUG] test\n"},
	}

	for _, tt := range tests {
		var buf bytes.Buffer
		w := NewWriter(&buf, tt.msgType)

		_, err := w.Write([]byte("test\n"))
		if err != nil {
			t.Fatalf("Write error for type %d: %v", tt.msgType, err)
		}

		if buf.String() != tt.expected {
			t.Errorf("type %d: got %q, want %q", tt.msgType, buf.String(), tt.expected)
		}
	}
}

func TestNewInfoWriter(t *testing.T) {
	config.SetPlain(true)
	defer config.ResetPlain()

	var buf bytes.Buffer
	w := NewInfoWriter(&buf)

	_, err := w.Write([]byte("info message\n"))
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}

	expected := "-> info message\n"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}

func TestNewErrorWriter(t *testing.T) {
	config.SetPlain(true)
	defer config.ResetPlain()

	var buf bytes.Buffer
	w := NewErrorWriter(&buf)

	_, err := w.Write([]byte("error message\n"))
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}

	expected := "Error: error message\n"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}
