package vendorsources

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	corearchive "github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/core/archive"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/core/commandline"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/core/packages"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/core/platform"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/core/project"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/infra/external"
)

type Service struct {
	Registry packages.Registry
	Paths    project.Paths
	Tools    external.Tooling
	Runner   platform.Runner
	Archiver corearchive.Archiver
	Stdout   io.Writer
}

func (s Service) Run(ctx context.Context, packageName string, worktree string) error {
	definition, err := s.Registry.Lookup(packageName)
	if err != nil {
		return err
	}

	worktree, err = filepath.Abs(worktree)
	if err != nil {
		return fmt.Errorf("resolve worktree path: %w", err)
	}

	if !directoryExists(worktree) {
		return fmt.Errorf("worktree not found: %s", worktree)
	}

	vendorRoot := s.Paths.VendorPackageRoot(packageName)
	if err := os.RemoveAll(vendorRoot); err != nil {
		return fmt.Errorf("clean vendor root: %w", err)
	}
	if err := os.MkdirAll(vendorRoot, 0o755); err != nil {
		return fmt.Errorf("create vendor root: %w", err)
	}

	switch definition.BuildSystem {
	case packages.BuildSystemGo:
		if err := s.vendorGo(ctx, definition, worktree, vendorRoot); err != nil {
			return err
		}
	case packages.BuildSystemCargo:
		if err := s.vendorCargo(ctx, definition, worktree, vendorRoot); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported build system: %s", definition.BuildSystem)
	}

	_, err = fmt.Fprintf(s.Stdout, "vendored %s with %s into %s\n", definition.Name, definition.BuildSystem, vendorRoot)
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}

func (s Service) vendorGo(ctx context.Context, definition packages.Definition, worktree string, vendorRoot string) error {
	if err := s.Runner.Run(ctx, s.Tools.Go.ModVendorCommand(worktree, s.Stdout, s.Stdout)); err != nil {
		return err
	}

	archivePath := filepath.Join(vendorRoot, definition.Name+"-vendor.tar.gz")
	return s.Archiver.CreateTarGz(archivePath, []corearchive.Entry{{
		SourcePath:  filepath.Join(worktree, "vendor"),
		ArchivePath: "vendor",
	}})
}

func (s Service) vendorCargo(ctx context.Context, definition packages.Definition, worktree string, vendorRoot string) error {
	vendorDirectory := filepath.Join(vendorRoot, "vendor")
	if err := s.Runner.Run(ctx, s.Tools.ContainerRuntime.BuildCommand(vendorImageTag(definition.Name), s.Paths.PackageBuildContainerDir(definition.Name), s.Stdout, s.Stdout)); err != nil {
		return err
	}

	configOutput, err := s.Runner.Capture(ctx, platform.Command{
		Name: s.Tools.ContainerRuntime.CommandName(),
		Args: []string{
			"run",
			"--rm",
			"-v", worktree + ":/src:ro,Z",
			"-v", vendorRoot + ":/vendor:Z",
			vendorImageTag(definition.Name),
			"bash",
			"-lc",
			cargoVendorCommandText(vendorDirectory),
		},
		Stderr: s.Stdout,
	})
	if err != nil {
		return err
	}
	configOutput = normalizeCargoVendorConfig(configOutput)

	configPath := filepath.Join(vendorRoot, "cargo-config.toml")
	if err := os.WriteFile(configPath, []byte(configOutput), 0o644); err != nil {
		return fmt.Errorf("write cargo config: %w", err)
	}

	archivePath := filepath.Join(vendorRoot, definition.Name+"-vendor.tar.gz")
	return s.Archiver.CreateTarGz(archivePath, []corearchive.Entry{
		{SourcePath: vendorDirectory, ArchivePath: "vendor"},
		{SourcePath: configPath, ArchivePath: "cargo-config.toml"},
	})
}

func vendorImageTag(packageName string) string {
	return fmt.Sprintf("localhost/fedora-package-builder-%s:latest", packageName)
}

func cargoVendorCommandText(vendorDirectory string) string {
	tokens := []string{"cd", "/src", "&&"}
	for _, token := range []string{"cargo", "vendor", filepath.ToSlash(filepath.Join("/vendor", filepath.Base(vendorDirectory)))} {
		tokens = append(tokens, commandline.Quote(token))
	}
	return strings.Join(tokens, " ")
}

func normalizeCargoVendorConfig(configOutput string) string {
	return strings.ReplaceAll(configOutput, "/vendor/vendor", "vendor")
}

func directoryExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
