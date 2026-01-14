package statemachine

import "github.com/sqve/kamaji/internal/domain"

// StuckThreshold is the number of consecutive failures before a task is considered stuck.
const StuckThreshold = 3

// TaskInfo contains the current task with its context.
type TaskInfo struct {
	TicketIndex int
	TaskIndex   int
	Ticket      *domain.Ticket
	Task        *domain.Task
}

// NextTask returns the current task info, or nil if all tasks are complete.
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

// Advance moves to the next task, handling ticket boundaries.
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

// RecordPass resets failure count and advances to the next task.
func RecordPass(state *domain.State, sprint *domain.Sprint) {
	state.FailureCount = 0
	Advance(state, sprint)
}

// RecordFail increments failure count without advancing.
func RecordFail(state *domain.State) {
	state.FailureCount++
}

// IsStuck returns true when failure count reaches the stuck threshold.
func IsStuck(state *domain.State) bool {
	return state.FailureCount >= StuckThreshold
}
