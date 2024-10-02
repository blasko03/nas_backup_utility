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
	wg.Add(1)
	go StratCompression(&wg, files)

	wg.Add(1)
	go SaveData(&wg, files)

	wg.Wait()
	fmt.Println("finished")
}

func StratCompression(wg *sync.WaitGroup, files chan *bytes.Buffer) {
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
	errors := compression.Compress(changedFiles, addFile)

	fmt.Println(errors)
	tarGz.Close()
	close(files)
	fmt.Println("Archives crearted")
}

func SaveData(wg *sync.WaitGroup, files chan *bytes.Buffer) {
	destination := destinations.NewHttpUploader()
	defer wg.Done()
	revieving := true
	for i := 0; revieving; i++ {
		archive, ok := <-files
		revieving = ok
		if !revieving {
			return
		}

		destination.Save(archive, i)
	}
}
