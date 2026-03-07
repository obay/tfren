package operations

import (
	"os"
	"path/filepath"
	"strconv"

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
	sourceName := filepath.Base(path)
	var targets []string
	var keepInPlace *hcl.Block // Block that should stay in the source file

	// First pass: collect targets and check for conflicts
	for i, block := range parsed.Blocks {
		newName := naming.GenerateFileName(block.Type, block.Labels, block.Alias)
		if newName == "" {
			continue
		}

		// If target name matches source, this block stays in place
		if newName == sourceName {
			keepInPlace = &parsed.Blocks[i]
			targets = append(targets, newName)
			continue
		}

		newPath := filepath.Join(dir, newName)
		if FileExists(newPath) {
			result.AddSkip(newName, "file already exists")
			if console != nil {
				console.Error("Cannot create " + newName + " - file already exists")
			}
			continue
		}
		targets = append(targets, newName)
	}

	if len(targets) == 0 {
		return nil
	}

	// Print header before creating files
	if console != nil {
		console.Success("Splitting " + filepath.Base(path) + " → " + strconv.Itoa(len(targets)) + " files:")
	}

	// Second pass: create new files (skip the block that stays in place)
	for _, block := range parsed.Blocks {
		newName := naming.GenerateFileName(block.Type, block.Labels, block.Alias)
		if newName == "" {
			continue
		}

		// Skip the block that stays in place - it will be handled at the end
		if newName == sourceName {
			if console != nil {
				console.Print("  → " + newName + " (rewritten)")
			}
			continue
		}

		newPath := filepath.Join(dir, newName)
		if FileExists(newPath) {
			continue // Already reported in first pass
		}

		if !opts.DryRun {
			content := block.ContentWithComments()
			if err := os.WriteFile(newPath, []byte(content), 0644); err != nil {
				return err
			}
		}

		if console != nil {
			console.Print("  → " + newName)
		}
	}

	result.AddSplit(filepath.Base(path), targets)

	// Handle the source file
	if !opts.DryRun && !opts.KeepOriginal {
		if keepInPlace != nil {
			// Rewrite source file with only the block that matches its name
			content := keepInPlace.ContentWithComments()
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				return err
			}
		} else {
			// No block matches source name, delete the source file
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}

	return nil
}
