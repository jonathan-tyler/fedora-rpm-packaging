package buildsrpm

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/your-github-username/fedora-package-infra/builder/src/core/packages"
	"github.com/your-github-username/fedora-package-infra/builder/src/core/platform"
	"github.com/your-github-username/fedora-package-infra/builder/src/core/project"
	"github.com/your-github-username/fedora-package-infra/builder/src/infra/external"
)

type Service struct {
	Registry packages.Registry
	Paths    project.Paths
	Tools    external.Tooling
	Runner   platform.Runner
	Stdout   io.Writer
}

func (s Service) Run(ctx context.Context, packageName string) error {
	if _, err := s.Registry.Lookup(packageName); err != nil {
		return err
	}

	if !fileExists(s.Paths.PackageSpec(packageName)) {
		return fmt.Errorf("spec file not found: %s", s.Paths.PackageSpec(packageName))
	}

	command := s.Tools.SRPMBuilder.BuildCommand(packageName, s.Stdout, s.Stdout)
	command.Dir = s.Paths.Root

	if err := s.Runner.Run(ctx, command); err != nil {
		return err
	}

	srpms, err := filepath.Glob(filepath.Join(s.Paths.PackageResultsRoot(packageName), "*.src.rpm"))
	if err != nil {
		return fmt.Errorf("glob srpms: %w", err)
	}
	if len(srpms) == 0 {
		return fmt.Errorf("SRPM not found in %s after build", s.Paths.PackageResultsRoot(packageName))
	}

	sort.Strings(srpms)
	if _, err := fmt.Fprintf(s.Stdout, "built SRPM %s\n", srpms[len(srpms)-1]); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
