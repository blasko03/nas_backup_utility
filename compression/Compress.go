package compression

import (
	"fmt"
)

type IAddFile interface {
	Write(string) error
}

func Compress(filePaths []string, addFile IAddFile) []error {
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
