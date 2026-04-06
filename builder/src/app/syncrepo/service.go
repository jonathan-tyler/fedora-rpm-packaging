package syncrepo

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/core/project"
)

type Service struct {
	Paths  project.Paths
	Stdout io.Writer
}

func (s Service) Run(stagedRepoRoot string, liveRepoRoot string, tempRoot string) error {
	if !directoryExists(stagedRepoRoot) {
		return fmt.Errorf("staged repo not found: %s", stagedRepoRoot)
	}

	if err := os.MkdirAll(liveRepoRoot, 0o755); err != nil {
		return fmt.Errorf("create live repo root: %w", err)
	}
	if err := os.RemoveAll(tempRoot); err != nil {
		return fmt.Errorf("clean temp root: %w", err)
	}
	if err := os.MkdirAll(tempRoot, 0o755); err != nil {
		return fmt.Errorf("create temp root: %w", err)
	}

	if err := copyTree(stagedRepoRoot, tempRoot); err != nil {
		return err
	}
	if err := normalizePermissions(tempRoot); err != nil {
		return err
	}

	if err := os.RemoveAll(liveRepoRoot); err != nil {
		return fmt.Errorf("replace live repo root: %w", err)
	}
	if err := os.MkdirAll(liveRepoRoot, 0o755); err != nil {
		return fmt.Errorf("recreate live repo root: %w", err)
	}
	if err := copyTree(tempRoot, liveRepoRoot); err != nil {
		return err
	}
	if err := os.RemoveAll(tempRoot); err != nil {
		return fmt.Errorf("remove temp root: %w", err)
	}

	_, err := fmt.Fprintf(s.Stdout, "synced staged repo from %s to %s\n", stagedRepoRoot, liveRepoRoot)
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}

func copyTree(sourceRoot string, targetRoot string) error {
	return filepath.Walk(sourceRoot, func(sourcePath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		relativePath, err := filepath.Rel(sourceRoot, sourcePath)
		if err != nil {
			return fmt.Errorf("relative path for %s: %w", sourcePath, err)
		}

		targetPath := targetRoot
		if relativePath != "." {
			targetPath = filepath.Join(targetRoot, relativePath)
		}

		if info.IsDir() {
			return os.MkdirAll(targetPath, 0o755)
		}

		if info.Mode()&os.ModeSymlink != 0 {
			linkTarget, err := os.Readlink(sourcePath)
			if err != nil {
				return fmt.Errorf("read symlink %s: %w", sourcePath, err)
			}
			if err := os.RemoveAll(targetPath); err != nil {
				return fmt.Errorf("replace symlink %s: %w", targetPath, err)
			}
			return os.Symlink(linkTarget, targetPath)
		}

		return copyFile(sourcePath, targetPath, info.Mode().Perm())
	})
}

func normalizePermissions(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		if info.IsDir() {
			return os.Chmod(path, 0o755)
		}
		return os.Chmod(path, 0o644)
	})
}

func copyFile(sourcePath string, targetPath string, mode os.FileMode) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open %s: %w", sourcePath, err)
	}
	defer sourceFile.Close()

	targetFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("create %s: %w", targetPath, err)
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		return fmt.Errorf("copy %s to %s: %w", sourcePath, targetPath, err)
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
