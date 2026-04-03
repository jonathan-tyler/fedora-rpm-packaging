package main

import (
	"context"
	"fmt"
	"os"

	"github.com/him/fedora-local-builder/src/app/bootstrap"
	buildcmd "github.com/him/fedora-local-builder/src/commands/build"
	"github.com/him/fedora-local-builder/src/commands/dispatch"
	"github.com/him/fedora-local-builder/src/commands/packagecmd"
	repocmd "github.com/him/fedora-local-builder/src/commands/repo"
	servicecmd "github.com/him/fedora-local-builder/src/commands/service"
	sourcecmd "github.com/him/fedora-local-builder/src/commands/source"
	upstreamcmd "github.com/him/fedora-local-builder/src/commands/upstream"
)

var version = "dev"

func main() {
	services, err := bootstrap.New(os.Stdout, os.Stderr)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	dispatcher := dispatch.New(
		packagecmd.NewListCommand(services, os.Stdout),
		upstreamcmd.NewFetchCommand(services),
		sourcecmd.NewVendorCommand(services),
		buildcmd.NewContainerBinaryCommand(services),
		buildcmd.NewMockRebuildCommand(services),
		repocmd.NewPublishCommand(services),
		repocmd.NewSyncServiceCommand(services),
		servicecmd.NewInstallQuadletCommand(services),
	)

	if len(os.Args) == 1 {
		_, _ = fmt.Fprint(os.Stderr, dispatcher.Usage())
		os.Exit(1)
	}

	if err := dispatcher.Execute(context.Background(), os.Args[1:]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		if err.Error() == dispatcher.Usage() {
			os.Exit(0)
		}
		os.Exit(1)
	}
}
