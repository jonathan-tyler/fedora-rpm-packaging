package buildcontainer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/him/fedora-local-builder/src/core/commandline"
	"github.com/him/fedora-local-builder/src/core/packages"
)

func prepCommandText(definition packages.Definition, names artifactNames, ref string) string {
	cloneArguments := []string{"git", "clone", "--depth", "1"}
	if ref != "" {
		cloneArguments = append(cloneArguments, "--branch", ref)
	}
	cloneArguments = append(cloneArguments, definition.RepoURL, names.CloneDir)

	quotedArguments := make([]string, 0, len(cloneArguments))
	for _, argument := range cloneArguments {
		quotedArguments = append(quotedArguments, commandline.Quote(argument))
	}

	return strings.Join([]string{
		strings.Join(quotedArguments, " "),
		"cd " + commandline.Quote(names.CloneDir),
		"go mod vendor",
	}, " && ")
}

func buildCommandText(definition packages.Definition, names artifactNames) string {
	return strings.Join([]string{
		"set -euo pipefail",
		"cd " + commandline.Quote(names.CloneDir),
		`version="$(git describe --tags --always 2>/dev/null || git rev-parse --short HEAD)"`,
		`commit="$(git rev-parse HEAD)"`,
		fmt.Sprintf("CGO_ENABLED=0 go build -mod=vendor -buildvcs=false -ldflags \"-X main.version=${version}\" -o %s ./", commandline.Quote(filepath.Join("/out", definition.BinaryName))),
		fmt.Sprintf("printf 'package=%%s\nversion=%%s\ncommit=%%s\n' %s \"$version\" \"$commit\" > %s", commandline.Quote(definition.Name), commandline.Quote("/out/build-info.txt")),
	}, "\n")
}
