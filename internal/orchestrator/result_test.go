package orchestrator_test

import (
	"testing"

	"github.com/sqve/kamaji/internal/mcp"
	"github.com/sqve/kamaji/internal/orchestrator"
)

func TestPassResult_CreatesPassStatus(t *testing.T) {
	result := orchestrator.PassResult("task completed successfully")

	if result.Status != orchestrator.StatusPass {
		t.Errorf("Status = %q, want %q", result.Status, orchestrator.StatusPass)
	}
	if result.Summary != "task completed successfully" {
		t.Errorf("Summary = %q, want %q", result.Summary, "task completed successfully")
	}
	if result.NoSignal {
		t.Error("NoSignal = true, want false")
	}
}

func TestPassResult_PassedReturnsTrue(t *testing.T) {
	result := orchestrator.PassResult("done")

	if !result.Passed() {
		t.Error("Passed() = false, want true")
	}
	if result.Failed() {
		t.Error("Failed() = true, want false")
	}
}

func TestFailResult_CreatesFailStatus(t *testing.T) {
	result := orchestrator.FailResult("validation error")

	if result.Status != orchestrator.StatusFail {
		t.Errorf("Status = %q, want %q", result.Status, orchestrator.StatusFail)
	}
	if result.Summary != "validation error" {
		t.Errorf("Summary = %q, want %q", result.Summary, "validation error")
	}
	if result.NoSignal {
		t.Error("NoSignal = true, want false")
	}
}

func TestFailResult_FailedReturnsTrue(t *testing.T) {
	result := orchestrator.FailResult("error")

	if result.Passed() {
		t.Error("Passed() = true, want false")
	}
	if !result.Failed() {
		t.Error("Failed() = false, want true")
	}
}

func TestNoSignalResult_CreatesFailWithNoSignalFlag(t *testing.T) {
	result := orchestrator.NoSignalResult()

	if result.Status != orchestrator.StatusFail {
		t.Errorf("Status = %q, want %q", result.Status, orchestrator.StatusFail)
	}
	if !result.NoSignal {
		t.Error("NoSignal = false, want true")
	}
	if result.Summary == "" {
		t.Error("Summary should not be empty")
	}
	if result.Summary != "process exited without signal" {
		t.Errorf("Summary = %q, want %q", result.Summary, "process exited without signal")
	}
}

func TestNoSignalResult_FailedReturnsTrue(t *testing.T) {
	result := orchestrator.NoSignalResult()

	if result.Passed() {
		t.Error("Passed() = true, want false")
	}
	if !result.Failed() {
		t.Error("Failed() = false, want true")
	}
}

func TestResultFromSignal_PassSignal(t *testing.T) {
	signal := mcp.Signal{
		Tool:    mcp.SignalToolTaskComplete,
		Status:  mcp.StatusPass,
		Summary: "all tests passed",
	}

	result := orchestrator.ResultFromSignal(signal)

	if result.Status != orchestrator.StatusPass {
		t.Errorf("Status = %q, want %q", result.Status, orchestrator.StatusPass)
	}
	if result.Summary != "all tests passed" {
		t.Errorf("Summary = %q, want %q", result.Summary, "all tests passed")
	}
	if result.NoSignal {
		t.Error("NoSignal = true, want false")
	}
	if !result.Passed() {
		t.Error("Passed() = false, want true")
	}
}

func TestResultFromSignal_FailSignal(t *testing.T) {
	signal := mcp.Signal{
		Tool:    mcp.SignalToolTaskComplete,
		Status:  mcp.StatusFail,
		Summary: "lint errors found",
	}

	result := orchestrator.ResultFromSignal(signal)

	if result.Status != orchestrator.StatusFail {
		t.Errorf("Status = %q, want %q", result.Status, orchestrator.StatusFail)
	}
	if result.Summary != "lint errors found" {
		t.Errorf("Summary = %q, want %q", result.Summary, "lint errors found")
	}
	if result.NoSignal {
		t.Error("NoSignal = true, want false")
	}
	if !result.Failed() {
		t.Error("Failed() = false, want true")
	}
}

func TestResultFromSignal_InvalidStatusNormalizesToFail(t *testing.T) {
	signal := mcp.Signal{
		Tool:    mcp.SignalToolTaskComplete,
		Status:  "invalid",
		Summary: "unknown status",
	}

	result := orchestrator.ResultFromSignal(signal)

	if result.Status != orchestrator.StatusFail {
		t.Errorf("Status = %q, want %q (invalid status should normalize to fail)", result.Status, orchestrator.StatusFail)
	}
	if !result.Failed() {
		t.Error("Failed() = false, want true for invalid status")
	}
}
