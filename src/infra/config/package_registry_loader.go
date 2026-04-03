package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/him/fedora-local-builder/src/core/packages"
)

type packageManifest struct {
	Packages []packages.Definition `json:"packages"`
}

type PackageRegistryLoader struct{}

func (PackageRegistryLoader) Load(path string) (packages.Registry, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read package manifest: %w", err)
	}

	var manifest packageManifest
	if err := json.Unmarshal(content, &manifest); err != nil {
		return nil, fmt.Errorf("decode package manifest: %w", err)
	}

	registry, err := packages.NewStaticRegistry(manifest.Packages)
	if err != nil {
		return nil, fmt.Errorf("build package registry: %w", err)
	}

	return registry, nil
}
