package external

import (
	"reflect"
	"testing"
)

func TestMirrorCloneCommandUsesConfiguredExtraArgs(t *testing.T) {
	git, err := NewGit(GitConfig{
		Command:              "git",
		MirrorCloneExtraArgs: []string{"--depth", "1"},
	})
	if err != nil {
		t.Fatalf("NewGit() error = %v", err)
	}

	command := git.MirrorCloneCommand("https://example.invalid/repo.git", "/tmp/repo.git", nil, nil)
	wantArgs := []string{"clone", "--depth", "1", "--mirror", "https://example.invalid/repo.git", "/tmp/repo.git"}
	if !reflect.DeepEqual(command.Args, wantArgs) {
		t.Fatalf("MirrorCloneCommand() args = %#v, want %#v", command.Args, wantArgs)
	}
}

func TestFetchMirrorCommandUsesConfiguredExtraArgs(t *testing.T) {
	git, err := NewGit(GitConfig{
		Command:              "git",
		MirrorFetchExtraArgs: []string{"--depth", "1"},
	})
	if err != nil {
		t.Fatalf("NewGit() error = %v", err)
	}

	command := git.FetchMirrorCommand("/tmp/repo.git", nil, nil)
	wantArgs := []string{"-C", "/tmp/repo.git", "fetch", "--depth", "1", "--prune", "--tags", "origin"}
	if !reflect.DeepEqual(command.Args, wantArgs) {
		t.Fatalf("FetchMirrorCommand() args = %#v, want %#v", command.Args, wantArgs)
	}
}
