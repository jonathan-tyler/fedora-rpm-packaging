package project

import "path/filepath"

type Paths struct {
	Root string
}

func (p Paths) PackageManifest() string {
	return filepath.Join(p.Root, "config", "packages.json")
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

func (p Paths) ResultsRoot() string {
	return filepath.Join(p.Root, "state", "results")
}

func (p Paths) PackageResultsRoot(packageName string) string {
	return filepath.Join(p.ResultsRoot(), packageName)
}

func (p Paths) RepoRoot() string {
	return filepath.Join(p.Root, "state", "repo")
}

func (p Paths) RepoPackageDownloadRoot(packageName string) string {
	return filepath.Join(p.RepoRoot(), "downloads", packageName)
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
	return filepath.Join(p.Root, "container", "build", packageName)
}
