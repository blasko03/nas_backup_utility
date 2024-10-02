package compression

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"strconv"
)

type AddFileChunked struct {
	archive   ITarGz
	chunkSize int
}

type ITarGz interface {
	AddFile(header *tar.Header, data []byte) (int, error)
	Close()
}

func NewAddFileChunked(tarWriter ITarGz, chunkSize int) *AddFileChunked {
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
	bytes := make([]byte, t.chunkSize)
	for i := 0; reading; i++ {

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
		_, err = t.archive.AddFile(header, bytes[0:n])

		if err != nil {
			return err
		}
	}

	return nil
}
