# journl

A minimal, fast CLI for capturing thoughts, notes, and incidents — before they slip away.

Pipe command output into it. Open your editor for longer entries. Switch contexts for work, personal, or incidents. Everything lands in a local file or database you own.

---

## Installation

```bash
go install github.com/lucasmaehn/journl@latest
journl init
```

`init` creates a default config at `~/.journl/config.yaml` and a default SQLite store at `~/.journl/db.sqlite`.

---

## Quick Start

```bash
# log a quick thought
journl log "redis pod OOMing again"

# open $EDITOR for a longer entry
journl log

# pipe command output directly
kubectl logs my-pod | journl log
```

---

## Commands

### `journl log [text]`

Create a new journal entry. If `text` is omitted, `$EDITOR` is opened. When the editor is closed, the file content is used as the entry — as long as it's not empty.

```bash
journl log "quick note"
journl log                          # opens $EDITOR
kubectl logs pod | journl log
```

**Title detection** — if the first line is followed by a completely empty second line, it is treated as the entry title:

```
This becomes the title

Everything below the empty line is the body.
```

### `journl context list`

List all configured contexts. The active context is marked with `*`.

```
* personal
  default
  work
```

### `journl context use [name]`

Switch the active context.

```bash
journl context use work
```

### `journl context add [name]`

Add a new context.

```bash
journl context add incident-foobar \
  --description "Redis OOM incident 2026-03-28" \
  --format sqlite
```

| Flag | Description |
|---|---|
| `--description` | Optional description for the context |
| `--format` | Store format: `sqlite`, `jsonl`, or `custom` |

### `journl init`

Initialise journl with a default config and context. Safe to run once — will not overwrite an existing config.

---

## Configuration

Config lives at `~/.journl/config.yaml`. Currently, advanced settings like store paths and entry templates must be edited by hand.

```yaml
current_context: personal

contexts:
  default:
    name: default
    description: The default context for journl
    store:
      format: sqlite
      path: ~/.journl/db.sqlite

  personal:
    name: personal
    description: Personal journal entries
    store:
      format: custom
      path: ~/.journl/2006-01-02.md   # time tokens are expanded automatically
      custom:
        template: |
          # {{if .Title}}{{.Title}}{{else}}{{.Timestamp.Format "2006-01-02 15:04"}}{{end}}
          > {{.Timestamp.Format "Monday, January 2 2006 at 15:04"}}{{if .Context}} · #{{.Context}}{{end}}
          {{.Body}}
          ---
```

### Store formats

| Format | Description |
|---|---|
| `sqlite` | Single SQLite database file. Good default for most contexts. |
| `jsonl` | Append-only JSONL file. One JSON object per line. |
| `custom` | Append to a file using a custom `text/template` entry format. Good, for example, to integrate into Obsidian|

### Path templates

The `path` field supports Go time format tokens, expanded at commit time:

```yaml
path: ~/.journl/2006-01-02.md      # → ~/.journl/2026-03-28.md
path: ~/.journl/work.db            # static path, no expansion
```

### Entry templates (`custom` format)

Templates use Go's `text/template` syntax. Available fields:

| Field | Type | Description |
|---|---|---|
| `.Title` | `string` | First line if followed by an empty line, otherwise empty |
| `.Body` | `string` | Entry body |
| `.Context` | `string` | Active context name |
| `.Timestamp` | `time.Time` | Entry timestamp |

---

## Future Work

- **`--file` flag** — attach files to log entries, copied into `~/.journl/attachments/`
- **`--clip` flag** — attach clipboard contents, either raw text or screenshot images
- **`--context` flag on `log`** — target a specific context inline without switching globally
- **`journl undo`** — delete the last entry in the active context
- **`journl amend`** — open the last entry in `$EDITOR` to correct or extend it

