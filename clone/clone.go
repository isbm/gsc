package gsc_clone

import (
	"fmt"
	"os"

	gsc_utils "github.com/isbm/gsc/utils"

	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
)

// GSCClone class
type GSCClone struct {
	project string
	repoUrl string
	pkg     string
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
	_, stderr := wzlib_subprocess.StreamedExec(NewFetcherProcessStream(oscPath), "osc", "bco", gw.project, gw.pkg)

	if stderr != "" {
		fmt.Println("ERR:\n", stderr)
	}
	wzlib_subprocess.BufferedExec("mv", GIT_PKG_REPO, fmt.Sprintf("%s/%s/", oscPath, gw.pkg))
	os.Chdir(fmt.Sprintf("%s/%s/", oscPath, gw.pkg))
	wzlib_subprocess.BufferedExec("osc", "add", GIT_PKG_REPO)
	return nil
}
