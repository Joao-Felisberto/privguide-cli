package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	attacktree "github.com/Joao-Felisberto/devprivops/attack_tree"
	"github.com/Joao-Felisberto/devprivops/database"
	"github.com/Joao-Felisberto/devprivops/fs"
	"github.com/Joao-Felisberto/devprivops/schema"
	"github.com/Joao-Felisberto/devprivops/util"
	"github.com/spf13/cobra"
)

const REPORT_ENDPOINT_FLAG_NAME = "report-endpoint"

func getURIMetadata() (*[]database.URIMetadata, error) {
	uriFile, err := fs.GetFile("uris.yml")
	if err != nil {
		return nil, err
	}

	return database.URIsFromFile(uriFile)
}

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
		// return slices.Contains(metadata.Files, file)
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
		// map[string]string{uri.Abreviation: uri.URI},
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
		/*
			schema, err := fs.GetFile(fmt.Sprintf("schemas/%s-schema.json", schemaIndicator))
			if err != nil {
				return err

			}
		*/

		if err := loadRep(dbManager, e, schema); err != nil {
			return err
		}
	}

	/*
		fmt.Println("===\nLoading DFD into DB\n===")
		if err := loadRep(dbManager, "descriptions/dfd.yml", "schemas/dfd-schema.json"); err != nil {
			return err
		}
		fmt.Println("===\nLoading DPIA into DB\n===")
		if err := loadRep(dbManager, "descriptions/dpia.yml", "schemas/dpia-schema.json"); err != nil {
			return err
		}
	*/

	return nil
}

func reasoner(dbManager *database.DBManager) error {
	slog.Info("===Reasoner Rules===")
	reasonDir, err := fs.GetFile("reasoner")
	if err != nil {
		return err
	}
	files, err := os.ReadDir(reasonDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		fPath, err := fs.GetFile(fmt.Sprintf("reasoner/%s", file.Name()))
		if err != nil {
			return err
		}
		if err := dbManager.ExecuteReasonerRule(fPath); err != nil {
			return fmt.Errorf("could not execute reasoner rule: %s", err)
		}
	}

	return nil
}

func policies(dbManager *database.DBManager, regulation string) (map[string]interface{}, error) {
	slog.Info("===Policy Compliance===")
	polFile, err := fs.GetFile(fmt.Sprintf("regulations/%s/policies.yml", regulation))
	if err != nil {
		return nil, err
	}
	polSchema, err := fs.GetFile("schemas/query-schema.json")
	if err != nil {
		return nil, err
	}
	yamlQueries, err := schema.ReadYAML(polFile, polSchema)
	if err != nil {
		return nil, err
	}

	yamlQueriesList := yamlQueries.([]interface{})
	queries := util.Map(yamlQueriesList, func(q1 interface{}) database.Query {
		q := q1.(map[interface{}]interface{})
		// format := q["format"].(map[interface{}]interface{})

		qFile, err := fs.GetFile(q["file"].(string))
		if err != nil {
			// very beautiful, isn't it?
			panic(err)
		}
		// maxViolations, err := strconv.Atoi(q["maximum violations"].(string))
		maxViolations := q["maximum violations"].(int)
		if err != nil {
			// will never happen because the schema has already been validated
			panic(err)
		}
		return database.NewQuery(
			// fmt.Sprintf("./.%s/%s", appName, q["file"].(string)),
			qFile,
			q["title"].(string),
			q["description"].(string),
			q["is consistency"].(bool),
			maxViolations,
			// 			format["heading whith results"].(string),
			// 			format["heading without results"].(string),
			// 			format["result line"].(string),
			q["mapping message"].(string),
		)
	})
	report := map[string]interface{}{}
	for _, pol := range queries {
		res, err := dbManager.ExecuteQueryFile(pol.File)
		if err != nil {
			return nil, fmt.Errorf("error executing query from '%s': %s", pol.File, err)
		}
		// TODO: operate on the results
		b, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			slog.Error("error parsing query results:", "error", err)
		}
		slog.Info("Violations:", "policy", pol.Title, "violations", b)
		report[pol.Title] = map[string]interface{}{
			"description":        pol.Description,
			"maximum violations": pol.MaxViolations,
			"is consistency":     pol.IsConsistency,
			"violations":         res,
			"mapping message":    pol.MappingMessage,
		}
	}

	return report, nil
}

func attackTrees(dbManager *database.DBManager) (map[string]interface{}, error) {
	slog.Info("===Attack Trees===")
	atkDir, err := fs.GetFile("attack_trees/descriptions/")
	if err != nil {
		return nil, err
	}
	files, err := os.ReadDir(atkDir)
	if err != nil {
		return nil, err
	}
	atkSchema, err := fs.GetFile("schemas/atk-tree-schema.json")
	if err != nil {
		return nil, err
	}
	report := map[string]interface{}{}
	for _, file := range files {
		fPath, err := fs.GetFile(fmt.Sprintf("attack_trees/descriptions/%s", file.Name()))
		if err != nil {
			return nil, err
		}
		tree, err := attacktree.NewAttackTreeFromYaml(fPath, atkSchema)
		if err != nil {
			return nil, err
		}

		// query code, failingNode, err
		_, failingNode, err := dbManager.ExecuteAttackTree(tree)
		if err != nil {
			return nil, fmt.Errorf("error at node '%s': %s", failingNode.Description, err)
		}

		report[file.Name()] = tree
	}
	return report, nil
}

