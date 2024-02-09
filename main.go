package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	appName := "devprivops"

	var rootCmd = &cobra.Command{
		Use:   appName,
		Short: fmt.Sprintf("A CLI application to analyze %s", appName),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("please specify a subcommand. Use '%s --help' for usage details", appName)
		},
	}

	var analyseCmd = &cobra.Command{
		Use:   "analyse <database endpoint>",
		Short: fmt.Sprintf("Analyse the specified database endpoint for %s", appName),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			databaseEndpoint := args[0]
			fmt.Printf("Analyzing database at endpoint: %s\n", databaseEndpoint)
			return nil
		},
	}

	rootCmd.AddCommand(analyseCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
