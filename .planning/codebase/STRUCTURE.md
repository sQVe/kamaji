# Codebase Structure

**Analysis Date:** 2026-01-23

## Directory Layout

```
kamaji/
├── bin/                    # Compiled binaries (gitignored)
├── cmd/                    # CLI entry points
│   └── kamaji/
│       ├── main.go         # Cobra CLI root
│       ├── start.go        # `kamaji start` command
│       └── testdata/       # Test fixtures
├── internal/               # Private packages (not for external import)
│   ├── config/             # File I/O and persistence (YAML, MCP config)
│   ├── domain/             # Pure data types (no I/O)
│   ├── git/                # Git operations (branch, commit, reset)
│   ├── mcp/                # MCP server and tool handlers
│   ├── orchestrator/       # Sprint execution loop and state workflows
│   ├── output/             # Terminal UI and signal formatting
│   ├── process/            # Agent spawning and process management
│   ├── prompt/             # XML context assembly for agents
│   ├── statemachine/       # Pure state transitions
│   ├── testutil/           # Shared testing helpers
│   └── version/            # Version/build info
├── .planning/              # Project planning and milestones (gitignored)
├── .changes/               # Changelog entries (managed by changie)
├── coverage/               # Test coverage output (gitignored)
├── research/               # Research and design notes
├── .github/                # GitHub workflows
├── go.mod, go.sum          # Go module definition
├── Makefile                # Build automation
├── DESIGN.md               # Architecture and execution flow
├── CLAUDE.md               # Project guidelines
├── CONTRIBUTING.md         # Contribution guidelines
└── README.md               # Project overview
```

## Directory Purposes

**`cmd/kamaji/`:**

- Purpose: CLI entry point and command handlers
- Contains: Main function, subcommands, integration tests
- Key files: `main.go` (root command via Cobra), `start.go` (sprint runner)
- Pattern: Thin wrapper around `orchestrator.Run()`, delegates orchestration logic

**`internal/domain/`:**

- Purpose: Data types only—no I/O, no business logic
- Contains: `Sprint`, `Ticket`, `Task`, `State`, `TicketHistory`, `CompletedTask`, `FailedAttempt`
- Key files: `sprint.go` (Sprint/Ticket/Task structs), `state.go` (State struct), `history.go` (TicketHistory struct)
- Pattern: Plain structs with YAML tags, used throughout all layers

**`internal/config/`:**

- Purpose: File I/O, serialization, and configuration management
- Contains: Sprint loading, state persistence, history tracking, MCP config generation, plain mode flag
- Key files:
    - `sprint.go`: LoadSprint() with validation
    - `state.go`: LoadState() and SaveState() to/from .kamaji/state.yaml
    - `history.go`: LoadTicketHistory(), SaveTicketHistory(), RecordCompleted(), RecordFailed(), RecordInsight()
    - `mcp.go`: WriteMCPConfig() generates temp MCP connection file
    - `plain.go`: IsPlain() flag for terminal output styling
- Pattern: Graceful degradation (missing state files return zero-values, not errors)

**`internal/statemachine/`:**

- Purpose: Pure state transition logic—deterministic, no I/O
- Contains: NextTask(), Advance(), RecordPass(), RecordFail(), IsStuck()
- Key files: `statemachine.go` (all state logic)
- Pattern: Functions take state and sprint, return new state values (in-place mutation), caller handles persistence
- Constants: StuckThreshold = 3 (max consecutive failures per task)

**`internal/orchestrator/`:**

- Purpose: Main execution loop and task outcome workflows
- Contains: Sprint runner, handler for pass/fail/stuck workflows, task result parsing
- Key files:
    - `run.go`: Run() (main loop), runTask() (single task executor), createTicketBranch()
    - `handler.go`: Handler class with OnPass(), OnFail(), OnStuck() methods
    - `result.go`: TaskResult and signal-to-result conversion
- Pattern: Coordinates all subsystems, manages MCP server, deferring to specialized packages for I/O
- Dependencies: Imports all internal packages except testutil

**`internal/git/`:**

- Purpose: Git command execution and branch/commit management
- Contains: CreateBranch(), CommitChanges(), ResetToHead(), BranchExists()
- Key files: `git.go`
- Pattern: Wraps os/exec.Command, handles git command-line semantics
- Error types: ErrBranchExists, ErrNothingToCommit (for workflow decisions)

**`internal/mcp/`:**

- Purpose: MCP server lifecycle and tool signal handling
- Contains: HTTP server with SSE, tool registration, signal channel
- Key files:
    - `server.go`: Server struct, Start(), Shutdown(), Signals() channel
    - `tools.go`: Tool handlers (task_complete, note_insight), signal marshalling
- Pattern: Server runs in background, signals flow via buffered channel
- Port strategy: Accepts 0 for dynamic assignment (used in production)

**`internal/process/`:**

- Purpose: Agent spawning and process lifecycle
- Contains: SpawnClaude() (real agent), SpawnCommand() (test override), process.Process wrapper
- Key files:
    - `spawn.go`: SpawnClaude and SpawnCommand functions
    - `process.go`: Process wrapper with Start(), Wait(), Kill() (abstraction over os/exec.Cmd)
- Pattern: Config-driven spawning, caller owns process lifecycle (Wait/Kill)
- MCP integration: Passes MCP config path via args to claude; passes port/dir via env vars for test commands

**`internal/prompt/`:**

- Purpose: XML context assembly for agent sessions
- Contains: AssembleContext() (main entry), BuildPrompt() (formatting)
- Key files: `prompt.go`
- Pattern: Builds XML strings with proper escaping, loads history, formats task/rules/history sections
- Format: XML with sections for task, ticket, current, steps, verify, rules, history, instructions

**`internal/output/`:**

