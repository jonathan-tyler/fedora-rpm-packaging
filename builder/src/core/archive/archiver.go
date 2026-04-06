package archive

type Entry struct {
	SourcePath  string
	ArchivePath string
}

type Archiver interface {
	CreateTarGz(destinationPath string, entries []Entry) error
}
