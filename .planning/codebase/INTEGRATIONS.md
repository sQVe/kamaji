# External Integrations

**Analysis Date:** 2026-01-23

## APIs & External Services

**Model Context Protocol (MCP):**

- MCP Server - Hosts tools for agent communication
    - SDK/Client: `github.com/mark3labs/mcp-go` v0.43.2
    - Protocol: HTTP on localhost (dynamic port assignment)
    - Served on: `/mcp` endpoint
    - Auth: None (localhost only, internal communication)

**Claude Code Integration:**

- Process: Spawned via `claude` CLI command with arguments
    - `--print <prompt>` - Pass context as initial prompt
    - `--dangerously-skip-permissions` - Auto-accept all tool calls
    - `--output-format stream-json` - Structured output parsing
    - Launched in working directory context
    - MCP server port passed via environment: `KAMAJI_MCP_PORT`

## Data Storage

**Databases:**

- None (no external database integration)

**File Storage:**

- Local filesystem only
- Stores state in: `.kamaji/state.yaml`
- Stores history in: `.kamaji/history/<ticket_name>.yaml`
- Stores history locks in: `.kamaji/history/<ticket_name>.lock`
- Generated config in: `.mcp.json`

**Caching:**

- None (stateless between executions except for .kamaji directory)

## Authentication & Identity

**Auth Provider:**

- None (local tool, no authentication required)
- Git operations use system git credentials (no explicit auth in kamaji)

## Monitoring & Observability

**Error Tracking:**

- None configured

**Logs:**

- Standard output/error streams via output package
- `output.PrintInfo()` - Information messages to stdout
- `output.PrintError()` - Error messages to stderr
- `output.PrintSignal()` - Agent signal events
- Structured logging via `log/slog` for MCP server debug logs

## CI/CD & Deployment

**Hosting:**

- GitHub (source repository)
- GitHub Releases (binary distribution)

**CI Pipeline:**

- GitHub Actions
- Triggers: Push to main, PR creation/updates
- Jobs:
    - Lint: golangci-lint v2.8.0
    - Test: Unit tests (all platforms) + Integration tests (all platforms)
    - Build: Linux/macOS release binaries
    - Changeset validation: Ensures changelog entries for code changes
    - CodeQL: Security scanning
    - Dependency audit: go list -u check

## Environment Configuration

**Required env vars for spawning agents:**

- `KAMAJI_MCP_PORT` - Set when spawning Claude Code (port number as string)
- `KAMAJI_WORK_DIR` - Set when spawning Claude Code (working directory path)

**Git config requirements (set by CI/CD):**

- `user.email` - For commits
- `user.name` - For commits
- `init.defaultBranch` - Default branch for new repos
- `commit.gpgsign false` - Disable GPG signing in CI

**Secrets location:**

- GitHub Secrets (used only for GITHUB_TOKEN in release workflow)
- No sensitive configuration stored in code

## Webhooks & Callbacks

**Incoming:**

- None (standalone CLI tool)

**Outgoing:**

- Git webhooks: Not used directly (agent commits trigger git hooks)
- MCP tool callbacks:
    - `task_complete(status, summary)` - Agent signals task completion
    - `note_insight(text)` - Agent records discoveries for future reference

## System Integration Points

**Git Operations:**

- `git.CreateBranch()` - Creates feature branches from base branch
- `git.CommitChanges()` - Stages and commits changes with message
- `git.ResetToHead()` - Discards changes and cleans untracked files
- All operations via `exec.Command("git", ...)` - Uses system git binary

**Process Management:**

- `process.SpawnClaude()` - Launches Claude Code CLI with:
    - MCP configuration written to `.mcp.json`
    - Prompt context passed via stdin
    - Working directory isolated per sprint
    - Process lifecycle management (Wait, Kill)

**File Locking:**

- File-based locks in `.kamaji/history/` for concurrent write prevention
- Uses `os.OpenFile` with `O_CREATE|O_EXCL` flags (atomic creation)
- 200-iteration retry loop with 10ms backoff (max 2s wait)

---

_Integration audit: 2026-01-23_
