package gsc_clone

import (
	"path"
	"strings"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
)

// FetcherProcessStream obj
type FetcherProcessStream struct {
	criterion string
	prefix    string
	wzlib_logger.WzLogger
}

// NewFetcherProcessStream constructor
func NewFetcherProcessStream(criterion string) *FetcherProcessStream {
	fps := new(FetcherProcessStream)
	fps.criterion = criterion
	return fps
}

// SetPrefix to the each output line
func (fps *FetcherProcessStream) SetPrefix(prefix string) *FetcherProcessStream {
	fps.prefix = prefix + " "
	return fps
}

// Filter out data
func (fps *FetcherProcessStream) Filter(line string) string {
	words := strings.Split(strings.ReplaceAll(line, fps.criterion, ""), " ")
	line = path.Base(words[len(words)-1])
	if line == "." {
		line = ""
	}
	return line
}

// Write data to the stdout
func (fps *FetcherProcessStream) Write(data []byte) (int, error) {
	line := strings.TrimSpace(string(data))
	if strings.Contains(line, fps.criterion) {
		line = fps.Filter(line)
		if line != "" {
			fps.GetLogger().Info(fps.prefix + line)
		}
	}
	return len(data), nil
}

// Close stream
func (fps *FetcherProcessStream) Close() error {
	return nil
}
