package gsc_clone

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
	gsc_push "github.com/isbm/gsc/push"
	gsc_utils "github.com/isbm/gsc/utils"
)

// GSCClone class
type GSCClone struct {
	project string
	pkg     string

	gpr      *gsc_utils.GitPkgRepo
	pkgutils *gsc_utils.GSCUtils
	git      *gsc_utils.GitCaller

	wzlib_logger.WzLogger
}

// NewGCSClone creates a package cloning tool
func NewGCSClone() *GSCClone {
	gw := new(GSCClone)
	gw.pkgutils = gsc_utils.NewGSCUtils()
	gw.git = gsc_utils.NewGitCaller()
	gw.gpr = new(gsc_utils.GitPkgRepo)

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
	gw.gpr.Url = repoUrl
	return gw
}

// SetGitBranch of the package release
func (gw *GSCClone) SetGitBranch(branch string) *GSCClone {
	gw.gpr.Branch = branch
	return gw
}

func (gw *GSCClone) getGitRepoUrl() string {
	repo := ""
	if strings.Contains(gw.gpr.Url, "@") {
		if !strings.HasPrefix(gw.gpr.Url, "git@") {
			gw.GetLogger().Panic("Wrong Git URL")
		}
		repo = gw.gpr.Url
	} else {
		repo = fmt.Sprintf("git@%s.git", gw.gpr.Url)
	}
	return repo
}

func (gw *GSCClone) addDefaultGitIgnore() error {
	var out bytes.Buffer
	out.WriteString(".osc/\n")
	out.WriteString(".gitignore\n")
	return ioutil.WriteFile(".gitignore", out.Bytes(), 0644)
}

func (gw *GSCClone) setupGit() error {
	var err error
	if gw.gpr.IsNew() {
		gw.GetLogger().Info("The package was not yet linked to the repository.")
		gw.addDefaultGitIgnore()
		gw.git.Call("init")
		gw.git.Call("add", "--all")
		gw.git.Call("reset", "--", ".osc")
		gw.git.Call("commit", "-m", "initial commit")
		gw.git.Call("remote", "add", "origin", gw.getGitRepoUrl())

		// This is the new package, so initial push to the repo required
		if err := gsc_push.NewGCSPush().Push(); err != nil {
			return err
		}

		if err := wzlib_subprocess.ExecCommand("osc", "commit", "-m", "initial link with Git repo").Run(); err != nil {
			gw.GetLogger().Errorf("Unable to link OSC repo with Git: %s", err.Error())
		} else {
			gw.GetLogger().Info("OSC repo has been linked with the Git")
		}
	} else {
		gitRepoURL := gw.getGitRepoUrl()
		tempGitRepo := gsc_utils.RandomString() + "-repo"

		gw.GetLogger().Infof("Getting package repository from %s", gitRepoURL)
		gw.git.Call("clone", gitRepoURL, tempGitRepo)

		files, err := ioutil.ReadDir("./" + tempGitRepo)
		if err != nil {
			return err
		}
		for _, nfo := range files {
			gw.GetLogger().Debugf("Moving %s", nfo.Name())
			os.Rename(path.Join(tempGitRepo, nfo.Name()), nfo.Name())
		}
		if err := wzlib_subprocess.ExecCommand("rm", "-rf", tempGitRepo).Run(); err != nil {
			return err
		}
	}
	// Branch to something
	var pkVer string
	var pkName string
	if pkVer, err = gw.pkgutils.GetPackageVersion(); err != nil {
		return err
	}
	if pkName, err = gw.pkgutils.GetPackageName(); err != nil {
		return err
	}

	if gw.gpr.Branch == "" {
		return fmt.Errorf("Branch is missing")
	}

	gw.GetLogger().Infof("Switching to the base branch '%s'", gw.gpr.Branch)
	gw.git.Call("checkout", gw.gpr.Branch)
	gw.git.Call("pull")

	branch := fmt.Sprintf("tmp-%s-%s", pkName, pkVer)
	gw.GetLogger().Infof("Creating and switching to a new Git branch: %s", branch)
	gw.git.Call("checkout", "-b", branch)

	return err
}

// Clone package with the bind to the Git repo
func (gw *GSCClone) Clone() error {
	usr, err := gsc_utils.GetOSCUser()
	if err != nil {
		return err
	}
	oscPath := fmt.Sprintf("home:%s:branches:%s", usr.Uid, gw.project)
	_, stderr := wzlib_subprocess.StreamedExec(NewFetcherProcessStream(oscPath).SetPrefix("Checking out"), "osc", "bco", gw.project, gw.pkg)

	if stderr != "" {
		gw.GetLogger().Debug(strings.TrimSpace(stderr))
	}
	wzlib_subprocess.BufferedExec("mv", gsc_utils.GIT_PKG_REPO, fmt.Sprintf("%s/%s/", oscPath, gw.pkg))
	os.Chdir(fmt.Sprintf("%s/%s/", oscPath, gw.pkg))

	_, err = wzlib_subprocess.BufferedExec("osc", "add", gsc_utils.GIT_PKG_REPO)
	if err != nil {
		return err
	}

	if err := gw.getRepoFromFile(); err != nil {
		return err
	}

	if err = gw.setupGit(); err != nil {
		return err
	}

	return nil
}
