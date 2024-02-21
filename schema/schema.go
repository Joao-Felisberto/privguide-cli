package schema

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ghodss/yaml"
	// "github.com/santhosh-tekuri/jsonschema"
	"github.com/xeipuuv/gojsonschema"
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

func convertToJSON(data interface{}) interface{} {
	switch v := data.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for key, value := range v {
			m[fmt.Sprintf("%v", key)] = convertToJSON(value)
		}
		return m
	case []interface{}:
		if len(v) == 0 {
			// Empty array represented as an array
			return []interface{}{}
		}
		var convertedArray []interface{}
		for _, value := range v {
			convertedArray = append(convertedArray, convertToJSON(value))
		}
		return convertedArray
	default:
		return data
	}
}

func ValidateYAMLAgainstSchema(yamlFile, schemaFile string) error {
	// Load JSON schema
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaFile)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		log.Fatalf("Failed to load JSON schema: %v", err)
	}

	// Load YAML data
	yamlData, err := os.ReadFile(yamlFile)
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	// Parse YAML data
	var yamlObj interface{}
	err = yaml.Unmarshal(yamlData, &yamlObj)
	if err != nil {
		log.Fatalf("Failed to parse YAML data: %v", err)
	}

	// Convert YAML data to JSON-like structure
	jsonData := convertToJSON(yamlObj)

	// Validate JSON-like data against JSON schema
	jsonLoader := gojsonschema.NewGoLoader(jsonData)
	result, err := schema.Validate(jsonLoader)
	if err != nil {
		log.Fatalf("Failed to validate YAML data: %v", err)
	}

	// Print validation result
	if result.Valid() {
		fmt.Println("YAML file abides by the schema.")
	} else {
		fmt.Println("YAML file does not abide by the schema. Validation errors:")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}

	return nil
}

/*
// ValidateYAMLAgainstSchema validates a YAML file against a JSON Schema
func ValidateYAMLAgainstSchema(yamlFile, schemaFile string) error {
	// Read YAML file
	yamlData, err := os.ReadFile(yamlFile)
	if err != nil {
		return err
	}

	// Read schema
	schemaText, err := os.ReadFile(schemaFile)
	if err != nil {
		return err
	}

	// Unmarshal YAML data
	var rawData interface{}
	if err := yaml.Unmarshal(yamlData, &rawData); err != nil {
		return err
	}
	data, err := toStringKeys(rawData)
	if err != nil {
		return err
	}

	// Load schema
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", strings.NewReader(string(schemaText))); err != nil {
		return err
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return err
	}
	if err := schema.ValidateInterface(data); err != nil {
		return err
	}

	return nil
}
*/
