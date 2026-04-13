Repository: Todo App (Go) — Hexagonal/Clean architecture

Build, test, and run

- Makefile targets (top-level):
  - make build    # builds binary to bin/todoapp
  - make test     # runs go test ./...
  - make run      # runs the app via go run ./cmd/app
- Single-test examples:
  - go test ./internal/adapters/http -run TestIntegration_TodosEndpoints -v
  - go test ./internal/adapters/http -run TestIntegration_TodosEndpoints
  - go test ./internal/adapters/http -run TestName
- Docker / compose:
  - ./scripts/run.sh        # runs docker-compose up --build
  - docker-compose up --build

Environment (for local run)

- The binary expects DB connection via env vars (see cmd/app/main.go):
  DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
- docker-compose.yaml provides a Postgres service and mounts deploy/postgres/init.sql to initialize the DB.

High-level architecture (big picture)

- Hexagonal / Clean structure under /internal:
  - internal/domain    -> core entities (Todo)
  - internal/repository -> port interfaces (TodoRepository)
  - internal/usecase   -> business logic / services (TodoService)
  - internal/infra     -> infrastructure adapters (postgres, inmemory)
  - internal/adapters  -> delivery adapters (http server)
- cmd/app is the application bootstrap: it reads env, creates DB connection, wires infra -> usecase -> adapters, starts HTTP server.
- Tests: adapters/http/server_test.go shows an integration-style test using the inmemory repo and httptest.Server.

Key conventions and patterns

- Ports & Adapters:
  - Define repository interfaces in internal/repository. Implement adapters in internal/infra/*.
  - Usecase constructors accept repository interfaces (dependency injection by constructor).
- Naming:
  - "Repo"/"Repository" types live in infra/inmemory or infra/postgres and provide NewTodoRepository constructors.
  - Usecase service is named TodoService and exposes Create and List methods.
- Tests:
  - Integration-like tests live next to adapters (see internal/adapters/http/server_test.go) and use inmemory repo when possible.
- Database migration/init:
  - deploy/postgres/init.sql is mounted into the Postgres container (docker-compose).

Existing automation & CI hooks

- No lint target or CI config detected in the repo root. Use make test / go test for verification.

Other assistant / AI configs found

- No CLAUDE.md, .cursorrules, AGENTS.md, CONVENTIONS.md, .windsurfrules, or existing .github/copilot-instructions.md were found.

Notes for future Copilot sessions

- Focus on the dependency flow: cmd -> infra (DB) -> usecase -> adapters (HTTP).
- When adding new persistence, add a new adapter under internal/infra and ensure it implements the internal/repository interface.
- For new endpoints, add delivery code under internal/adapters/* and test using inmemory.NewTodoRepository to avoid DB dependencies.

MCP servers

Would you like help configuring any MCP servers for this project (examples: Postgres test server, Playwright for web UI tests)?

Summary

Created .github/copilot-instructions.md with build/test commands, architecture overview, and repository-specific conventions. Want any adjustments or extra coverage (more test examples, CI/lint suggestions, or script docs)?