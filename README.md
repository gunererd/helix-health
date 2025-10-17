# helix-health

> Overengineered helix --health

Interactive TUI for viewing and searching Helix editor's health information.

![demo demo](demo.gif)

## Features

- **Interactive search** - Real-time fuzzy filtering as you type
- **Non-interactive mode** - Quick command-line queries
- **Multiple search terms** - Filter by multiple languages or tools at once

## Installation

### From source

```bash
git clone https://github.com/gunererd/helix-health.git
cd helix-health
go build
```

### Using go install

```bash
go install github.com/gunererd/helix-health@latest
```

## Usage

### Interactive mode

Launch the TUI:

```bash
helix-health
```

Type to search for languages or tools. Navigate with arrow keys or Vim/Emacs shortcuts.

### Non-interactive mode

Query specific languages from the command line:

```bash
helix-health python
helix-health rust go typescript
```

## Requirements

- [Helix editor](https://helix-editor.com/) installed and available in PATH
- Go 1.21+ (for building from source)

## How it works

helix-health wraps the `helix --health` command and provides:
- Real-time substring search across language names and tool names
- Highlighted matches in search results
- Clean, styled terminal output

## License

MIT