func validateReportInternal(report *map[string]interface{}) []string {
	regulations := (*report)["policies"].(map[string]interface{})
	violated := []string{}

	for _, policies := range regulations {
		for polName, policy := range policies.(map[string]interface{}) {
			// policy = policy.(map[string]interface{})
			maxViolations := policy.(map[string]interface{})["maximum violations"].(int)
			violations := len(policy.(map[string]interface{})["violations"].([]map[string]interface{}))

			if violations > maxViolations {
				violated = append(violated, polName)
			}
		}
	}

	return violated
}

func verifyRequirements(dbManager *database.DBManager) (*map[string]interface{}, error) {
	requirementsFile, err := fs.GetFile("requirements/requirements.yml")
	if err != nil {
		return nil, err
	}
	requirementsSchema, err := fs.GetFile("schemas/requirement-schema.json")
	if err != nil {
		return nil, err
	}
	requirementsRaw, err := schema.ReadYAML(requirementsFile, requirementsSchema)
	if err != nil {
		return nil, err
	}
	userStories, err := database.USFromYAML(requirementsRaw.([]interface{}))
	if err != nil {
		return nil, err
	}

	report := map[string]interface{}{}
	for _, us := range userStories {
		report[us.UseCase] = map[string]interface{}{
			"is misuse case": us.IsMisuseCase,
			"requirements":   []map[string]interface{}{},
		}
		for _, r := range us.Requirements {
			f, err := fs.GetFile(r.Query)
			if err != nil {
				return nil, err
			}
			res, err := dbManager.ExecuteQueryFile(f)
			if err != nil {
				return nil, err
			}
			if res == nil {
				res = []map[string]interface{}{}
			}
			usReport := report[us.UseCase].(map[string]interface{})
			usReport["requirements"] = append(usReport["requirements"].([]map[string]interface{}),
				map[string]interface{}{
					"title":       r.Title,
					"description": r.Description,
					"results":     res,
				},
			)
		}
	}

	return &report, nil
}

func getExtraData(dbManager *database.DBManager) (*[]map[string]interface{}, error) {
	slog.Info("===Extra Data===")

	extraDataFile, err := fs.GetFile("report_data/report_data.yml")
	if err != nil {
		return nil, err
	}
	extraDataSchema, err := fs.GetFile("schemas/report_data-schema.json")
	if err != nil {
		return nil, err
	}
	extraDataRaw, err := schema.ReadYAML(extraDataFile, extraDataSchema)
	if err != nil {
		return nil, err
	}

	extraData := extraDataRaw.([]interface{})
	report := util.Map(extraData, func(dRaw interface{}) map[string]interface{} {
		d := util.MapCast[string, interface{}](dRaw.(map[interface{}]interface{}))

		f, err := fs.GetFile(d["query"].(string))
		if err != nil {
			panic(fmt.Sprintf("Error getting query file %s", d["query"].(string)))
		}

		slog.Info("Getting extra information:", "query", f)
		d["results"], err = dbManager.ExecuteQueryFile(f)
		if err != nil {
			panic(fmt.Sprintf("Error processing query: %s", err))
		}
		resJson, err := json.Marshal(d["results"])
		if err != nil {
			panic(fmt.Sprintf("Error marshaling the results: %s", err))
		}
		slog.Info("Extra information extracted:", "info", resJson)

		delete(d, "query")
		return d
	})

	return &report, nil
}

func sendReport(url string, report *map[string]interface{}) error {
	// Define the URL
	// url := "http://localhost:8080/report"

	// Read report.json file
	reportData, err := os.ReadFile("report.json")
	if err != nil {
		return fmt.Errorf("error reading report.json: %s", err)
	}

	// Send HTTP POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reportData))
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %s", err)
	}
	defer resp.Body.Close()

	return nil
}

