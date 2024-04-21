package cmd

import (
	"fmt"
	"log/slog"
	"reflect"
	"strconv"

	"github.com/Joao-Felisberto/devprivops/database"
	"github.com/Joao-Felisberto/devprivops/fs"
	"github.com/spf13/cobra"
)

// Executes all tests for a given scenario
//
// `dbManager`: The DBManager connecting to the database
//
// `scenario`: The scenario whose tests are to be executed
//
// returns: error if tehre was any error reading files or validating their schemas, connecting to the database or executing queries
func runScenario(dbManager *database.DBManager, scenario database.TestScenario) error {
	dbManager.CleanDB()
	slog.Info("Loading scenario", "scenario", scenario.StateDir)

	err := loadRepresentations(dbManager, scenario.StateDir)
	if err != nil {
		return err
	}

	for _, t := range scenario.Tests {
		slog.Info("Running test", "test", t.Query)
		file, err := fs.GetFile(t.Query)
		if err != nil {
			return fmt.Errorf("error reading test file '%s': %s", t.Query, err)
		}
		res, err := dbManager.ExecuteQueryFile(file)
		if err != nil {
			return fmt.Errorf("error running test '%s': %s", file, err)
		}

		if !reflect.DeepEqual(t.ExpectedResult, res) {
			return fmt.Errorf("result of '%s' does not match expectations: got '%v', expected '%v'", file, res, t.ExpectedResult)
		}
	}

	slog.Info("All tests passed!", "scenario", scenario.StateDir)
	return nil
}

// Main entry point to the test command.
// Run all the tests of each scenario.
//
// `cmd`: The cobra command
//
// `args`: The args of said command
//
// returns: an error when reading any of the scenarios fails
func Test(cmd *cobra.Command, args []string) error {
	username := args[0]
	password := args[1]
	ip := args[2]
	port, err := strconv.Atoi(args[3])
	if err != nil {
		return err
	}
	dataset := args[4]

	dbManager := database.NewDBManager(
		username,
		password,
		ip,
		port,
		dataset,
	)

	// 1. Load data

	testFile, err := fs.GetFile("tests/spec.json")
	if err != nil {
		return err
	}

	tests, err := database.TestsFromFile(testFile)
	if err != nil {
		return err
	}

	// 2. For each scenario, run the tests
	for _, t := range tests {
		err := runScenario(&dbManager, t)
		if err != nil {
			return fmt.Errorf("test failed for scenario '%s': %s", t.StateDir, err)
		}
	}

	return nil
}
