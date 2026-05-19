// Package dms is the single choke point between dms-ui and a running
// docker-mailserver container. Every mutation goes through the container's
// own `setup` CLI (which keeps Postfix/Dovecot consistent and triggers DMS's
// graceful change-detector reload); reads use setup/doveadm/supervisorctl.
//
// We shell out to the `docker` CLI rather than the Docker Go SDK: it keeps
// the dependency tree tiny, and the CLI transparently supports
// DOCKER_HOST=ssh://host (used for dev/verification against the live host).
package dms

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Client targets one mail container via the local docker CLI.
type Client struct {
	dockerBin string
	container string
}

// New returns a client for the named container. The docker CLI must be on
// PATH; DOCKER_HOST in the environment is honoured by the CLI.
func New(containerName string) (*Client, error) {
	bin, err := exec.LookPath("docker")
	if err != nil {
		return nil, fmt.Errorf("docker CLI not found on PATH: %w", err)
	}
	return &Client{dockerBin: bin, container: containerName}, nil
}

// Result is the captured outcome of one exec.
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;?]*[a-zA-Z]`)

func stripANSI(s string) string { return ansiRe.ReplaceAllString(s, "") }

// docker runs `docker <args>` with a timeout and returns captured output.
func (c *Client) docker(ctx context.Context, args ...string) (*Result, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, c.dockerBin, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	res := &Result{
		Stdout: stripANSI(stdout.String()),
		Stderr: stripANSI(stderr.String()),
	}
	if cmd.ProcessState != nil {
		res.ExitCode = cmd.ProcessState.ExitCode()
	}
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			// Failed to even launch / context error.
			return res, fmt.Errorf("docker %s: %w", strings.Join(args, " "), err)
		}
	}
	return res, nil
}

// Exec runs argv inside the mail container without a TTY (so output carries
// no interactive prompts or colour). argv is passed verbatim, no shell, so
// no quoting/injection concerns.
func (c *Client) Exec(ctx context.Context, argv ...string) (*Result, error) {
	return c.docker(ctx, append([]string{"exec", c.container}, argv...)...)
}

// run is the helper for commands where a non-zero exit is an error.
func (c *Client) run(ctx context.Context, argv ...string) (string, error) {
	r, err := c.Exec(ctx, argv...)
	if err != nil {
		return "", err
	}
	if r.ExitCode != 0 {
		msg := strings.TrimSpace(r.Stderr)
		if msg == "" {
			msg = strings.TrimSpace(r.Stdout)
		}
		return "", fmt.Errorf("%s exited %d: %s", strings.Join(argv, " "), r.ExitCode, msg)
	}
	return r.Stdout, nil
}

// setup runs the DMS `setup` CLI.
func (c *Client) setup(ctx context.Context, args ...string) (string, error) {
	return c.run(ctx, append([]string{"setup"}, args...)...)
}

// ContainerName is the configured mail container (shown in the UI).
func (c *Client) ContainerName() string { return c.container }

// Ping verifies docker connectivity and that the container is running.
func (c *Client) Ping(ctx context.Context) error {
	r, err := c.docker(ctx, "inspect", "-f", "{{.State.Running}}", c.container)
	if err != nil {
		return err
	}
	if r.ExitCode != 0 {
		return fmt.Errorf("container %s not found: %s", c.container, strings.TrimSpace(r.Stderr))
	}
	if strings.TrimSpace(r.Stdout) != "true" {
		return fmt.Errorf("container %s is not running", c.container)
	}
	return nil
}
