package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Joao-Felisberto/devprivops/database"
	"github.com/Joao-Felisberto/devprivops/fs"
	"github.com/Joao-Felisberto/devprivops/schema"
	"github.com/Joao-Felisberto/devprivops/util"
)

func loadRep(dbManager *database.DBManager, file string, schemaFile string) error {
	repName, err := fs.GetFile(file)
	if err != nil {
		return err
	}
	repSchemaFname, err := fs.GetFile(schemaFile)
	if err != nil {
		return err
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
		return util.Any(metadata.Files, func(r *regexp.Regexp) bool { return r.MatchString(file) })
	})
	if len(uris) == 0 {
		return fmt.Errorf("no base uri for '%s', please add it to 'uris.yml'", file)
	}
	uri := uris[0]
	uriMap := util.MapToMap(*uriMetadata, func(uri_ database.URIMetadata) (string, string) {
		return uri_.Abreviation, uri_.URI
	})

	statusCode, err := dbManager.AddTriples(schema.YAMLtoRDF(
		fmt.Sprintf("%s/ROOT", uri.URI),
		rep,
		fmt.Sprintf("%s/ROOT", uri.URI),
		uri.URI,
		&uriMap,
	),
		uriMap,
	)
	if err != nil {
		return err
	}
	if statusCode != 204 {
		return fmt.Errorf("unexpected status code: %d", statusCode)
	}

	return nil
}

func loadRepresentations(dbManager *database.DBManager, root string) error {
	entries, err := fs.GetDescriptions(root)
	if err != nil {
		return err
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

func getURIMetadata() (*[]database.URIMetadata, error) {
	uriFile, err := fs.GetFile("uris.yml")
	if err != nil {
		return nil, err
	}

	return database.URIsFromFile(uriFile)
}
