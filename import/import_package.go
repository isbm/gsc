package gsc_import

import (
	"fmt"

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
	wzlib_logger.WzLogger
}

// NewGSCPackageImport constructor
func NewGSCPackageImport() *GSCPackageImport {
	imp := new(GSCPackageImport)
	return imp
}

// Automatically try to load Git repo URL from the _git_repo file, if any.
func (imp *GSCPackageImport) autoloadGitRepoUrl() {

}

// SetGitRepoUrl
func (imp *GSCPackageImport) SetGitRepoUrl(repo string) *GSCPackageImport {
	imp.repoUrl = repo
	return imp
}

// Import
func (imp *GSCPackageImport) Import() error {
	if imp.repoUrl == "" {
		err := fmt.Errorf("Repo to import from was not found")
		imp.GetLogger().Error(err)
		return err
	}
	return nil
}
