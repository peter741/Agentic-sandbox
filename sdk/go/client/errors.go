// Agentic Sandbox derivative: updated repository namespace; see NOTICE.
package client

import "github.com/peter741/Agentic-sandbox/sdk/go/rawclient"

// Type aliases so callers can use errors.As without importing rawclient directly.
type SandboxClientError = rawclient.SandboxClientError
type SandboxConflictError = rawclient.SandboxConflictError
type SandboxNotFoundError = rawclient.SandboxNotFoundError
type SandboxNotReadyError = rawclient.SandboxNotReadyError
type SandboxInvalidStateError = rawclient.SandboxInvalidStateError
type ExecNotFoundError = rawclient.ExecNotFoundError
type ExecAlreadyTerminalError = rawclient.ExecAlreadyTerminalError
type ExecNotRunningError = rawclient.ExecNotRunningError
type SandboxSequenceExpiredError = rawclient.SandboxSequenceExpiredError
