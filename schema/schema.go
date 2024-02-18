package schema

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/santhosh-tekuri/jsonschema"
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

	isURI = strings.HasPrefix(o, "https://")
	if isURI {
		o = fmt.Sprintf(`<%s>`, o)
	} else {
		o = fmt.Sprintf(`"%s"`, o)
	}

	return Triple{s, p, o}
}

var idCounter int

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
		return nil, err
	}

	// Unmarshal YAML data
	var rawData interface{}
	if err := yaml.Unmarshal(yamlData, &rawData); err != nil {
		return nil, err
	}
	data, err := toStringKeys(rawData)
	if err != nil {
		return nil, err
	}

	if schemaFile != "" {
		schemaText, err := os.ReadFile(schemaFile)
		if err != nil {
			return nil, err
		}

		compiler := jsonschema.NewCompiler()
		if err := compiler.AddResource("schema.json", strings.NewReader(string(schemaText))); err != nil {
			return nil, err
		}
		schema, err := compiler.Compile("schema.json")
		if err != nil {
			return nil, err
		}
		if err := schema.ValidateInterface(data); err != nil {
			return nil, err
		}
	}

	return data, nil
}

// ValidateYAMLAgainstSchema validates a YAML file against a JSON Schema
func ValidateYAMLAgainstSchema(yamlFile string, schemaFile string) error {
	// Read schema
	schemaText, err := os.ReadFile(schemaFile)
	if err != nil {
		return err
	}

	/*
		// Read YAML file
		yamlData, err := os.ReadFile(yamlFile)
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
	*/
	data, err := ReadYAML(yamlFile, "")
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
				id := generateAnonID()

				triples = append(triples, NewTriple(rootURI, fmt.Sprintf("https://example.com/%v", p), id))
				triples = append(triples, YAMLtoRDF(fmt.Sprintf("%v", p), t, id)...)
			case []interface{}:
				triples = append(triples, YAMLtoRDF(fmt.Sprintf("%v", p), t, rootURI)...)
			default:
				triples = append(triples, NewTriple(rootURI, fmt.Sprintf("https://example.com/%v", p), t.(string)))
			}
		}
	case map[string]interface{}:
		for p, value := range v {
			switch t := value.(type) {
			case map[interface{}]interface{}:
				id := generateAnonID()
				triples = append(triples, NewTriple(rootURI, fmt.Sprintf("https://example.com/%v", p), id))
				triples = append(triples, YAMLtoRDF(fmt.Sprintf("%v", p), t, id)...)
			case []interface{}:
				triples = append(triples, YAMLtoRDF(fmt.Sprintf("%v", p), t, rootURI)...)
			default:
				triples = append(triples, NewTriple(rootURI, fmt.Sprintf("https://example.com/%v", p), t.(string)))
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
				triples = append(triples, NewTriple(rootURI, fmt.Sprintf("https://example.com/%s", key), t.(string)))
			}
		}
	default:
		//triples = append(triples, Triple{rootURI, key, v})
		fmt.Printf("ERROR: %s: %s\n", key, v)
	}

	return triples
}
