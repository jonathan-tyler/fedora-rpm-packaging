package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/him/fedora-local-builder/src/app/bootstrap"
	"github.com/him/fedora-local-builder/src/commands/dispatch"
)

func NewInstallQuadletCommand(services *bootstrap.Services) dispatch.FuncCommand {
	return dispatch.NewFuncCommand([]string{"service", "install-quadlet"}, "install quadlet units for the repo service", func(ctx context.Context, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("usage: fpb service install-quadlet")
		}

		targetHome := os.Getenv("QUADLET_TARGET_HOME")
		if targetHome == "" {
			homeDirectory, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			targetHome = homeDirectory
		}

		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		repoServerRoot := os.Getenv("REPO_SERVER_ROOT")
		if repoServerRoot == "" {
			repoServerRoot = filepath.Join(targetHome, "fedora-package-repo", "repo")
		}

		return services.InstallQuadlet.Run(ctx, targetHome, xdgConfigHome, repoServerRoot)
	})
}
