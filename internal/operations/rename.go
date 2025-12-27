package operations

import (
	"os"
	"path/filepath"

	"github.com/obay/tfren/internal/hcl"
	"github.com/obay/tfren/internal/naming"
	"github.com/obay/tfren/internal/output"
)

func Rename(opts Options, result *output.Result, console *output.Console) error {
	files, err := GetTerraformFiles(opts.Directory, opts.Recursive)
	if err != nil {
		return err
	}

	for _, file := range files {
		result.FilesProcessed++
		if err := renameFile(file, opts, result, console); err != nil {
			result.AddError(err.Error())
		}
	}

	return nil
}

func renameFile(path string, opts Options, result *output.Result, console *output.Console) error {
	parsed, err := hcl.ParseFile(path)
	if err != nil {
		return err
	}

	if parsed.ParseError != nil {
		result.AddSkip(filepath.Base(path), "parse error")
		if console != nil {
			console.Warning("Skipping " + path + " - parse error")
		}
		return nil
	}

	if len(parsed.Blocks) == 0 {
		result.AddSkip(filepath.Base(path), "no blocks found")
		return nil
	}

	if len(parsed.Blocks) > 1 {
		result.AddSkip(filepath.Base(path), "contains multiple blocks")
		if console != nil {
			console.Warning("Skipping " + path + " - contains multiple blocks (use 'split' first)")
		}
		return nil
	}

	block := parsed.Blocks[0]
	newName := naming.GenerateFileName(block.Type, block.Labels, block.Alias)
	if newName == "" {
		result.AddSkip(filepath.Base(path), "cannot determine filename")
		return nil
	}

	currentName := filepath.Base(path)
	if currentName == newName {
		result.AddCompliant(currentName)
		if console != nil {
			console.Success("Already compliant: " + currentName)
		}
		return nil
	}

	dir := filepath.Dir(path)
	newPath := filepath.Join(dir, newName)

	if FileExists(newPath) {
		result.AddSkip(currentName, "target file already exists: "+newName)
		if console != nil {
			console.Error("Cannot rename " + currentName + " to " + newName + " - file already exists")
		}
		return nil
	}

	if opts.Backup && !opts.DryRun {
		if err := CreateBackup(path); err != nil {
			return err
		}
	}

	if !opts.DryRun {
		if err := os.Rename(path, newPath); err != nil {
			return err
		}
	}

	result.AddRename(currentName, newName)
	if console != nil {
		console.Warning("Renamed: " + currentName + " -> " + newName)
	}

	return nil
}
