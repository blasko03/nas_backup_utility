package compression

import (
	"archive/tar"
	"io"
	"os"
)

type AddFile struct {
	archive *TarGz
}

func NewAddFile(tarWriter *TarGz) *AddFile {
	return &AddFile{
		archive: tarWriter,
	}
}

func (t *AddFile) Write(filePath string) error {
	file, err := os.Open(filePath)

	if err != nil {
		return err
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name:    filePath,
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}

	err = t.archive.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(t.archive, file)
	if err != nil {
		return err
	}

	return nil
}
