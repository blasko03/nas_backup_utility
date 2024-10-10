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
	files := make(chan *bytes.Buffer)
	var compressedFiles *[]compression.CompressedFile
	wg.Add(1)
	go StratCompression(&wg, files, &compressedFiles)

	wg.Add(1)
	go SaveData(&wg, files)

	wg.Wait()
	fmt.Println(compressedFiles)
	fmt.Println("finished")
}

func StratCompression(wg *sync.WaitGroup, files chan *bytes.Buffer, pointer **[]compression.CompressedFile) {
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
	*pointer = compression.Compress(changedFiles, addFile)

	tarGz.Close()
	close(files)
	fmt.Println("Archives crearted")
}

func SaveData(wg *sync.WaitGroup, files chan *bytes.Buffer) {
	destination := destinations.NewHttpUploader()
	defer wg.Done()
	recieving := true
	for i := 0; recieving; i++ {
		archive, ok := <-files
		recieving = ok
		if !recieving {
			return
		}

		destination.Save(archive, i)
	}
}

func ConfirmBackup(compressedFiles *[]compression.CompressedFile) {

}
