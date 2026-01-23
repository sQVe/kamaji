# Phase 10: Integration - Context

**Gathered:** 2026-01-20
**Status:** Ready for planning

<domain>
## Phase Boundary

End-to-end sprint execution combining all components. The `kamaji start` command launches a sprint that: loads config, initializes state, spawns Claude Code sessions per task, handles pass/fail/stuck outcomes, and manages git operations throughout.

</domain>

<decisions>
## Implementation Decisions

### Claude's Discretion

User delegated all integration behavior decisions to Claude. Apply standard CLI patterns:

**Execution flow:**

- Start fresh or resume based on existing state file
- No confirmation prompts for normal operation
- Clear error if config missing or invalid

**Progress reporting:**

- Show current task and ticket position
- Surface pass/fail outcomes as they happen
- Summary at sprint completion

**Interruption handling:**

- Ctrl+C triggers graceful shutdown
- Save current state before exit
- Next run resumes from last position

**Error presentation:**

- Clear messages for config/state errors
- Show failure reason from Claude's task_complete call
- Stuck mode surfaced with context about consecutive failures

</decisions>

<specifics>
## Specific Ideas

User emphasized: "This is our first real command. We need to plan THOROUGH integration tests."

Testing priority is high — the planner should ensure comprehensive coverage of:

- Happy path (full sprint execution)
- State persistence and resume
- Error conditions (missing config, invalid state)
- Edge cases (empty sprint, all tasks fail, stuck detection)

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

_Phase: 10-integration_
_Context gathered: 2026-01-20_
