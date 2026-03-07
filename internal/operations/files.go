package operations

import (
	"os"
	"path/filepath"
	"strings"
)

func GetTerraformFiles(dir string, recursive bool) ([]string, error) {
	var files []string

	if recursive {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".tf") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tf") {
				files = append(files, filepath.Join(dir, entry.Name()))
			}
		}
	}

	return files, nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func CreateBackup(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return os.WriteFile(path+".bak", content, 0644)
}
