package cmd

import (
	"log"

	"github.com/norbix/demo4_cli_golang/internal/adapters/cli"

	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "file-filter-cli",
		Short: "CLI tool to monitor and backup files",
	}

	rootCmd.AddCommand(cli.MonitorCommand())
	rootCmd.AddCommand(cli.LogCommand())
	rootCmd.AddCommand(cli.LogFilterCommand())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
