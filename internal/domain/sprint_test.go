package domain

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSprint_YAMLRoundtrip(t *testing.T) {
	original := Sprint{
		Name:       "Test Sprint",
		BaseBranch: "main",
		Rules:      []string{"Rule 1", "Rule 2"},
		Tickets: []Ticket{
			{
				Name:        "login-form",
				Branch:      "feat/login-form",
				Description: "Create login form",
				Tasks: []Task{
					{
						Description: "Create component",
						Steps:       []string{"Add validation", "Add tests"},
						Verify:      "Component renders",
					},
				},
			},
		},
	}

	data, err := yaml.Marshal(&original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded Sprint
	if err := yaml.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.Name != original.Name {
		t.Errorf("Name: got %q, want %q", decoded.Name, original.Name)
	}
	if decoded.BaseBranch != original.BaseBranch {
		t.Errorf("BaseBranch: got %q, want %q", decoded.BaseBranch, original.BaseBranch)
	}
	if len(decoded.Rules) != len(original.Rules) {
		t.Errorf("Rules length: got %d, want %d", len(decoded.Rules), len(original.Rules))
	}
	if len(decoded.Tickets) != len(original.Tickets) {
		t.Fatalf("Tickets length: got %d, want %d", len(decoded.Tickets), len(original.Tickets))
	}

	ticket := decoded.Tickets[0]
	if ticket.Name != "login-form" {
		t.Errorf("Ticket.Name: got %q, want %q", ticket.Name, "login-form")
	}
	if ticket.Branch != "feat/login-form" {
		t.Errorf("Ticket.Branch: got %q, want %q", ticket.Branch, "feat/login-form")
	}
	if len(ticket.Tasks) != 1 {
		t.Fatalf("Tasks length: got %d, want 1", len(ticket.Tasks))
	}

	task := ticket.Tasks[0]
	if task.Description != "Create component" {
		t.Errorf("Task.Description: got %q, want %q", task.Description, "Create component")
	}
	if len(task.Steps) != 2 {
		t.Errorf("Task.Steps length: got %d, want 2", len(task.Steps))
	}
	if task.Verify != "Component renders" {
		t.Errorf("Task.Verify: got %q, want %q", task.Verify, "Component renders")
	}
}

func TestSprint_ZeroValue(t *testing.T) {
	var s Sprint
	if s.Name != "" {
		t.Errorf("Name zero value: got %q, want empty", s.Name)
	}
	if s.BaseBranch != "" {
		t.Errorf("BaseBranch zero value: got %q, want empty", s.BaseBranch)
	}
	if s.Rules != nil {
		t.Errorf("Rules zero value: got %v, want nil", s.Rules)
	}
	if s.Tickets != nil {
		t.Errorf("Tickets zero value: got %v, want nil", s.Tickets)
	}
}

func TestSprint_YAMLTags(t *testing.T) {
	yamlData := `name: "Sprint Name"
base_branch: develop
rules:
  - "Follow patterns"
tickets:
  - name: ticket-1
    branch: feat/ticket-1
    description: "Description"
    tasks:
      - description: "Task description"
        steps:
          - "Step 1"
        verify: "Verification"
`
	var s Sprint
	if err := yaml.Unmarshal([]byte(yamlData), &s); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if s.Name != "Sprint Name" {
		t.Errorf("Name: got %q, want %q", s.Name, "Sprint Name")
	}
	if s.BaseBranch != "develop" {
		t.Errorf("BaseBranch: got %q, want %q", s.BaseBranch, "develop")
	}
	if len(s.Rules) != 1 || s.Rules[0] != "Follow patterns" {
		t.Errorf("Rules: got %v, want [Follow patterns]", s.Rules)
	}
	if len(s.Tickets) != 1 {
		t.Fatalf("Tickets: got %d tickets, want 1", len(s.Tickets))
	}
	if s.Tickets[0].Name != "ticket-1" {
		t.Errorf("Ticket.Name: got %q, want %q", s.Tickets[0].Name, "ticket-1")
	}
}

func TestTicket_ZeroValue(t *testing.T) {
	var ticket Ticket
	if ticket.Name != "" {
		t.Errorf("Name zero value: got %q, want empty", ticket.Name)
	}
	if ticket.Branch != "" {
		t.Errorf("Branch zero value: got %q, want empty", ticket.Branch)
	}
	if ticket.Description != "" {
		t.Errorf("Description zero value: got %q, want empty", ticket.Description)
	}
	if ticket.Tasks != nil {
		t.Errorf("Tasks zero value: got %v, want nil", ticket.Tasks)
	}
}

func TestTask_ZeroValue(t *testing.T) {
	var task Task
	if task.Description != "" {
		t.Errorf("Description zero value: got %q, want empty", task.Description)
	}
	if task.Steps != nil {
		t.Errorf("Steps zero value: got %v, want nil", task.Steps)
	}
	if task.Verify != "" {
		t.Errorf("Verify zero value: got %q, want empty", task.Verify)
	}
}
