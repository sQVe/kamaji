package domain

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestTicketHistory_YAMLRoundtrip(t *testing.T) {
	original := TicketHistory{
		Ticket: "login-form",
		Completed: []CompletedTask{
			{Task: "Create component", Summary: "Created LoginForm.tsx"},
		},
		FailedAttempts: []FailedAttempt{
			{Task: "Add OAuth", Summary: "Conflicts with middleware"},
		},
		Insights: []string{"Uses Zustand for state"},
	}

	data, err := yaml.Marshal(&original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded TicketHistory
	if err := yaml.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.Ticket != original.Ticket {
		t.Errorf("Ticket: got %q, want %q", decoded.Ticket, original.Ticket)
	}
	if len(decoded.Completed) != 1 {
		t.Fatalf("Completed length: got %d, want 1", len(decoded.Completed))
	}
	if decoded.Completed[0].Task != "Create component" {
		t.Errorf("Completed[0].Task: got %q, want %q", decoded.Completed[0].Task, "Create component")
	}
	if decoded.Completed[0].Summary != "Created LoginForm.tsx" {
		t.Errorf("Completed[0].Summary: got %q, want %q", decoded.Completed[0].Summary, "Created LoginForm.tsx")
	}
	if len(decoded.FailedAttempts) != 1 {
		t.Fatalf("FailedAttempts length: got %d, want 1", len(decoded.FailedAttempts))
	}
	if decoded.FailedAttempts[0].Task != "Add OAuth" {
		t.Errorf("FailedAttempts[0].Task: got %q, want %q", decoded.FailedAttempts[0].Task, "Add OAuth")
	}
	if len(decoded.Insights) != 1 || decoded.Insights[0] != "Uses Zustand for state" {
		t.Errorf("Insights: got %v, want [Uses Zustand for state]", decoded.Insights)
	}
}

func TestTicketHistory_ZeroValue(t *testing.T) {
	var history TicketHistory
	if history.Ticket != "" {
		t.Errorf("Ticket zero value: got %q, want empty", history.Ticket)
	}
	if history.Completed != nil {
		t.Errorf("Completed zero value: got %v, want nil", history.Completed)
	}
	if history.FailedAttempts != nil {
		t.Errorf("FailedAttempts zero value: got %v, want nil", history.FailedAttempts)
	}
	if history.Insights != nil {
		t.Errorf("Insights zero value: got %v, want nil", history.Insights)
	}
}

func TestTicketHistory_YAMLTags(t *testing.T) {
	yamlData := `ticket: feature-x
completed:
  - task: "Task 1"
    summary: "Done task 1"
failed_attempts:
  - task: "Task 2"
    summary: "Failed task 2"
insights:
  - "Important insight"
`
	var history TicketHistory
	if err := yaml.Unmarshal([]byte(yamlData), &history); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if history.Ticket != "feature-x" {
		t.Errorf("Ticket: got %q, want %q", history.Ticket, "feature-x")
	}
	if len(history.Completed) != 1 || history.Completed[0].Task != "Task 1" {
		t.Errorf("Completed: unexpected value %v", history.Completed)
	}
	if len(history.FailedAttempts) != 1 || history.FailedAttempts[0].Task != "Task 2" {
		t.Errorf("FailedAttempts: unexpected value %v", history.FailedAttempts)
	}
	if len(history.Insights) != 1 || history.Insights[0] != "Important insight" {
		t.Errorf("Insights: got %v, want [Important insight]", history.Insights)
	}
}

func TestCompletedTask_ZeroValue(t *testing.T) {
	var ct CompletedTask
	if ct.Task != "" {
		t.Errorf("Task zero value: got %q, want empty", ct.Task)
	}
	if ct.Summary != "" {
		t.Errorf("Summary zero value: got %q, want empty", ct.Summary)
	}
}

func TestFailedAttempt_ZeroValue(t *testing.T) {
	var fa FailedAttempt
	if fa.Task != "" {
		t.Errorf("Task zero value: got %q, want empty", fa.Task)
	}
	if fa.Summary != "" {
		t.Errorf("Summary zero value: got %q, want empty", fa.Summary)
	}
}
