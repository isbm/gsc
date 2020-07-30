package gsc_clone

import (
	"os"
	"strings"
)

// FetcherProcessStream obj
type FetcherProcessStream struct {
	criterion string
}

// NewFetcherProcessStream constructor
func NewFetcherProcessStream(criterion string) *FetcherProcessStream {
	fps := new(FetcherProcessStream)
	fps.criterion = criterion
	return fps
}

// Write data to the stdout
func (fps *FetcherProcessStream) Write(data []byte) (int, error) {
	line := strings.TrimSpace(string(data))
	if strings.Contains(line, fps.criterion) {
		line = strings.ReplaceAll(line, fps.criterion, "")
		os.Stdout.WriteString(line + "\n")
	}
	return len(data), nil
}

// Close stream
func (fps *FetcherProcessStream) Close() error {
	return nil
}
