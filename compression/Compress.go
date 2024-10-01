package compression

import (
	"fmt"
)

func Compress(filePaths []string, addFile *AddFileChunked) []error {
	var e []error

	for _, filePath := range filePaths {
		fmt.Println("Compressing " + filePath)
		err := addFile.Write(filePath)
		if err != nil {
			e = append(e, err)
		}
	}
	return e
}
