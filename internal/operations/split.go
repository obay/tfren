package operations

import (
	"os"
	"path/filepath"

	"github.com/obay/tfren/internal/hcl"
	"github.com/obay/tfren/internal/naming"
	"github.com/obay/tfren/internal/output"
)

func Split(opts Options, result *output.Result, console *output.Console) error {
	files, err := GetTerraformFiles(opts.Directory, opts.Recursive)
	if err != nil {
		return err
	}

	for _, file := range files {
		result.FilesProcessed++
		if err := splitFile(file, opts, result, console); err != nil {
			result.AddError(err.Error())
		}
	}

	return nil
}

func splitFile(path string, opts Options, result *output.Result, console *output.Console) error {
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

	if len(parsed.Blocks) <= 1 {
		return nil
	}

	if opts.Backup && !opts.DryRun {
		if err := CreateBackup(path); err != nil {
			return err
		}
	}

	dir := filepath.Dir(path)
	var targets []string

	for _, block := range parsed.Blocks {
		newName := naming.GenerateFileName(block.Type, block.Labels, block.Alias)
		if newName == "" {
			continue
		}

		newPath := filepath.Join(dir, newName)
		targets = append(targets, newName)

		if FileExists(newPath) {
			result.AddSkip(newName, "file already exists")
			if console != nil {
				console.Error("Cannot create " + newName + " - file already exists")
			}
			continue
		}

		if !opts.DryRun {
			content := block.ContentWithComments()
			if err := os.WriteFile(newPath, []byte(content), 0644); err != nil {
				return err
			}
		}

		if console != nil {
			console.Success("Created: " + newName)
		}
	}

	if len(targets) > 0 {
		result.AddSplit(filepath.Base(path), targets)

		if !opts.DryRun && !opts.KeepOriginal {
			if err := os.Remove(path); err != nil {
				return err
			}
		}

		if console != nil {
			console.Success("Split " + filepath.Base(path) + " into " + string(rune(len(targets))) + " files")
		}
	}

	return nil
}
