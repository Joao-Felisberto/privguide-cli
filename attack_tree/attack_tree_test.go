package attacktree_test

import (
	"os"
	"testing"

	attacktree "github.com/Joao-Felisberto/devprivops/attack_tree"
)

func TestSetExecutionStatus(t *testing.T) {
	node := attacktree.AttackNode{
		Description:     "",
		Query:           "",
		Children:        []*attacktree.AttackNode{},
		ExecutionStatus: attacktree.NOT_EXECUTED,
	}

	node.SetExecutionStatus(attacktree.POSSIBLE)

	if node.ExecutionStatus != attacktree.POSSIBLE {
		t.Errorf("Actual execution status is not 'POSSIBLE' but '%#v'", node.ExecutionStatus)
	}
}

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
