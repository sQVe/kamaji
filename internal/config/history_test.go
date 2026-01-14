package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sqve/kamaji/internal/domain"
)

func TestLoadTicketHistory_Valid(t *testing.T) {
	dir := t.TempDir()
	historyDir := filepath.Join(dir, ".kamaji", "history")
	if err := os.MkdirAll(historyDir, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	content := `ticket: login-form
completed:
  - task: "Create component"
    summary: "Created LoginForm.tsx"
failed_attempts:
  - task: "Add OAuth"
    summary: "Conflicts with middleware"
insights:
  - "Uses Zustand for state"
`
	path := filepath.Join(historyDir, "login-form.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	history, err := LoadTicketHistory(dir, "login-form")
	if err != nil {
		t.Fatalf("LoadTicketHistory error: %v", err)
	}

	if history.Ticket != "login-form" {
		t.Errorf("Ticket: got %q, want %q", history.Ticket, "login-form")
	}
	if len(history.Completed) != 1 {
		t.Errorf("Completed: got %d, want 1", len(history.Completed))
	}
	if len(history.FailedAttempts) != 1 {
		t.Errorf("FailedAttempts: got %d, want 1", len(history.FailedAttempts))
	}
	if len(history.Insights) != 1 {
		t.Errorf("Insights: got %d, want 1", len(history.Insights))
	}
}

func TestLoadTicketHistory_MissingFile_ReturnsZeroValue(t *testing.T) {
	dir := t.TempDir()

	history, err := LoadTicketHistory(dir, "nonexistent")
	if err != nil {
		t.Fatalf("LoadTicketHistory should not error for missing file: %v", err)
	}

	if history.Ticket != "nonexistent" {
		t.Errorf("Ticket: got %q, want %q", history.Ticket, "nonexistent")
	}
	if history.Completed != nil {
		t.Errorf("Completed: got %v, want nil", history.Completed)
	}
	if history.FailedAttempts != nil {
		t.Errorf("FailedAttempts: got %v, want nil", history.FailedAttempts)
	}
	if history.Insights != nil {
		t.Errorf("Insights: got %v, want nil", history.Insights)
	}
}

func TestLoadTicketHistory_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	historyDir := filepath.Join(dir, ".kamaji", "history")
	if err := os.MkdirAll(historyDir, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	content := `ticket: login-form
  invalid: yaml: syntax
`
	path := filepath.Join(historyDir, "login-form.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, err := LoadTicketHistory(dir, "login-form")
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestSaveTicketHistory_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()

	history := &domain.TicketHistory{
		Ticket: "test-ticket",
		Completed: []domain.CompletedTask{
			{Task: "Task 1", Summary: "Done"},
		},
	}

	if err := SaveTicketHistory(dir, history); err != nil {
		t.Fatalf("SaveTicketHistory error: %v", err)
	}

	// Verify directory was created
	historyDir := filepath.Join(dir, ".kamaji", "history")
	info, err := os.Stat(historyDir)
	if err != nil {
		t.Fatalf("stat history dir: %v", err)
	}
	if !info.IsDir() {
		t.Error(".kamaji/history should be a directory")
	}

	// Verify file contents
	path := filepath.Join(historyDir, "test-ticket.yaml")
	data, err := os.ReadFile(path) //nolint:gosec // test file path from t.TempDir
	if err != nil {
		t.Fatalf("read ticket history: %v", err)
	}
	if !strings.Contains(string(data), "ticket: test-ticket") {
		t.Errorf("ticket history should contain 'ticket: test-ticket', got: %s", data)
	}
}

func TestSaveTicketHistory_SanitizesFilename(t *testing.T) {
	dir := t.TempDir()

	history := &domain.TicketHistory{
		Ticket: "feature/login",
	}

	if err := SaveTicketHistory(dir, history); err != nil {
		t.Fatalf("SaveTicketHistory error: %v", err)
	}

	// Check that "/" was replaced with "-"
	path := filepath.Join(dir, ".kamaji", "history", "feature-login.yaml")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file at %s, got error: %v", path, err)
	}

	// Ensure the original path with "/" doesn't exist (intentionally invalid)
	badPath := filepath.Join(dir, ".kamaji", "history", "feature-login.yaml", "..", "feature", "login.yaml") //nolint:gocritic // intentionally testing bad path
	if _, err := os.Stat(badPath); err == nil {
		t.Errorf("file should not exist at %s", badPath)
	}
}

func TestSaveTicketHistory_Roundtrip(t *testing.T) {
	dir := t.TempDir()

	original := &domain.TicketHistory{
		Ticket: "roundtrip-test",
		Completed: []domain.CompletedTask{
			{Task: "Task 1", Summary: "Done 1"},
			{Task: "Task 2", Summary: "Done 2"},
		},
		FailedAttempts: []domain.FailedAttempt{
			{Task: "Failed task", Summary: "Reason"},
		},
		Insights: []string{"Insight 1", "Insight 2"},
	}

	if err := SaveTicketHistory(dir, original); err != nil {
		t.Fatalf("SaveTicketHistory error: %v", err)
	}

	loaded, err := LoadTicketHistory(dir, "roundtrip-test")
	if err != nil {
		t.Fatalf("LoadTicketHistory error: %v", err)
	}

	if loaded.Ticket != original.Ticket {
		t.Errorf("Ticket: got %q, want %q", loaded.Ticket, original.Ticket)
	}
	if len(loaded.Completed) != len(original.Completed) {
		t.Errorf("Completed length: got %d, want %d", len(loaded.Completed), len(original.Completed))
	}
	if len(loaded.FailedAttempts) != len(original.FailedAttempts) {
		t.Errorf("FailedAttempts length: got %d, want %d", len(loaded.FailedAttempts), len(original.FailedAttempts))
	}
	if len(loaded.Insights) != len(original.Insights) {
		t.Errorf("Insights length: got %d, want %d", len(loaded.Insights), len(original.Insights))
	}
}

func TestLoadTicketHistory_WithSlashInName(t *testing.T) {
	dir := t.TempDir()
	historyDir := filepath.Join(dir, ".kamaji", "history")
	if err := os.MkdirAll(historyDir, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create file with sanitized name
	content := `ticket: feature/login
completed: []
`
	path := filepath.Join(historyDir, "feature-login.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	// Load using the original name with slash
	history, err := LoadTicketHistory(dir, "feature/login")
	if err != nil {
		t.Fatalf("LoadTicketHistory error: %v", err)
	}

	if history.Ticket != "feature/login" {
		t.Errorf("Ticket: got %q, want %q", history.Ticket, "feature/login")
	}
}

func TestSaveTicketHistory_SanitizesAllPlatformCharacters(t *testing.T) {
	dir := t.TempDir()

	// Test various problematic characters across platforms
	history := &domain.TicketHistory{
		Ticket: `fix:bug<v1>|"test"`,
	}

	if err := SaveTicketHistory(dir, history); err != nil {
		t.Fatalf("SaveTicketHistory error: %v", err)
	}

	// All problematic characters should be replaced with "-"
	expected := "fix-bug-v1---test-.yaml"
	path := filepath.Join(dir, ".kamaji", "history", expected)
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file at %s, got error: %v", path, err)
	}
}
