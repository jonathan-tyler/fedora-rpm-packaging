package system

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/him/fedora-local-builder/src/core/platform"
)

type CommandRunner struct {
	stdout io.Writer
	stderr io.Writer
}

func NewCommandRunner(stdout io.Writer, stderr io.Writer) *CommandRunner {
	return &CommandRunner{stdout: stdout, stderr: stderr}
}

func (r *CommandRunner) Run(ctx context.Context, command platform.Command) error {
	cmd := exec.CommandContext(ctx, command.Name, command.Args...)
	cmd.Dir = command.Dir
	cmd.Env = mergeEnvironment(command.Env)
	cmd.Stdout = fallbackWriter(command.Stdout, r.stdout)
	cmd.Stderr = fallbackWriter(command.Stderr, r.stderr)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run %s: %w", command.Name, err)
	}

	return nil
}

func (r *CommandRunner) Capture(ctx context.Context, command platform.Command) (string, error) {
	cmd := exec.CommandContext(ctx, command.Name, command.Args...)
	cmd.Dir = command.Dir
	cmd.Env = mergeEnvironment(command.Env)
	cmd.Stderr = fallbackWriter(command.Stderr, r.stderr)

	var buffer bytes.Buffer
	cmd.Stdout = &buffer

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("run %s: %w", command.Name, err)
	}

	return buffer.String(), nil
}

func fallbackWriter(primary io.Writer, fallback io.Writer) io.Writer {
	if primary != nil {
		return primary
	}
	if fallback != nil {
		return fallback
	}
	return io.Discard
}

func mergeEnvironment(extra map[string]string) []string {
	if len(extra) == 0 {
		return os.Environ()
	}

	environment := append([]string{}, os.Environ()...)
	for key, value := range extra {
		environment = append(environment, key+"="+value)
	}
	return environment
}
