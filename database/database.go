// Contains all the data types and utilities to communicate with the database
package database

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"text/template"

	attacktree "github.com/Joao-Felisberto/devprivops/attack_tree"
	"github.com/Joao-Felisberto/devprivops/fs"
	"github.com/Joao-Felisberto/devprivops/schema"
)

// todo: sanitization https://stackoverflow.com/a/55726984

// The types of query methods allowed by the SparQL1.1 over HTTP standard
type QueryMethod string

const (
	QUERY  QueryMethod = "query"  // A generic query method
	UPDATE QueryMethod = "update" // An update method
	DATA   QueryMethod = "data"   // A method that adds raw triples to the database
	UPLOAD QueryMethod = "upload" // A method that uploads a file with triples
)

// Models the data needed for each database connection to a triple store that supports HTTP authentication
type DBManager struct {
	username string // the username
	password string // the password
	ip       string // the triple store's IP
	port     int    // the triple store's port
	dataset  string // the dataset to which to connect
}

// Creates a new DBManager instance from which it is possible to communicate with the trile store
//
// `username`: the username
//
// `password`: the password
//
// `ip`: the triple store's IP
//
// `port`: the triple store's port
//
// `dataset`: the dataset to which to connect

func NewDBManager(
	username string,
	password string,
	ip string,
	port int,
	dataset string,
) DBManager {
	return DBManager{
		username,
		password,
		ip,
		port,
		dataset,
	}
}

// Sends a sparql query in a query with a specific method
//
// `query`: the query to send
//
// `method`: the method to send it with
//
// returns: the query response or the error that occured whrn executing the query
func (db *DBManager) sendSparqlQuery(query string, method QueryMethod) (*http.Response, error) {
	slog.Debug("Sending query", "query", query)
	endpoint := fmt.Sprintf("http://%s:%d/%s/%s", db.ip, db.port, db.dataset, method)
	client := &http.Client{}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(query)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", fmt.Sprintf("application/sparql-%s", method))
	req.Header.Set("Accept", "application/json")

	auth := db.username + ":" + db.password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", basicAuth)

	return client.Do(req)
}

// Removes all triples from the triple store
//
// returns: the query response or the error that occured whrn executing the query
func (db *DBManager) CleanDB() (*http.Response, error) {
	return db.sendSparqlQuery(`
		DELETE { 
			?s ?p ?o 
		} WHERE { 
			?s ?p ?o 
		}
	`, UPDATE,
	)
}

// Adds the list of triples to the triple store
//
// `triples`: the triples to add
//
// `prefixes`: the map of prefix abreviations to the full prefix URI
//
// returns: the status code or an error if the query failed
func (db *DBManager) AddTriples(triples []schema.Triple, prefixes map[string]string) (int, error) {
	sparqlTemplate := `
		{{ range $key, $value := .Prefixes }} PREFIX {{ $key }}: <{{ $value }}>
		{{ end }}

        INSERT DATA {
		{{ range .Triples }}{{ .Subject }} {{ .Predicate }} {{ .Object }} .
		{{ end }}
        }
    `
	var sparqlQuery strings.Builder

	tpl := template.Must(template.New("insert triples").Parse(sparqlTemplate))
	if err := tpl.Execute(&sparqlQuery, struct {
		Triples  []schema.Triple
		Prefixes map[string]string
	}{
		triples,
		prefixes,
	}); err != nil {
		return -1, err
	}

	// fmt.Printf("Sending %s\n", sparqlQuery.String())

	response, err := db.sendSparqlQuery(sparqlQuery.String(), UPDATE)
	if err != nil {
		return -1, fmt.Errorf("error sending SPARQL query: %s", err)
	}
	defer response.Body.Close()

	resTxt, err := io.ReadAll(response.Body)
	if err != nil {
		return -1, fmt.Errorf("failed to read result of AddTriples: %s", err)
	}

	slog.Debug("AddTriples resopnse", "body", resTxt)
	return response.StatusCode, nil
}

// Executes a single reasoner rule
//
// `file`: the file where the reasoner rule resides
//
// returns: an error if reading or validating the file or running the query result in an error
func (db *DBManager) ExecuteReasonerRule(file string) error {
	sparqlQueryBytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("could not read rule file '%s': %s", file, err)
	}

	sparqlQuery := string(sparqlQueryBytes)

	slog.Debug("Executing reasoner rule", "rule", sparqlQuery)

	response, err := db.sendSparqlQuery(sparqlQuery, UPDATE)
	if err != nil {
		return fmt.Errorf("query from '%s' had db errors: %s", file, err)
	}
	defer response.Body.Close()

	return nil
}

