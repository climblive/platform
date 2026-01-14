# AGENTS.md

This document provides AI coding agents and developers with essential information about the ClimbLive platform codebase.

## Project Overview

ClimbLive is a bouldering competition scoring system consisting of:

- **Backend**: Go-based REST API server with real-time event streaming
- **Frontend**: Multiple Svelte 5 applications served through a monorepo structure
  - `admin`: Competition administration interface
  - `scoreboard`: Live scoreboard display for competitions
  - `scorecard`: Scorecard interface for competitors
  - `www`: Public website

The platform manages bouldering climbing competitions, calculates scores in real-time, and provides live updates to multiple clients, including the scoreboard and scorecard. Scores are based on the information entered by the competitiors them self using the scorecard app during competitions.

## Architecture

### Backend (`/backend`)

**Language**: Go 1.25  
**Database**: MySQL  
**Key Features**: Event-driven architecture, real-time score calculation, JWT authentication

**Structure**:

```
backend/
├── cmd/api/              # API server entry point
├── cmd/simulator/        # Testing/simulation tools
├── internal/
│   ├── authorizer/       # JWT authentication & authorization
│   ├── database/         # SQLC-generated database layer
│   ├── domain/           # Core domain models and interfaces
│   ├── events/           # Event broker for real-time updates
│   ├── handlers/rest/    # HTTP handlers
│   ├── repository/       # Repository implementations
│   ├── scores/           # Score engine implementation
│   ├── usecases/         # Business logic layer
│   └── utils/            # Shared utilities
└── database/             # SQL schema and queries
```

**Key Patterns**:

- **Clean Architecture**: Innermost layers are made up of the domain followed by the use cases layer
- **Repository Pattern**: Data access abstracted through repository interfaces
- **Event Broker**: Pub/sub system for real-time scoring events, etc. (`events.EventBroker`)
- **Score Engine**: Separate engine that manages score calculations
- **Use Cases**: Business logic written almost completely using the Go standard library
- **SQLC**: Type-safe SQL queries generated from `queries.sql`

**Core Domain Entities**:

- `Contest`: A climbing competition
- `CompClass`: Competition class/category (e.g., Males, Females)
- `Contender`: A contender participating in a contest
- `Problem`: A boulder problem set for the competition
- `Tick`: An attempt at a problem by a contender
- `Score`: Calculated score for a contender

### Frontend (`/web`)

**Framework**: Svelte 5  
**Build Tool**: Vite  
**Package Manager**: pnpm  
**Node Version**: ≥20

**Structure**:

```
web/
├── admin/              # Admin application
├── scoreboard/         # Live scoreboard
├── scorecard/          # Contender scorecard
├── www/                # Public website (marketing)
├── packages/lib/       # Shared library (@climblive/lib)
└── e2e/                # End-to-end tests (Playwright)
```

**Shared Library** (`packages/lib`):

- Shared UI components
- API client (`Api.ts`)
- TypeScript models (generated from Go backend via tygo)
- TanStack Query queries
- Model validators built in zod
- Form utilities
- Theme and styling

**Key Technologies**:

- **Remote State Management**: TanStack Query
- **Routing**: svelte-routing
- **UI Components**: @awesome.me/webawesome
- **Validation**: Zod
- **Error Tracking**: Sentry

## Development Setup

### Prerequisites

- Go 1.25+
- Node.js 20+
- pnpm
- MySQL database
- Docker (for containerized development)

### Backend Setup

```bash
cd backend

# Install dependencies
go mod download

# Set environment variables
export DB_USERNAME=your_db_user
export DB_PASSWORD=your_db_password
export DB_HOST=localhost
export DB_PORT=3306
export DB_DATABASE=climblive

# Run migrations (automatic on startup)
# Run the API server
go run cmd/api/main.go
```

The API server runs on `http://0.0.0.0:8090` by default.

**Database Tools**:

