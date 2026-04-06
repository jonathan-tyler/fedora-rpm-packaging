package external

import (
	"io"

	"github.com/your-github-username/fedora-package-infra/builder/src/core/platform"
)

type GPGTool struct {
	command string
}

func NewGPGTool(config ToolConfig) (GPGTool, error) {
	if err := validateCommandName(config.Command, "gpg"); err != nil {
		return GPGTool{}, err
	}
	return GPGTool{command: config.Command}, nil
}

func (t GPGTool) ExportArmoredPublicKeyCommand(gpgKey string, stderr io.Writer) platform.Command {
	return platform.Command{Name: t.command, Args: []string{"--batch", "--yes", "--armor", "--export", gpgKey}, Stderr: stderr}
}

func (t GPGTool) DetachSignCommand(gpgKey string, outputPath string, inputPath string, stdout io.Writer, stderr io.Writer) platform.Command {
	return platform.Command{
		Name:   t.command,
		Args:   []string{"--batch", "--yes", "--armor", "--detach-sign", "--local-user", gpgKey, "--output", outputPath, inputPath},
		Stdout: stdout,
		Stderr: stderr,
	}
}
