# Technology Stack

<cite>
Source files referenced:
- [go.mod](file://go.mod)
- [Makefile](file://Makefile)
- [cmd/repowiki/main.go](file://cmd/repowiki/main.go)
- [internal/git/git.go](file://internal/git/git.go)
- [internal/wiki/engine.go](file://internal/wiki/engine.go)
</cite>

## Table of Contents

- [Programming Language](#programming-language)
- [Build System](#build-system)
- [Dependencies](#dependencies)
- [External Tools](#external-tools)
- [Project Structure](#project-structure)

## Programming Language

**Go 1.22+** — The entire project is written in Go, leveraging:
- Standard library for all core functionality
- No external dependencies (zero third-party imports)
- Cross-platform compatibility (macOS, Linux)
- Single binary deployment

```go
// go.mod
module github.com/IKrasnodymov/repowiki

go 1.22.0
```

## Build System

The project uses a simple **Makefile** for build automation:

```makefile
.PHONY: build install test clean

BINARY := repowiki
BUILD_DIR := bin

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/repowiki

install:
	go install ./cmd/repowiki

test:
	go test ./internal/... -v -race

clean:
	rm -rf $(BUILD_DIR)
```

### Build Commands

| Command | Description |
|---------|-------------|
| `make build` | Compile binary to `bin/repowiki` |
| `make install` | Install to `$GOPATH/bin` |
| `make test` | Run tests with race detection |
| `make clean` | Remove build artifacts |

## Dependencies

This project has **zero external dependencies**. All functionality is implemented using Go's standard library:

### Standard Library Packages Used

| Package | Purpose | Files Using |
|---------|---------|-------------|
| `os` | File operations, process management | All files |
| `os/exec` | External command execution | `git/git.go`, `wiki/qoder.go` |
| `path/filepath` | Cross-platform path manipulation | All files |
| `encoding/json` | Configuration serialization | `config/config.go` |
| `fmt` | Formatted I/O | All files |
| `strings` | String manipulation | Multiple files |
| `time` | Timestamp handling | `config/config.go`, `wiki/wiki.go` |
| `syscall` | Process detachment | `hooks.go` |
| `strconv` | String/number conversion | `wiki/commit.go`, `wiki/qoder.go` |
| `bytes` | Buffer management | `wiki/qoder.go` |

### Why Zero Dependencies?

The project intentionally avoids external dependencies for:
- **Simplicity**: No dependency management complexity
- **Reliability**: No risk of broken or compromised dependencies
- **Portability**: Single binary that just works
- **Security**: Minimal attack surface

## External Tools

While the binary itself has no dependencies, it integrates with external tools:

### Required

| Tool | Purpose | Detection |
|------|---------|-----------|
| `git` | Repository operations | Must be in PATH |
| `qodercli`, `claude`, or `codex` | AI-powered wiki generation | Configurable path |

### Engine Detection

The tool searches for the configured AI engine binary in multiple locations:

```go
func FindEngineBinary(cfg *config.Config) (string, error) {
    switch cfg.Engine {
    case config.EngineQoder:
        return findQoderBinary(cfg)
    case config.EngineClaudeCode:
        return findClaudeCodeBinary(cfg)
    case config.EngineCodex:
        return findCodexBinary(cfg)
    default:
        return "", fmt.Errorf("unknown engine: %s", cfg.Engine)
    }
}
```

For Qoder CLI, it checks:
1. Config `engine_path` override
2. PATH for `qodercli`
3. Known macOS locations (Qoder.app bundle)

For Claude Code, it checks:
1. Config `engine_path` override
2. PATH for `claude`
3. Common locations (`~/.local/bin/claude`, `~/.claude/bin/claude`, `/usr/local/bin/claude`)

For OpenAI Codex, it checks:
1. Config `engine_path` override
2. PATH for `codex`

## Project Structure

```
repowiki/
├── cmd/repowiki/          # CLI entry points
│   ├── main.go            # Command router
│   ├── enable.go          # Enable command
│   ├── disable.go         # Disable command
│   ├── status.go          # Status command
│   ├── generate.go        # Full generation
│   ├── update.go          # Incremental update
│   ├── hooks.go           # Hook entry point
│   └── logs.go            # Logs command
├── internal/              # Internal packages
│   ├── config/            # Configuration management
│   ├── git/               # Git operations
│   ├── hook/              # Git hook management
│   ├── lockfile/          # Process locking
│   └── wiki/              # Wiki generation
│       ├── wiki.go        # Core generation logic
│       ├── engine.go      # Multi-engine abstraction
│       ├── commit.go      # Auto-commit logic
│       ├── detect.go      # Change detection
│       └── prompt.go      # Prompt building
├── go.mod                 # Go module definition
├── Makefile               # Build automation
└── README.md              # User documentation
```

## Configuration Files

### Runtime Configuration

| File | Format | Purpose |
|------|--------|---------|
| `.repowiki/config.json` | JSON | Tool configuration |
| `.repowiki/.repowiki.lock` | Text | Process lock file |
| `.repowiki/.committing` | Text | Sentinel file (loop prevention) |
| `.repowiki/logs/*.log` | Text | Execution logs |

### Generated Output

| File/Directory | Format | Purpose |
|----------------|--------|---------|
| `.qoder/repowiki/en/content/*.md` | Markdown | Wiki documentation |
| `.qoder/repowiki/en/meta/repowiki-metadata.json` | JSON | Code snippet index |
