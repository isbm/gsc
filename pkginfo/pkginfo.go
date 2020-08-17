package gsc_info

import (
	gsc_utils "github.com/isbm/gsc/utils"
)

type GSCPackageInfo struct {
	Name    string
	Version string
	Project *gsc_utils.GSCProjectInfo
	GitRepo *gsc_utils.GitPkgRepo
	git     *gsc_utils.GitCaller
	utils   *gsc_utils.GSCUtils
}

func NewGSCPackageInfo() *GSCPackageInfo {
	pki := new(GSCPackageInfo)
	pki.git = gsc_utils.NewGitCaller()
	pki.utils = gsc_utils.NewGSCUtils()

	return pki
}

func (pki *GSCPackageInfo) ObtainInfo() error {
	var err error
	if pki.GitRepo, err = gsc_utils.GetRepoFromFile(""); err != nil {
		return err
	}
	pki.Name, err = pki.utils.GetPackageName()

	if pki.Version, err = pki.utils.GetPackageVersion(); err != nil {
		return err
	}

	if pki.Project, err = pki.utils.GetProjectInfo(); err != nil {
		return err
	}

	return nil
}
