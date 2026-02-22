# Repository Guidelines

## Project Structure & Module Organization
This repository is a Go CLI for Dida365.

- `cmd/dida/main.go`: application entrypoint.
- `internal/cli/`: Cobra command definitions and command handlers (`add`, `list`, `auth`, `projects`, etc.).
- `internal/dida/`: Dida365 API client and request/response logic.
- `internal/config/`: local config loading and persistence.
- `internal/output/`: output formatting helpers.
- `docs/research/`: reference notes and API mapping research.

Keep new runtime code under `internal/` and wire user-facing commands through `internal/cli/`.

## Build, Test, and Development Commands
Use the Makefile targets for common workflows:

- `make build`: compile `./cmd/dida` to local binary `./dida`.
- `make run`: run CLI help quickly (`go run ./cmd/dida --help`).
- `make test`: run all unit tests (`go test ./...`).
- `make fmt`: format source (`gofmt -w cmd internal`).
- `make tidy`: clean and sync module dependencies.

Direct equivalents also work, e.g. `go test ./...` during iteration.

## Coding Style & Naming Conventions
- Follow standard Go formatting and idioms; run `make fmt` before opening a PR.
- Use tabs/`gofmt` defaults (do not hand-align spacing).
- Package names are short, lowercase, and singular where practical.
- File names in `internal/cli/` should map to command purpose (for example, `done.go`, `projects.go`).
- Keep command flags and UX behavior consistent with existing Cobra patterns in `internal/cli/root.go`.

## Testing Guidelines
- Place tests next to implementation files with `_test.go` suffix.
- Prefer table-driven tests for command formatting, API behavior, and error handling.
- Keep tests deterministic; avoid live API calls in unit tests.
- Run `make test` before committing; add tests for bug fixes and new command behavior.

## Commit & Pull Request Guidelines
Current history follows Conventional Commit-style prefixes:
- `feat: ...`
- `fix: ...`
- `docs: ...`
- `chore: ...`

Use short imperative subjects and keep each commit scoped to one logical change.  
PRs should include:
- clear summary of behavior changes,
- linked issue or rationale,
- test evidence (for example, `make test` output),
- CLI usage examples when command UX changes.

## Security & Configuration Tips
- Never commit tokens, client secrets, or local config files.
- Default config path is `~/.config/dida365-cli/config.json`; prefer env vars for temporary credentials (for example, `DIDA_ACCESS_TOKEN`).
