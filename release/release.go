package gsc_release

import (
	"io/ioutil"
	"os"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	gsc_utils "github.com/isbm/gsc/utils"
)

type GSCPackageRelease struct {
	gpr *gsc_utils.GitPkgRepo
	git *gsc_utils.GitCaller
	wzlib_logger.WzLogger
}

func NewGSCPackageRelease() *GSCPackageRelease {
	rel := new(GSCPackageRelease)

	var err error
	rel.gpr, err = gsc_utils.GetRepoFromFile("")
	if err != nil {
		rel.GetLogger().Errorf("Error obtaining repository data: %s", err.Error())
		os.Exit(1)
	}

	rel.git = gsc_utils.NewGitCaller()

	return rel
}

// SetReleaseBranch to which package is going to be released
func (rel *GSCPackageRelease) SetReleaseBranch(branch string) *GSCPackageRelease {
	rel.gpr.Branch = branch
	return rel
}

// Release package
func (rel *GSCPackageRelease) Release() error {
	rel.GetLogger().Debugf("Updating '%s' with branch '%s'", gsc_utils.GIT_PKG_REPO, rel.gpr.Branch)
	if err := ioutil.WriteFile(gsc_utils.GIT_PKG_REPO, []byte(rel.gpr.ToXML()), 0644); err != nil {
		return err
	}

	// Save current branch
	tempBranch := rel.git.GetCurrentBranch()

	// Create release branch
	rel.GetLogger().Debugf("Creating branch '%s'", rel.gpr.Branch)
	rel.git.Call("checkout", "-b", rel.gpr.Branch)

	// Remove temporary branch
	rel.git.Call("branch", "-D", tempBranch)

	// Push the new branch
	rel.git.Call("push", "-u", "origin", rel.gpr.Branch)

	return nil
}
