package operations

import (
	"github.com/obay/tfren/internal/output"
)

func Organize(opts Options, result *output.Result, console *output.Console) error {
	if err := Split(opts, result, console); err != nil {
		return err
	}

	result.FilesProcessed = 0

	// Rename phase is silent during organize - pass nil console
	// Files just created by split will be "already compliant" which is expected
	if err := Rename(opts, result, nil); err != nil {
		return err
	}

	return nil
}
