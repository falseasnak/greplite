# greplite

A fast grep-like CLI with built-in support for structured log formats like JSON and logfmt.

---

## Installation

```bash
go install github.com/yourusername/greplite@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/greplite.git
cd greplite
go build -o greplite .
```

---

## Usage

```bash
# Basic pattern search (like grep)
greplite "error" app.log

# Search within JSON log fields
greplite --format json --field message "connection refused" app.log

# Search logfmt structured logs
greplite --format logfmt --field level "error" app.log

# Pipe support
cat app.log | greplite --format json --field status "500"

# Case-insensitive search
greplite -i "warning" app.log
```

### Flags

| Flag | Description |
|------|-------------|
| `--format` | Log format: `json`, `logfmt`, or `text` (default: `text`) |
| `--field` | Target a specific field in structured logs |
| `-i` | Case-insensitive matching |
| `-n` | Show line numbers |
| `-c` | Print match count only |

---

## Why greplite?

Standard `grep` treats every line as plain text. `greplite` understands structured log formats, letting you search specific fields without wrestling with `jq` pipelines or `awk` one-liners.

---

## License

MIT © 2024 yourusername