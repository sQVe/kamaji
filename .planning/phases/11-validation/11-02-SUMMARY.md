---
phase: 11-validation
plan: 02
subsystem: cli
tags: [cobra, validation, cli]

# Dependency graph
requires:
    - phase: 11-01
      provides: ValidateSprint function with field path error reporting
provides:
    - kamaji validate command integrated into CLI
    - Clear error reporting for validation failures
    - Success/failure exit codes for CI integration
affects: [11-03]

# Tech tracking
tech-stack:
    added: []
    patterns: [cobra command structure, sentinel error pattern]

key-files:
    created:
        - cmd/kamaji/validate.go
    modified:
        - cmd/kamaji/main.go

key-decisions:
    - "Use output package for consistent styled terminal output"
    - "Return errConfigInvalid sentinel to prevent double error printing"

patterns-established:
    - "Sentinel errors for command-level failures"
    - "SilenceUsage = true for expected error cases"

# Metrics
duration: 2min
completed: 2026-01-23
---

# Phase 11 Plan 02: Validation Command Summary

**CLI validate command with field-path error reporting and styled terminal output**

## Performance

- **Duration:** 2 min
- **Started:** 2026-01-23T14:44:16Z
- **Completed:** 2026-01-23T14:46:29Z
- **Tasks:** 3
- **Files modified:** 2

## Accomplishments

- Created kamaji validate command following existing CLI patterns
- Integrated ValidateSprint to show all validation errors at once
- Styled error output with field paths for clear debugging
- Success message for valid configurations
- Proper exit codes for CI integration (0 = valid, 1 = invalid)

## Task Commits

Each task was committed atomically:

1. **Tasks 1-2: Create and wire validate command** - `81eec21` (feat)

_Note: Tasks 1 and 2 were committed together to avoid linter errors from unused code._

## Files Created/Modified

- `cmd/kamaji/validate.go` - Validate command implementation with LoadSprint and ValidateSprint calls
- `cmd/kamaji/main.go` - Wire validate command into root, handle errConfigInvalid sentinel

## Decisions Made

- Used output package (PrintError/PrintSuccess) for consistent terminal styling
- Sentinel error pattern prevents double-printing of validation errors
- SilenceUsage = true matches existing start command pattern

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- Validate command complete and tested
- Ready for 11-03 to add unit tests for validate command
- All manual verification scenarios passed:
    - Valid config: success message, exit 0
    - Missing file: clear error, exit 1
    - Invalid YAML: parser error, exit 1
    - Semantic errors: field paths shown, exit 1

---

_Phase: 11-validation_
_Completed: 2026-01-23_
