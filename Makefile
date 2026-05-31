.PHONY: all build test lint fmt dev setup clean image push docker-dev docker-dev-down db-shell db-reset \
	generate generate-go generate-ts lint-contracts

# ─── Default ──────────────────────────────────────────────────────────────────
all: lint test build

# ─── Build ────────────────────────────────────────────────────────────────────
build:
	bazel build //...

build-server:
	bazel build //cmd/server

# ─── Test ─────────────────────────────────────────────────────────────────────
test:
	bazel test //...

test-go:
	bazel test //internal/... //cmd/...

test-web:
	cd apps/web && pnpm test

# ─── Lint ─────────────────────────────────────────────────────────────────────
lint:
	bazel test //:golangci_lint //:yamllint //:prettier //:eslint

lint-go:
	bazel test //:golangci_lint

lint-web:
	bazel test //:prettier //:eslint

lint-yaml:
	bazel test //:yamllint

lint-actions:
	bazel test //:actionlint

lint-local:
	pre-commit run --all-files

# ─── Format ───────────────────────────────────────────────────────────────────
fmt:
	go fmt ./...
	cd apps/web && pnpm format

fmt-check:
	bazel test //:prettier

# ─── Dev ──────────────────────────────────────────────────────────────────────
dev:
	@echo "Starting Go backend and Astro frontend..."
	@trap 'kill 0' EXIT; \
		(cd apps/web && pnpm dev) & \
		(air -c .air.toml 2>/dev/null || go run ./cmd/server) & \
		wait

dev-web:
	cd apps/web && pnpm dev

dev-server:
	air -c .air.toml 2>/dev/null || go run ./cmd/server

# ─── Container Image ─────────────────────────────────────────────────────────
image:
	bazel build //cmd/server:image

push:
	bazel run //cmd/server:push

# ─── Gazelle ──────────────────────────────────────────────────────────────────
gazelle:
	bazel run //:gazelle

gazelle-update:
	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies

# ─── Setup ────────────────────────────────────────────────────────────────────
setup:
	pip install pre-commit yamllint
	pre-commit install --hook-type commit-msg --hook-type pre-commit
	cd apps/web && pnpm install
	@echo "Setup complete. Run 'make dev' to start developing."

# ─── Clean ────────────────────────────────────────────────────────────────────
clean:
	bazel clean
	rm -rf apps/web/dist apps/web/.astro

# ─── Contracts / Code Generation ──────────────────────────────────────────────
generate: generate-go generate-ts

generate-go:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen \
		--config packages/contracts/oapi-codegen.yaml \
		packages/contracts/openapi/api.yaml

generate-ts:
	cd packages/contracts && pnpm generate:ts

lint-contracts:
	cd packages/contracts && pnpm lint

# ─── Docker Dev ───────────────────────────────────────────────────────────────
docker-dev:
	docker compose up --build

docker-dev-down:
	docker compose down

db-shell:
	docker compose exec db psql -U devops devops_course

db-reset:
	docker compose down -v
	docker compose up -d db
	@echo "Database volume wiped. Run 'make docker-dev' to restart all services."
