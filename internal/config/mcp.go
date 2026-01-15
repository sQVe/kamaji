package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type mcpConfig struct {
	MCPServers map[string]mcpServerConfig `json:"mcpServers"`
}

type mcpServerConfig struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// WriteMCPConfig creates a .mcp.json file for Claude Code to connect to the kamaji MCP server.
// If dir is empty, os.TempDir() is used. Returns the path to the created file.
func WriteMCPConfig(dir string, port int) (string, error) {
	if port <= 0 || port > 65535 {
		return "", errors.New("port must be between 1 and 65535")
	}
	if dir == "" {
		dir = os.TempDir()
	}

	cfg := mcpConfig{
		MCPServers: map[string]mcpServerConfig{
			"kamaji": {
				Type: "http",
				URL:  fmt.Sprintf("http://localhost:%d/mcp", port),
			},
		},
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal mcp config: %w", err)
	}

	path := filepath.Join(dir, ".mcp.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return "", fmt.Errorf("write mcp config: %w", err)
	}

	return path, nil
}
