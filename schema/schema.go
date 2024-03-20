package schema

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	// "github.com/santhosh-tekuri/jsonschema"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

type YAMLVal interface{ int | string }

type Triple struct {
	Subject   string
	Predicate string
	Object    interface{}
}

func NewTriple(s, p, o string) Triple {
	isURI := s[0] == '<' && s[len(s)-1] == '>'
	if !isURI {
		s = fmt.Sprintf(`<%s>`, s)
	}
	s = strings.ReplaceAll(s, " ", "_")
	isURI = p[0] == '<' && p[len(p)-1] == '>'
	if !isURI {
		p = strings.ReplaceAll(p, " ", "_")
		p = fmt.Sprintf(`<%s>`, p)
	}

	isURI = strings.HasPrefix(o, "https://") || strings.HasPrefix(o, "http://")
	if isURI {
		o = fmt.Sprintf(`<%s>`, o)
	} else {
		o = fmt.Sprintf(`"%s"`, o)
	}

	return Triple{s, p, o}
}

var idCounter int

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

func ReadYAML(yamlFile string, schemaFile string) (interface{}, error) {

	yamlData, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %s", err)
	}

	// Unmarshal YAML data
	var rawData interface{}
	if err := yaml.Unmarshal(yamlData, &rawData); err != nil {
		return nil, fmt.Errorf("error reading YAML file: %s", err)
	}
	// fmt.Printf("RAWDATA %s\n", rawData)
	/*
		data, err := toStringKeys(rawData)
		if err != nil {
			return nil, err
		}
		fmt.Printf("data: %s\n", data)
	*/

	if schemaFile != "" {
		res, err := ValidateYAMLAgainstSchema(yamlFile, schemaFile)
		if err != nil {
			return nil, fmt.Errorf("error validating schema: %s", err)
		}

		if !res.Valid() {
			// fmt.Println("YAML file does not abide by the schema. Validation errors:")
			//
			//	for _, desc := range res.Errors() {
			//		fmt.Printf("- %s\n", desc)
			//	}
			return nil, fmt.Errorf("the file '%s' does not abide by the schema: %s", yamlFile, res.Errors())
		}
	}

	//return data, nil
	return rawData, nil
}

func ValidateYAMLAgainstSchema(yamlFile string, schemaFile string) (*gojsonschema.Result, error) {
	// Load JSON schema
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaFile)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		log.Fatalf("Failed to load JSON schema: %v", err)
		return nil, fmt.Errorf("failed to load JSON schema: %s", err)
	}

	// Load YAML data
	yamlData, err := os.ReadFile(yamlFile)
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
		return nil, fmt.Errorf("failed to read YAML file: %s", err)
	}

	// Parse YAML data
	var yamlObj interface{}
	err = yaml.Unmarshal(yamlData, &yamlObj)
	if err != nil {
		log.Fatalf("Failed to parse YAML data: %v", err)
		return nil, fmt.Errorf("failed to parse YAML data: %s", err)
	}

	// Convert YAML data to JSON-like structure
	jsonData := convertToJSON(yamlObj)

	// Validate JSON-like data against JSON schema
	jsonLoader := gojsonschema.NewGoLoader(jsonData)
	result, err := schema.Validate(jsonLoader)
	if err != nil {
		log.Fatalf("YAML file does not abide by the schema: %v", err)
		return nil, fmt.Errorf("YAML file does not abide by the schema: %s", err)
	}

	// Print validation result
	/*
		if result.Valid() {
			fmt.Println("YAML file abides by the schema.")
		} else {
			fmt.Println("YAML file does not abide by the schema. Validation errors:")
			for _, desc := range result.Errors() {
				fmt.Printf("- %s\n", desc)
			}
		}
	*/

	return result, nil
}

// // ValidateYAMLAgainstSchema validates a YAML file against a JSON Schema
// func ValidateYAMLAgainstSchema(yamlFile string, schemaFile string) error {
// 	// Read schema
// 	schemaText, err := os.ReadFile(schemaFile)
// 	if err != nil {
// 		return err
// 	}
//
// 	/*
// 		// Read YAML file
// 		yamlData, err := os.ReadFile(yamlFile)
// 		if err != nil {
// 			return err
// 		}
//
// 		// Unmarshal YAML data
// 		var rawData interface{}
// 		if err := yaml.Unmarshal(yamlData, &rawData); err != nil {
// 			return err
// 		}
// 		data, err := toStringKeys(rawData)
// 		if err != nil {
// 			return err
// 		}
// 	*/
// 	data, err := ReadYAML(yamlFile, "")
// 	if err != nil {
// 		return err
// 	}
//
// 	// Load schema
// 	compiler := jsonschema.NewCompiler()
// 	if err := compiler.AddResource("schema.json", strings.NewReader(string(schemaText))); err != nil {
// 		return err
// 	}
// 	schema, err := compiler.Compile("schema.json")
// 	if err != nil {
// 		return err
// 	}
// 	if err := schema.ValidateInterface(data); err != nil {
// 		return err
// 	}
//
// 	return nil
// }

