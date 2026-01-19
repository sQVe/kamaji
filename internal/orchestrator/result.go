package orchestrator

import "github.com/sqve/kamaji/internal/mcp"

// Re-export status constants from mcp for convenience.
const (
	StatusPass = mcp.StatusPass
	StatusFail = mcp.StatusFail
)

// TaskResult normalizes pass/fail/no-signal outcomes from task execution.
// Processes that exit without calling task_complete are treated as failures.
type TaskResult struct {
	Status   string
	Summary  string
	NoSignal bool
}

// PassResult creates a pass result with the given summary.
func PassResult(summary string) TaskResult {
	return TaskResult{
		Status:  StatusPass,
		Summary: summary,
	}
}

// FailResult creates a fail result with the given summary.
func FailResult(summary string) TaskResult {
	return TaskResult{
		Status:  StatusFail,
		Summary: summary,
	}
}

// NoSignalResult creates a fail result for processes that exited without signaling.
func NoSignalResult() TaskResult {
	return TaskResult{
		Status:   StatusFail,
		Summary:  "process exited without signal",
		NoSignal: true,
	}
}

// ResultFromSignal converts an MCP signal to a TaskResult.
// Invalid status values are normalized to fail.
func ResultFromSignal(sig mcp.Signal) TaskResult {
	status := sig.Status
	if status != mcp.StatusPass && status != mcp.StatusFail {
		status = mcp.StatusFail
	}
	return TaskResult{
		Status:  status,
		Summary: sig.Summary,
	}
}

// Passed returns true if the task completed successfully.
func (r TaskResult) Passed() bool {
	return r.Status == StatusPass
}

// Failed returns true if the task failed or exited without signal.
func (r TaskResult) Failed() bool {
	return r.Status == StatusFail
}
