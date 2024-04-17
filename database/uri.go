package database

import (
	"fmt"
	"os"
	"regexp"

	"github.com/Joao-Felisberto/devprivops/util"
	"gopkg.in/yaml.v2"
)

type URIMetadata struct {
	Abreviation string
	URI         string
	Files       []*regexp.Regexp
}

func URIsFromFile(file string) (*[]URIMetadata, error) {
	rawData, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file '%s': %s", file, err)
	}

	// Unmarshal YAML data
	var URIstruct []struct {
		Abreviation string
		URI         string
		Files       []string
	}
	if err := yaml.Unmarshal(rawData, &URIstruct); err != nil {
		return nil, fmt.Errorf("error reading YAML file '%s': %s", file, err)
	}

	res := util.Map(URIstruct, func(uri struct {
		Abreviation string
		URI         string
		Files       []string
	}) URIMetadata {
		return URIMetadata{
			uri.Abreviation,
			uri.URI,
			util.Map(uri.Files, func(regex string) *regexp.Regexp {
				return regexp.MustCompile(regex)
			}),
		}
	})

	return &res, nil
}
