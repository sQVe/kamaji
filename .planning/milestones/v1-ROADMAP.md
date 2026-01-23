# Milestone v1: MVP

**Status:** ✅ SHIPPED 2026-01-23
**Phases:** 1-10
**Total Plans:** 23

## Overview

Build an external Go CLI that orchestrates autonomous coding sprints by spawning fresh Claude Code sessions per task. Starting with project foundation and config parsing, then implementing the core state machine, MCP server with Streamable HTTP transport, context injection, and process spawning. Follow with git operations, ticket logging, terminal output, and stuck detection. Conclude by integrating all components into the end-to-end sprint execution flow.

## Phases

### Phase 1: Foundation

**Goal**: Project structure with Cobra CLI, YAML config parsing, and test infrastructure
**Depends on**: Nothing (first phase)
**Plans**: 3 plans

Plans:

- [x] 01-01: Project scaffolding and Cobra CLI setup
- [x] 01-02: YAML config parsing (kamaji.yaml, .kamaji/state.yaml schemas)
- [x] 01-03: Testscript infrastructure setup

### Phase 2: State Machine

**Goal**: Core state tracking with persistence to .kamaji/state.yaml
**Depends on**: Phase 1
**Plans**: 2 plans

Plans:

- [x] 02-01: State machine types and core navigation
- [x] 02-02: Pass/fail handling (RecordPass, RecordFail, IsStuck)

### Phase 3: MCP Server

**Goal**: Streamable HTTP server exposing task_complete and note_insight tools
**Depends on**: Phase 2
**Plans**: 2 plans

Plans:

- [x] 03-01: MCP server infrastructure with Streamable HTTP transport
- [x] 03-02: Tool handlers (task_complete, note_insight)

### Phase 4: Context Injection

**Goal**: XML prompt structure generation following DESIGN.md format
**Depends on**: Phase 1
**Plans**: 2 plans

Plans:

- [x] 04-01: XML template structure (BuildPrompt function)
- [x] 04-02: Context assembly from state and config (AssembleContext function)

### Phase 5: Process Spawning

**Goal**: Launch Claude Code subprocess with injected context
**Depends on**: Phase 3, Phase 4
**Plans**: 3 plans

Plans:

- [x] 05-01: MCP Server Signal Channel
- [x] 05-02: Claude Process Manager
- [x] 05-03: Context Injection & Launch

### Phase 6: Git Operations

**Goal**: Branch creation, commits on pass, reset on fail
**Depends on**: Phase 2
**Plans**: 2 plans

Plans:

- [x] 06-01: Branch management
- [x] 06-02: Commit and reset operations

### Phase 7: Ticket Logging

**Goal**: History tracking with completed tasks, failed attempts, insights
**Depends on**: Phase 2
**Plans**: 2 plans

Plans:

- [x] 07-01: Log file structure and writing
- [x] 07-02: History queries and reporting

### Phase 8: Streaming Output

**Goal**: Terminal output with lipgloss styling
**Depends on**: Phase 5
**Plans**: 2 plans

Plans:

- [x] 08-01: Output styling and formatting
- [x] 08-02: Progress indicators and status display

### Phase 9: Stuck Detection

**Goal**: Detect and handle 3+ consecutive failures
**Depends on**: Phase 2, Phase 7
**Plans**: 2 plans

Plans:

- [x] 09-01: TaskResult type normalizing pass/fail/no-signal outcomes
- [x] 09-02: Handler orchestrating pass/fail/stuck workflows

### Phase 10: Integration

**Goal**: End-to-end sprint execution combining all components with comprehensive integration tests
**Depends on**: All previous phases
**Plans**: 3 plans

Plans:

- [x] 10-01: Sprint orchestration (Run function and main loop)
- [x] 10-02: Integration tests (MCP + handler + git + state verification)
- [x] 10-03: CLI start command with testscript coverage

---

## Milestone Summary

**Key Decisions:**

- External orchestrator vs internal plugin: Avoids context pollution, enables parallelization later
- Streamable HTTP transport: Replaces deprecated SSE-only, future-proofs for v2 parallelization
- YAML for config: Human-readable, easy to edit
- Orchestrator handles git: Claude focuses on coding, clean separation
- TDD mandatory: Ensures correctness, enables fearless refactoring
- Script tests (go-internal): Integration-focused, copy grove patterns
- Using github.com/mark3labs/mcp-go for MCP server implementation

**Issues Resolved:**

- Context pollution in long Claude sessions
- State loss on crashes
- No-signal handling (process exits without MCP call)

**Issues Deferred:**

- Parallel ticket execution (v2)
- Worktree support (v2)
- Service management (v2)
- kamaji status/stop/retry commands (v2)

**Technical Debt Incurred:**

- None identified

---

_For current project status, see .planning/ROADMAP.md_
