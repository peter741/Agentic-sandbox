# Changelog

Derivative changes only. The sandbox engine, SDKs, integrations, images, and
earlier tests were inherited from upstream.

## Unreleased

### Added
- `agbox preflight`: offline, metadata-only workspace exposure review with
  relative-path JSON, symlink warnings, bounded entry inspection, and
  incomplete-scan reporting. Includes four regression tests.
- Read-only `agbox doctor` setup diagnostics, optional workspace validation,
  bounded probes, actionable hints, and JSON output.
- Regression tests for failures, timeout handling, workspace path policy,
  gRPC Ping, and Docker HTTP probes.
- Provenance, contribution/security guidance, bug/PR templates, and
  application preparation notes.
- Linux/macOS CI for Go/Python unit and mock tests and Go vet.

### Changed
- Moved Go module/imports and generated protocol metadata to
  `github.com/peter741/Agentic-sandbox`.
- Updated website/docs/release-installer links to the derivative repository.
- Moved Docker ownership labels to the derivative namespace; existing upstream
  resources are not adopted automatically.
- Removed optional upstream-specific guardrail hooks, retaining generic hooks
  and gofmt. Preserved original attribution and license.
- README identifies the derivative and explains building its own CLI.
- CLI reference documents doctor and exit codes.
- Restored omitted development files and executable script permissions.
- Python smoke helper uses its explicit test socket even without
  XDG_RUNTIME_DIR; added a regression test.
- Added archive/credential ignore rules and removed the redundant source ZIP.

### Status
No tagged derivative release, external adoption, package publication, or live
container-isolation validation is claimed. Hosted CI must run after pushing.
See docs/validation.md for local results.
