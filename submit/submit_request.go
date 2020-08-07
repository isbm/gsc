package gsc_submit

import (
	"strings"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	gsc_utils "github.com/isbm/gsc/utils"
)

/*
	Submit request for the package acceptance
*/

type GSCSubmitRequest struct {
	git   *gsc_utils.GitCaller
	chlog *gsc_utils.GSCChangeLog

	wzlib_logger.WzLogger
}

// Create obj
func NewGSCSubmitRequest() *GSCSubmitRequest {
	sr := new(GSCSubmitRequest)
	sr.git = gsc_utils.NewGitCaller()
	sr.chlog = gsc_utils.NewGSCChangeLog()
	return sr
}

// Add to the sequence the whole changelog entry and reset WIP status
func (sr *GSCSubmitRequest) commitChangelog() error {
	err := gsc_utils.CallWithTTY("osc", "commit", sr.chlog.GetFilename())
	if err != nil {
		return err
	}

	sr.git.Call("add", sr.chlog.GetFilename())
	sr.git.Call("commit", "-m", "Set Changelog message")
	sr.git.Call("push", "--set-upstream", "origin", sr.git.GetCurrentBranch())
	return nil
}

// Format changelog entry and reset WIP status
func (sr *GSCSubmitRequest) formatChangelog() error {
	entry := sr.chlog.GetLatest()
	out := []string{}
	for _, line := range entry.Messages {
		if strings.HasPrefix(line, "#") && !strings.Contains(line, "WIP ENTRY") {
			out = append(out, strings.TrimSpace(strings.SplitN(line, "#", 2)[1]))
		}
	}
	entry.Messages = out
	if err := sr.chlog.Write(); err != nil {
		return err
	}

	return gsc_utils.CallWithTTY("osc", "vc", "-e")
}

// Submit request
func (sr *GSCSubmitRequest) Submit() error {
	sr.formatChangelog()
	sr.commitChangelog()

	return gsc_utils.CallWithTTY("osc", "sr", "--yes")
}
