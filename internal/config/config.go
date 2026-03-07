package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Directory  string
	Recursive  bool
	DryRun     bool
	Verbose    bool
	Quiet      bool
	JSONOutput bool
	NoGitCheck bool
	Force      bool
}

func Get() *Config {
	return &Config{
		Directory:  viper.GetString("directory"),
		Recursive:  viper.GetBool("recursive"),
		DryRun:     viper.GetBool("dry_run"),
		Verbose:    viper.GetBool("verbose"),
		Quiet:      viper.GetBool("quiet"),
		JSONOutput: viper.GetBool("json"),
		NoGitCheck: viper.GetBool("no_git_check"),
		Force:      viper.GetBool("force"),
	}
}
