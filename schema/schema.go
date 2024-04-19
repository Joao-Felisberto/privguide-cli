package schema

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

var prefixRe = regexp.MustCompile(`^([a-zA-Z]+):(.*)`)

type YAMLVal interface{ int | string }

type Triple struct {
	Subject   string
	Predicate string
	Object    interface{}
}

func NewTriple(s, p, o string, uriBase string, uriMap *map[string]string) Triple {
	isURI := s[0] == '<' && s[len(s)-1] == '>'
	if !isURI {
		parts := strings.Split(s, "/")
		new := parts[len(parts)-1]
		matches := prefixRe.FindStringSubmatch(new)

		if len(matches) > 2 {
			prefix := matches[1]
			id := matches[2]

			uri := (*uriMap)[prefix]

			// slog.Info("New Subject!", "s", s, "uri", uri, "id", id)
			new = strings.ReplaceAll(id, " ", "_")
			new = fmt.Sprintf(`<%s/%s>`, uri, new)
			// slog.Info("NEW", "uri", new)
			s = new
		} else {
			s = fmt.Sprintf(`<%s>`, s)
		}
	}
	s = strings.ReplaceAll(s, " ", "_")

	isURI = p[0] == '<' && p[len(p)-1] == '>'
	if !isURI {
		p = strings.ReplaceAll(p, " ", "_")
		p = fmt.Sprintf(`<%s>`, p)
	}

	isURI = strings.HasPrefix(o, "https://") || strings.HasPrefix(o, "http://")
	if isURI {
		o = strings.ReplaceAll(o, " ", "_")
		parts := strings.Split(o, "/")
		tmp := parts[len(parts)-1]
		matches := prefixRe.FindStringSubmatch(tmp)

		if len(tmp) > 0 && tmp[0] == ':' {
			o = fmt.Sprintf(`<%s/%s>`, uriBase, tmp[1:])
		} else if len(matches) > 2 {
			prefix := matches[1]
			id := matches[2]

			uri := (*uriMap)[prefix]

			o = fmt.Sprintf(`<%s/%s>`, uri, id)
		} else {
			o = fmt.Sprintf(`<%s>`, o)
		}
	} else {
		matches := prefixRe.FindStringSubmatch(o)
		if len(o) > 0 && o[0] == ':' {
			o = strings.ReplaceAll(o, " ", "_")
			o = fmt.Sprintf(`<%s/%s>`, uriBase, o[1:])
		} else if len(matches) > 2 {
			prefix := matches[1]
			id := matches[2]

			uri := (*uriMap)[prefix]

			o = strings.ReplaceAll(id, " ", "_")
			o = fmt.Sprintf(`<%s/%s>`, uri, o)
		} else if o != "true" && o != "false" {
			o = fmt.Sprintf(`"%s"`, o)
		}
	}

	return Triple{s, p, o}
}

var idCounter int

