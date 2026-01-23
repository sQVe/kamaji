# Requirements: Kamaji

**Defined:** 2026-01-23
**Core Value:** Reliable state machine. Never lose progress, survive crashes, always know exactly where you are in the sprint.

## v1.1 Requirements

Requirements for sprint planning feature. Each maps to roadmap phases.

### Init Command

- [ ] **INIT-01**: User can run `kamaji init` to create template config in CWD
- [ ] **INIT-02**: Template contains minimal scaffold with comments explaining each section
- [ ] **INIT-03**: Init fails gracefully if kamaji.yaml already exists

### Refine Command

- [ ] **REFN-01**: User can run `kamaji refine` to improve config with AI
- [ ] **REFN-02**: Refine spawns Claude Code to read and rewrite config
- [ ] **REFN-03**: AI uses YAML comments as context for improvements
- [ ] **REFN-04**: AI breaks down vague tasks into concrete ones
- [ ] **REFN-05**: AI fills missing structure and improves descriptions
- [ ] **REFN-06**: Refine validates config before AI processing
- [ ] **REFN-07**: Refine validates config after AI processing
- [ ] **REFN-08**: Config is written in place (no diff/approval flow)
- [ ] **REFN-09**: AI writes questions as comments when clarification needed

### Validation

- [ ] **VALD-01**: User can run `kamaji validate` for one-off config check
- [ ] **VALD-02**: Validation checks YAML schema (structure, required fields, types)
- [ ] **VALD-03**: Validation checks semantic heuristics (no empty descriptions, deps exist)
- [ ] **VALD-04**: Validation returns clear error messages with locations

### Integration

- [ ] **INTG-01**: Start command validates config before running sprint

## Future Requirements

Deferred to later milestones.

### Parallel Execution (v2)

- **PARA-01**: Multiple tickets can execute concurrently
- **PARA-02**: Worktree support for parallel isolation
- **PARA-03**: Service management (status/stop/retry)

## Out of Scope

Explicitly excluded from this milestone.

| Feature                       | Reason                                          |
| ----------------------------- | ----------------------------------------------- |
| Interactive refine prompts    | Config comments serve as communication medium   |
| Diff/approval flow for refine | Keeps workflow simple, user can git diff        |
| Multiple config file support  | One kamaji.yaml per project sufficient for v1.1 |
| Config migration tooling      | No breaking changes to config format            |

## Traceability

Which phases cover which requirements. Updated during roadmap creation.

| Requirement | Phase    | Status  |
| ----------- | -------- | ------- |
| INIT-01     | Phase 12 | Pending |
| INIT-02     | Phase 12 | Pending |
| INIT-03     | Phase 12 | Pending |
| REFN-01     | Phase 13 | Pending |
| REFN-02     | Phase 13 | Pending |
| REFN-03     | Phase 13 | Pending |
| REFN-04     | Phase 13 | Pending |
| REFN-05     | Phase 13 | Pending |
| REFN-06     | Phase 13 | Pending |
| REFN-07     | Phase 13 | Pending |
| REFN-08     | Phase 13 | Pending |
| REFN-09     | Phase 13 | Pending |
| VALD-01     | Phase 11 | Pending |
| VALD-02     | Phase 11 | Pending |
| VALD-03     | Phase 11 | Pending |
| VALD-04     | Phase 11 | Pending |
| INTG-01     | Phase 14 | Pending |

**Coverage:**

- v1.1 requirements: 17 total
- Mapped to phases: 17
- Unmapped: 0

---

_Requirements defined: 2026-01-23_
_Last updated: 2026-01-23 after roadmap creation_
