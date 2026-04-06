package build

import (
	"context"
	"fmt"

	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/app/bootstrap"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/commands/dispatch"
)

func NewSRPMCommand(services *bootstrap.Services) dispatch.FuncCommand {
	return dispatch.NewFuncCommand([]string{"build", "srpm"}, "build an SRPM in a one-shot helper container", func(ctx context.Context, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("usage: fpb build srpm PACKAGE")
		}

		return services.BuildSRPM.Run(ctx, args[0])
	})
}
