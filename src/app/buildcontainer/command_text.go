package buildcontainer

import (
	"fmt"
	"strings"

	"github.com/your-github-username/fedora-package-builder/src/core/commandline"
	"github.com/your-github-username/fedora-package-builder/src/core/packages"
)

const packageScriptsDir = "/pkg"
const prepareScriptPath = "/pkg/prepare.sh"
const buildScriptPath = "/pkg/build.sh"

func prepCommandText(definition packages.Definition, names artifactNames, ref string) (string, error) {
	return packageScriptCommandText(definition, names, ref, prepareScriptPath), nil
}

func buildCommandText(definition packages.Definition, names artifactNames, ref string) (string, error) {
	versionLine := `version="$(git describe --tags --always 2>/dev/null || git rev-parse --short HEAD)"`
	if _, ok := normalizeStableVersion(ref); ok {
		versionLine = fmt.Sprintf("version=%s", commandline.Quote(ref))
	}

	lines := []string{
		packageScriptEnvironmentText(definition, names, "", buildScriptPath),
		"cd " + commandline.Quote(names.CloneDir),
		versionLine,
		`commit="$(git rev-parse HEAD)"`,
		buildInfoCommandText(definition),
	}
	return strings.Join(lines, "\n"), nil
}

func packageScriptCommandText(definition packages.Definition, names artifactNames, ref string, scriptPath string) string {
	return packageScriptEnvironmentText(definition, names, ref, scriptPath)
}

func packageScriptEnvironmentText(definition packages.Definition, names artifactNames, ref string, scriptPath string) string {
	lines := []string{
		"set -euo pipefail",
		fmt.Sprintf("export FPB_PACKAGE=%s", commandline.Quote(definition.Name)),
		fmt.Sprintf("export FPB_REPO_URL=%s", commandline.Quote(definition.RepoURL)),
		fmt.Sprintf("export FPB_REF=%s", commandline.Quote(ref)),
		fmt.Sprintf("export FPB_CLONE_DIR=%s", commandline.Quote(names.CloneDir)),
		fmt.Sprintf("export FPB_OUT_DIR=%s", commandline.Quote("/out")),
		fmt.Sprintf("export FPB_OUTPUT_BINARIES=%s", commandline.Quote(strings.Join(definition.OutputBinaries, " "))),
		fmt.Sprintf("export FPB_PACKAGE_DIR=%s", commandline.Quote(packageScriptsDir)),
		fmt.Sprintf("bash %s", commandline.Quote(scriptPath)),
	}
	return strings.Join(lines, "\n")
}
