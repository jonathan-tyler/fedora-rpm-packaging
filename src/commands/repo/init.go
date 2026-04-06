package repo

import (
	"context"
	"fmt"

	"github.com/your-github-username/fedora-package-builder/src/app/bootstrap"
	"github.com/your-github-username/fedora-package-builder/src/commands/dispatch"
)

func NewInitCommand(services *bootstrap.Services) dispatch.FuncCommand {
	return dispatch.NewFuncCommand([]string{"repo", "init"}, "initialize staged repository metadata for DNF clients", func(ctx context.Context, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("usage: fpb repo init")
		}

		releasever, basearch, gpgKey := repoEnvironment()

		return services.PublishRepo.Init(ctx, releasever, basearch, gpgKey)
	})
}