// Executes a single query from a file
//
// `file`: the file where the reasoner rule resides
//
// returns: the execution results or an error if reading or validating the file or running the query result in an error
func (db *DBManager) ExecuteQueryFile(file string) ([]map[string]interface{}, error) {
	sparqlQueryBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file '%s': %s", file, err)
	}

	sparqlQuery := string(sparqlQueryBytes)

	response, err := db.sendSparqlQuery(sparqlQuery, QUERY)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query '%s': %s", file, err)
	}
	defer response.Body.Close()

	resTxt, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read result of '%s': %s", file, err)
	}

	var resJSON map[string]interface{}
	if err := json.Unmarshal(resTxt, &resJSON); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result of '%s', was there an error in the query? %s. Result was %s", file, err, resTxt)
	}

	results, ok := resJSON["results"].(map[string]interface{})
	if !ok {
		return nil, errors.New("results not found in response")
	}

	bindings, ok := results["bindings"].([]interface{})
	if !ok {
		return nil, errors.New("bindings not found in response")
	}

	binds := []map[string]interface{}{}
	for _, bind := range bindings {
		bindMap, ok := bind.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid binding format")
		}
		for k := range bindMap {
			bindMap[k] = bindMap[k].(map[string]interface{})["value"]
		}
		binds = append(binds, bindMap)
	}

	return binds, nil
}

// Executes the query of an attack/harm tree node if it is reachable.
//
// `attackNode`: The note whose query is to be executed
//
// returns: The execution results, the node that failed previously and the error that caused its failure.
// Errors can occur when reading or validating the query file or executing the query.
func (db *DBManager) executeAttackTreeNode(attackNode *attacktree.AttackNode) ([]map[string]interface{}, *attacktree.AttackNode, error) {
	// attackNode.ExecutionStatus = -1
	thisNodeIsReachable := len(attackNode.Children) == 0
	for _, node := range attackNode.Children {
		response, failingNode, err := db.executeAttackTreeNode(node)
		if err != nil {
			return response, failingNode, err
		}
		// fmt.Printf("- %s\n", response)
		if len(response) != 0 {
			thisNodeIsReachable = true
		}
	}
	if thisNodeIsReachable {
		slog.Info("Executing attack node:", "attack node", attackNode.Description)
		qFile, err := fs.GetFile(attackNode.Query)
		if err != nil {
			return nil, attackNode, err
		}
		binds, qErr := db.ExecuteQueryFile(qFile)

		if len(binds) == 0 {
			slog.Info("NOT POSSIBLE", "node", attackNode.Description)
			attackNode.SetExecutionStatus(attacktree.NOT_POSSIBLE)
		} else {
			slog.Info("POSSIBLE", "node", attackNode.Description)
			attackNode.SetExecutionStatus(attacktree.POSSIBLE)
		}

		return binds, attackNode, qErr
	}
	slog.Info("UNREACHABLE", "node", attackNode.Description)
	return nil, nil, nil
}

// Finds out whether the attack/harm described by the tree is possible in the system.
//
// `attackTree`: The tree to be executed
//
// returns: The execution results, the node that failed previously and the error that caused its failure.
// Errors can occur when reading or validating the query file or executing the query.
func (db *DBManager) ExecuteAttackTree(attackTree *attacktree.AttackTree) ([]map[string]interface{}, *attacktree.AttackNode, error) {
	return db.executeAttackTreeNode(&attackTree.Root)
}

// Applies the current configuration to the descrption already in the triple store.
//
// This means the queries do not have to take into account parts of the system that might be configurable,
// as the identifiers to configuration variables are replaced by the objects they point to in the config
//
// returns: The response of the query or an error, in case it fails to execute
func (db *DBManager) ApplyConfig() (*http.Response, error) {
	return db.sendSparqlQuery(`
PREFIX cfg: <https://devprivops.com/config/>

DELETE {
  ?s ?p ?o .
}
INSERT {
  ?s ?p ?newValue .
}
WHERE {
  ?s ?p ?o .
  ?o cfg:value ?newValue .
}
`,
		UPDATE,
	)
}
