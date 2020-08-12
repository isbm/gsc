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

func GetRepoFromFile(repoUrl string) (bool, string, error) {
	content, err := ioutil.ReadFile(GIT_PKG_REPO)
	if err != nil {
		if repoUrl == "" {
			return false, repoUrl, fmt.Errorf("Package link to a Git repository was not found.")
		}

		content = []byte(fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<git>\n  <url>%s</url>\n</git>\n", repoUrl))
		ioutil.WriteFile(GIT_PKG_REPO, content, 0644)
		return true, repoUrl, nil
	}

	var gitPkgRepo GitPkgRepo
	if err := xml.Unmarshal(content, &gitPkgRepo); err != nil {
		return false, "", err
	}

	if repoUrl != "" && repoUrl != gitPkgRepo.Url {
		return false, "", fmt.Errorf("Git repository is already linked to: %s", gitPkgRepo.Url)
	}

	return false, gitPkgRepo.Url, nil
}
