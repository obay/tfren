package main

import (
	"os"

	"github.com/obay/tfren/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
