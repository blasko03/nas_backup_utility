package compression

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"strconv"
)

type AddFileChunked struct {
	archive   *TarGz
	chunkSize int
}

func NewAddFileChunked(tarWriter *TarGz, chunkSize int) *AddFileChunked {
	return &AddFileChunked{
		archive:   tarWriter,
		chunkSize: chunkSize,
	}
}

func (t *AddFileChunked) Write(filePath string) error {
	file, err := os.Open(filePath)

	if err != nil {
		return err
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	reading := true

	for i := 0; reading; i++ {
		bytes := make([]byte, t.chunkSize)
		n, err := file.Read(bytes)

		if err != nil {
			if errors.Is(err, io.EOF) {
				reading = false
				return nil
			} else {
				return err
			}
		}

		header := &tar.Header{
			Name:    filePath + ".bck-chunk-" + strconv.Itoa(i),
			Size:    int64(n),
			Mode:    int64(stat.Mode()),
			ModTime: stat.ModTime(),
		}
		err = t.archive.WriteHeader(header)

		if err != nil {
			return err
		}

		t.archive.Write(bytes[0:n])
	}

	return nil
}
