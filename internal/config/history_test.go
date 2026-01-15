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

func TestRecordCompleted_EmptyHistory(t *testing.T) {
	dir := t.TempDir()

	if err := RecordCompleted(dir, "new-ticket", "Task 1", "Completed task one"); err != nil {
		t.Fatalf("RecordCompleted error: %v", err)
	}

	history, err := LoadTicketHistory(dir, "new-ticket")
	if err != nil {
		t.Fatalf("LoadTicketHistory error: %v", err)
	}

	if len(history.Completed) != 1 {
		t.Errorf("Completed length: got %d, want 1", len(history.Completed))
	}
	if history.Completed[0].Task != "Task 1" {
		t.Errorf("Task: got %q, want %q", history.Completed[0].Task, "Task 1")
	}
	if history.Completed[0].Summary != "Completed task one" {
		t.Errorf("Summary: got %q, want %q", history.Completed[0].Summary, "Completed task one")
	}
}

func TestRecordCompleted_ExistingHistory(t *testing.T) {
	dir := t.TempDir()

	initial := &domain.TicketHistory{
		Ticket: "existing-ticket",
		Completed: []domain.CompletedTask{
			{Task: "Task 1", Summary: "First task"},
		},
	}
	if err := SaveTicketHistory(dir, initial); err != nil {
		t.Fatalf("SaveTicketHistory error: %v", err)
	}

	if err := RecordCompleted(dir, "existing-ticket", "Task 2", "Second task"); err != nil {
		t.Fatalf("RecordCompleted error: %v", err)
	}

	history, err := LoadTicketHistory(dir, "existing-ticket")
	if err != nil {
		t.Fatalf("LoadTicketHistory error: %v", err)
	}

	if len(history.Completed) != 2 {
		t.Errorf("Completed length: got %d, want 2", len(history.Completed))
	}
	if history.Completed[0].Task != "Task 1" {
		t.Errorf("First task: got %q, want %q", history.Completed[0].Task, "Task 1")
	}
	if history.Completed[1].Task != "Task 2" {
		t.Errorf("Second task: got %q, want %q", history.Completed[1].Task, "Task 2")
	}
}

func TestRecordFailed_EmptyHistory(t *testing.T) {
	dir := t.TempDir()

	if err := RecordFailed(dir, "new-ticket", "Failed task", "Something went wrong"); err != nil {
		t.Fatalf("RecordFailed error: %v", err)
	}

	history, err := LoadTicketHistory(dir, "new-ticket")
	if err != nil {
		t.Fatalf("LoadTicketHistory error: %v", err)
	}

	if len(history.FailedAttempts) != 1 {
		t.Errorf("FailedAttempts length: got %d, want 1", len(history.FailedAttempts))
	}
	if history.FailedAttempts[0].Task != "Failed task" {
		t.Errorf("Task: got %q, want %q", history.FailedAttempts[0].Task, "Failed task")
	}
	if history.FailedAttempts[0].Summary != "Something went wrong" {
		t.Errorf("Summary: got %q, want %q", history.FailedAttempts[0].Summary, "Something went wrong")
	}
}

func TestRecordFailed_ExistingHistory(t *testing.T) {
	dir := t.TempDir()

	initial := &domain.TicketHistory{
		Ticket: "existing-ticket",
		FailedAttempts: []domain.FailedAttempt{
			{Task: "First failure", Summary: "Reason 1"},
		},
	}
	if err := SaveTicketHistory(dir, initial); err != nil {
		t.Fatalf("SaveTicketHistory error: %v", err)
	}

	if err := RecordFailed(dir, "existing-ticket", "Second failure", "Reason 2"); err != nil {
		t.Fatalf("RecordFailed error: %v", err)
	}

	history, err := LoadTicketHistory(dir, "existing-ticket")
	if err != nil {
		t.Fatalf("LoadTicketHistory error: %v", err)
	}

	if len(history.FailedAttempts) != 2 {
		t.Errorf("FailedAttempts length: got %d, want 2", len(history.FailedAttempts))
	}
	if history.FailedAttempts[0].Task != "First failure" {
		t.Errorf("First failure: got %q, want %q", history.FailedAttempts[0].Task, "First failure")
	}
	if history.FailedAttempts[1].Task != "Second failure" {
		t.Errorf("Second failure: got %q, want %q", history.FailedAttempts[1].Task, "Second failure")
	}
}

