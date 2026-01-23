# Phase 11: Validation - Research

**Researched:** 2026-01-23
**Domain:** CLI configuration validation, YAML schema validation, semantic validation
**Confidence:** HIGH

## Summary

Validation for CLI tools requires three layers: YAML syntax validation, schema validation (required fields, types), and semantic validation (business logic like dependency references). The Go ecosystem offers multiple approaches, but the current codebase already uses `gopkg.in/yaml.v3` which handles syntax and basic unmarshaling. For enhanced validation with better error reporting, `goccy/go-yaml` (published Jan 2026) provides superior line/column error reporting with source code snippets. For struct validation, `go-playground/validator/v10` is the industry standard.

The existing codebase already has basic validation in `config.validateSprint()` that checks required fields. The validate command should enhance this with better error formatting and add semantic checks for empty descriptions and dependency references.

**Primary recommendation:** Extend existing `config.validateSprint()` function with semantic checks, wrap in new validate command using Cobra patterns, output errors using existing `output.Style()` functions with clear location information.

## Standard Stack

The established libraries/tools for this domain:

### Core

| Library                           | Version | Purpose          | Why Standard                                            |
| --------------------------------- | ------- | ---------------- | ------------------------------------------------------- |
| gopkg.in/yaml.v3                  | v3      | YAML parsing     | Already in use, standard library, handles syntax errors |
| github.com/spf13/cobra            | v1.9.1  | CLI framework    | Already in use, de facto standard for Go CLIs           |
| github.com/charmbracelet/lipgloss | v1.0.0  | Terminal styling | Already in use, modern terminal output                  |

### Supporting

| Library                     | Version             | Purpose               | When to Use                                                |
| --------------------------- | ------------------- | --------------------- | ---------------------------------------------------------- |
| goccy/go-yaml               | latest (2026-01-08) | Enhanced YAML parsing | If line/column numbers needed, better error formatting     |
| go-playground/validator/v10 | v10.26+             | Struct validation     | If complex cross-field validation needed (likely overkill) |

### Alternatives Considered

| Instead of        | Could Use                    | Tradeoff                                                             |
| ----------------- | ---------------------------- | -------------------------------------------------------------------- |
| gopkg.in/yaml.v3  | goccy/go-yaml                | Better error reporting but adds dependency, migration effort         |
| Custom validation | go-playground/validator      | More powerful but overkill for simple checks, adds learning curve    |
| Plain text errors | Structured validation errors | More complex but better for tooling integration (future JSON output) |

**Installation:**
No new dependencies required for v1 scope. Existing stack sufficient.

## Architecture Patterns

### Recommended Project Structure

```
cmd/kamaji/
├── validate.go          # New validate command
internal/
├── config/
│   └── sprint.go        # Enhance validateSprint()
└── validator/           # Optional: semantic validation logic
    ├── validator.go     # Main validation orchestrator
    ├── schema.go        # Schema validation (required fields)
    └── semantic.go      # Semantic checks (deps, empty strings)
```

### Pattern 1: Command Structure

**What:** Cobra command with PreRunE for early validation, exit codes for automation
**When to use:** All CLI commands that can fail
**Example:**

```go
// Source: Terraform validate pattern + existing kamaji start.go
func validateCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "validate",
        Short: "Validate kamaji.yaml configuration",
        RunE: func(cmd *cobra.Command, _ []string) error {
            workDir, err := os.Getwd()
            if err != nil {
                return err
            }

            errors := ValidateConfig(filepath.Join(workDir, "kamaji.yaml"))
            if len(errors) > 0 {
                PrintValidationErrors(errors)
                return errConfigInvalid // Exit code 1
            }

            output.PrintSuccess("Configuration is valid")
            return nil // Exit code 0
        },
    }

    cmd.SilenceUsage = true
    return cmd
}
```

### Pattern 2: Layered Validation

**What:** Separate validation into syntax, schema, and semantic layers
**When to use:** Complex validation with different error types
**Example:**

