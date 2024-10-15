// Tests for the schema package
package schema_test

import (
	"fmt"
	"slices"
	"testing"

	"reflect"

	"github.com/Joao-Felisberto/devprivops/schema"
	"github.com/Joao-Felisberto/devprivops/util"
	"gopkg.in/yaml.v2"
)

// Tests if a new triple can be created passing either common URIs or literals
func TestNewTriple(t *testing.T) {
	triple := schema.NewTriple(
		"http://example.com/ex1",
		"http://example.com/ex2",
		"http://example.com/ex3",
		"http://example.com",
		&map[string]string{"ex": "https://example.com"},
	)

	if triple.Subject != "<http://example.com/ex1>" {
		t.Errorf("Triple subject does not match: expected '<http://example.com/ex1>', got '%s'", triple.Subject)
	}
	if triple.Predicate != "<http://example.com/ex2>" {
		t.Errorf("Triple predicate does not match: expected '<http://example.com/ex2>', got '%s'", triple.Predicate)
	}
	if triple.Object != "<http://example.com/ex3>" {
		t.Errorf("Triple object does not match: expected '<http://example.com/ex3>', got '%s'", triple.Object)
	}

	triple = schema.NewTriple(
		"http://example.com/ex4",
		"http://example.com/ex5",
		"6",
		"http://example.com",
		&map[string]string{"ex": "http://example.com"},
	)

	if triple.Subject != "<http://example.com/ex4>" {
		t.Errorf("Triple subject does not match: expected '<http://example.com/ex1>', got '%s'", triple.Subject)
	}
	if triple.Predicate != "<http://example.com/ex5>" {
		t.Errorf("Triple predicate does not match: expected '<http://example.com/ex2>', got '%s'", triple.Predicate)
	}
	if triple.Object != "\"6\"" {
		t.Errorf("Triple object does not match: expected '<http://example.com/ex3>', got '%s'", triple.Object)
	}
}

// Tests whether a new triple can be created with ids in the form `[file default URI]:<id>`
func TestNewTripleFileId(t *testing.T) {
	triple := schema.NewTriple(
		"http://example.com/ex1",
		"http://example.com/ex2",
		"http://example.com/:b",
		"http://example.com",
		&map[string]string{"ex": "http://example.com"},
	)

	if triple.Subject != "<http://example.com/ex1>" {
		t.Errorf("Triple subject does not match: expected '<http://example.com/ex1>', got '%s'", triple.Subject)
	}
	if triple.Predicate != "<http://example.com/ex2>" {
		t.Errorf("Triple predicate does not match: expected '<http://example.com/ex2>', got '%s'", triple.Predicate)
	}
	if triple.Object != "<http://example.com/b>" {
		t.Errorf("Triple object does not match: expected '<http://example.com/b>', got '%s'", triple.Object)
	}

	triple = schema.NewTriple(
		"http://example.com/ex1",
		"http://example.com/ex2",
		":b",
		"http://example.com",
		&map[string]string{"ex": "http://example.com"},
	)

	if triple.Subject != "<http://example.com/ex1>" {
		t.Errorf("Triple subject does not match: expected '<http://example.com/ex1>', got '%s'", triple.Subject)
	}
	if triple.Predicate != "<http://example.com/ex2>" {
		t.Errorf("Triple predicate does not match: expected '<http://example.com/ex2>', got '%s'", triple.Predicate)
	}
	if triple.Object != "<http://example.com/b>" {
		t.Errorf("Triple object does not match: expected '<http://example.com/b>', got '%s'", triple.Object)
	}
}

// Tests whether a new triple can be created with ids in the form `[file uri/]abreviation:<id>`
func TestNewTripleRemoteId(t *testing.T) {
	triple := schema.NewTriple(
		"http://example.com/ot:a",
		"http://example.com/ex2",
		"http://example.com/ot:b",
		"http://example.com",
		&map[string]string{"ex": "http://example.com", "ot": "http://another.org"},
	)

	if triple.Subject != "<http://another.org/a>" {
		t.Errorf("Triple subject does not match: expected '<http://another.org/a>', got '%s'", triple.Subject)
	}
	if triple.Predicate != "<http://example.com/ex2>" {
		t.Errorf("Triple predicate does not match: expected '<http://example.com/ex2>', got '%s'", triple.Predicate)
	}
	if triple.Object != "<http://another.org/b>" {
		t.Errorf("Triple object does not match: expected '<http://another.org/b>', got '%s'", triple.Object)
	}

	triple = schema.NewTriple(
		"ot:a",
		"http://example.com/ex2",
		"ot:b",
		"http://example.com",
		&map[string]string{"ex": "http://example.com", "ot": "http://another.org"},
	)

	if triple.Subject != "<http://another.org/a>" {
		t.Errorf("Triple subject does not match: expected '<http://another.org/a>', got '%s'", triple.Subject)
	}
	if triple.Predicate != "<http://example.com/ex2>" {
		t.Errorf("Triple predicate does not match: expected '<http://example.com/ex2>', got '%s'", triple.Predicate)
	}
	if triple.Object != "<http://another.org/b>" {
		t.Errorf("Triple object does not match: expected '<http://another.org/b>', got '%s'", triple.Object)
	}
}