func TestRecordInsight_EmptyHistory(t *testing.T) {
	dir := t.TempDir()

	if err := RecordInsight(dir, "new-ticket", "First insight"); err != nil {
		t.Fatalf("RecordInsight error: %v", err)
	}

	history, err := LoadTicketHistory(dir, "new-ticket")
	if err != nil {
		t.Fatalf("LoadTicketHistory error: %v", err)
	}

	if len(history.Insights) != 1 {
		t.Errorf("Insights length: got %d, want 1", len(history.Insights))
	}
	if history.Insights[0] != "First insight" {
		t.Errorf("Insight: got %q, want %q", history.Insights[0], "First insight")
	}
}

func TestRecordInsight_MultipleInsights(t *testing.T) {
	dir := t.TempDir()

	if err := RecordInsight(dir, "ticket", "Insight 1"); err != nil {
		t.Fatalf("RecordInsight error: %v", err)
	}
	if err := RecordInsight(dir, "ticket", "Insight 2"); err != nil {
		t.Fatalf("RecordInsight error: %v", err)
	}
	if err := RecordInsight(dir, "ticket", "Insight 3"); err != nil {
		t.Fatalf("RecordInsight error: %v", err)
	}

	history, err := LoadTicketHistory(dir, "ticket")
	if err != nil {
		t.Fatalf("LoadTicketHistory error: %v", err)
	}

	if len(history.Insights) != 3 {
		t.Errorf("Insights length: got %d, want 3", len(history.Insights))
	}
	for i, want := range []string{"Insight 1", "Insight 2", "Insight 3"} {
		if history.Insights[i] != want {
			t.Errorf("Insight[%d]: got %q, want %q", i, history.Insights[i], want)
		}
	}
}

func TestRecordInsight_DuplicateAppended(t *testing.T) {
	dir := t.TempDir()

	if err := RecordInsight(dir, "ticket", "Same insight"); err != nil {
		t.Fatalf("RecordInsight error: %v", err)
	}
	if err := RecordInsight(dir, "ticket", "Same insight"); err != nil {
		t.Fatalf("RecordInsight error: %v", err)
	}

	history, err := LoadTicketHistory(dir, "ticket")
	if err != nil {
		t.Fatalf("LoadTicketHistory error: %v", err)
	}

	if len(history.Insights) != 2 {
		t.Errorf("Insights length: got %d, want 2 (duplicates should be appended)", len(history.Insights))
	}
}

func TestListTicketHistories_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()
	historyDir := filepath.Join(dir, ".kamaji", "history")
	if err := os.MkdirAll(historyDir, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	histories, err := ListTicketHistories(dir)
	if err != nil {
		t.Fatalf("ListTicketHistories error: %v", err)
	}

	if len(histories) != 0 {
		t.Errorf("expected empty slice, got %d histories", len(histories))
	}
}

func TestListTicketHistories_DirectoryNotExists(t *testing.T) {
	dir := t.TempDir()

	histories, err := ListTicketHistories(dir)
	if err != nil {
		t.Fatalf("ListTicketHistories should not error for missing directory: %v", err)
	}

	if len(histories) != 0 {
		t.Errorf("expected empty slice, got %d histories", len(histories))
	}
}