- Purpose: Terminal UI, progress reporting, signal formatting
- Contains: Progress writer, styled output, signal display, plain-text converter
- Key files:
    - `progress.go`: Progress writer with formatting
    - `signal.go`: Signal display formatting
    - `styles.go`: Terminal styling (with plain-mode fallback)
    - `writer.go`: Info/Error writers that apply styling
- Pattern: Respects config.IsPlain() for ASCII-only output

**`internal/testutil/`:**

- Purpose: Shared testing utilities and helpers
- Contains: Assertions, temp directory management, git test setup, platform detection
- Key files: `assert.go`, `tempdir.go`, `git.go`, `platform.go`
- Pattern: Reusable test infrastructure to keep test files focused

**`internal/version/`:**

- Purpose: Version and build metadata
- Contains: Version constant and Full() string formatting
- Key files: `version.go`
- Pattern: Set at build time via ldflags

## Key File Locations

**Entry Points:**

- `cmd/kamaji/main.go`: CLI root, invokes startCmd()
- `cmd/kamaji/start.go`: `kamaji start` command, calls orchestrator.Run()

**Configuration Files (checked into git):**

- `go.mod`: Go module definition with dependencies
- `Makefile`: Build, test, lint, changelog commands
- `.golangci.yml`: Go linter configuration (strict rules)
- `.pre-commit-config.yaml`: Pre-commit hooks (gofmt, golangci-lint)
- `.changie.yaml`: Changelog tool configuration

**Runtime State Files (generated, in .gitignore):**

- `.kamaji/state.yaml`: Current sprint progress (current_ticket, current_task, failure_count)
- `.kamaji/history/<ticket-name>.yaml`: Per-ticket history (completed, failed_attempts, insights)
- `.kamaji/.mcp.json`: Temporary MCP config (deleted after agent exit)

**Build Artifacts:**

- `bin/kamaji`: Compiled binary (generated by `make build`)

**Test Coverage:**

- `coverage/`: Generated coverage reports (gitignored)

## Naming Conventions

**Files:**

- Command files: `{command_name}.go` (e.g., `start.go` for `start` command)
- Test files: `{name}_test.go` (e.g., `statemachine_test.go`)
- Integration tests: `script_test.go` in cmd/
- Utility helpers: Descriptive names (e.g., `plain.go` for plain mode utilities)

**Directories:**

- Package names: lowercase, no underscores (e.g., `orchestrator`, `statemachine`)
- Internal packages: All under `internal/`
- Test data: `testdata/` subdirectory under package

**Functions:**

- Public: PascalCase (e.g., `LoadSprint`, `RunTask`)
- Private: camelCase (e.g., `sanitizeFilename`, `handleTaskComplete`)
- Constructors: `New{Type}` (e.g., `NewHandler`, `NewServer`)
- Getters: `{Field}()` (e.g., `Signals()`)
- State mutators: `Record{Action}` (e.g., `RecordPass`, `RecordFail`)

**Constants:**

- SCREAMING_SNAKE_CASE (e.g., `StuckThreshold`)
- Error types: `Err{Name}` (e.g., `ErrBranchExists`)

**Package-level variables:**

- Limited use; prefer function returns
- When needed: unexported (camelCase)

## Where to Add New Code

**New Feature (e.g., retry logic, parallel tickets):**

- Primary code: Extend `internal/orchestrator/run.go` or add new workflow handler
- State changes: Update `internal/statemachine/statemachine.go`
- Tests: Add tests adjacent to code in `*_test.go`
- CLI args: Add flags in `cmd/kamaji/start.go`, pass via `RunConfig` struct

**New Agent Integration (e.g., different spawner):**

- Spawner interface: Extend `orchestrator.ProcessSpawner` interface
- Implementation: New package under `internal/process/` (e.g., `internal/process/anthropic.go`)
- Configuration: Add flag to `cmd/kamaji/start.go`
- Tests: Mock spawner via `RunConfig.Spawner` field (pattern established in tests)

**New Data Type:**

- Domain struct: Add to appropriate file in `internal/domain/` (e.g., new ticket type → sprint.go)
- Serialization: Add load/save functions in `internal/config/` (follow LoadTicketHistory pattern)
- Tests: Pair with `*_test.go` in domain, then config

**New Git Operation:**

- Add function to `internal/git/git.go` following existing pattern (validation, runGit, error handling)
- Call from `orchestrator/handler.go` or `orchestrator/run.go`

**Terminal Output:**

- Add formatting function to `internal/output/` (separate for readability)
- Use `output.Print{Feature}()` naming
- Respect `config.IsPlain()` in styles

**Test Fixtures:**

- Shared utilities: `internal/testutil/` (e.g., temp dirs, git setup)
- Per-package fixtures: `testdata/` subdirectory
- Integration test scripts: `cmd/kamaji/testdata/script/`

## Special Directories

**`.planning/`:**

- Purpose: Project planning, milestones, phase tracking (generated by orchestrator)
- Generated: Yes (managed by gsd commands)
- Committed: No (.gitignored)
- Contents: Project state, roadmap, phase plans, summaries

**`.changes/`:**

- Purpose: Changelog fragments
- Generated: Yes (via `make change` command, uses changie)
- Committed: Yes (checked in until release)
- Contents: YAML files describing PRs/commits by category

**`.github/workflows/`:**

- Purpose: CI/CD pipeline definition
- Generated: No (hand-maintained)
- Committed: Yes
- Contents: GitHub Actions workflows for testing, linting, releasing

**`coverage/`:**

- Purpose: Test coverage reports
- Generated: Yes (via `make coverage`)
- Committed: No (.gitignored)

**`research/`:**

- Purpose: Design notes and research
- Generated: No (hand-maintained)
- Committed: Yes
- Contents: Exploration of technologies, patterns, decisions
