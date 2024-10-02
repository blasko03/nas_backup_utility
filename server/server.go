package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func fileUpload(w http.ResponseWriter, r *http.Request) {
	uploadId := r.PathValue("uploadId")
	part := r.PathValue("part")
	fileName := fmt.Sprintf("/tmp/%s/tarballFilePath%s.tar.gz", uploadId, part)
	if part == "0" {
		fmt.Println("remove file")
		os.Remove(fileName)
	}

	out, _ := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	defer out.Close()
	defer r.Body.Close()

	data, _ := io.ReadAll(r.Body)

	out.Write(data)
}

func main() {
	http.HandleFunc("/fileUpload/{uploadId}/{part}", fileUpload)
	http.HandleFunc("/fileUploadCompleted", fileUpload)
	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		fmt.Println(err)
	}
}
