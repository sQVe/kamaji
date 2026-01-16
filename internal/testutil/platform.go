package testutil

import (
	"runtime"
	"testing"
)

func OnWindows() bool {
	return runtime.GOOS == "windows"
}

func OnNonWindows() bool {
	return runtime.GOOS != "windows"
}

func SkipOnWindows(t *testing.T, reason string) {
	t.Helper()
	if OnWindows() {
		t.Skip(reason)
	}
}

func SkipOnNonWindows(t *testing.T, reason string) {
	t.Helper()
	if OnNonWindows() {
		t.Skip(reason)
	}
}

func SkipIfNoUnixPermissions(t *testing.T) {
	t.Helper()
	if OnWindows() {
		t.Skip("Unix-style file permissions not supported on Windows")
	}
}
