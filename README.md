# chores

A command-line tool for tracking household chores using a plain-text markdown file. Like plain-text accounting, but for chores.

## Installation

```bash
go install github.com/kusha/chores-md/cmd/chores@latest
```

Or build from source:

```bash
git clone https://github.com/kusha/chores-md
cd chores-md
go build ./cmd/chores
```

## Usage

```bash
chores                          # Show what's due (default)
chores show                     # Same as above
chores list                     # List all defined chores
chores done "Chore Name"        # Mark a chore as completed today
chores -f ~/my-chores.md show   # Use a custom file path
chores --help                   # Show help
chores --version                # Show version
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-f PATH` | Path to chores file | `./chores.md` |

**Note:** Global flags must come BEFORE the subcommand.

## File Format

### Chore Definition

Define chores using level-2 markdown headers with a frequency blockquote:

```markdown
## Kitchen - Clean Stovetop
> 2w

Optional description text here.
```

### Frequency Codes

| Code | Meaning |
|------|---------|
| `1d` | Every day |
| `2d` | Every 2 days |
| `1w` | Every week |
| `2w` | Every 2 weeks |
| `1m` | Every month (30 days) |
| `3m` | Every 3 months (90 days) |
| `1y` | Every year (365 days) |

### Duration Estimation (Optional)

You can optionally add estimated duration to any chore after the frequency:

```markdown
## Kitchen - Clean Stovetop
> 2w 30m
```

Supported formats:
- `30m` - Minutes only
- `2h` - Hours only
- `1h30m` - Hours and minutes combined

The CLI will display estimated durations and totals in the `show` command.

### Completion Entries

Log completions anywhere in the file using ISO date format:

```markdown
2026-02-03 Kitchen - Clean Stovetop
2026-02-01 Take Out Trash  # optional comment
```

### Example File

```markdown
# My Household Chores

## Kitchen - Clean Stovetop
> 2w

Wipe down with degreaser, clean drip pans.

## Take Out Trash
> 2d

Take all trash and recycling to curb.

## Vacuum Living Room
> 1w

Move furniture to get under couch.

---

# Completion Log

2026-02-03 Kitchen - Clean Stovetop
2026-02-02 Take Out Trash
2026-01-28 Vacuum Living Room
```

## Output Examples

### `chores show`

```
OVERDUE
  Vacuum Living Room (~1h) (3 days overdue)
    Last: 2026-01-28

DUE TODAY
  Take Out Trash
    Last: 2026-02-02

UPCOMING (7 days)
  Kitchen - Clean Stovetop (~30m) (due in 5 days)
    Last: 2026-02-03
  Total: 1h 30m

ALL CLEAR
  Garage Cleanup (due in 45 days)
    Last: 2026-01-01
```

### `chores list`

```
Garage Cleanup	every 3m	never
Kitchen - Clean Stovetop	every 2w	Last: 2026-02-03
Take Out Trash	every 2d	Last: 2026-02-02
Vacuum Living Room	every 1w	Last: 2026-01-28
```

### `chores done`

```
Done: "Take Out Trash" (2026-02-04)
```

## Design Philosophy

- **Plain-text first**: Your data lives in a readable markdown file
- **Git-friendly**: Append-only completion log for clean diffs
- **No database**: Single file contains everything
- **Unix-friendly**: Pipe-friendly output, standard exit codes
- **Minimal**: Edit chores in your text editor; CLI only shows status and marks done

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (file not found, parse error, unknown chore) |

## License

MIT
