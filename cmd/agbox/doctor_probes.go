// Agentic Sandbox derivative: updated repository namespace; see NOTICE.
package main

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/peter741/Agentic-sandbox/internal/platform"
	"github.com/peter741/Agentic-sandbox/sdk/go/rawclient"
	"os/exec"
	"runtime"
)

type doctorDeps struct {
	goos          string
	lookPath      func(string) (string, error)
	dockerOS      func(context.Context) (string, error)
	socketPath    func(platform.LookupEnv) (string, error)
	daemonVersion func(context.Context, string) (string, error)
}

func defaultDoctorDeps() doctorDeps {
	return doctorDeps{runtime.GOOS, exec.LookPath, probeDockerOS, platform.SocketPath, probeDaemonVersion}
}

func probeDockerOS(ctx context.Context) (string, error) {
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer docker.Close()
	info, err := docker.Info(ctx)
	if err != nil {
		return "", err
	}
	return info.OSType, nil
}

func probeDaemonVersion(ctx context.Context, socketPath string) (string, error) {
	cli, err := rawclient.New(socketPath, rawclient.WithTimeout(0))
	if err != nil {
		return "", err
	}
	defer cli.Close()
	response, err := cli.Ping(ctx)
	if err != nil {
		return "", err
	}
	return response.GetVersion(), nil
}
