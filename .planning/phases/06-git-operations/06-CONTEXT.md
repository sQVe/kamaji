# Phase 6: Git Operations - Context

**Gathered:** 2026-01-15
**Status:** Ready for research

<vision>
## How This Should Work

Each task gets its own branch for clean isolation. When a task passes, work is committed (but not auto-merged — merging stays manual). When a task fails, a hard reset wipes all changes completely so the next attempt starts fresh.

The git layer should be invisible during normal operation — it just works. But the history should tell a clear story: what was attempted, what passed, what failed.

</vision>

<essential>
## What Must Be Nailed

- **Never lose good work** — Once a task passes and commits, that work is safe
- **Clean isolation** — Failed attempts can't pollute future attempts or other tasks
- **Traceable history** — Easy to see what was tried, what passed, what failed

</essential>

<boundaries>
## What's Out of Scope

- Conflict resolution — assumes single worker, no merge conflicts
- PR creation — no GitHub/GitLab integration, purely local git
- Auto-merging — commits happen, but merging is a separate concern

</boundaries>

<specifics>
## Specific Ideas

- Commit messages follow conventional commit format
- Reference Grove's git implementation at `/home/sqve/code/personal/grove/main` — research thoroughly for patterns and testing approaches

</specifics>

<notes>
## Additional Context

Grove's git operations implementation should be studied as a reference for both implementation patterns and testing strategies. The user considers this a valuable prior art to learn from.

</notes>

---

_Phase: 06-git-operations_
_Context gathered: 2026-01-15_
