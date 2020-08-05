package gsc_utils

import (
	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
)

type GitCaller struct {
	wzlib_logger.WzLogger
}

func NewGitCaller() *GitCaller {
	return new(GitCaller)
}

// Call git with specific params. All calls are blocking.
func (gitcall *GitCaller) Call(args ...string) {
	cmd := wzlib_subprocess.ExecCommand("git", args...)
	err := cmd.Run()
	if err != nil {
		gitcall.GetLogger().Errorf("Error calling Git: %s", err.Error())
	}
	if err != nil {
		gitcall.GetLogger().Errorf("Error completing Git call: %s", err.Error())
	}
}
