package gsc_add

import (
	wzlib_logger "github.com/infra-whizz/wzlib/logger"
)

/*
	Implements "add" command
*/

// GSCAdd class
type GSCAdd struct {
	mask string
	wzlib_logger.WzLogger
}

// NewGSCAdd constructs new instance of a GCSAdd class
func NewGSCAdd() *GSCAdd {
	return new(GSCAdd).SetMask("*")
}

// SetMask sets a mask of files to add
func (add *GSCAdd) SetMask(mask string) *GSCAdd {
	add.mask = mask
	return add
}

func (add *GSCAdd) Add() error {
	return nil
}
