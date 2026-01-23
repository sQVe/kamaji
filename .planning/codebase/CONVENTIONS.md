# Coding Conventions

**Analysis Date:** 2026-01-23

## Naming Patterns

**Files:**

- Lowercase with underscores for multi-word names
- Test files: `*_test.go` suffix (e.g., `git_test.go`, `state_test.go`)
- Package files correspond to their package name
- Examples: `git.go`, `statemachine.go`, `handler.go`, `progress_test.go`

**Functions:**

- PascalCase for exported functions and methods
- camelCase for unexported functions
- Descriptive names indicating action: `LoadState`, `SaveState`, `CommitChanges`, `CreateBranch`
- Test functions: `Test<FunctionName>_<Scenario>` (e.g., `TestRun_EmptySprint`, `TestCreateBranch_Success`)

**Variables:**

- camelCase for local variables
- PascalCase for exported package-level variables and constants
- Short names acceptable in tight scopes (e.g., `t` for `*testing.T`, `cfg` for config)
- State machine uses `prev*` prefix for rollback tracking: `prevTicket`, `prevTask`, `prevFailures`
- Private struct fields use camelCase: `workDir`, `state`, `sprint`, `spawner`

**Types:**

- PascalCase for all type names (struct, interface, custom types)
- Config structs end with "Config": `RunConfig`, `SpawnConfig`
- Result structs end with "Result": `RunResult`, `SpawnResult`, `TaskResult`
- Handler structs end with "Handler": `Handler`
- Interface names typically singular without "I" prefix: `ProcessSpawner` not `IProcessSpawner`

**Constants:**

- UPPER_SNAKE_CASE for exported constants
- Examples: `StuckThreshold`, `ErrNothingToCommit`, `ErrBranchExists`
- Error sentinels prefixed with `Err`: `var ErrBranchExists = errors.New("branch already exists")`

## Code Style

**Formatting:**

- gofumpt: Formatter for strict Go style with extra rules enabled
- goimports: Automatic import organization
- Tab indentation for Go files (indent_size=4)
- Space indentation for other files (indent_size=2)
- EditorConfig enforced via `.editorconfig`

**Linting:**

- golangci-lint v2 with comprehensive rulesets
- Enabled linters: `errcheck`, `govet`, `staticcheck`, `unused`, `bodyclose`, `goconst`, `gocritic`, `unparam`, `unconvert`, `wastedassign`, `gosec`, `errorlint`, `nilerr`, `nilnil`, `contextcheck`, `testifylint`
- File-specific exceptions: `_test.go` files exclude `goconst` linter (constants are allowed in tests)
- Max issues per linter: 0 (fail on any violations)

**Line width:** No specific limit enforced; prioritize readability

## Import Organization

**Order:**

1. Standard library imports (`fmt`, `os`, `io`, etc.)
2. Third-party imports (`github.com/...`, `gopkg.in/...`, etc.)
3. Local package imports from same module (`github.com/sqve/kamaji/internal/...`)
4. Blank line between groups

**Example from `orchestrator/run.go`:**

```go
import (
	"context"
	"errors"
	"os"

	"github.com/sqve/kamaji/internal/config"
	"github.com/sqve/kamaji/internal/domain"
	"github.com/sqve/kamaji/internal/git"
	"github.com/sqve/kamaji/internal/mcp"
	"github.com/sqve/kamaji/internal/output"
	"github.com/sqve/kamaji/internal/process"
	"github.com/sqve/kamaji/internal/prompt"
	"github.com/sqve/kamaji/internal/statemachine"
)
```

**Path aliases:** Not used; full import paths preferred for clarity

## Error Handling

**Patterns:**

1. **Error wrapping with context:**

    ```go
    if err != nil {
    	return fmt.Errorf("operation context: %w", err)
    }
    ```

    - All errors wrapped with `%w` for error chain preservation
    - Context included describing what failed
    - Example from `config/state.go`: `fmt.Errorf("reading state file: %w", err)`

2. **Error sentinel variables:**

    ```go
    var ErrNothingToCommit = errors.New("nothing to commit")
    var ErrBranchExists = errors.New("branch already exists")
    ```

    - Exported error sentinels in package scope
    - Used with `errors.Is()` for type checking
    - Example from `orchestrator/handler.go`: `if errors.Is(err, git.ErrNothingToCommit)`

