package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sqve/kamaji/internal/domain"
	"gopkg.in/yaml.v3"
)

// LoadTicketHistory reads the ticket history from .kamaji/history/<ticketName>.yaml.
// Returns an empty TicketHistory with the ticket name set if the file doesn't exist.
func LoadTicketHistory(dir, ticketName string) (*domain.TicketHistory, error) {
	filename := sanitizeFilename(ticketName)
	path := filepath.Join(dir, ".kamaji", "history", filename+".yaml")

	data, err := os.ReadFile(path) // #nosec G304 -- path derived from user working directory
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &domain.TicketHistory{Ticket: ticketName}, nil
		}
		return nil, fmt.Errorf("reading ticket history file: %w", err)
	}

	var history domain.TicketHistory
	if err := yaml.Unmarshal(data, &history); err != nil {
		return nil, fmt.Errorf("parsing ticket history YAML: %w", err)
	}

	return &history, nil
}

// SaveTicketHistory writes the ticket history to .kamaji/history/<ticketName>.yaml.
// Creates the .kamaji/history directory if it doesn't exist.
func SaveTicketHistory(dir string, history *domain.TicketHistory) error {
	historyDir := filepath.Join(dir, ".kamaji", "history")
	if err := os.MkdirAll(historyDir, 0o750); err != nil {
		return fmt.Errorf("creating history directory: %w", err)
	}

	data, err := yaml.Marshal(history)
	if err != nil {
		return fmt.Errorf("marshaling ticket history: %w", err)
	}

	filename := sanitizeFilename(history.Ticket)
	path := filepath.Join(historyDir, filename+".yaml")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("writing ticket history file: %w", err)
	}

	return nil
}

// sanitizeFilename replaces characters that are invalid in filenames across platforms.
func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "-",
		"?", "-",
		"\"", "-",
		"<", "-",
		">", "-",
		"|", "-",
	)
	return replacer.Replace(name)
}
