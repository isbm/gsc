package gsc_utils

import (
	"strings"

	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
)

type LogEntry struct {
	Revision string
	User     string
	Date     string
	Checksum string
	Message  string
}

func (le *LogEntry) Parse(lines []string) *LogEntry {
	if !strings.Contains(lines[0], "|") {
		panic("Unparse-able log entry.")
	}
	data := strings.Split(lines[0], "|")
	le.Revision = strings.TrimSpace(data[0])
	le.User = strings.TrimSpace(data[1])
	le.Date = strings.TrimSpace(data[2])
	le.Checksum = strings.TrimSpace(data[3])
	le.Message = strings.TrimSpace(strings.Join(lines[1:len(lines)-1], "\n"))

	return le
}

type GSCRevisionLog struct {
	entries []*LogEntry
}

// NewGSCRevisionLog constructor
func NewGSCRevisionLog() *GSCRevisionLog {
	return new(GSCRevisionLog).parse()
}

func (revlog *GSCRevisionLog) parse() *GSCRevisionLog {
	cmd, _ := wzlib_subprocess.BufferedExec("osc", "log")
	out := cmd.StdoutString()
	cmd.Wait()
	revlog.entries = []*LogEntry{}

	entry := []string{}
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "|") && len(strings.Split(line, "|")) == 6 {
			if len(entry) > 0 {
				if len(entry) > 3 {
					revlog.entries = append(revlog.entries, new(LogEntry).Parse(entry))
				}
				entry = []string{}
			} else {
				entry = append(entry, line)
			}
		}
		entry = append(entry, line)
	}
	revlog.entries = append(revlog.entries, new(LogEntry).Parse(entry))

	return revlog
}

// GetLatest revision log entry
func (revlog *GSCRevisionLog) GetLatest() *LogEntry {
	if len(revlog.entries) > 1 {
		return revlog.entries[0]
	}
	return nil
}
