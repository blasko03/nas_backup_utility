package main

import (
	"compress/gzip"
	"fmt"
	"os"
	"sync"
	"synchronizer/backup"
	"synchronizer/compression"
	"time"
)

type DestinationWriter interface {
	Write([]byte) (int, error)
	Close() error
}

func main() {
	var wg sync.WaitGroup
	files := make(chan []byte, 2)
	wg.Add(1)
	go StratCompression(&wg, files)

	wg.Add(1)
	go SaveData(&wg, files)

	wg.Wait()
	fmt.Println("finished")
}

func StratCompression(wg *sync.WaitGroup, files chan []byte) {
	defer wg.Done()

	config := backup.GetConfig()
	var dateFrom = time.Now().Add(-1 * time.Hour)

	changedFiles, err := backup.ChangedFiles(config.IncludedFolders[0], &dateFrom, config.ExcludedFolders)
	if err != nil {
		fmt.Println(err)
		return
	}

	tarGz := compression.NewTarGzChunked(files, gzip.BestCompression)

	addFile := compression.NewAddFileChunked(tarGz, config.ChunkSize)
	compression.Compress(changedFiles, addFile)
	tarGz.Close()
	close(files)
	fmt.Println("Archives crearted")
}

func SaveData(wg *sync.WaitGroup, files chan []byte) {
	//destination := destinations.NewHttpUploader()
	defer wg.Done()
	revieving := true
	for i := 0; revieving; i++ {
		archive, ok := <-files
		revieving = ok
		if !revieving {
			return
		}
		filename := fmt.Sprintf("/tmp/tarballFilePath.tar.gz %d", i)
		os.Remove(filename)
		destination, _ := os.Create(filename)
		n, err := destination.Write(archive)
		defer destination.Close()
		fmt.Println(n, err)
	}
}
