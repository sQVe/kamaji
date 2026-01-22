package orchestrator_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sqve/kamaji/internal/domain"
	"github.com/sqve/kamaji/internal/orchestrator"
	"github.com/sqve/kamaji/internal/testutil"
	"gopkg.in/yaml.v3"
)

func TestRun_EmptySprint(t *testing.T) {
	dir := t.TempDir()
	testutil.InitGitRepo(t, dir)

	sprint := &domain.Sprint{
		Name:       "empty",
		BaseBranch: "main",
		Tickets:    []domain.Ticket{},
	}
	sprintPath := writeSprintFile(t, dir, sprint)

	result, err := orchestrator.Run(context.Background(), orchestrator.RunConfig{
		WorkDir:    dir,
		SprintPath: sprintPath,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !result.Success {
		t.Error("expected Success=true for empty sprint")
	}
	if result.TasksRun != 0 {
		t.Errorf("expected TasksRun=0, got %d", result.TasksRun)
	}
}

func TestRun_ConfigValidation_EmptyWorkDir(t *testing.T) {
	_, err := orchestrator.Run(context.Background(), orchestrator.RunConfig{
		SprintPath: "/some/path.yaml",
	})
	if err == nil {
		t.Fatal("expected error for empty WorkDir")
	}
	if !strings.Contains(err.Error(), "WorkDir") {
		t.Errorf("expected error about WorkDir, got: %v", err)
	}
}

func TestRun_ConfigValidation_EmptySprintPath(t *testing.T) {
	_, err := orchestrator.Run(context.Background(), orchestrator.RunConfig{
		WorkDir: "/some/dir",
	})
	if err == nil {
		t.Fatal("expected error for empty SprintPath")
	}
	if !strings.Contains(err.Error(), "SprintPath") {
		t.Errorf("expected error about SprintPath, got: %v", err)
	}
}

func TestRun_InvalidSprintPath(t *testing.T) {
	dir := t.TempDir()

	_, err := orchestrator.Run(context.Background(), orchestrator.RunConfig{
		WorkDir:    dir,
		SprintPath: filepath.Join(dir, "nonexistent.yaml"),
	})
	if err == nil {
		t.Fatal("expected error for missing sprint file")
	}
}

func TestRun_SprintWithEmptyTicket(t *testing.T) {
	dir := t.TempDir()
	testutil.InitGitRepo(t, dir)

	sprint := &domain.Sprint{
		Name:       "test",
		BaseBranch: "main",
		Tickets: []domain.Ticket{
			{Name: "TICKET-1", Branch: "feat/empty", Tasks: []domain.Task{}},
		},
	}
	sprintPath := writeSprintFile(t, dir, sprint)

	result, err := orchestrator.Run(context.Background(), orchestrator.RunConfig{
		WorkDir:    dir,
		SprintPath: sprintPath,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !result.Success {
		t.Error("expected Success=true for sprint with empty ticket")
	}
	if result.TasksRun != 0 {
		t.Errorf("expected TasksRun=0, got %d", result.TasksRun)
	}
}

func TestRun_ContextCancellation_BeforeTaskStart(t *testing.T) {
	// Tests that a pre-cancelled context is detected at the start of the main
	// loop before any task processing begins.
	dir := t.TempDir()
	testutil.InitGitRepo(t, dir)

	sprint := &domain.Sprint{
		Name:       "test",
		BaseBranch: "main",
		Tickets: []domain.Ticket{
			{Name: "TICKET-1", Branch: "feat/test", Tasks: []domain.Task{{Description: "Task 1"}}},
		},
	}
	sprintPath := writeSprintFile(t, dir, sprint)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := orchestrator.Run(ctx, orchestrator.RunConfig{
		WorkDir:    dir,
		SprintPath: sprintPath,
	})
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
	if !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("expected context canceled error, got: %v", err)
	}
}

func writeSprintFile(t *testing.T, dir string, sprint *domain.Sprint) string {
	t.Helper()
	path := filepath.Join(dir, "kamaji.yaml")
	data, err := yaml.Marshal(sprint)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}
