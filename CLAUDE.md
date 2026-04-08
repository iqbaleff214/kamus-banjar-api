# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Project Is

A REST API for Kamus Banjar (Banjar language dictionary), written in Go using the Fiber web framework. Data is stored in JSON files per alphabet letter (`data/a.json`, `data/b.json`, etc.).

## Commands

```bash
# Run
go run .

# Build (default/embedded mode)
go build -ldflags "-s -w" -o ./bin/app .

# Build (filesystem mode, for Docker)
go build -ldflags "-s -w" -tags=fs -o server .

# Test
go test -v -cover ./...

# Run single test file
go test -v ./domain/dictionary/...
```

## Architecture

Clean architecture with a single domain (`domain/dictionary/`):

- **`main.go`** — app entry, Fiber setup, route registration
- **`handler.go`** — root `/api/v1/` handler
- **`domain/dictionary/handler.go`** — HTTP handlers for dictionary endpoints
- **`domain/dictionary/service.go`** — business logic (fuzzy search via Levenshtein distance)
- **`domain/dictionary/repository.go`** — repository interface + factory that selects implementation
- **`domain/dictionary/dictionary.go`** — in-memory index over JSON data with search
- **`domain/dictionary/model.go`** — data models (`Word`, `Letter`, etc.)
- **`domain/dictionary/utils/`** — string utilities

### Repository Implementations

The factory in `repository.go` picks one of four implementations:

| Implementation | File | When Used |
|---|---|---|
| Embedded FS | `repository_embed.go` | Default (no build tag, no `MYSQL_DSN`) |
| Filesystem | `repository_fs.go` | `fs` build tag, no `MYSQL_DSN` |
| MySQL | `repository_mysql.go` | `MYSQL_DSN` env var is set |
| Remote JSON (GitHub) | `repository_json.go` | Fallback when no data available |

### Build Tags

- **No tag** (default): JSON data embedded at compile time via `go:embed`
- **`fs` tag**: JSON data read from `data/` directory at runtime (Docker uses this)

## API Endpoints

```
GET /api/v1/                        # API info + memory stats
GET /api/v1/alphabets               # All letters with word counts
GET /api/v1/alphabets/:letter       # Words for a letter
GET /api/v1/entries/:word           # Word definition
GET /api/v1/entries?search=keyword  # Fuzzy search
```

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `:8001` | Server port |
| `MYSQL_DSN` | — | MySQL DSN; enables MySQL repository when set |
| `SOURCE_PATH` | `data/` | Path to JSON files (only with `fs` tag) |

## Tests

Unit tests in `domain/dictionary/service_test.go` load actual JSON data from `data/` and use a mock repository. Test cases are generated dynamically from the real data files — avoid modifying test data structures without understanding how test cases are auto-generated.
