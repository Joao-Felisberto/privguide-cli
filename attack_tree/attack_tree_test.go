// Tests for the attack_tree package
package attacktree_test

import (
	"os"
	"testing"

	attacktree "github.com/Joao-Felisberto/devprivops/attack_tree"
)

// Tests for the (*attacktree.AttackNode)SetExecutionStatus method
func TestSetExecutionStatus(t *testing.T) {
	node := attacktree.AttackNode{
		Description:     "",
		Query:           "",
		Children:        []*attacktree.AttackNode{},
		ExecutionStatus: attacktree.NOT_EXECUTED,
		ExecutionResult: nil,
	}

	m := []map[string]interface{}{
		{"a": 1, "b": 2},
	}
	node.SetExecutionResults(attacktree.POSSIBLE, &m)

	if node.ExecutionStatus != attacktree.POSSIBLE {
		t.Errorf("Actual execution status is not 'POSSIBLE' but '%#v'", node.ExecutionStatus)
	}

	if node.ExecutionResult != &m {
		t.Errorf("Actual execution result is not '[a:1 b:2]' but '%#v'", node.ExecutionResult)
	}
}

// Test whether attack trees can be adequately created from well formed YAML descriptions.
func TestNewAttackTreeFromYaml(t *testing.T) {
	fileData := `
description: R
query: master.rq
children:
  - description: C1 
    query: file1.rq
    children: []
  - description: C2
    query: file2.rq
    children: []
`

	if err := os.WriteFile("tmp.yml", []byte(fileData), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("tmp.yml")

	atkTree, err := attacktree.NewAttackTreeFromYaml("tmp.yml", "../.devprivops/schemas/atk-tree-schema.json")
	if err != nil {
		t.Fatal(err)
	}

	root := atkTree.Root
	if root.Description != "R" {
		t.Errorf("Description mismatch, expected 'R', got '%s'", root.Description)
	}
	if root.Query != "master.rq" {
		t.Errorf("Query mismatch, expected 'master.rq', got '%s'", root.Query)
	}
	if len(root.Children) != 2 {
		t.Errorf("Node 'R' should have 2 children")
	}

	c0 := root.Children[0]
	if c0.Description != "C1" {
		t.Errorf("Description mismatch, expected 'C1', got '%s'", c0.Description)
	}
	if c0.Query != "file1.rq" {
		t.Errorf("Query mismatch, expected 'file1.rq', got '%s'", c0.Query)
	}
	if len(c0.Children) > 0 {
		t.Errorf("Node 'C1' should not have children")
	}
	c1 := root.Children[1]
	if c1.Description != "C2" {
		t.Errorf("Description mismatch, expected 'C2', got '%s'", c1.Description)
	}
	if c1.Query != "file2.rq" {
		t.Errorf("Query mismatch, expected 'file2.rq', got '%s'", c1.Query)
	}
	if len(c1.Children) > 0 {
		t.Errorf("Node 'C2' should not have children")
	}
}

// Test whether the errors caused by creating attack trees from invalid descriptions are adequate.
// This test bypasses schema validation, which is tested in its own module and integration tests.
func TestNewAttackTreeFromInvalidYaml(t *testing.T) {
	fileData := `
description: R
query: master.rq
children:
  - description: C1 
    children: []
`

	if err := os.WriteFile("tmp.yml", []byte(fileData), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("tmp.yml")

	_, err := attacktree.NewAttackTreeFromYaml("tmp.yml", "") // No schema so invalid yaml can be passed
	if err.Error() != "error parsing child node: missing required fields in node" {
		t.Fatal(err)
	}
}

// Test whether the errors caused by creating attack trees from a description file with an array at the top level.
// This test bypasses schema validation, which is tested in its own module and integration tests.
func TestNewAttackTreeFromYamlArray(t *testing.T) {
	fileData := `
- description: R
  query: master.rq
  children:
    - description: C1 
      children: []
`

	if err := os.WriteFile("tmp.yml", []byte(fileData), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("tmp.yml")

	_, err := attacktree.NewAttackTreeFromYaml("tmp.yml", "") // No schema so invalid yaml can be passed
	if err.Error() != "invalid node data type: []interface {}" {
		t.Fatal(err)
	}
}

// Test whether the schema validation works by passing a wrong schema
func TestNewAttackTreeFromYamlInvalidSchema(t *testing.T) {
	fileData := `
description: R
query: master.rq
children:
  - description: C1 
    query: file1.rq
    children: []
  - description: C2
    query: file2.rq
    children: []
`

	invalidSchema := `
{
    "$schema": "http://json-schema.org/draft-06/schema#",
    "$ref": "#/definitions/InvalidObject",
    "definitions": {
        "InvalidObject": {
            "type": "object",
            "properties": {
                "a": {
                    "type": "integer"
                }
            },
            "required": [
                "a"
            ],
            "title": "InvalidObject"
        }
    }
}	
`

	if err := os.WriteFile("tmp.yml", []byte(fileData), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("tmp.yml")

	if err := os.WriteFile("tmp.json", []byte(invalidSchema), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("tmp.json")

	_, err := attacktree.NewAttackTreeFromYaml("tmp.yml", "./tmp.json")
	if err.Error() != "the file 'tmp.yml' does not abide by the schema: [(root): a is required]" {
		t.Fatal(err)
	}
}
