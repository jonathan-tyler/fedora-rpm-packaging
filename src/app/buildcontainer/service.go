package buildcontainer

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/him/fedora-local-builder/src/core/packages"
	"github.com/him/fedora-local-builder/src/core/platform"
	"github.com/him/fedora-local-builder/src/core/project"
	"github.com/him/fedora-local-builder/src/infra/external"
)

type Service struct {
	Registry packages.Registry
	Paths    project.Paths
	Tools    external.Tooling
	Runner   platform.Runner
	Clock    platform.Clock
	Stdout   io.Writer
}

func (s Service) Run(ctx context.Context, request Request) error {
	definition, err := s.Registry.Lookup(request.PackageName)
	if err != nil {
		return err
	}
	if !definition.SupportsContainerBinaryBuild() {
		return fmt.Errorf("containerized binary build is not implemented for %s", request.PackageName)
	}

	names := newArtifactNames(s.Paths, request.PackageName, s.Clock.Now())
	if err := os.MkdirAll(names.StagedRoot, 0o755); err != nil {
		return fmt.Errorf("create staged root: %w", err)
	}

	defer s.cleanup(context.Background(), names)

	if err := s.Runner.Run(ctx, s.Tools.ContainerRuntime.BuildCommand(names.ImageTag, s.Paths.PackageBuildContainerDir(request.PackageName), s.Stdout, s.Stdout)); err != nil {
		return err
	}

	prepText, err := prepCommandText(definition, names, request.Ref, s.Tools)
	if err != nil {
		return err
	}

	if err := s.Runner.Run(ctx, s.Tools.ContainerRuntime.PrepareRunCommand(names.PrepContainer, names.ImageTag, prepText, s.Stdout, s.Stdout)); err != nil {
		return err
	}

	if err := s.Runner.Run(ctx, s.Tools.ContainerRuntime.CommitCommand(names.PrepContainer, names.SnapshotImage, io.Discard, s.Stdout)); err != nil {
		return err
	}

	buildText, err := buildCommandText(definition, names, s.Tools)
	if err != nil {
		return err
	}

	if err := s.Runner.Run(ctx, s.Tools.ContainerRuntime.OfflineRunCommand(names.StagedRoot+":/out:Z", names.SnapshotImage, buildText, s.Stdout, s.Stdout)); err != nil {
		return err
	}

	_, err = fmt.Fprintf(s.Stdout, "built %s into %s\n", request.PackageName, names.StagedRoot)
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}

func (s Service) cleanup(ctx context.Context, names artifactNames) {
	_ = s.Runner.Run(ctx, s.Tools.ContainerRuntime.RemoveContainerCommand(names.PrepContainer, io.Discard, io.Discard))
	_ = s.Runner.Run(ctx, s.Tools.ContainerRuntime.RemoveImageCommand(names.SnapshotImage, io.Discard, io.Discard))
}
