## Tasks - command-catalog-cli

### 1) Objective

Break MVP delivery into atomic, executable tasks with explicit dependencies, done criteria, and release milestones for small PRs.

### 2) Conventions

- Initial status: `TODO`
- Priority levels: `P0`, `P1`, `P2`
- Task IDs: `T001`, `T002`, ...
- All tasks should be implementable in one focused PR.

### 3) Implementation Backlog

#### T001 - Bootstrap Go CLI project structure

- **Priority:** P0
- **Dependencies:** none
- **Scope:** Create baseline folders (`cmd/cs`, `internal/*`, `tests/*`), initialize module, and wire root CLI entrypoint.
- **Done when:** `cs` binary builds successfully and root help command runs.

#### T002 - Define core domain contracts and typed errors

- **Priority:** P0
- **Dependencies:** T001
- **Scope:** Implement domain types (`CatalogFile`, `CatalogCommand`, `ExecutionPlan`, `ExecutionResult`) and shared typed errors (`not found`, `validation`, `ambiguous match`, `danger denied`).
- **Done when:** Contracts compile and unit tests validate basic invariants/serialization tags.

#### T003 - Implement JSON catalog repository (load/save)

- **Priority:** P0
- **Dependencies:** T001, T002
- **Scope:** Build storage layer to load/save catalog JSON, initialize empty catalog shape, and return recovery-oriented error on invalid JSON.
- **Done when:** Repository integration tests cover load existing, load missing, and invalid JSON behavior.

#### T004 - Add atomic file write strategy for catalog mutations

- **Priority:** P0
- **Dependencies:** T003
- **Scope:** Implement temp file + fsync + atomic rename write path for `Save` and mutation operations.
- **Done when:** Tests verify file remains valid after mutation path and no partial write is produced in normal failure simulation.

#### T005 - Implement `create` validator and use case

- **Priority:** P0
- **Dependencies:** T002, T003, T004
- **Scope:** Validate `key`/`value`, generate unique `id`, persist new command, and return clear output/errors.
- **Done when:** Unit/integration tests cover valid create, empty input rejection, persistence failure handling.

#### T006 - Implement `list` use case and formatter

- **Priority:** P0
- **Dependencies:** T003
- **Scope:** Load catalog and render deterministic output with `id`, `key`, `value`, including empty state.
- **Done when:** Tests verify output schema and empty catalog response.

#### T007 - Implement `delete` validator and use case

- **Priority:** P0
- **Dependencies:** T003, T004
- **Scope:** Validate `id`, delete exact entry, persist updated catalog, return not-found error for missing id.
- **Done when:** Tests cover successful delete, missing id input, and unknown id behavior.

#### T008 - Build placeholder extractor and binder service

- **Priority:** P0
- **Dependencies:** T002
- **Scope:** Parse placeholders from template values and bind runtime args into a resolved command string.
- **Done when:** Unit tests cover single placeholder, multiple placeholders, and missing args rejection.

#### T009 - Implement best-match resolver (runtime)

- **Priority:** P0
- **Dependencies:** T002, T003
- **Scope:** Resolve input key tokens against stored commands with deterministic scoring and deterministic tie handling.
- **Done when:** Unit tests cover exact match, best-match, tie/disambiguation error, and not-found error.

#### T010 - Implement command execution runner

- **Priority:** P0
- **Dependencies:** T008, T009
- **Scope:** Execute resolved shell command as child process and return exit code/stdout/stderr consistently.
- **Done when:** Integration tests verify success/failure exit code propagation and output capture.

#### T011 - Implement dangerous-command classifier

- **Priority:** P0
- **Dependencies:** T002
- **Scope:** Add initial dangerous pattern rules and classification service returning reasons.
- **Done when:** Unit tests validate expected dangerous/non-dangerous classification for baseline pattern set.

#### T012 - Implement interactive confirmation gate

- **Priority:** P0
- **Dependencies:** T010, T011
- **Scope:** Prompt user with explicit confirmation (`[y/N]`) for dangerous commands and block execution on denial.
- **Done when:** Integration tests cover accepted flow, denied flow, and prompt text contract.

#### T013 - Implement non-interactive safe-fail behavior

- **Priority:** P0
- **Dependencies:** T012
- **Scope:** Detect non-interactive context and abort dangerous execution with actionable error unless explicit override contract is defined.
- **Done when:** Tests validate safe abort behavior in non-interactive mode.

