//go:build integration

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestScript(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/script",
		Cmds: map[string]func(ts *testscript.TestScript, neg bool, args []string){
			"gitinit": gitInitCmd,
		},
	})
}

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"kamaji":     main,
		"mock-agent": mockAgentMain,
	})
}

func gitInitCmd(ts *testscript.TestScript, neg bool, args []string) {
	if neg {
		ts.Fatalf("gitinit does not support negation")
	}

	dir := ts.Getenv("WORK")
	run := func(gitArgs ...string) {
		cmd := exec.Command("git", gitArgs...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test",
			"GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=Test",
			"GIT_COMMITTER_EMAIL=test@test.com",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			ts.Fatalf("git %v failed: %v\n%s", gitArgs, err, out)
		}
	}

	run("init", "-b", "main")
	run("config", "user.email", "test@test.com")
	run("config", "user.name", "Test")
	run("config", "commit.gpgsign", "false")
	run("config", "core.autocrlf", "false")
}

func mockAgentMain() {
	port, err := strconv.Atoi(os.Getenv("KAMAJI_MCP_PORT"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "KAMAJI_MCP_PORT invalid or not set")
		os.Exit(1)
	}

	script := os.Getenv("KAMAJI_AGENT_SCRIPT")
	if script == "" {
		return
	}

	if err := runMockAgent(port, script); err != nil {
		fmt.Fprintln(os.Stderr, "mock-agent:", err)
		os.Exit(1)
	}
}

func runMockAgent(port int, script string) error {
	c, err := client.NewStreamableHttpClient(fmt.Sprintf("http://localhost:%d/mcp", port))
	if err != nil {
		return err
	}
	defer func() { _ = c.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = c.Initialize(ctx, mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: "2024-11-05",
			ClientInfo:      mcp.Implementation{Name: "mock-agent", Version: "1.0.0"},
		},
	})
	if err != nil {
		return err
	}

	for _, line := range parseScriptLines(script) {
		tool, args, ok := parseScriptCommand(line)
		if !ok {
			continue
		}

		if tool == "note_insight" {
			// Small delay between note_insight and subsequent tool calls ensures
			// the HTTP response is fully processed. The MCP server uses a buffered
			// channel that preserves FIFO ordering, but rapid sequential HTTP
			// requests may interleave on slow systems.
			time.Sleep(10 * time.Millisecond)
		}

		_, err = c.CallTool(ctx, mcp.CallToolRequest{
			Params: mcp.CallToolParams{Name: tool, Arguments: args},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func parseScriptLines(s string) []string {
	s = strings.ReplaceAll(s, "\\n", "\n")
	lines := strings.Split(s, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			result = append(result, line)
		}
	}
	return result
}

func parseScriptCommand(line string) (tool string, args map[string]any, ok bool) {
	parts := strings.SplitN(line, " ", 2)
	if len(parts) < 2 {
		fmt.Fprintf(os.Stderr, "mock-agent: ignoring malformed command: %q\n", line)
		return "", nil, false
	}

	tool = parts[0]
	rest := parts[1]

	switch tool {
	case "task_complete":
		parts := strings.SplitN(rest, " ", 2)
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "mock-agent: ignoring malformed task_complete: %q\n", line)
			return "", nil, false
		}
		return tool, map[string]any{"status": parts[0], "summary": strings.Trim(parts[1], "\"")}, true
	case "note_insight":
		return tool, map[string]any{"text": strings.Trim(rest, "\"")}, true
	default:
		fmt.Fprintf(os.Stderr, "mock-agent: ignoring unknown command: %q\n", tool)
		return "", nil, false
	}
}
