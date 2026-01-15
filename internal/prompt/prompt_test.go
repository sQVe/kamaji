package prompt

import (
	"strings"
	"testing"

	"github.com/sqve/kamaji/internal/domain"
	"github.com/sqve/kamaji/internal/statemachine"
)

func TestBuildPrompt_FullData(t *testing.T) {
	taskInfo := &statemachine.TaskInfo{
		Ticket: &domain.Ticket{
			Name:        "login-form",
			Branch:      "feat/login-form",
			Description: "Create login form with email/password validation",
		},
		Task: &domain.Task{
			Description: "Create LoginForm component with email and password fields",
			Steps:       []string{"Add form validation using Zod", "Handle submit with loading state"},
			Verify:      "Component renders without errors, form validation rejects invalid email",
		},
	}
	rules := []string{"Use TypeScript strict mode.", "Follow existing patterns in src/."}
	history := &domain.TicketHistory{
		Completed:      []domain.CompletedTask{{Task: "Created auth utility", Summary: "Added loginUser to src/utils/auth.ts"}},
		FailedAttempts: []domain.FailedAttempt{{Task: "OAuth integration", Summary: "passport.js conflicts with session middleware"}},
		Insights:       []string{"Codebase uses Zustand for state management"},
	}

	result := BuildPrompt(taskInfo, rules, history)

	// Check task section
	if !strings.Contains(result, `<ticket name="login-form" branch="feat/login-form">`) {
		t.Error("missing ticket tag with name and branch attributes")
	}
	if !strings.Contains(result, "Create login form with email/password validation") {
		t.Error("missing ticket description")
	}
	if !strings.Contains(result, "<current>") {
		t.Error("missing current tag")
	}
	if !strings.Contains(result, "Create LoginForm component with email and password fields") {
		t.Error("missing task description")
	}
	if !strings.Contains(result, "<steps>") {
		t.Error("missing steps tag")
	}
	if !strings.Contains(result, "- Add form validation using Zod") {
		t.Error("missing step 1")
	}
	if !strings.Contains(result, "- Handle submit with loading state") {
		t.Error("missing step 2")
	}
	if !strings.Contains(result, "<verify>") {
		t.Error("missing verify tag")
	}
	if !strings.Contains(result, "Component renders without errors, form validation rejects invalid email") {
		t.Error("missing verify content")
	}

	// Check rules section
	if !strings.Contains(result, "<rules>") {
		t.Error("missing rules tag")
	}
	if !strings.Contains(result, "Use TypeScript strict mode.") {
		t.Error("missing rule 1")
	}
	if !strings.Contains(result, "Follow existing patterns in src/.") {
		t.Error("missing rule 2")
	}

	// Check history section
	if !strings.Contains(result, "<history>") {
		t.Error("missing history tag")
	}
	if !strings.Contains(result, "<completed>") {
		t.Error("missing completed tag")
	}
	if !strings.Contains(result, "- Created auth utility: Added loginUser to src/utils/auth.ts") {
		t.Error("missing completed task")
	}
	if !strings.Contains(result, "<failed_attempts>") {
		t.Error("missing failed_attempts tag")
	}
	if !strings.Contains(result, "- OAuth integration: passport.js conflicts with session middleware") {
		t.Error("missing failed attempt")
	}
	if !strings.Contains(result, "<insights>") {
		t.Error("missing insights tag")
	}
	if !strings.Contains(result, "- Codebase uses Zustand for state management") {
		t.Error("missing insight")
	}

	// Check instructions section
	if !strings.Contains(result, "<instructions>") {
		t.Error("missing instructions tag")
	}
	if !strings.Contains(result, "Complete the task. Call task_complete(pass/fail, summary) when done.") {
		t.Error("missing task_complete instruction")
	}
	if !strings.Contains(result, "Use note_insight() to record discoveries useful for future tasks.") {
		t.Error("missing note_insight instruction")
	}
}

func TestBuildPrompt_EmptyHistory(t *testing.T) {
	taskInfo := &statemachine.TaskInfo{
		Ticket: &domain.Ticket{
			Name:        "test-ticket",
			Branch:      "feat/test",
			Description: "Test description",
		},
		Task: &domain.Task{
			Description: "Test task",
			Steps:       []string{"Step 1"},
			Verify:      "Verify step",
		},
	}
	rules := []string{"Rule 1"}

	// Test with nil history
	result := BuildPrompt(taskInfo, rules, nil)
	if strings.Contains(result, "<history>") {
		t.Error("should not contain history tag when history is nil")
	}

	// Test with empty history struct
	emptyHistory := &domain.TicketHistory{}
	result = BuildPrompt(taskInfo, rules, emptyHistory)
	if strings.Contains(result, "<history>") {
		t.Error("should not contain history tag when history is empty")
	}
}

func TestBuildPrompt_EmptySteps(t *testing.T) {
	taskInfo := &statemachine.TaskInfo{
		Ticket: &domain.Ticket{
			Name:        "test-ticket",
			Branch:      "feat/test",
			Description: "Test description",
		},
		Task: &domain.Task{
			Description: "Test task",
			Steps:       nil,
			Verify:      "Verify step",
		},
	}
	rules := []string{"Rule 1"}

	result := BuildPrompt(taskInfo, rules, nil)

	if strings.Contains(result, "<steps>") {
		t.Error("should not contain steps tag when steps is nil")
	}

	// Test with empty slice
	taskInfo.Task.Steps = []string{}
	result = BuildPrompt(taskInfo, rules, nil)
	if strings.Contains(result, "<steps>") {
		t.Error("should not contain steps tag when steps is empty slice")
	}
}

