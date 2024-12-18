package main

import (
	"errors"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/hashicorp/hcl/v2/hclparse"
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
	content, err := os.ReadFile(file.Name())
	if err != nil {
		log.Fatalf("Failed to read file %s: %v", file.Name(), err)
	}

	parser := hclparse.NewParser()
	_, diags := parser.ParseHCL(content, file.Name())
	if diags.HasErrors() {
		log.Warnf("Failed to parse file %s: %v", file.Name(), diags)
		return
	}

	newFileName := generateNewFileNameFromHCL(string(content))
	if newFileName == "" {
		log.Warnf("Could not determine new name for %s", file.Name())
		return
	}

	if newFileName == file.Name() {
		PrintSuccess("Named correctly: " + file.Name() + ". Skipping...")
		return
	}

	handleFileRenaming(file.Name(), newFileName)
}

func generateNewFileNameFromHCL(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") || line == "" {
			continue
		}

		// Extract the first word to determine the block type
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "resource", "data":
			if len(parts) >= 3 {
				resourceType := strings.Trim(parts[1], "\"")
				resourceName := strings.Trim(parts[2], "\"")
				return parts[0] + "." + resourceType + "." + resourceName + ".tf"
			}
		case "variable", "module", "output":
			if len(parts) >= 2 {
				name := strings.Trim(parts[1], "\"")
				return parts[0] + "." + name + ".tf"
			}
		case "provider":
			if len(parts) >= 2 {
				name := strings.Trim(parts[1], "\"")
				return parts[0] + "." + name + ".tf"
			}
		case "terraform":
			return "terraform.tf"
		}
	}
	return ""
}

func handleFileRenaming(oldName, newName string) {
	if _, err := os.Stat(newName); err == nil {
		PrintError("A file with the name \"" + newName + "\" exists already. Unable to rename \"" + oldName + "\". Skipping...")
		return
	} else if errors.Is(err, os.ErrNotExist) {
		// Attempt renaming the file
		if err := os.Rename(oldName, newName); err != nil {
			log.Fatalf("Failed to rename file from %s to %s: %v", oldName, newName, err)
		}
		PrintWarning("Renamed file from " + oldName + " to " + newName)
	} else {
		log.Fatalf("Failed to check file existence: %v", err)
	}
}
