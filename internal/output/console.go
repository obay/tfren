package output

import (
	"fmt"

	"github.com/fatih/color"
)

type Console struct {
	verbose bool
	quiet   bool
}

func NewConsole(verbose, quiet bool) *Console {
	return &Console{
		verbose: verbose,
		quiet:   quiet,
	}
}

func (c *Console) Success(msg string) {
	if !c.quiet {
		color.Green(msg)
	}
}

func (c *Console) Warning(msg string) {
	if !c.quiet {
		color.Yellow(msg)
	}
}

func (c *Console) Error(msg string) {
	color.Red(msg)
}

func (c *Console) Info(msg string) {
	if c.verbose && !c.quiet {
		fmt.Println(msg)
	}
}

func (c *Console) Print(msg string) {
	if !c.quiet {
		fmt.Println(msg)
	}
}
