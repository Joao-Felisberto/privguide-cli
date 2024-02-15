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
		{"ex:ROOT", "a", "ex:_1"},
		{"ex:_1", "b", "ex:_2"},
		{"ex:_2", "c", "ex:_3"},
		{"ex:_3", "d", 1},
		{"ex:_3", "f", 1},
		{"ex:_3", "f", 2},
		{"ex:_3", "f", 3},
		{"ex:_2", "e", "ex:_7"},
		{"ex:_7", "a1", 1},
		{"ex:_7", "a2", 2},
		{"ex:_2", "e", "ex:_8"},
		{"ex:_8", "a1", 3},
		{"ex:_8", "a2", 4},
	}

	if lt, le := len(triples), len(expected); lt != le {
		t.Errorf("Number of triples generated does not match: expected %d, got %d", le, lt)
	}
	for i, v := range triples {
		found := false
		for j := 0; j < len(expected); j++ {
			if v == expected[j] {
				found = true
				break
			}
		}
		if found {
			continue
		}
		t.Errorf("The produced triples do not match @%d, expected (%s,%s,%s), got (%s,%s,%s)",
			i,
			expected[i].Subject, expected[i].Predicate, expected[i].Object,
			v.Subject, v.Predicate, v.Object,
		)
		for i, v := range triples {
			t.Logf(
				"%s,%s,%s\t%s,%s,%s",
				expected[i].Subject, expected[i].Predicate, expected[i].Object,
				v.Subject, v.Predicate, v.Object,
			)
		}
		break
	}
}
