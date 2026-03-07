package cli

import (
	"os"

	"github.com/obay/tfren/internal/config"
	"github.com/obay/tfren/internal/operations"
	"github.com/obay/tfren/internal/output"
	"github.com/spf13/cobra"
)

var azureNamingCmd = &cobra.Command{
	Use:   "azure-naming",
	Short: "Check Azure resource names against CAF naming convention",
	Long: `Validate that azurerm_* resource blocks have name attributes following
Microsoft's Cloud Adoption Framework (CAF) naming convention.

Checks prefix, structure, length, valid characters, and naming rules
for each Azure resource type.`,
	RunE: runAzureNaming,
}

func init() {
	azureNamingCmd.Flags().Bool("strict", false, "exit with error code if violations found")
}

func runAzureNaming(cmd *cobra.Command, args []string) error {
	cfg := config.Get()

	strict, _ := cmd.Flags().GetBool("strict")

	opts := operations.Options{
		Directory: cfg.Directory,
		Recursive: cfg.Recursive,
		Verbose:   cfg.Verbose,
		Strict:    strict,
	}

	result := output.NewResult(Version, "azure-naming", cfg.Directory, false)

	var console *output.Console
	if !cfg.JSONOutput {
		console = output.NewConsole(cfg.Verbose, cfg.Quiet)
	}

	if console != nil && !cfg.Quiet {
		console.Print("Checking Azure CAF naming convention...\n")
	}

	if err := operations.AzureNaming(opts, result, console); err != nil {
		result.AddError(err.Error())
		result.ExitCode = 3
	}

	if cfg.JSONOutput {
		if err := result.Print(); err != nil {
			return err
		}
	} else if console != nil && !cfg.Quiet {
		printAzureNamingSummary(result, console)
	}

	if result.ExitCode != 0 {
		os.Exit(result.ExitCode)
	}

	return nil
}

func printAzureNamingSummary(result *output.Result, console *output.Console) {
	console.Print("")
	console.Print("Azure CAF Naming Summary:")
	console.Print("  Resources checked: " + itoa(result.ResourcesChecked))
	console.Print("  Compliant: " + itoa(result.ResourcesCompliant))
	console.Print("  Violations: " + itoa(len(result.NamingViolations)))
	if result.ResourcesSkipped > 0 {
		console.Print("  Skipped (dynamic): " + itoa(result.ResourcesSkipped))
	}
	if result.ResourcesUnknown > 0 {
		console.Print("  Unknown (no CAF rule): " + itoa(result.ResourcesUnknown))
	}

	if len(result.NamingViolations) == 0 && result.ResourcesChecked > 0 {
		console.Success("\nAll Azure resource names are compliant with CAF naming convention.")
	}
}
