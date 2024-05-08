package cmd

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/Joao-Felisberto/devprivops/database"
	"github.com/Joao-Felisberto/devprivops/fs"
	"github.com/Joao-Felisberto/devprivops/schema"
	"github.com/Joao-Felisberto/devprivops/util"
)

// Validates and loads the representation in the given file into the database
// The representation must abide by the provided schema.
//
// `dbManager`: The DBManager connecting to the database
//
// `repFile`: file containing the representation
//
// `schemaFile`: file containing the schema
//
// returns: error if reading or validating any file or connecting to the database or running a query fails
func loadRep(dbManager *database.DBManager, repFile string, schemaFile string) error {
	repName, err := fs.GetFile(repFile)
	if err != nil {
		return err
	}
	repSchemaFname := ""
	if schemaFile != "" {
		repSchemaFname, err = fs.GetFile(schemaFile)
		if err != nil {
			return err
		}
	}
	rep, err := schema.ReadYAML(
		repName,
		repSchemaFname,
	)
	if err != nil {
		return err
	}

	uriMetadata, err := getURIMetadata()
	if err != nil {
		return err
	}
	uris := util.Filter(*uriMetadata, func(metadata database.URIMetadata) bool {
		return util.Any(metadata.Files, func(r *regexp.Regexp) bool { return r.MatchString(repFile) })
	})
	if len(uris) == 0 {
		return fmt.Errorf("no base uri for '%s', please add it to 'uris.yml'", repFile)
	}
	uri := uris[0]
	uriMap := util.ArrayToMap(*uriMetadata, func(uri_ database.URIMetadata) (string, string) {
		return uri_.Abreviation, uri_.URI
	})

	triples := schema.YAMLtoRDF(
		fmt.Sprintf("%s/ROOT", uri.URI),
		rep,
		fmt.Sprintf("%s/ROOT", uri.URI),
		uri.URI,
		&uriMap,
	)
	statusCode, err := dbManager.AddTriples(triples, uriMap)
	if err != nil {
		return err
	}
	if statusCode != 204 {
		return fmt.Errorf("unexpected status code: %d", statusCode)
	}

	return nil
}

// Loads all representations in the representations directories
//
// `dbManager`: The DBManager connecting to the database
//
// `root`: The directory with all the representations
//
// returns: error if reading or validating any file or connecting to the database or running a query fails
func loadRepresentations(dbManager *database.DBManager, root string) error {
	slog.Debug("Getting descriptions", "root", root)
	entries, err := fs.GetDescriptions(root)
	if err != nil {
		return fmt.Errorf("error fetching description: %s", err)
	}

	for _, e := range entries {
		fPath := strings.Split(e, "/")
		fname := fPath[len(fPath)-1]

		tmp := strings.Split(fname, ".")
		schemaIndicator := tmp[len(tmp)-2]

		schema := fmt.Sprintf("schemas/%s-schema.json", schemaIndicator)

		if err := loadRep(dbManager, e, schema); err != nil {
			return err
		}
	}

	return nil
}

// Reads all the URI metadata provided in the `uris.yml` file
//
// returns the list of metadata about each URI or an error if reading the file or serializing it fails
func getURIMetadata() (*[]database.URIMetadata, error) {
	uriFile, err := fs.GetFile("uris.yml")
	if err != nil {
		return nil, err
	}

	return database.URIsFromFile(uriFile)
}
