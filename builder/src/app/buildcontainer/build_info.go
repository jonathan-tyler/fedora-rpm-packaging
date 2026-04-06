package buildcontainer

import (
	"fmt"
	"os"
	"strings"
)

type buildInfo struct {
	PackageName string
	Version     string
	Commit      string
	Outputs     []string
}

func readBuildInfo(path string) (buildInfo, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return buildInfo{}, fmt.Errorf("read build info %s: %w", path, err)
	}

	return parseBuildInfo(string(content))
}

func parseBuildInfo(content string) (buildInfo, error) {
	info := buildInfo{}

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return buildInfo{}, fmt.Errorf("invalid build info line %q", line)
		}

		switch key {
		case "package":
			info.PackageName = value
		case "version":
			info.Version = value
		case "commit":
			info.Commit = value
		case "outputs":
			if value == "" {
				continue
			}
			info.Outputs = strings.Split(value, ",")
		}
	}

	if info.Version == "" {
		return buildInfo{}, fmt.Errorf("build info is missing version")
	}

	return info, nil
}
