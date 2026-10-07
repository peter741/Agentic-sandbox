from unittest.mock import patch

import agents_sandbox.client as client_module
from tests.smoke_support import _new_client


def test_explicit_test_socket_does_not_require_login_runtime_directory(monkeypatch):
    monkeypatch.delenv("XDG_RUNTIME_DIR", raising=False)
    with patch.object(client_module, "SandboxGrpcClient") as raw_client:
        client = _new_client("/tmp/isolated-test.sock")
        assert client.socket_path == "/tmp/isolated-test.sock"
        assert all(call.args[0] == "/tmp/isolated-test.sock" for call in raw_client.call_args_list)
        client.close()
