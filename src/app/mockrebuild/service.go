package mockrebuild

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/him/fedora-local-builder/src/core/packages"
	"github.com/him/fedora-local-builder/src/core/platform"
	"github.com/him/fedora-local-builder/src/core/project"
	"github.com/him/fedora-local-builder/src/infra/external"
)

type Service struct {
	Registry packages.Registry
	Paths    project.Paths
	Tools    external.Tooling
	Runner   platform.Runner
	Stdout   io.Writer
}

func (s Service) Run(ctx context.Context, packageName string, srpmPath string) error {
	if _, err := s.Registry.Lookup(packageName); err != nil {
		return err
	}

	if !fileExists(srpmPath) {
		return fmt.Errorf("SRPM not found: %s", srpmPath)
	}

	resultDir := s.Paths.PackageResultsRoot(packageName)
	if err := os.MkdirAll(resultDir, 0o755); err != nil {
		return fmt.Errorf("create result directory: %w", err)
	}

	command := s.Tools.Mock.RebuildCommand(s.Paths.MockConfigDir(), resultDir, srpmPath, "fedora-42-x86_64", s.Stdout, s.Stdout)
	command.Env = map[string]string{
		"FEDORA_PACKAGE_BUILDER_ROOT": s.Paths.Root,
	}

	return s.Runner.Run(ctx, command)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
