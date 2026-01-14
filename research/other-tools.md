# Other tools

Patterns from autonomous coding tools applicable to kamaji.

## Aider

Git-aware AI pair programming in the terminal.

**Repository map:** Aider builds a concise map of classes, functions, and relationships using tree-sitter for AST-aware analysis. Claude sees structure, not content, until files are explicitly added.

```
Repository → tree-sitter parse → function signatures + relationships → compressed map
```

**Architect/editor split:** Two-model pattern where the architect proposes and the editor implements.

| Role      | Purpose                                |
| --------- | -------------------------------------- |
| Architect | Describes how to solve the problem     |
| Editor    | Translates description into file edits |

**Incremental file inclusion:** Users add only files that need changes. The repo map pulls context from related files automatically.

**Prompt caching:** Aider caches unchanged context across requests.

## OpenHands

Open platform for autonomous AI development with sandboxed execution.

**Event-sourced state:** State represented through an event log with deterministic replay.

```
action_1 → observation_1 → action_2 → observation_2 → ...
```

**Workspace abstraction:** Same agent runs locally or containerized.

```python
workspace = LocalWorkspace()           # Local prototyping
workspace = ContainerizedWorkspace()   # Production
```

**Sandboxing:** Docker-based isolation. Agents restricted to own container, torn down post-session.

**Multi-agent delegation:** Hierarchical agent structures with explicit capability grants.

## SWE-agent

Princeton/Stanford research agent for autonomous GitHub issue fixing.

**Agent-computer interface (ACI):** Interface design matters as much as model capabilities. SWE-agent achieved 12.5% resolution rate on SWE-bench by optimizing agent-code interaction.

**Single YAML configuration:** Entire agent behavior governed by one config file.

```yaml
model: claude-sonnet-4
tools:
    - file_read
    - file_edit
    - terminal
max_iterations: 30
```

**Benchmark-driven design:** Given a GitHub issue, produce a patch. Clear success criteria enable iteration.

## Cline

Autonomous coding agent in VS Code with explicit permission steps.

**Plan + Act mode:** Two-phase approach where the model creates an outline, then works through it.

1. **Plan** — Create outline of needed changes
2. **Act** — Work through the outline

**Context engineering:** Dynamic context management, AST-based analysis, and memory bank for tribal knowledge.

**MCP integration:** Uses Model Context Protocol to extend capabilities and create tools dynamically.

## Claude Code hooks

Official lifecycle hooks for controlling Claude Code behavior.

**Hook events:**

| Event          | Timing                     | Use case                    |
| -------------- | -------------------------- | --------------------------- |
| `PreToolUse`   | Before tool calls          | Block or modify tool inputs |
| `PostToolUse`  | After tool completes       | React to tool output        |
| `Stop`         | Claude finishes responding | End-of-turn validation      |
| `SubagentStop` | Subagent finishes          | Subagent validation         |

**Matchers:** Filter which tools trigger hooks.

```json
{
    "hooks": {
        "PostToolUse": [
            {
                "matcher": "Edit|MultiEdit|Write",
                "command": "./validate-edits.sh"
            }
        ]
    }
}
```

**Input modification:** PreToolUse hooks can modify tool inputs before execution.

**Stop hooks for validation:** Ideal for end-of-turn quality gates (lint, typecheck, test).

---

## Kamaji implications

| Tool      | Pattern                | Application                              |
| --------- | ---------------------- | ---------------------------------------- |
| Aider     | Repo map               | Compress context for large codebases     |
| Aider     | Architect/editor       | Consider two-phase reasoning             |
| OpenHands | Event-sourced state    | Enable debugging and replay of actions   |
| OpenHands | Sandboxing             | Isolated execution environments          |
| SWE-agent | Single config          | Keep configuration simple                |
| Cline     | Plan + Act             | Explicit planning phase before execution |
| Hooks     | Lifecycle interception | Verify changes independently of Claude   |

**High-priority patterns:**

| Pattern              | Reason                                               |
| -------------------- | ---------------------------------------------------- |
| Repo maps            | Compress context for large codebases                 |
| Event-sourced state  | Enable debugging and replay                          |
| Plan + Act           | Consider explicit planning before task execution     |
| Hooks for validation | Verify changes independently of Claude's self-report |

---

## Sources

- [Aider Documentation](https://aider.chat/docs/)
- [OpenHands GitHub](https://github.com/OpenHands/OpenHands)
- [OpenHands Technical Report](https://arxiv.org/abs/2511.03690)
- [SWE-agent Documentation](https://swe-agent.com/latest/)
- [SWE-agent GitHub](https://github.com/SWE-agent/SWE-agent)
- [Cline GitHub](https://github.com/cline/cline)
- [Claude Code Hooks Guide](https://code.claude.com/docs/en/hooks-guide)
