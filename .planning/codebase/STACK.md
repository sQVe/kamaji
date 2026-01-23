# Technology Stack

**Analysis Date:** 2026-01-23

## Languages

**Primary:**

- Go 1.25 - Core language for the entire project
- YAML - Configuration format for sprint definitions and state

**Secondary:**

- Shell/Bash - For Makefile targets and git operations
- JSON - For MCP configuration files (.mcp.json)

## Runtime

**Environment:**

- Go runtime 1.25.5 (CI/CD uses this)
- Unix-like systems (Linux, macOS) and Windows support

**Package Manager:**

- Go modules (go.mod, go.sum)
- Lockfile: Present (go.sum)

## Frameworks

**Core:**

- Cobra v1.9.1 - CLI framework for command handling (`github.com/spf13/cobra`)

**MCP Server:**

- mcp-go v0.43.2 - Model Context Protocol server implementation (`github.com/mark3labs/mcp-go`)

**Testing:**

- Go standard library testing
- gotestsum v1.13.0 - Test runner with better output formatting
- Integration tags for categorized test runs

**Build/Dev:**

- GoReleaser v2 - Multi-platform binary builds and releases
- golangci-lint v2.8.0 - Code linting and static analysis
- Changie v1.24.0 - Changelog management
- pre-commit with golangci-lint hooks

## Key Dependencies

**Critical:**

- `github.com/mark3labs/mcp-go/mcp` v0.43.2 - MCP protocol implementation for agent communication
- `github.com/spf13/cobra` v1.9.1 - CLI command structure
- `gopkg.in/yaml.v3` v3.0.1 - YAML parsing for sprint/state configs

**Infrastructure:**

- `github.com/charmbracelet/lipgloss` v1.0.0 - Terminal UI styling for output
- `github.com/rogpeppe/go-internal` v1.14.1 - Go internal utilities
- `golang.org/x/tools` v0.26.0 - Go tooling (deadcode checker)
- `golang.org/x/sys` v0.26.0 - System-level operations

**Build/JSON Schema:**

- `github.com/invopop/jsonschema` v0.13.0 - JSON schema generation (used by mcp-go)
- `github.com/google/uuid` v1.6.0 - UUID generation
- `github.com/buger/jsonparser` v1.1.1 - JSON parsing utilities

## Configuration

**Environment:**

- Git configuration for CI: email, name, default branch, GPG signing disabled
- KAMAJI_MCP_PORT - Port for MCP server (dynamic or override)
- KAMAJI_WORK_DIR - Working directory for agent execution

**Build:**

- `.golangci.yml` - Linter configuration with strict rules
- `go.mod` - Module definition and version constraints
- `.goreleaser.yml` - Release build configuration for Linux/macOS (amd64, arm64)
- `Makefile` - Standard build targets (test, lint, build, ci, clean)
- `.pre-commit-config.yaml` - Pre-commit hooks for golangci-lint and prettier

**Runtime Artifacts:**

- `kamaji.yaml` - Sprint definition (user-provided, loaded from working directory)
- `.kamaji/state.yaml` - State machine persistence in project directory
- `.kamaji/history/*.yaml` - Per-ticket history files with task completion/failure records
- `.mcp.json` - Generated MCP server configuration for Claude Code connections

## Platform Requirements

**Development:**

- Go 1.25+ toolchain
- git (for git operations and CI/CD)
- npx/prettier (for code formatting in Makefile)
- golangci-lint v2.8.0+ (for linting)
- gotestsum v1.13.0+ (for test runs)
- changie v1.24.0+ (for changelog management)

**Production:**

- Unix-like OS (Linux, macOS) or Windows
- Git repository (for branch/commit operations)
- Claude Code CLI installed and in PATH (for spawning agents)
- Network access to MCP HTTP server (localhost, dynamic ports)

## Release & Deployment

**Distribution:**

- GoReleaser builds binaries for:
    - Linux: amd64, arm64
    - macOS: amd64, arm64
- Package formats:
    - tar.gz archives
    - Debian (.deb)
    - RPM (.rpm)
- GitHub Releases with checksums
- Installation via `go install github.com/sqve/kamaji/cmd/kamaji@latest`

**CI/CD Platforms:**

- GitHub Actions:
    - Matrix testing across ubuntu-latest, windows-latest, macos-latest
    - Linting via golangci-lint-action
    - GoReleaser action for release builds
    - Conditional changesets for PR validation

---

_Stack analysis: 2026-01-23_