```go
// Source: Industry pattern from Terraform, kubectl
type ValidationError struct {
    Field   string // e.g., "tickets[0].tasks[1].description"
    Message string // User-friendly message
    Line    int    // Optional: YAML line number
}

func ValidateConfig(path string) []ValidationError {
    errors := []ValidationError{}

    // Layer 1: Syntax (handled by yaml.Unmarshal)
    sprint, err := config.LoadSprint(path)
    if err != nil {
        // Parse YAML error for line numbers if available
        return []ValidationError{{Message: err.Error()}}
    }

    // Layer 2: Schema (enhance existing validateSprint)
    errors = append(errors, validateSchema(sprint)...)

    // Layer 3: Semantic
    errors = append(errors, validateSemantic(sprint)...)

    return errors
}
```

### Pattern 3: Error Formatting

**What:** Clear, actionable error messages with location context
**When to use:** All validation errors
**Example:**

```go
// Source: Existing output package patterns
func PrintValidationErrors(errors []ValidationError) {
    output.PrintError("Configuration validation failed:")
    fmt.Println()

    for _, err := range errors {
        if err.Field != "" {
            fmt.Printf("  %s %s: %s\n",
                output.ErrorMsg("✗"),
                err.Field,
                err.Message)
        } else {
            fmt.Printf("  %s %s\n",
                output.ErrorMsg("✗"),
                err.Message)
        }
    }

    fmt.Printf("\n%d error(s) found\n", len(errors))
}
```

### Anti-Patterns to Avoid

- **Parsing YAML twice:** Don't reparse for validation, use loaded struct
- **Cryptic error messages:** Always include field path and actionable guidance
- **Silent failures:** Validate command must exit non-zero on validation failure
- **Overly complex validation:** Don't add JSON Schema or complex framework for simple checks

## Don't Hand-Roll

Problems that look simple but have existing solutions:

| Problem                  | Don't Build              | Use Instead                             | Why                                                                           |
| ------------------------ | ------------------------ | --------------------------------------- | ----------------------------------------------------------------------------- |
| Struct field validation  | Custom field checkers    | go-playground/validator tags            | Handles all edge cases, tested, well-documented (BUT: likely overkill for v1) |
| Terminal color detection | Custom TTY detection     | lipgloss (already in use)               | Handles NO_COLOR, pipe detection, platform differences                        |
| CLI exit codes           | Magic numbers            | Named error types like errConfigInvalid | Self-documenting, consistent with Cobra patterns                              |
| YAML line numbers        | String parsing of errors | goccy/go-yaml FormatError               | Handles all YAML syntax cases, colorized output                               |

**Key insight:** The existing validation in `config.validateSprint()` is already well-structured. Don't rebuild it, extend it with semantic checks.

## Common Pitfalls

### Pitfall 1: Poor Error Location Context

**What goes wrong:** Validation errors say "description missing" without indicating which ticket/task
**Why it happens:** Using simple string matching or not tracking position during traversal
**How to avoid:** Always include array indices in error messages (e.g., "tickets[2].tasks[0]")
**Warning signs:** User has to guess where the error is, or must validate entire file manually

### Pitfall 2: Inconsistent Exit Codes

**What goes wrong:** Validate command exits 0 even with errors, or uses inconsistent codes
**Why it happens:** Not following CLI conventions or Cobra error handling patterns
**How to avoid:** Follow pattern: 0 = valid, 1 = invalid/error, use `errConfigInvalid` like `errSprintFailed`
**Warning signs:** Automation scripts can't detect validation failures

### Pitfall 3: YAML Syntax Error Line Numbers

**What goes wrong:** `gopkg.in/yaml.v3` has known issues with incorrect line numbers in errors
**Why it happens:** Library bug reported in go-yaml/yaml#1055
**How to avoid:** Accept that line numbers may be approximate, or migrate to `goccy/go-yaml`
**Warning signs:** User reports error location is wrong, can't find syntax issue

### Pitfall 4: Empty String vs Missing Field

**What goes wrong:** Treating `description: ""` same as missing `description:` field
**Why it happens:** Go's zero value for string is `""`
**How to avoid:** Both are invalid for required fields, catch with `if field == ""` check
**Warning signs:** User provides empty description, validation passes incorrectly

### Pitfall 5: Circular Dependency Detection

