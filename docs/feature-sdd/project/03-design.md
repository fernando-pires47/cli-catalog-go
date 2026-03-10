## Design - command-catalog-cli

### 1) Objective

Define an implementation-oriented technical design for a Golang CLI that manages a local JSON command catalog, resolves parameterized commands, supports TAB completion with deterministic best-match behavior, and enforces confirmation for dangerous execution.

### 2) Architecture Principles

- Keep domain rules pure and testable (validation, matching, danger classification).
- Isolate side effects behind interfaces (file system, process execution, terminal I/O).
- Prefer deterministic behavior over implicit guesses.
- Fail safe on ambiguous match, invalid args, and dangerous non-interactive execution.
- Keep MVP small: add only capabilities required by US01-US06.

### 3) Layer View

- `cmd/cs/main.go`: process bootstrap, root command setup, exit code handling.
- `internal/cli/*`: argument parsing and command handlers (`create`, `list`, `delete`, `run`, `completion`).
- `internal/domain/*`: entities and business rules (`CatalogCommand`, validation, matching score, danger rules).
- `internal/storage/*`: JSON repository, file lock/write strategy, load/save operations.
- `internal/execution/*`: placeholder binding, command plan builder, confirmation gate, shell process runner.
- `internal/completion/*`: shell completion adapter and best-match suggestion provider.
- `internal/contracts/*`: interfaces for repository, confirmer, command runner, clock/id generator.
- `tests/*`: unit, integration, and CLI contract tests.

### 4) Commands and Responsibilities

- `cs create "<key>" "<value>"`
  - Validate required args.
  - Parse placeholder tokens from `value`.
  - Persist command in JSON catalog with generated `id`.
- `cs list`
  - Load catalog and print table/list with `id`, `key`, `value`.
- `cs delete <id>`
  - Validate `id`, remove matching command, persist result.
- `cs <key tokens...> [arg tokens...]`
  - Resolve command by best-match strategy.
  - Bind runtime args to placeholders.
  - Apply dangerous-command confirmation when needed.
  - Execute resolved shell command and return subprocess exit code.
- `cs completion <bash|zsh>`
  - Emit shell completion script for installation.
  - Completion engine suggests/autofills best key match on TAB.

### 5) State Model

#### 5.1 Global Persistent State

- Source: local JSON file (example path: `$HOME/.cs/catalog.json`).
- Shape:
  - `version: string`
  - `commands: CatalogCommand[]`
- Invariants:
  - JSON must remain valid after each successful mutation.
  - `id` must be unique.

#### 5.2 Runtime Local State

- Parsed CLI command intent (`op`, `rawArgs`).
- Match candidates and score matrix.
- Placeholder binding map (example: `port -> 3040`).
- Execution plan (`resolvedCommand`, `isDangerous`, `requiresConfirmation`).
- Confirmation result (`accepted`, `rejected`, `nonInteractiveAbort`).

### 6) Data Contracts (Go)

```go
type CatalogFile struct {
    Version  string           `json:"version"`
    Commands []CatalogCommand `json:"commands"`
}

type CatalogCommand struct {
    ID        string `json:"id"`
    Key       string `json:"key"`
    Value     string `json:"value"`
    CreatedAt string `json:"createdAt,omitempty"`
    UpdatedAt string `json:"updatedAt,omitempty"`
}

type MatchCandidate struct {
    Command CatalogCommand
    Score   int
}

type ExecutionPlan struct {
    CommandID            string
    ResolvedCommand      string
    IsDangerous          bool
    RequiresConfirmation bool
}

type ExecutionResult struct {
    ExitCode int
    Stdout   string
    Stderr   string
}
```

### 7) Service Design and Interfaces

- `CatalogRepository`
  - `Load(ctx) (CatalogFile, error)`
  - `Save(ctx, catalog CatalogFile) error`
  - `Create(ctx, key, value string) (CatalogCommand, error)`
  - `DeleteByID(ctx, id string) error`
- `MatcherService`
  - `Resolve(inputTokens []string, commands []CatalogCommand) (CatalogCommand, error)`
  - `Suggest(prefixTokens []string, commands []CatalogCommand) ([]CatalogCommand, error)`
- `TemplateBinder`
  - `ExtractPlaceholders(template string) ([]string, error)`
  - `Bind(template string, args []string) (string, error)`
- `DangerPolicy`
  - `Classify(resolvedCommand string) (isDangerous bool, reasons []string)`
