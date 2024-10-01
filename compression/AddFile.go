package compression

import (
	"archive/tar"
	"io"
	"os"
)

func AddFile(filePath string, tarWriter *tar.Writer) error {
	file, err := os.Open(filePath)

	if err != nil {
		return err
	}

	defer file.Close()

	if err != nil {
		return err
	}

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

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return err
	}

	return nil
}
