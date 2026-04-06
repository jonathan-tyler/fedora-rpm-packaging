package external

import (
	"io"

	"github.com/your-github-username/fedora-package-infra/builder/src/core/platform"
)

type SRPMBuilderTool struct {
	command string
}

func NewSRPMBuilderTool(config ToolConfig) (SRPMBuilderTool, error) {
	if err := validateCommandName(config.Command, "srpm_builder"); err != nil {
		return SRPMBuilderTool{}, err
	}

	return SRPMBuilderTool{command: config.Command}, nil
}

func (t SRPMBuilderTool) BuildCommand(packageName string, stdout io.Writer, stderr io.Writer) platform.Command {
	return platform.Command{
		Name:   t.command,
		Args:   []string{packageName},
		Stdout: stdout,
		Stderr: stderr,
	}
}
