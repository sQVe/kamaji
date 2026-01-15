package process

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestSpawnClaude_ValidatesPrompt(t *testing.T) {
	dir := t.TempDir()
	_, err := SpawnClaude(SpawnConfig{
		Prompt:  "",
		MCPPort: 8080,
		WorkDir: dir,
	})
	if err == nil {
		t.Error("SpawnClaude() error = nil, want error for empty prompt")
	}
}

func TestSpawnClaude_ValidatesMCPPort(t *testing.T) {
	dir := t.TempDir()
	_, err := SpawnClaude(SpawnConfig{
		Prompt:  "test",
		MCPPort: 0,
		WorkDir: dir,
	})
	if err == nil {
		t.Error("SpawnClaude() error = nil, want error for zero MCPPort")
	}

	_, err = SpawnClaude(SpawnConfig{
		Prompt:  "test",
		MCPPort: -1,
		WorkDir: dir,
	})
	if err == nil {
		t.Error("SpawnClaude() error = nil, want error for negative MCPPort")
	}

	_, err = SpawnClaude(SpawnConfig{
		Prompt:  "test",
		MCPPort: 65536,
		WorkDir: dir,
	})
	if err == nil {
		t.Error("SpawnClaude() error = nil, want error for port > 65535")
	}
}

func TestSpawnClaude_ValidatesWorkDir(t *testing.T) {
	_, err := SpawnClaude(SpawnConfig{
		Prompt:  "test",
		MCPPort: 8080,
		WorkDir: "",
	})
	if err == nil {
		t.Error("SpawnClaude() error = nil, want error for empty WorkDir")
	}
}

func TestSpawnClaude_ValidatesWorkDirExists(t *testing.T) {
	_, err := SpawnClaude(SpawnConfig{
		Prompt:  "test",
		MCPPort: 8080,
		WorkDir: "/nonexistent/path/that/does/not/exist",
	})
	if err == nil {
		t.Error("SpawnClaude() error = nil, want error for nonexistent WorkDir")
	}
}

func TestSpawnClaude_CleansUpConfigOnStartFailure(t *testing.T) {
	// Set empty PATH to ensure claude binary not found
	t.Setenv("PATH", "")

	dir := t.TempDir()

	_, err := SpawnClaude(SpawnConfig{
		Prompt:  "test prompt",
		MCPPort: 9999,
		WorkDir: dir,
	})

	// Expect error when claude binary can't be found
	if err == nil {
		t.Fatal("SpawnClaude() error = nil, want error when claude not in PATH")
	}
	if !errors.Is(err, exec.ErrNotFound) {
		t.Errorf("SpawnClaude() error = %v, want exec.ErrNotFound", err)
	}

	// Config should be cleaned up on failure
	configPath := filepath.Join(dir, ".mcp.json")
	if _, statErr := os.Stat(configPath); !os.IsNotExist(statErr) {
		t.Error(".mcp.json should be cleaned up after Start() failure")
	}
}

// TestSpawnClaude_Integration verifies the full spawn flow.
// If claude is available, it starts and we clean up. If not, we verify cleanup behavior.
func TestSpawnClaude_Integration(t *testing.T) {
	dir := t.TempDir()

	result, err := SpawnClaude(SpawnConfig{
		Prompt:  "test integration prompt",
		MCPPort: 12345,
		WorkDir: dir,
	})

	configPath := filepath.Join(dir, ".mcp.json")

	if err == nil && result != nil {
		// Claude is available - config should exist and process should be running
		if _, statErr := os.Stat(configPath); os.IsNotExist(statErr) {
			t.Error(".mcp.json should exist when process starts successfully")
		}
		if result.ConfigPath != configPath {
			t.Errorf("ConfigPath = %q, want %q", result.ConfigPath, configPath)
		}
		_ = result.Process.Kill()
		_ = result.Process.Wait()
		_ = os.Remove(result.ConfigPath)
	} else {
		// Claude not available - config should be cleaned up
		if _, statErr := os.Stat(configPath); !os.IsNotExist(statErr) {
			t.Error(".mcp.json should be cleaned up after Start() failure")
		}
	}
}
