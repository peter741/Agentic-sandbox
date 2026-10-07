// Agentic Sandbox derivative: updated repository namespace; see NOTICE.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/peter741/Agentic-sandbox/internal/platform"
	"github.com/peter741/Agentic-sandbox/internal/version"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
	"time"
)

type doctorCheck struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Hint    string `json:"hint,omitempty"`
}
type doctorReport struct {
	Version string        `json:"version"`
	Healthy bool          `json:"healthy"`
	Checks  []doctorCheck `json:"checks"`
}

func newDoctorCommand() *cobra.Command { return newDoctorCommandWithDeps(defaultDoctorDeps()) }
func newDoctorCommandWithDeps(deps doctorDeps) *cobra.Command {
	var asJSON bool
	var timeout time.Duration
	var workspace string
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Diagnose local Docker and daemon setup without changing it",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if timeout <= 0 {
				return usageErrorf("--timeout must be greater than zero")
			}
			report := collectDoctorReport(cmd.Context(), lookupEnvFromCmd(cmd), deps, timeout, workspace)
			if err := writeDoctorReport(cmd.OutOrStdout(), report, asJSON); err != nil {
				return runtimeErrorf("write doctor report: %v", err)
			}
			if !report.Healthy {
				return exitCodeError(exitCodeRuntimeError)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Print a machine-readable report")
	cmd.Flags().DurationVar(&timeout, "timeout", 5*time.Second, "Timeout for each Docker or daemon probe")
	cmd.Flags().StringVar(&workspace, "workspace", "", "Also check a project directory without copying it")
	return cmd
}
func collectDoctorReport(ctx context.Context, lookupEnv platform.LookupEnv, deps doctorDeps, timeout time.Duration, workspace string) doctorReport {
	report := doctorReport{Version: version.Version, Healthy: true, Checks: []doctorCheck{}}
	add := func(name, status, message, hint string) {
		report.Checks = append(report.Checks, doctorCheck{name, status, message, hint})
		if status == "fail" {
			report.Healthy = false
		}
	}
	if deps.goos == "linux" || deps.goos == "darwin" {
		add("platform", "pass", deps.goos, "")
	} else {
		add("platform", "fail", "Unsupported host platform: "+deps.goos, "Use Linux or macOS; Windows users can run inside a configured Linux environment.")
	}
	if _, err := deps.lookPath("docker"); err != nil {
		add("docker_cli", "fail", "Docker CLI is not on PATH", "Install Docker and open a new terminal.")
	} else {
		add("docker_cli", "pass", "Docker CLI found", "")
	}
	dockerCtx, cancelDocker := context.WithTimeout(ctx, timeout)
	engineOS, err := deps.dockerOS(dockerCtx)
	deadline := dockerCtx.Err()
	cancelDocker()
	if err != nil {
		message := "Docker Engine could not be queried"
		if deadline != nil {
			message = "Docker Engine probe timed out or was cancelled"
		}
		add("docker_engine", "fail", message, "Start Docker and check daemon access, DOCKER_HOST, and Docker TLS settings.")
	} else if engineOS != "linux" {
		add("docker_engine", "fail", "Docker Engine is not running Linux containers", "Switch Docker to Linux containers.")
	} else {
		add("docker_engine", "pass", "Docker Engine reports Linux containers", "")
	}
	socketPath, err := deps.socketPath(lookupEnv)
	if err != nil {
		add("daemon", "fail", "Cannot resolve the daemon socket", "On Linux, use a login session with XDG_RUNTIME_DIR configured; see docs/doctor.md.")
	} else {
		probeCtx, cancel := context.WithTimeout(ctx, timeout)
		daemonVersion, pingErr := deps.daemonVersion(probeCtx, socketPath)
		cancel()
		if pingErr != nil {
			add("daemon", "fail", "agboxd did not answer its health probe", "Start agboxd and inspect its service logs; see docs/doctor.md.")
		} else {
			add("daemon", "pass", "agboxd responded", "")
			if daemonVersion == "" || version.Version == "dev" || daemonVersion == "dev" {
				add("daemon_version", "warn", "Version compatibility is not verified for development or unspecified versions", "Build and install CLI and daemon from the same checkout.")
			} else if daemonVersion != version.Version {
				add("daemon_version", "warn", "CLI and daemon versions differ", "Install CLI and daemon from the same release.")
			} else {
				add("daemon_version", "pass", "CLI and daemon versions match", "")
			}
		}
	}
	if workspace != "" {
		resolved, err := validateWorkspacePath(workspace)
		if err != nil {
			add("workspace", "fail", "Project path is unavailable or rejected by the workspace policy", "Choose an existing project directory, not your home directory or filesystem root.")
		} else if info, err := os.Stat(resolved); err != nil || !info.IsDir() {
			add("workspace", "fail", "Project path is not an accessible directory", "Pass a directory to --workspace.")
		} else {
			add("workspace", "pass", "Project directory exists and passes workspace path policy", "")
		}
	}
	return report
}
func writeDoctorReport(w io.Writer, report doctorReport, asJSON bool) error {
	if asJSON {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(report)
	}
	if _, err := fmt.Fprintf(w, "Agentic Sandbox diagnostics (agbox %s)\n", report.Version); err != nil {
		return err
	}
	for _, check := range report.Checks {
		if _, err := fmt.Fprintf(w, "[%s] %s: %s\n", strings.ToUpper(check.Status), check.Name, check.Message); err != nil {
			return err
		}
		if check.Hint != "" {
			if _, err := fmt.Fprintf(w, "  %s\n", check.Hint); err != nil {
				return err
			}
		}
	}
	_, err := fmt.Fprintln(w, "These connectivity checks do not certify sandbox security or agent authentication.")
	return err
}
