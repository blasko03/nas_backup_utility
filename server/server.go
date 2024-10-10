package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"synchronizer/compression"
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
	var compressedFiles []compression.CompressedFile
	dec := json.NewDecoder(r.Body)
	dec.Decode(&compressedFiles)
	fmt.Println(compressedFiles)
	baseDir := path.Join("/", "tmp", uploadId)
	archives, _ := os.ReadDir(baseDir)
	hashes := make(map[string]hash.Hash)
	for _, archive := range archives {
		compressed, err := os.Open(path.Join(baseDir, archive.Name()))
		if err != nil {
			log.Fatal(err)
			return
		}
		gzipStream, err := gzip.NewReader(compressed)
		if err != nil {
			log.Fatal(err)
		}

		tarReader := tar.NewReader(gzipStream)

		for {
			header, err := tarReader.Next()

			if err == io.EOF {
				break
			}
			filename := strings.Split(header.Name, ".bck-chunk-")
			if hashes[filename[0]] == nil {
				hashes[filename[0]] = sha256.New()
			}
			io.Copy(hashes[filename[0]], tarReader)
		}
	}
	for key, file := range hashes {
		fmt.Println(key, hex.EncodeToString(file.Sum(nil)))
	}
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
