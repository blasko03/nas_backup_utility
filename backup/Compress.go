package backup

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func Compress(filePaths []string, destination io.Writer, level int) []error {
	var e []error

	gzipWriter, _ := gzip.NewWriterLevel(destination, level)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	for _, filePath := range filePaths {
		fmt.Println("Compressing " + filePath)
		err := addFile(filePath, tarWriter)
		if err != nil {
			e = append(e, err)
		}
	}
	return e
}

func addFile(filePath string, tarWriter *tar.Writer) error {
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
