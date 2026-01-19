package orchestrator_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sqve/kamaji/internal/config"
	"github.com/sqve/kamaji/internal/domain"
	"github.com/sqve/kamaji/internal/orchestrator"
	"github.com/sqve/kamaji/internal/statemachine"
	"github.com/sqve/kamaji/internal/testutil"
)

func TestOnPass_CommitsAndAdvances(t *testing.T) {
	dir := t.TempDir()
	testutil.InitGitRepo(t, dir)

	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{Name: "TICKET-1", Tasks: []domain.Task{{Description: "task 1"}, {Description: "task 2"}}},
		},
	}
	state := &domain.State{CurrentTicket: 0, CurrentTask: 0, FailureCount: 0}

	writeFile(t, dir, "test.txt", "content")

	h := orchestrator.NewHandler(dir, state, sprint)
	err := h.OnPass("TICKET-1", "task 1", "Implement feature")
	if err != nil {
		t.Fatalf("OnPass failed: %v", err)
	}

	if state.CurrentTask != 1 {
		t.Errorf("expected task 1, got %d", state.CurrentTask)
	}
	if state.FailureCount != 0 {
		t.Errorf("expected failure count 0, got %d", state.FailureCount)
	}

	assertCommitExists(t, dir, "Implement feature")

	history, err := config.LoadTicketHistory(dir, "TICKET-1")
	if err != nil {
		t.Fatalf("LoadTicketHistory failed: %v", err)
	}
	if len(history.Completed) != 1 {
		t.Fatalf("expected 1 completed, got %d", len(history.Completed))
	}
	if history.Completed[0].Summary != "Implement feature" {
		t.Errorf("expected summary 'Implement feature', got %q", history.Completed[0].Summary)
	}

	savedState, err := config.LoadState(dir)
	if err != nil {
		t.Fatalf("LoadState failed: %v", err)
	}
	if savedState.CurrentTask != 1 {
		t.Errorf("saved state task mismatch: got %d", savedState.CurrentTask)
	}
}

func TestOnPass_NoChangesStillAdvances(t *testing.T) {
	dir := t.TempDir()
	testutil.InitGitRepo(t, dir)

	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{Name: "TICKET-1", Tasks: []domain.Task{{Description: "task 1"}, {Description: "task 2"}}},
		},
	}
	state := &domain.State{CurrentTicket: 0, CurrentTask: 0, FailureCount: 0}

	h := orchestrator.NewHandler(dir, state, sprint)
	err := h.OnPass("TICKET-1", "task 1", "Verification passed")
	if err != nil {
		t.Fatalf("OnPass failed: %v", err)
	}

	if state.CurrentTask != 1 {
		t.Errorf("expected task 1, got %d", state.CurrentTask)
	}

	history, err := config.LoadTicketHistory(dir, "TICKET-1")
	if err != nil {
		t.Fatalf("LoadTicketHistory failed: %v", err)
	}
	if len(history.Completed) != 1 {
		t.Fatalf("expected 1 completed, got %d", len(history.Completed))
	}
}

