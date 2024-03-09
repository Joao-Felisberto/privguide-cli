package database_test

import (
	"os"
	"reflect"
	"testing"

	attacktree "github.com/Joao-Felisberto/devprivops/attack_tree"
	"github.com/Joao-Felisberto/devprivops/database"
	"github.com/Joao-Felisberto/devprivops/schema"
)

const (
	USER = "user"
	PASS = "password"
	HOST = "localhost"
	PORT = 3030
	DB   = "tmp"
)

func TestAddTriples(t *testing.T) {
	db := database.NewDBManager(USER, PASS, HOST, PORT, DB)

	code, err := db.AddTriples([]schema.Triple{
		{Subject: "<https://example.com/1>", Predicate: "<https://example.com/2>", Object: "\"1\""},
		{Subject: "<https://example.com/3>", Predicate: "<https://example.com/4>", Object: "\"2\""},
	})

	if err != nil {
		t.Fatal(err)
	}

	if code != 204 {
		t.Fatalf("Unexpected status code: %d", code)
	}
}

func TestExecuteReasonerRule(t *testing.T) {
	db := database.NewDBManager(USER, PASS, HOST, PORT, DB)

	fileData := `
	INSERT DATA {
		<https://example.com/5> <https://example.com/6> "2"
	}
`

	if err := os.WriteFile("tmp.rq", []byte(fileData), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("tmp.rq")

	if err := db.ExecuteReasonerRule("tmp.rq"); err != nil {
		t.Fatal(err)
	}
}

func TestExecuteQueryFile(t *testing.T) {
	db := database.NewDBManager(USER, PASS, HOST, PORT, DB)

	fileData := `
	SELECT * 
	WHERE {
		# ?s ?p ?o .
		<https://example.com/5> <https://example.com/6> ?o .
	} # LIMIT 2
`

	if err := os.WriteFile("tmp.rq", []byte(fileData), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("tmp.rq")

	res, err := db.ExecuteQueryFile("tmp.rq")

	if err != nil {
		t.Fatal(err)
	}

	bind := res[0]["o"].(map[string]interface{})

	if bind["value"] != "2" {
		t.Errorf("Result did not match: %s", res)
	}
}

func TestFindQueryFiles(t *testing.T) {
	if err := os.Mkdir("root", os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("root")

	if err := os.WriteFile("root/f1.rq", []byte(""), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("root/f1.rq")

	if err := os.WriteFile("root/f2.rq", []byte(""), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("root/f2.rq")

	res, err := database.FindQueryFiles("root")
	if err != nil {
		t.Fatal(err)
	}

	files := []string{"root/f1.rq", "root/f2.rq"}
	if !reflect.DeepEqual(res, files) {
		t.Errorf("Found incorrect file lists: (%s)%s (%s)%s", reflect.TypeOf(files), files, reflect.TypeOf(res), res)
	}
}

func TestExecuteAttackTree(t *testing.T) {
	db := database.NewDBManager(USER, PASS, HOST, PORT, DB)

	if err := os.Mkdir("root", os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("root")

	if err := os.WriteFile("root/f1.rq", []byte("SELECT * WHERE {<https://no.exists> ?p ?o}"), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("root/f1.rq")

	if err := os.WriteFile("root/f2.rq", []byte("SELECT * WHERE {<https://no.exists> ?p ?o}"), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("root/f2.rq")

	if err := os.WriteFile("root/f3.rq", []byte("SELECT * WHERE {<https://no.exists> ?p ?o}"), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("root/f3.rq")

	atkTreeFile := `
description: R
query: root/f1.rq
children:
  - description: C1 
    query: root/f2.rq
    children: []
  - description: C2
    query: root/f3.rq
    children: []
`

	if err := os.WriteFile("atk_tree.yml", []byte(atkTreeFile), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("atk_tree.yml")

	atkTree, err := attacktree.NewAttackTreeFromYaml("atk_tree.yml", "")
	if err != nil {
		t.Fatal(err)
	}

	res, failNode, err := db.ExecuteAttackTree(atkTree)
	if err != nil {
		t.Fatal(err)
	}
	if failNode != nil {
		t.Errorf("Failed at node %v", &failNode)
	}
	if res != nil {
		t.Errorf("Failed with result %s", res)
	}
}

func ExecuteAttackTree(s string) {
	panic("unimplemented")
}
