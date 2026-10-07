# Publish these changes

GitHub is not changed until you commit and push the source deliverable.

1. Extract the updated ZIP into a new folder.
2. Clone https://github.com/peter741/Agentic-sandbox.git into a separate folder.
3. Copy extracted contents, including .github, .gitignore,
   .pre-commit-config.yaml and .githooks, into that clone. Keep its .git directory.
4. Remove the old agents-sandbox-main.zip if still tracked.
5. Review the diff, commit, and push normally; do not force-push.

Inside the clone after copying:

```bash
git rm --ignore-unmatch agents-sandbox-main.zip
git add -A
git diff --cached --stat
git diff --cached
git commit -m "Add setup diagnostics and document upstream provenance"
git push origin main
```

Confirm the Test workflow succeeds in Actions. No successful hosted CI run is
claimed yet. Run a locally built doctor binary on your own machine. Before a
runtime release, run live integration tests on a configured Linux host and
record limitations. Enable private vulnerability reporting and document a
monitored contact. Update application notes with real evidence only.

## Namespace transition

Build the CLI and daemon together after these changes. The Go import path is
`github.com/peter741/Agentic-sandbox`. Docker ownership labels now use
`io.github.peter741.agentic-sandbox`. Existing upstream-labeled resources will
not be adopted or cleaned up by this version. Stop and remove existing
sandboxes using the previous installation before switching; do not run two
daemon versions against the same socket/state store. The previous installation
remains responsible for its own resources. Preserve backups of existing state
before changing installations. Upstream runtime images remain dependencies;
no images have been published under peter741 automatically.

Website and documentation configurations no longer claim upstream domains or
Cloudflare projects. Set your actual deployment URL when hosting is configured.
The installer now targets this GitHub repository and needs release assets;
use the documented source-build path until those exist.
