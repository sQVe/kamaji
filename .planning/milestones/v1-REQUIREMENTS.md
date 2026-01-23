# Requirements Archive: v1 MVP

**Archived:** 2026-01-23
**Status:** ✅ SHIPPED

This is the archived requirements specification for v1.
For current requirements, see `.planning/REQUIREMENTS.md` (created for next milestone).

---

## What This Is

External Go CLI that orchestrates autonomous coding sprints by spawning fresh Claude Code sessions per task. Solves context pollution and exhaustion by owning the state machine externally and giving each task a clean context window.

## Core Value

Reliable state machine. Never lose progress, survive crashes, always know exactly where you are in the sprint.

## Requirements

### Validated

- [x] State machine that tracks current ticket/task position — v1
- [x] MCP server exposing task_complete and note_insight tools — v1
- [x] YAML config parsing (kamaji.yaml for sprints, .kamaji/state.yaml for runtime) — v1
- [x] Spawn Claude Code with fresh context per task — v1
- [x] Git operations (branch creation, commits on pass, reset on fail) — v1
- [x] Ticket logs with history (completed, failed attempts, insights) — v1
- [x] Context injection via XML prompt structure — v1
- [x] Streaming output to terminal — v1
- [x] Stuck detection (3+ consecutive failures) — v1

### Out of Scope

- Parallel ticket execution — v2
- Worktree support — v2
- Service management — v2
- kamaji status/stop/retry commands — v2
- Support for agents other than Claude Code — v2

## Constraints

- **Language**: Go - required
- **Agent**: Claude Code only for v1
- **Testing**: TDD mandatory - tests written before implementation
- **Dependencies**: cobra, lipgloss, go-internal/testscript, yaml.v3

---

## Milestone Summary

**Shipped:** 9 of 9 v1 requirements
**Adjusted:** None
**Dropped:** None

---

_Archived: 2026-01-23 as part of v1 milestone completion_
