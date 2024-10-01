package backup

type Config struct {
	ChunkSize       int
	BufferSize      int
	IncludedFolders []string
	ExcludedFolders []string
}

func GetConfig() Config {
	return Config{
		ChunkSize:       10 * 1024 * 1024,
		BufferSize:      10 * 1024 * 1024,
		IncludedFolders: []string{"/home/daniel/apps"},
		ExcludedFolders: []string{},
	}
}
