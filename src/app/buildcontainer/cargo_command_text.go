package buildcontainer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/him/fedora-local-builder/src/core/commandline"
	"github.com/him/fedora-local-builder/src/core/packages"
	"github.com/him/fedora-local-builder/src/infra/external"
)

func prepCargoCommandText(definition packages.Definition, names artifactNames, ref string, tools external.Tooling) string {
	return strings.Join([]string{
		cloneCommandText(definition, names, ref, tools),
		"cd " + commandline.Quote(names.CloneDir),
		"mkdir -p .cargo",
		strings.Join(quoteArguments(tools.Cargo.VendorCommandTokens("vendor")), " ") + " > .cargo/config.toml",
	}, " && ")
}

func buildCargoCommandText(definition packages.Definition, names artifactNames, tools external.Tooling) string {
	lines := []string{
		"set -euo pipefail",
		"cd " + commandline.Quote(names.CloneDir),
		`version="$(git describe --tags --always 2>/dev/null || git rev-parse --short HEAD)"`,
		`commit="$(git rev-parse HEAD)"`,
		cargoBuildCommandText(definition, tools),
	}

	for _, outputBinary := range definition.OutputBinaries {
		lines = append(lines, fmt.Sprintf("cp %s %s", commandline.Quote(filepath.Join("target", "release", outputBinary)), commandline.Quote(filepath.Join("/out", outputBinary))))
	}

	lines = append(lines, buildInfoCommandText(definition))
	return strings.Join(lines, "\n")
}

func cargoBuildCommandText(definition packages.Definition, tools external.Tooling) string {
	return strings.Join(quoteArguments(tools.Cargo.BuildReleaseOfflineCommandTokens(definition.OutputBinaries)), " ")
}
