# State persistence

How to survive context death and maintain progress across sessions.

## Core principle

Assume context will die. Design for recovery.

The filesystem is the single source of truth. Context comes from disk changes, not conversation history. When a new agent starts, it reads state from files and continues where the last one left off.

## Context window limits

| Model                        | Context size     |
| ---------------------------- | ---------------- |
| Claude Sonnet/Opus standard  | 200K tokens      |
| Claude Sonnet 4.5 Enterprise | 500K tokens      |
| Claude Sonnet 4/4.5 Tier 4+  | 1M tokens (beta) |

Claude Code summarizes earlier messages when approaching limits, enabling longer conversations but losing early context detail.

## Health thresholds

| Level    | Context used | Recommendation           |
| -------- | ------------ | ------------------------ |
| Healthy  | <60%         | Work normally            |
| Warning  | 60-80%       | Wrap up current subtask  |
| Critical | >80%         | Commit progress and exit |

Performance degrades significantly in the final 20% of the context window.

## Degradation signs

Signs of context exhaustion:

- Forgetting earlier conversation context
- Repeating completed work
- Losing track of file modifications
- Inconsistent task state understanding

## File-based memory

Each iteration loads the same state files:

| File                         | Purpose                                          |
| ---------------------------- | ------------------------------------------------ |
| `kamaji.yaml`                | Sprint definition (tasks, rules, tickets)        |
| `.kamaji/state.yaml`         | Runtime state (current position, failure count)  |
| `.kamaji/logs/<ticket>.yaml` | Per-ticket history (completed, failed, insights) |

The agent doesn't need to remember—it reads current state from disk.

## Guardrails as persistent memory

When an agent hits a blocker, record it so future iterations avoid the same mistake:

```yaml
# .kamaji/logs/login-form.yaml
failed_attempts:
    - task: "Add OAuth integration"
      summary: "passport.js conflicts with existing session middleware"
insights:
    - "Codebase uses Zustand for state management"
    - "Validation schemas are in src/schemas/"
```

The next agent reads these "signs" and adjusts its approach.

## Reactive guardrail tuning

Add guardrails reactively through observation:

1. **Start guardrail-free** — Let the agent build first
2. **Watch for failure patterns** — Observe where it goes wrong
3. **Add rules after failure** — Document what to avoid
4. **Iterate** — Repeat until failure modes are covered

This inverts upfront error prevention. Tune based on observed failures.

## Checkpointing strategies

**Commit early, commit often:**

```markdown
After completing each step, commit with descriptive message.
Do not batch multiple changes into single commits.
```

**Atomic task design:**

```yaml
# Bad: Large task that might exceed context
tasks:
  - description: "Implement entire authentication system"

# Good: Atomic subtasks
tasks:
  - description: "Create User model with password hashing"
    done: "npm test -- User.test.ts passes"
  - description: "Add login endpoint"
    done: "curl -X POST /login returns 200"
```

**External state verification:**

```yaml
tasks:
    - description: "Continue migration where left off"
      verify: "Check which migration files exist in db/migrations/"
      done: "All migrations applied, prisma db push succeeds"
```

## Escape hatches

For long-running tasks, include exit conditions:

```yaml
# In sprint rules or task definition
After 3 consecutive failures on the same task:
    - Document blockers in ticket log
    - List attempted approaches
    - Exit for manual intervention
```

This prevents infinite loops on impossible tasks.

## Self-referential feedback loop

Filesystem artifacts create emergent memory:

```
Claude's previous work persists in files
    ↓
Each iteration sees:
  - Modified source files
  - Test results and failures
  - Git history of changes
  - Updated ticket logs
    ↓
Claude reads state and continues
```

The prompt stays constant. Context evolves through disk changes.

---

## Kamaji implications

| Principle                | Application                                                          |
| ------------------------ | -------------------------------------------------------------------- |
| External state is truth  | `.kamaji/state.yaml` is authoritative; Claude's memory is disposable |
| Atomic task design       | Each task should complete in reasonable context                      |
| Verification over memory | Verify state from filesystem; distrust "I remember we did X"         |
| Failure logging          | Record failed attempts and insights for future iterations            |
| Reactive tuning          | Add sprint rules based on observed failure patterns                  |

Kamaji is well-positioned for state persistence:

- Fresh subagent per task (no accumulated pollution)
- State persists in `.kamaji/state.yaml`
- Task completion verified independently
- History injected into future task context

---

## Sources

- [The Ralph Playbook](https://claytonfarr.github.io/ralph-playbook/)
- [Claude context window documentation](https://platform.claude.com/docs/en/build-with-claude/context-windows)
- [Anthropic Claude Code Ralph Plugin](https://github.com/anthropics/claude-code/blob/main/plugins/ralph-wiggum/README.md)
