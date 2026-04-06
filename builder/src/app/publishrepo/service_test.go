package publishrepo

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/your-github-username/fedora-package-infra/builder/src/core/packages"
	"github.com/your-github-username/fedora-package-infra/builder/src/core/platform"
	"github.com/your-github-username/fedora-package-infra/builder/src/core/project"
	"github.com/your-github-username/fedora-package-infra/builder/src/infra/external"
)

type fakeRunner struct {
	runCommands []platform.Command
}

func (f *fakeRunner) Run(_ context.Context, command platform.Command) error {
	f.runCommands = append(f.runCommands, command)
	return nil
}

func (f *fakeRunner) Capture(_ context.Context, _ platform.Command) (string, error) {
	return "", nil
}

func TestInitRemovesStaleSigningArtifactsWhenUnsigned(t *testing.T) {
	t.Helper()

	root := t.TempDir()
	repoDir := filepath.Join(root, "state", "repo", "fedora", "42", "x86_64")
	repodataDir := filepath.Join(repoDir, "repodata")
	publicKeyPath := filepath.Join(root, "state", "repo", "keys", "RPM-GPG-KEY-local-fedora-builder")

	for _, dir := range []string{repodataDir, filepath.Dir(publicKeyPath)} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	staleFiles := []string{
		filepath.Join(repodataDir, "repomd.xml.asc"),
		publicKeyPath,
	}
	for _, path := range staleFiles {
		if err := os.WriteFile(path, []byte("stale"), 0o644); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	runner := &fakeRunner{}
	stdout := &bytes.Buffer{}
	service := Service{
		Paths: project.Paths{Root: root},
		Tools: external.Tooling{
			CreateRepo: external.CreateRepoTool{},
		},
		Runner: runner,
		Stdout: stdout,
	}

	if err := service.Init(context.Background(), "42", "x86_64", ""); err != nil {
		t.Fatalf("init repo: %v", err)
	}

	if len(runner.runCommands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(runner.runCommands))
	}

	for _, path := range staleFiles {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("expected %s to be removed, got err=%v", path, err)
		}
	}

	output := stdout.String()
	if !strings.Contains(output, "initialized repository metadata") {
		t.Fatalf("expected init output, got %q", output)
	}
	if !strings.Contains(output, "repository metadata is unsigned") {
		t.Fatalf("expected unsigned warning, got %q", output)
	}
}

func TestPublishReportsUnsignedRPMAndMetadataState(t *testing.T) {
	t.Helper()

	root := t.TempDir()
	resultsDir := filepath.Join(root, "state", "results", "demo")
	if err := os.MkdirAll(resultsDir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", resultsDir, err)
	}

	rpmPath := filepath.Join(resultsDir, "demo-1.0-1.x86_64.rpm")
	if err := os.WriteFile(rpmPath, []byte("rpm"), 0o644); err != nil {
		t.Fatalf("write %s: %v", rpmPath, err)
	}

	registry, err := packages.NewStaticRegistry([]packages.Definition{{
		Name:                 "demo",
		RepoURL:              "https://example.invalid/demo",
		BuildSystem:          packages.BuildSystemGo,
		OutputBinaries:       []string{"demo"},
		VendorHint:           "go mod vendor",
		ContainerBinaryBuild: packages.ContainerBinaryBuildSupported,
	}})
	if err != nil {
		t.Fatalf("new registry: %v", err)
	}

	runner := &fakeRunner{}
	stdout := &bytes.Buffer{}
	service := Service{
		Registry: registry,
		Paths:    project.Paths{Root: root},
		Tools: external.Tooling{
			CreateRepo: external.CreateRepoTool{},
		},
		Runner: runner,
		Stdout: stdout,
	}

	if err := service.Run(context.Background(), "demo", "42", "x86_64", ""); err != nil {
		t.Fatalf("publish repo: %v", err)
	}

	if len(runner.runCommands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(runner.runCommands))
	}

	output := stdout.String()
	if !strings.Contains(output, "published demo") {
		t.Fatalf("expected publish output, got %q", output)
	}
	if !strings.Contains(output, "package payload and repository metadata are unsigned") {
		t.Fatalf("expected unsigned publish warning, got %q", output)
	}
}

func TestPublishRemovesSignatureSidecars(t *testing.T) {
	t.Helper()

	root := t.TempDir()
	resultsDir := filepath.Join(root, "state", "results", "demo")
	if err := os.MkdirAll(resultsDir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", resultsDir, err)
	}

	rpmPath := filepath.Join(resultsDir, "demo-1.0-1.x86_64.rpm")
	if err := os.WriteFile(rpmPath, []byte("rpm"), 0o644); err != nil {
		t.Fatalf("write %s: %v", rpmPath, err)
	}

	registry, err := packages.NewStaticRegistry([]packages.Definition{{
		Name:                 "demo",
		RepoURL:              "https://example.invalid/demo",
		BuildSystem:          packages.BuildSystemGo,
		OutputBinaries:       []string{"demo"},
		VendorHint:           "go mod vendor",
		ContainerBinaryBuild: packages.ContainerBinaryBuildSupported,
	}})
	if err != nil {
		t.Fatalf("new registry: %v", err)
	}

	runner := &fakeRunner{}
	stdout := &bytes.Buffer{}
	service := Service{
		Registry: registry,
		Paths:    project.Paths{Root: root},
		Tools: external.Tooling{
			CreateRepo: external.CreateRepoTool{},
			RPMSign:    external.RPMSignTool{},
		},
		Runner: runner,
		Stdout: stdout,
	}

	repoDir := filepath.Join(root, "state", "repo", "fedora", "42", "x86_64")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", repoDir, err)
	}
	staleSidecar := filepath.Join(repoDir, "demo-1.0-1.x86_64.rpm.sig")
	if err := os.WriteFile(staleSidecar, []byte("stale"), 0o644); err != nil {
		t.Fatalf("write %s: %v", staleSidecar, err)
	}

	if err := service.Run(context.Background(), "demo", "42", "x86_64", "key"); err != nil {
		t.Fatalf("publish repo: %v", err)
	}

	if _, err := os.Stat(staleSidecar); !os.IsNotExist(err) {
		t.Fatalf("expected stale sidecar to be removed, got err=%v", err)
	}

	foundRPMsign := false
	for _, command := range runner.runCommands {
		if len(command.Args) > 0 && command.Args[0] == "--addsign" {
			foundRPMsign = true
			break
		}
	}
	if !foundRPMsign {
		t.Fatalf("expected an rpmsign command, got %#v", runner.runCommands)
	}
}