func generateAnonID() string {
	idCounter++
	return fmt.Sprintf("https://example.com/%d", idCounter)
}

func YAMLtoRDF(key string, val interface{}, rootURI string) []Triple {
	triples := []Triple{}

	switch v := val.(type) {
	case map[interface{}]interface{}:
		for p, value := range v {
			switch t := value.(type) {
			case map[interface{}]interface{}:
				// id := generateAnonID()

				id := fmt.Sprintf("https://example.com/%v", t["id"])
				if id == "https://example.com/<nil>" {
					id = generateAnonID()
				} else {
					delete(t, "id")
				}
				fmt.Printf("ID: %s\n", id)

				triples = append(triples, NewTriple(rootURI, fmt.Sprintf("https://example.com/%v", p), id))
				// triples = append(triples, NewTriple(rootURI, p.(string), id))
				triples = append(triples, YAMLtoRDF(fmt.Sprintf("%v", p), t, id)...)
			case []interface{}:
				triples = append(triples, YAMLtoRDF(fmt.Sprintf("%v", p), t, rootURI)...)
			default:
				switch tv := t.(type) {
				case int:
					tn := strconv.Itoa(tv)
					triples = append(triples, NewTriple(rootURI, fmt.Sprintf("https://example.com/%v", p), tn))
				case bool:
					tn := strconv.FormatBool(tv)
					triples = append(triples, NewTriple(rootURI, fmt.Sprintf("https://example.com/%v", p), tn))
				default: // string
					triples = append(triples, NewTriple(rootURI, fmt.Sprintf("https://example.com/%v", p), tv.(string)))
				}
				// triples = append(triples, NewTriple(rootURI, fmt.Sprintf("https://example.com/%v", p), t.(string)))
			}
		}
	case map[string]interface{}:
		for p, value := range v {
			switch t := value.(type) {
			case map[interface{}]interface{}:
				// id := generateAnonID()
				id := fmt.Sprintf("https://example.com/%v", t["id"])
				if id == "https://example.com/<nil>" {
					id = generateAnonID()
				} else {
					delete(t, "id")
				}

				triples = append(triples, NewTriple(rootURI, fmt.Sprintf("https://example.com/%v", p), id))
				triples = append(triples, YAMLtoRDF(fmt.Sprintf("%v", p), t, id)...)
			case []interface{}:
				triples = append(triples, YAMLtoRDF(fmt.Sprintf("%v", p), t, rootURI)...)
			default:
				switch tv := t.(type) {
				case int:
					tn := strconv.Itoa(tv)
					triples = append(triples, NewTriple(rootURI, fmt.Sprintf("https://example.com/%v", p), tn))
				default: // string
					triples = append(triples, NewTriple(rootURI, fmt.Sprintf("https://example.com/%v", p), tv.(string)))
				}
				// triples = append(triples, NewTriple(rootURI, fmt.Sprintf("https://example.com/%v", p), t.(string)))
			}
		}
	case []interface{}:
		for _, e := range v {
			id := generateAnonID()
			switch t := e.(type) {
			case map[interface{}]interface{}, []interface{}:
				triples = append(triples, NewTriple(rootURI, fmt.Sprintf("https://example.com/%s", key), id))
				triples = append(triples, YAMLtoRDF(id, t, id)...)
			default:
				switch tv := t.(type) {
				case int:
					tn := strconv.Itoa(tv)
					triples = append(triples, NewTriple(rootURI, fmt.Sprintf("https://example.com/%v", key), tn))
				default: // string
					triples = append(triples, NewTriple(rootURI, fmt.Sprintf("https://example.com/%v", key), tv.(string)))
				}
				// triples = append(triples, NewTriple(rootURI, fmt.Sprintf("https://example.com/%s", key), t.(string)))
			}
		}
	default:
		//triples = append(triples, Triple{rootURI, key, v})
		fmt.Printf("ERROR: %s: %s\n", key, v)
	}

	return triples
}
