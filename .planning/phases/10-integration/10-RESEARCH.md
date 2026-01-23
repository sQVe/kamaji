# Phase 10: Integration - Research

**Researched:** 2026-01-20
**Domain:** Go orchestration loop combining existing components
**Confidence:** HIGH

## Summary

Phase 10 integrates all existing components into a working `kamaji start` command. The research confirms all building blocks are complete and ready for wiring. The orchestration pattern is straightforward: a sequential loop that spawns agents, waits for signals, and handles outcomes.

Key findings from codebase review:

- All 9 prerequisite phases implemented and tested
- Handler already orchestrates pass/fail/stuck outcomes with git operations
- MCP server provides signal channel for task completion
- Process spawning creates Claude subprocess with MCP config
- Output module provides all progress/status formatting

**Primary recommendation:** Create a simple sequential loop in the `orchestrator` package. Use select on MCP signals and process exit. Wire to Cobra `start` command with minimal flags.

## Standard Stack

All components already exist. No new dependencies needed.

### Core (existing)

| Package                 | Purpose                                   | Status   |
| ----------------------- | ----------------------------------------- | -------- |
| `internal/orchestrator` | Handler for pass/fail/stuck               | Complete |
| `internal/mcp`          | Server with signals channel               | Complete |
| `internal/process`      | SpawnClaude with config                   | Complete |
| `internal/prompt`       | AssembleContext                           | Complete |
| `internal/git`          | CreateBranch, CommitChanges, ResetToHead  | Complete |
| `internal/config`       | LoadSprint, LoadState, SaveState, history | Complete |
| `internal/statemachine` | NextTask, Advance, IsStuck                | Complete |
| `internal/output`       | Progress, signals, status display         | Complete |

### CLI (existing)

| Package                  | Purpose           | Status                 |
| ------------------------ | ----------------- | ---------------------- |
| `github.com/spf13/cobra` | Command framework | In use                 |
| `cmd/kamaji/main.go`     | Root command only | Needs start subcommand |

### No new dependencies needed

The integration uses only existing code. No external libraries required.

## Architecture Patterns

### Pattern 1: Run Function

**What:** Single orchestration function that owns the execution loop
**Why:** Separates orchestration logic from CLI wiring, enables testing

```go
// internal/orchestrator/run.go
type RunConfig struct {
    WorkDir    string
    SprintPath string
    DryRun     bool
}

type RunResult struct {
    Success     bool
    TasksRun    int
    Stuck       bool
    StuckReason string
}

func Run(ctx context.Context, cfg RunConfig) (*RunResult, error) {
    // 1. Load sprint and state
    // 2. Loop: get next task, spawn agent, wait for signal, handle outcome
    // 3. Return result
}
```

### Pattern 2: Signal/Exit Race

**What:** Use select to handle MCP signal vs process exit
**Why:** Process may exit before signaling (crash) or signal without exiting

```go
// Wait for either signal or process exit
select {
case sig := <-server.Signals():
    // Got signal, process it
    result = orchestrator.ResultFromSignal(sig)
case <-processExited:
    // Process exited without signal
    result = orchestrator.NoSignalResult()
}

// Always wait for process cleanup
_ = spawnResult.Process.Wait()
```

### Pattern 3: Ticket Boundary Detection

**What:** Track ticket transitions to trigger branch creation
**Why:** git.CreateBranch only runs when starting a new ticket

```go
prevTicket := state.CurrentTicket
// ... handle outcome which may call Advance ...
if state.CurrentTicket != prevTicket {
    // Transitioned to new ticket
    if taskInfo := statemachine.NextTask(state, sprint); taskInfo != nil {
        git.CreateBranch(workDir, sprint.BaseBranch, taskInfo.Ticket.Branch)
    }
}
```

### Pattern 4: Context Cancellation

**What:** Respect context.Context for graceful shutdown
**Why:** Allows ctrl-c to cleanly terminate

```go
func Run(ctx context.Context, cfg RunConfig) (*RunResult, error) {
    for {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
        }
        // ... run task ...
    }
}
```

