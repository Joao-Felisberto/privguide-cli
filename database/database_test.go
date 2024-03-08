package database_test

import (
	"testing"

	"github.com/Joao-Felisberto/devprivops/database"
	"github.com/Joao-Felisberto/devprivops/schema"
)

func TestAddTriples(t *testing.T) {
	db := database.NewDBManager("user", "password", "localhost", 3030, "tmp")

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
