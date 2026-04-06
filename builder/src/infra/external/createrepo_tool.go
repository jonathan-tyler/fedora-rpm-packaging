package external

import (
	"io"

	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/core/platform"
)

type CreateRepoTool struct {
	command string
}

func NewCreateRepoTool(config ToolConfig) (CreateRepoTool, error) {
	if err := validateCommandName(config.Command, "createrepo"); err != nil {
		return CreateRepoTool{}, err
	}
	return CreateRepoTool{command: config.Command}, nil
}

func (t CreateRepoTool) UpdateCommand(repoDir string, stdout io.Writer, stderr io.Writer) platform.Command {
	return platform.Command{Name: t.command, Args: []string{"--update", repoDir}, Stdout: stdout, Stderr: stderr}
}
