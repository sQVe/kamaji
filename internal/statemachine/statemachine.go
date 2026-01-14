package statemachine

import "github.com/sqve/kamaji/internal/domain"

// StuckThreshold allows transient failures while preventing infinite retry loops.
const StuckThreshold = 3

type TaskInfo struct {
	TicketIndex int
	TaskIndex   int
	Ticket      *domain.Ticket
	Task        *domain.Task
}

// NextTask returns nil when the sprint is complete.
func NextTask(state *domain.State, sprint *domain.Sprint) *TaskInfo {
	if state.CurrentTicket >= len(sprint.Tickets) {
		return nil
	}

	ticket := &sprint.Tickets[state.CurrentTicket]
	if state.CurrentTask >= len(ticket.Tasks) {
		return nil
	}

	return &TaskInfo{
		TicketIndex: state.CurrentTicket,
		TaskIndex:   state.CurrentTask,
		Ticket:      ticket,
		Task:        &ticket.Tasks[state.CurrentTask],
	}
}

// Advance is a no-op when state is past the sprint end.
func Advance(state *domain.State, sprint *domain.Sprint) {
	if state.CurrentTicket >= len(sprint.Tickets) {
		return
	}

	ticket := &sprint.Tickets[state.CurrentTicket]
	state.CurrentTask++

	if state.CurrentTask >= len(ticket.Tasks) {
		state.CurrentTicket++
		state.CurrentTask = 0
	}
}

// RecordPass resets failure count because failures are tracked per-task.
func RecordPass(state *domain.State, sprint *domain.Sprint) {
	state.FailureCount = 0
	Advance(state, sprint)
}

// RecordFail stays on the task to allow retries before giving up.
func RecordFail(state *domain.State) {
	state.FailureCount++
}

func IsStuck(state *domain.State) bool {
	return state.FailureCount >= StuckThreshold
}
