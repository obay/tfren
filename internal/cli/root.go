package cli

import (
	"fmt"
	"os"

	"github.com/obay/tfren/internal/config"
	"github.com/obay/tfren/internal/git"
	"github.com/obay/tfren/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version = "dev"
	cfgFile string
	out     *output.Console
)

var rootCmd = &cobra.Command{
	Use:   "tfren",
	Short: "Terraform file organizer",
	Long: `tfren organizes Terraform configuration files according to
Obay's Terraform Naming Convention (OTN).

It splits multi-block files into individual files and renames them
to follow the pattern: <block_type>.<provider>.<name>.tf

Running 'tfren' without a subcommand executes 'organize' (split + rename).`,
	PersistentPreRunE: preRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		return organizeCmd.RunE(cmd, args)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is .tfren.json)")
	rootCmd.PersistentFlags().StringP("directory", "d", ".", "target directory")
	rootCmd.PersistentFlags().BoolP("recursive", "r", false, "process subdirectories recursively")
	rootCmd.PersistentFlags().BoolP("dry-run", "n", false, "show what would be done without making changes")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "suppress non-error output")
	rootCmd.PersistentFlags().Bool("json", false, "output in JSON format")
	rootCmd.PersistentFlags().Bool("no-git-check", false, "skip git repository status check")
	rootCmd.PersistentFlags().BoolP("force", "f", false, "skip confirmation prompts")

	viper.BindPFlag("directory", rootCmd.PersistentFlags().Lookup("directory"))
	viper.BindPFlag("recursive", rootCmd.PersistentFlags().Lookup("recursive"))
	viper.BindPFlag("dry_run", rootCmd.PersistentFlags().Lookup("dry-run"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
	viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))
	viper.BindPFlag("no_git_check", rootCmd.PersistentFlags().Lookup("no-git-check"))
	viper.BindPFlag("force", rootCmd.PersistentFlags().Lookup("force"))

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(splitCmd)
	rootCmd.AddCommand(renameCmd)
	rootCmd.AddCommand(organizeCmd)
	rootCmd.AddCommand(validateCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".tfren")
		viper.SetConfigType("json")
		viper.AddConfigPath(".")
		home, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(home)
		}
	}

	viper.SetEnvPrefix("TFREN")
	viper.AutomaticEnv()

	viper.ReadInConfig()
}

func preRun(cmd *cobra.Command, args []string) error {
	cfg := config.Get()
	out = output.NewConsole(cfg.Verbose, cfg.Quiet)

	if cfg.Quiet || cfg.JSONOutput {
		return nil
	}

	if cmd.Name() == "validate" || cmd.Name() == "version" || cfg.DryRun {
		return nil
	}

	// Check if git check is disabled
	noGitCheck, _ := cmd.Flags().GetBool("no-git-check")
	if noGitCheck {
		return nil
	}

	// Show progress while checking git status
	fmt.Print("Checking git status... ")
	status := git.CheckStatus(cfg.Directory, cfg.Recursive)
	fmt.Print("\r                       \r") // Clear the line

	if warning := status.Warning(); warning != "" {
		out.Warning(warning)
	}

	// If uncommitted .tf files in target dir, prompt for confirmation
	if status.RequiresConfirmation() {
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			if !git.PromptContinue() {
				return fmt.Errorf("operation cancelled by user")
			}
		}
	}

	return nil
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tfren version %s\n", Version)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func SetVersion(v string) {
	Version = v
}