- `Confirmer`
  - `ConfirmDanger(ctx, prompt string) (bool, error)`
- `CommandRunner`
  - `Run(ctx, command string) (ExecutionResult, error)`

### 8) Business Rules in Design

- BR01/BR02: enforced by `CreateValidator` before persistence.
- BR03: enforced by `ListFormatter` contract.
- BR04: enforced by repository delete operation with exact id match.
- BR05/BR07: enforced by `TemplateBinder` and runtime arg validator.
- BR06: enforced by matcher returning typed `ErrNotFound`.
- BR08: enforced by repository `Save` with atomic replace.
- BR09: enforced by `DangerPolicy` + `Confirmer` gate before `CommandRunner`.
- BR10: no `update` command exposed in CLI.

### 9) Async and Process Strategy

- CLI operations are synchronous in request/response style.
- No polling or background worker in MVP.
- Execution command runs as child process; stdout/stderr streamed or buffered for terminal output.
- Interruption:
  - SIGINT/SIGTERM propagates cancellation to child process when possible.
- Retry:
  - No automatic retry for storage or command execution in MVP (fail fast with clear message).

### 10) UX States and Accessibility Baseline (CLI)

- Success state: concise confirmation message and exit code `0`.
- Empty state: clear `no commands found` output for `list`.
- Validation error state: actionable message + usage hint + non-zero exit.
- Dangerous confirmation state: explicit prompt (default no), example `[y/N]`.
- Accessibility baseline:
  - Plain text output first; no color-only meaning.
  - Stable output format for screen-reader and script parsing compatibility.

### 11) Validation Rules

- `create`
  - Require exactly two args: `key`, `value`.
  - Reject blank/whitespace-only values.
- `delete`
  - Require one non-empty `id` arg.
- `run`
  - Resolve a key deterministically.
  - Reject ambiguous tie if tie-break cannot choose uniquely.
  - Ensure all placeholders are bound before execution.
- `catalog`
  - Reject invalid JSON with explicit recovery guidance.

### 12) Observability and Telemetry

- MVP default: no remote telemetry.
- Local optional debug logs behind flag/env (example: `CS_DEBUG=1`).
- Log events (debug mode):
  - `catalog_loaded`
  - `command_created`
  - `command_deleted`
  - `match_resolved`
  - `danger_confirmation_prompted`
  - `command_executed`

### 13) Security and Privacy

- Treat stored commands as user-trusted but potentially dangerous.
- Use confirmation gate for dangerous command patterns.
- Do not escalate privileges automatically; execution inherits current user context.
- Sanitize logs to avoid leaking sensitive command args in non-debug mode.
- Limit shell invocation surface to required execution path.

### 14) Test Strategy

- Unit tests
  - validators (`create`, `delete`, placeholder parsing)
  - matching score and deterministic tie-break
  - danger classification
  - template binder
- Integration tests
  - JSON repository load/save with atomic replace
  - create/list/delete end-to-end against temp catalog file
  - run flow with confirmation accepted/denied
- CLI contract tests
  - exit codes and stderr/stdout for success/failure paths
  - completion output generation for target shell(s)

### 15) Incremental Delivery Plan

- Phase 1: core catalog (`create`, `list`, `delete`) + JSON persistence.
- Phase 2: command execution + placeholder binding + error handling.
- Phase 3: dangerous detection + confirmation gate.
- Phase 4: TAB completion generation + matching refinements.
- Phase 5: hardening tests, docs, and release packaging.

### 16) Technical Risks and Mitigations

- Ambiguous matching causes wrong execution -> deterministic scoring and strict tie handling.
- Shell completion inconsistency across shells -> support one shell first (bash), add zsh after contract stabilizes.
- JSON corruption on interrupted writes -> temp file + fsync + atomic rename.
- Dangerous pattern false positives/negatives -> maintain explicit pattern list with tests and user feedback loop.
- Placeholder parsing edge cases -> restrict MVP syntax and document unsupported patterns.

### 17) Open Decisions

- OD01: Final best-match scoring formula and tie-break order.
- OD02: Initial dangerous command pattern catalog.
- OD03: Confirmation prompt behavior in non-interactive mode (flag/env override).
- OD04: Key uniqueness policy (`exact unique` vs `allow overlapping patterns`).
- OD05: Initial completion shell scope (`bash only` vs `bash + zsh`).
