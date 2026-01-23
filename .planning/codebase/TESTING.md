# Testing Patterns

**Analysis Date:** 2026-01-23

## Test Framework

**Runner:**

- Go built-in `testing` package
- Test runner: `gotestsum` v1.13.0 (from Makefile)
- Config: `Makefile` specifies test targets with tags and timeouts

**Assertion Library:**

- Standard Go testing: direct value comparison via `if actual != expected`
- Custom testutil helpers: `AssertContains()`, `AssertNotContains()`, `AssertPathEqual()`
- Test assertions formatted: `t.Errorf("message: got %v, want %v", got, want)`

**Run Commands:**

```bash
make test              # Run all unit tests (default target)
make test-unit         # Run unit tests, exclude integration tests
make test-coverage     # Generate coverage report to coverage/coverage.out
make test-integration # Run integration tests with 300s timeout
```

**Coverage:**

- Target: Not explicitly enforced
- Generation: `go test -coverprofile=coverage/coverage.out -covermode=atomic`
- View: `go tool cover -func=coverage/coverage.out`
- Location: `coverage/` directory (gitignored)

## Test File Organization

**Location:**

- Co-located with source: `*_test.go` files in same directory as implementation
- Examples:
    - `internal/git/git_test.go` tests `git.go`
    - `internal/orchestrator/run_test.go` tests `run.go`
    - `internal/config/state_test.go` tests `state.go`

**Naming:**

- Test file suffix: `_test.go`
- Test functions: `Test<FunctionName>_<Scenario>`
- Examples: `TestRun_EmptySprint`, `TestCreateBranch_Success`, `TestLoadState_Valid`
- Scenario names are descriptive: `_Success`, `_Failure`, `_NoChanges`, `_MissingFile_ReturnsZeroValue`

**Structure:**

```
internal/
├── domain/                 # Pure data types
│   └── state.go
│   └── state_test.go       # Tests for State type
├── config/
│   └── state.go            # Config load/save
│   └── state_test.go       # Tests for config operations
├── orchestrator/
│   └── run.go              # Main orchestration
│   └── run_test.go         # Tests for Run function
├── testutil/               # Shared test helpers
│   ├── assert.go           # Assertion helpers
│   ├── git.go              # Git setup helpers
│   └── tempdir.go          # Temp directory helpers
└── git/
    └── git_test.go
```

## Test Structure

**Suite Organization:**

Test functions are flat (no nested suites). Each test is independent:

```go
func TestRun_EmptySprint(t *testing.T) {
	// 1. Setup
	dir := t.TempDir()
	testutil.InitGitRepo(t, dir)

	// 2. Test data
	sprint := &domain.Sprint{
		Name:       "empty",
		BaseBranch: "main",
		Tickets:    []domain.Ticket{},
	}
	sprintPath := writeSprintFile(t, dir, sprint)

	// 3. Execute
	result, err := orchestrator.Run(context.Background(), orchestrator.RunConfig{
		WorkDir:    dir,
		SprintPath: sprintPath,
	})

	// 4. Assert
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !result.Success {
		t.Error("expected Success=true for empty sprint")
	}
	if result.TasksRun != 0 {
		t.Errorf("expected TasksRun=0, got %d", result.TasksRun)
	}
}
```

**Patterns:**

1. **Setup phase:**
    - Create temporary directories: `t.TempDir()`
    - Initialize git repos: `testutil.InitGitRepo(t, dir)`
    - Prepare test data structures

2. **Execution phase:**
    - Call function under test with known inputs
    - Assign return values: `result, err := function(...)`

3. **Assertion phase:**
    - Separate assertions for each aspect
    - Fatal errors for preconditions: `t.Fatalf()`
    - Comparative errors for assertions: `t.Errorf()`

4. **Error assertions:**

    ```go
    if err == nil {
        t.Fatal("expected error")
    }
    if !strings.Contains(err.Error(), "expected text") {
        t.Errorf("error should mention X, got: %v", err)
    }
    ```

5. **Cleanup:**
    - Implicit via `t.TempDir()` (cleaned up automatically)
    - Explicit deferred cleanup when needed: `defer config.ResetPlain()`

## Mocking

**Framework:** Interfaces with concrete mock implementations

**Patterns:**

1. **Interface-based mocking:**

    ```go
    type ProcessSpawner interface {
        Spawn(cfg process.SpawnConfig) (*process.SpawnResult, error)
    }

    type commandSpawner struct {
        cmd string
    }

    func (s commandSpawner) Spawn(cfg process.SpawnConfig) (*process.SpawnResult, error) {
        return process.SpawnCommand(s.cmd, cfg)
    }
    ```

    - Defined in production code, used for dependency injection
    - Implementations in `orchestrator/run.go`: `defaultSpawner`, `commandSpawner`

2. **Test injection:**

    ```go
    result, err := orchestrator.Run(context.Background(), orchestrator.RunConfig{
        WorkDir:    dir,
        SprintPath: sprintPath,
        Spawner:    &mockSpawner{},  // Inject mock
    })
    ```

    - Configuration struct includes optional spawner
    - Tests pass mock implementations via config

3. **Git operations mocking via real git:**
    - No git mocking; tests use real git commands
    - Tests use temporary directories and real git repositories
    - Example from `git_test.go`: calls actual `git` commands to verify behavior

