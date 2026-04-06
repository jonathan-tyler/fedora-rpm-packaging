package external

import (
	"io"

	"github.com/your-github-username/fedora-package-builder/src/core/platform"
)

type SystemctlTool struct {
	command string
}

func NewSystemctlTool(config ToolConfig) (SystemctlTool, error) {
	if err := validateCommandName(config.Command, "systemctl"); err != nil {
		return SystemctlTool{}, err
	}
	return SystemctlTool{command: config.Command}, nil
}

func (t SystemctlTool) UserDaemonReloadCommand(stdout io.Writer, stderr io.Writer) platform.Command {
	return platform.Command{Name: t.command, Args: []string{"--user", "daemon-reload"}, Stdout: stdout, Stderr: stderr}
}
