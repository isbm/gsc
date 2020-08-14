package gsc_utils

import (
	"encoding/xml"
	"fmt"
)

/*
	The _git_pkg_repo file structure.

	<?xml version="1.0" encoding="UTF-8"?>
	<git>
	  <url>git@github.com:somebody/my_great_package.git</url>
	  <branch>release-1.2</branch>
	</git>
*/
var GIT_PKG_REPO string = "_git_pkg_repo"

type GitPkgRepo struct {
	XMLName xml.Name `xml:"git"`
	Url     string   `xml:"url"`
	Branch  string   `xml:"branch"`
	isNew   bool
}

// SetIsNew to true if the repo XML file was just created
func (gpr *GitPkgRepo) SetIsNew() {
	gpr.isNew = true
}

// IsNew returns true if repo file was just created
func (gpr *GitPkgRepo) IsNew() bool {
	return gpr.isNew
}

// ToXML serialisation
func (gpr *GitPkgRepo) ToXML() string {
	return fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<git>\n  <url>%s</url>\n  <branch>%s</branch>\n</git>\n",
		gpr.Url, gpr.Branch)
}
