//go:build profile

// Program entry point
package main

import (
	"fmt"
	"log/slog"
	"os"

	"runtime"
	"runtime/pprof"

	"github.com/Joao-Felisberto/devprivops/cmd"
	"github.com/Joao-Felisberto/devprivops/fs"
	"github.com/Joao-Felisberto/devprivops/schema"
	"github.com/Joao-Felisberto/devprivops/util"
	"github.com/spf13/cobra"
)

var verbose = false // Whether the log should log more information or not
var profile = ""

// Builds the command and delegates execution to the appropriate function from the cmd package
func main() {
	/*
		cpuProfileFile, err := os.Create("cpu")
		if err != nil {
			fmt.Println(err)
		}
		memoryProfileFile, err := os.Create("mem")
		if err != nil {
			fmt.Println(err)
		}

		pprof.StartCPUProfile(cpuProfileFile)
		defer pprof.StopCPUProfile()
	*/

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
			memoryProfileFile, err := os.Create("mem.prof")
			if profile == "cpu" {
				cpuProfileFile, err := os.Create("cpu.prof")
				if err != nil {
					return err
				}
				pprof.StartCPUProfile(cpuProfileFile)
			} else if profile == "mem" {
				// memoryProfileFile, err := os.Create("mem.prof")
				if err != nil {
					return err
				}
			}

			logLevel := slog.LevelInfo
			if verbose {
				logLevel = slog.LevelDebug
			}
			util.SetupLogger(logLevel)
			res := cmd.Analyse(cmd_, args)

			if profile == "cpu" {
				pprof.StopCPUProfile()
			} else if profile == "mem" {
				pprof.WriteHeapProfile(memoryProfileFile)
			}

			return res
		},
	}

	var testCmd = &cobra.Command{
		Use:   "test <username> <password> <database ip> <database port> <dataset>",
		Short: "Tests the queries against user-defined scenarios",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd_ *cobra.Command, args []string) error {
			// memoryProfileFile, err := os.Create("mem.prof")
			if profile == "cpu" {
				runtime.SetCPUProfileRate(400)
				cpuProfileFile, err := os.Create("cpu.prof")
				if err != nil {
					return err
				}
				pprof.StartCPUProfile(cpuProfileFile)
			} /*else if profile == "mem" {
				memoryProfileFile, err := os.Create("mem.prof")
				if err != nil {
					return err
				}
			}*/

			logLevel := slog.LevelInfo
			if verbose {
				logLevel = slog.LevelDebug
			}
			util.SetupLogger(logLevel)
			res := cmd.Test(cmd_, args)

			if profile == "cpu" {
				pprof.StopCPUProfile()
			} else if profile == "mem" {
				memoryProfileFile, err := os.Create("mem.prof")
				if err != nil {
					return err
				}
				pprof.WriteHeapProfile(memoryProfileFile)
			}

			return res
		},
	}

	var schemaCmd = &cobra.Command{
		Use:   "schema <attack-tree|query|report|requirement>",
		Short: "Prints the internal json schema",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd_ *cobra.Command, args []string) error {
			logLevel := slog.LevelInfo
			if verbose {
				logLevel = slog.LevelDebug
			}
			util.SetupLogger(logLevel)
			schemaName := args[0]
			switch schemaName {
			case "attack-tree":
				fmt.Print(schema.ATK_TREE_SCHEMA)
			case "query":
				fmt.Print(schema.QUERY_SCHEMA)
			case "report":
				fmt.Print(schema.REPORT_DATA_SCHEMA)
			case "requirement":
				fmt.Print(schema.REQUIREMENT_SCHEMA)
			default:
				return fmt.Errorf("schema '%s' not found, valid possibilities: [attack-tree, query, report, requirement]", schemaName)
			}
			return nil
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

	analyseCmd.Flags().StringVar(&profile, "profile", "", "What to profile")
	testCmd.Flags().StringVar(&profile, "profile", "", "What to profile")

	rootCmd.AddCommand(analyseCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(schemaCmd)

	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
