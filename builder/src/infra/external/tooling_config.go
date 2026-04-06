package external

type ToolConfig struct {
	Command string `json:"command"`
}

type GitConfig struct {
	Command              string   `json:"command"`
	CloneExtraArgs       []string `json:"clone_extra_args"`
	MirrorCloneExtraArgs []string `json:"mirror_clone_extra_args"`
	MirrorFetchExtraArgs []string `json:"mirror_fetch_extra_args"`
}

type ContainerRuntimeConfig struct {
	Kind                     string   `json:"kind"`
	BuildExtraArgs           []string `json:"build_extra_args"`
	PrepareRunExtraArgs      []string `json:"prepare_run_extra_args"`
	OfflineRunExtraArgs      []string `json:"offline_run_extra_args"`
	CommitExtraArgs          []string `json:"commit_extra_args"`
	RemoveContainerExtraArgs []string `json:"remove_container_extra_args"`
	RemoveImageExtraArgs     []string `json:"remove_image_extra_args"`
}

type ToolingConfig struct {
	Git              GitConfig              `json:"git"`
	ContainerRuntime ContainerRuntimeConfig `json:"container_runtime"`
	Go               ToolConfig             `json:"go"`
	Cargo            ToolConfig             `json:"cargo"`
	SRPMBuilder      ToolConfig             `json:"srpm_builder"`
	Mock             ToolConfig             `json:"mock"`
	RPMSign          ToolConfig             `json:"rpmsign"`
	CreateRepo       ToolConfig             `json:"createrepo"`
	GPG              ToolConfig             `json:"gpg"`
	Systemctl        ToolConfig             `json:"systemctl"`
}

func DefaultToolingConfig() ToolingConfig {
	return ToolingConfig{
		Git: GitConfig{
			Command:              "git",
			CloneExtraArgs:       []string{"--depth", "1"},
			MirrorCloneExtraArgs: []string{"--depth", "1"},
			MirrorFetchExtraArgs: []string{"--depth", "1"},
		},
		ContainerRuntime: ContainerRuntimeConfig{
			Kind:                     "podman",
			BuildExtraArgs:           []string{},
			PrepareRunExtraArgs:      []string{},
			OfflineRunExtraArgs:      []string{"--network", "none"},
			CommitExtraArgs:          []string{},
			RemoveContainerExtraArgs: []string{},
			RemoveImageExtraArgs:     []string{},
		},
		Go:          ToolConfig{Command: "go"},
		Cargo:       ToolConfig{Command: "cargo"},
		SRPMBuilder: ToolConfig{Command: "./scripts/build-srpm-container.sh"},
		Mock:        ToolConfig{Command: "./scripts/mock-rebuild-container.sh"},
		RPMSign:     ToolConfig{Command: "rpmsign"},
		CreateRepo:  ToolConfig{Command: "./scripts/createrepo-container.sh"},
		GPG:         ToolConfig{Command: "gpg"},
		Systemctl:   ToolConfig{Command: "systemctl"},
	}
}
