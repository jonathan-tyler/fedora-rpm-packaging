package external

import (
	"io"

	"github.com/your-github-username/fedora-package-builder/src/core/platform"
)

type RPMSignTool struct {
	command string
}

func NewRPMSignTool(config ToolConfig) (RPMSignTool, error) {
	if err := validateCommandName(config.Command, "rpmsign"); err != nil {
		return RPMSignTool{}, err
	}
	return RPMSignTool{command: config.Command}, nil
}

func (t RPMSignTool) AddSignCommand(rpms []string, stdout io.Writer, stderr io.Writer) platform.Command {
	args := append([]string{"--addsign"}, rpms...)
	return platform.Command{Name: t.command, Args: args, Stdout: stdout, Stderr: stderr}
}
