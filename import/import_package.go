package gsc_import

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	gsc_utils "github.com/isbm/gsc/utils"
)

/*
	Implements package import.

	This operation assumes that the Git source is always correct,
	and whatever content is in the current package, it will be "flashed"
	by the content of the Git repo.
*/

type GSCPackageImport struct {
	repoUrl  string
	pkgutils *gsc_utils.GSCUtils
	git      *gsc_utils.GitCaller
	nfo      *gsc_utils.GSCProjectInfo
	wzlib_logger.WzLogger
}

// NewGSCPackageImport constructor
func NewGSCPackageImport() *GSCPackageImport {
	imp := new(GSCPackageImport)
	imp.pkgutils = gsc_utils.NewGSCUtils()
	imp.git = gsc_utils.NewGitCaller()
	imp.autoloadGitRepoUrl()

	var err error
	imp.nfo, err = imp.pkgutils.GetProjectInfo()
	if err != nil {
		imp.GetLogger().Errorf("Unable to get project info: %s", err.Error())
		os.Exit(1)
	}

	return imp
}

// Automatically try to load Git repo URL from the _git_repo file, if any.
func (imp *GSCPackageImport) autoloadGitRepoUrl() {
	var err error
	_, imp.repoUrl, err = gsc_utils.GetRepoFromFile(imp.repoUrl)
	if err != nil {
		imp.GetLogger().Warning(err.Error())
	}
}

// SetGitRepoUrl
func (imp *GSCPackageImport) SetGitRepoUrl(repo string) *GSCPackageImport {
	if strings.HasPrefix(repo, "//") {
		repo = strings.SplitN(repo, "//", 2)[1]
	}

	if !strings.HasSuffix(repo, ".git") {
		repo = fmt.Sprintf("git@%s.git", repo)
	}

	imp.repoUrl = repo
	imp.GetLogger().Debugf("Using %s repo", repo)

	return imp
}

// Checks if the working directory is correct
func (imp *GSCPackageImport) checkDir() error {
	var err error
	files, err := ioutil.ReadDir(".")
	if err != nil {
		imp.GetLogger().Error(err.Error())
		os.Exit(1)
	}

	expectedSpec := fmt.Sprintf("%s.spec", imp.nfo.PackageName)
	oscf, spec := false, false
	for _, nfo := range files {
		if nfo.Name() == expectedSpec {
			spec = true
		} else if nfo.Name() == ".osc" {
			oscf = true
		}
	}

	if !oscf {
		return fmt.Errorf("This is not package directory")
	} else if !spec {
		return fmt.Errorf("Spec file %s was not found", expectedSpec)
	}

	return nil
}

// Import
func (imp *GSCPackageImport) Import() error {
	if err := imp.checkDir(); err != nil {
		return err
	}

	if imp.repoUrl == "" {
		return fmt.Errorf("Repo to import from was not found or specified")
	}

	fmt.Println("Importing...")

	// Clone to a temporary directory
	tempCloneDir := fmt.Sprintf("%s-%s", imp.nfo.PackageName, gsc_utils.RandomString())
	imp.GetLogger().Debugf("Cloning into %s directory", tempCloneDir)
	imp.git.Call("clone", imp.repoUrl, tempCloneDir)

	// Remove everything in the current directory, except .osc and temporary clone
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	currentFiles, err := ioutil.ReadDir(workingDir)
	if err != nil {
		return err
	}

	for _, nfo := range currentFiles {
		if nfo.Name() != ".osc" && nfo.Name() != tempCloneDir {
			os.RemoveAll(nfo.Name())
		}
	}

	// Read all that
	clonedFilesPath := path.Join(workingDir, tempCloneDir)
	clonedFiles, err := ioutil.ReadDir(clonedFilesPath)
	if err != nil {
		return err
	}

	errors := 0
	for _, nfo := range clonedFiles {
		imp.GetLogger().Debugf("Fetching file: %s", nfo.Name())
		if err := os.Rename(path.Join(clonedFilesPath, nfo.Name()), path.Join(workingDir, nfo.Name())); err != nil {
			imp.GetLogger().Warningf("Error moving file %s: %s", nfo.Name(), err.Error())
			errors++
		}
	}

	if errors > 0 {
		return fmt.Errorf("Unfortunately errors had happened during the import.")
	}

	// Cleanup
	os.RemoveAll(clonedFilesPath)

	// Add stuff to the package, if any
	files, err := filepath.Glob("*")
	if err != nil {
		return err
	}
	args := []string{"add"}
	for _, fname := range files {
		if !strings.HasPrefix(fname, ".") {
			args = append(args, fname)
		}
	}

	// Update content of the package
	if err := wzlib_subprocess.ExecCommand("osc", args...).Run(); err != nil {
		return err
	}

	// Commit content to OSC
	if err := gsc_utils.CallWithTTY("osc", "commit"); err != nil {
		return err
	}

	// Make submit request
	return gsc_utils.CallWithTTY("osc", "sr", "--yes")
}
