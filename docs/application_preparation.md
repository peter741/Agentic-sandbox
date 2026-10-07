# Codex for Open Source application preparation

Program: https://developers.openai.com/community/codex-for-oss

The public page invites core maintainers and people running widely used public
projects. Smaller projects with an important ecosystem role can explain it.
No numerical minimum guaranteeing approval is published there. Check current
program terms when applying.

## Verified project facts
- Repository: https://github.com/peter741/Agentic-sandbox
- Independent Apache 2.0 derivative of Agents Sandbox, maintained as peter741/Agentic-sandbox; see provenance.md for the upstream source.
- New change set: read-only setup diagnostics, regression tests, Python smoke
  helper fix, repository configuration, and explicit provenance.
- Changes prepared with Codex assistance.
- Local checks and limitations: [validation.md](validation.md).

## Still needed from the maintainer
- Review changes and run them on your actual machine.
- Commit/push and confirm GitHub CI.
- Run appropriate live runtime checks before publishing a release.
- Record actual users, feedback, issues resolved, and contributions.
- Explain practical ecosystem value using verifiable examples.

No users, downloads, stars, external contributors, or accepted upstream PRs
have been verified for this derivative. No ongoing maintenance history is
claimed yet. Do not use upstream metrics or artificial activity as evidence.

## Application outline
1. Role: your actual maintenance responsibilities, with links to evidence.
2. Project: purpose, upstream attribution, and exactly which new behavior you maintain.
3. Impact: real adoption or the important unmet need; describe a new project honestly.
4. Codex plan: concrete tasks such as reproducing setup bugs, adding regressions,
   reviewing contributions, and improving platform diagnostics, plus validation.

Finalize answers after gathering the missing evidence. This document does not
establish eligibility or say that six months of ChatGPT Pro has been granted.

## Additional functional contribution

`agbox preflight` adds an offline workspace exposure review before copying a
project into an agent sandbox. It flags credential-looking names and symlinks,
reports incomplete scans, supports JSON, and has four regression tests. This
is a derivative addition, not a claim of novelty across all open-source tools.
It does not inspect secret contents or automatically prevent agent launches.
