package prompt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sqve/kamaji/internal/domain"
)

func TestAssembleContext_FullSprint(t *testing.T) {
	dir := t.TempDir()
	setupHistory(t, dir, "ticket-1", `ticket: ticket-1
completed:
  - task: "task 0"
    summary: "done"
insights:
  - "useful insight"
`)

	sprint := &domain.Sprint{
		Name:  "test",
		Rules: []string{"rule one"},
		Tickets: []domain.Ticket{{
			Name:        "ticket-1",
			Branch:      "feat/ticket-1",
			Description: "First ticket",
			Tasks: []domain.Task{{
				Description: "Do something",
				Steps:       []string{"step one"},
				Verify:      "check it",
			}},
		}},
	}
	state := &domain.State{CurrentTicket: 0, CurrentTask: 0}

	result, err := AssembleContext(sprint, state, dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertContains(t, result, "<task>")
	assertContains(t, result, `<ticket name="ticket-1"`)
	assertContains(t, result, "<current>")
	assertContains(t, result, "<steps>")
	assertContains(t, result, "<verify>")
	assertContains(t, result, "<rules>")
	assertContains(t, result, "<history>")
	assertContains(t, result, "<completed>")
	assertContains(t, result, "<insights>")
	assertContains(t, result, "<instructions>")
}

func TestAssembleContext_NoHistoryFile(t *testing.T) {
	dir := t.TempDir()

	sprint := &domain.Sprint{
		Name: "test",
		Tickets: []domain.Ticket{{
			Name:        "ticket-2",
			Branch:      "feat/ticket-2",
			Description: "Second ticket",
			Tasks: []domain.Task{{
				Description: "Do something",
			}},
		}},
	}
	state := &domain.State{CurrentTicket: 0, CurrentTask: 0}

	result, err := AssembleContext(sprint, state, dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertContains(t, result, "<task>")
	assertContains(t, result, `<ticket name="ticket-2"`)
	assertNotContains(t, result, "<history>")
}

func TestAssembleContext_SprintComplete(t *testing.T) {
	dir := t.TempDir()

	sprint := &domain.Sprint{
		Name: "test",
		Tickets: []domain.Ticket{{
			Name:   "ticket-1",
			Branch: "feat/ticket-1",
			Tasks: []domain.Task{{
				Description: "task",
			}},
		}},
	}
	state := &domain.State{CurrentTicket: 1, CurrentTask: 0}

	result, err := AssembleContext(sprint, state, dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty string for complete sprint, got: %q", result)
	}
}

func TestAssembleContext_NilSprint(t *testing.T) {
	dir := t.TempDir()
	state := &domain.State{}

	_, err := AssembleContext(nil, state, dir)
	if err == nil {
		t.Error("expected error for nil sprint")
	}
	if !strings.Contains(err.Error(), "sprint is nil") {
		t.Errorf("expected 'sprint is nil' error, got: %v", err)
	}
}

func TestAssembleContext_NilState(t *testing.T) {
	dir := t.TempDir()
	sprint := &domain.Sprint{Name: "test"}

	_, err := AssembleContext(sprint, nil, dir)
	if err == nil {
		t.Error("expected error for nil state")
	}
	if !strings.Contains(err.Error(), "state is nil") {
		t.Errorf("expected 'state is nil' error, got: %v", err)
	}
}

func setupHistory(t *testing.T, dir, ticketName, content string) {
	t.Helper()
	historyDir := filepath.Join(dir, ".kamaji", "history")
	if err := os.MkdirAll(historyDir, 0o750); err != nil {
		t.Fatalf("failed to create history dir: %v", err)
	}
	path := filepath.Join(historyDir, ticketName+".yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write history file: %v", err)
	}
}

func assertContains(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("expected output to contain %q", substr)
	}
}

func assertNotContains(t *testing.T, s, substr string) {
	t.Helper()
	if strings.Contains(s, substr) {
		t.Errorf("expected output to NOT contain %q", substr)
	}
}
