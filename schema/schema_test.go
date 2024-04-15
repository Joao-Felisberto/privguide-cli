package schema_test

import (
	"fmt"
	"slices"
	"testing"

	"reflect"

	"github.com/Joao-Felisberto/devprivops/schema"
	"gopkg.in/yaml.v2"
)

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
		&map[string]string{"ex": "https://example.com"},
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

func TestConvertToJSON(t *testing.T) {
	orig := map[string]interface{}{
		"Obj": map[string]interface{}{
			"Value1": 1,
			"Value2": 2,
		},
	}

	res := schema.ConvertToJSON(orig).(map[string]interface{})

	if !reflect.DeepEqual(orig, res) {
		t.Errorf("Failed to compare '%s' with '%s'", orig, res)
	}
}

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
      f: 2
    e: 3
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
		{"<https://example.com/ROOT>", "<https://example.com/a>", "<https://example.com/aId>"},
		{"<https://example.com/aId>", "<https://example.com/b>", "<https://example.com/bId>"},
		{"<https://example.com/bId>", "<https://example.com/c>", "<https://example.com/cId>"},
		{"<https://example.com/cId>", "<https://example.com/d>", "\"1\""},
		{"<https://example.com/cId>", "<https://example.com/f>", "\"2\""},
		{"<https://example.com/bId>", "<https://example.com/e>", "\"3\""},
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
