package operations

import (
	"path/filepath"

	"github.com/obay/tfren/internal/hcl"
	"github.com/obay/tfren/internal/naming"
	"github.com/obay/tfren/internal/output"
)

func Validate(opts Options, result *output.Result, console *output.Console) error {
	files, err := GetTerraformFiles(opts.Directory, opts.Recursive)
	if err != nil {
		return err
	}

	for _, file := range files {
		result.FilesProcessed++
		if err := validateFile(file, opts, result, console); err != nil {
			result.AddError(err.Error())
		}
	}

	if len(result.Violations) > 0 && opts.Strict {
		result.ExitCode = 2
	}

	return nil
}

func validateFile(path string, opts Options, result *output.Result, console *output.Console) error {
	parsed, err := hcl.ParseFile(path)
	if err != nil {
		return err
	}

	if parsed.ParseError != nil {
		result.AddError("Failed to parse: " + path)
		return nil
	}

	currentName := filepath.Base(path)

	if len(parsed.Blocks) == 0 {
		return nil
	}

	if len(parsed.Blocks) > 1 {
		result.AddViolation(path, currentName, "(should be split)", "multiple")
		if console != nil {
			console.Warning(currentName + " contains multiple blocks - should be split")
		}
		return nil
	}

	block := parsed.Blocks[0]
	expectedName := naming.GenerateFileName(block.Type, block.Labels, block.Alias)
	if expectedName == "" {
		return nil
	}

	if currentName != expectedName {
		result.AddViolation(path, currentName, expectedName, block.Type)
		if console != nil {
			console.Warning(currentName + " should be " + expectedName)
		}
	} else {
		result.AddCompliant(currentName)
		if console != nil && opts.Verbose {
			console.Success(currentName + " - OK")
		}
	}

	return nil
}
