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

func main() {
	config := backup.GetConfig()
	var dateFrom = time.Now().Add(-1 * time.Hour)
	files, err := backup.ChangedFiles(config.IncludedFolders[0], &dateFrom, config.ExcludedFolders)

	if err != nil {
		fmt.Println(err)
		return
	}

	bytes := backup.NewLimitedBuffer(config.BufferSize)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		compression.Compress(files, bytes, gzip.NoCompression, config.ChunkSize, compression.AddFileChunked)
		fmt.Println("backup.Compress completed")
		bytes.Completed()
	}()
	wg.Add(1)
	largest := 0
	smallest := 100000000
	i := 0
	var total int64 = 0
	go func() {
		defer wg.Done()
		//os.Remove("/tmp/tarballFilePath.tar.gz")
		//file, _ := os.Create("/tmp/tarballFilePath.tar.gz")
		destination := destinations.NewHttpUploader()
		for !(bytes.Len() == 0 && bytes.IsCompleted()) {
			res := make([]byte, bytes.Len())
			n, _ := bytes.Read(res)
			largest = max(largest, n)
			if n > 0 {
				smallest = min(smallest, n)
			}
			total += int64(n)
			i++
			destination.Write(res[0:n])
			//file.Write(res[0:n])
		}
		defer destination.Close()
		fmt.Println(smallest, largest, total, i, total/int64(i))
	}()

	wg.Wait()
	fmt.Println("finished")
}
