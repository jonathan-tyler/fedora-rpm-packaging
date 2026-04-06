package buildcontainer

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/your-github-username/fedora-package-builder/src/core/project"
)

type artifactNames struct {
	Timestamp     string
	ImageTag      string
	PrepContainer string
	SnapshotImage string
	StagedRoot    string
	CloneDir      string
}

func newArtifactNames(paths project.Paths, packageName string, now time.Time) artifactNames {
	timestamp := now.Format("20060102-150405")
	randSuffix := fmt.Sprintf("%06d", now.UnixNano()%1_000_000)

	return artifactNames{
		Timestamp:     timestamp,
		ImageTag:      fmt.Sprintf("localhost/fedora-package-builder-%s:latest", packageName),
		PrepContainer: fmt.Sprintf("fpb-%s-prep-%s", packageName, randSuffix),
		SnapshotImage: fmt.Sprintf("localhost/fedora-package-builder-%s-snapshot:%s", packageName, timestamp),
		StagedRoot:    filepath.Join(paths.PackageArtifactsRoot(packageName), timestamp),
		CloneDir:      filepath.Join("/src", packageName),
	}
}
