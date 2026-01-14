package statemachine

import (
	"testing"

	"github.com/sqve/kamaji/internal/domain"
)

func TestNextTask_ReturnsTaskInfoForValidPosition(t *testing.T) {
	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{
				Name: "ticket-1",
				Tasks: []domain.Task{
					{Description: "task-0"},
					{Description: "task-1"},
				},
			},
			{
				Name: "ticket-2",
				Tasks: []domain.Task{
					{Description: "task-0"},
				},
			},
		},
	}
	state := &domain.State{CurrentTicket: 0, CurrentTask: 0}

	info := NextTask(state, sprint)

	if info == nil {
		t.Fatal("expected TaskInfo, got nil")
	}
	if info.TicketIndex != 0 {
		t.Errorf("TicketIndex = %d, want 0", info.TicketIndex)
	}
	if info.TaskIndex != 0 {
		t.Errorf("TaskIndex = %d, want 0", info.TaskIndex)
	}
	if info.Ticket.Name != "ticket-1" {
		t.Errorf("Ticket.Name = %q, want %q", info.Ticket.Name, "ticket-1")
	}
	if info.Task.Description != "task-0" {
		t.Errorf("Task.Description = %q, want %q", info.Task.Description, "task-0")
	}
}

func TestNextTask_ReturnsTaskInfoAtMidTicket(t *testing.T) {
	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{
				Name: "ticket-1",
				Tasks: []domain.Task{
					{Description: "task-0"},
					{Description: "task-1"},
				},
			},
		},
	}
	state := &domain.State{CurrentTicket: 0, CurrentTask: 1}

	info := NextTask(state, sprint)

	if info == nil {
		t.Fatal("expected TaskInfo, got nil")
	}
	if info.TaskIndex != 1 {
		t.Errorf("TaskIndex = %d, want 1", info.TaskIndex)
	}
	if info.Task.Description != "task-1" {
		t.Errorf("Task.Description = %q, want %q", info.Task.Description, "task-1")
	}
}

func TestNextTask_ReturnsNilWhenPastEnd(t *testing.T) {
	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{
				Name: "ticket-1",
				Tasks: []domain.Task{
					{Description: "task-0"},
				},
			},
		},
	}
	state := &domain.State{CurrentTicket: 2, CurrentTask: 0}

	info := NextTask(state, sprint)

	if info != nil {
		t.Errorf("expected nil, got TaskInfo{TicketIndex: %d, TaskIndex: %d}", info.TicketIndex, info.TaskIndex)
	}
}

func TestNextTask_ReturnsNilForEmptyTicket(t *testing.T) {
	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{
				Name:  "empty-ticket",
				Tasks: []domain.Task{},
			},
		},
	}
	state := &domain.State{CurrentTicket: 0, CurrentTask: 0}

	info := NextTask(state, sprint)

	if info != nil {
		t.Errorf("expected nil for empty ticket, got TaskInfo{TicketIndex: %d, TaskIndex: %d}", info.TicketIndex, info.TaskIndex)
	}
}

func TestNextTask_ReturnsNilForEmptySprint(t *testing.T) {
	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{},
	}
	state := &domain.State{CurrentTicket: 0, CurrentTask: 0}

	info := NextTask(state, sprint)

	if info != nil {
		t.Errorf("expected nil for empty sprint, got TaskInfo{TicketIndex: %d, TaskIndex: %d}", info.TicketIndex, info.TaskIndex)
	}
}

func TestNextTask_HandlesOutOfBoundsTaskIndex(t *testing.T) {
	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{
				Name: "ticket-1",
				Tasks: []domain.Task{
					{Description: "task-0"},
				},
			},
		},
	}
	state := &domain.State{CurrentTicket: 0, CurrentTask: 5}

	info := NextTask(state, sprint)

	if info != nil {
		t.Errorf("expected nil for out-of-bounds task index, got TaskInfo{TicketIndex: %d, TaskIndex: %d}", info.TicketIndex, info.TaskIndex)
	}
}

func TestAdvance_IncrementsTask(t *testing.T) {
	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{
				Name: "ticket-1",
				Tasks: []domain.Task{
					{Description: "task-0"},
					{Description: "task-1"},
					{Description: "task-2"},
				},
			},
		},
	}
	state := &domain.State{CurrentTicket: 0, CurrentTask: 0}

	Advance(state, sprint)

	if state.CurrentTicket != 0 {
		t.Errorf("CurrentTicket = %d, want 0", state.CurrentTicket)
	}
	if state.CurrentTask != 1 {
		t.Errorf("CurrentTask = %d, want 1", state.CurrentTask)
	}
}

func TestAdvance_HandlesTicketBoundary(t *testing.T) {
	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{
				Name: "ticket-1",
				Tasks: []domain.Task{
					{Description: "task-0"},
					{Description: "task-1"},
					{Description: "task-2"},
				},
			},
			{
				Name: "ticket-2",
				Tasks: []domain.Task{
					{Description: "task-0"},
				},
			},
		},
	}
	state := &domain.State{CurrentTicket: 0, CurrentTask: 2}

	Advance(state, sprint)

	if state.CurrentTicket != 1 {
		t.Errorf("CurrentTicket = %d, want 1", state.CurrentTicket)
	}
	if state.CurrentTask != 0 {
		t.Errorf("CurrentTask = %d, want 0", state.CurrentTask)
	}
}