func TestListTicketHistories_MultipleFiles(t *testing.T) {
	dir := t.TempDir()

	h1 := &domain.TicketHistory{
		Ticket:    "ticket-1",
		Completed: []domain.CompletedTask{{Task: "T1", Summary: "S1"}},
	}
	h2 := &domain.TicketHistory{
		Ticket:         "ticket-2",
		FailedAttempts: []domain.FailedAttempt{{Task: "F1", Summary: "SF1"}},
	}
	h3 := &domain.TicketHistory{
		Ticket:   "ticket-3",
		Insights: []string{"Insight 1"},
	}

	if err := SaveTicketHistory(dir, h1); err != nil {
		t.Fatalf("SaveTicketHistory error: %v", err)
	}
	if err := SaveTicketHistory(dir, h2); err != nil {
		t.Fatalf("SaveTicketHistory error: %v", err)
	}
	if err := SaveTicketHistory(dir, h3); err != nil {
		t.Fatalf("SaveTicketHistory error: %v", err)
	}

	histories, err := ListTicketHistories(dir)
	if err != nil {
		t.Fatalf("ListTicketHistories error: %v", err)
	}

	if len(histories) != 3 {
		t.Errorf("expected 3 histories, got %d", len(histories))
	}

	tickets := make(map[string]bool)
	for _, h := range histories {
		tickets[h.Ticket] = true
	}

	for _, want := range []string{"ticket-1", "ticket-2", "ticket-3"} {
		if !tickets[want] {
			t.Errorf("expected ticket %q in results", want)
		}
	}
}

