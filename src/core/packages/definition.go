package packages

import "fmt"

type BuildSystem string

const (
	BuildSystemGo    BuildSystem = "go"
	BuildSystemCargo BuildSystem = "cargo"
)

type ContainerBinaryBuildSupport string

const (
	ContainerBinaryBuildSupported ContainerBinaryBuildSupport = "yes"
	ContainerBinaryBuildNotYet    ContainerBinaryBuildSupport = "not-yet"
)

type Definition struct {
	Name                 string                      `json:"name"`
	RepoURL              string                      `json:"repo_url"`
	BuildSystem          BuildSystem                 `json:"build_system"`
	BinaryName           string                      `json:"binary_name"`
	VendorHint           string                      `json:"vendor_hint"`
	ContainerBinaryBuild ContainerBinaryBuildSupport `json:"container_binary_build"`
}

func (d Definition) SupportsContainerBinaryBuild() bool {
	return d.ContainerBinaryBuild == ContainerBinaryBuildSupported
}

func (d Definition) Validate() error {
	if d.Name == "" {
		return fmt.Errorf("package name is required")
	}
	if d.RepoURL == "" {
		return fmt.Errorf("package %s repo URL is required", d.Name)
	}
	if d.BinaryName == "" {
		return fmt.Errorf("package %s binary name is required", d.Name)
	}
	if d.VendorHint == "" {
		return fmt.Errorf("package %s vendor hint is required", d.Name)
	}

	switch d.BuildSystem {
	case BuildSystemGo, BuildSystemCargo:
	default:
		return fmt.Errorf("package %s has unsupported build system %q", d.Name, d.BuildSystem)
	}

	switch d.ContainerBinaryBuild {
	case ContainerBinaryBuildSupported, ContainerBinaryBuildNotYet:
	default:
		return fmt.Errorf("package %s has unsupported container binary build state %q", d.Name, d.ContainerBinaryBuild)
	}

	return nil
}