// Test for the `convertToJSON` function
func TestConvertToJSON(t *testing.T) {
	orig := map[interface{}]interface{}{
		"Obj1": map[string]int{
			"Value1": 1,
		},
	}

	res := schema.ConvertToJSON(orig)

	// TODO what's going on here and how do I adequately test this?
	if !reflect.DeepEqual(util.MapCast[string, interface{}](orig), res) {
		t.Errorf("Failed to compare '%#v' with '%#v'", orig, res)
	}

	orig2 := []interface{}{
		map[string]interface{}{"Obj2": 2},
		map[string]interface{}{"Obj3": 3},
	}

	res2 := schema.ConvertToJSON(orig2)

	if !reflect.DeepEqual(orig2, res2) {
		t.Errorf("Failed to compare '%s' with '%s'", orig2, res2)
	}

	orig3 := []interface{}{}

	res3 := schema.ConvertToJSON(orig3)

	if !reflect.DeepEqual(orig3, res3) {
		t.Errorf("Failed to compare '%s' with '%s'", orig3, res3)
	}
}

// Test conversion of YAML without arrays to RDF
func TestYAMLtoRDF(t *testing.T) {
	// Example YAML input with multiple addresses
	yamlInput := `
a:
  id: aId
  b:
    id: bId
    c:
      id: cId
      d: 1
      f: true
    e: 3
`

	// Parse YAML into map
	var data map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(yamlInput), &data)
	if err != nil {
		t.Fatalf("Could not parse static YAML: %s", err)
	}
	fmt.Println(data)

	// Root URI
	rootURI := "https://example.com/ROOT"

	// Convert YAML to RDF triples
	triples := schema.YAMLtoRDF(rootURI, data, rootURI, "https://example.com", &map[string]string{"ex": "https://example.com"})

	expected := []schema.Triple{
		{"<https://example.com/ROOT>", "<https://example.com/a>", "<https://example.com/aId>"},
		{"<https://example.com/aId>", "<https://example.com/b>", "<https://example.com/bId>"},
		{"<https://example.com/bId>", "<https://example.com/c>", "<https://example.com/cId>"},
		{"<https://example.com/cId>", "<https://example.com/d>", "\"1\""},
		{"<https://example.com/cId>", "<https://example.com/f>", "true"},
		{"<https://example.com/bId>", "<https://example.com/e>", "\"3\""},
	}

	if lt, le := len(triples), len(expected); lt != le {
		t.Errorf("Number of triples generated does not match: expected %d, got %d", le, lt)
	}

	for _, v := range triples {
		if !slices.Contains(expected, v) {
			t.Errorf("'%s' not expected.", v)
			break
		}
	}
}

// Test conversion of YAML arrays to RDF
func TestYAMLArrayToRDF(t *testing.T) {
	// Example YAML input with multiple addresses
	yamlInput := `
main:
  - id: e1
    a: 1
    b: 2
  - id: e2
    a: 3
    b: 4
other:
  - 10
  - 20
  - 30
`

	// Parse YAML into map
	var data map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(yamlInput), &data)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)

	// Root URI
	rootURI := "https://example.com/ROOT"

	// Convert YAML to RDF triples
	triples := schema.YAMLtoRDF(rootURI, data, rootURI, "https://example.com", &map[string]string{"ex": "https://example.com"})

	expected := []schema.Triple{
		{"<https://example.com/ROOT>", "<https://example.com/main>", "<https://example.com/e1>"},
		{"<https://example.com/ROOT>", "<https://example.com/main>", "<https://example.com/e2>"},
		{"<https://example.com/e1>", "<https://example.com/a>", "\"1\""},
		{"<https://example.com/e1>", "<https://example.com/b>", "\"2\""},
		{"<https://example.com/e2>", "<https://example.com/a>", "\"3\""},
		{"<https://example.com/e2>", "<https://example.com/b>", "\"4\""},
		{"<https://example.com/ROOT>", "<https://example.com/other>", "\"10\""},
		{"<https://example.com/ROOT>", "<https://example.com/other>", "\"20\""},
		{"<https://example.com/ROOT>", "<https://example.com/other>", "\"30\""},
	}

	if lt, le := len(triples), len(expected); lt != le {
		t.Errorf("Number of triples generated does not match: expected %d, got %d", le, lt)
	}
	for _, e := range triples {
		t.Logf("%s", e)
	}
	for _, v := range triples {
		if !slices.Contains(expected, v) {
			t.Errorf("'%s' not expected.", v)
			for i, e := range expected {
				t.Logf("%s: %s\t%s", reflect.TypeOf(e.Object), e, triples[i])
			}
			break
		}
	}
}

