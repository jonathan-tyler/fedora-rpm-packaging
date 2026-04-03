package external

type ToolConfig struct {
	Command string `json:"command"`
}

type GitConfig struct {
	Command              string   `json:"command"`
	CloneExtraArgs       []string `json:"clone_extra_args"`
	MirrorCloneExtraArgs []string `json:"mirror_clone_extra_args"`
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
			MirrorCloneExtraArgs: []string{},
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
		Go:         ToolConfig{Command: "go"},
		Cargo:      ToolConfig{Command: "cargo"},
		Mock:       ToolConfig{Command: "mock"},
		RPMSign:    ToolConfig{Command: "rpmsign"},
		CreateRepo: ToolConfig{Command: "createrepo_c"},
		GPG:        ToolConfig{Command: "gpg"},
		Systemctl:  ToolConfig{Command: "systemctl"},
	}
}
