package process

import (
	"errors"
	"io"
	"os"

	"github.com/sqve/kamaji/internal/config"
)

// SpawnConfig configures SpawnClaude.
type SpawnConfig struct {
	Prompt  string    // Required: context from AssembleContext
	MCPPort int       // Required: port from MCP server
	WorkDir string    // Required: project directory to run in
	Stdout  io.Writer // Optional: defaults to os.Stdout
	Stderr  io.Writer // Optional: defaults to os.Stderr
}

// SpawnResult contains the spawned process and paths for cleanup.
type SpawnResult struct {
	Process    *Process
	ConfigPath string // Path to .mcp.json, caller should remove after process exits
}

// SpawnClaude creates MCP config, launches Claude Code, and returns the result.
// The caller owns the process lifecycle (Wait/Kill) and should remove ConfigPath after the process exits.
func SpawnClaude(cfg SpawnConfig) (*SpawnResult, error) {
	if cfg.Prompt == "" {
		return nil, errors.New("prompt is required")
	}
	if cfg.MCPPort <= 0 || cfg.MCPPort > 65535 {
		return nil, errors.New("MCPPort must be between 1 and 65535")
	}
	if cfg.WorkDir == "" {
		return nil, errors.New("WorkDir is required")
	}
	info, err := os.Stat(cfg.WorkDir)
	if err != nil || !info.IsDir() {
		return nil, errors.New("WorkDir must be an existing directory")
	}

	configPath, err := config.WriteMCPConfig(cfg.WorkDir, cfg.MCPPort)
	if err != nil {
		return nil, err
	}

	p := NewProcess("claude",
		"-p", cfg.Prompt,
		"--dangerously-skip-permissions",
		"--output-format", "stream-json",
	)

	p.Apply(WithDir(cfg.WorkDir))
	if cfg.Stdout != nil {
		p.Apply(WithStdout(cfg.Stdout))
	}
	if cfg.Stderr != nil {
		p.Apply(WithStderr(cfg.Stderr))
	}

	if err := p.Start(); err != nil {
		_ = os.Remove(configPath)
		return nil, err
	}

	return &SpawnResult{Process: p, ConfigPath: configPath}, nil
}
