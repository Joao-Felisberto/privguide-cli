package database

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
)

// todo: sanitization https://stackoverflow.com/a/55726984

func sendSparqlQuery(endpoint, query, username, password string) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(query)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/sparql-update")
	req.Header.Set("Accept", "application/json")

	auth := username + ":" + password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", basicAuth)

	return client.Do(req)
}

func TestDB(dataset string, ip string, port int64, username string, password string) {
	fmt.Printf("ip=%s port=%d username=%s password=%s\n", ip, port, username, password)

	endpointURL := fmt.Sprintf("http://%s:%d/%s/update", ip, port, dataset)
	sparqlQuery := `
        PREFIX foaf: <http://xmlns.com/foaf/0.1/>
        INSERT DATA {
            <http://example.org/JaneDane> foaf:name "Jane Dane" ;
                                          foaf:email <mailto:jane@example.org> .
        }
    `

	response, err := sendSparqlQuery(endpointURL, sparqlQuery, username, password)
	if err != nil {
		fmt.Println("Error sending SPARQL query:", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Triple inserted successfully.")
	} else {
		fmt.Println("Error inserting triple. Status code:", response.StatusCode)
	}
}
