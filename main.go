// Program entry point
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Joao-Felisberto/devprivops/cmd"
	"github.com/spf13/cobra"
)

// const REPORT_ENDPOINT_FLAG_NAME = "report-endpoint"

// Builds the command and delegates execution to the appropriate function from the cmd package
func main() {
	appName := "devprivops"
	reportEndpoint := ""

	var rootCmd = &cobra.Command{
		Use:   appName,
		Short: fmt.Sprintf("A CLI application to analyze %s", appName),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("please specify a subcommand. Use '%s --help' for usage details", appName)
		},
	}

	var analyseCmd = &cobra.Command{
		Use:   "analyse <username> <password> <database ip> <database port> <dataset>",
		Short: fmt.Sprintf("Analyse the specified database endpoint for %s", appName),
		Args:  cobra.ExactArgs(5),
		RunE:  cmd.Analyse,
	}

	var testCmd = &cobra.Command{
		Use:   "test",
		Short: "Tests the queries against user-defined scenarios",
		Args:  cobra.ExactArgs(5),
		RunE:  cmd.Test,
	}

	analyseCmd.Flags().StringVar(&reportEndpoint, "report-endpoint", "", "Endpoint where to send the final report")

	rootCmd.AddCommand(analyseCmd)
	rootCmd.AddCommand(testCmd)

	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
