package gsc_utils

import (
	"os"

	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
)

// Add to the sequence the whole changelog entry and reset WIP status
func CallWithTTY(name string, args ...string) error {
	cmd := wzlib_subprocess.ExecCommand(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
