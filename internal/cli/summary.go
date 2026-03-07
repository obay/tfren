package cli

import (
	"strconv"

	"github.com/obay/tfren/internal/output"
)

func printSummary(result *output.Result, console *output.Console) {
	console.Print("")
	console.Print("Summary:")
	console.Print("  Files processed: " + itoa(result.FilesProcessed))
	if result.FilesSplit > 0 {
		console.Print("  Files split: " + itoa(result.FilesSplit))
	}
	if result.FilesRenamed > 0 {
		console.Print("  Files renamed: " + itoa(result.FilesRenamed))
	}
	if result.FilesSkipped > 0 {
		console.Print("  Files skipped: " + itoa(result.FilesSkipped))
	}
	if result.FilesAlreadyCompliant > 0 {
		console.Print("  Already compliant: " + itoa(result.FilesAlreadyCompliant))
	}
	if len(result.Errors) > 0 {
		console.Print("  Errors: " + itoa(len(result.Errors)))
	}

	if result.DryRun {
		console.Warning("\nDry run - no changes were made.")
	}
}

func itoa(i int) string {
	return strconv.Itoa(i)
}
