# Codebase Concerns

**Analysis Date:** 2026-01-23

## Signal Channel Overflow

**Issue:** Dropped signals on full channel buffer

- **Files:** `internal/mcp/server.go` (lines 64-67, 79-82)
- **Problem:** Signal channel has fixed buffer size of 10. If MCP handlers send signals faster than orchestrator consumes them, signals are silently dropped with only debug logs
- **Impact:** Lost task completion signals or insights during high-frequency signaling could cause tasks to appear hung or insights to disappear without detection
- **Current state:** Signals are non-blocking; dropped signals only log at debug level—no retry or notification to caller
- **Improvement path:** Either increase buffer size, use unbuffered channel with explicit blocking/timeout handling, or add explicit error return when signals drop. Consider adding observability metric for dropped signals

## Process Exit Race Condition

**Issue:** Process exit and concurrent signal receiving

- **Files:** `internal/orchestrator/run.go` (lines 214-264)
- **Problem:** Complex race between process exit (done channel) and signal handling. After process exits, code tries to drain pending signals with default select case that could exit prematurely if no signals arrive in first iteration
- **Symptoms:** If process exits with pending signals in MCP channel, those signals may not be captured
- **Impact:** Task completion signal arriving concurrently with process exit could be lost, leaving state inconsistent (no commit made but task shown as pending retry)
- **Fix approach:** Use bounded timeout for draining signals or ensure signals are read before checking process exit status. Add test for this race condition

## Ignored RecordInsight Errors in Hot Path

**Issue:** Silent failure when recording insights during task execution

- **Files:** `internal/orchestrator/run.go` (lines 234, 256)
- **Problem:** Calls to `config.RecordInsight()` ignore errors with `_ =` in the hot execution loop. If insight recording fails (file lock timeout, I/O error), the error is silently dropped and execution continues
- **Impact:** Loss of insights without any indication to orchestrator or agent. Failed insights could indicate permission issues or disk problems that operator never sees
- **Fix approach:** Log errors from RecordInsight, but don't fail the entire task. Consider implementing retry with exponential backoff for lock acquisition

## File Lock Timeout Without Backoff Guarantee

**Issue:** History file lock acquisition may fail under concurrent pressure

- **Files:** `internal/config/history.go` (lines 124-154)
- **Problem:** Loop tries 200 times with 10ms sleep between attempts, totaling max 2 seconds. Under high concurrency or on slow filesystems, legitimate lock holders may not release in time, causing spurious "timeout acquiring history lock" errors
- **Symptoms:** Orchestrator crashes when multiple concurrent insight/completion recordings happen during process exit signal draining
- **Impact:** Sprint execution fails mid-task with no recovery path. State left at last saved checkpoint, requiring manual intervention
- **Trigger:** Run task that sends many insights rapidly, then fails/completes just as another signal arrives
- **Workaround:** Retry RunConfig at orchestrator level
- **Fix approach:** Implement exponential backoff instead of fixed intervals, or use context with longer deadline for lock acquisition. Add metric for lock contention

## Missing Process Cleanup on Context Cancellation

**Issue:** Zombie process if orchestrator context is cancelled

- **Files:** `internal/orchestrator/run.go` (lines 223-226, 240-243)
- **Problem:** When context is cancelled (e.g., parent timeout), code kills process and waits. However, the deferred cleanup for ConfigPath (line 210) may not execute if the wait channel blocks. If process.Wait() doesn't return quickly, cleanup leaks
- **Impact:** Temporary MCP config file left in project directory if task forcibly terminated
- **Fix approach:** Ensure cleanup happens via explicit goroutine or non-deferred cleanup wrapper, not relying on defer across goroutine boundaries

## State Rollback Not Atomic

**Issue:** Incomplete state recovery on SaveState failure

- **Files:** `internal/orchestrator/handler.go` (lines 47-58, 76-83)
- **Problem:** Code rolls back state in memory after SaveState fails, but if this rollback itself is interrupted or the handler is destroyed, the in-memory state stays corrupted. Next Run() will load from disk (which was never updated), creating divergence
- **Symptoms:** State machine advances but save fails; next run resets to old state, duplicating work
- **Impact:** Duplicate task execution or lost progress if operator restarts without noticing the error
- **Current state:** Error is returned to caller; caller should retry, but no automatic retry exists
- **Fix approach:** Either make SaveState transactional with temp file + rename, or log a specific error code that suggests state is inconsistent and needs manual review

## Unbounded Sprint Completion Loop

**Issue:** No timeout or heartbeat in task execution loop

