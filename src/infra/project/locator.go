package project

import (
	"fmt"
	"os"
	"path/filepath"

	coreproject "github.com/him/fedora-local-builder/src/core/project"
)

const rootOverrideEnv = "FEDORA_PACKAGE_BUILDER_ROOT"

type Locator struct{}

func (Locator) Locate() (coreproject.Paths, error) {
	if override := os.Getenv(rootOverrideEnv); override != "" {
		absoluteRoot, err := filepath.Abs(override)
		if err != nil {
			return coreproject.Paths{}, fmt.Errorf("resolve %s: %w", rootOverrideEnv, err)
		}
		paths := coreproject.Paths{Root: absoluteRoot}
		if !hasManifest(paths.PackageManifest()) {
			return coreproject.Paths{}, fmt.Errorf("%s does not point at a fedora-local-builder repo", rootOverrideEnv)
		}
		return paths, nil
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		return coreproject.Paths{}, fmt.Errorf("get working directory: %w", err)
	}

	current, err := filepath.Abs(workingDirectory)
	if err != nil {
		return coreproject.Paths{}, fmt.Errorf("resolve working directory: %w", err)
	}

	for {
		paths := coreproject.Paths{Root: current}
		if hasManifest(paths.PackageManifest()) {
			return paths, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return coreproject.Paths{}, fmt.Errorf("could not locate fedora-local-builder project root")
}

func hasManifest(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
