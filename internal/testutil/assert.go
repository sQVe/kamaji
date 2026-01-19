package testutil

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func AssertContains(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("expected %q to contain %q", truncate(s, 200), substr)
	}
}

func AssertNotContains(t *testing.T, s, substr string) {
	t.Helper()
	if strings.Contains(s, substr) {
		t.Errorf("expected %q to NOT contain %q", truncate(s, 200), substr)
	}
}

func truncate(s string, limit int) string {
	if len(s) <= limit {
		return s
	}
	return s[:limit] + "..."
}

func AssertPathEqual(t *testing.T, got, want string) {
	t.Helper()
	gotClean := filepath.Clean(got)
	wantClean := filepath.Clean(want)
	if gotClean != wantClean {
		t.Errorf("path = %q, want %q", gotClean, wantClean)
	}
}

var outputMu sync.Mutex

// CaptureStdout captures stdout during fn execution and returns the output.
func CaptureStdout(t *testing.T, fn func()) string {
	return captureOutput(t, &os.Stdout, fn)
}

// CaptureStderr captures stderr during fn execution and returns the output.
func CaptureStderr(t *testing.T, fn func()) string {
	return captureOutput(t, &os.Stderr, fn)
}

func captureOutput(t *testing.T, target **os.File, fn func()) string {
	t.Helper()
	outputMu.Lock()
	defer outputMu.Unlock()

	old := *target
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	*target = w
	defer func() {
		*target = old
	}()

	fn()

	_ = w.Close()

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()

	return buf.String()
}
