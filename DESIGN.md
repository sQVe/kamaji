# Kamaji

External Go CLI that orchestrates autonomous coding sprints by spawning fresh Claude Code sessions per task.

## Problem

Running the loop inside Claude Code causes:

- **Context pollution** - past work bleeds into future tasks
- **Context exhaustion** - less cognitive room as session grows
- **No parallelization** - single session bottleneck
- **Vendor lock-in** - tied to Claude Code hooks

## Solution

External orchestrator that:

- Owns the state machine
- Spawns fresh Claude Code sessions per task
- Manages git operations (commits, branches, resets)
- Exposes MCP server for completion signals

```
┌─────────────────────────────────────────┐
│              Kamaji (Go CLI)            │
│  - State machine                        │
│  - Reads/writes kamaji.yaml             │
│  - Spawns Claude Code sessions          │
│  - MCP server for signals               │
└─────────────────────────────────────────┘
           │                    ▲
           │ spawns             │ task_complete / note_insight
           ▼                    │
┌─────────────────────────────────────────┐
│         Claude Code Session             │
│  (fresh context per task)               │
└─────────────────────────────────────────┘
```

## File Structure

```
project/
  kamaji.yaml              # Sprint definition (checked into git)
  .kamaji/
    state.yaml             # Runtime state (current position, failure count)
    logs/
      <ticket-name>.yaml   # Per-ticket log (completed, failed, insights)
```

## Schema (state.yaml)

```yaml
current_ticket: 0          # Index into tickets array
current_task: 0            # Index into current ticket's tasks array
failure_count: 0           # Consecutive failures on current task (resets on pass)
```

## Schema (ticket log)

Stored in `.kamaji/logs/<ticket-name>.yaml`:

```yaml
ticket: login-form
completed:
  - task: "Create LoginForm component"
    summary: "Created LoginForm.tsx with Zod validation, loading state"
  - task: "Add unit tests"
    summary: "Added 5 tests covering validation and submit flow"
failed_attempts:
  - task: "Add OAuth integration"
    summary: "Tried passport.js but conflicts with existing session middleware"
insights:
  - "Codebase uses Zustand for state management"
  - "Validation schemas are in src/schemas/"
```

## Schema (kamaji.yaml)

```yaml
name: "Sprint name"
base_branch: main

rules:
  - "Use TypeScript strict mode"
  - "Follow existing patterns"

tickets:
  - name: login-form
    branch: feat/login-form
    description: "Create login form with validation"
    tasks:
      - description: "Create LoginForm component"
        steps:
          - "Add form validation using Zod"
          - "Handle submit with loading state"
        verify: "Component renders, validation works"

      - description: "Add unit tests"
        verify: "All tests pass"
```

## MCP Server

Kamaji runs an SSE-based MCP server that Claude Code connects to.

- **Transport**: Server-Sent Events (SSE)
- **Endpoint**: `http://localhost:<port>/mcp`
- **Port**: Dynamically assigned (or configurable via `--port`)

Claude Code is spawned with `--mcp-config` pointing to a temp file:

```json
{
  "mcpServers": {
    "kamaji": {
      "url": "http://localhost:9999/mcp"
    }
  }
}
```

## MCP Tools

**task_complete(status, summary)**

- status: "pass" | "fail"
- summary: what was done or why it failed
- Stored in ticket log, injected into future tasks

**note_insight(text)**

- Record discoveries useful for future tasks
- Stored in ticket log, injected into future tasks

## CLI

```bash
kamaji start           # Run sprint until done or stuck
kamaji start --dry-run # Show what would run
```

## Execution Flow

```
1. Read kamaji.yaml
2. Load state from .kamaji/state.yaml (or initialize)
3. Determine next task (first incomplete in first incomplete ticket)
4. If new ticket:
   a. git checkout <base_branch>
   b. git pull origin <base_branch>
   c. git checkout -b <ticket_branch>
5. Start MCP server
6. Build XML context (task + ticket + rules + history)
7. Spawn: claude -p "<context>" --mcp-config <kamaji-mcp> --dangerously-skip-permissions
8. Stream output to terminal
9. Wait for signal:
   a. task_complete(pass) → commit changes, store summary, next task
   b. task_complete(fail) → reset to HEAD, increment failures, store attempt, retry or stuck
   c. Process exits without signal → treat as fail
10. When all tasks done → exit success
    When stuck (3+ failures) → exit failure
```

## Context Injection (XML)

```xml
<task>
<ticket name="login-form" branch="feat/login-form">
Create login form with email/password validation
</ticket>

<current>
Create LoginForm component with email and password fields
</current>

<steps>
- Add form validation using Zod
- Handle submit with loading state
</steps>

<verify>
Component renders without errors, form validation rejects invalid email
</verify>
</task>

<rules>
Use TypeScript strict mode.
Follow existing patterns in src/.
</rules>

<history>
<completed>
- Created auth utility: Added loginUser to src/utils/auth.ts
</completed>

<failed_attempts>
- OAuth integration: passport.js conflicts with session middleware
</failed_attempts>

<insights>
- Codebase uses Zustand for state management
</insights>
</history>

<instructions>
Complete the task. Call task_complete(pass/fail, summary) when done.
Use note_insight() to record discoveries useful for future tasks.
</instructions>
```

## Git Handling

- **On pass**: Orchestrator commits with task summary as message
- **On fail**: `git reset --hard HEAD` (clean slate for retry)
- **On ticket start**: Create branch from latest base_branch

Claude focuses on coding. Orchestrator handles git.

## Failure Handling

- **Failure count**: Consecutive failures on the current task (resets to 0 on pass)
- **Stuck threshold**: 3 consecutive failures on the same task
- **On stuck**: Exit with failure, leave state intact for manual intervention
- **Exit without signal**: Treated as a failure (Claude crashed or forgot to call task_complete)

## V1 Scope (Minimal)

**Included:**

- Single `kamaji start` command
- Sequential task execution
- MCP tools (task_complete, note_insight)
- Ticket logs with history
- Git operations (branch, commit, reset)
- Streaming output

**Excluded (future):**

- Parallel ticket execution
- Worktree support
- Service management
- `kamaji status/stop/retry` commands

---

## Implementation Details

### Package Architecture

Two-package strategy separates concerns:

- `internal/domain/` — Pure data types (Sprint, Ticket, Task, State, TicketLog, CompletedTask, FailedAttempt)
- `internal/config/` — File I/O and persistence (LoadSprint, LoadState, SaveState, LoadTicketLog, SaveTicketLog)

Domain owns structures. Config owns serialization.

### Thread Safety

- Config state: `sync.RWMutex` for thread-safe access to global config
- Logger: `atomic.Bool` for debug flag

### State Machine

Pure functions with in-place mutation. Caller owns persistence.

```go
NextTask(sprint, state) *TaskInfo  // Returns nil at end, skips empty tickets
Advance(state, sprint)             // Increments task, handles ticket boundaries
RecordPass(state, sprint)          // Resets failure_count, calls Advance
RecordFail(state)                  // Increments failure_count
IsStuck(state) bool                // Returns failure_count >= 3
```

`TaskInfo` contains both domain objects and indices for orchestration context.

### Error Handling

- **Missing state files**: Return zero-value, not error (graceful fresh start)
- **Config errors**: Include context (indices, file paths) in messages
- **Filename sanitization**: `/` in ticket names becomes `-` in log paths

### Plain Mode

Logger and styles respect `config.IsPlain()` flag. ASCII-only output when enabled:

- `[ok]` — Success
- `Error:` — Error
- `->` — Info
- `Warning:` — Warning
- `[DEBUG]` — Debug