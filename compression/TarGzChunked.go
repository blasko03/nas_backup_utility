package compression

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
)

type TarGzChunked struct {
	gzipWriter     *gzip.Writer
	tarWriter      *tar.Writer
	files          chan *bytes.Buffer
	level          int
	buffer         *bytes.Buffer
	archiveMaxSize int
}

func NewTarGzChunked(files chan *bytes.Buffer, level int, archiveMaxSize int) *TarGzChunked {
	t := &TarGzChunked{
		files:          files,
		level:          level,
		archiveMaxSize: archiveMaxSize,
	}

	return t
}

func (t *TarGzChunked) AddFile(header *tar.Header, data []byte) (int, error) {
	if t.buffer == nil || t.buffer.Len() >= t.archiveMaxSize {
		t.Close()
		t.newArchive()
	}
	err := t.tarWriter.WriteHeader(header)
	if err != nil {
		return 0, err
	}
	n, err := t.tarWriter.Write(data)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (t *TarGzChunked) newArchive() error {
	t.buffer = new(bytes.Buffer)
	gzipWriter, err := gzip.NewWriterLevel(t.buffer, t.level)
	tarWriter := tar.NewWriter(gzipWriter)
	t.gzipWriter = gzipWriter
	t.tarWriter = tarWriter
	return err
}

func (t *TarGzChunked) Close() {
	if t.tarWriter != nil {
		t.tarWriter.Close()
		t.gzipWriter.Close()
		t.files <- t.buffer
	}
}
