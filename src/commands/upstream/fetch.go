package upstream

import (
	"context"

	"github.com/him/fedora-local-builder/src/app/bootstrap"
	"github.com/him/fedora-local-builder/src/commands/dispatch"
)

func NewFetchCommand(services *bootstrap.Services) dispatch.FuncCommand {
	return dispatch.NewFuncCommand([]string{"upstream", "fetch"}, "mirror upstream repositories", func(ctx context.Context, args []string) error {
		return services.FetchUpstreams.Run(ctx, args)
	})
}