4. **Filesystem isolation:**
    - All file I/O tests use `t.TempDir()`
    - No cleanup needed (automatic)
    - Prevents test pollution and side effects

**What to mock:**

- External processes via `ProcessSpawner` interface
- Configuration systems via optional config fields
- Platform-specific code via abstraction interfaces

**What NOT to mock:**

- Standard library functions (os, io packages)
- Git operations (use real git in temp directories)
- YAML/JSON marshaling (test real serialization)
- File system (use `t.TempDir()` for isolation)

## Fixtures and Factories

**Test Data:**

Helper functions create test fixtures:

```go
func writeSprintFile(t *testing.T, dir string, sprint *domain.Sprint) string {
	t.Helper()
	path := filepath.Join(dir, "kamaji.yaml")
	data, err := yaml.Marshal(sprint)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}
```

**Factory pattern:**

- Helper functions take `*testing.T` as first parameter
- Mark as helper: `t.Helper()` (improves error reporting)
- Return prepared fixtures or error on setup failure
- Example: `testutil.InitGitRepo(t, dir, branches...)` prepares git repos

**Location:**

- Helpers in test files themselves: `run_test.go` contains `writeSprintFile()`
- Shared helpers in `internal/testutil/`:
    - `git.go`: `InitGitRepo()` for setting up git repositories
    - `assert.go`: `AssertContains()`, `AssertNotContains()`, `AssertPathEqual()`
    - `tempdir.go`: Temporary directory utilities

## Coverage

**Requirements:** No explicit minimum enforced

**View Coverage:**

```bash
make test-coverage    # Generate coverage/coverage.out
go tool cover -func=coverage/coverage.out  # Show function coverage
go tool cover -html=coverage/coverage.out  # Open in browser
```

**Strategy:**

- Focus on behavior coverage, not line coverage
- Test error paths and edge cases
- Example: `TestCreateBranch_ExistingTicketBranch` tests error scenario

## Test Types

**Unit Tests:**

- Scope: Single function or method
- Approach: Fast, isolated, deterministic
- Examples: `TestLoadState_Valid`, `TestCreateBranch_Success`
- Typical: 30-50 tests across codebase

**Integration Tests:**

- Scope: Multiple components working together
- Approach: Real git operations, real file I/O, real MCP server
- Tag: `integration` (can be skipped with `-tags=!integration`)
- Location: Primarily in `cmd/kamaji/` per Makefile
- Command: `make test-integration` with 300s timeout

**E2E Tests:**

- Not present in current codebase
- Integration tests serve this role

## Common Patterns

**Async Testing:**

Channels and goroutines:

```go
done := make(chan struct{})
go func() {
	_ = spawnResult.Process.Wait()
	close(done)
}()

select {
case <-ctx.Done():
	return &RunResult{}, ctx.Err()
case sig, ok := <-tc.server.Signals():
	// handle signal
case <-done:
	// process exited
}
```

- Use channels with select for synchronization
- Close channels to signal completion
- Example from `orchestrator/run.go` runTask function

**Context Testing:**

```go
ctx, cancel := context.WithCancel(context.Background())
cancel()

_, err := orchestrator.Run(ctx, orchestrator.RunConfig{...})
if err == nil {
	t.Fatal("expected error for cancelled context")
}
```

- Pre-cancel context to test early termination
- Verify context.Err() in error message
- Example: `TestRun_ContextCancellation_BeforeTaskStart`

**Error Testing:**

Multiple assertion styles:

```go
// 1. Nil check
if err == nil {
	t.Error("expected error")
}

// 2. Error contains text
if !strings.Contains(err.Error(), "expected phrase") {
	t.Errorf("got: %v", err)
}

// 3. Error is sentinel (errors.Is)
if errors.Is(err, git.ErrNothingToCommit) {
	// expected
}
```

**Roundtrip Testing:**

Save and load verification:

```go
func TestSaveState_Roundtrip(t *testing.T) {
	original := &domain.State{
		CurrentTicket: 3,
		CurrentTask:   7,
		FailureCount:  2,
	}

	if err := SaveState(dir, original); err != nil {
		t.Fatalf("SaveState error: %v", err)
	}

	loaded, err := LoadState(dir)
	if err != nil {
		t.Fatalf("LoadState error: %v", err)
	}

	// Compare fields
	if loaded.CurrentTicket != original.CurrentTicket {
		t.Errorf("CurrentTicket: got %d, want %d", loaded.CurrentTicket, original.CurrentTicket)
	}
}
```

- Verify data survives serialization round-trip
- Used for config and state tests

**Setup/Teardown:**

Deferred cleanup:

```go
config.SetPlain(true)
defer config.ResetPlain()
```

- Use defer for cleanup
- Test setup in function body, not separate methods
- Example: `output/writer_test.go` tests manage plain mode state

## Test Coverage Insights

**Well-tested areas:**

- `internal/git/`: Full coverage of git operations (9 tests)
- `internal/config/`: Load/save operations tested comprehensively
- `internal/orchestrator/`: Run function lifecycle tested
- `internal/output/`: Writer and styling tested

**Common test scenarios:**

- Success paths
- Error scenarios (missing files, invalid input, non-git directories)
- Edge cases (empty sprints, no changes)
- State transitions and rollbacks

---

_Testing analysis: 2026-01-23_
