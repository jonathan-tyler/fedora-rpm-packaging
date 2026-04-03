package buildcontainer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/him/fedora-local-builder/src/core/commandline"
	"github.com/him/fedora-local-builder/src/core/packages"
	"github.com/him/fedora-local-builder/src/infra/external"
)

func prepGoCommandText(definition packages.Definition, names artifactNames, ref string, tools external.Tooling) string {
	return strings.Join([]string{
		cloneCommandText(definition, names, ref, tools),
		"cd " + commandline.Quote(names.CloneDir),
		strings.Join(quoteArguments(tools.Go.ModVendorCommandTokens()), " "),
	}, " && ")
}

func buildGoCommandText(definition packages.Definition, names artifactNames, tools external.Tooling) string {
	return strings.Join([]string{
		"set -euo pipefail",
		"cd " + commandline.Quote(names.CloneDir),
		`version="$(git describe --tags --always 2>/dev/null || git rev-parse --short HEAD)"`,
		`commit="$(git rev-parse HEAD)"`,
		fmt.Sprintf("CGO_ENABLED=0 %s build -mod=vendor -buildvcs=false -ldflags \"-X main.version=${version}\" -o %s ./", commandline.Quote(tools.Go.CommandName()), commandline.Quote(filepath.Join("/out", definition.OutputBinaries[0]))),
		buildInfoCommandText(definition),
	}, "\n")
}
