# DayPlanner — Design Document

## Overview

A terminal-based daily task planner for a single person. Built around dependency-first thinking: the tool surfaces *what you can actually work on right now*, not just a flat list of todos. Lightweight Jira-like structure without the overhead.

---

## Tech Stack

| Concern | Choice |
|---|---|
| Language | Go |
| TUI framework | Bubble Tea |
| Styling | Lipgloss |
| LLM integration | LangchainGo |
| Persistence | JSON flat file |
| Config/path resolution | XDG base dir spec |

---

## File Structure

```
config/
  config.go           - XDG path resolution, Dir() method → ~/.config/dayplanner/

internal/
  domain/             - Pure data types only. No logic touching other packages.
  graph/              - DAG construction, cycle detection, topological sort
  scheduler/          - Ordering logic: graph + priority + creation time tiebreaker
  repository/         - JSON read/write → ~/.config/dayplanner/state.json
  history/            - Undo/redo stack, max 25 entries, covers delete and edit
  inference/          - LangchainGo integration for LLM-assisted task drafting
  model/              - Bubble Tea root model, page state, keybindings (single source of truth)
  update/             - Bubble Tea Update handlers, one per page
  ui/                 - Stateless rendering: Lipgloss styles, headers, footers, components

app.go                - Bubble Tea wiring, avoids import cycles
main.go               - Entry point
```

---

## Domain Model

```go
type Task struct {
    ID        string    // Semantic: "AUTH-login-flow" (TAG-slug)
    Tag       string    // Area prefix e.g. "AUTH", "API", "INFRA" (2–6 chars, uppercase)
    Name      string
    Goal      string    // Acceptance criteria. LLM-assisted draft.
    Priority  Priority  // high | medium | low
    Status    Status    // todo | in-progress | blocked | done
    DependsOn []string  // IDs of tasks that must be done before this one can start
    Notes     string    // Single freetext field (link, context, blocker reason)
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Priority string  // "high" | "medium" | "low"
type Status   string  // "todo" | "in-progress" | "blocked" | "done"
```

**Key design rules:**
- `DependsOn` is the only stored dependency direction. "Blocks" is derived at read time (reverse lookup).
- `blocked` status is set automatically when any dependency is not `done`.
- IDs are user-authored or LLM-suggested. Format enforced: `[TAG]-[slug]`. Slug is lowercase, hyphen-separated.
- Tags constrained to 2–6 uppercase characters.

### JSON Representation

Tasks serialize directly to/from JSON. The state file is a single array:

```json
[
  {
    "id": "AUTH-login-flow",
    "tag": "AUTH",
    "name": "Implement login flow",
    "goal": "User can log in with email/password. Invalid credentials show an error. Session persists on refresh.",
    "priority": "high",
    "status": "in-progress",
    "depends_on": ["AUTH-setup"],
    "notes": "See RFC doc at https://...",
    "created_at": "2025-01-15T08:00:00Z",
    "updated_at": "2025-01-15T09:30:00Z"
  }
]
```

---

## Scheduler

Ordering rules applied in sequence:

1. **Dependency layer first** — tasks at layer 0 (no dependencies) before layer 1, etc.
2. **Priority within layer** — high → medium → low
3. **Creation time as tiebreaker** — earlier created tasks surface first

The scheduler output feeds directly into the Dashboard's "Ready now" section.

---

## Tag Color System

Each unique tag gets a persistent color assigned at first creation. Colors are:
- Generated deterministically from the tag string (e.g. hash-based) so they're stable without storing extra state
- Drawn from a curated palette of Lipgloss-compatible ANSI colors that are visually distinct and readable on dark and light terminals
- Used consistently across all views: dashboard rows, task view header, graph nodes

---

## Pages

### 1. Dashboard (startup view)

Three sections, always visible on launch:

```
┌─ READY NOW ──────────────────────────────────────────────┐
│ ● AUTH-login-flow    [HIGH]  Implement login flow         │
│ ● API-endpoints      [MED]   Define REST endpoints        │
├─ BLOCKED ────────────────────────────────────────────────┤
│ ○ AUTH-deploy        blocked by: AUTH-login-flow          │
│ ○ API-tests          blocked by: API-endpoints            │
├─ DONE ───────────────────────────────────────────────────┤
│ ✓ AUTH-setup         Initial project setup               │
└──────────────────────────────────────────────────────────┘
```

- **Done** section is collapsed by default, expandable with a keypress
- Carried-over tasks from previous days are flagged visually
- Tasks carried over 3+ times prompt: *"Still relevant? [k]eep / [d]rop"*
- `space` advances task status inline (no need to open the task)
- `g` toggles between scheduler order and grouped-by-tag view
- `enter` opens Task View for the selected task
- `n` opens Add Task page
- `d` opens Dependency Graph page

### 2. Add / Edit Task Page

Two modes selectable at the top: **Single Task** and **Chain**.

**Single Task mode:**

