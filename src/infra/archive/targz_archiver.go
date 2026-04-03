package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	corearchive "github.com/him/fedora-local-builder/src/core/archive"
)

type TarGzArchiver struct{}

func (TarGzArchiver) CreateTarGz(destinationPath string, entries []corearchive.Entry) error {
	file, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("create archive %s: %w", destinationPath, err)
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	for _, entry := range entries {
		if err := addEntry(tarWriter, entry); err != nil {
			return err
		}
	}

	return nil
}

func addEntry(writer *tar.Writer, entry corearchive.Entry) error {
	info, err := os.Stat(entry.SourcePath)
	if err != nil {
		return fmt.Errorf("stat %s: %w", entry.SourcePath, err)
	}

	if !info.IsDir() {
		return writeFile(writer, entry.SourcePath, entry.ArchivePath, info)
	}

	return filepath.Walk(entry.SourcePath, func(path string, walkInfo os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		relativePath, err := filepath.Rel(entry.SourcePath, path)
		if err != nil {
			return fmt.Errorf("relative path for %s: %w", path, err)
		}

		archivePath := entry.ArchivePath
		if relativePath != "." {
			archivePath = filepath.Join(entry.ArchivePath, relativePath)
		}

		if walkInfo.IsDir() {
			return writeDirectory(writer, archivePath, walkInfo)
		}

		return writeFile(writer, path, archivePath, walkInfo)
	})
}

func writeDirectory(writer *tar.Writer, archivePath string, info os.FileInfo) error {
	if archivePath == "" || archivePath == "." {
		return nil
	}

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return fmt.Errorf("build tar header for %s: %w", archivePath, err)
	}
	header.Name = toArchivePath(archivePath) + "/"
	return writer.WriteHeader(header)
}

func writeFile(writer *tar.Writer, sourcePath string, archivePath string, info os.FileInfo) error {
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return fmt.Errorf("build tar header for %s: %w", sourcePath, err)
	}
	header.Name = toArchivePath(archivePath)

	if err := writer.WriteHeader(header); err != nil {
		return fmt.Errorf("write tar header for %s: %w", sourcePath, err)
	}

	file, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open %s: %w", sourcePath, err)
	}
	defer file.Close()

	if _, err := io.Copy(writer, file); err != nil {
		return fmt.Errorf("write tar body for %s: %w", sourcePath, err)
	}

	return nil
}

func toArchivePath(path string) string {
	return strings.ReplaceAll(path, string(filepath.Separator), "/")
}
