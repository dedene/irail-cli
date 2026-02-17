# Repository Guidelines

## Project Structure

- `cmd/irail/`: CLI entrypoint
- `internal/`: implementation
  - `cmd/`: command routing (kong CLI framework, 13 command files)
  - `api/`: iRail API client + rate limiting
  - `tui/`: terminal UI (Charmbracelet bubbletea)
  - `output/`: time formatting + output rendering
  - `errfmt/`: error formatting
- `bin/`: build outputs

## Build, Test, and Development Commands

- `make build`: compile to `bin/irail`
- `make install`: install to system
- `make fmt` / `make lint` / `make test` / `make ci`: format, lint, test, full local gate
- `make tools`: install pinned dev tools into `.tools/`
- `make clean`: remove bin/ and .tools/

**Note**: Tests run with `-race` flag for race condition detection.

## Coding Style & Naming Conventions

- Formatting: `make fmt` (goimports local prefix `github.com/dedene/irail-cli` + gofumpt)
- Output: keep stdout parseable (`--json`); send human hints/progress to stderr
- Linting: golangci-lint with project config
- TUI: use Charmbracelet ecosystem (bubbletea)

## Testing Guidelines

- Unit tests: stdlib `testing` with `-race` flag
- Coverage areas: stations, time formatting, rate limiting
- 3 test files; focus on core logic

## Config & Secrets

- **No authentication**: iRail is a public API
- **Stateless**: no config file needed
- **Rate limiting**: built-in to prevent API abuse

## Key Commands

- `stations`: station lookups and search
- `liveboard`: live train departures/arrivals
- `connections`: train connections between stations
- `composition`: train composition details
- `vehicle`: vehicle/train information
- `disturbances`: service disruptions

## Commit & Pull Request Guidelines

- Conventional Commits: `feat|fix|refactor|build|ci|chore|docs|style|perf|test`
- Group related changes; avoid bundling unrelated refactors
- PR review: use `gh pr view` / `gh pr diff`; don't switch branches

## Security Tips

- No credentials to manage (public API)
- Rate limiting protects against accidental API abuse
