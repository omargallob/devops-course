# DevOps Course

An interactive, open-source web course for learning Go as a DevOps engineer. Built with a Go backend, Astro frontend, and an in-browser Go playground for hands-on exercises.

## Overview

This project provides a structured curriculum covering Go topics relevant to DevOps: CLI tools, Kubernetes operators, Terraform providers, observability, and more. Lessons are authored in MDX and served by Astro, while a Go backend proxies code execution to the Go Playground so learners can run code directly in the browser.

## Architecture

```
devops-course/
├── apps/web/            # Astro frontend (MDX content, Tailwind, Vitest)
├── cmd/server/          # Go API server entrypoint
├── internal/
│   ├── api/             # chi router, middleware, handlers
│   ├── auth/            # GitHub OAuth + JWT sessions (milestone 5)
│   ├── database/        # GORM models and Postgres connection
│   ├── exercises/       # Exercise validation (milestone 2)
│   └── playground/      # Go Playground compile proxy
│   └── platform/        # Build constraints (linux/darwin × amd64/arm64)
├── bazel/               # Custom Bazel rules and macros
├── infra/docker/        # Production Dockerfiles, nginx, compose
├── tools/               # Build scripts (workspace_status.sh)
├── docker-compose.yml   # Dev environment (hot reload + Postgres)
├── Makefile             # Convenience targets
└── MODULE.bazel         # Bazel 9 module configuration
```

**Backend:** Go with chi router, structured logging (slog), graceful shutdown, and GORM for Postgres.

**Frontend:** Astro with MDX content collections, Tailwind CSS v4, and Vite dev server with API proxy.

**Build:** Bazel 9 for hermetic builds, cross-compilation, and container images (`rules_img`). Makefile as a convenience layer.

## Prerequisites

- **Go** 1.24+
- **Node.js** 22 LTS (see `.nvmrc`)
- **pnpm** (enabled via corepack)
- **Docker** and **Docker Compose** (for the dev environment)
- **Bazel** 9.1.0 (optional, pinned in `.bazelversion`)

## Quick Start

### Docker Compose (recommended)

```bash
# Copy env file
cp .env.example .env

# Start all services (Go API + Astro dev + Postgres)
docker compose up

# Or via Make
make docker-dev
```

- **Web:** http://localhost:4321
- **API:** http://localhost:8080
- **Postgres:** localhost:5432 (user: `devops`, password: `devops`, db: `devops_course`)

Go file changes trigger automatic rebuild via [Air](https://github.com/air-verse/air). Astro file changes trigger HMR.

### Local (without Docker)

```bash
# Install dependencies
make setup

# Start both servers (Go backend + Astro frontend)
make dev
```

Without `DATABASE_URL` set, the Go server starts in DB-less mode. Set it to connect to a local Postgres instance:

```bash
export DATABASE_URL="postgres://devops:devops@localhost:5432/devops_course?sslmode=disable"
```

## Make Targets

| Target | Description |
|--------|-------------|
| `make dev` | Start Go backend + Astro frontend locally |
| `make docker-dev` | Start full dev environment via Docker Compose |
| `make docker-dev-down` | Stop Docker Compose services |
| `make build` | Build all targets with Bazel |
| `make test` | Run all tests with Bazel |
| `make test-go` | Run Go tests only |
| `make test-web` | Run frontend tests (Vitest) |
| `make lint` | Run all linters via Bazel |
| `make fmt` | Format Go and frontend code |
| `make image` | Build container image with Bazel |
| `make push` | Push container image to ghcr.io |
| `make gazelle` | Regenerate BUILD.bazel files |
| `make db-shell` | Open psql shell in the Postgres container |
| `make db-reset` | Wipe database volume and restart |
| `make setup` | Install pre-commit hooks, pnpm deps |
| `make clean` | Clean Bazel cache and build artifacts |

## Development

### Project Layout

| Directory | Purpose |
|-----------|---------|
| `apps/web/src/content/` | MDX lesson files, organised by module |
| `apps/web/src/pages/` | Astro page routes |
| `apps/web/src/utils/` | Frontend utilities (playground client, module data) |
| `cmd/server/` | Server entrypoint with graceful shutdown |
| `internal/api/` | HTTP router, middleware, route handlers |
| `internal/database/` | GORM models (User, Session, LessonProgress, ExerciseSubmission) |
| `internal/playground/` | Go Playground compile proxy |
| `bazel/rules/` | Custom lint rules (golangci-lint, shellcheck, yamllint, actionlint, eslint) |

### Database

The backend uses PostgreSQL with GORM. Tables are auto-migrated on startup:

- **users** -- GitHub-authenticated course participants
- **sessions** -- JWT token references for revocation
- **lesson_progress** -- per-user lesson completion tracking
- **exercise_submissions** -- submitted code and results

### Linting

Linters run through Bazel test targets and pre-commit hooks:

- **Go:** golangci-lint v2
- **YAML:** yamllint
- **GitHub Actions:** actionlint
- **Frontend:** ESLint (typescript-eslint + eslint-plugin-astro)
- **Commits:** commitlint (conventional commits)

### CI

GitHub Actions runs on every push and PR: lint, test, build, Bazel verification, and container image build. Images are pushed to `ghcr.io/omargallob/devops-course/server` on merges to `main`.

## Contributing

1. Fork the repo and create a feature branch
2. Follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages
3. Run `make lint` and `make test` before submitting a PR
4. Open a PR against `main`

See [open issues](https://github.com/omargallob/devops-course/issues) for planned work.

## License

[MIT](LICENSE)
