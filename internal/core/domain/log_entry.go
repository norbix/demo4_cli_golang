package domain

import "time"

type LogEntry struct {
	Timestamp time.Time
	Action    string
	FilePath  string
}
