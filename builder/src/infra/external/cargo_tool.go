package external

import (
	"io"

	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/core/platform"
)

type CargoTool struct {
	command string
}

func NewCargoTool(config ToolConfig) (CargoTool, error) {
	if err := validateCommandName(config.Command, "cargo"); err != nil {
		return CargoTool{}, err
	}
	return CargoTool{command: config.Command}, nil
}

func (t CargoTool) VendorCommand(dir string, vendorDirectory string, stderr io.Writer) platform.Command {
	return platform.Command{Name: t.command, Args: []string{"vendor", vendorDirectory}, Dir: dir, Stderr: stderr}
}

func (t CargoTool) CommandName() string {
	return t.command
}

func (t CargoTool) VendorCommandTokens(vendorDirectory string) []string {
	return []string{t.command, "vendor", vendorDirectory}
}

func (t CargoTool) BuildReleaseOfflineCommandTokens(outputBinaries []string) []string {
	args := []string{t.command, "build", "--locked", "--offline", "--release"}
	for _, outputBinary := range outputBinaries {
		args = append(args, "--bin", outputBinary)
	}
	return args
}
