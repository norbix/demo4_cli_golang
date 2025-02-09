package ports

import (
	"github.com/norbix/demo4_cli_golang/internal/core/domain"
)

type LogRepository interface {
	SaveLog(entry domain.LogEntry) error
	ReadLogs() ([]domain.LogEntry, error)
	FilterLogs(filter string) ([]domain.LogEntry, error)
}
