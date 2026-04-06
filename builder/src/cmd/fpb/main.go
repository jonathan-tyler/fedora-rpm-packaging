package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/app/bootstrap"
	buildcmd "github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/commands/build"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/commands/dispatch"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/commands/packagecmd"
	repocmd "github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/commands/repo"
	sourcecmd "github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/commands/source"
	upstreamcmd "github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/commands/upstream"
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
		packagecmd.NewArtifactTagCommand(services),
		upstreamcmd.NewFetchCommand(services),
		sourcecmd.NewVendorCommand(services),
		buildcmd.NewContainerBinaryCommand(services),
		buildcmd.NewSRPMCommand(services),
		buildcmd.NewMockRebuildCommand(services),
		repocmd.NewInitCommand(services),
		repocmd.NewPublishCommand(services),
		repocmd.NewSyncServiceCommand(services),
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
