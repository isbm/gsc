package gsc_utils

import (
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
