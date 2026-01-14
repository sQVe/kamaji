package domain

// TicketHistory holds the history for a single ticket (.kamaji/history/<ticket>.yaml).
type TicketHistory struct {
	Ticket         string          `yaml:"ticket"`
	Completed      []CompletedTask `yaml:"completed"`
	FailedAttempts []FailedAttempt `yaml:"failed_attempts"`
	Insights       []string        `yaml:"insights"`
}

// CompletedTask records a successfully completed task.
type CompletedTask struct {
	Task    string `yaml:"task"`
	Summary string `yaml:"summary"`
}

// FailedAttempt records a failed task attempt.
type FailedAttempt struct {
	Task    string `yaml:"task"`
	Summary string `yaml:"summary"`
}
