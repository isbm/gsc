package gsc_utils

import (
	"encoding/xml"
)

/*
	The _git_pkg_repo file structure.

	<?xml version="1.0" encoding="UTF-8"?>
	<git>
	  <url>git@github.com:somebody/my_great_package.git</url>
	</git>
*/
var GIT_PKG_REPO string = "_git_pkg_repo"

type GitPkgRepo struct {
	XMLName xml.Name `xml:"git"`
	Url     string   `xml:"url"`
}
