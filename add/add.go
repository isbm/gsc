package gsc_add

import (
	"fmt"
	"os"
	"strings"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
	gsc_utils "github.com/isbm/gsc/utils"
)

/*
	Implements "add" command
*/

// GSCAdd class
type GSCAdd struct {
	pathspec string
	git      *gsc_utils.GitCaller
	wzlib_logger.WzLogger
}

// NewGSCAdd constructs new instance of a GCSAdd class
func NewGSCAdd() *GSCAdd {
	add := new(GSCAdd).SetPathspec("*")
	add.git = gsc_utils.NewGitCaller()
	return add
}

// SetMask sets a mask of files to add
func (add *GSCAdd) SetPathspec(pathspec string) *GSCAdd {
	add.pathspec = pathspec
	return add
}

func (add *GSCAdd) Add() error {
	cmd := wzlib_subprocess.ExecCommand("osc", "commit", add.pathspec)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		add.GetLogger().Error(err.Error())
		return err
	}

	entry := gsc_utils.NewGSCRevisionLog().GetLatest() // Should be never nil
	if entry == nil {
		return fmt.Errorf("Unable to find an initial log message")
	}

	add.git.Call("add", add.pathspec)
	add.git.Call("commit", "-m", fmt.Sprintf("%s", strings.ReplaceAll(entry.Message, "'", "\\'")))

	return nil
}
