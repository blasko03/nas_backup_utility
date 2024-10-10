package backup

type Config struct {
	ChunkSize       int
	BufferSize      int
	ArchiveMaxSize  int
	IncludedFolders []string
	ExcludedFolders []string
}

func GetConfig() Config {
	return Config{
		ChunkSize:       2 * 1024 * 1024,
		ArchiveMaxSize:  4 * 1024 * 1024,
		IncludedFolders: []string{"/home/daniel/Downloads/"},
		ExcludedFolders: []string{},
	}
}
