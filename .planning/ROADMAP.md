# Roadmap: Kamaji

## Milestones

- **v1.0 MVP** - Phases 1-10 (shipped 2026-01-23)
- **v1.1 Sprint Planning** - Phases 11-14 (in progress)

## Phases

<details>
<summary>v1.0 MVP (Phases 1-10) - SHIPPED 2026-01-23</summary>

23 plans across 10 phases. Core orchestration, state machine, MCP server, Claude Code spawning, git operations, streaming output. See git history for details.

</details>

### v1.1 Sprint Planning (In Progress)

**Milestone Goal:** Add commands to help users create and refine sprint configs before execution.

- [ ] **Phase 11: Validation** - Schema and semantic validation for kamaji.yaml
- [ ] **Phase 12: Init** - Create template config with explanatory comments
- [ ] **Phase 13: Refine** - AI-assisted config improvement via Claude Code
- [ ] **Phase 14: Integration** - Wire validation into start command

## Phase Details

### Phase 11: Validation

**Goal**: Users can validate their kamaji.yaml config with clear feedback
**Depends on**: Phase 10 (v1.0 complete)
**Requirements**: VALD-01, VALD-02, VALD-03, VALD-04
**Success Criteria** (what must be TRUE):

1. User can run `kamaji validate` and see pass/fail result
2. Invalid YAML structure produces error with line number
3. Missing required fields produce errors naming the field
4. Empty descriptions and invalid dependency references are caught
5. All errors include enough context to locate and fix the problem
   **Plans**: TBD

Plans:

- [ ] 11-01: TBD
- [ ] 11-02: TBD

### Phase 12: Init

**Goal**: Users can bootstrap a new kamaji.yaml with a single command
**Depends on**: Phase 10 (independent of Phase 11)
**Requirements**: INIT-01, INIT-02, INIT-03
**Success Criteria** (what must be TRUE):

1. User can run `kamaji init` in any directory and get kamaji.yaml created
2. Generated file contains valid YAML with comments explaining each section
3. Running `kamaji init` when kamaji.yaml exists shows error and does not overwrite
   **Plans**: TBD

Plans:

- [ ] 12-01: TBD

### Phase 13: Refine

**Goal**: Users can improve their config by having Claude analyze and rewrite it
**Depends on**: Phase 11 (validation used before/after AI processing)
**Requirements**: REFN-01, REFN-02, REFN-03, REFN-04, REFN-05, REFN-06, REFN-07, REFN-08, REFN-09
**Success Criteria** (what must be TRUE):

1. User can run `kamaji refine` and see Claude spawned to process config
2. Config is validated before Claude processes it (invalid configs rejected)
3. AI improves vague tasks, fills missing structure, uses comments as context
4. AI adds clarifying questions as YAML comments when information is missing
5. Modified config is written in place and passes validation
   **Plans**: TBD

Plans:

- [ ] 13-01: TBD
- [ ] 13-02: TBD
- [ ] 13-03: TBD

### Phase 14: Integration

**Goal**: Start command validates config before running to catch errors early
**Depends on**: Phase 11 (uses validation logic)
**Requirements**: INTG-01
**Success Criteria** (what must be TRUE):

1. Running `kamaji start` with invalid config shows validation errors
2. Running `kamaji start` with valid config proceeds to sprint execution
   **Plans**: TBD

Plans:

- [ ] 14-01: TBD

## Progress

**Execution Order:** 11 -> 12 -> 13 -> 14 (12 can run parallel to 11 if desired)

| Phase           | Milestone | Plans Complete | Status      | Completed  |
| --------------- | --------- | -------------- | ----------- | ---------- |
| 1-10            | v1.0      | 23/23          | Complete    | 2026-01-23 |
| 11. Validation  | v1.1      | 0/2            | Not started | -          |
| 12. Init        | v1.1      | 0/1            | Not started | -          |
| 13. Refine      | v1.1      | 0/3            | Not started | -          |
| 14. Integration | v1.1      | 0/1            | Not started | -          |

---

_Roadmap created: 2026-01-23_
_Last updated: 2026-01-23_
