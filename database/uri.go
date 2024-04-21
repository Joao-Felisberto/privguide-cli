package database

import (
	"fmt"
	"os"
	"regexp"

	"github.com/Joao-Felisberto/devprivops/util"
	"gopkg.in/yaml.v2"
)

// Represents each URI and the metadata associated with it
type URIMetadata struct {
	Abreviation string           // The abreviated form of the URI
	URI         string           // The complete URI
	Files       []*regexp.Regexp // The regex that files in which all entities should have this uri by default match
}

// Fetches all URI information from the given file
//
// `file`: The file with the URIs
//
// returns: the list of URIs or an error when the file could not be read or parsed
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
