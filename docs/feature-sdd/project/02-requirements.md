## Requirements - command-catalog-cli

### 1) Objective

Define testable MVP requirements for a Golang CLI (`cs`) that manages a local command catalog in JSON, supports parameterized command execution, provides TAB-assisted key completion, and enforces confirmation for dangerous commands.

### 2) MVP Scope

- Create command templates: `cs create "<key>" "<value>"`.
- List catalog commands: `cs list` with `id`, `key`, `value`.
- Delete catalog command by id: `cs delete <id>`.
- Execute commands through key-based invocation with runtime args.
- Resolve placeholders from template values at runtime.
- Persist catalog in a local JSON file.
- Provide shell completion for TAB to suggest/autofill best key match.
- Require explicit user confirmation for dangerous commands.

### 3) Personas and Needs

- Developer/Operator: needs to avoid retyping long shell commands.
- Linux power user: needs fast and predictable command reuse with arguments.
- Security-conscious user: needs a confirmation gate before risky command execution.

### 4) User Stories (MVP)

#### US01 - Create catalog command

As a Linux user, I want to create a reusable command template so I can execute repetitive shell commands faster.

**Acceptance Criteria**

- Dado que informo `key` e `value` validos, quando executo `cs create "<key>" "<value>"`, entao o sistema salva um novo comando com `id` unico no catalogo JSON.
- Dado que `key` ou `value` esta vazio, quando executo `cs create`, entao o sistema rejeita a operacao com erro de validacao e codigo de saida nao zero.
- Dado que a persistencia falha, quando executo `cs create`, entao o sistema nao grava estado parcial e retorna erro acionavel.

#### US02 - List catalog commands

As a user, I want to list registered commands so I can inspect what is available.

**Acceptance Criteria**

- Dado que existem comandos no catalogo, quando executo `cs list`, entao o sistema exibe cada entrada com `id`, `key` e `value`.
- Dado que nao existem comandos, quando executo `cs list`, entao o sistema exibe estado vazio sem erro de execucao.

#### US03 - Delete command by id

As a user, I want to delete a command by id so I can remove obsolete entries.

**Acceptance Criteria**

- Dado um `id` existente, quando executo `cs delete <id>`, entao o sistema remove apenas a entrada correspondente e persiste o catalogo atualizado.
- Dado um `id` inexistente, quando executo `cs delete <id>`, entao o sistema retorna erro de nao encontrado e codigo de saida nao zero.

#### US04 - Execute command with parameters

As a user, I want to execute a stored command with runtime arguments so placeholders are replaced and the underlying shell command runs.

**Acceptance Criteria**

- Dado um comando com placeholder (ex.: `$port`), quando executo `cs kill port 3040`, entao o sistema resolve `$port=3040` e monta o comando final para execucao.
- Dado que faltam argumentos obrigatorios para placeholders, quando executo o comando, entao o sistema bloqueia a execucao e retorna erro de validacao.
- Dado que a chave informada nao corresponde a nenhum comando, quando executo `cs ...`, entao o sistema retorna erro de comando nao encontrado com codigo nao zero.

#### US05 - TAB completion for command keys

As a user, I want TAB completion to suggest/autofill the best key match so command discovery and execution are faster.

**Acceptance Criteria**

- Dado uma chave parcial e catalogo com correspondencias, quando pressiono TAB, entao o shell sugere/autopreenche a melhor chave de forma deterministica.
- Dado multiplas correspondencias com mesma pontuacao, quando pressiono TAB, entao o sistema aplica desempate deterministico definido na estrategia de matching.

#### US06 - Confirmation for dangerous commands

As a safety-conscious user, I want confirmation before dangerous command execution so accidental destructive actions are reduced.

**Acceptance Criteria**

- Dado que um comando foi classificado como perigoso, quando tento executa-lo, entao o sistema solicita confirmacao explicita antes de rodar.
- Dado que respondo negativamente a confirmacao, quando o prompt e finalizado, entao o comando nao e executado.
- Dado que respondo positivamente a confirmacao, quando o prompt e finalizado, entao o comando e executado normalmente.

