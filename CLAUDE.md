# Project Guidelines

## Repository Overview

This is an interactive web course for learning Go as a DevOps engineer. It's a monorepo with:

- **`apps/web/`** — Astro frontend (MDX lessons, Tailwind CSS v4, TypeScript)
- **`cmd/server/`** — Go backend (chi router, GORM, PostgreSQL)
- **`internal/`** — Go packages (api, auth, database, exercises, playground, platform)
- **`packages/contracts/`** — API contracts (OpenAPI specs, JSON schemas, generated code)
- **`infra/`** — Production Dockerfiles, nginx, compose
- **`e2e/`** — Cypress end-to-end tests
- **`bazel/`** — Custom Bazel rules and macros

**Stack:** Go 1.26, Node 22, pnpm 11, Bazel 9, PostgreSQL, Docker Compose

---

## Build & Development Commands

| Task | Command |
|------|---------|
| Full build | `make build` (runs `bazel build //...`) |
| Run all tests | `make test` |
| Go tests only | `make test-go` |
| Frontend tests | `make test-web` |
| Lint everything | `make lint` |
| Lint locally (pre-commit) | `make lint-local` |
| Format code | `make fmt` |
| Dev servers (Go + Astro) | `make dev` |
| Setup environment | `make setup` |
| Docker dev environment | `make docker-dev` |
| Build container image | `make image` |
| Regenerate BUILD files | `make gazelle` |

---

## Commit Conventions

This repo uses **Conventional Commits** enforced by commitlint.

Format: `<type>(<scope>): <description>`

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`

Examples:
```
feat(api): add exercise submission endpoint
fix(web): correct lesson navigation ordering
docs: update contract spec for auth endpoints
ci: add contract validation step
```

---

## Contract Driven Development (CDD)

All **new** API endpoints and features MUST follow Contract Driven Development. Existing endpoints are migrated opportunistically.

### Principles

1. **Contract First** — Write or update the OpenAPI spec before writing implementation code.
2. **Single Source of Truth** — The contract in `packages/contracts/` defines the API. Implementation must conform to it, not the other way around.
3. **Generated Types** — Server and client types are generated from contracts. Never hand-write types that duplicate contract definitions.
4. **CI Enforced** — Contracts are validated, generated code freshness is checked, and breaking changes are detected in CI.

### Contract Location

Contracts live in `packages/contracts/` (a pnpm workspace package):

```
packages/contracts/
├── openapi/
│   ├── api.yaml              # Main OpenAPI 3.x spec (or split per domain)
│   └── components/
│       └── schemas/          # Reusable schema definitions
├── schemas/                  # Standalone JSON schemas (non-API validation)
├── generated/
│   ├── go/                   # oapi-codegen output (server types + interfaces)
│   └── ts/                   # openapi-typescript output (client types)
├── package.json
└── README.md
```

### Workflow for New Endpoints

1. **Define the contract** — Add/update the OpenAPI spec in `packages/contracts/openapi/`.
2. **Run code generation** — Generate Go server types and TypeScript client types.
3. **Implement the handler** — Write the Go handler conforming to the generated interface.
4. **Consume from frontend** — Import generated TS types in the Astro frontend.
5. **Write tests** — Validate request/response against the contract schema.

### Code Generation

| Target | Tool | Output |
|--------|------|--------|
| Go server types & interfaces | `oapi-codegen` | `packages/contracts/generated/go/` |
| TypeScript client types | `openapi-typescript` | `packages/contracts/generated/ts/` |

Generation commands (to be added to Makefile):
```bash
# Generate Go types from OpenAPI spec
oapi-codegen -generate types,server -package api \
  packages/contracts/openapi/api.yaml > packages/contracts/generated/go/api.gen.go

# Generate TypeScript types from OpenAPI spec
pnpm openapi-typescript packages/contracts/openapi/api.yaml \
  -o packages/contracts/generated/ts/api.ts
```

### JSON Schema Usage

Use JSON schemas in `packages/contracts/schemas/` for:
- Request body validation (middleware-level, runtime)
- Configuration file validation
- Exercise input/output format definitions

Validate at runtime using a JSON Schema validator (e.g., `gojsonschema` in Go, `ajv` in TypeScript).

### CI Enforcement

The CI pipeline enforces contract compliance:

1. **Spec Linting** — Validate OpenAPI specs with `spectral` or `vacuum`.
2. **Codegen Freshness** — Regenerate code in CI and fail if there's a diff (ensures generated files are committed and up-to-date).
3. **Breaking Change Detection** — Use `oasdiff` to detect breaking changes on PRs. Breaking changes require explicit approval.
4. **Schema Validation in Tests** — Integration tests validate actual responses against the OpenAPI spec.

### Rules

- Never modify generated files by hand. Always update the contract and regenerate.
- Every PR that adds/changes an API endpoint MUST include the corresponding contract update.
- Breaking changes to contracts require a migration plan and explicit reviewer approval.
- Keep contracts as the shared language between frontend and backend engineers.

---

## Pull Request Guidelines

- PRs target `main` branch.
- CI must pass (lint, test, build, bazel, e2e).
- For API changes: include contract diff, regenerated types, and implementation in the same PR.
- Keep PRs focused — separate contract additions from unrelated refactors.

---

## Architecture Decisions

- **Bazel** is the primary build system. Use `make` targets as convenience wrappers.
- **pnpm workspaces** manage the monorepo's Node packages.
- **PostgreSQL** with GORM for persistence; auto-migration in dev, explicit migrations in prod.
- **Air** for Go hot-reload in development.
- **Pre-commit hooks** enforce formatting and linting before commits land.
