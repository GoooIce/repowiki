# CLI Commands

<cite>
Source files referenced:
- [cmd/repowiki/main.go](/to/cmd/repowiki/main.go)
- [cmd/repowiki/enable.go](/to/cmd/repowiki/enable.go)
- [cmd/repowiki/disable.go](/to/cmd/repowiki/disable.go)
- [cmd/repowiki/status.go](/to/cmd/repowiki/status.go)
- [cmd/repowiki/generate.go](/to/cmd/repowiki/generate.go)
- [cmd/repowiki/update.go](/to/cmd/repowiki/update.go)
- [cmd/repowiki/logs.go](/to/cmd/repowiki/logs.go)
</cite>

## Table of Contents

- [Command Reference](#command-reference)
- [Global Options](#global-options)
- [enable](#enable)
- [disable](#disable)
- [status](#status)
- [generate](#generate)
- [update](#update)
- [version](#version)
- [help](#help)
- [logs](#logs)

## Command Reference

```
repowiki v0.1.0 — Auto-generate Qoder repo wiki on git commits

Usage:
  repowiki <command> [flags]

Commands:
  enable      Enable repowiki in current project (install git hook)
  disable     Disable repowiki (remove git hook)
  status      Show current status and configuration
  generate    Run full wiki generation
  update      Run incremental wiki update for recent changes
  logs        Show the most recent log file
  version     Show version

Flags for 'enable':
  --force             Reinstall hook even if already present
  --qodercli-path     Path to qodercli binary
  --model             Qoder model level (auto, efficient, performance, ultimate)
  --no-auto-commit    Don't auto-commit wiki changes

Flags for 'update':
  --commit            Specific commit hash to process
  --from-hook         Internal: indicates hook-triggered run
```

## Global Options

| Option | Description |
|--------|-------------|
| `--version`, `-v` | Show version information |
| `--help`, `-h` | Show help message |

## enable

Enable repowiki in the current git repository.

### Synopsis

```
repowiki enable [flags]
```

### Description

Installs the post-commit git hook and creates the configuration file. This is the first command to run when setting up repowiki in a project.

### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--force` | bool | false | Reinstall hook even if already present |
| `--qodercli-path` | string | "" | Path to qodercli binary |
| `--model` | string | "" | Qoder model level (auto, efficient, performance, ultimate) |
| `--no-auto-commit` | bool | false | Don't auto-commit wiki changes |

### Examples

```bash
# Enable with defaults
repowiki enable

# Force reinstall
repowiki enable --force

# Specify Qoder CLI path
repowiki enable --qodercli-path /Applications/Qoder.app/.../qodercli

# Use performance model
repowiki enable --model performance

# Disable auto-commit
repowiki enable --no-auto-commit
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (not a git repository, hook already installed, etc.) |

### Output Example

```
repowiki enabled in /path/to/project

  Config:  .repowiki/config.json
  Hook:    .git/hooks/post-commit
  Qoder:   found

Every commit will now auto-update the repo wiki.
Run 'repowiki generate' for initial full wiki generation.
```

## disable

Disable repowiki in the current repository.

### Synopsis

```
repowiki disable
```

### Description

Removes the post-commit git hook and updates the configuration to disable repowiki. Wiki files in `.qoder/repowiki/` are preserved.

### Examples

```bash
repowiki disable
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (not a git repository) |

### Output Example

```
repowiki disabled in /path/to/project
Wiki files in .qoder/repowiki/ are preserved.
```

## status

Show current status and configuration.

### Synopsis

```
repowiki status
```

### Description

Displays the current repowiki status including:
- Whether repowiki is enabled/disabled
- Hook installation status
- Qoder CLI location
- Wiki page count
- Configuration settings

### Examples

```bash
repowiki status
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (not a git repository) |

### Output Example

```
repowiki v0.1.0

  Status:       enabled
  Hook:         installed (.git/hooks/post-commit)
  Qoder CLI:    /Applications/Qoder.app/.../qodercli
  Wiki path:    .qoder/repowiki/en/content/ (12 pages)
  Model:        auto
  Auto-commit:  true
  Max turns:    50
  Last run:     2026-02-19T15:30:00Z
  Last commit:  abc123def456
```

## generate

Run full wiki generation from scratch.

### Synopsis

```
repowiki generate
```

### Description

Analyzes the entire codebase and generates comprehensive documentation. This command is useful for:
- Initial wiki creation after `repowiki enable`
- Complete documentation refresh
- Recovery from corrupted or missing wiki files

**Note**: This may take several minutes depending on codebase size.

### Examples

```bash
repowiki generate
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (not configured, generation failed) |

### Output Example

```
Starting full wiki generation... (this may take several minutes)
Wiki generation complete.
```

## update

Run incremental wiki update for recent changes.

### Synopsis

```
repowiki update [flags]
```

### Description

Updates the wiki based on files changed since the last processed commit. Uses incremental updates for small changes or full regeneration if the change threshold is exceeded.

This command is typically called automatically by the post-commit hook.

### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--commit` | string | "" | Specific commit hash to process |
| `--from-hook` | bool | false | Internal: indicates hook-triggered run |

### Examples

```bash
# Update for changes since last run
repowiki update

# Update for specific commit
repowiki update --commit abc123
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success (or no relevant changes) |
| 1 | Error (not configured, update failed) |

### Output Examples

**With changes:**
```
Updating wiki for 5 changed files...
Wiki update complete.
```

**Full regeneration triggered:**
```
Running full wiki generation (25 files changed)...
Wiki update complete.
```

**No changes:**
```
No relevant file changes detected.
```

## version

Show version information.

### Synopsis

```
repowiki version
repowiki --version
repowiki -v
```

### Examples

```bash
repowiki version
repowiki --version
repowiki -v
```

### Output Example

```
repowiki v0.1.0
```

## help

Show help message.

### Synopsis

```
repowiki help
repowiki --help
repowiki -h
repowiki <command> --help
```

### Examples

```bash
repowiki help
repowiki --help
repowiki enable --help
```

### Output Example

```
repowiki v0.1.0 — Auto-generate Qoder repo wiki on git commits

Usage:
  repowiki <command> [flags]

Commands:
  enable      Enable repowiki in current project (install git hook)
  disable     Disable repowiki (remove git hook)
  status      Show current status and configuration
  generate    Run full wiki generation
  update      Run incremental wiki update for recent changes
  version     Show version

Flags for 'enable':
  --force             Reinstall hook even if already present
  --qodercli-path     Path to qodercli binary
  --model             Qoder model level (auto, efficient, performance, ultimate)
  --no-auto-commit    Don't auto-commit wiki changes

Flags for 'update':
  --commit            Specific commit hash to process
  --from-hook         Internal: indicates hook-triggered run

Examples:
  repowiki enable                    # Enable in current project
  repowiki enable --force            # Reinstall hook
  repowiki generate                  # Full wiki generation
  repowiki update --commit abc123    # Update for specific commit
  repowiki logs                      # View most recent log
  repowiki disable                   # Remove hook
```

## Internal Command: hooks

**Note**: This command is for internal use by the git hook and should not be called directly.

### Synopsis

```
repowiki hooks post-commit
```

### Description

Entry point called by the git post-commit hook. Runs loop prevention checks and spawns a background update process.

### Flow

1. Validate argument is `post-commit`
2. Find git repository root
3. Check sentinel file (loop prevention)
4. Check lock file (concurrency prevention)
5. Load configuration
6. Get current commit hash
7. Check commit message prefix (loop prevention)
8. Spawn background update process

### Background Process

```go
cmd := exec.Command(self, "update", "--from-hook", "--commit", commitHash)
cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
cmd.Start()
```

The background process:
- Runs detached from the parent (using `Setsid: true`)
- Outputs to `.repowiki/logs/hook.log`
- Does not block the user's terminal

## logs

Show the most recent log file from hook executions.

### Synopsis

```
repowiki logs
```

### Description

Displays the contents of the most recent log file. Logs are sorted by filename in descending order, showing the newest log first.

### Examples

```bash
repowiki logs
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (not a git repository, failed to read log) |

### Output Example

```
=== hook.log ===
2026-02-19T19:30:00Z Starting wiki update for commit abc123...
2026-02-19T19:30:05Z Changed files: 5
2026-02-19T19:30:10Z Running incremental update...
2026-02-19T19:30:45Z Wiki update complete
2026-02-19T19:30:46Z Auto-committing changes...
2026-02-19T19:30:47Z Done
```

### Log Location

Logs are stored in `.repowiki/logs/` directory. Common log files include:
- `hook.log` - Output from hook-triggered updates
- Date-stamped logs for specific runs (e.g., `2026-02-19.log`)

### Implementation

The command reads all files from `.repowiki/logs/`, sorts them by name descending (newest first), and displays the contents of the first (most recent) log file.
