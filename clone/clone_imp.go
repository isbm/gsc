package gsc_clone

import (
	gsc_utils "github.com/isbm/gsc/utils"
)

func (gw *GSCClone) getRepoFromFile() error {
	var err error
	gw.gpr, err = gsc_utils.GetRepoFromFile(gw.gpr.Url)
	return err
}
