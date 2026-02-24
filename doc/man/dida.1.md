# dida365-cli 1 "Feb 2026" dida365-cli "User Manuals"

## NAME

`dida365-cli` - Dida365 command line interface

## SYNOPSIS

`dida365-cli` [GLOBAL OPTIONS] _COMMAND_ [OPTIONS]

## DESCRIPTION

`dida365-cli` is a Go CLI for Dida365 OpenAPI.

Current implementation provides:
- OAuth setup/login/token exchange/status
- project list/create
- task list/show/add/update/done/delete
- token refresh command (best-effort; depends on server response capabilities)
- read-command cache (10-second window for repeated GET requests)
- write-command debounce (3-second window for duplicate write requests)

## GLOBAL OPTIONS

`--config <PATH>`
: Path to config file. Default is `~/.config/dida365-cli/config.json`.

`--debug`
: Enable debug output.

`--dry-run`
: Print request intent without sending real API requests.

`--no-cache`
: Disable local read cache. Same as setting `DIDA_NO_CACHE=1`.

`-h, --help`
: Show help.

## COMMANDS

## `dida365-cli auth`

OAuth setup and token management.

### `dida365-cli auth init`

Save OAuth client settings.

Required flags:
- `--client-id`
- `--client-secret`
- `--redirect-uri`

Optional flags:
- `--json`

### `dida365-cli auth login`

Print authorization URL (user opens it in browser to authorize).

Optional flags:
- `--scope` (default: `tasks:read tasks:write`)
- `--state` (default: auto-generated)
- `--json`

### `dida365-cli auth token`

Exchange authorization code for access token and save into config.

Required flags:
- `--code`

Optional flags:
- `--scope`
- `--json`

### `dida365-cli auth status`

Show current config/token status:
- config path
- api base url
- oauth fields present status
- access token presence
- token expiry/scope (if available)

Optional flags:
- `--json`

### `dida365-cli auth refresh`

Try refreshing access token with refresh token.

Optional flags:
- `--refresh-token` (default: read from config)
- `--scope`
- `--json`

Note:
- If server response does not include a refresh token, this command may be unavailable in practice.
- In that case, use `dida365-cli auth login` + `dida365-cli auth token --code ...` again.

### `dida365-cli auth logout`

Clear stored token in config.

Optional flags:
- `--json`

## `dida365-cli projects list`

List projects from Dida API.

Optional flags:
- `--json`

## `dida365-cli projects create`

Create project in Dida API.

Required flags:
- `--name`

Optional flags:
- `--color`
- `--sort-order`
- `--view-mode` (`list|kanban|timeline`)
- `--kind` (`TASK|NOTE`)
- `--json`

## `dida365-cli list`

List tasks in a project.

Default table columns include: ID, title, completed, due date, priority, content preview.

Required flags:
- `--project`

Optional flags:
- `--format table|json`
- `--json` (alias style; compatible with `--format json`)

## `dida365-cli show`

Show one task detail in a project.

Default table output includes: content, completed state, completed time, due date, priority.

Required flags:
- `--project`
- `--id`

Optional flags:
- `--json`

Optional flags:
- `--format table|json`
- `--json`

## `dida365-cli add`

Create task.

Input limits:
- `--title`, `--content`, `--desc` must be `< 500` characters.
- Limit applies to create input only, not display.

Required flags:
- `--project`
- `--title`

Optional flags:
- `--content`
- `--desc`
- `--start`
- `--due`
- `--repeat`
- `--timezone`
- `--priority`
- `--all-day`
- `--json`

## `dida365-cli update`

Update task.

Required flags:
- `--project`
- `--id`

Optional flags:
- `--json`

Optional flags:
- `--title`
- `--content`
- `--desc`
- `--start`
- `--due`
- `--repeat`
- `--timezone`
- `--priority`
- `--all-day`
- `--json`

## `dida365-cli done`

Complete a task.

Required flags:
- `--project`
- `--id`

## `dida365-cli delete`

Delete a task.

Required flags:
- `--project`
- `--id`

## `dida365-cli version`

Print CLI version.

Optional flags:
- `--json`

## EXAMPLES

```bash
# init oauth settings
dida365-cli auth init --client-id <id> --client-secret <secret> --redirect-uri http://127.0.0.1:53682/callback

# print auth url
dida365-cli auth login

# exchange code
dida365-cli auth token --code <authorization_code>

# check status
dida365-cli auth status

# list projects
dida365-cli projects list

# create project
dida365-cli projects create --name "Roadmap" --view-mode list --kind TASK

# list tasks in project
dida365-cli list --project <project_id>

# show task
dida365-cli show --project <project_id> --id <task_id>

# create task
dida365-cli add --project <project_id> --title "Buy milk"

# complete task
dida365-cli done --project <project_id> --id <task_id>
```

## FILES

- Default config: `~/.config/dida365-cli/config.json`

## SEE ALSO

- `README.md`
- `docs/guides/dida-auth-token-flow.md`
- `docs/research/dida-api-overview.md`
