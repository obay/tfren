package main

import (
	"os"

	"github.com/obay/tfrn/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
