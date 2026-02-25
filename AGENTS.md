# Repository Guidelines

## Project Structure & Module Organization
- `cmd/structify`: CLI entrypoint for schema/repository generation.
- `cmd/server`: HTTP server exposing generation endpoints used by the web UI.
- `internal/`: core implementation (`parser`, `adapter`, `application`, `generator`, `domain`, `mapper`, `util`).
- `pkg/cli`: CLI command parsing and app wiring.
- `web/nextjs`: frontend app (Next.js + React + TypeScript).
- `test/fixtures`, `test/expected`, `test/integration`: backend fixtures, golden outputs, and integration tests.
- `examples/`: sample domains (`blog`, `ecommerce`) for generated SQL and repository output.

## Build, Test, and Development Commands
- `go build ./cmd/structify`: build the CLI binary.
- `go test ./...`: run all Go unit tests.
- `go test -tags=integration ./test/integration/...`: run backend integration tests.
- `go run ./cmd/server`: start the backend API locally.
- `npm --prefix web/nextjs install`: install frontend dependencies.
- `npm --prefix web/nextjs run dev`: run frontend in development mode.
- `npm --prefix web/nextjs run test`: run frontend unit tests (Vitest).
- `npm --prefix web/nextjs run lint`: run frontend ESLint checks.

## Coding Style & Naming Conventions
- Go: use `gofmt` formatting, idiomatic Go naming (`CamelCase` exports, `camelCase` internals), and table-driven tests where useful.
- TypeScript/React: follow ESLint rules in `web/nextjs/eslint.config.mjs`; prefer functional components and explicit prop types.
- Generated repository files use `*.gen.go` (for example `user_repository.gen.go`).
- Tests use `*_test.go` (Go) and `*.test.ts(x)` (frontend).

## Testing Guidelines
- Backend: `testing` package with unit tests across `internal/`, `pkg/`, and `cmd/`.
- Frontend: Vitest + Testing Library in `web/nextjs/__tests__`.
- Keep or improve coverage when changing parsing, SQL generation, or docs/UI behavior.
- Add regression tests with any bug fix.

## Commit & Pull Request Guidelines
- Follow Conventional Commit style seen in history: `feat:`, `fix:`, `refactor:`, `test:`, optionally scoped (for example `feat(frontend): ...`).
- PRs should include a concise problem/solution summary, linked issue (if available), and test evidence (`go test`, `vitest`, or both).
- Include UI screenshots or GIFs for frontend changes under `web/nextjs`.
- Keep diffs focused; avoid mixing unrelated refactors with functional changes.

## Security & Configuration Tips
- Do not commit secrets. Use environment variables for deployment/runtime config.
- Validate untrusted inputs in parser/server paths and keep dependency updates current via CI/Dependabot.
