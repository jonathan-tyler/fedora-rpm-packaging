package external

import (
	"io"

	"github.com/your-github-username/fedora-package-builder/src/core/platform"
)

type MockTool struct {
	command string
}

func NewMockTool(config ToolConfig) (MockTool, error) {
	if err := validateCommandName(config.Command, "mock"); err != nil {
		return MockTool{}, err
	}
	return MockTool{command: config.Command}, nil
}

func (t MockTool) RebuildCommand(configDir string, resultDir string, srpmPath string, rootPath string, stdout io.Writer, stderr io.Writer) platform.Command {
	return platform.Command{
		Name:   t.command,
		Args:   []string{"--configdir", configDir, "--root", rootPath, "--resultdir", resultDir, "--rebuild", srpmPath},
		Stdout: stdout,
		Stderr: stderr,
	}
}
