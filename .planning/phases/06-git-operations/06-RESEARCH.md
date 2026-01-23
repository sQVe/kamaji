# Phase 6: Git Operations - Research

**Researched:** 2026-01-15
**Domain:** Go git command execution via os/exec
**Confidence:** HIGH

<research_summary>

## Summary

Researched git operations implementation patterns for Go using `os/exec` (not go-git). The standard approach uses a `runGit` helper that captures stdout/stderr, sets working directory, and wraps errors with stderr context. Testing uses real git commands in temp directories with test helper utilities.

Key finding: Grove's implementation at `/home/sqve/code/personal/grove/main/internal/git/` provides excellent reference patterns. Their `testutil/git` package with `TestRepo` struct is directly applicable.

**Primary recommendation:** Follow Grove's patterns - create a simple git package with `runGit` helper, use `t.TempDir()` + test helper struct for testing with real git commands.
</research_summary>

<standard_stack>

## Standard Stack

### Core

| Library      | Version | Purpose                 | Why Standard                  |
| ------------ | ------- | ----------------------- | ----------------------------- |
| os/exec      | stdlib  | Git command execution   | No dependencies, full control |
| bytes.Buffer | stdlib  | Capture stdout/stderr   | Standard pattern for exec     |
| errors       | stdlib  | Error wrapping/checking | errors.As for ExitError       |

### Supporting

| Library     | Version | Purpose             | When to Use    |
| ----------- | ------- | ------------------- | -------------- |
| testing     | stdlib  | Test framework      | All tests      |
| t.TempDir() | stdlib  | Isolated test repos | Every git test |

### Alternatives Considered

| Instead of     | Could Use | Tradeoff                                                             |
| -------------- | --------- | -------------------------------------------------------------------- |
| os/exec        | go-git    | go-git adds dependency, more abstraction, but plan specifies os/exec |
| Real git tests | Mocks     | Mocks don't catch real git behavior edge cases                       |

**Installation:**

```bash
# No installation needed - stdlib only
```

</standard_stack>

<architecture_patterns>

## Architecture Patterns

### Recommended Project Structure

```
internal/
├── git/
│   ├── git.go        # runGit helper + public functions
│   └── git_test.go   # Tests with temp repo fixtures
```

### Pattern 1: runGit Helper

**What:** Central function for all git command execution
**When to use:** Every git operation
**Example (from Grove):**

```go
func runGit(workDir string, args ...string) (stdout, stderr string, err error) {
    cmd := exec.Command("git", args...)
    cmd.Dir = workDir

    var stdoutBuf, stderrBuf bytes.Buffer
    cmd.Stdout = &stdoutBuf
    cmd.Stderr = &stderrBuf

    err = cmd.Run()
    stdout = strings.TrimSpace(stdoutBuf.String())
    stderr = strings.TrimSpace(stderrBuf.String())

    if err != nil {
        if stderr != "" {
            return stdout, stderr, fmt.Errorf("git %s: %w: %s", args[0], err, stderr)
        }
        return stdout, stderr, fmt.Errorf("git %s: %w", args[0], err)
    }
    return stdout, stderr, nil
}
```

### Pattern 2: Test Repository Helper

**What:** Struct that creates isolated git repos for testing
**When to use:** All git operation tests
**Example (adapted from Grove):**

```go
type testRepo struct {
    t    *testing.T
    path string
}

func newTestRepo(t *testing.T) *testRepo {
    t.Helper()
    dir := t.TempDir()

    // git init
    cmd := exec.Command("git", "init")
    cmd.Dir = dir
    if err := cmd.Run(); err != nil {
        t.Fatalf("git init: %v", err)
    }

    // Configure git for tests
    for _, cfg := range [][]string{
        {"user.email", "test@example.com"},
        {"user.name", "Test User"},
        {"commit.gpgsign", "false"},
    } {
        cmd := exec.Command("git", "config", cfg[0], cfg[1])
        cmd.Dir = dir
        if err := cmd.Run(); err != nil {
            t.Fatalf("git config %s: %v", cfg[0], err)
        }
    }

    return &testRepo{t: t, path: dir}
}

func (r *testRepo) writeFile(name, content string) {
    r.t.Helper()
    if err := os.WriteFile(filepath.Join(r.path, name), []byte(content), 0644); err != nil {
        r.t.Fatalf("write file: %v", err)
    }
}
```

### Anti-Patterns to Avoid

- **Mocking git:** Real git in temp dirs catches edge cases mocks miss
- **Not capturing stderr:** Errors without stderr context are hard to debug
- **Using cmd.Output():** ExitError.Stderr can be truncated; use explicit buffers
- **Forgetting git config:** Tests fail without user.email/user.name configured
  </architecture_patterns>

<dont_hand_roll>

## Don't Hand-Roll

| Problem               | Don't Build              | Use Instead                         | Why                             |
| --------------------- | ------------------------ | ----------------------------------- | ------------------------------- |
| Stdout/stderr capture | Custom pipes             | bytes.Buffer with cmd.Stdout/Stderr | Simpler, avoids race conditions |
| Exit code checking    | Manual exit code parsing | errors.As with \*exec.ExitError     | Standard Go pattern             |
| Test repo setup       | Ad-hoc git commands      | testRepo helper struct              | DRY, consistent config          |

**Key insight:** Git command execution is straightforward with os/exec. The complexity is in error handling (including stderr) and test isolation (temp dirs with proper git config).
</dont_hand_roll>

