# repowiki

Auto-generate [Qoder](https://qoder.com) repo wiki on git commits — like [Entire](https://entire.io) but for documentation.

## How It Works

1. Install `repowiki` and run `repowiki enable` in your project
2. Every git commit automatically triggers wiki generation in the background
3. Updated wiki is auto-committed with `[repowiki]` prefix
4. Your `.qoder/repowiki/` stays in sync with your code

## Requirements

- Git repository with at least one commit
- [Qoder](https://qoder.com) IDE or CLI installed (`qodercli` binary)

## Installation

```bash
# From source
go install github.com/ikrasnodymov/repowiki/cmd/repowiki@latest

# Or build locally
git clone https://github.com/ikrasnodymov/repowiki.git
cd repowiki
make install
```

## Usage

```bash
# Enable in your project
cd /path/to/your/project
repowiki enable

# Generate wiki for the first time
repowiki generate

# Check status
repowiki status

# Disable
repowiki disable
```

After `enable`, every commit auto-updates the wiki in the background. No action needed.

## Commands

| Command | Description |
|---------|-------------|
| `repowiki enable` | Install post-commit hook and create config |
| `repowiki disable` | Remove hook, preserve wiki files |
| `repowiki status` | Show configuration and status |
| `repowiki generate` | Full wiki generation from scratch |
| `repowiki update` | Incremental update for recent changes |
| `repowiki version` | Show version |

## Configuration

Config is stored in `.repowiki/config.json`:

```json
{
  "enabled": true,
  "qodercli_path": "qodercli",
  "model": "auto",
  "max_turns": 50,
  "language": "en",
  "auto_commit": true,
  "commit_prefix": "[repowiki]",
  "wiki_path": ".qoder/repowiki",
  "full_generate_threshold": 20
}
```

## How It Prevents Infinite Loops

Wiki commits could trigger the hook again. Three layers prevent this:

1. **Sentinel file** — `.repowiki/.committing` exists during wiki commit
2. **Lock file** — `.repowiki/.repowiki.lock` prevents concurrent runs
3. **Commit prefix** — commits starting with `[repowiki]` are skipped

## License

MIT
