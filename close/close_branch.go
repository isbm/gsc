package gsc_close

import (
	"fmt"
	"os"
	"path"
	"strings"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
	gsc_utils "github.com/isbm/gsc/utils"
)

type GSCCloseBranch struct {
	git   *gsc_utils.GitCaller
	utils *gsc_utils.GSCUtils
	wzlib_logger.WzLogger
}

func NewGSCCloseBranch() *GSCCloseBranch {
	cb := new(GSCCloseBranch)
	cb.git = gsc_utils.NewGitCaller()
	cb.utils = gsc_utils.NewGSCUtils()
	return cb
}

// Cleanup subproject, remove everything
func (cb *GSCCloseBranch) Cleanup() error {
	// XXX: A bit of crude way... The code needs a bit better refactoring later.

	projName, err := cb.utils.GetProjectInfo()
	if err != nil {
		return err
	}

	pkgName, err := cb.utils.GetPackageName()
	if err != nil {
		return err
	}

	if !strings.Contains(projName, ":branches:") {
		return fmt.Errorf("You should never commit directly to the master branch of the package!")
	}

	cb.GetLogger().Infof("Marking content of the sub-project '%s' as removed", projName)
	branchName := fmt.Sprintf("%s/%s", projName, pkgName)
	cb.GetLogger().Debugf("Branch name: %s", branchName)
	parentDir := path.Dir(path.Dir(branchName))
	cb.GetLogger().Debugf("Moving to %s", parentDir)

	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}
	os.Chdir("../../")

	cmd, err := wzlib_subprocess.BufferedExec("osc", "delete", branchName)
	if err != nil {
		return err
	}
	out := cmd.StdoutString()
	srr := cmd.StderrString()

	cmd.Wait()

	cb.GetLogger().Debug(out)
	cb.GetLogger().Debug(srr)
	cb.GetLogger().Infof("Removing content...")

	cb.GetLogger().Debugf("Chdir to %s for commit", workingDir)
	os.Chdir(workingDir)
	if err := wzlib_subprocess.ExecCommand("osc", "commit").Run(); err != nil {
		cb.GetLogger().Warnf("Error while committing to the OBS: %s", err.Error())
	}
	os.Chdir("../")

	cb.GetLogger().Infof("Removing the entire local content of the %s", pkgName)
	if err := wzlib_subprocess.ExecCommand("rm", "-rf", pkgName).Run(); err != nil {
		return err
	}

	return nil
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
	cb.GetLogger().Infof("Closing %s branch and discarding all the changes.", currentBranch)
	if err := wzlib_subprocess.ExecCommand("git", "checkout", cb.git.GetDefaultBranch()).Run(); err != nil {
		return err
	}
	return wzlib_subprocess.ExecCommand("git", "branch", "--delete", "--force", currentBranch).Run()
}
