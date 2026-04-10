# json-sheets

A terminal UI for viewing and editing flat JSON files.

> This project started as a fork of [maaslalani/sheets](https://github.com/maaslalani/sheets)
> and was refactored to exclusively support flat JSON key-value documents,
> removing CSV, markdown, and formula support in favour of a simpler, focused tool.

## Usage

```bash
sheets <file.json>
```

Opens the JSON file in an interactive two-column table:

| KEYS | VALUES |
|------|--------|
| name | Alice  |
| age  | 30     |
| ...  | ...    |

The program requires a flat JSON object as input:

```json
{
  "name": "Alice",
  "age": 30,
  "active": true,
  "city": "New York"
}
```

## Installation

```bash
make build
```

Produces a static binary `sheets` in the project root.

## Keybindings

### Navigation

| Key | Action |
|-----|--------|
| `h j k l` or arrow keys | Move cursor |
| `gg` | Go to first row |
| `G` | Go to last row |
| `Ctrl+D` | Half page down |
| `Ctrl+U` | Half page up |
| `0` | Jump to KEYS column |
| `$` | Jump to last non-blank column |

### Editing

| Key | Action |
|-----|--------|
| `i` | Edit current cell |
| `I` | Edit current cell from start |
| `c` | Clear and edit current cell |
| `o` | Insert row below |
| `O` | Insert row above |
| `dd` | Delete current row |
| `u` | Undo |
| `Ctrl+R` or `U` | Redo |

### Visual mode

| Key | Action |
|-----|--------|
| `v` | Enter visual mode |
| `V` | Enter row visual mode |
| `y` | Yank selection |
| `x` | Cut selection |
| `p` | Paste |
| `Esc` | Exit visual mode |

### Search

| Key | Action |
|-----|--------|
| `/` | Search forward |
| `?` | Search backward |
| `n` | Next match |
| `N` | Previous match |

### Commands

| Command | Action |
|---------|--------|
| `:w` | Save file |
| `:w <path>` | Save to path |
| `:e <path>` | Open another JSON file |
| `:set width=N` | Set column width |
| `:q` | Quit (prompts if unsaved) |
| `:q!` | Quit without saving |
| `:wq` or `:x` | Save and quit |

## Options

```
sheets [options] <file.json>

Options:
  -h, --help        Show help
  -v, --version     Show version
  -w, --width N     Set column cell width (default: auto)
```

## Acknowledgements

This project is based on [maaslalani/sheets](https://github.com/maaslalani/sheets)
by [Maas Lalani](https://github.com/maaslalani), originally a general-purpose
terminal spreadsheet supporting CSV, TSV, and Markdown.

The codebase was refactored to strip out everything except JSON support,
simplify the data model to flat key-value pairs, and reorient the UI
around a fixed two-column KEYS/VALUES layout.
