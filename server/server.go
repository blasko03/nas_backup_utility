package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func fileUpload(w http.ResponseWriter, r *http.Request) {
	uploadId := r.PathValue("uploadId")
	chunk := r.PathValue("chunk")

	fileName := fmt.Sprintf("/tmp/%s/backup-%s.tar.gz", uploadId, chunk)
	file, _ := os.Create(fileName)
	io.Copy(file, r.Body)

	defer file.Close()
	defer r.Body.Close()
}

func fileUploadCompleted(w http.ResponseWriter, r *http.Request) {
	uploadId := r.PathValue("uploadId")
	checksum := sha256.New()
	for i := 0; i <= 3; i++ {
		fileName := fmt.Sprintf("/tmp/%s/backup-%s.tar.gz", uploadId, strconv.Itoa(i))
		archive, _ := os.Open(fileName)

		gzipStream, err := gzip.NewReader(archive)
		if err != nil {
			log.Fatal("failed")
		}

		tarReader := tar.NewReader(gzipStream)

		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			fmt.Println(header)
			io.Copy(checksum, tarReader)
		}
	}
	fmt.Println(hex.EncodeToString(checksum.Sum(nil)))

}

func inventory(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/fileUpload/{uploadId}/{chunk}", fileUpload)
	http.HandleFunc("/fileUploadCompleted/{uploadId}", fileUploadCompleted)
	http.HandleFunc("/inventory", inventory)
	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		fmt.Println(err)
	}
}
