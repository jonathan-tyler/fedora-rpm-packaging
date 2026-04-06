package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/your-github-username/fedora-package-infra/builder/src/infra/external"
)

type ToolingLoader struct{}

func (ToolingLoader) Load(path string) (external.Tooling, error) {
	config := external.DefaultToolingConfig()

	content, err := os.ReadFile(path)
	if err != nil {
		return external.Tooling{}, fmt.Errorf("read tooling config: %w", err)
	}

	if err := json.Unmarshal(content, &config); err != nil {
		return external.Tooling{}, fmt.Errorf("decode tooling config: %w", err)
	}

	tooling, err := external.NewTooling(config)
	if err != nil {
		return external.Tooling{}, fmt.Errorf("build tooling: %w", err)
	}

	return tooling, nil
}
