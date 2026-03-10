## Feature Context: command-catalog-cli

### 1) Overview

This feature creates a local CLI command catalog so a user can define reusable command aliases with parameters and execute them as custom commands.
The main problem it solves is repetitive shell usage for common operational tasks (for example, killing a process by port).

### 2) MVP Objective

Deliver a CLI tool (`cs`) that allows users to create, list, run, and delete custom commands stored locally, with parameter substitution at execution time, focused on Linux usage.

### 3) MVP Scope

- Provide `cs create "<key>" "<value>"` to register a new command template.
- Persist commands in a local JSON file-based catalog.
- Provide `cs list` to display all stored commands (`id`, `key`, `value`).
- Provide `cs delete <id>` to remove a command.
- Allow execution by key pattern with arguments (example: `cs kill port 3040`).
- Support parameter placeholders in command values (example: `$port`).
- Provide shell completion so TAB suggests and autofills the best command key match.
- Require explicit confirmation before executing dangerous commands.

### 4) Out of Scope (Non-Goals)

- Remote/shared catalog across machines.
- GUI/web interface.
- Full cross-platform support parity (Windows/macOS) in MVP.
- Role-based access control or multi-user permissions.
- Command history analytics and usage telemetry.
- In-place command update/edit operation.

### 5) Assumptions and Constraints

- Primary OS target is Linux.
- Runtime/language is Golang.
- Storage must be local file-based JSON (no database for MVP).
- Command model keeps `id`, `key`, and `value` as minimum required fields.
- Execution runs shell commands on user machine and inherits user permissions.
- Update flow in MVP is delete + create new command.
- Shell completion integration must be supported for TAB key behavior.

### 6) User Journey (Flow)

1. User identifies a repetitive shell command they want to reuse.
2. User creates a catalog entry with key and command template:
   `cs create "kill port $port" "sudo kill -9 $(sudo lsof -t -i:$port)"`.
3. System validates and stores the command in local JSON catalog.
4. User types a partial command key and presses TAB.
5. Shell completion suggests and autofills the best key match.
6. User runs the command with runtime args:
   `cs kill port 3040`.
7. System resolves placeholders, requests confirmation for dangerous commands, and executes.
8. User lists commands with `cs list` to inspect catalog.
9. User removes obsolete command with `cs delete <id>`.

### 7) Business Rules (MVP)

- BR01: Every command must have a unique `id`.
- BR02: A command entry must include non-empty `key` and `value`.
- BR03: `cs list` returns at least `id`, `key`, and `value` for each entry.
- BR04: `cs delete <id>` removes only the targeted entry.
- BR05: Execution resolves placeholders in the stored template using provided runtime args.
- BR06: If command key cannot be resolved, system returns a clear error.
- BR07: If required placeholder args are missing, system does not execute and returns validation error.
- BR08: Catalog changes (create/delete) persist immediately to local JSON storage.
- BR09: Dangerous commands require explicit confirmation before execution.
- BR10: Command update is not available in MVP; user must delete and recreate.

### 8) Initial Project Structure (Proposed)

- `cmd/cs/*`: CLI entrypoint and command wiring.
- `internal/cli/*`: argument parsing and command handlers (`create`, `list`, `delete`, `run`).
- `internal/domain/*`: command entity and validation rules.
- `internal/storage/*`: local JSON repository (read/write/atomic replace).
- `internal/execution/*`: placeholder resolution, danger detection, and process execution.
- `internal/completion/*`: shell completion generation and best-match selection.
- `internal/types/*`: contracts for command records and execution requests.
- `tests/*`: unit and integration tests for CLI behavior.

### 9) Data Model (Initial Contracts)

- `CatalogCommand`
  - `id: string`
  - `key: string`
  - `value: string`
  - `createdAt?: string`
  - `updatedAt?: string`
- `CatalogFile`
  - `version: string`
  - `commands: CatalogCommand[]`
- `ExecutionRequest`
  - `inputTokens: string[]`
- `ExecutionPlan`
  - `commandId: string`
  - `resolvedCommand: string`
  - `isDangerous: bool`
- `ExecutionResult`
  - `exitCode: int`
  - `stdout?: string`
  - `stderr?: string`

### 10) Functional Requirements Summary

- FR01: User can create catalog commands with template placeholders.
- FR02: User can list all catalog commands.
- FR03: User can delete command by id.
- FR04: User can execute a command via key-based invocation with args.
- FR05: System validates required inputs before storage and execution.
- FR06: System persists catalog locally between sessions.
- FR07: System supports shell TAB completion that suggests and autofills best key match.
- FR08: System asks for explicit confirmation before dangerous command execution.
- FR09: System returns clear disambiguation errors for unresolved or ambiguous command selection.

### 11) Non-Functional Requirements Summary

- NFR01: CLI interactions provide deterministic and human-readable output.
- NFR02: Common operations (`create`, `list`, `delete`) complete quickly for small/medium catalogs.
- NFR03: Failures return non-zero exit code and actionable error text.
- NFR04: Storage operations use an atomic write strategy to reduce corruption risk.
- NFR05: MVP runs on Linux shell environments.
- NFR06: Completion and matching behavior is deterministic for same input/catalog state.

### 12) Acceptance Criteria (By Journey Step)

- AC01: Given a valid `create` input, when user runs `cs create ...`, then a new entry is persisted with `id`, `key`, and `value` in JSON catalog.
- AC02: Given stored entries, when user runs `cs list`, then all entries are shown including `id`, `key`, and `value`.
- AC03: Given an existing id, when user runs `cs delete <id>`, then only that entry is removed and no longer appears in `cs list`.
- AC04: Given a command template with `$port`, when user runs `cs kill port 3040`, then `$port` resolves to `3040` in executed command.
- AC05: Given missing required runtime arguments, when execution is attempted, then command is not executed and validation error is shown.
- AC06: Given unknown command key, when execution is attempted, then system returns command-not-found style error with non-zero exit.
- AC07: Given a partial key, when user presses TAB, then shell completion suggests and autofills the best command key match.
- AC08: Given a command detected as dangerous, when user executes it, then system requests explicit confirmation before running.
- AC09: Given a dangerous command and user declines confirmation, when prompt is answered negatively, then command is not executed.

### 13) Risks and Mitigations

- Risk: Unsafe command execution may run destructive shell instructions.
  - Mitigation: dangerous-command detection + explicit confirmation gate.
- Risk: Ambiguous key matching may execute unintended command.
  - Mitigation: deterministic best-match scoring and explicit ambiguity errors.
- Risk: Placeholder parsing errors can produce invalid commands.
  - Mitigation: validate placeholders at create time and before execution.
- Risk: Local file corruption due to interrupted writes.
  - Mitigation: write temp file + atomic replace strategy.
- Risk: Shell completion behavior differs by shell.
  - Mitigation: provide shell-specific completion scripts and tested install steps.

### 14) Decisions and Open Items

- D01: Runtime/language is Golang.
- D02: Catalog persistence format is JSON (local file).
- D03: Shell completion suggests and autofills the best command key when TAB is pressed.
- D04: Dangerous commands require explicit user confirmation before execution.
- D05: Command update/edit is out of MVP; users must delete and create a new command.
- OD01: Define exact best-match algorithm (score rules and tie-breakers).
- OD02: Define dangerous-command detection strategy (initial pattern rules for MVP).
- OD03: Define confirmation UX contract for interactive and non-interactive contexts.
- OD04: Define key uniqueness policy (`exact unique` vs `allow overlapping patterns`).
- OD05: Define initial completion shell scope (`bash only` vs `bash + zsh`).
