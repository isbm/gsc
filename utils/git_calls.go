package gsc_utils

import (
	"bytes"
	"fmt"
	"strings"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
)

type GitCaller struct {
	wzlib_logger.WzLogger
}

func NewGitCaller() *GitCaller {
	return new(GitCaller)
}

// Call git with specific params. All calls are blocking.
func (gitcall *GitCaller) Call(args ...string) error {
	if err := wzlib_subprocess.ExecCommand("git", args...).Run(); err != nil {
		gitcall.GetLogger().Errorf("Error calling Git: %s", err.Error())
		return err
	}
	return nil
}

// GetDefaultBranch from Git. Note, "default branch" is GitHub's terminology.
func (gitcall *GitCaller) GetDefaultBranch() string {
	cmd := wzlib_subprocess.ExecCommand("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	branchPath := strings.Split(strings.TrimSpace(out.String()), "/")
	return branchPath[len(branchPath)-1]
}

func (gitcall *GitCaller) GetCurrentBranch() string {

	cmd := wzlib_subprocess.ExecCommand("git", "symbolic-ref", "--short", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	return strings.TrimSpace(out.String())
}

// GetProjectStatus is similar to what OSC does, but this returns status for the Git repo
func (gitcall *GitCaller) GetProjectStatus() (*GSCProjectStatus, error) {
	cmd := wzlib_subprocess.ExecCommand("git", "status", "--porcelain=v1")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	stat := NewGSCProjectStatus()
	for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		fstat := strings.SplitN(line, " ", 2)
		if len(fstat) == 2 {
			fileStatus := strings.ToLower(fstat[0])
			fileName := strings.TrimSpace(fstat[1])
			fmt.Println("Status:", fileStatus, "Name:", fileName)
			switch fileStatus {
			case "m":
				stat.Modified = append(stat.Modified, fileName)
			case "d":
				stat.Deleted = append(stat.Deleted, fileName)
			case "??":
				if !strings.HasPrefix(fileName, ".osc/") {
					stat.New = append(stat.New, fileName)
				}
			default:
				gitcall.GetLogger().Errorf("Unknown Git status '%s' for file '%s'. Please report a bug", fileStatus, fileName)
			}

		}
	}

	return stat, nil
}
