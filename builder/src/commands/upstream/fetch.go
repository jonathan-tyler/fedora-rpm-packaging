package upstream

import (
	"context"

	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/app/bootstrap"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/commands/dispatch"
)

func NewFetchCommand(services *bootstrap.Services) dispatch.FuncCommand {
	return dispatch.NewFuncCommand([]string{"upstream", "fetch"}, "mirror upstream repositories", func(ctx context.Context, args []string) error {
		return services.FetchUpstreams.Run(ctx, args)
	})
}
