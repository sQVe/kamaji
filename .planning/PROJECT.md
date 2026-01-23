# Kamaji

## What This Is

External Go CLI that orchestrates autonomous coding sprints by spawning fresh Claude Code sessions per task. Solves context pollution and exhaustion by owning the state machine externally and giving each task a clean context window.

## Core Value

Reliable state machine. Never lose progress, survive crashes, always know exactly where you are in the sprint.

## Current State

**Shipped:** v1 MVP (2026-01-23)
**Codebase:** 7,570 LOC Go
**Tech stack:** Cobra CLI, lipgloss, testscript, yaml.v3, mcp-go

## Current Milestone: v1.1 Sprint Planning

**Goal:** Add commands to help users create and refine sprint configs before execution.

**Target features:**

- `kamaji init` — Create minimal template config with explanatory comments
- `kamaji refine` — AI-assisted config improvement (spawns Claude to rewrite)
- `kamaji validate` — Schema + heuristic validation (deterministic)
- Integration with `kamaji start` (validate before running)

## Requirements

### Validated

- ✓ State machine that tracks current ticket/task position — v1
- ✓ MCP server exposing task_complete and note_insight tools — v1
- ✓ YAML config parsing (kamaji.yaml for sprints, .kamaji/state.yaml for runtime) — v1
- ✓ Spawn Claude Code with fresh context per task — v1
- ✓ Git operations (branch creation, commits on pass, reset on fail) — v1
- ✓ Ticket logs with history (completed, failed attempts, insights) — v1
- ✓ Context injection via XML prompt structure — v1
- ✓ Streaming output to terminal — v1
- ✓ Stuck detection (3+ consecutive failures) — v1

### Active

- [ ] Init command creates template config with explanatory comments
- [ ] Refine command spawns Claude to improve config based on comments
- [ ] Validate command checks schema and semantic heuristics
- [ ] Start command validates config before running sprint

### Out of Scope

- Parallel ticket execution — deferred to v2
- Worktree support — deferred to v2
- Service management — deferred to v2
- kamaji status/stop/retry commands — deferred to v2
- Support for agents other than Claude Code — deferred to v2

## Context

DESIGN.md contains detailed technical specification including:

- File structure (.kamaji/, kamaji.yaml schemas)
- MCP server transport (Streamable HTTP on localhost)
- Execution flow (10-step process)
- Context injection format (XML structure)
- Git handling strategy
- Failure handling rules

## Constraints

- **Language**: Go - required
- **Agent**: Claude Code only for v1
- **Testing**: TDD mandatory - tests written before implementation
- **Dependencies**: cobra, lipgloss, go-internal/testscript, yaml.v3, mcp-go

## Key Decisions

| Decision                                 | Rationale                                                          | Outcome |
| ---------------------------------------- | ------------------------------------------------------------------ | ------- |
| External orchestrator vs internal plugin | Avoids context pollution, enables parallelization later            | ✓ Good  |
| Streamable HTTP transport                | Replaces deprecated SSE-only, future-proofs for v2 parallelization | ✓ Good  |
| YAML for config                          | Human-readable, easy to edit                                       | ✓ Good  |
| Orchestrator handles git                 | Claude focuses on coding, clean separation                         | ✓ Good  |
| TDD mandatory                            | Ensures correctness, enables fearless refactoring                  | ✓ Good  |
| Script tests (go-internal)               | Integration-focused, proven patterns                               | ✓ Good  |
| github.com/mark3labs/mcp-go              | Mature library, Streamable HTTP support                            | ✓ Good  |

---

_Last updated: 2026-01-23 after v1.1 milestone start_
