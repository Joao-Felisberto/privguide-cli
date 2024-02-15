package database

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"

	attacktree "github.com/Joao-Felisberto/devprivops/attack_tree"
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

func New(
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

/*
func (db *DBManager) TestDB() (int, error) {
	sparqlQuery := `
        PREFIX foaf: <http://xmlns.com/foaf/0.1/>
        INSERT DATA {
            <http://example.org/JaneDane> foaf:name "Jane Dane" ;
                                          foaf:email <mailto:jane@example.org> .
        }
    `

	response, err := db.sendSparqlQuery(sparqlQuery, "update")
	if err != nil {
		fmt.Println("Error sending SPARQL query:", err)
		return -1, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Triple inserted successfully.")
	} else {
		fmt.Println("Error inserting triple. Status code:", response.StatusCode)
	}

	return response.StatusCode, nil
}
*/

func (db *DBManager) AddTriples(triples []schema.Triple) (int, error) {
	sparqlTemplate := `
        INSERT DATA {
			{{ range . }}
				{{ .Subject }} {{ .Predicate }} {{ .Object }} .
			{{ end }}
        }
    `
	var sparqlQuery strings.Builder

	tpl := template.Must(template.New("insert triples").Parse(sparqlTemplate))
	if err := tpl.Execute(&sparqlQuery, triples); err != nil {
		fmt.Println("ERROR could not instantiate template")
	}

	response, err := db.sendSparqlQuery(sparqlQuery.String(), UPDATE)
	if err != nil {
		fmt.Println("Error sending SPARQL query:", err)
		return -1, err
	}
	defer response.Body.Close()

	return response.StatusCode, nil
}

func (db *DBManager) executeQueryFile(file string, method QueryMethod) (int, error) {

	sparqlQueryBytes, err := os.ReadFile(file)
	if err != nil {
		return -1, err
	}

	sparqlQuery := string(sparqlQueryBytes)

	response, err := db.sendSparqlQuery(sparqlQuery, method)
	if err != nil {
		fmt.Println("Error sending SPARQL query:", err)
		return -1, err
	}
	defer response.Body.Close()

	return response.StatusCode, nil
}

func (db *DBManager) executeAttackTreeNode(attackNode *attacktree.AttackNode) (int, *attacktree.AttackNode, error) {
	for _, node := range attackNode.Children {
		code, failingNode, err := db.executeAttackTreeNode(&node)
		if err != nil {
			return code, failingNode, err
		}
		qCode, qErr := db.executeQueryFile(attackNode.Query, QUERY)
		if qErr != nil {
			return qCode, attackNode, qErr
		}
	}
	return -1, nil, nil
}

func (db *DBManager) ExecuteAttackTree(attackTree *attacktree.AttackTree) (int, *attacktree.AttackNode, error) {
	return db.executeAttackTreeNode(&attackTree.Root)
}
