package database

import (
	"encoding/json"
	"fmt"
	"os"
)

// Represents a scenario on which to run tests
type TestScenario struct {
	StateDir string `json:"stateDir"` // The directory where the state is stored
	Tests    []Test `json:"tests"`    // The list of tests to run
}

// Represents a single test to be executed in a specific scenario
type Test struct {
	Query          string                   `json:"query"`          // The query to test
	ExpectedResult []map[string]interface{} `json:"expectedResult"` // The expected result for the scenario
}

// Get all test scenarios and associated tests from a file where theya re specified
//
// `file`: The file containing the test scenarios
//
// returns: the list of test scenarios and tests for each or an error if the file could not be read or parsed
func TestsFromFile(file string) ([]TestScenario, error) {
	jsonData, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %s", err)
	}

	// Unmarshal YAML data
	var tests []TestScenario
	if err := json.Unmarshal(jsonData, &tests); err != nil {
		return nil, fmt.Errorf("error reading JSON file: %s", err)
	}

	return tests, nil
}
