package destinations

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
)

type HttpUploader struct {
	buffer    bytes.Buffer
	UploadId  string
	LastChunk int
	lock      sync.Mutex
}

func NewHttpUploader() *HttpUploader {
	return &HttpUploader{
		LastChunk: 0,
		UploadId:  "randomid",
	}
}

func (upload *HttpUploader) Write(data []byte) (int, error) {
	upload.lock.Lock()
	defer upload.lock.Unlock()

	upload.buffer.Write(data)
	url := fmt.Sprintf("http://localhost:3000/fileUpload/%s/%d", upload.UploadId, upload.LastChunk)
	resp, err := http.Post(url, "application/tar+gzip", &upload.buffer)
	fmt.Println(url)
	fmt.Println(resp)
	upload.LastChunk++
	return upload.buffer.Len(), err
}

func (server *HttpUploader) Close() error {
	return nil
}
