package backup

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ChangedFiles(path string, dateFrom *time.Time, excluded []string) ([]string, error) {
	var changedFiles []string

	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if containsAny(path, excluded[:]) {
			return nil
		}

		if dateFrom != nil && info.ModTime().Before(*dateFrom) {
			return nil
		}

		changedFiles = append(changedFiles, path)
		return nil
	})

	return changedFiles, nil
}

func containsAny(path string, arr []string) bool {
	for _, str := range arr {
		if strings.HasPrefix(path, str) {
			return true
		}
	}
	return false
}
