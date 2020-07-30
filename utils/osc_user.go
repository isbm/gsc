package gsc_utils

import (
	"fmt"
	"strings"

	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
)

type OSCUser struct {
	Uid   string
	Name  string
	Email string
}

func GetOSCUser() (*OSCUser, error) {
	cmd, err := wzlib_subprocess.BufferedExec("osc", "user")
	if err != nil {
		return nil, err
	}

	stdout := strings.TrimSpace(cmd.StdoutString())
	cmd.Wait()

	if stdout == "" {
		return nil, fmt.Errorf("Unable to get OSC user. Have you setup 'osc' correctly?")
	}

	buff := strings.Split(stdout, ":")
	user := &OSCUser{
		Uid:   buff[0],
		Name:  strings.TrimSpace(strings.Split(buff[1], "<")[0]),
		Email: strings.TrimSpace(strings.Split(buff[1], "<")[1]),
	}

	return user, nil
}
