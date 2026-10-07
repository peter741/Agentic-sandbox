# Validation of this change set

Validated on Linux with Go 1.25.0 and Python 3.12.14.

| Check | Result |
| --- | --- |
| CLI and daemon builds | Passed using go build -buildvcs=false |
| Go vet ./... | Passed |
| Go formatting | Passed |
| New doctor tests, excluding live Unix-socket Ping | 11 passed |
| Python SDK tests excluding socket-server tests and real-runtime suites | 68 passed, 1 pre-existing skip, 9 deselected |
| Workflow/template YAML parsing | Passed |
| CLI JSON failure report | Valid JSON, exit 1, no extra stderr |
| Invalid zero timeout | Exit 2 |

## Full-suite limits

The full Go suite was attempted and failed on Unix-socket listener restrictions
(operation not permitted), including one new live Ping test and inherited CLI
and socket validation tests. Other Go packages completed successfully. The full
non-Docker Python suite was also attempted: 9 socket-server tests failed to bind,
68 passed, and one existing real-daemon test was skipped. These are recorded
failures, not a claim that the full suite passed.

Restricted runs above exclude socket-dependent tests explicitly. No production
security checks or test assertions were disabled in the delivered source.
GitHub CI still runs the full unit/mock suites without these exclusions.
Confirm a green CI run after pushing.

No real Docker/container/network-isolation tests or macOS runtime checks were
performed here. Run make integration-test on a configured host before a runtime
release. Go/Python mocks and doctor connectivity do not certify isolation.
No tagged release or external adoption has been verified for this derivative.

## Subsequent preflight addition

All four new preflight regression tests pass. They cover metadata findings,
unchanged file contents, relative-path JSON, symlink non-traversal, scan limits,
cancellation, clean exit, usage errors, and invalid file workspaces. The CLI
was rebuilt and Go vet rerun successfully. Built-CLI JSON smoke checks pass
for both a clean directory (exit 0) and a `.env` filename (exit 1), with no
secret contents or extra stderr. The earlier full-suite limitations still apply.

## Repository namespace update

Both Go binaries rebuild and Go vet passes after changing the Go module/imports
and regenerating Go/Python protocol package metadata. Docker-label tests and
both Go SDK packages pass. All 15 runnable doctor/preflight tests pass (the live
Unix-socket Ping test remains excluded). A targeted Python configuration,
conversion, and environment-helper suite passes: 37 tests. Bash syntax checks
and website JavaScript parsing pass. Source links for original attribution
remain in NOTICE and provenance.md; operational repository references use the
derivative repository. No GitHub push, deployment, release upload, or image
publication was performed. Earlier full-suite/runtime limits still apply.
