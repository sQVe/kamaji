package version

import "testing"

func TestFull_DevVersion(t *testing.T) {
	// Save original values
	origVersion, origCommit, origDate := Version, Commit, Date
	t.Cleanup(func() {
		Version, Commit, Date = origVersion, origCommit, origDate
	})

	Version = "dev"
	Commit = "unknown"
	Date = "unknown"

	got := Full()
	want := "dev"
	if got != want {
		t.Errorf("Full() = %q, want %q", got, want)
	}
}

func TestFull_ReleaseVersion(t *testing.T) {
	// Save original values
	origVersion, origCommit, origDate := Version, Commit, Date
	t.Cleanup(func() {
		Version, Commit, Date = origVersion, origCommit, origDate
	})

	Version = "v1.0.0"
	Commit = "abc1234"
	Date = "2025-01-14"

	got := Full()
	want := "v1.0.0 (abc1234, 2025-01-14)"
	if got != want {
		t.Errorf("Full() = %q, want %q", got, want)
	}
}
