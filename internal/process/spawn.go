package process

import (
	"errors"
	"io"
	"os"
	"strconv"

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

type Waiter interface {
	Wait() error
	Kill() error
}

type SpawnResult struct {
	Process    Waiter // Process that can be waited on or killed
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
		"--print", cfg.Prompt, // non-interactive mode with initial prompt
		"--dangerously-skip-permissions", // auto-accept all tool calls without user confirmation
		"--output-format", "stream-json", // structured output for programmatic parsing
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

// SpawnCommand runs an arbitrary command for testing purposes.
// Unlike SpawnClaude, it ignores cfg.Prompt since the command receives context
// via environment variables (KAMAJI_MCP_PORT, KAMAJI_WORK_DIR).
func SpawnCommand(cmd string, cfg SpawnConfig) (*SpawnResult, error) {
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

	env := append(os.Environ(),
		"KAMAJI_MCP_PORT="+strconv.Itoa(cfg.MCPPort),
		"KAMAJI_WORK_DIR="+cfg.WorkDir,
	)

	p := NewProcess(cmd)
	p.Apply(
		WithDir(cfg.WorkDir),
		WithEnv(env),
	)
	if cfg.Stdout != nil {
		p.Apply(WithStdout(cfg.Stdout))
	}
	if cfg.Stderr != nil {
		p.Apply(WithStderr(cfg.Stderr))
	}

	if err := p.Start(); err != nil {
		return nil, err
	}

	return &SpawnResult{Process: p}, nil
}
