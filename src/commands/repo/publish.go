package repo

import (
	"context"
	"fmt"
	"os"

	"github.com/him/fedora-local-builder/src/app/bootstrap"
	"github.com/him/fedora-local-builder/src/commands/dispatch"
)

func NewPublishCommand(services *bootstrap.Services) dispatch.FuncCommand {
	return dispatch.NewFuncCommand([]string{"repo", "publish"}, "publish built RPMs into the staged repository", func(ctx context.Context, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("usage: fpb repo publish PACKAGE")
		}

		releasever := getenvDefault("FEDORA_PACKAGE_RELEASEVER", "42")
		basearch := getenvDefault("FEDORA_PACKAGE_BASEARCH", "x86_64")
		gpgKey := os.Getenv("FEDORA_PACKAGE_GPG_KEY")

		return services.PublishRepo.Run(ctx, args[0], releasever, basearch, gpgKey)
	})
}

func getenvDefault(name string, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
