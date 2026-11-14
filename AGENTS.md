# AGENTS.md

## Git Worktrees

- Always create worktrees inside the main repository at `.local/worktrees/<name>`, never in `/tmp` or outside the repository.
- Use relative Git worktree links (`--relative-paths`) so worktrees remain accessible from both the host and the devcontainer.

## Coding Conventions

- **Comments**: Keep comments to an absolute minimum, preferably none.

### Backend (Go)

- **Error Handling**: Use custom error types in `domain/error.go`; wrap errors with `go-errors/errors` for stack traces.
- **Context**: Always pass `context.Context` as the first parameter to functions that perform I/O.
- **Interfaces**: Define dependencies as interfaces; implement in other packages.
- **Testing**: Unit tests alongside implementation files (`*_test.go`).
- **Logging**: Use structured logging with `slog`.

### Frontend (Svelte)

- **Svelte 5**: Use runes API (`$state`, `$derived`).
- **Effects**: Always avoid `$effect` unless absolutely necessary; prefer lifecycle hooks such as `onMount` for setup and cleanup.
- **API Calls**: Use TanStack Query through queries in the shared library.
- **Conditionals**: Always use braces for one-line conditionals; do not use shorthand `if (...) return ...` style.
- **Styling**: Component-scoped styles in plain CSS (prefer nesting).
- **Web Awesome**: Use the `webawesome` skill for Web Awesome components and styling.
- **Design Tokens**: Use Web Awesome design tokens whenever possible.
- **Theme**: Shared theme in `packages/lib/src/theme.css`.
- **Zod**: Always import as `import { z } from "@climblive/lib/utils"`.

## Testing Guidelines

### Backend

- Unit tests for use cases, domain logic, and utilities.
- Table-driven tests preferred for multiple test cases.
- Use `stretchr/testify` for assertions.
- Mock interfaces using `stretchr/testify/mock`.
- Use `testing/synctest` (via `synctest.Test`) whenever a test involves `time.Now()`.

### Frontend

- E2E tests using Playwright in `web/e2e`.
- Run `make test` for full E2E suite.

## Common Tasks

### Modifying database schema

1. Update model in MySQL Workbench (`backend/database/climblive.mwb`) and forward export to `backend/database/climblive.sql`.
2. Create a new Goose migration in `backend/cmd/api/migrations/`.
3. Update `backend/database/queries.sql` with new queries.
4. Run `sqlc generate`.

**Important:** Always update the MySQL Workbench model file (`.mwb`) when making database schema changes.

### Modify domain models

1. Update domain models in `backend/internal/domain/public.go`.
2. Run `tygo generate` to generate updated TypeScript types.

## When making changes

1. Follow existing code style and patterns.
2. Ensure database changes include migrations.
3. Run `sqlc generate` if database model or queries are updated.
4. Run `tygo generate` if either of `public.go` or `id.go` is updated.
5. Keep comments to a minimum.
6. Make sure that `golangci-lint run` and `go test ./...` pass.
7. Make sure that `pnpm check` and `pnpm lint` pass.
8. Make sure to format front-end code using `pnpm format`.

For questions about specific subsystems, examine existing tests and implementation files in the relevant package.

## Pull Requests

- **Titles**: Should be formatted in the Conventional Commits format using only lowercase letters.
- **Descriptions**: Keep descriptions brief and on a high-level.
