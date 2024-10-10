package compression

import (
	"fmt"
)

type IAddFile interface {
	Write(string) ([]byte, error)
}

type CompressedFile struct {
	Path string
	Hash []byte
	err  error
}

func Compress(filePaths *[]string, addFile IAddFile) *[]CompressedFile {
	var compressedFiles []CompressedFile
	for _, filePath := range *filePaths {
		fmt.Println("Compressing " + filePath)
		hash, err := addFile.Write(filePath)
		compressedFiles = append(compressedFiles, CompressedFile{Path: filePath, Hash: hash, err: err})
	}
	return &compressedFiles
}