### Anti-Patterns to Avoid

- **Multiple goroutines for simple sequential loop:** Unnecessary complexity
- **Spawning next agent before current exits:** Resource leak, MCP port conflict
- **Ignoring MCP server shutdown:** Leaves port bound
- **Not removing MCP config file:** Leaves `.mcp.json` in project

## Don't Hand-Roll

| Problem               | Don't Build    | Use Instead                                | Why                                 |
| --------------------- | -------------- | ------------------------------------------ | ----------------------------------- |
| Task outcome handling | Custom logic   | `orchestrator.Handler`                     | Already handles git, history, state |
| Signal conversion     | Manual parsing | `orchestrator.ResultFromSignal`            | Normalizes edge cases               |
| Progress calculation  | Loop counting  | `output.SprintStatus`                      | Already implemented                 |
| Ticket boundary logic | Manual checks  | Compare `state.CurrentTicket` before/after | Built into statemachine             |

**Key insight:** All complex logic exists. Integration is primarily wiring and sequencing.

## Common Pitfalls

### Pitfall 1: MCP Server Lifecycle

**What goes wrong:** Port conflicts, signals lost, server not cleaned up
**Why it happens:** Server started once, used across multiple tasks
**How to avoid:** Start server before loop, shutdown in defer, one server for entire sprint
**Warning signs:** "address already in use", signals not received

```go
server := mcp.NewServer(mcp.WithPort(0))
port, err := server.Start()
if err != nil {
    return nil, err
}
defer server.Shutdown(ctx)
```

### Pitfall 2: Process Exit Without Signal

**What goes wrong:** Sprint hangs waiting for signal that never comes
**Why it happens:** Claude crashes, max turns, or forgets to call task_complete
**How to avoid:** Use select with process exit channel, treat as failure
**Warning signs:** Sprint hangs indefinitely on failing task

```go
done := make(chan struct{})
go func() {
    _ = spawnResult.Process.Wait()
    close(done)
}()

select {
case sig := <-server.Signals():
    result = orchestrator.ResultFromSignal(sig)
case <-done:
    result = orchestrator.NoSignalResult()
}
```

### Pitfall 3: MCP Config File Cleanup

**What goes wrong:** `.mcp.json` left in project directory
**Why it happens:** Forgetting to remove after process exits
**How to avoid:** Always remove in defer, handle errors gracefully
**Warning signs:** Spurious `.mcp.json` files in target projects

```go
spawnResult, err := process.SpawnClaude(cfg)
if err != nil {
    return nil, err
}
defer os.Remove(spawnResult.ConfigPath)
```

### Pitfall 4: First Task on Fresh Sprint

**What goes wrong:** Missing branch creation on first task
**Why it happens:** Only checking ticket transitions, not fresh start
**How to avoid:** Create branch before first task if state.CurrentTask == 0
**Warning signs:** First task runs on wrong branch

```go
taskInfo := statemachine.NextTask(state, sprint)
if state.CurrentTask == 0 && state.FailureCount == 0 {
    // First task of ticket (fresh or new ticket)
    if err := git.CreateBranch(workDir, sprint.BaseBranch, taskInfo.Ticket.Branch); err != nil {
        return nil, err
    }
}
```

### Pitfall 5: Insight Recording During Task

**What goes wrong:** Insights not persisted, or not associated with correct ticket
**Why it happens:** note_insight signals arrive during task execution, need immediate handling
**How to avoid:** Process signals in loop, record immediately
**Warning signs:** Insights lost or in wrong ticket log

```go
// While waiting for task_complete, also handle note_insight
for {
    select {
    case sig := <-server.Signals():
        if sig.Tool == mcp.SignalToolNoteInsight {
            config.RecordInsight(workDir, ticketName, sig.Summary)
            output.PrintSignal(sig)
            continue
        }
        // task_complete - break out
        return orchestrator.ResultFromSignal(sig)
    case <-done:
        return orchestrator.NoSignalResult()
    }
}
```

## Code Examples

### Main Loop Structure