func TestListTicketHistories_IgnoresNonYAMLFiles(t *testing.T) {
	dir := t.TempDir()
	historyDir := filepath.Join(dir, ".kamaji", "history")
	if err := os.MkdirAll(historyDir, 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	h := &domain.TicketHistory{Ticket: "valid-ticket"}
	if err := SaveTicketHistory(dir, h); err != nil {
		t.Fatalf("SaveTicketHistory error: %v", err)
	}

	txtPath := filepath.Join(historyDir, "notes.txt")
	if err := os.WriteFile(txtPath, []byte("some notes"), 0o600); err != nil {
		t.Fatalf("write txt file: %v", err)
	}

	jsonPath := filepath.Join(historyDir, "data.json")
	if err := os.WriteFile(jsonPath, []byte("{}"), 0o600); err != nil {
		t.Fatalf("write json file: %v", err)
	}

	histories, err := ListTicketHistories(dir)
	if err != nil {
		t.Fatalf("ListTicketHistories error: %v", err)
	}

	if len(histories) != 1 {
		t.Errorf("expected 1 history (ignoring non-yaml), got %d", len(histories))
	}
	if histories[0].Ticket != "valid-ticket" {
		t.Errorf("expected ticket 'valid-ticket', got %q", histories[0].Ticket)
	}
}

func TestGetHistorySummary_SingleHistory(t *testing.T) {
	history := &domain.TicketHistory{
		Ticket: "test-ticket",
		Completed: []domain.CompletedTask{
			{Task: "T1", Summary: "S1"},
			{Task: "T2", Summary: "S2"},
		},
		FailedAttempts: []domain.FailedAttempt{
			{Task: "F1", Summary: "SF1"},
		},
		Insights: []string{"I1", "I2", "I3"},
	}

	summary := GetHistorySummary(history)

	if summary.TotalCompleted != 2 {
		t.Errorf("TotalCompleted: got %d, want 2", summary.TotalCompleted)
	}
	if summary.TotalFailed != 1 {
		t.Errorf("TotalFailed: got %d, want 1", summary.TotalFailed)
	}
	if summary.TotalInsights != 3 {
		t.Errorf("TotalInsights: got %d, want 3", summary.TotalInsights)
	}
	if summary.TicketCount != 1 {
		t.Errorf("TicketCount: got %d, want 1", summary.TicketCount)
	}
}

func TestGetHistorySummary_EmptyHistory(t *testing.T) {
	history := &domain.TicketHistory{Ticket: "empty-ticket"}

	summary := GetHistorySummary(history)

	if summary.TotalCompleted != 0 {
		t.Errorf("TotalCompleted: got %d, want 0", summary.TotalCompleted)
	}
	if summary.TotalFailed != 0 {
		t.Errorf("TotalFailed: got %d, want 0", summary.TotalFailed)
	}
	if summary.TotalInsights != 0 {
		t.Errorf("TotalInsights: got %d, want 0", summary.TotalInsights)
	}
	if summary.TicketCount != 1 {
		t.Errorf("TicketCount: got %d, want 1", summary.TicketCount)
	}
}

func TestGetHistorySummary_NilHistory(t *testing.T) {
	summary := GetHistorySummary(nil)

	if summary.TotalCompleted != 0 {
		t.Errorf("TotalCompleted: got %d, want 0", summary.TotalCompleted)
	}
	if summary.TotalFailed != 0 {
		t.Errorf("TotalFailed: got %d, want 0", summary.TotalFailed)
	}
	if summary.TotalInsights != 0 {
		t.Errorf("TotalInsights: got %d, want 0", summary.TotalInsights)
	}
	if summary.TicketCount != 0 {
		t.Errorf("TicketCount: got %d, want 0", summary.TicketCount)
	}
}

func TestGetAllHistoriesSummary_MultipleHistories(t *testing.T) {
	histories := []*domain.TicketHistory{
		{
			Ticket:    "ticket-1",
			Completed: []domain.CompletedTask{{Task: "T1", Summary: "S1"}},
			Insights:  []string{"I1"},
		},
		{
			Ticket:         "ticket-2",
			Completed:      []domain.CompletedTask{{Task: "T2", Summary: "S2"}, {Task: "T3", Summary: "S3"}},
			FailedAttempts: []domain.FailedAttempt{{Task: "F1", Summary: "SF1"}},
		},
		{
			Ticket:   "ticket-3",
			Insights: []string{"I2", "I3"},
		},
	}

	summary := GetAllHistoriesSummary(histories)

	if summary.TotalCompleted != 3 {
		t.Errorf("TotalCompleted: got %d, want 3", summary.TotalCompleted)
	}
	if summary.TotalFailed != 1 {
		t.Errorf("TotalFailed: got %d, want 1", summary.TotalFailed)
	}
	if summary.TotalInsights != 3 {
		t.Errorf("TotalInsights: got %d, want 3", summary.TotalInsights)
	}
	if summary.TicketCount != 3 {
		t.Errorf("TicketCount: got %d, want 3", summary.TicketCount)
	}
}

func TestGetAllHistoriesSummary_EmptySlice(t *testing.T) {
	summary := GetAllHistoriesSummary([]*domain.TicketHistory{})

	if summary.TotalCompleted != 0 {
		t.Errorf("TotalCompleted: got %d, want 0", summary.TotalCompleted)
	}
	if summary.TotalFailed != 0 {
		t.Errorf("TotalFailed: got %d, want 0", summary.TotalFailed)
	}
	if summary.TotalInsights != 0 {
		t.Errorf("TotalInsights: got %d, want 0", summary.TotalInsights)
	}
	if summary.TicketCount != 0 {
		t.Errorf("TicketCount: got %d, want 0", summary.TicketCount)
	}
}

func TestGetAllHistoriesSummary_NilSlice(t *testing.T) {
	summary := GetAllHistoriesSummary(nil)

	if summary.TotalCompleted != 0 {
		t.Errorf("TotalCompleted: got %d, want 0", summary.TotalCompleted)
	}
	if summary.TotalFailed != 0 {
		t.Errorf("TotalFailed: got %d, want 0", summary.TotalFailed)
	}
	if summary.TotalInsights != 0 {
		t.Errorf("TotalInsights: got %d, want 0", summary.TotalInsights)
	}
	if summary.TicketCount != 0 {
		t.Errorf("TicketCount: got %d, want 0", summary.TicketCount)
	}
}

func TestGetAllHistoriesSummary_WithNilElements(t *testing.T) {
	histories := []*domain.TicketHistory{
		{Ticket: "ticket-1", Completed: []domain.CompletedTask{{Task: "T1", Summary: "S1"}}},
		nil,
		{Ticket: "ticket-2", Insights: []string{"I1"}},
	}

	summary := GetAllHistoriesSummary(histories)

	if summary.TotalCompleted != 1 {
		t.Errorf("TotalCompleted: got %d, want 1", summary.TotalCompleted)
	}
	if summary.TotalInsights != 1 {
		t.Errorf("TotalInsights: got %d, want 1", summary.TotalInsights)
	}
	if summary.TicketCount != 2 {
		t.Errorf("TicketCount: got %d, want 2 (skipping nil)", summary.TicketCount)
	}
}
