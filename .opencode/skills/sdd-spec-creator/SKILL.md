---
name: sdd-spec-creator
description: Create SDD specs for a feature using 4 files (01-context, 02-requirements, 03-design, 04-tasks) with testable acceptance criteria, architecture decisions, and atomic execution tasks.
license: MIT
metadata:
  version: '1.0.0'
  category: documentation
  type: simple-task
  triggers:
    - 'create sdd'
    - 'create spec docs'
    - 'feature design doc'
    - 'specify feature'
    - 'requirements design tasks'
---

# SDD Spec Creator

Create feature specs in a consistent SDD flow with 4 files:

1. `01-context.md`
2. `02-requirements.md`
3. `03-design.md`
4. `04-tasks.md`

Target output path:

`docs/feature-sdd/<feature-slug>/`

If `00-context.md` exists, use it as source of truth and evolve it into the 4 files above.

## Primary Goal

Transform product intent into implementation-ready documentation:

- clear scope and non-goals
- testable requirements
- technical architecture and contracts
- atomic tasks with dependencies and done criteria

## Mandatory Rules

1. Do not invent business rules when they are truly unknown; mark open decisions explicitly.
2. Keep acceptance criteria testable (`Dado / Quando / Entao` or equivalent).
3. Keep design implementation-oriented (state, services, routes, data contracts, error states).
4. Keep tasks atomic and executable in small PRs.
5. Preserve existing repository conventions (naming, folder layout, stack).
6. Keep docs in ASCII unless file already uses Unicode intentionally.

## Input Signals

When user provides any of these, this skill should run:

- initial context doc (ex.: `00-context.md`)
- rough journey in bullets
- list of business constraints
- request like "evolve my design doc" or "generate specs"

## Required Workflow

### Phase 1 - Context (`01-context.md`)

Build or refine:

- overview and MVP objective
- in-scope and out-of-scope
- assumptions/constraints
- user journey/steps
- business rules
- initial data model contracts
- functional/non-functional requirements summary
- risks and open decisions

Reference template: `references/01-context-template.md`

### Phase 2 - Requirements (`02-requirements.md`)

Create:

- user stories (MVP)
- acceptance criteria per story (testable)
- business rules catalog
- non-functional requirements
- exception flows
- definition of done
- pending decisions

Reference template: `references/02-requirements-template.md`

### Phase 3 - Design (`03-design.md`)

Define:

- architecture layers and responsibilities
- routes and flow responsibilities
- state model (global/local)
- service boundaries and interfaces
- domain types/contracts
- async strategy (polling/retry/status transitions if needed)
- validation, UX states, accessibility baseline
- test strategy and incremental delivery plan

Reference template: `references/03-design-template.md`

### Phase 4 - Tasks (`04-tasks.md`)

Break into implementation tasks with:

- id (`T001`, `T002`, ...)
- priority (`P0`, `P1`, `P2`)
- dependencies
- scope
- done criteria
- recommended execution order
- release milestones

Reference template: `references/04-tasks-template.md`

## Quality Checklist Before Finalizing

- Every requirement maps to at least one task.
- Every key journey step has acceptance criteria.
- Design covers data flow, state transitions, and error handling.
- Open decisions are visible and actionable.
- Wording is concise and implementation-focused.

## Output Contract

Create or update exactly these files:

- `docs/feature-sdd/<feature-slug>/01-context.md`
- `docs/feature-sdd/<feature-slug>/02-requirements.md`
- `docs/feature-sdd/<feature-slug>/03-design.md`
- `docs/feature-sdd/<feature-slug>/04-tasks.md`

Optional:

- if project has reusable templates, also update `docs/feature-sdd/_template/*` when asked.

## Suggested Command Patterns

- "Create SDD for `<feature-slug>` based on `00-context.md`."
- "Evolve my feature doc to full SDD (context, requirements, design, tasks)."
- "Generate requirements and tasks from this journey description."
