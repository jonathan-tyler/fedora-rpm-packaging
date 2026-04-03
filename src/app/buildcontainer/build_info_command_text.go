package buildcontainer

import (
	"fmt"
	"strings"

	"github.com/him/fedora-local-builder/src/core/commandline"
	"github.com/him/fedora-local-builder/src/core/packages"
)

func buildInfoCommandText(definition packages.Definition) string {
	return fmt.Sprintf(
		"printf 'package=%%s\nversion=%%s\ncommit=%%s\noutputs=%%s\n' %s \"$version\" \"$commit\" %s > %s",
		commandline.Quote(definition.Name),
		commandline.Quote(strings.Join(definition.OutputBinaries, ",")),
		commandline.Quote("/out/build-info.txt"),
	)
}
