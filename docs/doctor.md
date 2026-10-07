# Diagnose setup with agbox doctor

When an agent cannot start, check setup before changing it. Doctor queries Docker
and sends a daemon Ping; it does not create containers, copy workspaces, start
services, modify permissions, or run an agent.

```bash
agbox doctor
agbox doctor --workspace ./my-project
agbox doctor --json --timeout 3s
```

## Checks

| Check | Verifies |
| --- | --- |
| platform | Linux or macOS host |
| docker_cli | Docker executable on PATH |
| docker_engine | Docker API responds and reports Linux containers |
| daemon | Configured Unix-socket daemon answers Ping |
| daemon_version | Matching release versions; development or differences warn |
| workspace | Optional directory exists and passes root/home path policy, including symlinks |

Docker discovery uses environment-based API configuration, including DOCKER_HOST
and Docker TLS settings, as the runtime does. It does not infer the active Docker
CLI context. Ensure the daemon's own environment reaches the intended engine.
A remote Linux engine does not prove local host-mount compatibility. Workspace
checks do not verify writable permissions or inspect contents. Agent
authentication is not tested.

Each network probe has its own positive timeout (default 5s). Probes run
sequentially, so total network time may be approximately twice that timeout.
Cancellation is honored. Raw probe errors, endpoint details, and credentials are
not included. Review accompanying logs before posting them.

## Output and exit codes

- 0: no check failed; warnings can remain.
- 1: a check failed or the report could not be written.
- 2: invalid usage, including zero/negative timeout.

JSON contains version, healthy, and checks. Each check has name, status
(pass/warn/fail), message, and optional hint. A failing JSON report still
appears on stdout, without an additional diagnostic error on stderr.
Healthy means connectivity checks passed, not that sandbox security is certified.

## Troubleshooting

- Missing CLI: install Docker and reopen the terminal.
- Unreachable engine: start Docker and inspect daemon access and Docker settings.
  Do not make the Docker socket world-writable.
- Linux runtime path: use a normal login session with XDG_RUNTIME_DIR.
- Linux daemon: inspect `systemctl --user status agboxd` and
  `journalctl --user -u agboxd`. CAP_NET_ADMIN is required; review capability
  requirements and the installer before granting it.
- macOS daemon: inspect launchctl and ~/Library/Logs/agboxd.log if installed
  by the local installer.
- Version warning: build/install CLI and daemon from the same checkout.

See [security limitations](../SECURITY.md).
