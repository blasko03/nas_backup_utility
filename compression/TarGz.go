package compression

import (
	"archive/tar"
	"compress/gzip"
	"io"
)

type TarGz struct {
	gzipWriter *gzip.Writer
	tarWriter  *tar.Writer
}

func NewTarGz(destination io.Writer, level int) *TarGz {
	gzipWriter, _ := gzip.NewWriterLevel(destination, level)
	tarWriter := tar.NewWriter(gzipWriter)
	return &TarGz{
		gzipWriter: gzipWriter,
		tarWriter:  tarWriter,
	}
}

func (t *TarGz) Write(data []byte) (int, error) {
	return t.tarWriter.Write(data)
}

func (t *TarGz) WriteHeader(header *tar.Header) error {
	return t.tarWriter.WriteHeader(header)
}

func (t *TarGz) Close() {
	t.tarWriter.Close()
	t.gzipWriter.Close()
}
