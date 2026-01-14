package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sqve/kamaji/internal/domain"
	"gopkg.in/yaml.v3"
)

// LoadState reads the state from .kamaji/state.yaml in the given directory.
// Returns a zero-value State if the file doesn't exist.
func LoadState(dir string) (*domain.State, error) {
	path := filepath.Join(dir, ".kamaji", "state.yaml")

	data, err := os.ReadFile(path) // #nosec G304 -- path derived from user working directory
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &domain.State{}, nil
		}
		return nil, fmt.Errorf("reading state file: %w", err)
	}

	var state domain.State
	if err := yaml.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("parsing state YAML: %w", err)
	}

	return &state, nil
}

// SaveState writes the state to .kamaji/state.yaml in the given directory.
// Creates the .kamaji directory if it doesn't exist.
func SaveState(dir string, state *domain.State) error {
	kamajiDir := filepath.Join(dir, ".kamaji")
	if err := os.MkdirAll(kamajiDir, 0o750); err != nil {
		return fmt.Errorf("creating .kamaji directory: %w", err)
	}

	data, err := yaml.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshaling state: %w", err)
	}

	path := filepath.Join(kamajiDir, "state.yaml")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("writing state file: %w", err)
	}

	return nil
}
