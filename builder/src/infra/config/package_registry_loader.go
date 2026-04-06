package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/core/packages"
)

type PackageRegistryLoader struct{}

func (PackageRegistryLoader) Load(root string) (packages.Registry, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read package directory: %w", err)
	}

	definitions := make([]packages.Definition, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		definition, err := loadDefinition(filepath.Join(root, entry.Name(), "manifest.json"), entry.Name())
		if err != nil {
			return nil, err
		}
		definitions = append(definitions, definition)
	}

	if len(definitions) == 0 {
		return nil, fmt.Errorf("no package manifests found in %s", root)
	}

	registry, err := packages.NewStaticRegistry(definitions)
	if err != nil {
		return nil, fmt.Errorf("build package registry: %w", err)
	}

	return registry, nil
}

func loadDefinition(path string, directoryName string) (packages.Definition, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return packages.Definition{}, fmt.Errorf("read package manifest %s: %w", directoryName, err)
	}

	var definition packages.Definition
	if err := json.Unmarshal(content, &definition); err != nil {
		return packages.Definition{}, fmt.Errorf("decode package manifest %s: %w", directoryName, err)
	}

	if definition.Name != "" && definition.Name != directoryName {
		return packages.Definition{}, fmt.Errorf("package manifest %s has mismatched name %q", path, definition.Name)
	}

	definition.Name = directoryName
	return definition, nil
}
