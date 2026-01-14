package config

import (
	"fmt"
	"os"

	"github.com/sqve/kamaji/internal/domain"
	"gopkg.in/yaml.v3"
)

// LoadSprint reads and parses a sprint configuration from the given path.
func LoadSprint(path string) (*domain.Sprint, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- user-provided config path is intentional
	if err != nil {
		return nil, fmt.Errorf("reading sprint file: %w", err)
	}

	var sprint domain.Sprint
	if err := yaml.Unmarshal(data, &sprint); err != nil {
		return nil, fmt.Errorf("parsing sprint YAML: %w", err)
	}

	if err := validateSprint(&sprint); err != nil {
		return nil, err
	}

	return &sprint, nil
}

// validateSprint checks required fields in the sprint configuration.
func validateSprint(s *domain.Sprint) error {
	if s.Name == "" {
		return fmt.Errorf("sprint missing required field: name")
	}

	for i, ticket := range s.Tickets {
		if ticket.Name == "" {
			return fmt.Errorf("ticket[%d] missing required field: name", i)
		}
		for j, task := range ticket.Tasks {
			if task.Description == "" {
				return fmt.Errorf("ticket[%d].task[%d] missing required field: description", i, j)
			}
		}
	}

	return nil
}
