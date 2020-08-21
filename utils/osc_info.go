package gsc_utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
)

// Definition of OSC status
type GSCProjectStatus struct {
	Modified  []string
	New       []string
	Deleted   []string
	Untracked []string
}

func NewGSCProjectStatus() *GSCProjectStatus {
	stat := new(GSCProjectStatus)
	stat.Modified = make([]string, 0)
	stat.New = make([]string, 0)
	stat.Deleted = make([]string, 0)
	stat.Untracked = make([]string, 0)
	return stat
}

type GSCProjectInfo struct {
	ProjectName string
	PackageName string
	Path        string
	ApiUrl      string
	SourceUrl   string
	Md5         string
	Revision    string
	LinkInfo    string
}

type GSCUtils struct {
	spec []string
	wzlib_logger.WzLogger
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
func (utl *GSCUtils) GetProjectInfo() (*GSCProjectInfo, error) {
	cmd, err := wzlib_subprocess.BufferedExec("osc", "info")
	if err != nil {
		return nil, err
	}
	out := cmd.StdoutString()
	cmd.Wait()

	nfo := new(GSCProjectInfo)
	for _, line := range strings.Split(out, "\n") {
		keyset := strings.SplitN(line, ":", 2)
		if len(keyset) == 2 {
			key, value := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(keyset[0]), " ", "")), strings.TrimSpace(keyset[1])
			switch key {
			case "projectname":
				nfo.ProjectName = value
			case "packagename":
				nfo.PackageName = value
			case "apiurl":
				nfo.ApiUrl = value
			case "sourceurl":
				nfo.SourceUrl = value
			case "srcmd5":
				nfo.Md5 = value
			case "revision":
				nfo.Revision = value
			case "linkinfo":
				nfo.LinkInfo = value
			case "path":
				nfo.Path = value
			default:
				utl.GetLogger().Warningf("Unknown OSC info section: %s", key)
			}
		}
	}
	if nfo.PackageName == "" {
		return nil, fmt.Errorf("missing package information")
	}
	return nfo, nil
}

// GetStatus about the OSC current project (list deleted, modified and added files)
func (utl *GSCUtils) GetStatus() (*GSCProjectStatus, error) {
	cmd, err := wzlib_subprocess.BufferedExec("osc", "status")
	if err != nil {
		return nil, err
	}
	out := cmd.StdoutString()
	cmd.Wait()

	status := NewGSCProjectStatus()

	for _, line := range strings.Split(out, "\n") {
		statFname := strings.SplitN(strings.TrimSpace(line), " ", 2)
		if len(statFname) == 2 {
			stat := strings.ToLower(strings.TrimSpace(statFname[0]))
			fname := strings.TrimSpace(statFname[1])
			switch stat {
			case "!":
				status.Deleted = append(status.Deleted, fname)
			case "a":
				status.New = append(status.New, fname)
			case "m":
				status.Modified = append(status.Modified, fname)
			case "?":
				status.Untracked = append(status.Untracked, fname)
			default:
				utl.GetLogger().Errorf("Unknown OSC status '%s' for file '%s'. Please report a bug.", stat, fname)
				os.Exit(1)
			}
		}

	}

	return status, nil
}
