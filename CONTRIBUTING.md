# Contributing to Agentic Sandbox

Contribute to [Agentic Sandbox](https://github.com/peter741/Agentic-sandbox), an independent derivative of Agents Sandbox. See [provenance](docs/provenance.md).
Read [provenance](docs/provenance.md) before describing inherited work as a new contribution.

## Development setup

- Go 1.25+, as required by go.mod.
- Python 3.12+ and uv for Python SDK tests.
- Docker and a configured agboxd service only for live integration tests.

Build without installing or modifying services:

```bash
mkdir -p .build
go build -buildvcs=false -o .build/agbox ./cmd/agbox
go build -buildvcs=false -o .build/agboxd ./cmd/agboxd
.build/agbox doctor --help
```

## Before opening a pull request

1. Describe a concrete user problem and keep the change focused.
2. Add regression tests for meaningful behavior changes.
3. Update user documentation and CHANGELOG.md.
4. Retain licenses and credit outside work.
5. Report checks run, failures, and skips.

```bash
go test ./... -count=1
go vet ./...
cd sdk/python
uv run pytest tests/ --ignore=tests/test_real_runtime.py --ignore=tests/test_network_isolation.py
```

Format changed Go files with gofmt. CI runs unit/mock tests on Linux and macOS.
Passing mocks does not prove container or network isolation. Before a runtime
release, maintainers should run `make integration-test` on a configured Linux
host. It starts real containers and may require CAP_NET_ADMIN for the test daemon;
inspect the suite before running it.

Use this repository's bug form with a minimal reproduction. Remove credentials
and private paths from logs. Link upstream issues if the problem reproduces there.

Maintainers must review AI-assisted code, reproduce issues, respond to contributors,
and publish accurate release notes. Inherited history and upstream users are not
this derivative's own maintenance or adoption.
