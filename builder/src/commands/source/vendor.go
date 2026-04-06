package source

import (
	"context"
	"fmt"

	"github.com/your-github-username/fedora-package-infra/builder/src/app/bootstrap"
	"github.com/your-github-username/fedora-package-infra/builder/src/commands/dispatch"
)

func NewVendorCommand(services *bootstrap.Services) dispatch.FuncCommand {
	return dispatch.NewFuncCommand([]string{"source", "vendor"}, "vendor package dependencies", func(ctx context.Context, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("usage: fpb source vendor PACKAGE /path/to/package-worktree")
		}
		return services.VendorSources.Run(ctx, args[0], args[1])
	})
}
