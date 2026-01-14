package domain

// State holds runtime state for the orchestrator (.kamaji/state.yaml).
type State struct {
	CurrentTicket int `yaml:"current_ticket"`
	CurrentTask   int `yaml:"current_task"`
	FailureCount  int `yaml:"failure_count"`
}
