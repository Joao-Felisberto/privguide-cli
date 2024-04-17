package database

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type TestScenario struct {
	StateDir string `json:"stateDir"`
	Tests    []Test `json:"tests"`
}

type Test struct {
	Query          string      `json:"query"`
	ExpectedResult interface{} `json:"expectedResult"`
}

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

func (t *Test) IsValidResult(actual interface{}) bool {
	return reflect.DeepEqual(t.ExpectedResult, actual)
}
