package config

import (
	"testing"

	"github.com/sqve/kamaji/internal/domain"
)

func TestValidateSprint_ValidConfig(t *testing.T) {
	sprint := &domain.Sprint{
		Name:       "Test Sprint",
		BaseBranch: "main",
		Tickets: []domain.Ticket{
			{
				Name:        "ticket-1",
				Branch:      "feat/ticket-1",
				Description: "Ticket description",
				Tasks: []domain.Task{
					{
						Description: "Task description",
						Steps:       []string{"step 1"},
						Verify:      "verify it works",
					},
				},
			},
		},
	}

	errs := ValidateSprint(sprint)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %d: %v", len(errs), errs)
	}
}

func TestValidateSprint_MissingSprintName(t *testing.T) {
	sprint := &domain.Sprint{
		Name:    "",
		Tickets: []domain.Ticket{},
	}

	errs := ValidateSprint(sprint)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Field != "name" {
		t.Errorf("Field: got %q, want %q", errs[0].Field, "name")
	}
	if errs[0].Message != "required" {
		t.Errorf("Message: got %q, want %q", errs[0].Message, "required")
	}
}

func TestValidateSprint_MissingTicketName(t *testing.T) {
	sprint := &domain.Sprint{
		Name: "Test Sprint",
		Tickets: []domain.Ticket{
			{
				Name:  "",
				Tasks: []domain.Task{},
			},
		},
	}

	errs := ValidateSprint(sprint)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Field != "tickets[0].name" {
		t.Errorf("Field: got %q, want %q", errs[0].Field, "tickets[0].name")
	}
	if errs[0].Message != "required" {
		t.Errorf("Message: got %q, want %q", errs[0].Message, "required")
	}
}

func TestValidateSprint_MissingTaskDescription(t *testing.T) {
	sprint := &domain.Sprint{
		Name: "Test Sprint",
		Tickets: []domain.Ticket{
			{
				Name: "ticket-1",
				Tasks: []domain.Task{
					{
						Description: "",
					},
				},
			},
		},
	}

	errs := ValidateSprint(sprint)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Field != "tickets[0].tasks[0].description" {
		t.Errorf("Field: got %q, want %q", errs[0].Field, "tickets[0].tasks[0].description")
	}
	if errs[0].Message != "required" {
		t.Errorf("Message: got %q, want %q", errs[0].Message, "required")
	}
}

func TestValidateSprint_WhitespaceOnlyTicketDescription(t *testing.T) {
	sprint := &domain.Sprint{
		Name: "Test Sprint",
		Tickets: []domain.Ticket{
			{
				Name:        "ticket-1",
				Description: "   \n\t  ",
				Tasks:       []domain.Task{},
			},
		},
	}

	errs := ValidateSprint(sprint)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Field != "tickets[0].description" {
		t.Errorf("Field: got %q, want %q", errs[0].Field, "tickets[0].description")
	}
	if errs[0].Message != "cannot be empty" {
		t.Errorf("Message: got %q, want %q", errs[0].Message, "cannot be empty")
	}
}

func TestValidateSprint_WhitespaceOnlyTaskDescription(t *testing.T) {
	sprint := &domain.Sprint{
		Name: "Test Sprint",
		Tickets: []domain.Ticket{
			{
				Name: "ticket-1",
				Tasks: []domain.Task{
					{
						Description: "   \n\t  ",
					},
				},
			},
		},
	}

	errs := ValidateSprint(sprint)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Field != "tickets[0].tasks[0].description" {
		t.Errorf("Field: got %q, want %q", errs[0].Field, "tickets[0].tasks[0].description")
	}
	if errs[0].Message != "cannot be empty" {
		t.Errorf("Message: got %q, want %q", errs[0].Message, "cannot be empty")
	}
}

func TestValidateSprint_MultipleErrors(t *testing.T) {
	sprint := &domain.Sprint{
		Name: "",
		Tickets: []domain.Ticket{
			{
				Name:        "",
				Description: "  ",
				Tasks: []domain.Task{
					{
						Description: "",
					},
					{
						Description: "  \t  ",
					},
				},
			},
			{
				Name: "ticket-2",
				Tasks: []domain.Task{
					{
						Description: "",
					},
				},
			},
		},
	}

	errs := ValidateSprint(sprint)
	if len(errs) != 6 {
		t.Errorf("expected 6 errors, got %d", len(errs))
		for i, err := range errs {
			t.Logf("  error %d: %s - %s", i, err.Field, err.Message)
		}
	}

	expectedFields := map[string]bool{
		"name":                            true,
		"tickets[0].name":                 true,
		"tickets[0].description":          true,
		"tickets[0].tasks[0].description": true,
		"tickets[0].tasks[1].description": true,
		"tickets[1].tasks[0].description": true,
	}

	for _, err := range errs {
		if !expectedFields[err.Field] {
			t.Errorf("unexpected error field: %s", err.Field)
		}
	}
}
