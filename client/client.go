package main

import (
	"compress/gzip"
	"fmt"
	"sync"
	"synchronizer/backup"
	"synchronizer/compression"
	"synchronizer/destinations"
	"time"
)

type DestinationWriter interface {
	Write([]byte) (int, error)
	Close() error
}

func main() {
	config := backup.GetConfig()
	buffer := backup.NewLimitedBuffer(config.BufferSize)
	//os.Remove("/tmp/tarballFilePath.tar.gz")
	//destination, _ := os.Create("/tmp/tarballFilePath.tar.gz")
	destination := destinations.NewHttpUploader()

	var wg sync.WaitGroup

	wg.Add(1)
	go StartCompression(&wg, buffer)

	wg.Add(1)
	go SaveData(&wg, buffer, destination)

	wg.Wait()
	fmt.Println("finished")
}

func StartCompression(wg *sync.WaitGroup, buffer *backup.LimitedBuffer) {
	config := backup.GetConfig()
	var dateFrom = time.Now().Add(-1 * time.Hour)
	files, err := backup.ChangedFiles(config.IncludedFolders[0], &dateFrom, config.ExcludedFolders)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer wg.Done()
	tarGz := compression.NewTarGz(buffer, gzip.NoCompression)
	defer tarGz.Close()
	addFile := compression.NewAddFileChunked(tarGz, config.ChunkSize)
	compression.Compress(files, addFile)
	buffer.Close()
	fmt.Println("backup.Compress completed")
}

func SaveData(wg *sync.WaitGroup, buffer *backup.LimitedBuffer, destination DestinationWriter) {
	defer wg.Done()

	for !(buffer.Len() == 0 && buffer.IsClosed()) {
		res := make([]byte, buffer.Len())
		n, _ := buffer.Read(res)
		destination.Write(res[0:n])
	}
	defer destination.Close()
}
