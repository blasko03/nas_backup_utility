package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"sync"
	"synchronizer/backup"
	"synchronizer/compression"
	"synchronizer/destinations"
	"time"
)

func main() {
	var wg sync.WaitGroup
	files := make(chan *bytes.Buffer, 1)
	var compressedFiles *[]compression.CompressedFile
	destination := destinations.NewHttpUploader()

	wg.Add(1)
	go StratCompression(&wg, files, &compressedFiles)

	wg.Add(1)
	//go SaveData(&wg, destination, files)

	wg.Wait()
	destination.Close(compressedFiles)
	fmt.Println(compressedFiles)
	fmt.Println("finished")
}

func StratCompression(wg *sync.WaitGroup, files chan *bytes.Buffer, compressedFiles **[]compression.CompressedFile) {
	defer wg.Done()

	config := backup.GetConfig()
	var dateFrom = time.Now().Add(-10000 * time.Hour)

	changedFiles, err := backup.ChangedFiles(config.IncludedFolders[0], &dateFrom, config.ExcludedFolders)
	if err != nil {
		fmt.Println(err)
		return
	}

	tarGz := compression.NewTarGzChunked(files, gzip.NoCompression, config.ArchiveMaxSize)
	addFile := compression.NewAddFileChunked(tarGz, config.ChunkSize)
	*compressedFiles = compression.Compress(&changedFiles, addFile)

	tarGz.Close()
	close(files)
	fmt.Println("Archives crearted")
}

func SaveData(wg *sync.WaitGroup, destination *destinations.HttpUploader, files chan *bytes.Buffer) {
	defer wg.Done()
	recieving := true
	for i := 0; recieving; i++ {
		_, ok := <-files
		recieving = ok
		if !recieving {
			return
		}
		time.Sleep(time.Millisecond * 10000)
		//destination.Save(archive, i)
	}
}
