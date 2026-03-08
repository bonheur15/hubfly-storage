# Progress

- [x] Inspect current codebase and release workflow.
- [ ] Add/maintain task tracker files (`progress.md`, `todo.md`).
- [ ] Update release pipeline to publish zip artifact containing binary.
- [ ] Inject Git tag into build and wire `hubfly-storage version` to print version only.
- [ ] Implement resilient FileBrowser bootstrap and `.env` bootstrapping.
- [ ] Extend `/health` endpoint with FileBrowser status/version details.
- [ ] Run validation/build checks.
- [ ] Create atomic commits per logical change.

## Log

- 2026-03-08: Read `cmd/hubfly-storage/main.go`, `handlers/handlers.go`, and `.github/workflows/release.yml`.