```go
// internal/orchestrator/run.go
func Run(ctx context.Context, cfg RunConfig) (*RunResult, error) {
    sprint, err := config.LoadSprint(cfg.SprintPath)
    if err != nil {
        return nil, err
    }

    state, err := config.LoadState(cfg.WorkDir)
    if err != nil {
        return nil, err
    }

    server := mcp.NewServer(mcp.WithPort(0))
    port, err := server.Start()
    if err != nil {
        return nil, err
    }
    defer server.Shutdown(ctx)

    handler := NewHandler(cfg.WorkDir, state, sprint)
    var tasksRun int

    for {
        select {
        case <-ctx.Done():
            return &RunResult{Success: false}, ctx.Err()
        default:
        }

        taskInfo := statemachine.NextTask(state, sprint)
        if taskInfo == nil {
            output.PrintSprintComplete(sprint, state)
            return &RunResult{Success: true, TasksRun: tasksRun}, nil
        }

        // Branch creation for new ticket
        if state.CurrentTask == 0 && state.FailureCount == 0 {
            if err := git.CreateBranch(cfg.WorkDir, sprint.BaseBranch, taskInfo.Ticket.Branch); err != nil {
                return nil, err
            }
            output.PrintBranchCreated(taskInfo.Ticket.Branch)
        }

        output.PrintTaskStart(taskInfo, sprint)

        result, err := runTask(ctx, cfg.WorkDir, port, sprint, state, taskInfo, server)
        if err != nil {
            return nil, err
        }
        tasksRun++

        output.PrintSignal(mcp.Signal{
            Tool:    mcp.SignalToolTaskComplete,
            Status:  result.Status,
            Summary: result.Summary,
        })

        if result.Passed() {
            if err := handler.OnPass(taskInfo.Ticket.Name, taskInfo.Task.Description, result.Summary); err != nil {
                return nil, err
            }
        } else {
            if err := handler.OnFail(taskInfo.Ticket.Name, taskInfo.Task.Description, result.Summary); err != nil {
                return nil, err
            }
            if handler.IsStuck() {
                if err := handler.OnStuck(); err != nil {
                    return nil, err
                }
                return &RunResult{
                    Success:     false,
                    TasksRun:    tasksRun,
                    Stuck:       true,
                    StuckReason: result.Summary,
                }, nil
            }
        }
    }
}
```

### Task Execution with Signal Handling

```go
func runTask(ctx context.Context, workDir string, port int, sprint *domain.Sprint, state *domain.State, taskInfo *statemachine.TaskInfo, server *mcp.Server) (orchestrator.TaskResult, error) {
    prompt, err := prompt.AssembleContext(sprint, state, workDir)
    if err != nil {
        return orchestrator.TaskResult{}, err
    }

    spawnResult, err := process.SpawnClaude(process.SpawnConfig{
        Prompt:  prompt,
        MCPPort: port,
        WorkDir: workDir,
        Stdout:  os.Stdout,
        Stderr:  os.Stderr,
    })
    if err != nil {
        return orchestrator.TaskResult{}, err
    }
    defer os.Remove(spawnResult.ConfigPath)

    done := make(chan struct{})
    go func() {
        _ = spawnResult.Process.Wait()
        close(done)
    }()

    ticketName := taskInfo.Ticket.Name

    for {
        select {
        case <-ctx.Done():
            _ = spawnResult.Process.Kill()
            return orchestrator.TaskResult{}, ctx.Err()
        case sig := <-server.Signals():
            if sig.Tool == mcp.SignalToolNoteInsight {
                _ = config.RecordInsight(workDir, ticketName, sig.Summary)
                output.PrintSignal(sig)
                continue
            }
            <-done // Wait for process to exit
            return orchestrator.ResultFromSignal(sig), nil
        case <-done:
            return orchestrator.NoSignalResult(), nil
        }
    }
}
```

### CLI Command

