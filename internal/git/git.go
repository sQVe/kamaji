package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

// runGit executes a git command in the specified directory.
//
//nolint:unparam // stdout will be used by future operations
func runGit(workDir string, args ...string) (stdout, stderr string, err error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = workDir

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

// CreateBranch creates a new branch from the base branch.
// It checks out the base branch, pulls latest (continues if offline), then creates the ticket branch.
func CreateBranch(workDir, baseBranch, ticketBranch string) error {
	if workDir == "" {
		return errors.New("workDir required")
	}
	if baseBranch == "" {
		return errors.New("baseBranch required")
	}
	if ticketBranch == "" {
		return errors.New("ticketBranch required")
	}

	_, stderr, err := runGit(workDir, "checkout", baseBranch)
	if err != nil {
		return fmt.Errorf("checkout %s (%s): %w", baseBranch, stderr, err)
	}

	// Pull latest - continue even if this fails (offline scenario)
	_, _, _ = runGit(workDir, "pull", "origin", baseBranch)

	_, stderr, err = runGit(workDir, "checkout", "-b", ticketBranch)
	if err != nil {
		return fmt.Errorf("create branch %s (%s): %w", ticketBranch, stderr, err)
	}

	return nil
}

// CommitChanges stages all changes and commits with the provided message.
func CommitChanges(workDir, message string) error {
	if workDir == "" {
		return errors.New("workDir required")
	}
	if message == "" {
		return errors.New("commit message required")
	}

	_, stderr, err := runGit(workDir, "add", "-A")
	if err != nil {
		return fmt.Errorf("git add (%s): %w", stderr, err)
	}

	// Check if anything is staged
	_, _, err = runGit(workDir, "diff", "--cached", "--quiet")
	if err == nil {
		// Exit 0 means no differences (nothing staged)
		return errors.New("nothing to commit")
	}

	_, stderr, err = runGit(workDir, "commit", "-m", message)
	if err != nil {
		return fmt.Errorf("git commit (%s): %w", stderr, err)
	}

	return nil
}

// ResetToHead discards all uncommitted changes and removes untracked files.
func ResetToHead(workDir string) error {
	if workDir == "" {
		return errors.New("workDir required")
	}

	_, stderr, err := runGit(workDir, "reset", "--hard", "HEAD")
	if err != nil {
		return fmt.Errorf("git reset (%s): %w", stderr, err)
	}

	// Remove untracked files and directories
	_, stderr, err = runGit(workDir, "clean", "-fd")
	if err != nil {
		return fmt.Errorf("git clean (%s): %w", stderr, err)
	}

	return nil
}
