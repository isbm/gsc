package gsc_merge

import (
	"fmt"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	gsc_utils "github.com/isbm/gsc/utils"
)

type GSCMergeBranch struct {
	git *gsc_utils.GitCaller
	wzlib_logger.WzLogger
}

func NewGSCMergeBranch() *GSCMergeBranch {
	cb := new(GSCMergeBranch)
	cb.git = gsc_utils.NewGitCaller()
	return cb
}

// Merge current working branch and delete it
func (cb *GSCMergeBranch) Merge() error {
	currentBranch := cb.git.GetCurrentBranch()
	defaultBranch := cb.git.GetDefaultBranch()
	if currentBranch == defaultBranch {
		cb.GetLogger().Error("Not forked state.")
		return fmt.Errorf("You cannot merge the main branch to the main branch")
	}

	cb.GetLogger().Infof("Switching to %s branch from %s", defaultBranch, currentBranch)
	if err := cb.git.Call("checkout", defaultBranch); err != nil {
		return err
	}

	cb.GetLogger().Infof("Merging branch %s to %s", currentBranch, defaultBranch)
	if err := cb.git.Call("merge", "--no-ff", "--no-squash", "--message",
		fmt.Sprintf("Merge branch '%s'", currentBranch), currentBranch); err != nil {
		return err
	}

	cb.GetLogger().Infof("Pushing %s to the Git repo", defaultBranch)
	if err := cb.git.Call("push", "-u", "origin", defaultBranch); err != nil {
		return err
	}

	return nil
}
