package gsc_utils

import (
	"fmt"
	"io/ioutil"
	"strings"

	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
)

type GSCUtils struct {
	spec []string
}

func NewGSCUtils() *GSCUtils {
	return new(GSCUtils)
}

// Read the specfile. First wins...
func (utl *GSCUtils) getSpec() error {
	if utl.spec != nil {
		return nil
	}
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return err
	}
	var spec []byte
	for _, fname := range files {
		if strings.HasSuffix(fname.Name(), ".spec") {
			spec, err = ioutil.ReadFile(fname.Name())
			if err != nil {
				return err
			}
		}
	}
	utl.spec = strings.Split(string(spec), "\n")
	return nil
}

// Get a specific key from the spec
func (utl *GSCUtils) getKey(key string) (string, error) {
	if err := utl.getSpec(); err != nil {
		return "", err
	}
	for _, line := range utl.spec {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), strings.ToLower(key)+":") {
			return strings.TrimSpace(strings.Split(line, ":")[1]), nil
		}
	}
	return "", fmt.Errorf("Nothing found for %s", key)
}

// GetPackageVersion from the spec
func (utl *GSCUtils) GetPackageVersion() (string, error) {
	return utl.getKey("version")
}

// GetPackageName from the spec
func (utl *GSCUtils) GetPackageName() (string, error) {
	return utl.getKey("name")
}

// GetProjectName returns current project name
func (utl *GSCUtils) GetProjectInfo() (string, error) {
	cmd, err := wzlib_subprocess.BufferedExec("osc", "info")
	if err != nil {
		return "", err
	}
	out := cmd.StdoutString()
	cmd.Wait()
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(strings.ToLower(line), "project name:") {
			return strings.TrimSpace(strings.SplitN(line, ":", 2)[1]), nil
		}
	}
	return "", fmt.Errorf("Unable to find a project name in the current info")
}
