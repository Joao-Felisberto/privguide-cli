package schema

import (
	"errors"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/santhosh-tekuri/jsonschema"
)

func toStringKeys(val interface{}) (interface{}, error) {
	switch val := val.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for k, v := range val {
			k, ok := k.(string)
			if !ok {
				return nil, errors.New("found non-string key")
			}
			m[k] = v
		}
		return m, nil
	case []interface{}:
		var err error
		var l = make([]interface{}, len(val))
		for i, v := range l {
			l[i], err = toStringKeys(v)
			if err != nil {
				return nil, err
			}
		}
		return l, nil
	default:
		return val, nil
	}
}

// ValidateYAMLAgainstSchema validates a YAML file against a JSON Schema
func ValidateYAMLAgainstSchema(yamlFile, schemaFile string) (bool, error) {
	// Read YAML file
	yamlData, err := os.ReadFile(yamlFile)
	if err != nil {
		return false, err
	}

	// Read schema
	schemaText, err := os.ReadFile(schemaFile)
	if err != nil {
		return false, err
	}

	// Unmarshal YAML data
	var rawData interface{}
	if err := yaml.Unmarshal(yamlData, &rawData); err != nil {
		return false, err
	}
	data, err := toStringKeys(rawData)
	if err != nil {
		return false, err
	}

	// Load schema
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", strings.NewReader(string(schemaText))); err != nil {
		return false, err
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return false, err
	}
	if err := schema.ValidateInterface(data); err != nil {
		return false, err
	}

	return true, nil
}