func TestAdvance_StopsAtEnd(t *testing.T) {
	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{
				Name: "ticket-1",
				Tasks: []domain.Task{
					{Description: "task-0"},
				},
			},
		},
	}
	state := &domain.State{CurrentTicket: 0, CurrentTask: 0}

	Advance(state, sprint)

	// Now at position (1, 0) - past end, NextTask returns nil
	if state.CurrentTicket != 1 {
		t.Errorf("CurrentTicket = %d, want 1", state.CurrentTicket)
	}
	if state.CurrentTask != 0 {
		t.Errorf("CurrentTask = %d, want 0", state.CurrentTask)
	}

	// Verify NextTask returns nil
	info := NextTask(state, sprint)
	if info != nil {
		t.Error("expected NextTask to return nil after advancing past end")
	}
}

func TestAdvance_EmptySprint(t *testing.T) {
	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{},
	}
	state := &domain.State{CurrentTicket: 0, CurrentTask: 0}

	// Should not panic
	Advance(state, sprint)

	// State unchanged
	if state.CurrentTicket != 0 {
		t.Errorf("CurrentTicket = %d, want 0", state.CurrentTicket)
	}
	if state.CurrentTask != 0 {
		t.Errorf("CurrentTask = %d, want 0", state.CurrentTask)
	}
}

func TestRecordPass_ResetsFailureCountAndAdvances(t *testing.T) {
	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{
				Name: "ticket-1",
				Tasks: []domain.Task{
					{Description: "task-0"},
					{Description: "task-1"},
				},
			},
		},
	}
	state := &domain.State{CurrentTicket: 0, CurrentTask: 0, FailureCount: 2}

	RecordPass(state, sprint)

	if state.FailureCount != 0 {
		t.Errorf("FailureCount = %d, want 0", state.FailureCount)
	}
	if state.CurrentTicket != 0 {
		t.Errorf("CurrentTicket = %d, want 0", state.CurrentTicket)
	}
	if state.CurrentTask != 1 {
		t.Errorf("CurrentTask = %d, want 1", state.CurrentTask)
	}
}

func TestRecordPass_AdvancesAcrossTicketBoundary(t *testing.T) {
	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{
				Name: "ticket-1",
				Tasks: []domain.Task{
					{Description: "task-0"},
					{Description: "task-1"},
					{Description: "task-2"},
				},
			},
			{
				Name: "ticket-2",
				Tasks: []domain.Task{
					{Description: "task-0"},
				},
			},
		},
	}
	state := &domain.State{CurrentTicket: 0, CurrentTask: 2, FailureCount: 1}

	RecordPass(state, sprint)

	if state.FailureCount != 0 {
		t.Errorf("FailureCount = %d, want 0", state.FailureCount)
	}
	if state.CurrentTicket != 1 {
		t.Errorf("CurrentTicket = %d, want 1", state.CurrentTicket)
	}
	if state.CurrentTask != 0 {
		t.Errorf("CurrentTask = %d, want 0", state.CurrentTask)
	}
}

func TestRecordFail_IncrementsFailureCount(t *testing.T) {
	state := &domain.State{CurrentTicket: 0, CurrentTask: 0, FailureCount: 0}

	RecordFail(state)

	if state.FailureCount != 1 {
		t.Errorf("FailureCount = %d, want 1", state.FailureCount)
	}
}

func TestRecordFail_IncrementsFromExistingCount(t *testing.T) {
	state := &domain.State{CurrentTicket: 0, CurrentTask: 0, FailureCount: 2}

	RecordFail(state)

	if state.FailureCount != 3 {
		t.Errorf("FailureCount = %d, want 3", state.FailureCount)
	}
}

func TestRecordFail_PositionUnchanged(t *testing.T) {
	state := &domain.State{CurrentTicket: 1, CurrentTask: 2, FailureCount: 0}

	RecordFail(state)

	if state.CurrentTicket != 1 {
		t.Errorf("CurrentTicket = %d, want 1", state.CurrentTicket)
	}
	if state.CurrentTask != 2 {
		t.Errorf("CurrentTask = %d, want 2", state.CurrentTask)
	}
}

func TestIsStuck_FalseWhenZeroFailures(t *testing.T) {
	state := &domain.State{FailureCount: 0}

	if IsStuck(state) {
		t.Error("IsStuck = true, want false")
	}
}

func TestIsStuck_FalseWhenBelowThreshold(t *testing.T) {
	state := &domain.State{FailureCount: 2}

	if IsStuck(state) {
		t.Error("IsStuck = true, want false")
	}
}

func TestIsStuck_TrueAtThreshold(t *testing.T) {
	state := &domain.State{FailureCount: 3}

	if !IsStuck(state) {
		t.Error("IsStuck = false, want true")
	}
}

func TestIsStuck_TrueAboveThreshold(t *testing.T) {
	state := &domain.State{FailureCount: 5}

	if !IsStuck(state) {
		t.Error("IsStuck = false, want true")
	}
}
