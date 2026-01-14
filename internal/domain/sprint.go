package domain

// Sprint is loaded from kamaji.yaml.
type Sprint struct {
	Name       string   `yaml:"name"`
	BaseBranch string   `yaml:"base_branch"`
	Rules      []string `yaml:"rules"`
	Tickets    []Ticket `yaml:"tickets"`
}

type Ticket struct {
	Name        string `yaml:"name"`
	Branch      string `yaml:"branch"`
	Description string `yaml:"description"`
	Tasks       []Task `yaml:"tasks"`
}

type Task struct {
	Description string   `yaml:"description"`
	Steps       []string `yaml:"steps"`
	Verify      string   `yaml:"verify"`
}
