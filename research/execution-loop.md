# Execution loop

How to orchestrate autonomous coding sessions.

## Core insight

Individual LLM runs are unreliable, but the filesystem is deterministic. Code, tests, and git commits survive context rotation. Each agent inherits accumulated work through files, not conversation history.

The canonical implementation (Ralph Loop):

```bash
while :; do cat PROMPT.md | claude-code; done
```

## How it works

Each iteration:

1. Feed task definition to AI agent
2. Agent works until context fills or task completes
3. Agent exits (or orchestrator detects completion)
4. Loop restarts with fresh context
5. New agent reads filesystem state and continues

Progress persists in the filesystem and git history, not the LLM's context window.

## Internal vs external orchestration

**Internal (GSD pattern):** Loop runs inside Claude Code using hooks and subagents. Context accumulates across iterations until the window fills. Benefits from iterative memory but risks context pollution.

**External (kamaji pattern):** Orchestrator spawns fresh Claude Code sessions per task. Each session starts clean. No accumulated confusion, but no iterative memory either.

| Aspect           | Internal                | External                            |
| ---------------- | ----------------------- | ----------------------------------- |
| Context per task | Accumulates             | Fresh                               |
| State management | Claude tracks via files | Orchestrator owns state             |
| Verification     | Claude self-reports     | Orchestrator verifies independently |
| Failure recovery | Log and continue        | Reset and retry                     |

Kamaji uses external orchestration—it owns the state machine, spawns fresh contexts, and can verify completion independently.

## Binary success criteria

Tasks must have machine-verifiable completion conditions. The loop needs to know when to stop.

**Good criteria:**

- All tests pass
- Build succeeds with zero errors
- Coverage exceeds 85%
- Linter reports zero warnings

**Bad criteria:**

- "Make it better"
- "Improve code quality"
- "Refactor appropriately"

## Commit pattern

Each task gets its own commit immediately after completion:

```
task 1 complete → commit
task 2 complete → commit
task 3 complete → commit
```

The orchestrator handles git, not Claude. On pass, commit with task summary. On fail, reset to HEAD for clean retry.

## Failure handling

- **Failure count:** Consecutive failures on current task (resets on pass)
- **Stuck threshold:** N consecutive failures triggers exit
- **Exit without signal:** Treated as failure (Claude crashed or forgot to signal)

When stuck, exit with failure and leave state intact for manual intervention.

## Subagent isolation

Each task spawns in a fresh context. This keeps Claude in its "smart zone"—early context with no accumulated confusion from previous work or failed attempts.

Benefits:

- No context pollution from past failures
- Predictable token budget per task
- Independent tasks can run in parallel

## When to use autonomous loops

Tasks with clear, automatable verification:

- **Test coverage expansion** — Run tests, check percentage
- **Code refactoring** — Existing tests define correctness
- **Database migrations** — Up/down scripts work or don't
- **Dependency upgrades** — Build passes or fails
- **API implementations** — Contract tests verify behavior

## When NOT to use autonomous loops

- **Vague goals** — "Make it better" has no endpoint
- **Deep codebase archaeology** — Context-dependent understanding gets lost
- **UI/UX polish** — Requires human judgment
- **Manual verification required** — Loop can't check what scripts can't
- **Subjective quality** — No binary pass/fail

## State machine

```
┌─────────────┐
│  PICK_TASK  │
└──────┬──────┘
       │ task selected
       ▼
┌──────────────┐
│  WORK_TASK   │◄─────────┐
│  (run/verify)│          │
└──────┬───────┘          │
       │ complete         │ fail (retry)
       ▼                  │
┌──────────────┐          │
│  CHECKPOINT  ├──────────┘
│  (pass/fail?)│
└──────┬───────┘
       │ all done
       ▼
┌──────────────┐
│     DONE     │
└──────────────┘
```

---

## Kamaji implications

| Aspect                 | Application                                      |
| ---------------------- | ------------------------------------------------ |
| External orchestration | Kamaji owns state, spawns fresh Claude sessions  |
| Binary verification    | `verify` + `done` fields enable machine checking |
| Atomic commits         | Orchestrator commits on pass, resets on fail     |
| Stuck detection        | Exit after N consecutive failures                |
| Parallel execution     | Independent tickets can run simultaneously       |

---

## Sources

- [The Ralph Playbook](https://claytonfarr.github.io/ralph-playbook/)
- [Geoffrey Huntley's how-to-ralph-wiggum](https://github.com/ghuntley/how-to-ralph-wiggum)
- [glittercowboy/get-shit-done](https://github.com/glittercowboy/get-shit-done)
