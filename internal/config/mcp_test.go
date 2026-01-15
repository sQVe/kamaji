package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteMCPConfig_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path, err := WriteMCPConfig(dir, 8080)
	if err != nil {
		t.Fatalf("WriteMCPConfig() error = %v", err)
	}

	expectedPath := filepath.Join(dir, ".mcp.json")
	if path != expectedPath {
		t.Errorf("path = %q, want %q", path, expectedPath)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("File was not created")
	}
}

func TestWriteMCPConfig_CorrectJSON(t *testing.T) {
	dir := t.TempDir()
	path, err := WriteMCPConfig(dir, 9000)
	if err != nil {
		t.Fatalf("WriteMCPConfig() error = %v", err)
	}

	data, err := os.ReadFile(path) //nolint:gosec // path from WriteMCPConfig is safe
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	var cfg mcpConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	kamaji, ok := cfg.MCPServers["kamaji"]
	if !ok {
		t.Fatal("mcpServers.kamaji not found")
	}

	if kamaji.Type != "http" {
		t.Errorf("type = %q, want %q", kamaji.Type, "http")
	}

	expectedURL := "http://localhost:9000/mcp"
	if kamaji.URL != expectedURL {
		t.Errorf("url = %q, want %q", kamaji.URL, expectedURL)
	}
}

func TestWriteMCPConfig_PortInterpolation(t *testing.T) {
	dir := t.TempDir()
	path, err := WriteMCPConfig(dir, 12345)
	if err != nil {
		t.Fatalf("WriteMCPConfig() error = %v", err)
	}

	data, err := os.ReadFile(path) //nolint:gosec // path from WriteMCPConfig is safe
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	var cfg mcpConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	expectedURL := "http://localhost:12345/mcp"
	if cfg.MCPServers["kamaji"].URL != expectedURL {
		t.Errorf("url = %q, want %q", cfg.MCPServers["kamaji"].URL, expectedURL)
	}
}

func TestWriteMCPConfig_EmptyDirUsesTempDir(t *testing.T) {
	path, err := WriteMCPConfig("", 8080)
	if err != nil {
		t.Fatalf("WriteMCPConfig() error = %v", err)
	}
	defer func() { _ = os.Remove(path) }()

	expectedDir := os.TempDir()
	if filepath.Dir(path) != expectedDir {
		t.Errorf("dir = %q, want %q", filepath.Dir(path), expectedDir)
	}
}

func TestWriteMCPConfig_ValidatesPort(t *testing.T) {
	dir := t.TempDir()

	if _, err := WriteMCPConfig(dir, 0); err == nil {
		t.Error("WriteMCPConfig() error = nil, want error for zero port")
	}

	if _, err := WriteMCPConfig(dir, -1); err == nil {
		t.Error("WriteMCPConfig() error = nil, want error for negative port")
	}

	if _, err := WriteMCPConfig(dir, 65536); err == nil {
		t.Error("WriteMCPConfig() error = nil, want error for port > 65535")
	}
}
