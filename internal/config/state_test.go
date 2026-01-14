package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sqve/kamaji/internal/domain"
)

func TestLoadState_Valid(t *testing.T) {
	dir := t.TempDir()
	kamajiDir := filepath.Join(dir, ".kamaji")
	if err := os.MkdirAll(kamajiDir, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	content := `current_ticket: 2
current_task: 5
failure_count: 1
`
	path := filepath.Join(kamajiDir, "state.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	state, err := LoadState(dir)
	if err != nil {
		t.Fatalf("LoadState error: %v", err)
	}

	if state.CurrentTicket != 2 {
		t.Errorf("CurrentTicket: got %d, want 2", state.CurrentTicket)
	}
	if state.CurrentTask != 5 {
		t.Errorf("CurrentTask: got %d, want 5", state.CurrentTask)
	}
	if state.FailureCount != 1 {
		t.Errorf("FailureCount: got %d, want 1", state.FailureCount)
	}
}

func TestLoadState_MissingFile_ReturnsZeroValue(t *testing.T) {
	dir := t.TempDir()

	state, err := LoadState(dir)
	if err != nil {
		t.Fatalf("LoadState should not error for missing file: %v", err)
	}

	if state.CurrentTicket != 0 {
		t.Errorf("CurrentTicket: got %d, want 0", state.CurrentTicket)
	}
	if state.CurrentTask != 0 {
		t.Errorf("CurrentTask: got %d, want 0", state.CurrentTask)
	}
	if state.FailureCount != 0 {
		t.Errorf("FailureCount: got %d, want 0", state.FailureCount)
	}
}

func TestLoadState_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	kamajiDir := filepath.Join(dir, ".kamaji")
	if err := os.MkdirAll(kamajiDir, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	content := `current_ticket: 2
  invalid: yaml: syntax
`
	path := filepath.Join(kamajiDir, "state.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, err := LoadState(dir)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestSaveState_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()

	state := &domain.State{
		CurrentTicket: 1,
		CurrentTask:   2,
		FailureCount:  0,
	}

	if err := SaveState(dir, state); err != nil {
		t.Fatalf("SaveState error: %v", err)
	}

	// Verify directory was created
	kamajiDir := filepath.Join(dir, ".kamaji")
	info, err := os.Stat(kamajiDir)
	if err != nil {
		t.Fatalf("stat .kamaji: %v", err)
	}
	if !info.IsDir() {
		t.Error(".kamaji should be a directory")
	}

	// Verify file contents
	path := filepath.Join(kamajiDir, "state.yaml")
	data, err := os.ReadFile(path) //nolint:gosec // test file path from t.TempDir
	if err != nil {
		t.Fatalf("read state.yaml: %v", err)
	}
	if !strings.Contains(string(data), "current_ticket: 1") {
		t.Errorf("state.yaml should contain 'current_ticket: 1', got: %s", data)
	}
}

func TestSaveState_Roundtrip(t *testing.T) {
	dir := t.TempDir()

	original := &domain.State{
		CurrentTicket: 3,
		CurrentTask:   7,
		FailureCount:  2,
	}

	if err := SaveState(dir, original); err != nil {
		t.Fatalf("SaveState error: %v", err)
	}

	loaded, err := LoadState(dir)
	if err != nil {
		t.Fatalf("LoadState error: %v", err)
	}

	if loaded.CurrentTicket != original.CurrentTicket {
		t.Errorf("CurrentTicket: got %d, want %d", loaded.CurrentTicket, original.CurrentTicket)
	}
	if loaded.CurrentTask != original.CurrentTask {
		t.Errorf("CurrentTask: got %d, want %d", loaded.CurrentTask, original.CurrentTask)
	}
	if loaded.FailureCount != original.FailureCount {
		t.Errorf("FailureCount: got %d, want %d", loaded.FailureCount, original.FailureCount)
	}
}
