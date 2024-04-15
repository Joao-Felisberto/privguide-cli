package database

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type URIMetadata struct {
	Abreviation string
	URI         string
	Files       []string
}

func FromFile(file string) (*[]URIMetadata, error) {
	rawData, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file '%s': %s", file, err)
	}

	// Unmarshal YAML data
	var res []URIMetadata
	if err := yaml.Unmarshal(rawData, &res); err != nil {
		return nil, fmt.Errorf("error reading YAML file '%s': %s", file, err)
	}

	return &res, nil
}
