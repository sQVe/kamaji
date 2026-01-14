package domain

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestState_YAMLRoundtrip(t *testing.T) {
	original := State{
		CurrentTicket: 2,
		CurrentTask:   5,
		FailureCount:  1,
	}

	data, err := yaml.Marshal(&original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded State
	if err := yaml.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.CurrentTicket != original.CurrentTicket {
		t.Errorf("CurrentTicket: got %d, want %d", decoded.CurrentTicket, original.CurrentTicket)
	}
	if decoded.CurrentTask != original.CurrentTask {
		t.Errorf("CurrentTask: got %d, want %d", decoded.CurrentTask, original.CurrentTask)
	}
	if decoded.FailureCount != original.FailureCount {
		t.Errorf("FailureCount: got %d, want %d", decoded.FailureCount, original.FailureCount)
	}
}

func TestState_ZeroValue(t *testing.T) {
	var s State
	if s.CurrentTicket != 0 {
		t.Errorf("CurrentTicket zero value: got %d, want 0", s.CurrentTicket)
	}
	if s.CurrentTask != 0 {
		t.Errorf("CurrentTask zero value: got %d, want 0", s.CurrentTask)
	}
	if s.FailureCount != 0 {
		t.Errorf("FailureCount zero value: got %d, want 0", s.FailureCount)
	}
}

func TestState_YAMLTags(t *testing.T) {
	yamlData := `current_ticket: 3
current_task: 7
failure_count: 2
`
	var s State
	if err := yaml.Unmarshal([]byte(yamlData), &s); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if s.CurrentTicket != 3 {
		t.Errorf("CurrentTicket: got %d, want 3", s.CurrentTicket)
	}
	if s.CurrentTask != 7 {
		t.Errorf("CurrentTask: got %d, want 7", s.CurrentTask)
	}
	if s.FailureCount != 2 {
		t.Errorf("FailureCount: got %d, want 2", s.FailureCount)
	}
}