// Test for the ReadYAML function
func TestReadYAML(t *testing.T) {
	schemaName := ".test_read_yaml/schema1.json"
	err := util.CreateFileWithData(schemaName, `
{
    "$schema": "http://json-schema.org/draft-06/schema#",
    "$ref": "#/definitions/Welcome4",
    "definitions": {
        "Welcome4": {
            "type": "object",
            "additionalProperties": false,
            "properties": {
                "test": {
                    "type": "integer"
                },
                "another": {
                    "type": "string"
                },
                "some more": {
                    "type": "boolean"
                },
                "final": {
                    "type": "null"
                }
            },
            "required": [
                "another",
                "final",
                "some more",
                "test"
            ],
            "title": "Welcome4"
        }
    }
}
	`)
	defer util.DeleteFileAndParentPath(schemaName)
	if err != nil {
		t.Fatalf("Could not create schema '%s': %s", schemaName, err)
	}

	badSchemaName := ".test_read_yaml/badschema1.json"
	err = util.CreateFileWithData(badSchemaName, `
{
    "$schema": "http://json-schema.org/draft-06/schema#",
    "type": "array",
    "items": {},
    "definitions": {}
}
	`)
	defer util.DeleteFileAndParentPath(badSchemaName)
	if err != nil {
		t.Fatalf("Could not create schema '%s': %s", badSchemaName, err)
	}

	fileName := ".test_read_yaml/file1.yml"
	err = util.CreateFileWithData(fileName, `
test: 1
another: a
some more: true
final: null
`)
	defer util.DeleteFileAndParentPath(fileName)
	if err != nil {
		t.Fatalf("Could not create file '%s': %s", fileName, err)
	}

	data, err := schema.ReadYAML(fileName, schemaName)
	if err != nil {
		t.Errorf("Could not read YAML file '%s' with schema '%s': %s", fileName, schemaName, err)
	}

	expected := map[interface{}]interface{}{"test": 1, "another": "a", "some more": true, "final": nil}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("The data read from '%s' did not match the expected data: got '%#v', expected '%#v'", fileName, data, expected)
	}

	data, err = schema.ReadYAML(fileName, "")
	if err != nil {
		t.Errorf("Could not read YAML file '%s' WITHOUT SCHEMA: %s", fileName, err)
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("The data read from '%s' WITHOUT SCHEMA did not match the expected data: got '%#v', expected '%#v'", fileName, data, expected)
	}

	data, err = schema.ReadYAML(fileName, badSchemaName)
	if err == nil && data != nil {
		t.Errorf("Schema validation did not work for file '%s' with bad schema '%s'", fileName, badSchemaName)
	}
}

func TestValidateYAMLWithJSONSchemaString(t *testing.T) {
	schemaStr := `
{
    "$schema": "http://json-schema.org/draft-06/schema#",
    "$ref": "#/definitions/Welcome4",
    "definitions": {
        "Welcome4": {
            "type": "object",
            "additionalProperties": false,
            "properties": {
                "test": {
                    "type": "integer"
                },
                "another": {
                    "type": "string"
                },
                "some more": {
                    "type": "boolean"
                },
                "final": {
                    "type": "null"
                }
            },
            "required": [
                "another",
                "final",
                "some more",
                "test"
            ],
            "title": "Welcome4"
        }
    }
}
	`

	badSchema := `
{
    "$schema": "http://json-schema.org/draft-06/schema#",
    "type": "array",
    "items": {},
    "definitions": {}
}
	`

	fileName := ".test_read_yaml/file1.yml"
	err := util.CreateFileWithData(fileName, `
test: 1
another: a
some more: true
final: null
`)
	defer util.DeleteFileAndParentPath(fileName)
	if err != nil {
		t.Fatalf("Could not create file '%s': %s", fileName, err)
	}

	data, err := schema.ReadYAMLWithStringSchema(fileName, &schemaStr)
	if err != nil {
		t.Errorf("Could not read YAML file '%s' with schema '%s': %s", fileName, schemaStr, err)
	}

	expected := map[interface{}]interface{}{"test": 1, "another": "a", "some more": true, "final": nil}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("The data read from '%s' did not match the expected data: got '%#v', expected '%#v'", fileName, data, expected)
	}

	emptySchema := ""
	data, err = schema.ReadYAMLWithStringSchema(fileName, &emptySchema)
	if err != nil {
		t.Errorf("Could not read YAML file '%s' WITHOUT SCHEMA: %s", fileName, err)
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("The data read from '%s' WITHOUT SCHEMA did not match the expected data: got '%#v', expected '%#v'", fileName, data, expected)
	}

	data, err = schema.ReadYAMLWithStringSchema(fileName, &badSchema)
	if err == nil && data != nil {
		t.Errorf("Schema validation did not work for file '%s' with bad schema '%s'", fileName, badSchema)
	}
}
