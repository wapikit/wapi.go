# Repository Guidelines

## Project Structure & Module Organization
Wapi.go targets Go 1.21. Core SDK packages live under `pkg/` (messaging, events, business), shared HTTP plumbing sits in `internal/`, and orchestration helpers live in `manager/`. Examples are in `examples/` for quick sandbox runs. Generated references belong in `docs/` with templates in `docs/templates`. Tests sit beside implementations as `*_test.go`.

## Build, Test, and Development Commands
- `go build ./...`: compiles all packages to catch cross‑package breakage early.
- `go test ./...`: runs the standard test suite; narrow with `-run`.
- `make format`: invokes `go fmt ./...` to match canonical Go style.
- `make docs`: installs `gomarkdoc` if absent and regenerates `docs/api-reference/`.
- `go run ./examples/chat-bot` or `go run ./examples/http-backend-integration`: validate changes against sample flows.

## Coding Style & Naming Conventions
- Idiomatic Go: tabs and `gofmt` output.
- Exported identifiers use PascalCase with package‑level doc comments.
- Unexported helpers use camelCase and remain package‑local.
- Files/dirs are lowercase with hyphens. Prefer narrow, WhatsApp‑specific structs and reuse constructors from `manager/` and `pkg/`.

## Testing Guidelines
- Use the standard `testing` package with table‑driven `TestXxx` functions.
- Stub outbound HTTP by wrapping `internal/request_client` (no live endpoints).
- Cover serialization, validation, and webhook dispatch; ensure `go test ./...` is clean.

## Commit & Pull Request Guidelines
- Conventional Commits per `COMMIT_CONVENTION.md`, e.g., `feat(messaging): add template sender`.
- Branch names: `type/topic` (e.g., `fix/webhook-validation`).
- PRs summarize behavior, link issues, attach payloads/screenshots for API changes, and mention any doc updates (`make docs`).

## Documentation Workflow
- API references are generated artifacts. Update doc comments and templates, not generated files.
- After changes, run `make docs` and commit both source and generated markdown.

## Security & Configuration Tips (Optional)
- Provide API tokens via environment or secret management; never commit credentials.
- Store webhook secrets securely and rotate regularly.

