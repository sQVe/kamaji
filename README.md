# Kamaji

> [!WARNING]
> This project is under active development and not yet ready for use.

External Go CLI that orchestrates autonomous coding sprints by spawning fresh AI agent sessions per task.

## Why

**AI agents are unreliable planners.** When AI decides what tasks to execute, it misses things. Tests, edge cases, conventions. You only discover the gaps after the work is done.

**Kamaji flips the model.** You define the tasks. The AI executes but doesn't decide scope. Nothing gets skipped unless you consciously omit it.

**Opinionated by design.** Kamaji provides structure: research, implementation, testing, verification. You decide what applies per task. Good practices are built into the workflow, not left to the AI's judgment.

**Context stays clean.** Running loops inside AI agents causes context pollution. Kamaji owns the state machine externally, spawning fresh sessions per task.

## Features

- State machine that survives crashes
- Fresh agent session per task
- MCP server for completion signals
- Git operations (branch, commit, reset)
- Per-ticket history and insights

## Installation

```sh
go install github.com/sqve/kamaji/cmd/kamaji@latest
```

## Usage

```sh
kamaji start
```

## License

[MIT](LICENSE)
