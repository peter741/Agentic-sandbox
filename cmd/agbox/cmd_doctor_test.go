// Agentic Sandbox derivative: updated repository namespace; see NOTICE.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/peter741/Agentic-sandbox/internal/platform"
	"github.com/peter741/Agentic-sandbox/internal/version"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func healthyDoctorDeps() doctorDeps {
	return doctorDeps{goos: "linux", lookPath: func(string) (string, error) { return "/usr/bin/docker", nil }, dockerOS: func(context.Context) (string, error) { return "linux", nil }, socketPath: func(platform.LookupEnv) (string, error) { return "/fake/agboxd.sock", nil }, daemonVersion: func(context.Context, string) (string, error) { return version.Version, nil }}
}
func executeDoctor(t *testing.T, deps doctorDeps, args ...string) (doctorReport, error, string) {
	t.Helper()
	cmd := newDoctorCommandWithDeps(deps)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(append([]string{"--json"}, args...))
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	err := cmd.ExecuteContext(context.Background())
	var report doctorReport
	if e := json.Unmarshal(out.Bytes(), &report); e != nil {
		t.Fatalf("invalid JSON %q: %v", out.String(), e)
	}
	return report, err, out.String()
}
func doctorCheckByName(t *testing.T, r doctorReport, name string) doctorCheck {
	t.Helper()
	for _, c := range r.Checks {
		if c.Name == name {
			return c
		}
	}
	t.Fatalf("missing check %q", name)
	return doctorCheck{}
}
func TestDoctorHealthyAndVersionWarning(t *testing.T) {
	d := healthyDoctorDeps()
	d.daemonVersion = func(context.Context, string) (string, error) { return "different-version", nil }
	r, e, _ := executeDoctor(t, d)
	if e != nil || !r.Healthy || doctorCheckByName(t, r, "daemon_version").Status != "warn" {
		t.Fatalf("unexpected result %+v %v", r, e)
	}
}
func TestDoctorReportsIndependentFailuresWithoutLeakingProbeErrors(t *testing.T) {
	d := healthyDoctorDeps()
	d.lookPath = func(string) (string, error) { return "", errors.New("missing") }
	d.dockerOS = func(context.Context) (string, error) { return "", errors.New("secret-endpoint") }
	d.daemonVersion = func(context.Context, string) (string, error) { return "", errors.New("secret-details") }
	r, e, out := executeDoctor(t, d)
	if r.Healthy || exitCodeForError(e) != 1 || shouldPrintError(e) {
		t.Fatalf("expected silent runtime failure %+v %v", r, e)
	}
	for _, n := range []string{"docker_cli", "docker_engine", "daemon"} {
		c := doctorCheckByName(t, r, n)
		if c.Status != "fail" || c.Hint == "" {
			t.Fatalf("no actionable failure %+v", c)
		}
	}
	if strings.Contains(out, "secret-") {
		t.Fatal("probe details leaked")
	}
}
func TestDoctorMissingRuntimeDoesNotPing(t *testing.T) {
	d := healthyDoctorDeps()
	d.socketPath = func(platform.LookupEnv) (string, error) { return "", errors.New("missing") }
	d.daemonVersion = func(context.Context, string) (string, error) { t.Fatal("unexpected ping"); return "", nil }
	r, e, _ := executeDoctor(t, d)
	if e == nil || doctorCheckByName(t, r, "daemon").Status != "fail" {
		t.Fatalf("unexpected result %+v %v", r, e)
	}
}
func TestDoctorRejectsUnsupportedPlatformAndEngine(t *testing.T) {
	d := healthyDoctorDeps()
	d.goos = "windows"
	d.dockerOS = func(context.Context) (string, error) { return "windows", nil }
	r, e, _ := executeDoctor(t, d)
	if e == nil || doctorCheckByName(t, r, "platform").Status != "fail" || doctorCheckByName(t, r, "docker_engine").Status != "fail" {
		t.Fatalf("unexpected result %+v %v", r, e)
	}
}
func TestDoctorEachProbeHasIndependentTimeout(t *testing.T) {
	d := healthyDoctorDeps()
	d.dockerOS = func(ctx context.Context) (string, error) { <-ctx.Done(); return "", ctx.Err() }
	d.daemonVersion = func(ctx context.Context, _ string) (string, error) {
		if ctx.Err() != nil {
			t.Fatal("expired Docker deadline reused")
		}
		if _, ok := ctx.Deadline(); !ok {
			t.Fatal("missing deadline")
		}
		return version.Version, nil
	}
	r, e, _ := executeDoctor(t, d, "--timeout", "10ms")
	if e == nil || doctorCheckByName(t, r, "daemon").Status != "pass" {
		t.Fatalf("unexpected result %+v %v", r, e)
	}
}
func TestDoctorWorkspacePolicy(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "file")
	if e := os.WriteFile(file, []byte("unchanged"), 0600); e != nil {
		t.Fatal(e)
	}
	home, e := os.UserHomeDir()
	if e != nil {
		t.Fatal(e)
	}
	link := filepath.Join(dir, "root-link")
	if e := os.Symlink("/", link); e != nil {
		t.Fatal(e)
	}
	for _, tc := range []struct {
		path string
		pass bool
	}{{dir, true}, {file, false}, {filepath.Join(dir, "missing"), false}, {"/", false}, {home, false}, {link, false}} {
		t.Run(tc.path, func(t *testing.T) {
			r, e, _ := executeDoctor(t, healthyDoctorDeps(), "--workspace", tc.path)
			if r.Healthy != tc.pass || (e == nil) != tc.pass {
				t.Fatalf("unexpected result %+v %v", r, e)
			}
		})
	}
	data, e := os.ReadFile(file)
	if e != nil || string(data) != "unchanged" {
		t.Fatal("workspace modified")
	}
}
func TestDoctorRejectsInvalidTimeoutBeforeProbing(t *testing.T) {
	for _, v := range []string{"0", "-1s"} {
		d := healthyDoctorDeps()
		d.dockerOS = func(context.Context) (string, error) { t.Fatal("unexpected probe"); return "", nil }
		c := newDoctorCommandWithDeps(d)
		c.SetArgs([]string{"--timeout", v})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})
		c.SilenceErrors = true
		c.SilenceUsage = true
		if e := c.Execute(); exitCodeForError(e) != 2 {
			t.Fatalf("unexpected result %v", e)
		}
	}
}
func TestDoctorTextAndRootHelp(t *testing.T) {
	r := collectDoctorReport(context.Background(), nil, healthyDoctorDeps(), time.Second, "")
	var out bytes.Buffer
	if e := writeDoctorReport(&out, r, false); e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(out.String(), "[PASS] daemon") || !strings.Contains(out.String(), "do not certify sandbox security") {
		t.Fatalf("unexpected report %q", out.String())
	}
	out.Reset()
	var stderr bytes.Buffer
	code := run(context.Background(), []string{"doctor", "--help"}, &out, &stderr, func(string) (string, bool) { return "", false })
	if code != 0 || stderr.Len() != 0 || !strings.Contains(out.String(), "--json") || !strings.Contains(out.String(), "--timeout") {
		t.Fatalf("unexpected help %d %q %q", code, out.String(), stderr.String())
	}
}
