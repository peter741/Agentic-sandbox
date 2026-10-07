// Agentic Sandbox derivative: updated repository namespace; see NOTICE.
package main

import (
	"context"
	agboxv1 "github.com/peter741/Agentic-sandbox/api/generated/agboxv1"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"
)

type doctorPingServer struct {
	agboxv1.UnimplementedSandboxServiceServer
}

func (doctorPingServer) Ping(context.Context, *agboxv1.PingRequest) (*agboxv1.PingResponse, error) {
	return &agboxv1.PingResponse{Version: "doctor-test"}, nil
}
func TestDoctorProbeDaemonUsesActualPing(t *testing.T) {
	socket, _ := startSandboxTestServer(t, doctorPingServer{})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	got, e := probeDaemonVersion(ctx, socket)
	if e != nil || got != "doctor-test" {
		t.Fatalf("unexpected result %q %v", got, e)
	}
}
func TestDoctorProbeUnavailableDaemonRespectsDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	if _, e := probeDaemonVersion(ctx, filepath.Join(t.TempDir(), "missing.sock")); e == nil {
		t.Fatal("missing daemon must fail")
	}
}
func TestDoctorProbeDockerReadsEngineOS(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("API-Version", "1.51")
		if r.URL.Path == "/_ping" {
			_, _ = w.Write([]byte("OK"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"OSType":"linux"}`))
	}))
	defer server.Close()
	t.Setenv("DOCKER_HOST", server.URL)
	t.Setenv("DOCKER_TLS_VERIFY", "")
	t.Setenv("DOCKER_CERT_PATH", "")
	t.Setenv("DOCKER_API_VERSION", "")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	got, e := probeDockerOS(ctx)
	if e != nil || got != "linux" {
		t.Fatalf("unexpected result %q %v", got, e)
	}
}
func TestDoctorProbeDockerRespectsDeadline(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { <-r.Context().Done() }))
	defer server.Close()
	t.Setenv("DOCKER_HOST", server.URL)
	t.Setenv("DOCKER_TLS_VERIFY", "")
	t.Setenv("DOCKER_CERT_PATH", "")
	t.Setenv("DOCKER_API_VERSION", "1.51")
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	if _, e := probeDockerOS(ctx); e == nil || ctx.Err() == nil {
		t.Fatalf("expected timeout %v", e)
	}
}
