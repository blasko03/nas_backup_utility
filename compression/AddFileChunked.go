package compression

import (
	"archive/tar"
	"os"
	"strconv"
)

func AddFileChunked(filePath string, tarWriter *tar.Writer, chunkSize int) error {
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
