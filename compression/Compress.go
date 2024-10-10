package compression

import (
	"encoding/hex"
	"fmt"
)

type IAddFile interface {
	Write(string) ([]byte, error)
}

type CompressedFile struct {
	path string
	hash []byte
	err  error
}

func Compress(filePaths []string, addFile IAddFile) *[]CompressedFile {
	var compressedFiles []CompressedFile
	for _, filePath := range filePaths {
		fmt.Println("Compressing " + filePath)
		hash, err := addFile.Write(filePath)
		fmt.Println(hex.EncodeToString(hash))
		compressedFiles = append(compressedFiles, CompressedFile{path: filePath, hash: hash, err: err})
	}
	return &compressedFiles
}
