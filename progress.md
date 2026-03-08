# Progress

- [x] Inspect current codebase and release workflow.
- [x] Add/maintain task tracker files (`progress.md`, `todo.md`).
- [x] Update release pipeline to publish zip artifact containing binary.
- [x] Inject Git tag into build and wire `hubfly-storage version` to print version only.
- [x] Implement resilient FileBrowser bootstrap and `.env` bootstrapping.
- [x] Extend `/health` endpoint with FileBrowser status/version details.
- [x] Run validation/build checks.
- [x] Create atomic commits per logical change.

## Log

- 2026-03-08: Read `cmd/hubfly-storage/main.go`, `handlers/handlers.go`, and `.github/workflows/release.yml`.
- 2026-03-08: Added tracker files (`progress.md`, `todo.md`).
- 2026-03-08: Updated release workflow to build with injected tag version and upload zip artifact.
- 2026-03-08: Added startup `.env` bootstrap and non-fatal FileBrowser admin password rotation with PM2 integration.
- 2026-03-08: Extended `/health` response to include Hubfly and FileBrowser runtime info.
- 2026-03-08: Validation passed via `go build ./...` and explicit `version` command check.
