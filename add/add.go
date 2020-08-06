package gsc_add

import (
	"fmt"
	"os"
	"path/filepath"
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
	chlog    *gsc_utils.GSCChangeLog

	wzlib_logger.WzLogger
}

// NewGSCAdd constructs new instance of a GCSAdd class
func NewGSCAdd() *GSCAdd {
	add := new(GSCAdd).SetPathspec("*")
	add.git = gsc_utils.NewGitCaller()
	add.chlog = gsc_utils.NewGSCChangeLog()
	return add
}

// SetMask sets a mask of files to add
func (add *GSCAdd) SetPathspec(pathspec string) *GSCAdd {
	add.pathspec = pathspec
	return add
}

func (add *GSCAdd) expandPathspec() []string {
	if add.pathspec == "*" {
		files, _ := filepath.Glob(add.pathspec)
		out := []string{}
		for _, fname := range files {
			if !strings.HasPrefix(fname, ".") {
				out = append(out, fname)
			}
		}
		return out
	} else {
		return []string{add.pathspec}
	}
}

func (add *GSCAdd) Add() error {
	files := add.expandPathspec()
	cmd := wzlib_subprocess.ExecCommand("osc", append(files[:0], append([]string{"add"}, files[0:]...)...)...)
	if err := cmd.Run(); err != nil {
		add.GetLogger().Error(err.Error())
		return err
	}

	files = add.expandPathspec()
	cmd = wzlib_subprocess.ExecCommand("osc", append(files[:0], append([]string{"commit"}, files[0:]...)...)...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		add.GetLogger().Error(err.Error())
		return err
	}

	entry := gsc_utils.NewGSCRevisionLog().GetLatest() // Should be never nil
	if entry == nil {
		return fmt.Errorf("Unable to find an initial log message")
	}

	files = add.expandPathspec()
	add.git.Call(append(files[:0], append([]string{"add"}, files[0:]...)...)...)
	add.git.Call("commit", "-m", fmt.Sprintf("%s", strings.ReplaceAll(entry.Message, "'", "\\'")))

	return nil
}
