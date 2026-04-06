package vendorsources

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	corearchive "github.com/your-github-username/fedora-package-builder/src/core/archive"
	"github.com/your-github-username/fedora-package-builder/src/core/packages"
	"github.com/your-github-username/fedora-package-builder/src/core/platform"
	"github.com/your-github-username/fedora-package-builder/src/core/project"
	"github.com/your-github-username/fedora-package-builder/src/infra/external"
)

type fakeRunner struct {
	runCommands   []platform.Command
	captureOutput string
	captureCalls  []platform.Command
}

func (f *fakeRunner) Run(_ context.Context, command platform.Command) error {
	f.runCommands = append(f.runCommands, command)
	return nil
}

func (f *fakeRunner) Capture(_ context.Context, command platform.Command) (string, error) {
	f.captureCalls = append(f.captureCalls, command)
	return f.captureOutput, nil
}

type fakeArchiver struct {
	destination string
	entries     []corearchive.Entry
}

func (f *fakeArchiver) CreateTarGz(destinationPath string, entries []corearchive.Entry) error {
	f.destination = destinationPath
	f.entries = append([]corearchive.Entry{}, entries...)
	return nil
}

func TestVendorCargoUsesPackageContainer(t *testing.T) {
	root := t.TempDir()
	worktree := filepath.Join(root, "worktree")
	if err := os.MkdirAll(worktree, 0o755); err != nil {
		t.Fatalf("mkdir worktree: %v", err)
	}

	registry, err := packages.NewStaticRegistry([]packages.Definition{{
		Name:                 "television",
		RepoURL:              "https://example.invalid/television.git",
		BuildSystem:          packages.BuildSystemCargo,
		OutputBinaries:       []string{"tv"},
		VendorHint:           "cargo vendor",
		ContainerBinaryBuild: packages.ContainerBinaryBuildSupported,
	}})
	if err != nil {
		t.Fatalf("new registry: %v", err)
	}

	runner := &fakeRunner{captureOutput: "[source.crates-io]\nreplace-with = \"vendored-sources\"\n\n[source.vendored-sources]\ndirectory = \"/vendor/vendor\"\n"}
	archiver := &fakeArchiver{}
	service := Service{
		Registry: registry,
		Paths:    project.Paths{Root: root},
		Tools: external.Tooling{
			ContainerRuntime: mustNewContainerRuntime(t),
		},
		Runner:   runner,
		Archiver: archiver,
		Stdout:   &bytes.Buffer{},
	}

	if err := service.Run(context.Background(), "television", worktree); err != nil {
		t.Fatalf("run: %v", err)
	}

	if len(runner.runCommands) != 1 {
		t.Fatalf("run commands = %d, want %d", len(runner.runCommands), 1)
	}
	if len(runner.captureCalls) != 1 {
		t.Fatalf("capture calls = %d, want %d", len(runner.captureCalls), 1)
	}
	if got := runner.runCommands[0].Args[0]; got != "build" {
		t.Fatalf("first run arg = %q, want %q", got, "build")
	}
	if got := runner.captureCalls[0].Args[0]; got != "run" {
		t.Fatalf("first capture arg = %q, want %q", got, "run")
	}
	if archiver.destination == "" {
		t.Fatal("expected vendor archive to be created")
	}
	configPath := filepath.Join(root, "state", "vendor", "television", "cargo-config.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	wantConfig := "[source.crates-io]\nreplace-with = \"vendored-sources\"\n\n[source.vendored-sources]\ndirectory = \"vendor\"\n"
	if string(content) != wantConfig {
		t.Fatalf("config output = %q, want %q", string(content), wantConfig)
	}
}

func mustNewContainerRuntime(t *testing.T) external.ContainerRuntime {
	t.Helper()

	runtime, err := external.NewContainerRuntime(external.ContainerRuntimeConfig{Kind: "podman"})
	if err != nil {
		t.Fatalf("new container runtime: %v", err)
	}

	return runtime
}
