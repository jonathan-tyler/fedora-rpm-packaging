package buildcontainer

import (
	"fmt"

	"github.com/him/fedora-local-builder/src/core/packages"
	"github.com/him/fedora-local-builder/src/infra/external"
)

func prepCommandText(definition packages.Definition, names artifactNames, ref string, tools external.Tooling) (string, error) {
	switch definition.BuildSystem {
	case packages.BuildSystemGo:
		return prepGoCommandText(definition, names, ref, tools), nil
	case packages.BuildSystemCargo:
		return prepCargoCommandText(definition, names, ref, tools), nil
	default:
		return "", fmt.Errorf("unsupported build system for prep: %s", definition.BuildSystem)
	}
}

func buildCommandText(definition packages.Definition, names artifactNames, tools external.Tooling) (string, error) {
	switch definition.BuildSystem {
	case packages.BuildSystemGo:
		return buildGoCommandText(definition, names, tools), nil
	case packages.BuildSystemCargo:
		return buildCargoCommandText(definition, names, tools), nil
	default:
		return "", fmt.Errorf("unsupported build system for build: %s", definition.BuildSystem)
	}
}
