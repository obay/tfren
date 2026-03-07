package cli

import (
	"github.com/obay/tfren/internal/config"
	"github.com/obay/tfren/internal/operations"
	"github.com/obay/tfren/internal/output"
	"github.com/spf13/cobra"
)

var organizeCmd = &cobra.Command{
	Use:   "organize",
	Short: "Split and rename files (default command)",
	Long: `Organize Terraform files by first splitting multi-block files,
then renaming all files to follow OTN naming convention.

This is the default command when running 'tfrn' without a subcommand.`,
	RunE: runOrganize,
}

func init() {
	organizeCmd.Flags().Bool("keep-original", false, "keep original files after splitting")
	organizeCmd.Flags().Bool("backup", false, "create .bak backups before modifying")
}

func runOrganize(cmd *cobra.Command, args []string) error {
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

	result := output.NewResult(Version, "organize", cfg.Directory, cfg.DryRun)

	var console *output.Console
	if !cfg.JSONOutput {
		console = output.NewConsole(cfg.Verbose, cfg.Quiet)
	}

	if err := operations.Organize(opts, result, console); err != nil {
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
