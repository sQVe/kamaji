package prompt

import (
	"strings"
	"testing"

	"github.com/sqve/kamaji/internal/domain"
	"github.com/sqve/kamaji/internal/testutil"
)

func TestAssembleContext_FullSprint(t *testing.T) {
	dir := t.TempDir()
	testutil.WriteHistoryFile(t, dir, "ticket-1", `ticket: ticket-1
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

	testutil.AssertContains(t, result, "<task>")
	testutil.AssertContains(t, result, `<ticket name="ticket-1"`)
	testutil.AssertContains(t, result, "<current>")
	testutil.AssertContains(t, result, "<steps>")
	testutil.AssertContains(t, result, "<verify>")
	testutil.AssertContains(t, result, "<rules>")
	testutil.AssertContains(t, result, "<history>")
	testutil.AssertContains(t, result, "<completed>")
	testutil.AssertContains(t, result, "<insights>")
	testutil.AssertContains(t, result, "<instructions>")
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

	testutil.AssertContains(t, result, "<task>")
	testutil.AssertContains(t, result, `<ticket name="ticket-2"`)
	testutil.AssertNotContains(t, result, "<history>")
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