3. **Early returns:**
    - Functions return errors immediately after operations
    - No error accumulation or suppression
    - Single path for success case when possible

4. **Argument validation:**

    ```go
    if workDir == "" {
    	return errors.New("workDir required")
    }
    ```

    - Validate required parameters at function entry
    - Use simple `errors.New()` for validation errors
    - Examples from `git/git.go` and `orchestrator/run.go`

5. **Selective error handling:**
    - Some errors are expected (e.g., offline git pulls)
    - Comment explains when errors are ignored: `// Pull latest - continue even if this fails (offline scenario)`
    - Blank assignment used: `_, _, _ = runGit(...)`

6. **State rollback on persistence failure:**

    ```go
    if err := config.SaveState(h.workDir, h.state); err != nil {
    	h.state.CurrentTicket = prevTicket
    	h.state.CurrentTask = prevTask
    	h.state.FailureCount = prevFailures
    	return err
    }
    ```

    - Save original state before modification
    - Restore if persistence fails
    - Used in `orchestrator/handler.go` OnPass and OnFail methods

## Logging

**Framework:** stdlib `fmt` with output package helpers

**Patterns:**

- No structured logging used
- Output package provides styled messages: `PrintInfo()`, `PrintError()`, `PrintSuccess()`
- Signal logging for task completion: `output.PrintSignal(mcp.Signal{...})`
- Debug output via `output.PrintDebug()` (controlled by plain mode)
- All output to stdout/stderr via explicit `fmt.Fprintln()`

**Example from `output/styles.go`:**

```go
func PrintSuccess(msg string) {
	_, _ = fmt.Fprintln(os.Stdout, Style(Success, msg))
}
```

## Comments

**When to comment:**

- Explain WHY, not WHAT (code shows WHAT)
- Document non-obvious state transitions
- Explain domain rules and thresholds
- Clarify workarounds and trade-offs

**Examples:**

- `// Handler does not own state or sprint; the caller owns these values` (lifecycle ownership)
- `// Fresh context needed since caller context may be cancelled` (context lifecycle)
- `// git rev-parse exits non-zero when the ref doesn't exist, which is expected.` (external tool behavior)
- `// Pull latest - continue even if this fails (offline scenario)` (error tolerance rationale)

**Doc comments:**

- Exported functions use doc comments
- Format: `// FunctionName does X.` (capital letter, period at end)
- Example: `// LoadState reads the state from .kamaji/state.yaml in the given directory.`
- Return values documented in comment when non-obvious

**No comments:**

- Obvious variable assignments
- Self-documenting code with clear naming
- Standard library usage patterns

## Function Design

**Size:** Generally 30-50 lines for orchestration functions, 10-20 for helpers

- Example: `config/state.go` LoadState is ~20 lines
- Example: `orchestrator/run.go` Run is ~110 lines (complex orchestration allowed)

**Parameters:**

- Prefer struct parameters for multiple related arguments
- Examples: `RunConfig`, `SpawnConfig`, `process.SpawnConfig`
- Function signature from `Run(ctx context.Context, cfg RunConfig)` shows config grouping pattern
- Context always first parameter when used: `ctx context.Context`

**Return values:**

- Pointer returns for mutable types: `*RunResult`, `*domain.State`
- Error as last return value: `(*RunResult, error)` pattern
- Tuple returns for success/failure: `(result, error)`
- Empty struct OK when nothing to return: `struct{}` used in channels

**Receiver methods:**

- Value receivers for read-only operations
- Pointer receivers for state-mutating operations
- Example: `(h *Handler) OnPass()` mutates handler state via methods it calls

## Module Design

**Exports:**

- Exported names: package-level functions and types are public API
- Unexported helpers start with lowercase
- Config package exports Load/Save functions for all file types

**Barrel Files:**

- Not used in this codebase
- Each package is focused (e.g., `git/`, `config/`, `orchestrator/`)

**Package structure:**

- `internal/` directory prefix indicates internal packages (not public API)
- Packages grouped by concern: `git`, `config`, `orchestrator`, `domain`, `output`, `mcp`, `process`, `statemachine`, `testutil`
- `domain/` contains pure data types with no I/O

**Interfaces:**

- Small, focused interfaces (1-3 methods typical)
- Example: `ProcessSpawner` has single `Spawn()` method
- Enables testing via mock implementations

---

_Convention analysis: 2026-01-23_
