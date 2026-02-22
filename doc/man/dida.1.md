# dida 1 "Feb 2026" dida365-cli "User Manuals"

## NAME

`dida` - Dida365 command line interface

## SYNOPSIS

`dida` [GLOBAL OPTIONS] _COMMAND_ [OPTIONS]

## DESCRIPTION

`dida` is a Go CLI for Dida365 OpenAPI.

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

## `dida auth`

OAuth setup and token management.

### `dida auth init`

Save OAuth client settings.

Required flags:
- `--client-id`
- `--client-secret`
- `--redirect-uri`

Optional flags:
- `--json`

### `dida auth login`

Print authorization URL (user opens it in browser to authorize).

Optional flags:
- `--scope` (default: `tasks:read tasks:write`)
- `--state` (default: auto-generated)
- `--json`

### `dida auth token`

Exchange authorization code for access token and save into config.

Required flags:
- `--code`

Optional flags:
- `--scope`
- `--json`

### `dida auth status`

Show current config/token status:
- config path
- api base url
- oauth fields present status
- access token presence
- token expiry/scope (if available)

Optional flags:
- `--json`

### `dida auth refresh`

Try refreshing access token with refresh token.

Optional flags:
- `--refresh-token` (default: read from config)
- `--scope`
- `--json`

Note:
- If server response does not include a refresh token, this command may be unavailable in practice.
- In that case, use `dida auth login` + `dida auth token --code ...` again.

### `dida auth logout`

Clear stored token in config.

Optional flags:
- `--json`

## `dida projects list`

List projects from Dida API.

Optional flags:
- `--json`

## `dida projects create`

Create project in Dida API.

Required flags:
- `--name`

Optional flags:
- `--color`
- `--sort-order`
- `--view-mode` (`list|kanban|timeline`)
- `--kind` (`TASK|NOTE`)
- `--json`

## `dida list`

List tasks in a project.

Default table columns include: ID, title, completed, due date, priority, content preview.

Required flags:
- `--project`

Optional flags:
- `--format table|json`
- `--json` (alias style; compatible with `--format json`)

## `dida show`

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

## `dida add`

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

## `dida update`

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

## `dida done`

Complete a task.

Required flags:
- `--project`
- `--id`

## `dida delete`

Delete a task.

Required flags:
- `--project`
- `--id`

## `dida version`

Print CLI version.

Optional flags:
- `--json`

## EXAMPLES

```bash
# init oauth settings
dida auth init --client-id <id> --client-secret <secret> --redirect-uri http://127.0.0.1:53682/callback

# print auth url
dida auth login

# exchange code
dida auth token --code <authorization_code>

# check status
dida auth status

# list projects
dida projects list

# create project
dida projects create --name "Roadmap" --view-mode list --kind TASK

# list tasks in project
dida list --project <project_id>

# show task
dida show --project <project_id> --id <task_id>

# create task
dida add --project <project_id> --title "Buy milk"

# complete task
dida done --project <project_id> --id <task_id>
```

## FILES

- Default config: `~/.config/dida365-cli/config.json`

## SEE ALSO

- `README.md`
- `docs/guides/dida-auth-token-flow.md`
- `docs/research/dida-api-overview.md`
