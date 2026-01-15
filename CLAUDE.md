# CLAUDE.md

## Project Context

Kamaji orchestrates autonomous coding sprints by managing state and spawning fresh AI agent sessions per task to prevent context pollution.

**Essential reading:**

- [Design](DESIGN.md) — Architecture and rationale
- [Contributing](CONTRIBUTING.md) — Setup and standards

## Commands

- `make test` — run unit tests
- `make lint` — lint and fix violations
- `make build` — build binary
- `make ci` — run full CI pipeline

## Guidelines

- Follow existing patterns in `internal/`
- Domain types in `internal/domain/` hold data only—no I/O
- Config package loads and saves all files
- Add changelog entry via `make change` for significant changes
