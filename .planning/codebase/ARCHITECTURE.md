# Architecture

**Analysis Date:** 2026-01-23

## Pattern Overview

**Overall:** External state machine orchestrator with MCP-based agent integration

**Key Characteristics:**

- State machine owns sprint progress independently of agent sessions
- Each task spawns a fresh agent context to prevent context pollution
- MCP server acts as bidirectional signal channel between orchestrator and agents
- Pure domain layer separated from I/O and persistence
- Sequential task execution with git-based state recovery

## Layers

**Domain Layer:**

- Purpose: Holds immutable data structures representing sprint state, tickets, tasks, and history
- Location: `internal/domain/`
- Contains: `Sprint`, `Ticket`, `Task`, `State`, `TicketHistory`, `CompletedTask`, `FailedAttempt`
- Depends on: Nothing (no external dependencies)
- Used by: All other layers for data structures

**State Machine Layer:**

- Purpose: Pure state transitions for sprint progress tracking
- Location: `internal/statemachine/`
- Contains: `NextTask()`, `Advance()`, `RecordPass()`, `RecordFail()`, `IsStuck()`
- Depends on: `domain/`
- Used by: `orchestrator/` for deterministic state updates

**Persistence Layer (Config):**

- Purpose: File I/O and serialization to/from YAML
- Location: `internal/config/`
- Contains: Sprint/state loading, history recording, MCP config generation, plain mode flag
- Depends on: `domain/`, standard library
- Used by: `orchestrator/`, `process/` for file operations

**Orchestration Layer:**

- Purpose: Manages task execution loop, coordinates subsystems
- Location: `internal/orchestrator/`
- Contains: `Run()` (main loop), `Handler` (pass/fail/stuck workflows), `TaskResult` (outcome handling)
- Depends on: `domain/`, `statemachine/`, `config/`, `git/`, `mcp/`, `process/`, `prompt/`, `output/`
- Used by: `cmd/kamaji/`

**Git Integration Layer:**

- Purpose: Executes git operations (branch creation, commits, resets)
- Location: `internal/git/`
- Contains: `CreateBranch()`, `CommitChanges()`, `ResetToHead()`, branch existence checks
- Depends on: Standard library (os/exec)
- Used by: `orchestrator/` for branch/commit/reset operations

**MCP Server Layer:**

- Purpose: Manages MCP server lifecycle and tool signal handling
- Location: `internal/mcp/`
- Contains: `Server` (HTTP-based MCP with SSE), tool handlers, signal channel
- Depends on: `mark3labs/mcp-go`, standard library
- Used by: `orchestrator/` to receive task completion and insight signals

**Process Spawning Layer:**

- Purpose: Launches Claude agent session with MCP configuration
- Location: `internal/process/`
- Contains: `SpawnClaude()` (real agent), `SpawnCommand()` (test override), process lifecycle management
- Depends on: `config/`, standard library
- Used by: `orchestrator/` to spawn agent tasks

**Prompt Assembly Layer:**

- Purpose: Constructs XML context injection for agent sessions
- Location: `internal/prompt/`
- Contains: `AssembleContext()`, `BuildPrompt()`, XML formatting for task/rules/history
- Depends on: `domain/`, `statemachine/`, `config/`
- Used by: `orchestrator/` to generate agent prompt

**Output/UI Layer:**

- Purpose: Formats terminal output with optional plain-text mode
- Location: `internal/output/`
- Contains: Progress reporting, signal formatting, styled writers, plain-text conversion
- Depends on: `config/` (plain mode flag), standard library
- Used by: All layers for user-facing messages

**Entry Point:**

- Location: `cmd/kamaji/`
- Contains: `main.go` (Cobra CLI root), `start.go` (sprint command), `script_test.go` (integration tests)
- Depends on: `orchestrator/`, `spf13/cobra`
- Responsibilities: Parse CLI flags, invoke Run(), exit with appropriate code

## Data Flow

**Sprint Execution Loop:**

1. User runs `kamaji start`
2. `cmd/kamaji/start.go` loads current directory as workDir and kamaji.yaml path
3. `orchestrator.Run()` begins:
    - Load sprint definition from kamaji.yaml via `config.LoadSprint()`
    - Load execution state from .kamaji/state.yaml via `config.LoadState()` (or zero-value if missing)
    - Start MCP server on dynamic port
    - Enter sequential task loop
