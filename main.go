package main

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

func main() {
	files := getCurrentDirectoryFiles()
	for _, file := range files {
		if isValidFile(file) {
			renameFileBasedOnContent(file)
		}
	}
}

func getCurrentDirectoryFiles() []os.FileInfo {
	currentDirectory, err := os.Open(".")
	if err != nil {
		log.Fatal(err)
	}
	defer currentDirectory.Close()

	files, err := currentDirectory.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	return files
}

func isValidFile(file os.FileInfo) bool {
	return strings.HasSuffix(file.Name(), ".tf") && !file.IsDir()
}

func renameFileBasedOnContent(file os.FileInfo) {
	f, err := os.Open(file.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// If the line is a comment, skip it
		if strings.HasPrefix(line, "#") {
			continue
		}

		prefixes := []string{"resource", "provider", "variable", "terraform", "module", "data"} // Add any other prefixes to this slice
		matches := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(line, prefix) {
				matches = true
				break
			}
		}

		if matches {
			newFileName := generateNewFileName(line)
			if newFileName == file.Name() {
				PrintSuccess("Named correctly: " + file.Name() + ". Skipping...")
				break
			}
			handleFileRenaming(file.Name(), newFileName)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func generateNewFileName(line string) string {
	parts := strings.Split(line, " ")
	parts = removeQuotesFromSlice(parts)
	if strings.HasPrefix(line, "resource") && len(parts) == 4 {
		return parts[0] + "." + parts[1] + "." + parts[2] + ".tf"
	}
	if strings.HasPrefix(line, "data") && len(parts) == 4 {
		return parts[0] + "." + parts[1] + "." + parts[2] + ".tf"
	}
	if strings.HasPrefix(line, "provider") && len(parts) == 3 {
		return parts[0] + "." + parts[1] + ".tf"
	}
	if strings.HasPrefix(line, "variable") && len(parts) == 3 {
		return parts[0] + "." + parts[1] + ".tf"
	}
	if strings.HasPrefix(line, "module") && len(parts) == 3 {
		return parts[0] + "." + parts[1] + ".tf"
	}
	if strings.HasPrefix(line, "terraform") && len(parts) == 2 {
		return parts[0] + ".tf"
	}
	return ""
}

func removeQuotesFromSlice(slice []string) []string {
	for i, s := range slice {
		slice[i] = strings.Replace(s, "\"", "", -1)
	}
	return slice
}

func handleFileRenaming(oldName, newName string) {
	if _, err := os.Stat(newName); err == nil {
		PrintError("A file with the name \"" + newName + "\" exists already. Unable to rename \"" + oldName + "\". Skipping...")
		return
	} else if errors.Is(err, os.ErrNotExist) {
		e := os.Rename(oldName, newName)
		if e != nil {
			log.Fatal(e)
		}
		PrintWarning("Renamed file from " + oldName + " to " + newName)
	} else {
		log.Fatal(err)
	}
}
