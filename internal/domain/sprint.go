package domain

// Sprint defines a coding sprint configuration (kamaji.yaml).
type Sprint struct {
	Name       string   `yaml:"name"`
	BaseBranch string   `yaml:"base_branch"`
	Rules      []string `yaml:"rules"`
	Tickets    []Ticket `yaml:"tickets"`
}

// Ticket defines a unit of work with a branch and tasks.
type Ticket struct {
	Name        string `yaml:"name"`
	Branch      string `yaml:"branch"`
	Description string `yaml:"description"`
	Tasks       []Task `yaml:"tasks"`
}

// Task defines a single task within a ticket.
type Task struct {
	Description string   `yaml:"description"`
	Steps       []string `yaml:"steps"`
	Verify      string   `yaml:"verify"`
}
