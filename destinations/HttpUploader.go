package destinations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"synchronizer/compression"
)

type HttpUploader struct {
	UploadId string
}

func NewHttpUploader() *HttpUploader {
	return &HttpUploader{
		UploadId: "test",
	}
}

func (upload *HttpUploader) Save(data *bytes.Buffer, chunk int) (int, error) {
	url := fmt.Sprintf("http://localhost:3000/fileUpload/%s/%d", upload.UploadId, chunk)
	resp, err := http.Post(url, "application/tar+gzip", data)
	fmt.Println(url)
	fmt.Println(resp)
	return data.Len(), err
}

func (upload *HttpUploader) Close(compressedFiles *[]compression.CompressedFile) error {
	url := fmt.Sprintf("http://localhost:3000/fileUploadCompleted/%s", upload.UploadId)
	buff := &bytes.Buffer{}
	enc := json.NewEncoder(buff)
	err := enc.Encode(*compressedFiles)
	if err != nil {
		return err
	}

	_, errHttp := http.Post(url, "application/json", buff)

	return errHttp
}
