package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadSprint_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "kamaji.yaml")

	content := `name: "Test Sprint"
base_branch: main
rules:
  - "Use TypeScript"
tickets:
  - name: login-form
    branch: feat/login-form
    description: "Create login form"
    tasks:
      - description: "Create component"
        steps:
          - "Add validation"
        verify: "Component renders"
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	sprint, err := LoadSprint(path)
	if err != nil {
		t.Fatalf("LoadSprint error: %v", err)
	}

	if sprint.Name != "Test Sprint" {
		t.Errorf("Name: got %q, want %q", sprint.Name, "Test Sprint")
	}
	if sprint.BaseBranch != "main" {
		t.Errorf("BaseBranch: got %q, want %q", sprint.BaseBranch, "main")
	}
	if len(sprint.Rules) != 1 {
		t.Errorf("Rules: got %d, want 1", len(sprint.Rules))
	}
	if len(sprint.Tickets) != 1 {
		t.Fatalf("Tickets: got %d, want 1", len(sprint.Tickets))
	}
	if sprint.Tickets[0].Name != "login-form" {
		t.Errorf("Ticket.Name: got %q, want %q", sprint.Tickets[0].Name, "login-form")
	}
}

func TestLoadSprint_MissingFile(t *testing.T) {
	_, err := LoadSprint("/nonexistent/path/kamaji.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadSprint_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "kamaji.yaml")

	content := `name: "Test Sprint"
  invalid: yaml: syntax
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, err := LoadSprint(path)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoadSprint_ValidationError_MissingName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "kamaji.yaml")

	content := `base_branch: main
tickets: []
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, err := LoadSprint(path)
	if err == nil {
		t.Error("expected error for missing name")
	}
	if !strings.Contains(err.Error(), "name") {
		t.Errorf("error should mention 'name', got: %v", err)
	}
}

func TestLoadSprint_ValidationError_MissingTicketName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "kamaji.yaml")

	content := `name: "Test Sprint"
tickets:
  - branch: feat/something
    tasks: []
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, err := LoadSprint(path)
	if err == nil {
		t.Error("expected error for missing ticket name")
	}
	if !strings.Contains(err.Error(), "ticket[0]") {
		t.Errorf("error should mention 'ticket[0]', got: %v", err)
	}
}

func TestLoadSprint_ValidationError_MissingTaskDescription(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "kamaji.yaml")

	content := `name: "Test Sprint"
tickets:
  - name: ticket-1
    tasks:
      - verify: "Something"
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, err := LoadSprint(path)
	if err == nil {
		t.Error("expected error for missing task description")
	}
	if !strings.Contains(err.Error(), "ticket[0].task[0]") {
		t.Errorf("error should mention 'ticket[0].task[0]', got: %v", err)
	}
}