- Migrations are embedded in the binary (`cmd/api/migrations/*.sql`) and run automatically using Goose
- SQLC generates type-safe Go code from SQL queries in `database/queries.sql`
- Run `sqlc generate` to regenerate database code after modifying queries

### Frontend Setup

```bash
cd web

# Install dependencies
pnpm -r i

# Run specific app in development mode
cd admin && pnpm dev
# or
cd scoreboard && pnpm dev
# or
cd scorecard && pnpm dev
# or
cd www && pnpm dev

# Build for production
pnpm build
```

Development servers run with `--host` flag for network access.

### End-to-End Testing

```bash
cd web/e2e

# Install Playwright
pnpm exec playwright install --with-deps

# Run tests
make test
```

## Coding Conventions

- **Comments**: Keep comments to an absolute minimum, preferrably none

### Backend (Go)

- **Error Handling**: Use custom error types in `domain/error.go`; wrap errors with `go-errors/errors` for stack traces
- **Context**: Always pass `context.Context` as the first parameter to functions that perform I/O
- **Interfaces**: Define dependencies as interfaces; implement in other packages
- **Testing**: Unit tests alongside implementation files (`*_test.go`)
- **Logging**: Use structured logging with `slog`

### Frontend (Svelte)

- **Svelte 5**: Use runes API (`$state`, `$derived`, `$effect`)
- **TypeScript**: Strict mode enabled
- **API Calls**: Use TanStack Query through queries in the shared library
- **Styling**: Component-scoped styles in plain CSS (prefer nesting)
- **Tokens**: Use Web Awesome design tokens for all styling (`@awesome.me/webawesome/dist/styles/themes/default.css`)
- **Theme**: Shared theme in `packages/lib/src/theme.css`
- **Formatting**: Prettier with project-specific config
- **Linting**: ESLint with Svelte plugin

**Shared Code**: Common utilities, types, and components belong in `packages/lib` for reuse across apps.

## Testing Guidelines

### Backend

- Unit tests for use cases, domain logic, and utilities
- Table-driven tests preferred for multiple test cases
- Use `stretchr/testify` for assertions
- Mock interfaces using `stretchr/testify/mock`

### Frontend

- E2E tests using Playwright in `web/e2e`
- Run `make test` for full E2E suite

## Deployment

- Everything packaged as a single Debian package (see `packageroot/DEBIAN/`)
- Backend lifecycle managed as systemd service
- Static frontend builds served via Nginx (see `packageroot/etc/nginx/sites-available`)

## Common Tasks

### Modifying database schema

1. Update model in MySQL Workbench (`backend/database/climblive.mwb`) and forward export to `backend/database/climblive.sql`
2. Create a new Goose migration in `backend/cmd/api/migrations/`
3. Update `backend/database/queries.sql` with new queries
4. Run `sqlc generate`

**Important:** Always update the MySQL Workbench model file (`.mwb`) when making database schema changes. This ensures the visual model stays in sync with the actual database schema.

### Modify domain models

1. Update domain models in `backend/internal/domain/public.go`
2. Run `tygo generate` to generate updated TypeScript types

## Authentication

- **Organizers**: Admin users authenticate using JWT tokens issued by Cognito
- **Contenders**: Scorecard users authenticate using their unique alphanumeric registration code

## When making changes

1. Follow existing code style and patterns
2. Ensure database changes include migrations
3. Run `sqlc generate` if database model or queries are updated
4. Run `tygo generate` if either of `public.go` or `id.go` is updated
5. Keep comments to a minimum
6. Make sure that `golangci-lint run` and `go test ./...` pass
7. Make sure that `pnpm check` and `pnpm lint` pass
8. Make sure to format front-end code using `pnpm format`

For questions about specific subsystems, examine existing tests and implementation files in the relevant package.

## Pull Requests

- **Titles**: Should be formatted in the Conventional Commits format using only lowercase letters
- **Descriptions**: Keep descriptions brief and on a high-level