### 5) Business Rules Catalog

- BR01: Cada comando possui `id` unico.
- BR02: `key` e `value` sao obrigatorios e nao podem ser vazios.
- BR03: `cs list` sempre expoe `id`, `key`, `value`.
- BR04: `cs delete <id>` remove apenas o comando do `id` informado.
- BR05: Placeholders devem ser resolvidos com argumentos de runtime antes da execucao.
- BR06: Comando nao resolvido deve retornar erro claro de nao encontrado.
- BR07: Ausencia de argumentos obrigatorios bloqueia execucao.
- BR08: Operacoes de create/delete devem persistir imediatamente em JSON.
- BR09: Comandos perigosos exigem confirmacao explicita.
- BR10: Atualizacao in-place nao faz parte do MVP (fluxo e delete + create).

### 6) Non-Functional Requirements

- NFR01: Saidas do CLI devem ser legiveis e deterministicas.
- NFR02: Operacoes de catalogo devem ser eficientes para volumes pequenos/medios.
- NFR03: Erros devem retornar codigo de saida nao zero e mensagem acionavel.
- NFR04: Persistencia deve usar estrategia de escrita atomica (temp + replace).
- NFR05: Solucao deve operar em Linux no MVP.
- NFR06: Matching e TAB completion devem ser deterministas para mesma entrada e mesmo estado do catalogo.

### 7) Data Requirements

- Main entities:
  - `CatalogFile { version: string, commands: CatalogCommand[] }`
  - `CatalogCommand { id, key, value, createdAt?, updatedAt? }`
  - `ExecutionPlan { commandId, resolvedCommand, isDangerous }`
  - `ExecutionResult { exitCode, stdout?, stderr? }`
- Storage:
  - JSON file on local filesystem.
  - File must remain valid JSON after each successful mutation.
- Data integrity:
  - No duplicated `id`.
  - `key` uniqueness policy is pending decision (see section 11).

### 8) Exception Flows

- EF01: Invalid CLI args for `create` -> return usage + validation error.
- EF02: Missing catalog file -> initialize empty catalog on first write or return empty state on list.
- EF03: Corrupted JSON catalog -> abort operation and return recovery-oriented error.
- EF04: Placeholder mismatch (missing args) -> block execution.
- EF05: Dangerous command denied by user -> abort without side effects.
- EF06: Ambiguous or unresolved key match at runtime -> return deterministic disambiguation/not-found error.
- EF07: Non-interactive execution for dangerous command without explicit override -> abort and return actionable error.

### 9) Technical Dependencies

- Golang toolchain (version to be defined in implementation phase).
- Standard library for file I/O, process execution, JSON encoding/decoding.
- Shell completion integration scripts (bash/zsh target to be finalized).
- Test tooling: Go `testing` package (+ optional assertion helper library).

### 10) Definition of Done (MVP)

- All MVP stories (US01-US06) are implemented and validated.
- Acceptance criteria are covered by automated tests (unit/integration as applicable).
- `create`, `list`, `delete`, and execute commands run successfully in Linux.
- JSON persistence is stable and uses atomic write strategy.
- Dangerous command confirmation works in interactive mode and fails safely in non-interactive mode.
- TAB completion works with deterministic best-match behavior in at least one supported shell.
- Documentation for installation and completion setup is available.

### 11) Open Decisions

- OD01: Best-match scoring formula and deterministic tie-break order.
- OD02: Dangerous-command detection rule set (initial patterns and extensibility).
- OD03: Confirmation UX contract details (`[y/N]`, retry behavior, timeout behavior), including non-interactive mode behavior.
- OD04: `key` uniqueness rule (strict unique vs allow similar patterns).
- OD05: Initial completion shell scope (`bash only` vs `bash + zsh`).
