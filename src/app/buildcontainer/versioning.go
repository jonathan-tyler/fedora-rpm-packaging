package buildcontainer

import (
	"regexp"
	"strconv"
	"strings"
)

var stableVersionPattern = regexp.MustCompile(`^v?\d+(?:\.\d+)+$`)

func latestStableVersionTag(remoteOutput string) string {
	latestTag := ""
	latestParts := []int(nil)

	for _, line := range strings.Split(remoteOutput, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		ref := fields[len(fields)-1]
		if !strings.HasPrefix(ref, "refs/tags/") {
			continue
		}

		tag := strings.TrimPrefix(ref, "refs/tags/")
		normalized, ok := normalizeStableVersion(tag)
		if !ok {
			continue
		}

		parts := versionParts(normalized)
		if compareVersionParts(parts, latestParts) > 0 {
			latestTag = tag
			latestParts = parts
		}
	}

	return latestTag
}

func artifactDirectoryName(version string, fallback string) string {
	normalized, ok := normalizeStableVersion(version)
	if !ok {
		return fallback
	}

	return normalized
}

func normalizeStableVersion(version string) (string, bool) {
	trimmed := strings.TrimSpace(version)
	if !stableVersionPattern.MatchString(trimmed) {
		return "", false
	}

	return strings.TrimPrefix(trimmed, "v"), true
}

func versionParts(version string) []int {
	segments := strings.Split(version, ".")
	parts := make([]int, 0, len(segments))
	for _, segment := range segments {
		value, err := strconv.Atoi(segment)
		if err != nil {
			return nil
		}
		parts = append(parts, value)
	}

	return parts
}

func compareVersionParts(left []int, right []int) int {
	if len(left) == 0 {
		if len(right) == 0 {
			return 0
		}
		return -1
	}
	if len(right) == 0 {
		return 1
	}

	maxLen := len(left)
	if len(right) > maxLen {
		maxLen = len(right)
	}

	for i := 0; i < maxLen; i++ {
		leftPart := 0
		if i < len(left) {
			leftPart = left[i]
		}

		rightPart := 0
		if i < len(right) {
			rightPart = right[i]
		}

		if leftPart < rightPart {
			return -1
		}
		if leftPart > rightPart {
			return 1
		}
	}

	return 0
}
