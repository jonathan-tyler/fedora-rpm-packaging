package buildcontainer

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/him/fedora-local-builder/src/core/packages"
	"github.com/him/fedora-local-builder/src/core/platform"
	"github.com/him/fedora-local-builder/src/core/project"
)

type Service struct {
	Registry packages.Registry
	Paths    project.Paths
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

	if err := s.Runner.Run(ctx, platform.Command{
		Name:   "podman",
		Args:   []string{"build", "-t", names.ImageTag, s.Paths.PackageBuildContainerDir(request.PackageName)},
		Stdout: s.Stdout,
		Stderr: s.Stdout,
	}); err != nil {
		return err
	}

	if err := s.Runner.Run(ctx, platform.Command{
		Name:   "podman",
		Args:   []string{"run", "--name", names.PrepContainer, names.ImageTag, "bash", "-lc", prepCommandText(definition, names, request.Ref)},
		Stdout: s.Stdout,
		Stderr: s.Stdout,
	}); err != nil {
		return err
	}

	if err := s.Runner.Run(ctx, platform.Command{
		Name:   "podman",
		Args:   []string{"commit", names.PrepContainer, names.SnapshotImage},
		Stdout: io.Discard,
		Stderr: s.Stdout,
	}); err != nil {
		return err
	}

	if err := s.Runner.Run(ctx, platform.Command{
		Name: "podman",
		Args: []string{
			"run", "--rm", "--network", "none",
			"-v", names.StagedRoot + ":/out:Z",
			names.SnapshotImage,
			"bash", "-lc", buildCommandText(definition, names),
		},
		Stdout: s.Stdout,
		Stderr: s.Stdout,
	}); err != nil {
		return err
	}

	_, err = fmt.Fprintf(s.Stdout, "built %s into %s\n", request.PackageName, names.StagedRoot)
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}

func (s Service) cleanup(ctx context.Context, names artifactNames) {
	_ = s.Runner.Run(ctx, platform.Command{
		Name:   "podman",
		Args:   []string{"rm", "-f", names.PrepContainer},
		Stdout: io.Discard,
		Stderr: io.Discard,
	})
	_ = s.Runner.Run(ctx, platform.Command{
		Name:   "podman",
		Args:   []string{"image", "rm", "-f", names.SnapshotImage},
		Stdout: io.Discard,
		Stderr: io.Discard,
	})
}
