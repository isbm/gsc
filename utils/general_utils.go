package gsc_utils

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"
)

var mux sync.Mutex
var rnd uint32

// RandomString generates a... random string :)
func RandomString() string {
	mux.Lock()
	r := rnd
	if r == 0 {
		r = uint32(time.Now().UnixNano() + int64(os.Getpid()))
	}
	r = r*1664525 + 1013904223
	rnd = r
	mux.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

func GetRepoFromFile(repoUrl string) (*GitPkgRepo, error) {
	content, err := ioutil.ReadFile(GIT_PKG_REPO)
	gpr := new(GitPkgRepo)

	if err != nil {
		if repoUrl == "" {
			return nil, fmt.Errorf("Package link to a Git repository was not found.")
		}

		gpr.SetIsNew()
		gpr.Url = repoUrl
		gpr.Branch = "master" // No known releases yet

		if err := ioutil.WriteFile(GIT_PKG_REPO, []byte(gpr.ToXML()), 0644); err != nil {
			return nil, err
		}

		return gpr, nil
	}

	if err := xml.Unmarshal(content, gpr); err != nil {
		return nil, err
	}

	if repoUrl != "" && repoUrl != gpr.Url {
		return nil, fmt.Errorf("Git repository is already linked to: %s", gpr.Url)
	}

	return gpr, nil
}
