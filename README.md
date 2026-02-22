# dida365-cli

`dida365-cli` is a Go CLI for Dida365 OpenAPI.

## Requirements
- Go >= 1.23 (tested with Go 1.24.x)
- A Dida365 developer app (`client_id`, `client_secret`, `redirect_uri`)

## Build

```bash
go build -o dida ./cmd/dida
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

- `dida auth init|login|token|refresh|status|logout`
- `dida projects list`
- `dida projects create --name "<project_name>" [flags]`
- `dida list --project <project_id>`
- `dida show --project <project_id> --id <task_id>`
- `dida add --project <project_id> --title "..."`
- `dida update --project <project_id> --id <task_id> [flags]`
- `dida done --project <project_id> --id <task_id>`
- `dida delete --project <project_id> --id <task_id>`
- `dida version`

Use `--dry-run` to preview requests without executing.

## Config

Default config path:

- `~/.config/dida365-cli/config.json`
- This project uses a fixed HOME-based default path on all platforms (does not use OS-specific `UserConfigDir`).

Override config path:

```bash
dida --config /path/to/config.json auth status
```

Override access token directly:

```bash
DIDA_ACCESS_TOKEN=<token> dida projects list
```