func TestBuildPrompt_EmptyVerify(t *testing.T) {
	taskInfo := &statemachine.TaskInfo{
		Ticket: &domain.Ticket{
			Name:        "test-ticket",
			Branch:      "feat/test",
			Description: "Test description",
		},
		Task: &domain.Task{
			Description: "Test task",
			Steps:       []string{"Step 1"},
			Verify:      "",
		},
	}
	rules := []string{"Rule 1"}

	result := BuildPrompt(taskInfo, rules, nil)

	if strings.Contains(result, "<verify>") {
		t.Error("should not contain verify tag when verify is empty")
	}
}

func TestBuildPrompt_XMLEscaping(t *testing.T) {
	taskInfo := &statemachine.TaskInfo{
		Ticket: &domain.Ticket{
			Name:        "test<>ticket",
			Branch:      "feat/test&branch",
			Description: "Description with <html> & \"quotes\"",
		},
		Task: &domain.Task{
			Description: "Task with <script> & 'quotes'",
			Steps:       []string{"Step with <tag> & ampersand"},
			Verify:      "Verify <element> & check",
		},
	}
	rules := []string{"Rule with <brackets> & ampersand"}
	history := &domain.TicketHistory{
		Completed:      []domain.CompletedTask{{Task: "Task <1>", Summary: "Summary & details"}},
		FailedAttempts: []domain.FailedAttempt{{Task: "Task <2>", Summary: "Error & info"}},
		Insights:       []string{"Insight with <code> & symbols"},
	}

	result := BuildPrompt(taskInfo, rules, history)

	// Check that unsafe characters are escaped
	if strings.Contains(result, "<html>") {
		t.Error("<html> should be escaped")
	}
	if strings.Contains(result, "<script>") {
		t.Error("<script> should be escaped")
	}
	if !strings.Contains(result, "&lt;html&gt;") {
		t.Error("should contain escaped <html>")
	}
	if !strings.Contains(result, "&amp;") {
		t.Error("should contain escaped ampersand")
	}
	// Check ticket attributes are escaped
	if !strings.Contains(result, `name="test&lt;&gt;ticket"`) {
		t.Error("ticket name should be escaped")
	}
	if !strings.Contains(result, `branch="feat/test&amp;branch"`) {
		t.Error("ticket branch should be escaped")
	}
}

func TestBuildPrompt_NilTaskInfo(t *testing.T) {
	rules := []string{"Rule 1"}

	result := BuildPrompt(nil, rules, nil)

	if result != "" {
		t.Errorf("expected empty string for nil taskInfo, got %q", result)
	}
}

func TestBuildPrompt_EmptyRules(t *testing.T) {
	taskInfo := &statemachine.TaskInfo{
		Ticket: &domain.Ticket{
			Name:        "test-ticket",
			Branch:      "feat/test",
			Description: "Test description",
		},
		Task: &domain.Task{
			Description: "Test task",
			Steps:       []string{"Step 1"},
			Verify:      "Verify step",
		},
	}

	// Test with nil rules
	result := BuildPrompt(taskInfo, nil, nil)
	if strings.Contains(result, "<rules>") {
		t.Error("should not contain rules tag when rules is nil")
	}

	// Test with empty slice
	result = BuildPrompt(taskInfo, []string{}, nil)
	if strings.Contains(result, "<rules>") {
		t.Error("should not contain rules tag when rules is empty slice")
	}
}

func TestBuildPrompt_PartialHistory(t *testing.T) {
	taskInfo := &statemachine.TaskInfo{
		Ticket: &domain.Ticket{
			Name:        "test-ticket",
			Branch:      "feat/test",
			Description: "Test description",
		},
		Task: &domain.Task{
			Description: "Test task",
		},
	}
	rules := []string{"Rule 1"}

	// Only completed tasks
	history := &domain.TicketHistory{
		Completed: []domain.CompletedTask{{Task: "Task 1", Summary: "Done"}},
	}
	result := BuildPrompt(taskInfo, rules, history)

	if !strings.Contains(result, "<history>") {
		t.Error("should contain history tag")
	}
	if !strings.Contains(result, "<completed>") {
		t.Error("should contain completed tag")
	}
	if strings.Contains(result, "<failed_attempts>") {
		t.Error("should not contain failed_attempts tag when empty")
	}
	if strings.Contains(result, "<insights>") {
		t.Error("should not contain insights tag when empty")
	}
}

func TestBuildPrompt_EmptyDescription(t *testing.T) {
	taskInfo := &statemachine.TaskInfo{
		Ticket: &domain.Ticket{
			Name:        "test-ticket",
			Branch:      "feat/test",
			Description: "",
		},
		Task: &domain.Task{
			Description: "Test task",
		},
	}

	result := BuildPrompt(taskInfo, nil, nil)

	// Should not have double newlines in ticket tag
	if strings.Contains(result, ">\n\n</ticket>") {
		t.Error("empty description should not produce double newline")
	}
	if !strings.Contains(result, `<ticket name="test-ticket" branch="feat/test">`) {
		t.Error("missing ticket tag")
	}
}
