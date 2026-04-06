package build

import (
	"context"
	"fmt"

	"github.com/your-github-username/fedora-package-infra/builder/src/app/bootstrap"
	"github.com/your-github-username/fedora-package-infra/builder/src/commands/dispatch"
)

func NewMockRebuildCommand(services *bootstrap.Services) dispatch.FuncCommand {
	return dispatch.NewFuncCommand([]string{"build", "mock-rebuild"}, "rebuild an SRPM with mock", func(ctx context.Context, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("usage: fpb build mock-rebuild PACKAGE path/to/package.src.rpm")
		}
		return services.MockRebuild.Run(ctx, args[0], args[1])
	})
}
