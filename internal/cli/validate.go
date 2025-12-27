package cli

import (
	"os"

	"github.com/obay/tfren/internal/config"
	"github.com/obay/tfren/internal/operations"
	"github.com/obay/tfren/internal/output"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Check files against naming convention",
	Long: `Validate that Terraform files follow OTN naming convention.
Reports violations without making any changes.`,
	RunE: runValidate,
}

func init() {
	validateCmd.Flags().Bool("strict", false, "exit with error code if violations found")
}

func runValidate(cmd *cobra.Command, args []string) error {
	cfg := config.Get()

	strict, _ := cmd.Flags().GetBool("strict")

	opts := operations.Options{
		Directory: cfg.Directory,
		Recursive: cfg.Recursive,
		Verbose:   cfg.Verbose,
		Strict:    strict,
	}

	result := output.NewResult(Version, "validate", cfg.Directory, false)

	var console *output.Console
	if !cfg.JSONOutput {
		console = output.NewConsole(cfg.Verbose, cfg.Quiet)
	}

	if err := operations.Validate(opts, result, console); err != nil {
		result.AddError(err.Error())
		result.ExitCode = 3
	}

	if cfg.JSONOutput {
		if err := result.Print(); err != nil {
			return err
		}
	} else if console != nil && !cfg.Quiet {
		printValidationSummary(result, console)
	}

	if result.ExitCode != 0 {
		os.Exit(result.ExitCode)
	}

	return nil
}

func printValidationSummary(result *output.Result, console *output.Console) {
	console.Print("")
	if len(result.Violations) == 0 {
		console.Success("All files are compliant with OTN naming convention.")
	} else {
		console.Warning("Found " + itoa(len(result.Violations)) + " naming violations.")
	}
	console.Print("Files checked: " + itoa(result.FilesProcessed))
	console.Print("Compliant: " + itoa(result.FilesAlreadyCompliant))
	console.Print("Violations: " + itoa(len(result.Violations)))
}