<common_pitfalls>

## Common Pitfalls

### Pitfall 1: Missing git config in tests

**What goes wrong:** Tests fail with "please tell me who you are"
**Why it happens:** Git requires user.email and user.name for commits
**How to avoid:** Always configure in test setup: user.email, user.name, commit.gpgsign=false
**Warning signs:** Flaky tests, "Author identity unknown" errors

### Pitfall 2: Truncated stderr in errors

**What goes wrong:** Error messages lack useful git output
**Why it happens:** Using cmd.Output() truncates ExitError.Stderr
**How to avoid:** Use explicit bytes.Buffer for stderr, include in error wrapping
**Warning signs:** Errors like "exit status 1" with no context

### Pitfall 3: Race conditions with StdoutPipe/StderrPipe

**What goes wrong:** Sporadic test failures, race detector warnings
**Why it happens:** Reading from pipes while calling Wait() is incorrect
**How to avoid:** Use bytes.Buffer assigned to cmd.Stdout/Stderr instead
**Warning signs:** `go test -race` failures

### Pitfall 4: Pull failures blocking branch creation

**What goes wrong:** Offline/no-remote scenarios fail entirely
**Why it happens:** Treating pull failure as fatal
**How to avoid:** Log warning on pull fail, continue with local state (per plan spec)
**Warning signs:** CreateBranch fails when working offline
</common_pitfalls>

<code_examples>

## Code Examples

### CreateBranch (from plan 06-01)

```go
// Source: Plan specification + Grove patterns
func CreateBranch(workDir, baseBranch, ticketBranch string) error {
    // Checkout base branch
    if _, _, err := runGit(workDir, "checkout", baseBranch); err != nil {
        return fmt.Errorf("checkout %s: %w", baseBranch, err)
    }

    // Pull latest (warning if fails, not error)
    if _, stderr, err := runGit(workDir, "pull", "origin", baseBranch); err != nil {
        // Log warning but continue - allows offline work
        log.Printf("warning: pull failed (continuing): %s", stderr)
    }

    // Create ticket branch
    if _, _, err := runGit(workDir, "checkout", "-b", ticketBranch); err != nil {
        return fmt.Errorf("create branch %s: %w", ticketBranch, err)
    }

    return nil
}
```

### CommitChanges (from plan 06-02)

```go
// Source: Plan specification + Grove patterns
func CommitChanges(workDir, message string) error {
    if message == "" {
        return errors.New("commit message required")
    }

    // Stage all changes
    if _, _, err := runGit(workDir, "add", "-A"); err != nil {
        return err
    }

    // Check if anything staged
    _, _, err := runGit(workDir, "diff", "--cached", "--quiet")
    if err == nil {
        // Exit 0 means nothing staged
        return errors.New("nothing to commit")
    }
    // Exit 1 means changes exist - continue

    // Commit
    if _, _, err := runGit(workDir, "commit", "-m", message); err != nil {
        return err
    }

    return nil
}
```

### ResetToHead (from plan 06-02)

```go
// Source: Plan specification
func ResetToHead(workDir string) error {
    _, _, err := runGit(workDir, "reset", "--hard", "HEAD")
    return err
}
```

</code_examples>

<sota_updates>

## State of the Art (2025-2026)

| Old Approach     | Current Approach        | When Changed | Impact                                             |
| ---------------- | ----------------------- | ------------ | -------------------------------------------------- |
| go-git library   | os/exec for simple ops  | -            | Less dependency, more control for simple use cases |
| Mocked git tests | Real git in t.TempDir() | -            | Better test fidelity                               |

**New tools/patterns to consider:**

- t.TempDir() (Go 1.15+): Auto-cleanup, no manual defer needed
- errors.As (Go 1.13+): Proper error unwrapping for ExitError

**Deprecated/outdated:**

- StdoutPipe/StderrPipe for simple capture: Use bytes.Buffer instead
  </sota_updates>

<open_questions>

## Open Questions

None - this domain is well-understood and plans are already detailed.
</open_questions>

<sources>
## Sources

### Primary (HIGH confidence)

- Grove implementation at /home/sqve/code/personal/grove/main/internal/git/ - Complete reference implementation
- Grove testutil at /home/sqve/code/personal/grove/main/internal/testutil/git/git.go - Test helper patterns
- https://pkg.go.dev/os/exec - Official Go documentation

### Secondary (MEDIUM confidence)

- https://www.dolthub.com/blog/2022-11-28-go-os-exec-patterns/ - Verified patterns
- https://gopheradvent.com/calendar/2021/gotchas-in-exec-errors/ - ExitError.Stderr caveat
  </sources>

<metadata>
## Metadata

**Research scope:**

- Core technology: os/exec for git command execution
- Ecosystem: Standard library only
- Patterns: runGit helper, testRepo fixture
- Pitfalls: stderr capture, git config in tests

**Confidence breakdown:**

- Standard stack: HIGH - stdlib only, well-documented
- Architecture: HIGH - directly from Grove reference
- Pitfalls: HIGH - documented in official docs and verified
- Code examples: HIGH - from plan specs + Grove patterns

**Research date:** 2026-01-15
**Valid until:** 2026-02-15 (30 days - stable patterns)
</metadata>

---

_Phase: 06-git-operations_
_Research completed: 2026-01-15_
_Ready for planning: yes (plans already exist)_
