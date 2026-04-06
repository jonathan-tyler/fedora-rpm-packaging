package buildcontainer

import (
	"bytes"
	"context"
	"testing"

	"github.com/your-github-username/fedora-package-infra/builder/src/core/packages"
	"github.com/your-github-username/fedora-package-infra/builder/src/core/platform"
	"github.com/your-github-username/fedora-package-infra/builder/src/infra/external"
)

type fakeTagRunner struct {
	captureOutput  string
	captureCommand platform.Command
	captureCalls   int
}

func (f *fakeTagRunner) Run(_ context.Context, _ platform.Command) error {
	return nil
}

func (f *fakeTagRunner) Capture(_ context.Context, command platform.Command) (string, error) {
	f.captureCommand = command
	f.captureCalls++
	return f.captureOutput, nil
}

func TestResolveArtifactTag(t *testing.T) {
	registry, err := packages.NewStaticRegistry([]packages.Definition{{
		Name:                 "demo",
		RepoURL:              "https://example.invalid/demo.git",
		BuildSystem:          packages.BuildSystemCargo,
		OutputBinaries:       []string{"demo"},
		VendorHint:           "cargo vendor",
		ContainerBinaryBuild: packages.ContainerBinaryBuildSupported,
	}})
	if err != nil {
		t.Fatalf("new registry: %v", err)
	}

	runner := &fakeTagRunner{captureOutput: "aaa\trefs/tags/nightly\nbbb\trefs/tags/v26.1.4\nccc\trefs/tags/v26.1.22\n"}
	service := Service{
		Registry: registry,
		Tools:    external.Tooling{Git: external.Git{}},
		Runner:   runner,
		Stdout:   &bytes.Buffer{},
	}

	resolved, err := service.ResolveArtifactTag(context.Background(), "demo")
	if err != nil {
		t.Fatalf("resolve artifact tag: %v", err)
	}

	if resolved.UpstreamTag != "v26.1.22" {
		t.Fatalf("upstream tag = %q, want %q", resolved.UpstreamTag, "v26.1.22")
	}
	if resolved.ArtifactTag != "26.1.22" {
		t.Fatalf("artifact tag = %q, want %q", resolved.ArtifactTag, "26.1.22")
	}
	if runner.captureCalls != 1 {
		t.Fatalf("capture calls = %d, want %d", runner.captureCalls, 1)
	}
	if len(runner.captureCommand.Args) < 4 || runner.captureCommand.Args[0] != "ls-remote" {
		t.Fatalf("unexpected capture command: %#v", runner.captureCommand)
	}
}

func TestResolveArtifactTagWithoutStableTag(t *testing.T) {
	registry, err := packages.NewStaticRegistry([]packages.Definition{{
		Name:                 "demo",
		RepoURL:              "https://example.invalid/demo.git",
		BuildSystem:          packages.BuildSystemCargo,
		OutputBinaries:       []string{"demo"},
		VendorHint:           "cargo vendor",
		ContainerBinaryBuild: packages.ContainerBinaryBuildSupported,
	}})
	if err != nil {
		t.Fatalf("new registry: %v", err)
	}

	service := Service{
		Registry: registry,
		Tools:    external.Tooling{Git: external.Git{}},
		Runner:   &fakeTagRunner{captureOutput: "aaa\trefs/tags/nightly\n"},
		Stdout:   &bytes.Buffer{},
	}

	if _, err := service.ResolveArtifactTag(context.Background(), "demo"); err == nil {
		t.Fatal("expected error when no stable tag exists")
	}
}
