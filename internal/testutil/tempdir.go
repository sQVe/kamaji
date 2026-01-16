package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

func SetupKamajiDir(t *testing.T, baseDir string) string {
	t.Helper()
	kamajiDir := filepath.Join(baseDir, ".kamaji")
	if err := os.MkdirAll(kamajiDir, 0o750); err != nil {
		t.Fatalf("failed to create kamaji dir: %v", err)
	}
	return kamajiDir
}

func SetupHistoryDir(t *testing.T, baseDir string) string {
	t.Helper()
	historyDir := filepath.Join(baseDir, ".kamaji", "history")
	if err := os.MkdirAll(historyDir, 0o750); err != nil {
		t.Fatalf("failed to create history dir: %v", err)
	}
	return historyDir
}

func WriteHistoryFile(t *testing.T, baseDir, ticketName, content string) string {
	t.Helper()
	historyDir := SetupHistoryDir(t, baseDir)
	path := filepath.Join(historyDir, ticketName+".yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write history file: %v", err)
	}
	return path
}
