---
phase: 11-validation
plan: 01
subsystem: validation
tags: [go, testing, tdd, validation, config]

# Dependency graph
requires:
    - phase: v1.0
      provides: domain.Sprint type and existing validateSprint pattern
provides:
    - ValidationError type for multi-error collection
    - ValidateSprint function that collects all errors instead of short-circuiting
    - Semantic validation with whitespace detection
    - Field path format (tickets[0].tasks[1].description)
affects: [11-02-cli-integration, 12-init-command, 13-refine-command]

# Tech tracking
tech-stack:
    added: []
    patterns:
        - Multi-error validation with error collection
        - Helper functions for validation logic (validateRequired, validateNotEmpty)
        - TDD with RED-GREEN-REFACTOR cycle

key-files:
    created:
        - internal/config/validate.go
        - internal/config/validate_test.go
    modified: []

key-decisions:
    - "Separate ValidateSprint from existing validateSprint to support different use cases"
    - "Use strings.TrimSpace for semantic empty checks to catch whitespace-only content"
    - "Extract validation helpers to reduce duplication"

patterns-established:
    - "ValidationError struct pattern for field-level errors"
    - "Multi-error collection pattern for comprehensive validation feedback"

# Metrics
duration: 2min
completed: 2026-01-23
---

# Phase 11 Plan 01: Core validation logic Summary

**Multi-error validation with semantic whitespace checks and field path reporting**

## Performance

- **Duration:** 2 min
- **Started:** 2026-01-23T14:39:18Z
- **Completed:** 2026-01-23T14:41:34Z
- **Tasks:** 1 (TDD: 3 commits)
- **Files modified:** 2

## Accomplishments

- ValidationError type with Field and Message for structured error reporting
- ValidateSprint function that collects ALL errors, never short-circuits
- Semantic validation using strings.TrimSpace to catch whitespace-only descriptions
- Field path format like tickets[0].tasks[1].description for precise error location
- Comprehensive test coverage with 203 lines across 7 test cases

## Task Commits

TDD task produced 3 atomic commits:

1. **RED: Write failing tests** - `9e36a8a` (test)
    - ValidationError struct and ValidateSprint stub
    - 7 test cases covering all validation scenarios
    - All tests fail as expected

2. **GREEN: Implement validation** - `a5d9bcd` (feat)
    - Full ValidateSprint implementation
    - Collects all errors, validates required fields and whitespace
    - All tests pass

3. **REFACTOR: Extract helpers** - `1b433a5` (refactor)
    - validateRequired and validateNotEmpty helper functions
    - Simplified main logic
    - Tests still pass

## Files Created/Modified

- `internal/config/validate.go` - ValidationError type and ValidateSprint function with helper functions
- `internal/config/validate_test.go` - Comprehensive test suite (203 lines)

## Decisions Made

**1. Keep existing validateSprint separate**

- Rationale: LoadSprint needs fail-fast behavior, validate command needs multi-error collection
- Two different use cases, two different functions
- No modification to existing sprint.go validation

**2. Use strings.TrimSpace for semantic checks**

- Rationale: Catches whitespace-only content that would pass empty string checks
- Separate "required" vs "cannot be empty" messages for clarity

**3. Extract validation helpers**

- Rationale: Reduces duplication, makes main logic more readable
- Helper functions return updated error slice for consistency

## Deviations from Plan

None - plan executed exactly as written following TDD methodology.

## Issues Encountered

None - TDD workflow proceeded smoothly with RED-GREEN-REFACTOR cycle.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- Core validation logic complete and tested
- Ready for CLI integration in 11-02
- ValidationError type exported and ready for use
- Field path format established for error reporting

---

_Phase: 11-validation_
_Completed: 2026-01-23_
