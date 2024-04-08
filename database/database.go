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
	"path/filepath"
	"strings"
	"text/template"

	attacktree "github.com/Joao-Felisberto/devprivops/attack_tree"
	"github.com/Joao-Felisberto/devprivops/fs"
	"github.com/Joao-Felisberto/devprivops/schema"
)

// todo: sanitization https://stackoverflow.com/a/55726984

type QueryMethod string

const (
	QUERY  QueryMethod = "query"
	UPDATE QueryMethod = "update"
	DATA   QueryMethod = "data"
	UPLOAD QueryMethod = "upload"
)

type DBManager struct {
	username string
	password string
	ip       string
	port     int
	dataset  string
}

// var id_cnt = 0

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

func (db *DBManager) AddTriples(triples []schema.Triple) (int, error) {
	sparqlTemplate := `
		PREFIX ex: <https://example.com/>
        INSERT DATA {
		{{ range . }}{{ .Subject }} {{ .Predicate }} {{ .Object }} .
		{{ end }}
        }
    `
	var sparqlQuery strings.Builder

	tpl := template.Must(template.New("insert triples").Parse(sparqlTemplate))
	if err := tpl.Execute(&sparqlQuery, triples); err != nil {
		return -1, err
	}

	// fmt.Printf("Sending %s\n", sparqlQuery.String())

	response, err := db.sendSparqlQuery(sparqlQuery.String(), UPDATE)
	if err != nil {
		return -1, fmt.Errorf("error sending SPARQL query: %s", err)
	}
	defer response.Body.Close()

	return response.StatusCode, nil
}

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

	var binds []map[string]interface{}
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

	// fmt.Printf("BEFORE: %d\n", len(binds))
	return binds, nil
}

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

func (db *DBManager) ExecuteAttackTree(attackTree *attacktree.AttackTree) ([]map[string]interface{}, *attacktree.AttackNode, error) {
	return db.executeAttackTreeNode(&attackTree.Root)
}

func FindQueryFiles(rootDir string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".rq" {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
