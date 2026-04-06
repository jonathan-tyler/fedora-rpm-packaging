package external

import (
	"fmt"
	"io"

	"github.com/your-github-username/fedora-package-infra/builder/src/core/platform"
)

type ContainerRuntime interface {
	Kind() string
	CommandName() string
	BuildCommand(imageTag string, contextDir string, stdout io.Writer, stderr io.Writer) platform.Command
	PrepareRunCommand(containerName string, image string, commandText string, stdout io.Writer, stderr io.Writer) platform.Command
	OfflineRunCommand(mount string, image string, commandText string, stdout io.Writer, stderr io.Writer) platform.Command
	CommitCommand(containerName string, image string, stdout io.Writer, stderr io.Writer) platform.Command
	RemoveContainerCommand(containerName string, stdout io.Writer, stderr io.Writer) platform.Command
	RemoveImageCommand(image string, stdout io.Writer, stderr io.Writer) platform.Command
}

type podmanCompatibleRuntime struct {
	kind                     string
	command                  string
	buildExtraArgs           []string
	prepareRunExtraArgs      []string
	offlineRunExtraArgs      []string
	commitExtraArgs          []string
	removeContainerExtraArgs []string
	removeImageExtraArgs     []string
}

func NewContainerRuntime(config ContainerRuntimeConfig) (ContainerRuntime, error) {
	kind := config.Kind
	if kind == "" {
		kind = "podman"
	}

	switch kind {
	case "podman":
		return newPodmanCompatibleRuntime("podman", config), nil
	case "docker":
		return newPodmanCompatibleRuntime("docker", config), nil
	default:
		return nil, fmt.Errorf("unsupported container runtime kind %q", kind)
	}
}

func newPodmanCompatibleRuntime(command string, config ContainerRuntimeConfig) ContainerRuntime {
	return podmanCompatibleRuntime{
		kind:                     config.Kind,
		command:                  command,
		buildExtraArgs:           append([]string{}, config.BuildExtraArgs...),
		prepareRunExtraArgs:      append([]string{}, config.PrepareRunExtraArgs...),
		offlineRunExtraArgs:      append([]string{}, config.OfflineRunExtraArgs...),
		commitExtraArgs:          append([]string{}, config.CommitExtraArgs...),
		removeContainerExtraArgs: append([]string{}, config.RemoveContainerExtraArgs...),
		removeImageExtraArgs:     append([]string{}, config.RemoveImageExtraArgs...),
	}
}

func (r podmanCompatibleRuntime) Kind() string {
	if r.kind == "" {
		return "podman"
	}
	return r.kind
}

func (r podmanCompatibleRuntime) CommandName() string {
	return r.command
}

func (r podmanCompatibleRuntime) BuildCommand(imageTag string, contextDir string, stdout io.Writer, stderr io.Writer) platform.Command {
	args := []string{"build"}
	args = append(args, r.buildExtraArgs...)
	args = append(args, "-t", imageTag, contextDir)

	return platform.Command{Name: r.command, Args: args, Stdout: stdout, Stderr: stderr}
}

func (r podmanCompatibleRuntime) PrepareRunCommand(containerName string, image string, commandText string, stdout io.Writer, stderr io.Writer) platform.Command {
	args := []string{"run"}
	args = append(args, r.prepareRunExtraArgs...)
	args = append(args, "--name", containerName, image, "bash", "-lc", commandText)

	return platform.Command{Name: r.command, Args: args, Stdout: stdout, Stderr: stderr}
}

func (r podmanCompatibleRuntime) OfflineRunCommand(mount string, image string, commandText string, stdout io.Writer, stderr io.Writer) platform.Command {
	args := []string{"run", "--rm"}
	args = append(args, r.offlineRunExtraArgs...)
	args = append(args, "-v", mount, image, "bash", "-lc", commandText)

	return platform.Command{Name: r.command, Args: args, Stdout: stdout, Stderr: stderr}
}

func (r podmanCompatibleRuntime) CommitCommand(containerName string, image string, stdout io.Writer, stderr io.Writer) platform.Command {
	args := []string{"commit"}
	args = append(args, r.commitExtraArgs...)
	args = append(args, containerName, image)

	return platform.Command{Name: r.command, Args: args, Stdout: stdout, Stderr: stderr}
}

func (r podmanCompatibleRuntime) RemoveContainerCommand(containerName string, stdout io.Writer, stderr io.Writer) platform.Command {
	args := []string{"rm"}
	args = append(args, r.removeContainerExtraArgs...)
	args = append(args, "-f", containerName)

	return platform.Command{Name: r.command, Args: args, Stdout: stdout, Stderr: stderr}
}

func (r podmanCompatibleRuntime) RemoveImageCommand(image string, stdout io.Writer, stderr io.Writer) platform.Command {
	args := []string{"image", "rm"}
	args = append(args, r.removeImageExtraArgs...)
	args = append(args, "-f", image)

	return platform.Command{Name: r.command, Args: args, Stdout: stdout, Stderr: stderr}
}
