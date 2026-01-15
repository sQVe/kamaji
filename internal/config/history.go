package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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

// RecordCompleted loads the ticket history, appends a completed task, and saves.
// Uses file locking to prevent concurrent write races.
func RecordCompleted(dir, ticketName, taskDesc, summary string) error {
	unlock, err := acquireHistoryLock(dir, ticketName)
	if err != nil {
		return err
	}
	defer unlock()

	history, err := LoadTicketHistory(dir, ticketName)
	if err != nil {
		return err
	}

	history.Completed = append(history.Completed, domain.CompletedTask{
		Task:    taskDesc,
		Summary: summary,
	})

	return SaveTicketHistory(dir, history)
}

// RecordFailed loads the ticket history, appends a failed attempt, and saves.
// Uses file locking to prevent concurrent write races.
func RecordFailed(dir, ticketName, taskDesc, summary string) error {
	unlock, err := acquireHistoryLock(dir, ticketName)
	if err != nil {
		return err
	}
	defer unlock()

	history, err := LoadTicketHistory(dir, ticketName)
	if err != nil {
		return err
	}

	history.FailedAttempts = append(history.FailedAttempts, domain.FailedAttempt{
		Task:    taskDesc,
		Summary: summary,
	})

	return SaveTicketHistory(dir, history)
}

// RecordInsight loads the ticket history, appends an insight, and saves.
// Uses file locking to prevent concurrent write races.
func RecordInsight(dir, ticketName, text string) error {
	unlock, err := acquireHistoryLock(dir, ticketName)
	if err != nil {
		return err
	}
	defer unlock()

	history, err := LoadTicketHistory(dir, ticketName)
	if err != nil {
		return err
	}

	history.Insights = append(history.Insights, text)

	return SaveTicketHistory(dir, history)
}

// acquireHistoryLock creates a lock file for the given ticket to prevent concurrent writes.
// Returns an unlock function that must be called when done.
func acquireHistoryLock(dir, ticketName string) (func(), error) {
	lockDir := filepath.Join(dir, ".kamaji", "history")
	if err := os.MkdirAll(lockDir, 0o750); err != nil {
		return nil, fmt.Errorf("creating history directory: %w", err)
	}

	filename := sanitizeFilename(ticketName)
	lockPath := filepath.Join(lockDir, filename+".lock")

	var lockFile *os.File
	for range 200 {
		f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600) // #nosec G304 -- path derived from user working directory
		if err == nil {
			lockFile = f
			break
		}
		if !errors.Is(err, os.ErrExist) {
			return nil, fmt.Errorf("acquiring history lock: %w", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	if lockFile == nil {
		return nil, fmt.Errorf("timeout acquiring history lock for %s", ticketName)
	}

	return func() {
		_ = lockFile.Close()
		_ = os.Remove(lockPath)
	}, nil
}

// ListTicketHistories returns all ticket histories from .kamaji/history/.
// Returns empty slice if directory doesn't exist.
func ListTicketHistories(dir string) ([]*domain.TicketHistory, error) {
	historyDir := filepath.Join(dir, ".kamaji", "history")

	entries, err := os.ReadDir(historyDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []*domain.TicketHistory{}, nil
		}
		return nil, fmt.Errorf("reading history directory: %w", err)
	}

	var histories []*domain.TicketHistory
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		ticketName := strings.TrimSuffix(entry.Name(), ".yaml")

		history, err := LoadTicketHistory(dir, ticketName)
		if err != nil {
			return nil, fmt.Errorf("loading history %s: %w", entry.Name(), err)
		}

		histories = append(histories, history)
	}

	return histories, nil
}

// GetHistorySummary returns aggregate statistics for a single history.
func GetHistorySummary(history *domain.TicketHistory) domain.HistorySummary {
	if history == nil {
		return domain.HistorySummary{}
	}
	return domain.HistorySummary{
		TotalCompleted: len(history.Completed),
		TotalFailed:    len(history.FailedAttempts),
		TotalInsights:  len(history.Insights),
		TicketCount:    1,
	}
}

// GetAllHistoriesSummary returns aggregate statistics across all histories.
func GetAllHistoriesSummary(histories []*domain.TicketHistory) domain.HistorySummary {
	var summary domain.HistorySummary
	for _, h := range histories {
		if h == nil {
			continue
		}
		summary.TotalCompleted += len(h.Completed)
		summary.TotalFailed += len(h.FailedAttempts)
		summary.TotalInsights += len(h.Insights)
		summary.TicketCount++
	}
	return summary
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
