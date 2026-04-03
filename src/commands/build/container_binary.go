package build

import (
	"context"
	"fmt"

	"github.com/him/fedora-local-builder/src/app/bootstrap"
	"github.com/him/fedora-local-builder/src/app/buildcontainer"
	"github.com/him/fedora-local-builder/src/commands/dispatch"
)

func NewContainerBinaryCommand(services *bootstrap.Services) dispatch.FuncCommand {
	return dispatch.NewFuncCommand([]string{"build", "container-binary"}, "build a raw binary in an offline container", func(ctx context.Context, args []string) error {
		if len(args) < 1 || len(args) > 2 {
			return fmt.Errorf("usage: fpb build container-binary PACKAGE [REF]")
		}

		request := buildcontainer.Request{PackageName: args[0]}
		if len(args) == 2 {
			request.Ref = args[1]
		}

		return services.BuildContainer.Run(ctx, request)
	})
}
