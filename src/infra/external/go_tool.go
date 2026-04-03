package external

import (
	"io"

	"github.com/him/fedora-local-builder/src/core/platform"
)

type GoTool struct {
	command string
}

func NewGoTool(config ToolConfig) (GoTool, error) {
	if err := validateCommandName(config.Command, "go"); err != nil {
		return GoTool{}, err
	}
	return GoTool{command: config.Command}, nil
}

func (t GoTool) ModVendorCommand(dir string, stdout io.Writer, stderr io.Writer) platform.Command {
	return platform.Command{Name: t.command, Args: []string{"mod", "vendor"}, Dir: dir, Stdout: stdout, Stderr: stderr}
}

func (t GoTool) CommandName() string {
	return t.command
}

func (t GoTool) ModVendorCommandTokens() []string {
	return []string{t.command, "mod", "vendor"}
}
