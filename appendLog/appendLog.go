package appendLog

import (
	"time"
	uuidv7 "github.com/gofrs/uuid/v5"
	"strings"
	"errors"
)

type LogEntry struct {
	Id uuidv7.UUID
	Type string
	Date time.Time
	Hostname string
	Source string
	Content string
}

// encoder function to encode in one line
func (l *LogEntry) Encode() string {
	return l.Id.String() + " " + l.Type + " " + l.Date.Format(time.RFC3339) + " " + l.Hostname + " " + l.Source + " " + l.Content
}

// decoder function to decode from one line
func Decode(s string) (*LogEntry, error) {
	l := &LogEntry{}
	parts := strings.SplitN(s, " ", 6)
	if len(parts) != 6 {
		return nil, errors.New("invalid log entry format")
	}
	id, err := uuidv7.FromString(parts[0])
	if err != nil {
		return nil, err
	}
	l.Id = id
	l.Type = parts[1]
	date, err := time.Parse(time.RFC3339, parts[2])
	if err != nil {
		return nil, err
	}
	l.Date = date
	l.Hostname = parts[3]
	l.Source = parts[4]
	l.Content = parts[5]
	return l, nil
}
