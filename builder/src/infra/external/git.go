package external

import (
	"io"

	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/core/platform"
)

type Git struct {
	command              string
	cloneExtraArgs       []string
	mirrorCloneExtraArgs []string
	mirrorFetchExtraArgs []string
}

func NewGit(config GitConfig) (Git, error) {
	if err := validateCommandName(config.Command, "git"); err != nil {
		return Git{}, err
	}

	return Git{
		command:              config.Command,
		cloneExtraArgs:       append([]string{}, config.CloneExtraArgs...),
		mirrorCloneExtraArgs: append([]string{}, config.MirrorCloneExtraArgs...),
		mirrorFetchExtraArgs: append([]string{}, config.MirrorFetchExtraArgs...),
	}, nil
}

func (g Git) CloneCommandTokens(repoURL string, destination string, ref string) []string {
	args := []string{g.command, "clone"}
	args = append(args, g.cloneExtraArgs...)
	if ref != "" {
		args = append(args, "--branch", ref)
	}
	args = append(args, repoURL, destination)
	return args
}

func (g Git) MirrorCloneCommand(repoURL string, destination string, stdout io.Writer, stderr io.Writer) platform.Command {
	args := []string{"clone"}
	args = append(args, g.mirrorCloneExtraArgs...)
	args = append(args, "--mirror", repoURL, destination)

	return platform.Command{
		Name:   g.command,
		Args:   args,
		Stdout: stdout,
		Stderr: stderr,
	}
}

func (g Git) FetchMirrorCommand(target string, stdout io.Writer, stderr io.Writer) platform.Command {
	args := []string{"-C", target, "fetch"}
	args = append(args, g.mirrorFetchExtraArgs...)
	args = append(args, "--prune", "--tags", "origin")

	return platform.Command{
		Name:   g.command,
		Args:   args,
		Stdout: stdout,
		Stderr: stderr,
	}
}

func (g Git) ListRemoteTagsCommand(repoURL string, stderr io.Writer) platform.Command {
	return platform.Command{
		Name:   g.command,
		Args:   []string{"ls-remote", "--refs", "--tags", repoURL},
		Stderr: stderr,
	}
}

func (g Git) CommandName() string {
	return g.command
}
