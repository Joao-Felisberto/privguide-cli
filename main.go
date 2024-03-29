package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	attacktree "github.com/Joao-Felisberto/devprivops/attack_tree"
	"github.com/Joao-Felisberto/devprivops/database"
	"github.com/Joao-Felisberto/devprivops/fs"
	"github.com/Joao-Felisberto/devprivops/schema"
	"github.com/Joao-Felisberto/devprivops/util"
	"github.com/spf13/cobra"
)

func loadRepresentations(dbManager *database.DBManager) error {
	fmt.Println("===\nLoading DFD into DB\n===")
	dfdFname, err := fs.GetFile("descriptions/dfd.yml")
	if err != nil {
		return err
	}
	dfdSchemaFname, err := fs.GetFile("schemas/dfd-schema.json")
	if err != nil {
		return err
	}
	dfd, err := schema.ReadYAML(
		//		fmt.Sprintf("./.%s/dfd/dfd.yml", appName),
		//		fmt.Sprintf("./.%s/dfd/dfd-schema.json", appName),
		dfdFname,
		dfdSchemaFname,
	)
	if err != nil {
		return err
	}
	statusCode, err := dbManager.AddTriples(schema.YAMLtoRDF("https://example.com/ROOT", dfd, "https://example.com/ROOT"))
	if err != nil {
		return err
	}
	if statusCode != 204 {
		return fmt.Errorf("unexpected status code: %d", statusCode)
	}

	return nil
}

/*
func reasoner(dbManager *database.DBManager) error {
	fmt.Println("===\nReasoner Rules\n===")
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
*/

func policies(dbManager *database.DBManager, regulation string) (map[string]interface{}, error) {
	fmt.Println("===\nPolicy Compliance\n===")
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
		format := q["format"].(map[interface{}]interface{})

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
			format["heading whith results"].(string),
			format["heading without results"].(string),
			format["result line"].(string),
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
			fmt.Println("error parsing query results:", err)
		}
		fmt.Printf("Violations of '%s': %s\n", pol.Title, b)
		report[pol.Title] = map[string]interface{}{
			"maximum violations": pol.MaxViolations,
			"is consistency":     pol.IsConsistency,
			"violations":         res,
		}
	}

	return report, nil
}

func attackTrees(dbManager *database.DBManager) (map[string]interface{}, error) {
	fmt.Println("===\nAttack Trees\n===")
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

func analyse(cmd *cobra.Command, args []string) error {
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
	dbManager.CleanDB()

	report := map[string]interface{}{}

	// 1. Load DFD into DB
	if err = loadRepresentations(&dbManager); err != nil {
		return err
	}

	// 2. Run all the reasoner rules
	/*
		if err = reasoner(&dbManager); err != nil {
			return err
		}
	*/

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
	dbManager.CleanDB()

	// 6. Print and store report
	jsonReport, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Println("error parsing report:", err)
	}
	fmt.Printf("==============\nReport: %s\n", jsonReport)

	if err := os.WriteFile("report.json", []byte(jsonReport), 0666); err != nil {
		return err
	}

	// 7. Check whether the violations are acceptable
	violations := validateReportInternal(&report)
	if len(violations) != 0 {
		fmt.Fprintf(os.Stderr, "There are policies with too many violations\n")
		for _, v := range violations {
			fmt.Fprintf(os.Stderr, "\t- %s\n", v)
		}
		os.Exit(1)
	}

	// 8. Send the report to the site
	gitCommit := exec.Command("git", "rev-parse", "HEAD")
	var commitOut bytes.Buffer
	gitCommit.Stdout = &commitOut

	gitBranch := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	var branchOut bytes.Buffer
	gitBranch.Stdout = &branchOut

	gitCommit.Run()
	gitBranch.Run()

	fmt.Printf("%s:%s\n", strings.Trim(branchOut.String(), "\n"), commitOut.String())

	return nil
}

func main() {
	appName := "devprivops"

	//	dfdSchema := fmt.Sprintf("./.%s/schemas/dfd-schema.json", appName)
	//	attackTreeSchema := fmt.Sprintf("./.%s/schemas/atk-tree-schema.json", appName)
	//	dfdSchema, err := fs.GetFile("schemas/dfd-schema.json")
	//	if err != nil {
	//		return fmt.Errorf("Could not find the DFD schema, is the program correctly installed? %s", err)
	//	}
	//	attackTreeSchema, err := fs.GetFile("schemas/atk-tree-schema.json")

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

	var devCmd = &cobra.Command{
		Use:   "dev",
		Short: "Development tests only",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Running development command")
			return nil
		},
	}

	//	analyseCmd.Flags().StringVar(&dfdSchema, "dfd-schema", dfdSchema, "Custom DFD schema file")
	//	analyseCmd.Flags().StringVar(&attackTreeSchema, "attack-tree-schema", attackTreeSchema, "Custom attack tree schema file")

	rootCmd.AddCommand(analyseCmd)
	rootCmd.AddCommand(devCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
