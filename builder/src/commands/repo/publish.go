package repo

import (
	"context"
	"fmt"

	"github.com/your-github-username/fedora-package-builder/src/app/bootstrap"
	"github.com/your-github-username/fedora-package-builder/src/commands/dispatch"
)

func NewPublishCommand(services *bootstrap.Services) dispatch.FuncCommand {
	return dispatch.NewFuncCommand([]string{"repo", "publish"}, "publish built RPMs into the staged repository", func(ctx context.Context, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("usage: fpb repo publish PACKAGE")
		}

		releasever, basearch, gpgKey := repoEnvironment()

		return services.PublishRepo.Run(ctx, args[0], releasever, basearch, gpgKey)
	})
}
