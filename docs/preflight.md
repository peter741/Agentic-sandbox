# Workspace exposure preflight

`agbox preflight` reviews filenames before you copy a workspace into an agent
sandbox. Unlike doctor, it needs no Docker, daemon, or network connection.

```bash
agbox preflight --workspace ./my-project
agbox preflight --workspace ./my-project --json
agbox preflight --workspace ./my-project --max-entries 50000
```

The default workspace is the current directory. Filesystem root and your home
directory are rejected using the existing workspace policy. The command does
not start an agent or change files. Review findings, prepare a clean copy, rerun
preflight, and then explicitly choose that copy when starting an agent.

## What it detects

- Environment-file names (`.env`, `.env.*`), including template/example files.
- Common private-key names and `.pem`, `.key`, `.p12`, `.pfx` suffixes.
- Credential filenames and authentication configuration such as `.npmrc`.
- Credential directories: `.ssh`, `.aws`, `.azure`, `.kube`, `.gnupg`.
- All symlinks, without following their targets.
- Unreadable entries, interrupted scans, and entry-limit exhaustion.

Findings are conservative prompts for human review. For example, a public
certificate or placeholder `.env.example` may be harmless. The scan reads
filesystem metadata and directory listings, not file contents. It includes
hidden and ignored directories; it does not assume Git ignore rules match the
runtime's workspace copy behavior. Reports use paths relative to the workspace
and never include file contents. Reports may still disclose sensitive filenames.

## Exit codes and JSON

Exit 0 means the metadata scan completed with no findings. Exit 1 means findings
need review or the scan was incomplete. Exit 2 indicates invalid arguments or
workspace path policy rejection. Output failures are runtime errors.

JSON contains `complete`, `needs_review`, `entries_scanned`, and `findings`.
Each finding contains `path`, `kind`, and `advice`. Incomplete reports never
represent a clean scan. The default limit is 10,000 inspected entries;
`--max-entries` accepts 1 through 1,000,000. Directory enumeration itself may
read additional names before the entry limit is reached.

## Limits

This is an advisory command, not an automatic launch gate. It does not scan
secret values, prove that a file is sensitive, or certify sandbox isolation.
Credentials embedded in ordinary source files can go undetected. Files can
change between review and launch. External mounts, extra `--copy` paths, agent
authentication sources, and Docker configuration are outside this workspace
review. Symlink targets are neither resolved nor certified safe; actual copy
behavior must be reviewed separately. Run doctor for setup diagnostics and use
a content-aware secret scanner when you need content inspection.
