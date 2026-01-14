# Task format

How to structure task definitions for autonomous execution.

## Core principle

Each task needs explicit action, verification, and completion criteria. The separation enables machine-checkable task completion—the orchestrator runs the verify command and compares output against the done condition.

## Task schema

From GSD's XML format, adapted for kamaji's YAML:

```yaml
tasks:
    - description: "Create login endpoint"
      files:
          - src/app/api/auth/login/route.ts
      steps:
          - "Use jose for JWT (not jsonwebtoken - CommonJS issues)"
          - "Validate credentials against users table"
          - "Return httpOnly cookie on success"
      verify: "curl -X POST localhost:3000/api/auth/login returns 200 + Set-Cookie"
      done: "Valid credentials return cookie, invalid return 401"
```

## Element definitions

| Element       | Purpose                                          |
| ------------- | ------------------------------------------------ |
| `description` | Task title—what to accomplish                    |
| `files`       | Target file paths (optional, helps Claude focus) |
| `steps`       | Implementation instructions                      |
| `verify`      | Command to run (the HOW)                         |
| `done`        | Expected output (the WHAT)                       |

The `verify` field should be a machine-checkable command when possible. The `done` field defines the exit condition.

## Verify vs done separation

**Problem with single verify field:**

```yaml
verify: "Component renders without errors"
```

This describes what to check but not the explicit pass/fail signal.

**With verify/done separation:**

```yaml
verify: "npm test -- LoginForm"
done: "All tests pass with 0 failures"
```

The orchestrator can run `verify` and match output against `done`.

## Context injection

Tasks can specify files Claude should read before execution:

```yaml
tickets:
    - name: login-form
      context:
          - "@src/components/auth/README.md"
          - "@src/hooks/useAuth.ts"
      tasks: [...]
```

Purpose:

1. Reduce discovery overhead—Claude doesn't hunt for patterns
2. Ensure consistency—Claude sees the right examples
3. Inject previous work summaries for continuity

## Dependencies

Tasks can declare dependencies for parallel execution:

```yaml
tickets:
    - name: auth-context
      branch: feat/auth-context
      tasks: [...]

    - name: protected-routes
      branch: feat/protected-routes
      depends_on: [auth-context]
      tasks: [...]
```

Independent tickets (`depends_on: []`) execute in parallel. Dependencies enforce ordering.

## Cascading rules

Rules cascade from sprint to ticket to task:

```yaml
rules:
    - "Follow existing patterns in the codebase"

tickets:
    - name: login-form
      rules:
          - "Use Zod for validation (not yup)"
      tasks:
          - description: "Create LoginForm"
            rules:
                - "Match existing form components in src/components/forms/"
```

More specific rules override general ones without duplication.

## Task atomicity

Each task should:

- Complete in reasonable context (won't exhaust the window)
- Produce a meaningful commit
- Have binary pass/fail verification

**Bad—too large:**

```yaml
tasks:
    - description: "Implement entire authentication system"
```

**Good—atomic subtasks:**

```yaml
tasks:
    - description: "Create User model with password hashing"
      done: "npm test -- User.test.ts passes"
    - description: "Add login endpoint"
      done: "curl -X POST /login returns 200"
```

## Output specification

Tasks can specify required artifacts:

```yaml
tasks:
    - description: "Create LoginForm component"
      verify: "npm test -- LoginForm"
      done: "All tests pass with 0 failures"
      output: "Document usage in src/components/auth/README.md"
```

---

## Kamaji implications

| Pattern                | Application                                                           |
| ---------------------- | --------------------------------------------------------------------- |
| Verify/done separation | Enables orchestrator verification independent of Claude's self-report |
| Context injection      | Add `context` field to tickets                                        |
| Dependencies           | Add `depends_on` field for parallel execution                         |
| Cascading rules        | Allow rules at sprint, ticket, and task levels                        |
| Atomic tasks           | Design tasks that survive context exhaustion                          |

---

## Sources

- [glittercowboy/get-shit-done](https://github.com/glittercowboy/get-shit-done)
