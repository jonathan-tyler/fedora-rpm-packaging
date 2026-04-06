package bootstrap

import (
	"fmt"
	"io"

	apparchive "github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/infra/archive"
	appconfig "github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/infra/config"
	appexternal "github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/infra/external"
	appproject "github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/infra/project"
	appsystem "github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/infra/system"

	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/app/buildcontainer"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/app/buildsrpm"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/app/fetchupstreams"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/app/mockrebuild"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/app/publishrepo"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/app/syncrepo"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/app/vendorsources"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/core/packages"
	coreproject "github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/core/project"
)

type Services struct {
	Paths          coreproject.Paths
	Registry       packages.Registry
	Tooling        appexternal.Tooling
	FetchUpstreams fetchupstreams.Service
	VendorSources  vendorsources.Service
	BuildSRPM      buildsrpm.Service
	MockRebuild    mockrebuild.Service
	PublishRepo    publishrepo.Service
	SyncRepo       syncrepo.Service
	BuildContainer buildcontainer.Service
}

func New(stdout io.Writer, stderr io.Writer) (*Services, error) {
	paths, err := appproject.Locator{}.Locate()
	if err != nil {
		return nil, err
	}

	registry, err := appconfig.PackageRegistryLoader{}.Load(paths.PackagesRoot())
	if err != nil {
		return nil, fmt.Errorf("load package registry: %w", err)
	}

	tooling, err := appconfig.ToolingLoader{}.Load(paths.ToolingConfig())
	if err != nil {
		return nil, fmt.Errorf("load tooling config: %w", err)
	}

	runner := appsystem.NewCommandRunner(stdout, stderr)
	archiver := apparchive.TarGzArchiver{}
	clock := appsystem.Clock{}

	services := &Services{
		Paths:    paths,
		Registry: registry,
		Tooling:  tooling,
		FetchUpstreams: fetchupstreams.Service{
			Registry: registry,
			Paths:    paths,
			Tools:    tooling,
			Runner:   runner,
			Stdout:   stdout,
		},
		VendorSources: vendorsources.Service{
			Registry: registry,
			Paths:    paths,
			Tools:    tooling,
			Runner:   runner,
			Archiver: archiver,
			Stdout:   stdout,
		},
		BuildSRPM: buildsrpm.Service{
			Registry: registry,
			Paths:    paths,
			Tools:    tooling,
			Runner:   runner,
			Stdout:   stdout,
		},
		MockRebuild: mockrebuild.Service{
			Registry: registry,
			Paths:    paths,
			Tools:    tooling,
			Runner:   runner,
			Stdout:   stdout,
		},
		PublishRepo: publishrepo.Service{
			Registry: registry,
			Paths:    paths,
			Tools:    tooling,
			Runner:   runner,
			Stdout:   stdout,
		},
		SyncRepo: syncrepo.Service{
			Paths:  paths,
			Stdout: stdout,
		},
		BuildContainer: buildcontainer.Service{
			Registry: registry,
			Paths:    paths,
			Tools:    tooling,
			Runner:   runner,
			Clock:    clock,
			Stdout:   stdout,
		},
	}

	return services, nil
}
