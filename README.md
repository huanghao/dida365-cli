# dida365-cli

`dida365-cli` is a Go CLI for Dida365 OpenAPI.

## Requirements
- Go >= 1.23 (tested with Go 1.24.x)
- A Dida365 developer app (`client_id`, `client_secret`, `redirect_uri`)
- This repo uses `goenv` to manage Go versions. Initialize `goenv` in your shell before running build/test commands:

```bash
export GOENV_ROOT="$HOME/.goenv"
export PATH="$GOENV_ROOT/bin:$PATH"
eval "$(goenv init -)"
go version
```

## Build

```bash
go build -o dida365-cli ./cmd/dida
```

## Run

```bash
go run ./cmd/dida --help
```

## Quick Start (OAuth)

1. Initialize OAuth app config:

```bash
go run ./cmd/dida auth init \
  --client-id <client_id> \
  --client-secret <client_secret> \
  --redirect-uri <redirect_uri>
```

2. Get authorization URL and open it in browser:

```bash
go run ./cmd/dida auth login
```

3. Copy `code` from redirect URL and exchange token:

```bash
go run ./cmd/dida auth token --code <authorization_code>
```

4. Check auth status:

```bash
go run ./cmd/dida auth status
```

## Commands

- `dida365-cli auth init|login|token|refresh|status|logout`
- `dida365-cli projects list`
- `dida365-cli projects create --name "<project_name>" [flags]`
- `dida365-cli list --project <project_id>`
- `dida365-cli show --project <project_id> --id <task_id>`
- `dida365-cli add --project <project_id> --title "..."`
- `dida365-cli update --project <project_id> --id <task_id> [flags]`
- `dida365-cli done --project <project_id> --id <task_id> [--json]`
- `dida365-cli delete --project <project_id> --id <task_id> [--json]`
- `dida365-cli version [--json]`

Use `--dry-run` to preview requests without executing.
For agent workflows, prefer `--json` on all actionable commands.
Use `--no-cache` (or `DIDA_NO_CACHE=1`) to disable local read cache.

`dida365-cli list` table output includes frequently-used fields by default:
- completion state
- due date
- priority
- content (truncated preview)

Create input limit:
- `dida365-cli add` enforces text length `< 500` characters for `--title`, `--content`, `--desc`.
- This only limits CLI create input; existing longer task content can still be displayed.

Write debounce:
- write commands (`add/update/done/delete/projects create`) are debounced for 3 seconds.
- repeated identical write requests in the short window will be blocked to reduce accidental duplicate operations.

Read cache:
- read commands (`projects list`, `list`, `show`) use a local 10-second cache.
- any successful write command clears read cache to reduce stale reads.

## Agent Docs

- Quick manual: `docs/guides/agent-cli-quick-manual.md`
- Maintenance rules: `docs/guides/agent-manual-rules.md`
- Onboarding: `docs/guides/agent-project-onboarding.md`
- Release: `docs/guides/release-homebrew.md`

## Collaboration Process

- Workflow: `docs/process/COLLAB-WORKFLOW.md`
- Task board: `docs/process/TASK-BOARD.md`
- Handoff: `docs/process/HANDOFF.md`
- Design gates: `docs/process/design-gates/README.md`

## Config

Default config path:

- `~/.config/dida365-cli/config.json`
- `~/.config/dida365-cli/cache.json` (read cache state)
- `~/.config/dida365-cli/debounce.json` (write debounce state)
- This project uses a fixed HOME-based default path on all platforms (does not use OS-specific `UserConfigDir`).

Override config path:

```bash
dida365-cli --config /path/to/config.json auth status
```

Override access token directly:

```bash
DIDA_ACCESS_TOKEN=<token> dida365-cli projects list
```

## License

MIT. See `LICENSE`.
