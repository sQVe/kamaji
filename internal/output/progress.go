package output

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/sqve/kamaji/internal/config"
	"github.com/sqve/kamaji/internal/domain"
	"github.com/sqve/kamaji/internal/statemachine"
)

var (
	boldStyle = lipgloss.NewStyle().Bold(true)
	dimStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

// TaskProgress formats task progress indicator.
func TaskProgress(info *statemachine.TaskInfo, sprint *domain.Sprint) string {
	if info == nil || sprint == nil || info.Ticket == nil || info.Task == nil {
		return ""
	}

	ticketTotal := len(sprint.Tickets)
	ticketNum := info.TicketIndex + 1
	taskTotal := len(info.Ticket.Tasks)
	taskNum := info.TaskIndex + 1

	if config.IsPlain() {
		return fmt.Sprintf("[Ticket %d/%d] [Task %d/%d] %s",
			ticketNum, ticketTotal, taskNum, taskTotal, info.Task.Description)
	}

	ticketPart := boldStyle.Render(fmt.Sprintf("Ticket %d/%d", ticketNum, ticketTotal))
	taskPart := boldStyle.Render(fmt.Sprintf("Task %d/%d", taskNum, taskTotal))
	return fmt.Sprintf("%s > %s: %s", ticketPart, taskPart, info.Task.Description)
}

// PrintTaskStart outputs task progress when starting a task.
func PrintTaskStart(info *statemachine.TaskInfo, sprint *domain.Sprint) {
	_, _ = fmt.Fprintln(os.Stdout, TaskProgress(info, sprint))
}

// PrintTicketStart outputs when starting a new ticket.
func PrintTicketStart(ticket *domain.Ticket) {
	if ticket == nil {
		return
	}
	var msg string
	if config.IsPlain() {
		msg = fmt.Sprintf("Starting ticket: %s (%s)", ticket.Name, ticket.Branch)
	} else {
		branchPart := dimStyle.Render(fmt.Sprintf("(%s)", ticket.Branch))
		msg = fmt.Sprintf("Starting ticket: %s %s", ticket.Name, branchPart)
	}
	PrintInfo(msg)
}

// PrintBranchCreated outputs branch creation success.
func PrintBranchCreated(branch string) {
	PrintInfo("Created branch: " + branch)
}

// PrintCommitCreated outputs commit success.
func PrintCommitCreated(message string) {
	summary := truncate(message, 50)
	PrintInfo("Committed: " + summary)
}

// PrintResetPerformed outputs reset notification.
func PrintResetPerformed() {
	PrintInfo("Reset to HEAD (discarding changes)")
}

func truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-3]) + "..."
}

// SprintStatus returns a formatted sprint status string.
func SprintStatus(sprint *domain.Sprint, state *domain.State) string {
	if sprint == nil || state == nil {
		return ""
	}
	ticketsDone, tasksDone, totalTasks := calculateProgress(sprint, state)
	totalTickets := len(sprint.Tickets)

	var currentLine string
	if state.CurrentTicket < totalTickets {
		ticket := sprint.Tickets[state.CurrentTicket]
		taskNum := state.CurrentTask + 1
		taskTotal := len(ticket.Tasks)
		currentLine = fmt.Sprintf("Current: Ticket %d (%s) > Task %d/%d",
			state.CurrentTicket+1, ticket.Name, taskNum, taskTotal)
	} else {
		currentLine = "Current: Complete"
	}

	return fmt.Sprintf("Sprint: %q\nProgress: %d/%d tickets, %d/%d tasks\n%s",
		sprint.Name, ticketsDone, totalTickets, tasksDone, totalTasks, currentLine)
}

// PrintSprintStatus outputs the current sprint status.
func PrintSprintStatus(sprint *domain.Sprint, state *domain.State) {
	_, _ = fmt.Fprintln(os.Stdout, SprintStatus(sprint, state))
}

// PrintSprintComplete outputs sprint completion message.
func PrintSprintComplete(sprint *domain.Sprint, state *domain.State) {
	if sprint == nil || state == nil {
		return
	}
	_, _, totalTasks := calculateProgress(sprint, state)
	msg := fmt.Sprintf("Sprint %q complete: %d tickets, %d tasks", sprint.Name, len(sprint.Tickets), totalTasks)
	PrintSuccess(msg)
}

// PrintSprintStuck outputs stuck state message.
func PrintSprintStuck(sprint *domain.Sprint, state *domain.State) {
	if sprint == nil || state == nil {
		return
	}
	if state.CurrentTicket >= len(sprint.Tickets) {
		PrintError("Sprint stuck after completion")
		return
	}

	ticket := sprint.Tickets[state.CurrentTicket]
	taskDesc := "unknown"
	if state.CurrentTask < len(ticket.Tasks) {
		taskDesc = ticket.Tasks[state.CurrentTask].Description
	}

	msg := fmt.Sprintf("Sprint stuck after %d failures on: %s", state.FailureCount, taskDesc)
	PrintError(msg)
}

func calculateProgress(sprint *domain.Sprint, state *domain.State) (ticketsDone, tasksDone, totalTasks int) {
	for i, ticket := range sprint.Tickets {
		totalTasks += len(ticket.Tasks)
		if i < state.CurrentTicket {
			ticketsDone++
			tasksDone += len(ticket.Tasks)
		} else if i == state.CurrentTicket {
			tasksDone += state.CurrentTask
		}
	}
	return ticketsDone, tasksDone, totalTasks
}
