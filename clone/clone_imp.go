package gsc_clone

import (
	gsc_utils "github.com/isbm/gsc/utils"
)

func (gw *GSCClone) getRepoFromFile() error {
	var err error
	gw.initGit, gw.repoUrl, err = gsc_utils.GetRepoFromFile(gw.repoUrl)
	return err
}
