package buildcontainer

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/core/packages"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/core/platform"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/core/project"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/infra/external"
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

	resolvedRef, pruneExistingArtifacts, err := s.resolveRef(ctx, definition, request.Ref)
	if err != nil {
		return err
	}

	names := newArtifactNames(s.Paths, request.PackageName, s.Clock.Now())
	if err := os.MkdirAll(names.StagedRoot, 0o755); err != nil {
		return fmt.Errorf("create staged root: %w", err)
	}

	keepArtifacts := false
	defer s.cleanup(context.Background(), names)
	defer func() {
		if !keepArtifacts {
			_ = os.RemoveAll(names.StagedRoot)
		}
	}()

	if err := s.Runner.Run(ctx, s.Tools.ContainerRuntime.BuildCommand(names.ImageTag, s.Paths.PackageBuildContainerDir(request.PackageName), s.Stdout, s.Stdout)); err != nil {
		return err
	}

	prepText, err := prepCommandText(definition, names, resolvedRef)
	if err != nil {
		return err
	}

	if err := s.Runner.Run(ctx, s.Tools.ContainerRuntime.PrepareRunCommand(names.PrepContainer, names.ImageTag, prepText, s.Stdout, s.Stdout)); err != nil {
		return err
	}

	if err := s.Runner.Run(ctx, s.Tools.ContainerRuntime.CommitCommand(names.PrepContainer, names.SnapshotImage, io.Discard, s.Stdout)); err != nil {
		return err
	}

	buildText, err := buildCommandText(definition, names, resolvedRef)
	if err != nil {
		return err
	}

	if err := s.Runner.Run(ctx, s.Tools.ContainerRuntime.OfflineRunCommand(names.StagedRoot+":/out:Z", names.SnapshotImage, buildText, s.Stdout, s.Stdout)); err != nil {
		return err
	}

	buildInfo, err := readBuildInfo(filepath.Join(names.StagedRoot, "build-info.txt"))
	if err != nil {
		return err
	}

	artifactRoot, err := s.finalizeArtifacts(request.PackageName, names.StagedRoot, buildInfo.Version, pruneExistingArtifacts)
	if err != nil {
		return err
	}
	keepArtifacts = true

	_, err = fmt.Fprintf(s.Stdout, "built %s version %s into %s\n", request.PackageName, buildInfo.Version, artifactRoot)
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}

func (s Service) resolveRef(ctx context.Context, definition packages.Definition, requestedRef string) (string, bool, error) {
	if requestedRef != "" {
		return requestedRef, false, nil
	}

	resolvedRef, err := s.resolveLatestStableTag(ctx, definition)
	if err != nil {
		return "", false, err
	}
	if resolvedRef == "" {
		return "", false, nil
	}

	return resolvedRef, true, nil
}

func (s Service) finalizeArtifacts(packageName string, stagedRoot string, version string, pruneExistingArtifacts bool) (string, error) {
	targetName := artifactDirectoryName(version, filepath.Base(stagedRoot))
	finalRoot := filepath.Join(s.Paths.PackageArtifactsRoot(packageName), targetName)
	if finalRoot == stagedRoot {
		return stagedRoot, nil
	}

	if err := os.RemoveAll(finalRoot); err != nil {
		return "", fmt.Errorf("remove existing artifact directory %s: %w", finalRoot, err)
	}
	if err := os.Rename(stagedRoot, finalRoot); err != nil {
		return "", fmt.Errorf("promote artifact directory to %s: %w", finalRoot, err)
	}

	if pruneExistingArtifacts {
		if err := pruneArtifactDirectories(s.Paths.PackageArtifactsRoot(packageName), targetName); err != nil {
			return "", err
		}
	}

	return finalRoot, nil
}

func pruneArtifactDirectories(packageRoot string, keepName string) error {
	entries, err := os.ReadDir(packageRoot)
	if err != nil {
		return fmt.Errorf("read artifact directory %s: %w", packageRoot, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == keepName {
			continue
		}
		if err := os.RemoveAll(filepath.Join(packageRoot, entry.Name())); err != nil {
			return fmt.Errorf("remove artifact directory %s: %w", entry.Name(), err)
		}
	}

	return nil
}

func (s Service) cleanup(ctx context.Context, names artifactNames) {
	_ = s.Runner.Run(ctx, s.Tools.ContainerRuntime.RemoveContainerCommand(names.PrepContainer, io.Discard, io.Discard))
	_ = s.Runner.Run(ctx, s.Tools.ContainerRuntime.RemoveImageCommand(names.SnapshotImage, io.Discard, io.Discard))
}
