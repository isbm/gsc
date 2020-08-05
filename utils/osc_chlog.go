package gsc_utils

import (
	"fmt"
	"io/ioutil"
	"net/mail"
	"strings"
	"time"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
)

/*
	Operations with the changelog
*/
type ChangeLogEntry struct {
	Date     time.Time
	Address  *mail.Address
	Messages []string
}

type GSCChangeLog struct {
	entries  []*ChangeLogEntry
	filename string

	wzlib_logger.WzLogger
}

const CHLOG_SEP = "-------------------------------------------------------------------"

func NewGSCChangeLog() *GSCChangeLog {
	cl := new(GSCChangeLog).parse()
	return cl
}

func (cl *GSCChangeLog) parse() *GSCChangeLog {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		cl.GetLogger().Panicf("Unable to read current directory: %s", err.Error())
	}

	for _, fname := range files {
		if strings.HasSuffix(fname.Name(), ".changes") {
			cl.filename = fname.Name()
			break
		}
	}

	if cl.filename == "" {
		cl.GetLogger().Panic("Changelog file was not found")
	}

	data, err := ioutil.ReadFile(cl.filename)
	if err != nil {
		cl.GetLogger().Panicf("Unable to read changelog file at %s: %s", cl.filename, err.Error())
	}

	for _, entry := range strings.Split(string(data), CHLOG_SEP) {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		// parse entry
		lines := strings.Split(entry, "\n")

		dateAddrLine := strings.SplitN(lines[0], "-", 2)
		clEntry := new(ChangeLogEntry)
		clEntry.Date, err = time.Parse("Mon Jan 2 15:04:05 MST 2006", strings.TrimSpace(dateAddrLine[0]))
		if err != nil {
			cl.GetLogger().Panicf("Unable to parse date/time: %s", dateAddrLine[0])
		}
		clEntry.Address, err = mail.ParseAddress(strings.TrimSpace(dateAddrLine[1]))
		if err != nil {
			cl.GetLogger().Panicf("Unable to parse address: %s", dateAddrLine[1])
		}
		// Messages
		for _, msg := range lines[1:] {
			msg = strings.TrimSpace(msg)
			if msg != "" {
				clEntry.Messages = append(clEntry.Messages, msg)
			}
		}
		cl.entries = append(cl.entries, clEntry)

	}
	// Reverse
	for i, j := 0, len(cl.entries)-1; i < j; i, j = i+1, j-1 {
		cl.entries[i], cl.entries[j] = cl.entries[j], cl.entries[i]
	}
	return cl
}

func (cl *GSCChangeLog) AddEntry(entry *ChangeLogEntry) {
	cl.entries = append(cl.entries, entry)
}

// GetLatest changelog entry
func (cl *GSCChangeLog) GetLatest() *ChangeLogEntry {
	if len(cl.entries) > 0 {
		return cl.entries[len(cl.entries)-1]
	}
	return nil
}

// GetAll changelog entries
func (cl *GSCChangeLog) GetAll() []*ChangeLogEntry {
	return cl.entries
}

// Dump changelog back to the file
func (cl *GSCChangeLog) Dump() {
	for i := len(cl.entries); i > 0; i-- {
		entry := cl.entries[i-1]
		fmt.Println(CHLOG_SEP)
		ts := entry.Date.Format("Mon Jan 2 15:04:05 MST 2006") // Incompatible with buggy SUSE's formatting!
		fmt.Printf("%s - %s <%s>\n\n", ts,
			entry.Address.Name, entry.Address.Address)
		for _, msg := range entry.Messages {
			fmt.Printf("%s\n", msg)
		}
		fmt.Println("")
	}

}
