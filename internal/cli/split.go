package cli

import (
	"github.com/obay/tfren/internal/config"
	"github.com/obay/tfren/internal/operations"
	"github.com/obay/tfren/internal/output"
	"github.com/spf13/cobra"
)

var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "Split multi-block files into individual files",
	Long: `Split Terraform files containing multiple blocks into separate files,
one block per file, named according to OTN convention.`,
	RunE: runSplit,
}

func init() {
	splitCmd.Flags().Bool("keep-original", false, "keep the original file after splitting")
	splitCmd.Flags().Bool("backup", false, "create .bak backup before splitting")
}

func runSplit(cmd *cobra.Command, args []string) error {
	cfg := config.Get()

	keepOriginal, _ := cmd.Flags().GetBool("keep-original")
	backup, _ := cmd.Flags().GetBool("backup")

	opts := operations.Options{
		Directory:    cfg.Directory,
		Recursive:    cfg.Recursive,
		DryRun:       cfg.DryRun,
		Verbose:      cfg.Verbose,
		KeepOriginal: keepOriginal,
		Backup:       backup,
	}

	result := output.NewResult(Version, "split", cfg.Directory, cfg.DryRun)

	var console *output.Console
	if !cfg.JSONOutput {
		console = output.NewConsole(cfg.Verbose, cfg.Quiet)
	}

	if err := operations.Split(opts, result, console); err != nil {
		result.AddError(err.Error())
		result.ExitCode = 3
	}

	if cfg.JSONOutput {
		return result.Print()
	}

	if console != nil && !cfg.Quiet {
		printSummary(result, console)
	}

	return nil
}