func analyse(cmd *cobra.Command, args []string) error {
	username := args[0]
	password := args[1]
	ip := args[2]
	port, err := strconv.Atoi(args[3])
	if err != nil {
		return err
	}
	dataset := args[4]
	reportEndpoint := cmd.Flag(REPORT_ENDPOINT_FLAG_NAME).Value.String()

	dbManager := database.NewDBManager(
		username,
		password,
		ip,
		port,
		dataset,
	)
	dbManager.CleanDB()

	report := map[string]interface{}{}

	// 1. Load DFD into DB
	if err = loadRepresentations(&dbManager, "descriptions"); err != nil {
		return err
	}

	// 2. Run all the reasoner rules
	if err = reasoner(&dbManager); err != nil {
		return err
	}

	// 3. Verify policy compliance
	report["policies"] = map[string]interface{}{}
	regulations, err := fs.GetRegulations()
	if err != nil {
		return err
	}
	for _, regulation := range regulations {
		reg := report["policies"].(map[string]interface{})
		reg[regulation] = map[string]interface{}{}
		polReport, err := policies(&dbManager, regulation)
		if err != nil {
			return err
		}
		reg[regulation] = polReport
	}

	// 4. Run all attack trees
	atkReport, err := attackTrees(&dbManager)
	if err != nil {
		return err
	}
	report["attack_trees"] = atkReport

	// 5. Clean database
	// dbManager.CleanDB()

	// 6. Print and store report
	//	gitCommit := exec.Command("git", "rev-parse", "HEAD")
	//	var commitOut bytes.Buffer
	//	gitCommit.Stdout = &commitOut

	gitBranch := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	var branchOut bytes.Buffer
	gitBranch.Stdout = &branchOut

	//	gitCommit.Run()
	gitBranch.Run()

	time := fmt.Sprint(time.Now().Unix())

	projDir, err := os.Getwd()
	if err != nil {
		return err
	}
	projPath := strings.Split(projDir, "/")
	projDir = projPath[len(projPath)-1]

	report["branch"] = strings.Trim(branchOut.String(), "\n")
	// report["time"] = commitOut.String()
	report["time"] = time
	report["project"] = projDir

	// jsonReport, err := json.MarshalIndent(report, "", "  ")

	// 7. Check whether the violations are acceptable
	violations := validateReportInternal(&report)
	if len(violations) != 0 {
		slog.Error("There are policies with too many violations\n")
		for _, v := range violations {
			slog.Error(fmt.Sprintf("\t- %s\n", v))
		}
		// os.Exit(1)
	}

	// 8. Validate whether requirements are met
	usReport, err := verifyRequirements(&dbManager)
	if err != nil {
		slog.Error("Error validating requirements", "error", err)
	}
	report["user stories"] = usReport

	// 9. Get extra data
	extraData, err := getExtraData(&dbManager)
	if err != nil {
		slog.Error("Error fetching extra report data", "error", err)
	}
	report["extra data"] = extraData

	// 10. Send the report to the site
	jsonReport, err := json.Marshal(report)
	if err != nil {
		slog.Error("error parsing report:", "error", err)
	}
	/*
		slog.Info("Report")
		clean, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(clean))
	*/
	slog.Info("Report written to 'report.json'")
	if err := os.WriteFile("report.json", []byte(jsonReport), 0666); err != nil {
		return err
	}
	if reportEndpoint != "" {
		if err := sendReport(reportEndpoint, &report); err != nil {
			return err
		}
	}

	return nil
}

func runScenario(dbManager *database.DBManager, scenario database.TestScenario) error {
	dbManager.CleanDB()
	slog.Info("Loading scenario", "scenario", scenario.StateDir)

	err := loadRepresentations(dbManager, scenario.StateDir)
	if err != nil {
		return err
	}

	for _, t := range scenario.Tests {
		slog.Info("Running test", "test", t.Query)
		res, err := dbManager.ExecuteQueryFile(t.Query)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(res, t.ExpectedResult) {
			return fmt.Errorf("result of '%s' does not match expectations: got '%s', expected '%s'", t.Query, res, t.ExpectedResult)
		}
	}

	slog.Info("All tests passed!", "scenario", scenario.StateDir)
	return nil
}

func test(cmd *cobra.Command, args []string) error {
	username := args[0]
	password := args[1]
	ip := args[2]
	port, err := strconv.Atoi(args[3])
	if err != nil {
		return err
	}
	dataset := args[4]

	dbManager := database.NewDBManager(
		username,
		password,
		ip,
		port,
		dataset,
	)

	// 1. Load data

	testFile, err := fs.GetFile("tests/spec.json")
	if err != nil {
		return err
	}

	tests, err := database.TestsFromFile(testFile)
	if err != nil {
		return err
	}

	/*
		a, err := json.MarshalIndent(tests, "", " ")
		if err != nil {
			return err
		}
		fmt.Println(string(a))
	*/

	// 2. For each scenario, run the tests
	for _, t := range tests {
		err := runScenario(&dbManager, t)
		if err != nil {
			return fmt.Errorf("test failed for scenario '%s': %s", t.StateDir, err)
		}
	}

	return nil
}

func main() {
	appName := "devprivops"
	reportEndpoint := ""

	var rootCmd = &cobra.Command{
		Use:   appName,
		Short: fmt.Sprintf("A CLI application to analyze %s", appName),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("please specify a subcommand. Use '%s --help' for usage details", appName)
		},
	}

	var analyseCmd = &cobra.Command{
		Use:   "analyse <username> <password> <database ip> <database port> <dataset>",
		Short: fmt.Sprintf("Analyse the specified database endpoint for %s", appName),
		Args:  cobra.ExactArgs(5),
		RunE:  analyse,
	}

	var testCmd = &cobra.Command{
		Use:   "test",
		Short: "Tests the queries against user-defined scenarios",
		Args:  cobra.ExactArgs(5),
		RunE:  test,
	}

	analyseCmd.Flags().StringVar(&reportEndpoint, REPORT_ENDPOINT_FLAG_NAME, "", "Endpoint where to send the final report")

	rootCmd.AddCommand(analyseCmd)
	rootCmd.AddCommand(testCmd)

	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
