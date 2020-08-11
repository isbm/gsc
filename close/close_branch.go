package gsc_close

import (
	"fmt"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
	gsc_utils "github.com/isbm/gsc/utils"
)

type GSCCloseBranch struct {
	git *gsc_utils.GitCaller
	wzlib_logger.WzLogger
}

func NewGSCCloseBranch() *GSCCloseBranch {
	cb := new(GSCCloseBranch)
	cb.git = gsc_utils.NewGitCaller()
	return cb
}

// Close current branch and delete it
func (cb *GSCCloseBranch) Close() error {
	currentBranch := cb.git.GetCurrentBranch()
	defaultBranch := cb.git.GetDefaultBranch()
	if currentBranch == defaultBranch {
		err := fmt.Errorf("Cannot close main branch")
		cb.GetLogger().Error(err.Error())
		return err
	}
	if err := wzlib_subprocess.ExecCommand("git", "checkout", cb.git.GetDefaultBranch()).Run(); err != nil {
		return err
	}
	return wzlib_subprocess.ExecCommand("git", "branch", "--delete", "--force", currentBranch).Run()
}
