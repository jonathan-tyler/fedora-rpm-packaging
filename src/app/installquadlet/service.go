package installquadlet

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/him/fedora-local-builder/src/core/platform"
	"github.com/him/fedora-local-builder/src/core/project"
	"github.com/him/fedora-local-builder/src/infra/external"
)

type Service struct {
	Paths  project.Paths
	Tools  external.Tooling
	Runner platform.Runner
	Stdout io.Writer
}

func (s Service) Run(ctx context.Context, targetHome string, xdgConfigHome string, repoServerRoot string) error {
	unitDir := filepath.Join(targetHome, ".config", "containers", "systemd")
	if xdgConfigHome != "" {
		unitDir = filepath.Join(xdgConfigHome, "containers", "systemd")
	}

	if err := os.MkdirAll(unitDir, 0o755); err != nil {
		return fmt.Errorf("create unit directory: %w", err)
	}
	if err := os.MkdirAll(repoServerRoot, 0o755); err != nil {
		return fmt.Errorf("create repo server root: %w", err)
	}

	if err := s.Runner.Run(ctx, s.Tools.ContainerRuntime.BuildCommand("localhost/fedora-package-repo-nginx:latest", s.Paths.NginxContainerDir(), s.Stdout, s.Stdout)); err != nil {
		return err
	}

	if err := replaceSymlink(s.Paths.QuadletNetworkUnit(), filepath.Join(unitDir, "fedora-package-repo.network")); err != nil {
		return err
	}
	if err := replaceSymlink(s.Paths.QuadletContainerUnit(), filepath.Join(unitDir, "fedora-package-repo.container")); err != nil {
		return err
	}

	if err := s.Runner.Run(ctx, s.Tools.Systemctl.UserDaemonReloadCommand(s.Stdout, s.Stdout)); err != nil {
		return err
	}

	_, err := fmt.Fprint(s.Stdout, "Installed Quadlet files.\n\nNext step:\n  systemctl --user enable --now fedora-package-repo.service\n")
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}

func replaceSymlink(sourcePath string, targetPath string) error {
	if err := os.RemoveAll(targetPath); err != nil {
		return fmt.Errorf("replace symlink %s: %w", targetPath, err)
	}
	if err := os.Symlink(sourcePath, targetPath); err != nil {
		return fmt.Errorf("create symlink %s: %w", targetPath, err)
	}
	return nil
}
