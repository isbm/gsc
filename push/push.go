package gsc_push

import (
	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
)

/*
	Implements "push" command to the OBS
*/

// GSCPush class
type GSCPush struct {
	origin string
	wzlib_logger.WzLogger
}

// NewGCSPush creates a package cloning tool
func NewGCSPush() *GSCPush {
	push := new(GSCPush).SetOrigin("master")
	return push
}

func (push *GSCPush) toOrigin() error {
	cmd, err := wzlib_subprocess.BufferedExec("git", "push", "-u", "origin", push.origin)
	cmd.Wait()
	return err
}

// SetOrigin where to push the sources
func (push *GSCPush) SetOrigin(origin string) *GSCPush {
	push.origin = origin
	return push
}

// Push all the sources to the Git and OBS, synchronising among
func (push *GSCPush) Push() error {
	return push.toOrigin()
}
