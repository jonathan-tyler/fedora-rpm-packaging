package packagecmd

import (
	"context"
	"fmt"

	"github.com/your-github-username/fedora-package-infra/builder/src/app/bootstrap"
	"github.com/your-github-username/fedora-package-infra/builder/src/commands/dispatch"
)

func NewArtifactTagCommand(services *bootstrap.Services) dispatch.FuncCommand {
	return dispatch.NewFuncCommand([]string{"package", "artifact-tag"}, "resolve the latest stable upstream tag and normalized artifact tag", func(ctx context.Context, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("usage: fpb package artifact-tag PACKAGE")
		}

		resolved, err := services.BuildContainer.ResolveArtifactTag(ctx, args[0])
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(services.BuildContainer.Stdout, "package=%s\nupstream_tag=%s\nartifact_tag=%s\n", resolved.PackageName, resolved.UpstreamTag, resolved.ArtifactTag)
		return err
	})
}