```go
// cmd/kamaji/start.go
func startCmd() *cobra.Command {
    var dryRun bool

    cmd := &cobra.Command{
        Use:   "start",
        Short: "Run sprint until done or stuck",
        RunE: func(cmd *cobra.Command, args []string) error {
            workDir, err := os.Getwd()
            if err != nil {
                return err
            }

            result, err := orchestrator.Run(cmd.Context(), orchestrator.RunConfig{
                WorkDir:    workDir,
                SprintPath: filepath.Join(workDir, "kamaji.yaml"),
                DryRun:     dryRun,
            })
            if err != nil {
                return err
            }

            if !result.Success {
                os.Exit(1)
            }
            return nil
        },
    }

    cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would run without executing")
    return cmd
}
```

## Edge Cases to Handle

### 1. Empty Sprint (no tickets)

```go
if len(sprint.Tickets) == 0 {
    output.PrintInfo("Sprint has no tickets")
    return &RunResult{Success: true}, nil
}
```

### 2. Empty Ticket (no tasks)

The statemachine.NextTask already handles this by advancing past empty tickets.

### 3. Resume After Crash

State is persisted after every outcome. On restart:

- Load state from `.kamaji/state.yaml`
- Continue from `CurrentTicket`/`CurrentTask`
- Failure count preserved for stuck detection

### 4. Branch Already Exists

```go
if err := git.CreateBranch(...); err != nil {
    // Check if branch exists, checkout instead of create
    if strings.Contains(err.Error(), "already exists") {
        if _, _, err := runGit(workDir, "checkout", ticketBranch); err != nil {
            return err
        }
    } else {
        return err
    }
}
```

### 5. Dry Run Mode

```go
if cfg.DryRun {
    for {
        taskInfo := statemachine.NextTask(state, sprint)
        if taskInfo == nil {
            break
        }
        output.PrintInfo(fmt.Sprintf("Would run: %s", taskInfo.Task.Description))
        statemachine.Advance(state, sprint)
    }
    return &RunResult{Success: true}, nil
}
```

## Testing Strategy

### Unit Tests

| Test                            | What it verifies        |
| ------------------------------- | ----------------------- |
| `TestRun_CompletesAllTasks`     | Loop runs to completion |
| `TestRun_HandlesStuck`          | Exits after 3 failures  |
| `TestRun_RespectsContext`       | Cancellation works      |
| `TestRun_CreatesTicketBranches` | Branch per ticket       |
| `TestRun_DryRunMode`            | No actual execution     |

### Integration Approach

Use mocked components since real Claude is external:

- Mock process that exits with controllable signal
- Real MCP server (in-process)
- Real git operations in temp dir
- Real state/history persistence

### Script Test

The existing `cmd/kamaji/script_test.go` pattern can test the full CLI:

```go
func TestKamaji_Start_DryRun(t *testing.T) {
    // Create temp dir with kamaji.yaml
    // Run: kamaji start --dry-run
    // Verify output shows what would run
}
```

## Open Questions

### 1. Signal Timeout

**What we know:** Process can hang indefinitely if Claude gets stuck in a loop without calling task_complete
**What's unclear:** Should there be a timeout per task?
**Recommendation:** V1 has no timeout. V2 could add `--timeout` flag.

### 2. Parallel note_insight Signals

**What we know:** Multiple note_insight calls can occur during a task
**What's unclear:** Order preservation and race conditions
**Recommendation:** The signal channel is buffered (10). Record in order received. File locking in `config.RecordInsight` handles concurrent writes.

## Sources

### Primary (HIGH confidence)

- Codebase analysis of all internal packages
- DESIGN.md execution flow specification
- Existing phase summaries (01-09)

### Implementation References

- `internal/orchestrator/handler.go` - Handler pattern
- `internal/orchestrator/result.go` - TaskResult type
- `internal/mcp/server.go` - Signal channel pattern
- `internal/process/spawn.go` - SpawnClaude interface

## Metadata

**Confidence breakdown:**

- Loop structure: HIGH - follows DESIGN.md exactly
- Signal handling: HIGH - MCP server already tested
- Edge cases: HIGH - all covered in handler tests
- CLI wiring: HIGH - Cobra patterns established

**Research date:** 2026-01-20
**Valid until:** 2026-02-20 (30 days - integration of stable components)
