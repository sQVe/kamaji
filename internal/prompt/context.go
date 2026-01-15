package prompt

import (
	"errors"
	"fmt"

	"github.com/sqve/kamaji/internal/config"
	"github.com/sqve/kamaji/internal/domain"
	"github.com/sqve/kamaji/internal/statemachine"
)

// AssembleContext orchestrates context generation for agent session injection.
// Returns empty string with no error when sprint is complete.
func AssembleContext(sprint *domain.Sprint, state *domain.State, kamajiDir string) (string, error) {
	if sprint == nil {
		return "", errors.New("sprint is nil")
	}
	if state == nil {
		return "", errors.New("state is nil")
	}

	taskInfo := statemachine.NextTask(state, sprint)
	if taskInfo == nil {
		return "", nil
	}

	history, err := config.LoadTicketHistory(kamajiDir, taskInfo.Ticket.Name)
	if err != nil {
		return "", fmt.Errorf("load ticket history: %w", err)
	}

	return BuildPrompt(taskInfo, sprint.Rules, history), nil
}
