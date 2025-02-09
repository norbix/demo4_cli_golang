package cli

import (
	"github.com/spf13/cobra"

	"github.com/norbix/demo4_cli_golang/internal/core/service"
)

var (
	logFileName string
)

func MonitorCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "monitor",
		Short: "Start monitoring the hot folder",
		Run: func(cmd *cobra.Command, args []string) {
			service.StartMonitoring()
		},
	}
}

func LogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "View logs",
		Run: func(cmd *cobra.Command, args []string) {
			service.ViewLogs(logFileName)
		},
	}
	cmd.Flags().StringVarP(&logFileName, "file", "f", "demo4_cli_golang.out", "Log file to read")
	return cmd
}

func LogFilterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs-filter",
		Short: "Filter logs by date or filename",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			service.FilterLogs(logFileName, args[0])
		},
	}
	cmd.Flags().StringVarP(&logFileName, "file", "f", "demo4_cli_golang.out", "Log file to filter")
	return cmd
}
