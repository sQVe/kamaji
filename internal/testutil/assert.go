package testutil

import (
	"path/filepath"
	"strings"
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
		t.Errorf("path = %q, want %q", got, want)
	}
}