Sequential form fields navigated with `tab` / `shift+tab`:
- Tag (autocomplete from existing tags)
- Name
- Priority (select: high / medium / low)
- Goal (LLM drafts this when you tab away from Name — non-blocking, appears as suggestion to accept/edit/ignore)
- Dependency (optional — type existing task IDs, autocomplete supported)
- Notes (optional)

**Chain mode:**

Each entry is a mini form (Tag, Name, Priority). Pressing `enter` saves the current task and opens the next, with `DependsOn` automatically set to the previous task's ID. The chain is shown inline as you build it:

```
AUTH-write-tests → AUTH-review → AUTH-merge → [next...]
```

Press `esc` to close the chain and save all created tasks.

**LLM assistance (Goal field):**
- Triggered automatically when focus leaves the Name field
- Small model call: `Given task name "[name]", draft 2–3 acceptance criteria as a single concise paragraph`
- Result appears in Goal field as editable draft
- User can accept (keep editing), regenerate (`ctrl+r`), or clear and type manually

### 3. Dependency Graph Page

**Layout:** Left-to-right layered columns. Each column = one dependency depth level.

```
[ Layer 0 ]        [ Layer 1 ]        [ Layer 2 ]
┌──────────┐       ┌──────────┐       ┌──────────┐
│AUTH-setup│──────►│AUTH-login│──────►│AUTH-deplo│
└──────────┘       └──────────┘       └──────────┘
┌──────────┐       ┌──────────┐
│API-schema│──────►│API-endpts│
└──────────┘       └──────────┘
                   ┌──────────┐
                   │API-tests │
                   └──────────┘
```

- Critical path highlighted in a distinct accent color
- Task nodes colored by tag
- Arrow characters use Unicode box-drawing (`─`, `►`, `│`)
- Arrow routing: straight horizontal for same-layer deps, vertical connectors for multi-target nodes
- Navigation: arrow keys move between nodes, `enter` opens Task View for selected node, `esc` returns to Dashboard
- Scrollable viewport if graph exceeds terminal size

### 4. Task View Page

Two-section layout:

**Top — Task details:**
```
AUTH-login-flow                              [HIGH] [IN PROGRESS]
──────────────────────────────────────────────────────────────────
Name:    Implement login flow
Goal:    User can log in with email/password. Invalid credentials
         show an error. Session persists on refresh.
Notes:   See RFC doc at https://...
Created: 2025-01-15 08:00   Updated: 2025-01-15 09:30
```

**Bottom — Relationships:**
```
Dependencies
  Blocked by:  AUTH-setup [DONE]

Dependents
  Blocks:      AUTH-deploy [BLOCKED]
               AUTH-e2e-tests [TODO]
```

- Related task IDs are navigable: `enter` on a listed task to jump to its Task View
- `e` to open Edit mode for this task (reuses Add/Edit page in edit mode)
- `s` to cycle status
- `esc` to return to previous page

---

## History (Undo/Redo)

- Covers: task edits and task deletions
- Horizon: last 25 operations
- `ctrl+z` to undo, `ctrl+y` to redo
- Stored in memory only (not persisted to disk — resets on app restart)
- Each entry records the full previous task state (simple snapshot, not diff)

---

## Keybinding Summary

| Key | Action |
|---|---|
| `n` | New task |
| `enter` | Open selected task |
| `space` | Advance status (todo → in-progress → done) |
| `d` | Open dependency graph |
| `g` | Toggle grouping (scheduler / by tag) |
| `e` | Edit task (from Task View) |
| `s` | Cycle status (from Task View) |
| `ctrl+z` | Undo |
| `ctrl+y` | Redo |
| `ctrl+r` | Regenerate LLM suggestion (in Add/Edit) |
| `esc` | Back / cancel |
| `q` | Quit |

---

## LLM Integration Points

All LLM calls are:
- Explicitly triggered (never background/automatic except Goal draft on Name tab-out)
- Non-blocking — UI remains responsive while waiting
- Handled via LangchainGo with structured output parsing (JSON)
- Targeting a small, fast model (no deep reasoning required)

| Trigger | Input | Output |
|---|---|---|
| Tab out of Name (Add Task) | Task name | Draft acceptance criteria for Goal field |
| `ctrl+space` on Dashboard | Full task list | Suggested re-ordering or flag if overcommitted |
| Chain mode completion | List of task names | Suggested IDs and dependency linkage |

---

## Persistence

- State file: `~/.config/dayplanner/state.json`
- Written on every mutation (add, edit, delete, status change)
- Read on startup
- No migration strategy needed at v1 — schema is append-friendly (new optional fields safe to add)
- Old day's tasks carry over automatically; no archiving at v1

---

## Out of Scope (v1)

- Subtasks (by design — use tags for grouping instead)
- Scheduled times (sequence only, no calendar integration)
- Cross-day analytics / historical queries
- Multi-user / sync
- Notifications or reminders
