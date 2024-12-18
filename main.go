package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func main() {
	versionFlag := flag.Bool("version", false, "Print version information")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("v%s\n", Version)
		return
	}

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
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(content), "temp.tf")
	if diags.HasErrors() {
		return ""
	}

	// Get the root body of the HCL file
	body := file.Body.(*hclsyntax.Body)

	// Look through all blocks in the file
	for _, block := range body.Blocks {
		switch block.Type {
		case "resource", "data":
			if len(block.Labels) >= 2 {
				return block.Type + "." + block.Labels[0] + "." + block.Labels[1] + ".tf"
			}
		case "provider":
			if len(block.Labels) >= 1 {
				providerName := block.Labels[0]
				// Look for alias attribute in the block
				for name, attr := range block.Body.Attributes {
					if name == "alias" {
						// Get the alias value
						if val, diags := attr.Expr.Value(nil); !diags.HasErrors() {
							alias := val.AsString()
							if alias != "" {
								return block.Type + "." + providerName + "." + alias + ".tf"
							}
						}
					}
				}
				return block.Type + "." + providerName + ".tf"
			}
		case "variable", "module", "output":
			if len(block.Labels) >= 1 {
				return block.Type + "." + block.Labels[0] + ".tf"
			}
		case "locals":
			return "locals.tf"
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
