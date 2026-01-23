# Phase 5 Discovery: Claude Code CLI Interface

## Research Date: 2026-01-15

## Findings

### 1. Prompt Injection Flag

**Official syntax:** `claude -p "prompt"`

The `-p` flag enables headless/SDK mode - runs Claude Code programmatically without interactive mode, then exits.

### 2. MCP Configuration

**Recommended approach:** Use `.mcp.json` file in project root (not `--mcp-config` flag).

The `--mcp-config` flag has known bugs in v1.0.73+ where arguments after it are treated as config files.

**Correct `.mcp.json` format:**

```json
{
    "mcpServers": {
        "kamaji": {
            "type": "http",
            "url": "http://localhost:PORT/mcp"
        }
    }
}
```

### 3. Skip Permissions Flag

**Flag:** `--dangerously-skip-permissions`

- Skips all permission prompts
- Allows Claude to execute commands and file operations without confirmation
- Required for non-interactive/headless execution

### 4. Exit Codes

- Exit code 0: Success
- Exit code non-zero: Failure (max turns reached, error, or crash)

### 5. Output Capture

Stdout/stderr can be captured directly. For structured output:

- `--output-format json` for JSON responses
- `--output-format stream-json` for streaming events

## Command Structure

Based on research, the correct command is:

```bash
claude -p "<context>" --dangerously-skip-permissions
```

With `.mcp.json` in the working directory for MCP server configuration.

## Sources

- Claude Code CLI reference documentation
- Claude Code headless/SDK documentation
- MCP configuration documentation
- GitHub issues regarding `--mcp-config` flag bugs
