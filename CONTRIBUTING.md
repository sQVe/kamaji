# Contributing

## Setup

```bash
go mod download
```

**Prerequisites:** Go 1.25+, golangci-lint

## Development

- `make test` — run unit tests
- `make lint` — lint and fix violations
- `make build` — build binary
- `make ci` — run full CI pipeline

## Pull requests

1. Branch from `main`
2. Add a changelog entry via `make change`
3. Verify `make ci` passes
4. Open a PR to `main`

## Changelog

| Change type             | Changeset required?         |
| ----------------------- | --------------------------- |
| New feature             | Yes                         |
| Bug fix                 | Yes                         |
| Breaking change         | Yes                         |
| Performance improvement | Yes                         |
| Documentation only      | No                          |
| CI/workflow changes     | No                          |
| Test-only changes       | No                          |
| Dependency updates      | No (`dependencies` label)   |
| Internal refactors      | No (`skip-changelog` label) |
