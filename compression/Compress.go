package compression

import (
	"encoding/hex"
	"fmt"
)

type IAddFile interface {
	Write(string) ([]byte, error)
}

func Compress(filePaths []string, addFile IAddFile) []error {
	var e []error

	for _, filePath := range filePaths {
		fmt.Println("Compressing " + filePath)
		hash, err := addFile.Write(filePath)
		fmt.Println(hex.EncodeToString(hash))
		if err != nil {
			e = append(e, err)
		}
	}
	return e
}
