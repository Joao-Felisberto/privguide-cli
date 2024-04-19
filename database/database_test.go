package database_test

import (
	"encoding/json"
	"io"
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

type BindingVar struct {
	Type     string `json:"type"`
	DataType string `json:"datatype"`
	Value    string `json:"value"`
}

type HeadVars struct {
	Vars []string `json:"vars"`
}

type ResultBindings struct {
	Bindings []CountBinding `json:"bindings"`
}

type CountBinding struct {
	Cnt BindingVar `json:"cnt"`
}

type CountResult struct {
	Head    HeadVars       `json:"head"`
	Results ResultBindings `json:"results"`
}

func TestCleanDB(t *testing.T) {
	db := database.NewDBManager(USER, PASS, HOST, PORT, DB)
	db.CleanDB()

	code, err := db.AddTriples([]schema.Triple{
		{Subject: "<https://example.com/1>", Predicate: "<https://example.com/2>", Object: "\"1\""},
		{Subject: "<https://example.com/3>", Predicate: "<https://example.com/4>", Object: "\"2\""},
	},
		map[string]string{"ex": "https://example.com"},
	)

	if err != nil {
		t.Fatal(err)
	}

	if code != 204 {
		t.Fatalf("Unexpected status code: %d", code)
	}

	response, err := database.SendSparqlQuery(&db, `SELECT (COUNT(*) as ?cnt) WHERE {?s ?p ?o}`, database.QUERY)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	resTxt, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("failed to read result: %s", err)
	}

	var resJSON CountResult
	if err := json.Unmarshal(resTxt, &resJSON); err != nil {
		t.Fatalf("failed to unmarshal result, was there an error in the query? %s. Result was %s", err, resTxt)
	}

	expected := CountResult{
		Head: HeadVars{Vars: []string{"cnt"}},
		Results: ResultBindings{
			Bindings: []CountBinding{
				{
					Cnt: BindingVar{
						Type:     "literal",
						DataType: "http://www.w3.org/2001/XMLSchema#integer",
						Value:    "2",
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(resJSON, expected) {
		t.Fatalf("Insertion results did not match expectations, expected %v, got %v", expected, resJSON)
	}

	db.CleanDB()

	response, err = database.SendSparqlQuery(&db, `SELECT (COUNT(*) as ?cnt) WHERE {?s ?p ?o}`, database.QUERY)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	resTxt, err = io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("failed to read result: %s", err)
	}

	// var resJSON CountResult
	if err := json.Unmarshal(resTxt, &resJSON); err != nil {
		t.Fatalf("failed to unmarshal result, was there an error in the query? %s. Result was %s", err, resTxt)
	}

	expected = CountResult{
		Head: HeadVars{Vars: []string{"cnt"}},
		Results: ResultBindings{
			Bindings: []CountBinding{
				{
					Cnt: BindingVar{
						Type:     "literal",
						DataType: "http://www.w3.org/2001/XMLSchema#integer",
						Value:    "0",
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(resJSON, expected) {
		t.Fatalf("Insertion results did not match expectations, expected %v, got %v", expected, resJSON)
	}
}

func TestAddTriples(t *testing.T) {
	db := database.NewDBManager(USER, PASS, HOST, PORT, DB)
	db.CleanDB()

	code, err := db.AddTriples([]schema.Triple{
		{Subject: "<https://example.com/1>", Predicate: "<https://example.com/2>", Object: "\"1\""},
		{Subject: "<https://example.com/3>", Predicate: "<https://example.com/4>", Object: "\"2\""},
	},
		map[string]string{"ex": "https://example.com"},
	)

	if err != nil {
		t.Error(err)
	}

	if code != 204 {
		t.Errorf("Unexpected status code: %d", code)
	}
	db.CleanDB()
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
		<https://example.com/5> <https://example.com/6> ?o .
	} 
`

	if err := os.WriteFile("tmp.rq", []byte(fileData), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("tmp.rq")

	res, err := db.ExecuteQueryFile("tmp.rq")

	if err != nil {
		t.Fatal(err)
	}

	t.Log(json.MarshalIndent(res, "", "  "))

	bind := res[0]["o"].(string)

	if bind != "2" {
		t.Errorf("Result did not match: %s", res)
	}
}

func TestExecuteAttackTree(t *testing.T) {
	db := database.NewDBManager(USER, PASS, HOST, PORT, DB)
	db.CleanDB()

	db.AddTriples([]schema.Triple{
		{
			Subject:   "<http://example.com/1>",
			Predicate: "<http://example.com/2>",
			Object:    "<http://example.com/3>",
		},
	},
		map[string]string{"ex": "http://example.com"},
	)

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not get CWD: %s", err)
	}
	t.Logf("CWD: %s", dir)

	if err := os.Mkdir(".devprivops/", os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(".devprivops/")

	if err := os.Mkdir(".devprivops/test/", os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(".devprivops/test/")

	if err := os.Mkdir(".devprivops/test/root", os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(".devprivops/test/root")

	if err := os.WriteFile(".devprivops/test/root/f1.rq", []byte("SELECT * WHERE {<https://no.exists> ?p ?o}"), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(".devprivops/test/root/f1.rq")

	if err := os.WriteFile(".devprivops/test/root/f2.rq", []byte("SELECT * WHERE {<https://no.exists> ?p ?o}"), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(".devprivops/test/root/f2.rq")

	if err := os.WriteFile(".devprivops/test/root/f3.rq", []byte("SELECT * WHERE {<https://no.exists> ?p ?o}"), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(".devprivops/test/root/f3.rq")

	if err := os.WriteFile(".devprivops/test/root/f4.rq", []byte("SELECT * WHERE {?s ?p ?o}"), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(".devprivops/test/root/f4.rq")
	atkTreeFile := `
description: R
query: test/root/f1.rq
children:
  - description: C1
    query: test/root/f2.rq
    children: []
  - description: C2
    query: test/root/f3.rq
    children:
      - description: C21
        query: test/root/f4.rq
        children: []
`

	if err := os.WriteFile(".devprivops/test/atk_tree.yml", []byte(atkTreeFile), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(".devprivops/test/atk_tree.yml")

	atkTree, err := attacktree.NewAttackTreeFromYaml(".devprivops/test/atk_tree.yml", "")
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

	db.CleanDB()
}