**What goes wrong:** Validating that ticket dependencies form a DAG is complex
**Why it happens:** Dependencies aren't implemented yet, premature validation
**How to avoid:** Wait until dependencies are added to schema (not in v1), then use topological sort
**Warning signs:** Adding validation for features that don't exist yet

### Pitfall 6: Over-validation

**What goes wrong:** Checking things like branch name format, git branch existence
**Why it happens:** Trying to be too helpful, validation becomes runtime checks
**How to avoid:** Validate schema/structure only, not runtime state. Git operations fail naturally with clear errors.
**Warning signs:** Validation needs network/filesystem access, takes long time

## Code Examples

Verified patterns from official sources:

### Semantic Validation: Empty Description Check

```go
// Source: Existing config/sprint.go validateSprint pattern
func validateSemantic(s *domain.Sprint) []ValidationError {
    errors := []ValidationError{}

    // Check for empty descriptions (VALD-03)
    for i, ticket := range s.Tickets {
        if strings.TrimSpace(ticket.Description) == "" {
            errors = append(errors, ValidationError{
                Field:   fmt.Sprintf("tickets[%d].description", i),
                Message: "description cannot be empty",
            })
        }

        for j, task := range ticket.Tasks {
            if strings.TrimSpace(task.Description) == "" {
                errors = append(errors, ValidationError{
                    Field:   fmt.Sprintf("tickets[%d].tasks[%d].description", i, j),
                    Message: "description cannot be empty",
                })
            }
        }
    }

    return errors
}
```

### Semantic Validation: Dependency References (Future)

```go
// Source: Dependency graph pattern from https://dnaeon.github.io/dependency-graph-resolution-algorithm-in-go/
// NOTE: Dependencies not in v1 schema, this is for future reference
func validateDependencies(s *domain.Sprint) []ValidationError {
    errors := []ValidationError{}
    ticketNames := make(map[string]bool)

    // Build ticket name index
    for _, ticket := range s.Tickets {
        ticketNames[ticket.Name] = true
    }

    // Check each ticket's dependencies exist
    for i, ticket := range s.Tickets {
        for j, dep := range ticket.DependsOn { // Future field
            if !ticketNames[dep] {
                errors = append(errors, ValidationError{
                    Field:   fmt.Sprintf("tickets[%d].depends_on[%d]", i, j),
                    Message: fmt.Sprintf("dependency '%s' does not exist", dep),
                })
            }
        }
    }

    return errors
}
```

### Error Output with Existing Styles

```go
// Source: Existing internal/output/styles.go
func PrintValidationErrors(errors []ValidationError) {
    output.PrintError(fmt.Sprintf("Configuration validation failed with %d error(s):", len(errors)))
    fmt.Println()

    for _, err := range errors {
        var msg string
        if err.Field != "" {
            msg = fmt.Sprintf("%s: %s", err.Field, err.Message)
        } else {
            msg = err.Message
        }

        // Uses existing output package (respects --plain mode)
        fmt.Fprintf(os.Stderr, "  %s\n", output.ErrorMsg(msg))
    }
}
```

### Integrating with Existing Validation

```go
// Source: Enhance existing config/sprint.go validateSprint
func validateSprint(s *domain.Sprint) error {
    // Existing required field checks
    if s.Name == "" {
        return fmt.Errorf("sprint missing required field: name")
    }

    for i, ticket := range s.Tickets {
        if ticket.Name == "" {
            return fmt.Errorf("ticket[%d] missing required field: name", i)
        }

        // NEW: Add semantic check for empty description
        if strings.TrimSpace(ticket.Description) == "" {
            return fmt.Errorf("ticket[%d] description cannot be empty", i)
        }

        for j, task := range ticket.Tasks {
            if task.Description == "" {
                return fmt.Errorf("ticket[%d].task[%d] missing required field: description", i, j)
            }

            // NEW: Add semantic check for empty description
            if strings.TrimSpace(task.Description) == "" {
                return fmt.Errorf("ticket[%d].task[%d] description cannot be empty", i, j)
            }
        }
    }

    return nil
}
```

## State of the Art