func convertToJSON(data interface{}) interface{} {
	switch v := data.(type) {
	case map[interface{}]interface{}:
		m := map[string]interface{}{}
		for key, value := range v {
			m[fmt.Sprintf("%v", key)] = convertToJSON(value)
		}
		return m
	case []interface{}:
		if len(v) == 0 {
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

	if schemaFile != "" {
		res, err := ValidateYAMLAgainstSchema(yamlFile, schemaFile)
		if err != nil {
			return nil, fmt.Errorf("error validating schema: %s", err)
		}

		if !res.Valid() {
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
		return nil, fmt.Errorf("failed to load JSON schema: %s", err)
	}

	// Load YAML data
	yamlData, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %s", err)
	}

	// Parse YAML data
	var yamlObj interface{}
	err = yaml.Unmarshal(yamlData, &yamlObj)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML data: %s", err)
	}

	// Convert YAML data to JSON-like structure
	jsonData := convertToJSON(yamlObj)

	// Validate JSON-like data against JSON schema
	jsonLoader := gojsonschema.NewGoLoader(jsonData)
	result, err := schema.Validate(jsonLoader)
	if err != nil {
		return nil, fmt.Errorf("YAML file does not abide by the schema: %s", err)
	}

	return result, nil
}

func generateAnonID(uriBase string) string {
	idCounter++
	return fmt.Sprintf("%s/%d", uriBase, idCounter)
}

func YAMLtoRDF(key string, rawData interface{}, rootURI string, uriBase string, uriMap *map[string]string) []Triple {
	triples := []Triple{}
	switch data := rawData.(type) {
	case map[interface{}]interface{}:
		for key, rawValue := range data {
			// subject := tmp(uriBase, key, uriMap)
			switch value := rawValue.(type) {
			case map[interface{}]interface{}:
				// id := generateAnonID()

				id := fmt.Sprintf("%s/%v", uriBase, value["id"])
				if id == fmt.Sprintf("%s/<nil>", uriBase) {
					id = generateAnonID(uriBase)
				} else {
					delete(value, "id")
				}

				triples = append(triples, NewTriple(rootURI, fmt.Sprintf("%s/%v", uriBase, key), id, uriBase, uriMap))
				triples = append(triples, YAMLtoRDF(fmt.Sprintf("%v", key), value, id, uriBase, uriMap)...)
			case []interface{}:
				triples = append(triples, YAMLtoRDF(fmt.Sprintf("%v", key), value, rootURI, uriBase, uriMap)...)
			case int:
				tn := strconv.Itoa(value)
				triples = append(triples, NewTriple(rootURI, fmt.Sprintf("%s/%v", uriBase, key), tn, uriBase, uriMap))
			case bool:
				tn := strconv.FormatBool(value)
				triples = append(triples, NewTriple(rootURI, fmt.Sprintf("%s/%v", uriBase, key), tn, uriBase, uriMap))
			case nil:
				continue
			default: // string
				triples = append(triples, NewTriple(rootURI, fmt.Sprintf("%s/%v", uriBase, key), value.(string), uriBase, uriMap))
			}
		}
		/*
			case map[string]interface{}:
				for key, rawValue := range data {
					// subject := tmp(uriBase, key, uriMap)
					switch value := rawValue.(type) {
					case map[interface{}]interface{}:
						// id := generateAnonID()

						id := fmt.Sprintf("%s/%v", uriBase, value["id"])
						if id == fmt.Sprintf("%s/<nil>", uriBase) {
							id = generateAnonID(uriBase)
						} else {
							delete(value, "id")
						}
						// fmt.Printf("ID: %s\n", id)

						triples = append(triples, NewTriple(rootURI, fmt.Sprintf("%s/%v", uriBase, key), id, uriBase, uriMap))
						triples = append(triples, YAMLtoRDF(fmt.Sprintf("%v", key), value, id, uriBase, uriMap)...)
					case []interface{}:
						triples = append(triples, YAMLtoRDF(fmt.Sprintf("%v", key), value, rootURI, uriBase, uriMap)...)
					case int:
						tn := strconv.Itoa(value)
						triples = append(triples, NewTriple(rootURI, fmt.Sprintf("%s/%v", uriBase, key), tn, uriBase, uriMap))
					default: // string
						triples = append(triples, NewTriple(rootURI, fmt.Sprintf("%s/%v", uriBase, key), value.(string), uriBase, uriMap))
					}
				}
		*/
	case []interface{}:
		// subject := tmp(uriBase, key, uriMap)
		for _, rawElement := range data {
			// id := generateAnonID()

			switch e := rawElement.(type) {
			case []interface{}:
				// TODO: support array of arrays
				panic("Cannot have array of arrays")
			case map[interface{}]interface{}:
				id := fmt.Sprintf("%s/%v", uriBase, e["id"])
				if id == fmt.Sprintf("%s/<nil>", uriBase) {
					id = generateAnonID(uriBase)
				} else {
					delete(e, "id")
				}

				triples = append(triples, NewTriple(rootURI, fmt.Sprintf("%s/%s", uriBase, key), id, uriBase, uriMap))
				triples = append(triples, YAMLtoRDF(id, e, id, uriBase, uriMap)...)
			case int:
				eInt := strconv.Itoa(e)
				triples = append(triples, NewTriple(rootURI, fmt.Sprintf("%s/%v", uriBase, key), eInt, uriBase, uriMap))
			default: // string
				triples = append(triples, NewTriple(rootURI, fmt.Sprintf("%s/%v", uriBase, key), e.(string), uriBase, uriMap))
			}
		}
	default:
		slog.Error("Unparseable key-value pair", key, data)
	}

	return triples
}
