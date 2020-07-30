package gsc_clone

// GSCClone class
type GSCClone struct {
	project string
	repoUrl string
	pkg     string
}

// NewGCSClone creates a package cloning tool
func NewGCSClone() *GSCClone {
	gw := new(GSCClone)
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

// Clone package with the bind to the Git repo
func (gw *GSCClone) Clone() error {
	if err := gw.getRepoFromFile(); err != nil {
		return err
	}
	return nil
}