| Old Approach             | Current Approach                 | When Changed       | Impact                                                       |
| ------------------------ | -------------------------------- | ------------------ | ------------------------------------------------------------ |
| gopkg.in/yaml.v3         | goccy/go-yaml                    | 2026-01-08         | Better error messages with line numbers, source snippets     |
| Manual struct validation | go-playground/validator          | v10+ stable        | Tag-based validation, but may be overkill for simple schemas |
| Plain error strings      | Structured ValidationError types | Modern CLI pattern | Enables future JSON output, better UX                        |
| Single error return      | Multiple error collection        | Modern pattern     | Show all errors at once, not first-failure-only              |

**Deprecated/outdated:**

- gopkg.in/yaml.v2: Use v3 (already doing this)
- go-playground/validator v9: Use v10+ with WithRequiredStructEnabled
- Binary exit codes without named errors: Use named sentinel errors

## Open Questions

Things that couldn't be fully resolved:

1. **Line number accuracy with gopkg.in/yaml.v3**
    - What we know: Known issue with incorrect line numbers (go-yaml/yaml#1055)
    - What's unclear: How severe in practice, worth migration to goccy/go-yaml?
    - Recommendation: Accept approximate line numbers for syntax errors, focus on clear field paths for validation errors

2. **Validation error output format**
    - What we know: Plain text sufficient for v1, terraform validate supports JSON output
    - What's unclear: Should we plan for --format=json from start?
    - Recommendation: Design ValidationError struct to support JSON serialization but only implement text output for v1

3. **Validation thoroughness**
    - What we know: Requirements specify empty descriptions and dependency checking
    - What's unclear: Should we validate branch name format, base_branch existence, etc.?
    - Recommendation: Validate structure/content only (VALD-02, VALD-03), not runtime state. Keep it fast and deterministic.

4. **Dependency validation scope**
    - What we know: VALD-03 says "deps exist" but dependencies not in v1 schema
    - What's unclear: Is this for future feature or misunderstanding of requirements?
    - Recommendation: Skip dependency validation unless schema adds `depends_on` field. Focus on empty description checks.

## Sources

### Primary (HIGH confidence)

- [goccy/go-yaml documentation](https://pkg.go.dev/github.com/goccy/go-yaml) - Official docs for enhanced YAML library with error reporting
- [go-playground/validator/v10 documentation](https://pkg.go.dev/github.com/go-playground/validator/v10) - Official validation library docs
- [Cobra user guide](https://github.com/spf13/cobra/blob/main/site/content/user_guide.md) - Official Cobra framework patterns
- [Terraform validate command](https://developer.hashicorp.com/terraform/cli/commands/validate) - Industry standard validate command patterns
- Existing kamaji codebase: `internal/config/sprint.go`, `internal/output/styles.go` - Current patterns

### Secondary (MEDIUM confidence)

- [A Guide to Input Validation in Go with Validator V10](https://dev.to/kittipat1413/a-guide-to-input-validation-in-go-with-validator-v10-56bp) - Validation patterns
- [Dependency graph resolution algorithm in Go](https://dnaeon.github.io/dependency-graph-resolution-algorithm-in-go/) - DAG validation patterns
- [How to Build a CLI Tool in Go with Cobra](https://oneuptime.com/blog/post/2026-01-07-go-cobra-cli/view) - Recent Cobra best practices
- [Managing Circular Dependencies in Go](https://medium.com/@cosmicray001/managing-circular-dependencies-in-go-best-practices-and-solutions-723532f04dde) - Dependency validation patterns

### Tertiary (LOW confidence)

- WebSearch results on YAML validation - General patterns, marked for validation in implementation
- go-yaml/yaml GitHub issues - Known bugs and limitations, not officially resolved

## Metadata

**Confidence breakdown:**

- Standard stack: HIGH - Existing dependencies sufficient, well-tested patterns
- Architecture: HIGH - Clear command structure from existing start.go, validation patterns from config/sprint.go
- Pitfalls: HIGH - Based on known library issues, Cobra patterns, and existing codebase conventions

**Research date:** 2026-01-23
**Valid until:** 2026-02-23 (30 days - stable domain, minimal churn expected)
