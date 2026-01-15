package domain

// TicketHistory persists to .kamaji/history/<ticket>.yaml.
type TicketHistory struct {
	Ticket         string          `yaml:"ticket"`
	Completed      []CompletedTask `yaml:"completed"`
	FailedAttempts []FailedAttempt `yaml:"failed_attempts"`
	Insights       []string        `yaml:"insights"`
}

type CompletedTask struct {
	Task    string `yaml:"task"`
	Summary string `yaml:"summary"`
}

type FailedAttempt struct {
	Task    string `yaml:"task"`
	Summary string `yaml:"summary"`
}

// HistorySummary provides aggregate statistics for ticket history.
type HistorySummary struct {
	TotalCompleted int
	TotalFailed    int
	TotalInsights  int
	TicketCount    int
}
