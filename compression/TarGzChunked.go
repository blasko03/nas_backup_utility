package compression

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
)

type TarGzChunked struct {
	gzipWriter *gzip.Writer
	tarWriter  *tar.Writer
	files      chan []byte
	level      int
	buffer     bytes.Buffer
	counter    int
	bytes      int64
}

func NewTarGzChunked(files chan []byte, level int) *TarGzChunked {
	t := &TarGzChunked{
		files:   files,
		level:   level,
		counter: 0,
		bytes:   0,
	}

	return t
}

func (t *TarGzChunked) AddFile(header *tar.Header, data []byte) (int, error) {
	if t.counter%10 == 0 {
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
	t.counter++
	t.bytes += int64(n)
	return n, nil
}

func (t *TarGzChunked) newArchive() error {
	gzipWriter, err := gzip.NewWriterLevel(&t.buffer, t.level)
	tarWriter := tar.NewWriter(gzipWriter)
	t.gzipWriter = gzipWriter
	t.tarWriter = tarWriter
	return err
}

func (t *TarGzChunked) Close() {
	if t.tarWriter != nil {
		t.tarWriter.Close()
		t.gzipWriter.Close()
		t.files <- t.buffer.Bytes()
	}
}
