package gsc_clone

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

func (gw *GSCClone) getRepoFromFile() error {
	content, err := ioutil.ReadFile(GIT_PKG_REPO)
	if err != nil {
		return fmt.Errorf("Git repository is missing. Clone with the Git repo included instead.")
	}

	var gitPkgRepo GitPkgRepo
	if err := xml.Unmarshal(content, &gitPkgRepo); err != nil {
		return err
	}

	if gw.repoUrl != "" {
		return fmt.Errorf("Git repository is already specified as: %s", gitPkgRepo.Url)
	}

	gw.repoUrl = gitPkgRepo.Url
	return nil
}
