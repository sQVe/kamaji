package orchestrator

import (
	"errors"

	"github.com/sqve/kamaji/internal/config"
	"github.com/sqve/kamaji/internal/domain"
	"github.com/sqve/kamaji/internal/git"
	"github.com/sqve/kamaji/internal/output"
	"github.com/sqve/kamaji/internal/statemachine"
)

// Handler orchestrates pass/fail/stuck workflows for task outcomes.
// Handler does not own state or sprint; the caller owns these values and is
// responsible for their lifecycle. Handler expects single-threaded access.
type Handler struct {
	workDir string
	state   *domain.State
	sprint  *domain.Sprint
}

// NewHandler creates a Handler with the required dependencies.
func NewHandler(workDir string, state *domain.State, sprint *domain.Sprint) *Handler {
	return &Handler{
		workDir: workDir,
		state:   state,
		sprint:  sprint,
	}
}

// OnPass commits changes, records completion, advances state, and persists.
// If no files were changed, the commit is skipped but the task still advances.
func (h *Handler) OnPass(ticketName, taskDesc, summary string) error {
	committed := true
	if err := git.CommitChanges(h.workDir, summary); err != nil {
		if errors.Is(err, git.ErrNothingToCommit) {
			committed = false
		} else {
			return err
		}
	}

	if err := config.RecordCompleted(h.workDir, ticketName, taskDesc, summary); err != nil {
		return err
	}

	prevTicket := h.state.CurrentTicket
	prevTask := h.state.CurrentTask
	prevFailures := h.state.FailureCount

	statemachine.RecordPass(h.state, h.sprint)

	if err := config.SaveState(h.workDir, h.state); err != nil {
		h.state.CurrentTicket = prevTicket
		h.state.CurrentTask = prevTask
		h.state.FailureCount = prevFailures
		return err
	}

	if committed {
		output.PrintCommitCreated(summary)
	}
	return nil
}

// OnFail resets changes, records failure, increments failure count, and persists.
func (h *Handler) OnFail(ticketName, taskDesc, summary string) error {
	if err := git.ResetToHead(h.workDir); err != nil {
		return err
	}

	if err := config.RecordFailed(h.workDir, ticketName, taskDesc, summary); err != nil {
		return err
	}

	prevFailures := h.state.FailureCount

	statemachine.RecordFail(h.state)

	if err := config.SaveState(h.workDir, h.state); err != nil {
		h.state.FailureCount = prevFailures
		return err
	}

	output.PrintResetPerformed()
	return nil
}

// OnStuck outputs the stuck message and preserves state for manual intervention.
func (h *Handler) OnStuck() error {
	output.PrintSprintStuck(h.sprint, h.state)
	return config.SaveState(h.workDir, h.state)
}

// IsStuck returns true if the failure count has reached the stuck threshold.
func (h *Handler) IsStuck() bool {
	return statemachine.IsStuck(h.state)
}
