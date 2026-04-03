package buildcontainer

import (
	"strings"

	"github.com/him/fedora-local-builder/src/core/commandline"
	"github.com/him/fedora-local-builder/src/core/packages"
	"github.com/him/fedora-local-builder/src/infra/external"
)

func cloneCommandText(definition packages.Definition, names artifactNames, ref string, tools external.Tooling) string {
	return strings.Join(quoteArguments(tools.Git.CloneCommandTokens(definition.RepoURL, names.CloneDir, ref)), " ")
}

func quoteArguments(arguments []string) []string {
	quotedArguments := make([]string, 0, len(arguments))
	for _, argument := range arguments {
		quotedArguments = append(quotedArguments, commandline.Quote(argument))
	}
	return quotedArguments
}
