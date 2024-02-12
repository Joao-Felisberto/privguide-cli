package schema_test

import (
	"fmt"
	"testing"

	"github.com/Joao-Felisberto/devprivops/schema"
	"gopkg.in/yaml.v2"
)

func TestYAMLtoRDF(t *testing.T) {
	// Example YAML input with multiple addresses
	yamlInput := `
a:
  b:
    c:
      d: 1
      f:
        - 1
        - 2
        - 3
    e:
      - a1: 1
        a2: 2
      - a1: 3
        a2: 4
`

	// Parse YAML into map
	var data map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(yamlInput), &data)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)

	// Root URI
	rootURI := "ex:ROOT"

	// Convert YAML to RDF triples
	triples := schema.YAMLtoRDF(rootURI, data, rootURI)

	expected := []schema.Triple{
		{Subject: "ex:ROOT", Predicate: "a", Object: "ex:_1"},
		{Subject: "ex:_1", Predicate: "b", Object: "ex:_2"},
		{Subject: "ex:_2", Predicate: "c", Object: "ex:_3"},
		{Subject: "ex:_3", Predicate: "d", Object: 1},
		{Subject: "ex:_3", Predicate: "f", Object: 1},
		{Subject: "ex:_3", Predicate: "f", Object: 2},
		{Subject: "ex:_3", Predicate: "f", Object: 3},
		{Subject: "ex:_2", Predicate: "e", Object: "ex:_7"},
		{Subject: "ex:_7", Predicate: "a1", Object: 1},
		{Subject: "ex:_7", Predicate: "a2", Object: 2},
		{Subject: "ex:_2", Predicate: "e", Object: "ex:_8"},
		{Subject: "ex:_8", Predicate: "a1", Object: 3},
		{Subject: "ex:_8", Predicate: "a2", Object: 4},
	}

	if lt, le := len(triples), len(expected); lt != le {
		t.Errorf("Number of triples generated does not match: expected %d, got %d", le, lt)
	}
	for i, v := range triples {
		if v != expected[i] {
			t.Errorf("The produced triples do not match @%d, expected (%s,%s,%s), got (%s,%s,%s)",
				i,
				expected[i].Subject, expected[i].Predicate, expected[i].Object,
				v.Subject, v.Predicate, v.Object,
			)
			break
		}
	}
}
