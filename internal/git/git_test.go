package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// initTestRepo creates a git repo in dir with an initial commit and optional branches.
func initTestRepo(t *testing.T, dir string, branches ...string) {
	t.Helper()

	run := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test",
			"GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=Test",
			"GIT_COMMITTER_EMAIL=test@test.com",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v failed: %v\n%s", args, err, out)
		}
	}

	run("init", "-b", "main")
	run("config", "user.email", "test@test.com")
	run("config", "user.name", "Test")
	run("config", "commit.gpgsign", "false")
	run("config", "core.autocrlf", "false")

	// Create initial commit so we have a valid HEAD
	readme := filepath.Join(dir, "README.md")
	if err := os.WriteFile(readme, []byte("# Test\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-m", "initial commit")

	// Create additional branches
	for _, branch := range branches {
		run("branch", branch)
	}
}

func TestCreateBranch_Success(t *testing.T) {
	dir := t.TempDir()
	initTestRepo(t, dir)

	err := CreateBranch(dir, "main", "feature/test-123")
	if err != nil {
		t.Errorf("CreateBranch() error = %v, want nil", err)
	}

	// Verify we're on the new branch
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git branch --show-current failed: %v", err)
	}
	branch := strings.TrimSpace(string(out))
	if branch != "feature/test-123" {
		t.Errorf("current branch = %q, want %q", branch, "feature/test-123")
	}
}

func TestCreateBranch_MissingBaseBranch(t *testing.T) {
	dir := t.TempDir()
	initTestRepo(t, dir)

	err := CreateBranch(dir, "nonexistent-branch", "feature/test-123")
	if err == nil {
		t.Error("CreateBranch() error = nil, want error for missing base branch")
	}
	if !strings.Contains(err.Error(), "nonexistent-branch") {
		t.Errorf("error should mention branch name, got: %v", err)
	}
}

func TestCreateBranch_ExistingTicketBranch(t *testing.T) {
	dir := t.TempDir()
	initTestRepo(t, dir, "feature/existing")

	err := CreateBranch(dir, "main", "feature/existing")
	if err == nil {
		t.Error("CreateBranch() error = nil, want error for existing ticket branch")
	}
	if !strings.Contains(err.Error(), "feature/existing") {
		t.Errorf("error should mention branch name, got: %v", err)
	}
}

func TestCreateBranch_NotAGitRepo(t *testing.T) {
	dir := t.TempDir()
	// Don't init git repo

	err := CreateBranch(dir, "main", "feature/test-123")
	if err == nil {
		t.Error("CreateBranch() error = nil, want error for non-git directory")
	}
}

func TestCreateBranch_FromDifferentBranch(t *testing.T) {
	dir := t.TempDir()
	initTestRepo(t, dir, "develop")

	// Start on a different branch
	cmd := exec.Command("git", "checkout", "develop")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git checkout develop failed: %v\n%s", err, out)
	}

	// Create branch from main while on develop
	err := CreateBranch(dir, "main", "feature/from-main")
	if err != nil {
		t.Errorf("CreateBranch() error = %v, want nil", err)
	}

	// Verify we're on the new branch
	cmd = exec.Command("git", "branch", "--show-current")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git branch --show-current failed: %v", err)
	}
	branch := strings.TrimSpace(string(out))
	if branch != "feature/from-main" {
		t.Errorf("current branch = %q, want %q", branch, "feature/from-main")
	}
}

// CommitChanges tests

func TestCommitChanges_Success(t *testing.T) {
	dir := t.TempDir()
	initTestRepo(t, dir)

	newFile := filepath.Join(dir, "new.txt")
	if err := os.WriteFile(newFile, []byte("new content\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := CommitChanges(dir, "test: add new file")
	if err != nil {
		t.Errorf("CommitChanges() error = %v, want nil", err)
	}

	// Verify commit was created
	cmd := exec.Command("git", "log", "--oneline", "-1")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git log failed: %v", err)
	}
	if !strings.Contains(string(out), "test: add new file") {
		t.Errorf("commit message not found in log: %s", out)
	}
}

func TestCommitChanges_NoChanges(t *testing.T) {
	dir := t.TempDir()
	initTestRepo(t, dir)

	err := CommitChanges(dir, "test: nothing to commit")
	if err == nil {
		t.Error("CommitChanges() error = nil, want error for nothing to commit")
	}
	if !strings.Contains(err.Error(), "nothing to commit") {
		t.Errorf("error should mention nothing to commit, got: %v", err)
	}
}

func TestCommitChanges_EmptyMessage(t *testing.T) {
	dir := t.TempDir()
	initTestRepo(t, dir)

	// Create a file so there's something to commit
	newFile := filepath.Join(dir, "new.txt")
	if err := os.WriteFile(newFile, []byte("content\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := CommitChanges(dir, "")
	if err == nil {
		t.Error("CommitChanges() error = nil, want error for empty message")
	}
	if !strings.Contains(err.Error(), "commit message required") {
		t.Errorf("error should mention commit message required, got: %v", err)
	}
}

func TestCommitChanges_NotAGitRepo(t *testing.T) {
	dir := t.TempDir()

	err := CommitChanges(dir, "test: should fail")
	if err == nil {
		t.Error("CommitChanges() error = nil, want error for non-git directory")
	}
}

// ResetToHead tests

func TestResetToHead_Success(t *testing.T) {
	dir := t.TempDir()
	initTestRepo(t, dir)

	// Create uncommitted changes to tracked file
	readme := filepath.Join(dir, "README.md")
	if err := os.WriteFile(readme, []byte("modified content\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Create untracked file
	untracked := filepath.Join(dir, "untracked.txt")
	if err := os.WriteFile(untracked, []byte("untracked\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := ResetToHead(dir)
	if err != nil {
		t.Errorf("ResetToHead() error = %v, want nil", err)
	}

	// Verify tracked file changes were discarded
	content, err := os.ReadFile(readme) //nolint:gosec // test code with temp dir
	if err != nil {
		t.Fatalf("failed to read README.md: %v", err)
	}
	if string(content) != "# Test\n" {
		t.Errorf("README.md content = %q, want %q", content, "# Test\n")
	}

	// Verify untracked file was removed
	if _, err := os.Stat(untracked); !os.IsNotExist(err) {
		t.Error("untracked file should have been removed by ResetToHead")
	}
}

func TestResetToHead_NoChanges(t *testing.T) {
	dir := t.TempDir()
	initTestRepo(t, dir)

	// Reset with no changes should be idempotent
	err := ResetToHead(dir)
	if err != nil {
		t.Errorf("ResetToHead() error = %v, want nil (idempotent)", err)
	}
}

func TestResetToHead_NotAGitRepo(t *testing.T) {
	dir := t.TempDir()

	err := ResetToHead(dir)
	if err == nil {
		t.Error("ResetToHead() error = nil, want error for non-git directory")
	}
}
