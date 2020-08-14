package gsc_merge

import (
	"fmt"
	"os"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	gsc_utils "github.com/isbm/gsc/utils"
)

type GSCMergeBranch struct {
	gpr *gsc_utils.GitPkgRepo
	git *gsc_utils.GitCaller
	wzlib_logger.WzLogger
}

func NewGSCMergeBranch() *GSCMergeBranch {
	cb := new(GSCMergeBranch)
	cb.git = gsc_utils.NewGitCaller()

	var err error
	cb.gpr, err = gsc_utils.GetRepoFromFile("")
	if err != nil {
		cb.GetLogger().Errorf("Error getting Git repository: %s", err.Error())
		os.Exit(1)
	}
	return cb
}

// Merge current working branch and delete it
func (cb *GSCMergeBranch) Merge(devel bool) error {
	currentBranch := cb.git.GetCurrentBranch()
	var targetBranch string
	if devel {
		targetBranch = cb.git.GetDefaultBranch()
		cb.GetLogger().Debugf("Merging as development branch to '%s'", targetBranch)
	} else {
		targetBranch = cb.gpr.Branch
		cb.GetLogger().Debugf("Merging as release branch to '%s'", targetBranch)
	}

	if currentBranch == targetBranch {
		cb.GetLogger().Error("Not forked state.")
		return fmt.Errorf("You cannot merge the main branch to the main branch")
	} else if targetBranch == "" {
		return fmt.Errorf("Branch is missing. Please define one manually in '%s' file", gsc_utils.GIT_PKG_REPO)
	}

	cb.GetLogger().Infof("Switching to %s branch from %s", targetBranch, currentBranch)
	if err := cb.git.Call("checkout", targetBranch); err != nil {
		return err
	}

	cb.GetLogger().Infof("Merging branch %s to %s", currentBranch, targetBranch)
	if err := cb.git.Call("merge", "--no-ff", "--no-squash", "--message",
		fmt.Sprintf("Merge branch '%s'", currentBranch), currentBranch); err != nil {
		return err
	}

	cb.GetLogger().Infof("Pushing %s to the Git repo", targetBranch)
	if err := cb.git.Call("push", "-u", "origin", targetBranch); err != nil {
		return err
	}

	return nil
}
