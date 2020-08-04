package gsc_clone

import (
	"fmt"
	"os"
	"strings"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
	gsc_utils "github.com/isbm/gsc/utils"
)

// GSCClone class
type GSCClone struct {
	project string
	repoUrl string
	pkg     string
	initGit bool
	wzlib_logger.WzLogger
}

// NewGCSClone creates a package cloning tool
func NewGCSClone() *GSCClone {
	gw := new(GSCClone)
	return gw
}

// SetProject name
func (gw *GSCClone) SetProject(project string) *GSCClone {
	gw.project = project
	return gw
}

// SetPackage name
func (gw *GSCClone) SetPackage(pkg string) *GSCClone {
	gw.pkg = pkg
	return gw
}

// SetGitRepo name
func (gw *GSCClone) SetGitRepoUrl(repoUrl string) *GSCClone {
	gw.repoUrl = repoUrl
	return gw
}

func (gw *GSCClone) getGitRepoUrl() string {
	repo := ""
	if strings.Contains(gw.repoUrl, "@") {
		if !strings.HasPrefix(gw.repoUrl, "git@") {
			panic("Wrong git url!")
		}
		repo = gw.repoUrl
	} else {
		repo = fmt.Sprintf("git@%s.git", gw.repoUrl)
	}
	return repo
}

func (gw *GSCClone) setupGit() error {
	var err error

	if gw.initGit {
		wzlib_subprocess.BufferedExec("git", "init")
		wzlib_subprocess.BufferedExec("git", "add", "--all", "--force")
		wzlib_subprocess.BufferedExec("git", "commit", "-m", "initial commit")
		wzlib_subprocess.BufferedExec("git", "remote", "add", "origin", gw.getGitRepoUrl())
	} else {
		// Then:
		// - move all git files to the current package, overwriting everything
		_, err = wzlib_subprocess.BufferedExec("git", "clone", gw.getGitRepoUrl())
	}
	return err
}

// Clone package with the bind to the Git repo
func (gw *GSCClone) Clone() error {
	if err := gw.getRepoFromFile(); err != nil {
		return err
	}

	usr, err := gsc_utils.GetOSCUser()
	if err != nil {
		return err
	}
	oscPath := fmt.Sprintf("home:%s:branches:%s", usr.Uid, gw.project)
	_, stderr := wzlib_subprocess.StreamedExec(NewFetcherProcessStream(oscPath).SetPrefix("Checking out"), "osc", "bco", gw.project, gw.pkg)

	if stderr != "" {
		gw.GetLogger().Debug(strings.TrimSpace(stderr))
	}
	wzlib_subprocess.BufferedExec("mv", GIT_PKG_REPO, fmt.Sprintf("%s/%s/", oscPath, gw.pkg))
	os.Chdir(fmt.Sprintf("%s/%s/", oscPath, gw.pkg))

	_, err = wzlib_subprocess.BufferedExec("osc", "add", GIT_PKG_REPO)
	if err != nil {
		return err
	}

	// - git-clone from the repo
	if err = gw.setupGit(); err != nil {
		return err
	}

	return nil
}
