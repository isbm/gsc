package gsc_clone

import (
	"fmt"
	"os"
	"strings"

	gsc_push "github.com/isbm/gsc/push"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
	gsc_utils "github.com/isbm/gsc/utils"
)

// GSCClone class
type GSCClone struct {
	project  string
	repoUrl  string
	pkg      string
	initGit  bool
	pkgutils *gsc_utils.GSCUtils
	git      *gsc_utils.GitCaller
	wzlib_logger.WzLogger
}

// NewGCSClone creates a package cloning tool
func NewGCSClone() *GSCClone {
	gw := new(GSCClone)
	gw.pkgutils = gsc_utils.NewGSCUtils()
	gw.git = gsc_utils.NewGitCaller()

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
		gw.git.Call("init")
		gw.git.Call("add", "--all")
		gw.git.Call("reset", "--", ".osc")
		gw.git.Call("commit", "-m", "initial commit")
		gw.git.Call("remote", "add", "origin", gw.getGitRepoUrl())

		var pkVer string
		var pkName string
		if pkVer, err = gw.pkgutils.GetPackageVersion(); err != nil {
			return err
		}
		if pkName, err = gw.pkgutils.GetPackageName(); err != nil {
			return err
		}

		// This is the new package, so initial push to the repo required
		if err := gsc_push.NewGCSPush().Push(); err != nil {
			return err
		}

		if err := wzlib_subprocess.ExecCommand("osc", "commit", "-m", "initial link with Git repo").Run(); err != nil {
			gw.GetLogger().Errorf("Unable to link OSC repo with Git: %s", err.Error())
		} else {
			gw.GetLogger().Info("OSC repo has been linked with the Git")
		}

		// Branch to something
		branch := fmt.Sprintf("tmp-%s-%s", pkName, pkVer)
		gw.git.Call("checkout", "-b", branch)
		gw.GetLogger().Infof("New working Git branch created: %s", branch)
	} else {
		// Then:
		// - move all git files to the current package, overwriting everything
		gw.git.Call("clone", gw.getGitRepoUrl())
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
