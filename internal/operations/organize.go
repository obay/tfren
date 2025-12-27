package operations

import (
	"github.com/obay/tfrn/internal/output"
)

func Organize(opts Options, result *output.Result, console *output.Console) error {
	if err := Split(opts, result, console); err != nil {
		return err
	}

	result.FilesProcessed = 0

	if err := Rename(opts, result, console); err != nil {
		return err
	}

	return nil
}
