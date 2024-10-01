package backup

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strconv"
)

func CompressChunked(filePaths []string, destination io.Writer, level int, chunkSize int) []error {
	var e []error

	gzipWriter, _ := gzip.NewWriterLevel(destination, level)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	for _, filePath := range filePaths {
		fmt.Println("Compressing " + filePath)
		err := addFileChunked(filePath, tarWriter, chunkSize)
		if err != nil {
			e = append(e, err)
		}
	}
	return e
}

func addFileChunked(filePath string, tarWriter *tar.Writer, chunkSize int) error {
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

	reading := true

	for i := 0; reading; i++ {
		bytes := make([]byte, chunkSize)
		n, err := file.Read(bytes)

		if err != nil {
			return err
		}

		header := &tar.Header{
			Name:    filePath + ".bck-chunk-" + strconv.Itoa(i),
			Size:    int64(n),
			Mode:    int64(stat.Mode()),
			ModTime: stat.ModTime(),
		}

		err = tarWriter.WriteHeader(header)

		if err != nil {
			return err
		}
		tarWriter.Write(bytes[0:n])
		reading = n > 0
	}

	return nil
}
