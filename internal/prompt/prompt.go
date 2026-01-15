package prompt

import (
	"html"
	"strings"

	"github.com/sqve/kamaji/internal/domain"
	"github.com/sqve/kamaji/internal/statemachine"
)

// BuildPrompt generates XML prompt structure for agent session injection.
func BuildPrompt(taskInfo *statemachine.TaskInfo, rules []string, history *domain.TicketHistory) string {
	if taskInfo == nil {
		return ""
	}

	var b strings.Builder

	b.WriteString("<task>\n")
	writeTicket(&b, taskInfo.Ticket)
	writeCurrent(&b, taskInfo.Task)
	writeSteps(&b, taskInfo.Task.Steps)
	writeVerify(&b, taskInfo.Task.Verify)
	b.WriteString("</task>\n")

	writeRules(&b, rules)
	writeHistory(&b, history)
	writeInstructions(&b)

	return b.String()
}

func writeTicket(b *strings.Builder, ticket *domain.Ticket) {
	b.WriteString(`<ticket name="`)
	b.WriteString(html.EscapeString(ticket.Name))
	b.WriteString(`" branch="`)
	b.WriteString(html.EscapeString(ticket.Branch))
	b.WriteString(`">`)
	if ticket.Description != "" {
		b.WriteString("\n")
		b.WriteString(html.EscapeString(ticket.Description))
	}
	b.WriteString("\n</ticket>\n")
}

func writeCurrent(b *strings.Builder, task *domain.Task) {
	b.WriteString("\n<current>\n")
	b.WriteString(html.EscapeString(task.Description))
	b.WriteString("\n</current>\n")
}

func writeSteps(b *strings.Builder, steps []string) {
	if len(steps) == 0 {
		return
	}
	b.WriteString("\n<steps>\n")
	for _, step := range steps {
		b.WriteString("- ")
		b.WriteString(html.EscapeString(step))
		b.WriteString("\n")
	}
	b.WriteString("</steps>\n")
}

func writeVerify(b *strings.Builder, verify string) {
	if verify == "" {
		return
	}
	b.WriteString("\n<verify>\n")
	b.WriteString(html.EscapeString(verify))
	b.WriteString("\n</verify>\n")
}

func writeRules(b *strings.Builder, rules []string) {
	if len(rules) == 0 {
		return
	}
	b.WriteString("\n<rules>\n")
	for _, rule := range rules {
		b.WriteString(html.EscapeString(rule))
		b.WriteString("\n")
	}
	b.WriteString("</rules>\n")
}

func writeHistory(b *strings.Builder, history *domain.TicketHistory) {
	if history == nil {
		return
	}
	if len(history.Completed) == 0 && len(history.FailedAttempts) == 0 && len(history.Insights) == 0 {
		return
	}

	b.WriteString("\n<history>\n")

	if len(history.Completed) > 0 {
		b.WriteString("<completed>\n")
		for _, c := range history.Completed {
			b.WriteString("- ")
			b.WriteString(html.EscapeString(c.Task))
			b.WriteString(": ")
			b.WriteString(html.EscapeString(c.Summary))
			b.WriteString("\n")
		}
		b.WriteString("</completed>\n")
	}

	if len(history.FailedAttempts) > 0 {
		b.WriteString("\n<failed_attempts>\n")
		for _, f := range history.FailedAttempts {
			b.WriteString("- ")
			b.WriteString(html.EscapeString(f.Task))
			b.WriteString(": ")
			b.WriteString(html.EscapeString(f.Summary))
			b.WriteString("\n")
		}
		b.WriteString("</failed_attempts>\n")
	}

	if len(history.Insights) > 0 {
		b.WriteString("\n<insights>\n")
		for _, insight := range history.Insights {
			b.WriteString("- ")
			b.WriteString(html.EscapeString(insight))
			b.WriteString("\n")
		}
		b.WriteString("</insights>\n")
	}

	b.WriteString("</history>\n")
}

func writeInstructions(b *strings.Builder) {
	b.WriteString("\n<instructions>\n")
	b.WriteString("Complete the task. Call task_complete(pass/fail, summary) when done.\n")
	b.WriteString("Use note_insight() to record discoveries useful for future tasks.\n")
	b.WriteString("</instructions>\n")
}
