// Program entry point
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Joao-Felisberto/devprivops/cmd"
	"github.com/Joao-Felisberto/devprivops/fs"
	"github.com/Joao-Felisberto/devprivops/util"
	"github.com/spf13/cobra"
)

var verbose = false

// var logLevel = slog.LevelDebug

// Builds the command and delegates execution to the appropriate function from the cmd package
func main() {
	// slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	var rootCmd = &cobra.Command{
		Use:   util.AppName,
		Short: fmt.Sprintf("A CLI application to analyze %s", util.AppName),
		RunE: func(cmd_ *cobra.Command, args []string) error {
			return fmt.Errorf("please specify a subcommand. Use '%s --help' for usage details", util.AppName)
		},
	}

	var analyseCmd = &cobra.Command{
		Use:   "analyse <username> <password> <database ip> <database port> <dataset>",
		Short: fmt.Sprintf("Analyse the specified database endpoint for %s", util.AppName),
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd_ *cobra.Command, args []string) error {
			logLevel := slog.LevelInfo
			if verbose {
				logLevel = slog.LevelDebug
			}
			util.SetupLogger(logLevel)
			return cmd.Analyse(cmd_, args)
		},
	}

	var testCmd = &cobra.Command{
		Use:   "test",
		Short: "Tests the queries against user-defined scenarios",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd_ *cobra.Command, args []string) error {
			logLevel := slog.LevelInfo
			if verbose {
				logLevel = slog.LevelDebug
			}
			util.SetupLogger(logLevel)
			return cmd.Test(cmd_, args)
		},
	}

	analyseCmd.Flags().StringVar(&util.ReportEndpoint, "report-endpoint", "", "Endpoint where to send the final report")

	analyseCmd.Flags().BoolVar(&util.Pipeline, "pipeline", false, "whether to format the output for pipeline usage")
	testCmd.Flags().BoolVar(&util.Pipeline, "pipeline", false, "whether to format the output for pipeline usage")

	analyseCmd.Flags().StringVar(&fs.GlobalDir, "global-dir", fmt.Sprintf("/etc/%s", util.AppName), "The path to the global configurations")
	testCmd.Flags().StringVar(&fs.GlobalDir, "global-dir", fmt.Sprintf("/etc/%s", util.AppName), "The path to the global configurations")

	analyseCmd.Flags().StringVar(&fs.LocalDir, "local-dir", fmt.Sprintf("./.%s", util.AppName), "The path to the local configurations")
	testCmd.Flags().StringVar(&fs.LocalDir, "local-dir", fmt.Sprintf("./.%s", util.AppName), "The path to the local configurations")

	analyseCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "whether to display debug messages")
	testCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "whether to display debug messages")

	rootCmd.AddCommand(analyseCmd)
	rootCmd.AddCommand(testCmd)

	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