#### T014 - Wire runtime `cs <key...> [args...]` command flow

- **Priority:** P0
- **Dependencies:** T008, T009, T010, T012, T013
- **Scope:** Connect parser, resolver, binder, danger policy, confirmer, and runner into end-to-end run path.
- **Done when:** CLI integration tests pass for happy path and key error branches.

#### T015 - Implement shell completion command output

- **Priority:** P1
- **Dependencies:** T001, T009
- **Scope:** Add `cs completion <shell>` command to emit installable completion scripts (bash first).
- **Done when:** Script generation test passes and installation instructions are documented.

#### T016 - Implement TAB suggestion provider with best-match autofill contract

- **Priority:** P1
- **Dependencies:** T009, T015
- **Scope:** Provide completion-time suggestion logic aligned with runtime scoring and deterministic tie-break behavior.
- **Done when:** Completion tests verify best suggestion and deterministic tie handling.

#### T017 - Implement CLI error model and exit code standardization

- **Priority:** P0
- **Dependencies:** T002, T014
- **Scope:** Map typed errors to consistent stderr messages and exit codes for validation, not-found, ambiguous, and danger-denied cases.
- **Done when:** Contract tests assert stable exit codes/messages for each error type.

#### T018 - Add debug logging hooks (local only)

- **Priority:** P2
- **Dependencies:** T014
- **Scope:** Add optional debug logs via env/flag for key events (`catalog_loaded`, `command_created`, `command_deleted`, `match_resolved`, `danger_confirmation_prompted`, `command_executed`).
- **Done when:** Logging can be toggled without affecting default output contracts.

#### T019 - Build test suite coverage for MVP stories

- **Priority:** P0
- **Dependencies:** T005, T006, T007, T014, T016, T017
- **Scope:** Add/complete unit + integration + CLI contract tests to cover US01-US06 acceptance criteria.
- **Done when:** Automated tests pass and each user story has at least one explicit test mapping.

#### T020 - Prepare docs and release checklist for MVP

- **Priority:** P1
- **Dependencies:** T015, T019
- **Scope:** Document install/build, catalog path, command usage, completion setup, confirmation behavior, and known limitations.
- **Done when:** README/docs are updated and MVP release checklist is complete.

### 4) Requirement-to-Task Mapping

- FR01 (create command): T005, T019
- FR02 (list commands): T006, T019
- FR03 (delete by id): T007, T019
- FR04 (execute with args): T008, T009, T010, T014, T019
- FR05 (input validation): T005, T007, T008, T017, T019
- FR06 (local persistence): T003, T004, T019
- FR07 (TAB completion): T015, T016, T019
- FR08 (danger confirmation): T011, T012, T013, T019
- FR09 (disambiguation errors): T009, T017, T019
- NFR01 (deterministic/readable CLI output): T006, T017, T019
- NFR02 (catalog operation efficiency): T003, T006, T007, T019
- NFR03 (actionable errors and non-zero exits): T017, T019
- NFR04 (atomic writes): T004, T019
- NFR05 (Linux MVP operation): T001, T014, T019
- NFR06 (deterministic matching/completion): T009, T016, T019

### 5) Recommended Execution Order

1. Foundation: T001 -> T002 -> T003 -> T004
2. Catalog CRUD: T005 -> T006 -> T007
3. Execution core: T008 -> T009 -> T010
4. Safety flow: T011 -> T012 -> T013
5. Runtime wiring and contracts: T014 -> T017
6. Completion: T015 -> T016
7. Quality and release: T019 -> T020
8. Optional enhancement: T018

### 6) MVP Release Milestones

- **M1 - Catalog Base:** T001, T002, T003, T004, T005, T006, T007
- **M2 - Executable Commands:** T008, T009, T010, T014, T017
- **M3 - Safety Controls:** T011, T012, T013
- **M4 - Completion UX:** T015, T016
- **M5 - Hardening and Release:** T019, T020

### 7) Operational Risks

- Matching algorithm changes can break completion and runtime parity.
- Incomplete dangerous pattern list can allow risky commands without prompt.
- Shell-specific completion differences can impact user setup experience.
- Error message drift can break CLI contract tests and script integrations.

### 8) SDD Completion Criteria

- Files `01-context.md`, `02-requirements.md`, `03-design.md`, and `04-tasks.md` are aligned and approved.
- Every MVP requirement maps to at least one implementation task.
- All `P0` tasks have owner, estimate, and planned sprint assignment.
