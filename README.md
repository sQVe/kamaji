# Kamaji

> [!WARNING]
> This project is under active development and not yet ready for use.

External Go CLI that orchestrates autonomous coding sprints by spawning fresh AI agent sessions per task.

## Why

Running loops inside AI agents causes context pollution and exhaustion. Kamaji owns the state machine externally, giving each task a clean context window.

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
