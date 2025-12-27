package cli

import (
	"github.com/obay/tfren/internal/config"
	"github.com/obay/tfren/internal/operations"
	"github.com/obay/tfren/internal/output"
	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename files to match naming convention",
	Long: `Rename single-block Terraform files to follow OTN naming convention.
Files with multiple blocks are skipped (use 'split' first).`,
	RunE: runRename,
}

func init() {
	renameCmd.Flags().Bool("backup", false, "create .bak backup before renaming")
}

func runRename(cmd *cobra.Command, args []string) error {
	cfg := config.Get()

	backup, _ := cmd.Flags().GetBool("backup")

	opts := operations.Options{
		Directory: cfg.Directory,
		Recursive: cfg.Recursive,
		DryRun:    cfg.DryRun,
		Verbose:   cfg.Verbose,
		Backup:    backup,
	}

	result := output.NewResult(Version, "rename", cfg.Directory, cfg.DryRun)

	var console *output.Console
	if !cfg.JSONOutput {
		console = output.NewConsole(cfg.Verbose, cfg.Quiet)
	}

	if err := operations.Rename(opts, result, console); err != nil {
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
