# Project provenance

Agentic Sandbox is an independent derivative of
[1996fanrui/agents-sandbox](https://github.com/1996fanrui/agents-sandbox), Apache 2.0.

The initial uploaded source was compared to upstream commit
`f53c0c0413c579962cd2ad2c50e5edd8f9e2494d`: all 216 shared files were byte-identical.
This identifies the comparison baseline; it does not claim the original download
had recorded that commit.

Inherited work includes the control plane, Docker isolation logic, Go/Python
SDKs, agent integrations, runtime images, website, and existing tests/docs.
The original LICENSE is unchanged, and the [original README](upstream_readme.md)
is retained.

See [CHANGELOG.md](../CHANGELOG.md) for derivative changes. The first functional
addition is read-only setup diagnostics through agbox doctor. Repository setup
and documentation were repaired, and the Python smoke-test environment helper
was fixed. A subsequent addition, agbox preflight, reviews workspace filenames
for potential credential exposure without reading contents or following links.

Changes were prepared with OpenAI Codex assistance. The maintainer must review,
test on their own setup, publish, and maintain them. No claim of independent
handwritten authorship is made.

The Go module path now uses github.com/peter741/Agentic-sandbox. The Python
package name, executable names, and runtime image references remain for compatibility; they do not establish ownership of
upstream distribution channels. Build this checkout to get doctor. Upstream
installers and package channels do not distribute it automatically. Do not
publish derivative packages into upstream namespaces. No upstream publishing
workflows are enabled.

Applications must distinguish new changes from inherited work, this project's
users from upstream users, actual maintenance from plans, and local checks from
hosted CI and live Docker validation.

Repository links, Go imports/protocol package metadata, release installer target,
and Docker ownership labels were adapted to the derivative namespace.
Inherited notices and this source history retain original attribution.
