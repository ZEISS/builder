package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// ScanDir scans a directory for folders and files, ignoring specified paths.
func ScanDir(path string, ignore []string) []string {
	files := []string{}

	// Scan
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		for _, i := range ignore {
			if strings.Contains(path, i) {
				return nil
			}
		}

		if f.IsDir() {
			return nil
		}

		files = append(files, path)

		return nil
	})

	return files
}
