# Agentic Sandbox

An independent, early-stage open-source project maintained at [peter741/Agentic-sandbox](https://github.com/peter741/Agentic-sandbox), derived from Agents Sandbox, for running coding agents in local Docker sandboxes.

**Upstream provides the sandbox engine and agent integrations. This derivative adds setup diagnostics through `agbox doctor` and a metadata-only workspace exposure review through `agbox preflight`.** See [provenance](docs/provenance.md) and [CHANGELOG.md](CHANGELOG.md).

## Diagnose setup before launching an agent

Check missing Docker, an unreachable engine, and daemon setup in one command:

```bash
agbox doctor
agbox doctor --workspace ./my-project
agbox doctor --json --timeout 3s
```

Doctor checks Linux/macOS host support, Docker CLI availability, Docker Engine connectivity and Linux-container mode, daemon Ping, version differences, and optionally the workspace path. It does not start services, create containers, or run an agent. See [doctor documentation](docs/doctor.md).

## Review a workspace before giving it to an agent

```bash
agbox preflight --workspace ./my-project
agbox preflight --workspace ./my-project --json
```

Preflight flags credential-looking filenames, private-key names, credential
folders, and symlinks without reading file contents or changing files. It works
without Docker. Incomplete scans and findings exit 1 for review. This is an
advisory filename scan, not a security guarantee or an automatic launch gate.
See [preflight documentation](docs/preflight.md).

## Build this derivative

Requirements: Go 1.25 or newer. Docker and a configured agboxd daemon are needed for live sandbox use, but not for help or unit/mock tests.

```bash
git clone https://github.com/peter741/Agentic-sandbox.git
cd Agentic-sandbox
mkdir -p .build
go build -buildvcs=false -o .build/agbox ./cmd/agbox
go build -buildvcs=false -o .build/agboxd ./cmd/agboxd
.build/agbox doctor --help
.build/agbox doctor --json
```

Without Docker and the daemon configured, failed diagnostics are expected. For local installation, review `scripts/install_local.sh`: it builds both executables, sets up a user service, can edit a shell profile, and on Linux requests CAP_NET_ADMIN for the daemon. Doctor itself makes no such changes.

Build this checkout to get doctor. Upstream installers and packages do not automatically include derivative changes. The Go module uses this repository’s namespace. Existing runtime image references remain upstream dependencies; derivative images and packages have not been published.

## Inherited Python SDK

With the SDK installed from `sdk/python` and the daemon configured:

```python
import asyncio
from agents_sandbox import AgentsSandboxClient

async def main():
    async with AgentsSandboxClient() as client:
        print(await client.ping())

asyncio.run(main())
```

See [Python SDK usage](docs/sdk_python_usage.md). This API is inherited from upstream; doctor is a CLI addition.

## Validation and contributing

- [Contributing](CONTRIBUTING.md): setup, tests, and review expectations.
- [Validation](docs/validation.md): checks performed and limitations.
- [CLI reference](docs/cli_reference.md): commands and flags.
- [Security](SECURITY.md): reporting arrangements and limitations.
- [Publish these changes](docs/publish_changes.md): commit and push without losing history.

```bash
go test ./... -count=1
go vet ./...
```

Python SDK checks require Python 3.12+ and uv. Unit/mock tests do not validate live container or network isolation.

## Attribution and license

The inherited code, SDKs, images, website, and earlier tests remain upstream work. New changes were prepared with Codex assistance and require maintainer review.

Licensed under [Apache 2.0](LICENSE). See [NOTICE](NOTICE), [provenance](docs/provenance.md), and the [original README](docs/upstream_readme.md). This is not an official upstream release.
