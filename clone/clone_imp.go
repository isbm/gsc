package gsc_clone

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

func (gw *GSCClone) getRepoFromFile() error {
	content, err := ioutil.ReadFile(GIT_PKG_REPO)
	if err != nil {
		if gw.repoUrl == "" {
			return fmt.Errorf("Git repository is missing. Clone with the Git repo included instead.")
		}

		content = []byte(fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<git>\n  <url>%s</url>\n</git>\n", gw.repoUrl))
		ioutil.WriteFile(GIT_PKG_REPO, content, 0644)
		gw.initGit = true
		return nil
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
