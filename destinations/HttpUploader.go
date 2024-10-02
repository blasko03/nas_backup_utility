package destinations

import (
	"bytes"
	"fmt"
	"net/http"
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

func (server *HttpUploader) Close() error {
	return nil
}