4. For each task:
    - Get next task via `statemachine.NextTask()`
    - If new ticket (CurrentTask==0 && FailureCount==0): create branch via `git.CreateBranch()`
    - Load ticket history via `config.LoadTicketHistory()`
    - Assemble prompt XML via `prompt.AssembleContext()`
    - Spawn Claude with `process.SpawnClaude()` and MCP config
    - Stream output to terminal via `output.NewInfoWriter()`
    - Wait for MCP signals on `server.Signals()` channel
5. Signal handling:
    - `task_complete(pass, summary)` → `handler.OnPass()` → commit, record, advance, save state
    - `task_complete(fail, summary)` → `handler.OnFail()` → reset, record, increment failures, save state
    - `note_insight(text)` → `config.RecordInsight()` → injected into next task
    - No signal on exit → treated as fail
6. Stuck detection:
    - If failure_count >= 3 → `handler.OnStuck()` → exit with failure, state preserved
7. Completion:
    - If no more tasks → exit success

**State Persistence:**

- State lives in `.kamaji/state.yaml` with three fields: current_ticket, current_task, failure_count
- State is loaded at start and saved after every signal (pass/fail/stuck)
- On agent crash or no signal: state remains unchanged, next run retries same task
- Ticket history lives in `.kamaji/history/<ticket-name>.yaml` with completed/failed/insights

**MCP Server Communication:**

- Orchestrator: Exposes SSE endpoint at `http://localhost:<port>/mcp`
- Agent: Connects via MCP, calls `task_complete()` or `note_insight()` tools
- Channel: Signal objects flow back to orchestrator via buffered channel
- Signals drain: On process exit, orchestrator drains pending signals for insights/completion that arrived concurrently

## Key Abstractions

**TaskInfo:**

- Purpose: Groups task metadata with indices for orchestration context
- Fields: TicketIndex, TaskIndex, Ticket*, Task*
- Used by: State machine returns it, orchestrator uses it for prompt and logging
- Pattern: Separates domain objects (for agent context) from indices (for state updates)

**Handler:**

- Purpose: Stateful workflow manager for task outcomes
- Pattern: Single-threaded, owns none of its data, caller must manage state lifecycle
- Methods: OnPass(), OnFail(), OnStuck(), IsStuck()
- Responsibility: Commits, records history, updates state, resets on fail

**ProcessSpawner Interface:**

- Purpose: Abstracts process creation for testing
- Implementations: defaultSpawner (real Claude), commandSpawner (test override)
- Used by: Orchestrator accepts optional spawner for injection
- Pattern: Dependency injection for testability

**TaskResult:**

- Purpose: Outcome of single task execution
- Fields: Status (pass/fail/none), Summary, ExitCode
- Methods: Passed(), ResultFromSignal()
- Pattern: Decouples signal parsing from state updates

## Entry Points

**CLI Entry:**

- Location: `cmd/kamaji/main.go`
- Triggers: User invocation `kamaji start`
- Responsibilities: Parse args, run orchestrator, exit with code

**Sprint Start:**

- Location: `cmd/kamaji/start.go::startCmd()`
- Triggers: `kamaji start` command
- Responsibilities: Resolve workDir, call `orchestrator.Run()`, handle result

**Orchestration Loop:**

- Location: `internal/orchestrator/run.go::Run()`
- Triggers: Called from start.go
- Responsibilities: Load config, spawn tasks, coordinate handler, manage MCP server

## Error Handling

**Strategy:** Error propagation with context for debugging

**Patterns:**

- Missing state file: Return zero-value State (graceful fresh start)
- Missing kamaji.yaml: Return error (required file)
- Git failures: Propagate with stderr context (e.g., "checkout main (fatal: ...): error")
- Stuck task: Exit gracefully with state preserved (not an error, handled result)
- Config validation: Validate Sprint on load; return formatted errors with indices

**Special cases:**

- ErrNothingToCommit: Not an error in OnPass flow; task advances anyway
- ErrBranchExists: Not an error on retry; branch reused with existing content
- Process exit without signal: Treated as fail (not error), increments failure_count

## Cross-Cutting Concerns

**Logging:**

- Structured via `slog` in MCP server for tool debugging
- User-facing output via `output/` package with styled writers
- Plain mode controlled by `config.IsPlain()` flag

**Validation:**

- Sprint: Required fields checked on load (name, ticket names, task descriptions)
- Files: Existence checks via os.Stat before directory operations
- State: Zero-value acceptable; indices bounds-checked in state machine

**Authentication:**

- None required; Kamaji is subprocess. Agent session inherits user's claude auth.

**File permissions:**

- State files: 0o600 (owner read/write only)
- Directory .kamaji: 0o750 (owner full, group/other read+execute)
- MCP config temp file: Deleted after process exit
