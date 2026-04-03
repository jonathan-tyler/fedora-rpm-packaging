package external

import "fmt"

type Tooling struct {
	Git              Git
	ContainerRuntime ContainerRuntime
	Go               GoTool
	Cargo            CargoTool
	Mock             MockTool
	RPMSign          RPMSignTool
	CreateRepo       CreateRepoTool
	GPG              GPGTool
	Systemctl        SystemctlTool
}

func NewTooling(config ToolingConfig) (Tooling, error) {
	git, err := NewGit(config.Git)
	if err != nil {
		return Tooling{}, fmt.Errorf("build git tooling: %w", err)
	}

	containerRuntime, err := NewContainerRuntime(config.ContainerRuntime)
	if err != nil {
		return Tooling{}, fmt.Errorf("build container runtime tooling: %w", err)
	}

	goTool, err := NewGoTool(config.Go)
	if err != nil {
		return Tooling{}, fmt.Errorf("build go tooling: %w", err)
	}

	cargoTool, err := NewCargoTool(config.Cargo)
	if err != nil {
		return Tooling{}, fmt.Errorf("build cargo tooling: %w", err)
	}

	mockTool, err := NewMockTool(config.Mock)
	if err != nil {
		return Tooling{}, fmt.Errorf("build mock tooling: %w", err)
	}

	rpmSignTool, err := NewRPMSignTool(config.RPMSign)
	if err != nil {
		return Tooling{}, fmt.Errorf("build rpmsign tooling: %w", err)
	}

	createRepoTool, err := NewCreateRepoTool(config.CreateRepo)
	if err != nil {
		return Tooling{}, fmt.Errorf("build createrepo tooling: %w", err)
	}

	gpgTool, err := NewGPGTool(config.GPG)
	if err != nil {
		return Tooling{}, fmt.Errorf("build gpg tooling: %w", err)
	}

	systemctlTool, err := NewSystemctlTool(config.Systemctl)
	if err != nil {
		return Tooling{}, fmt.Errorf("build systemctl tooling: %w", err)
	}

	return Tooling{
		Git:              git,
		ContainerRuntime: containerRuntime,
		Go:               goTool,
		Cargo:            cargoTool,
		Mock:             mockTool,
		RPMSign:          rpmSignTool,
		CreateRepo:       createRepoTool,
		GPG:              gpgTool,
		Systemctl:        systemctlTool,
	}, nil
}

func validateCommandName(command string, name string) error {
	if command == "" {
		return fmt.Errorf("%s command is required", name)
	}
	return nil
}
