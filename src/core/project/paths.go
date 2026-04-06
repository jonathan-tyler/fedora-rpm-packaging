package project

import "path/filepath"

type Paths struct {
	Root string
}

func (p Paths) PackagesRoot() string {
	return filepath.Join(p.Root, "packages")
}

func (p Paths) PackageDir(packageName string) string {
	return filepath.Join(p.PackagesRoot(), packageName)
}

func (p Paths) PackageManifest(packageName string) string {
	return filepath.Join(p.PackageDir(packageName), "manifest.json")
}

func (p Paths) ToolingConfig() string {
	return filepath.Join(p.Root, "config", "tooling.json")
}

func (p Paths) UpstreamRoot() string {
	return filepath.Join(p.Root, "state", "upstream")
}

func (p Paths) VendorRoot() string {
	return filepath.Join(p.Root, "state", "vendor")
}

func (p Paths) VendorPackageRoot(packageName string) string {
	return filepath.Join(p.VendorRoot(), packageName)
}

func (p Paths) PackagingRoot() string {
	return filepath.Join(p.Root, "packaging")
}

func (p Paths) PackagingDir(packageName string) string {
	return filepath.Join(p.PackagingRoot(), packageName)
}

func (p Paths) PackageSpec(packageName string) string {
	return filepath.Join(p.PackagingDir(packageName), packageName+".spec")
}

func (p Paths) PackageSourcesDir(packageName string) string {
	return filepath.Join(p.PackagingDir(packageName), "SOURCES")
}

func (p Paths) ResultsRoot() string {
	return filepath.Join(p.Root, "state", "results")
}

func (p Paths) PackageResultsRoot(packageName string) string {
	return filepath.Join(p.ResultsRoot(), packageName)
}

func (p Paths) ArtifactsRoot() string {
	return filepath.Join(p.Root, "state", "artifacts")
}

func (p Paths) PackageArtifactsRoot(packageName string) string {
	return filepath.Join(p.ArtifactsRoot(), packageName)
}

func (p Paths) RepoRoot() string {
	return filepath.Join(p.Root, "state", "repo")
}

func (p Paths) RepoFedoraArchRoot(releasever string, basearch string) string {
	return filepath.Join(p.RepoRoot(), "fedora", releasever, basearch)
}

func (p Paths) RepoPublicKeyPath() string {
	return filepath.Join(p.RepoRoot(), "keys", "RPM-GPG-KEY-local-fedora-builder")
}

func (p Paths) MockConfigDir() string {
	return filepath.Join(p.Root, "mock")
}

func (p Paths) PackageBuildContainerDir(packageName string) string {
	return p.PackageDir(packageName)
}
