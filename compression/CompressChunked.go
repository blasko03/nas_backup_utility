package compression

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
)

type AddFileFunc func(string, *tar.Writer, int) error

func Compress(filePaths []string, destination io.Writer, level int, chunkSize int, addFile AddFileFunc) []error {
	var e []error

	gzipWriter, _ := gzip.NewWriterLevel(destination, level)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	for _, filePath := range filePaths {
		fmt.Println("Compressing " + filePath)
		err := addFile(filePath, tarWriter, chunkSize)
		if err != nil {
			e = append(e, err)
		}
	}
	return e
}
