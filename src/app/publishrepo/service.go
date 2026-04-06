package publishrepo

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/your-github-username/fedora-package-builder/src/core/packages"
	"github.com/your-github-username/fedora-package-builder/src/core/platform"
	"github.com/your-github-username/fedora-package-builder/src/core/project"
	"github.com/your-github-username/fedora-package-builder/src/infra/external"
)

type Service struct {
	Registry packages.Registry
	Paths    project.Paths
	Tools    external.Tooling
	Runner   platform.Runner
	Stdout   io.Writer
}

func (s Service) Init(ctx context.Context, releasever string, basearch string, gpgKey string) error {
	repoDir, err := s.prepareRepoDir(releasever, basearch)
	if err != nil {
		return err
	}

	return s.finalizeRepository(ctx, repoDir, gpgKey, "initialized repository metadata in %s\n", "initialized signed repository metadata in %s\n", "repository metadata is unsigned; set FEDORA_PACKAGE_GPG_KEY to export the public key and generate repomd.xml.asc\n")
}

func (s Service) Run(ctx context.Context, packageName string, releasever string, basearch string, gpgKey string) error {
	if _, err := s.Registry.Lookup(packageName); err != nil {
		return err
	}

	sourceDir := s.Paths.PackageResultsRoot(packageName)
	if !directoryExists(sourceDir) {
		return fmt.Errorf("Result directory not found: %s", sourceDir)
	}

	repoDir, err := s.prepareRepoDir(releasever, basearch)
	if err != nil {
		return err
	}

	rpms, err := filepath.Glob(filepath.Join(sourceDir, "*.rpm"))
	if err != nil {
		return fmt.Errorf("glob rpms: %w", err)
	}
	if len(rpms) == 0 {
		return fmt.Errorf("No RPMs found in %s", sourceDir)
	}

	binaryRPMs := make([]string, 0, len(rpms))
	for _, rpmPath := range rpms {
		if strings.HasSuffix(rpmPath, ".src.rpm") {
			continue
		}
		targetPath := filepath.Join(repoDir, filepath.Base(rpmPath))
		if err := removeOptionalFile(targetPath + ".sig"); err != nil {
			return err
		}
		if err := copyFile(rpmPath, targetPath); err != nil {
			return err
		}
		binaryRPMs = append(binaryRPMs, targetPath)
	}

	if len(binaryRPMs) == 0 {
		return fmt.Errorf("No binary RPMs found in %s", sourceDir)
	}

	if gpgKey != "" {
		if err := s.Runner.Run(ctx, s.Tools.RPMSign.AddSignCommand(binaryRPMs, s.Stdout, s.Stdout)); err != nil {
			return err
		}
		if err := removeSignatureSidecars(binaryRPMs); err != nil {
			return err
		}
	}

	return s.finalizeRepository(ctx, repoDir, gpgKey, "published %s into %s\n", "published %s into %s with signed RPMs and repository metadata\n", "package payload and repository metadata are unsigned; set FEDORA_PACKAGE_GPG_KEY and ~/.rpmmacros before publishing for DNF clients\n", packageName)
}

func (s Service) prepareRepoDir(releasever string, basearch string) (string, error) {
	repoDir := s.Paths.RepoFedoraArchRoot(releasever, basearch)
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		return "", fmt.Errorf("create repo directory: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(s.Paths.RepoPublicKeyPath()), 0o755); err != nil {
		return "", fmt.Errorf("create key directory: %w", err)
	}

	return repoDir, nil
}

func (s Service) finalizeRepository(ctx context.Context, repoDir string, gpgKey string, unsignedMessage string, signedMessage string, unsignedNote string, args ...any) error {
	if err := s.Runner.Run(ctx, s.Tools.CreateRepo.UpdateCommand(repoDir, s.Stdout, s.Stdout)); err != nil {
		return err
	}

	if gpgKey != "" {
		if err := s.exportPublicKey(ctx, gpgKey); err != nil {
			return err
		}
		if err := s.signRepositoryMetadata(ctx, gpgKey, repoDir); err != nil {
			return err
		}
	} else {
		if err := s.clearSigningArtifacts(repoDir); err != nil {
			return err
		}
	}

	message := unsignedMessage
	if gpgKey != "" {
		message = signedMessage
	}

	formatArgs := append(args, repoDir)
	_, err := fmt.Fprintf(s.Stdout, message, formatArgs...)
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	if gpgKey == "" {
		if _, err := fmt.Fprint(s.Stdout, unsignedNote); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
	}

	return nil
}

func (s Service) exportPublicKey(ctx context.Context, gpgKey string) error {
	publicKey, err := s.Runner.Capture(ctx, s.Tools.GPG.ExportArmoredPublicKeyCommand(gpgKey, s.Stdout))
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.Paths.RepoPublicKeyPath(), []byte(publicKey), 0o644); err != nil {
		return fmt.Errorf("write public key: %w", err)
	}

	return nil
}

func (s Service) signRepositoryMetadata(ctx context.Context, gpgKey string, repoDir string) error {
	return s.Runner.Run(ctx, s.Tools.GPG.DetachSignCommand(
		gpgKey,
		filepath.Join(repoDir, "repodata", "repomd.xml.asc"),
		filepath.Join(repoDir, "repodata", "repomd.xml"),
		s.Stdout,
		s.Stdout,
	))
}

func (s Service) clearSigningArtifacts(repoDir string) error {
	paths := []string{
		filepath.Join(repoDir, "repodata", "repomd.xml.asc"),
		s.Paths.RepoPublicKeyPath(),
	}

	for _, path := range paths {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove stale signing artifact %s: %w", path, err)
		}
	}

	return nil
}

func copyFile(sourcePath string, targetPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open %s: %w", sourcePath, err)
	}
	defer sourceFile.Close()

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("create %s: %w", targetPath, err)
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		return fmt.Errorf("copy %s to %s: %w", sourcePath, targetPath, err)
	}

	return nil
}

func removeSignatureSidecars(rpms []string) error {
	for _, rpmPath := range rpms {
		if err := removeOptionalFile(rpmPath + ".sig"); err != nil {
			return err
		}
	}

	return nil
}

func removeOptionalFile(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove %s: %w", path, err)
	}

	return nil
}

func directoryExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