func TestOnFail_ResetsAndRecords(t *testing.T) {
	dir := t.TempDir()
	testutil.InitGitRepo(t, dir)

	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{Name: "TICKET-1", Tasks: []domain.Task{{Description: "task 1"}}},
		},
	}
	state := &domain.State{CurrentTicket: 0, CurrentTask: 0, FailureCount: 0}

	testFile := filepath.Join(dir, "uncommitted.txt")
	if err := os.WriteFile(testFile, []byte("changes"), 0o600); err != nil {
		t.Fatal(err)
	}

	h := orchestrator.NewHandler(dir, state, sprint)
	err := h.OnFail("TICKET-1", "task 1", "Tests failed")
	if err != nil {
		t.Fatalf("OnFail failed: %v", err)
	}

	if state.FailureCount != 1 {
		t.Errorf("expected failure count 1, got %d", state.FailureCount)
	}
	if state.CurrentTask != 0 {
		t.Errorf("expected to stay on task 0, got %d", state.CurrentTask)
	}

	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("expected uncommitted file to be removed after reset")
	}

	history, err := config.LoadTicketHistory(dir, "TICKET-1")
	if err != nil {
		t.Fatalf("LoadTicketHistory failed: %v", err)
	}
	if len(history.FailedAttempts) != 1 {
		t.Fatalf("expected 1 failed attempt, got %d", len(history.FailedAttempts))
	}
	if history.FailedAttempts[0].Summary != "Tests failed" {
		t.Errorf("expected summary 'Tests failed', got %q", history.FailedAttempts[0].Summary)
	}

	savedState, err := config.LoadState(dir)
	if err != nil {
		t.Fatalf("LoadState failed: %v", err)
	}
	if savedState.FailureCount != 1 {
		t.Errorf("saved failure count mismatch: got %d", savedState.FailureCount)
	}
}

func TestOnFail_IncrementsToStuck(t *testing.T) {
	dir := t.TempDir()
	testutil.InitGitRepo(t, dir)

	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{Name: "TICKET-1", Tasks: []domain.Task{{Description: "task 1"}}},
		},
	}
	state := &domain.State{CurrentTicket: 0, CurrentTask: 0, FailureCount: 2}

	h := orchestrator.NewHandler(dir, state, sprint)

	if h.IsStuck() {
		t.Error("should not be stuck with FailureCount=2")
	}

	err := h.OnFail("TICKET-1", "task 1", "third failure")
	if err != nil {
		t.Fatalf("OnFail failed: %v", err)
	}

	if state.FailureCount != 3 {
		t.Errorf("expected failure count 3, got %d", state.FailureCount)
	}
	if !h.IsStuck() {
		t.Error("expected IsStuck() to return true after 3 failures")
	}
}

func TestOnStuck_PreservesState(t *testing.T) {
	dir := t.TempDir()
	testutil.InitGitRepo(t, dir)

	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{Name: "TICKET-1", Tasks: []domain.Task{{Description: "task 1"}}},
		},
	}
	state := &domain.State{CurrentTicket: 0, CurrentTask: 0, FailureCount: statemachine.StuckThreshold}

	h := orchestrator.NewHandler(dir, state, sprint)
	err := h.OnStuck()
	if err != nil {
		t.Fatalf("OnStuck failed: %v", err)
	}

	savedState, err := config.LoadState(dir)
	if err != nil {
		t.Fatalf("LoadState failed: %v", err)
	}
	if savedState.FailureCount != statemachine.StuckThreshold {
		t.Errorf("expected FailureCount=%d preserved, got %d", statemachine.StuckThreshold, savedState.FailureCount)
	}
	if savedState.CurrentTicket != 0 || savedState.CurrentTask != 0 {
		t.Error("state position should be preserved")
	}
}

// TestIsStuck_DelegatesToStateMachine uses minimal setup (empty workDir and sprint)
// because IsStuck only reads state.FailureCount and doesn't use other Handler fields.
func TestIsStuck_DelegatesToStateMachine(t *testing.T) {
	sprint := &domain.Sprint{}

	tests := []struct {
		name         string
		failureCount int
		want         bool
	}{
		{"zero failures", 0, false},
		{"one failure", 1, false},
		{"two failures", 2, false},
		{"at threshold", statemachine.StuckThreshold, true},
		{"above threshold", statemachine.StuckThreshold + 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &domain.State{FailureCount: tt.failureCount}
			h := orchestrator.NewHandler("", state, sprint)
			if got := h.IsStuck(); got != tt.want {
				t.Errorf("IsStuck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
}

func assertCommitExists(t *testing.T, dir, message string) {
	t.Helper()
	cmd := exec.Command("git", "log", "--oneline", "-1")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git log failed: %v", err)
	}
	if !strings.Contains(string(out), message) {
		t.Errorf("expected commit with message containing %q, got: %s", message, out)
	}
}
