package repo

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/your-github-username/fedora-package-builder/src/app/bootstrap"
	"github.com/your-github-username/fedora-package-builder/src/commands/dispatch"
)

func NewSyncServiceCommand(services *bootstrap.Services) dispatch.FuncCommand {
	return dispatch.NewFuncCommand([]string{"repo", "sync-service"}, "copy the staged repo into the live service tree", func(ctx context.Context, args []string) error {
		_ = ctx
		if len(args) != 0 {
			return fmt.Errorf("usage: fpb repo sync-service")
		}

		homeDirectory, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		stagedRepoRoot := getenvDefault("STAGED_REPO_ROOT", services.Paths.RepoRoot())
		liveRepoRoot := getenvDefault("LIVE_REPO_ROOT", filepath.Join(homeDirectory, "fedora-package-repo", "repo"))
		tempRoot := getenvDefault("LIVE_REPO_TMP_ROOT", filepath.Join(homeDirectory, "fedora-package-repo", ".repo-sync"))

		return services.SyncRepo.Run(stagedRepoRoot, liveRepoRoot, tempRoot)
	})
}