- **Files:** `internal/orchestrator/run.go` (lines 100-162)
- **Problem:** Main loop has no timeout or progress monitoring. If a single task never signals completion and the process hangs silently, orchestrator waits forever with no indication of stuck state
- **Impact:** Sprint runner appears frozen; CI pipeline hangs indefinitely; operator has no way to detect hang without external monitoring
- **Current state:** Only way to break is to kill the process
- **Fix approach:** Add context timeout to Run(), or add task-level timeout that forces task_complete(fail) if deadline passes. Add debug logging of task state transitions

## Git Pull Silently Fails

**Issue:** Network error during branch creation is ignored

- **Files:** `internal/git/git.go` (line 78)
- **Problem:** `git pull origin <branch>` output and error are explicitly discarded. If pull fails (network error, auth failure), branch creation continues with stale base branch. No logging or indication that pull failed
- **Impact:** Agent creates changes on old codebase, leading to merge conflicts or built against old code. Operator unaware of network issues
- **Fix approach:** Log pull errors at info/warning level (don't fail the task, since pull can fail offline), but expose the fact that pull was attempted and failed

## Signal Status Validation Missing

**Issue:** task_complete status values not validated

- **Files:** `internal/mcp/tools.go` (HandleTaskComplete)
- **Problem:** MCP handler accepts any string for status parameter. Invalid status values like "maybe" or "unknown" are passed through without validation, only formatted by output layer
- **Impact:** Ambiguous task state if agent sends invalid status. State machine uses status to decide pass/fail, so invalid status creates undefined behavior
- **Fix approach:** Add explicit validation in HandleTaskComplete to reject non-"pass"/"fail" status values with error. Document expected values in MCP tool description

## Lock File Cleanup Race

**Issue:** Lock file removal could fail silently

- **Files:** `internal/config/history.go` (lines 150-153)
- **Problem:** Lock file removal is ignored (`_ = os.Remove(lockPath)`). If removal fails (permission denied, file locked by antivirus), next acquisition loop will wait for timeout before discovering file exists
- **Impact:** Lock file accumulation if cleanup fails, causing subsequent operations to timeout unnecessarily
- **Fix approach:** Log removal errors at debug level (don't fail the operation since lock is released), but increment metric to detect cleanup issues

## No Validation of Sprint/State Consistency

**Issue:** No validation when NextTask called with invalid state

- **Files:** `internal/statemachine/statemachine.go` (lines 16-32)
- **Problem:** NextTask silently returns nil if CurrentTicket or CurrentTask out of bounds. No check that these indices are even valid. If state file is corrupted with huge index values, NextTask returns nil (appears as sprint complete) rather than error
- **Impact:** Sprint appears complete when state is actually corrupted. Operator doesn't realize state file needs recovery
- **Fix approach:** Return error from NextTask when indices are invalid, or add validation at state load time in `config.LoadState()`

## Test Coverage Gaps for Error Paths

**Issue:** Edge case error conditions not tested

- **Files:** Multiple, but particularly `internal/orchestrator/run.go`, `internal/mcp/server.go`
- **Problem:** Many error paths are not covered by tests (process kill failures, signal channel handling under backpressure, concurrent history writes)
- **Risk:** Untested error paths may behave differently in production under stress
- **Priority:** Medium—low-probability paths, but critical when they do occur
- **Improvement:** Add stress tests with concurrent signal generation, large signal payloads, and process termination races

## Fragile: MCP Server Shutdown Ordering

**Issue:** Server shutdown could lose pending signals

- **Files:** `internal/mcp/server.go` (lines 135-146)
- **Problem:** Shutdown closes signals channel with sync.Once. If signals are still being sent to channel while it's being closed, or if orchestrator is reading from channel while it closes, race condition could occur
- **Impact:** Last-minute insights/task_complete signals could cause panic or be lost
- **Current state:** sync.Once protects the close, but doesn't synchronize with sender goroutines
- **Safe modification:** Add explicit signal to stop sending before closing channel, or use context cancellation to coordinate shutdown
- **Test coverage:** shutdown_test.go exists but doesn't test concurrent sends during shutdown

## Transient Dependencies on External Process

**Issue:** Cannot run if `claude` or `git` commands missing

- **Files:** `internal/process/spawn.go` (line 53), `internal/git/git.go` (line 21)
- **Problem:** Spawning fails if `claude` CLI not in PATH, or if `git` not available. No fallback or graceful degradation. Error messages don't distinguish "command not found" from "command failed"
- **Impact:** Entire orchestrator unusable if dependencies missing. Hard to debug in CI environment
- **Fix approach:** Check for required commands during initialization (RunConfig validation), provide clear error message about missing dependency

---

_Concerns audit: 2026-01-23_
