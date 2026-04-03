package fetchupstreams

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/him/fedora-local-builder/src/core/packages"
	"github.com/him/fedora-local-builder/src/core/platform"
	"github.com/him/fedora-local-builder/src/core/project"
)

type Service struct {
	Registry packages.Registry
	Paths    project.Paths
	Runner   platform.Runner
	Stdout   io.Writer
}

func (s Service) Run(ctx context.Context, packageNames []string) error {
	definitions, err := s.resolvePackages(packageNames)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(s.Paths.UpstreamRoot(), 0o755); err != nil {
		return fmt.Errorf("create upstream root: %w", err)
	}

	for _, definition := range definitions {
		target := filepath.Join(s.Paths.UpstreamRoot(), definition.Name+".git")
		if directoryExists(target) {
			err = s.Runner.Run(ctx, platform.Command{
				Name:   "git",
				Args:   []string{"-C", target, "fetch", "--prune", "--tags", "origin"},
				Stdout: s.Stdout,
				Stderr: s.Stdout,
			})
		} else {
			err = s.Runner.Run(ctx, platform.Command{
				Name:   "git",
				Args:   []string{"clone", "--mirror", definition.RepoURL, target},
				Stdout: s.Stdout,
				Stderr: s.Stdout,
			})
		}
		if err != nil {
			return err
		}

		if _, err := fmt.Fprintf(s.Stdout, "mirrored %s from %s\n", definition.Name, definition.RepoURL); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
	}

	return nil
}

func (s Service) resolvePackages(packageNames []string) ([]packages.Definition, error) {
	if len(packageNames) == 0 {
		return s.Registry.List(), nil
	}

	definitions := make([]packages.Definition, 0, len(packageNames))
	for _, packageName := range packageNames {
		definition, err := s.Registry.Lookup(packageName)
		if err != nil {
			return nil, err
		}
		definitions = append(definitions, definition)
	}

	return definitions, nil
}

func directoryExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
