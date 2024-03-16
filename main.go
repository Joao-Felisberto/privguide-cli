package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	attacktree "github.com/Joao-Felisberto/devprivops/attack_tree"
	"github.com/Joao-Felisberto/devprivops/database"
	"github.com/Joao-Felisberto/devprivops/fs"
	"github.com/Joao-Felisberto/devprivops/schema"
	"github.com/Joao-Felisberto/devprivops/util"
	"github.com/spf13/cobra"
)

func execute(cmd *cobra.Command, args []string) error {
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

	report := map[string]interface{}{}

	// 1. Load DFD into DB
	fmt.Println("===\nLoading DFD into DB\n===")
	dfdFname, err := fs.GetFile("dfd/dfd.yml")
	if err != nil {
		return err
	}
	dfdSchemaFname, err := fs.GetFile("dfd-schema.json")
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

	// 2. Run all the reasoner rules
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

	// 3. Verify policy compliance
	fmt.Println("===\nPolicy Compliance\n===")
	/*
		polDir := fmt.Sprintf("./.%s/policies/", appName)
		polFiles, err := database.FindQueryFiles(polDir)
		if err != nil {
			return err
		}
		for _, pol := range polFiles {
			res, err := dbManager.ExecuteQueryFile(pol)
			if err != nil {
				return err
			}
			// TODO: operate on the results
			fmt.Printf("%s\n", res)
		}
	*/
	//	polFile := fmt.Sprintf("policies/policies.yml", appName)
	//	polschema := fmt.Sprintf("query-schema.json", appName)
	polFile, err := fs.GetFile("policies/policies.yml")
	if err != nil {
		return err
	}
	polSchema, err := fs.GetFile("query-schema.json")
	if err != nil {
		return err
	}
	yamlQueries, err := schema.ReadYAML(polFile, polSchema)
	if err != nil {
		return err
	}
	// fmt.Printf("MAIN: %s\n", yamlQueries)
	// yamlQueriesList := yamlQueries.([]map[string]interface{})
	// queries := util.Map(yamlQueriesList, func(q map[string]interface{}) database.Query {
	yamlQueriesList := yamlQueries.([]interface{})
	queries := util.Map(yamlQueriesList, func(q1 interface{}) database.Query {
		q := q1.(map[interface{}]interface{})
		format := q["format"].(map[interface{}]interface{})

		qFile, err := fs.GetFile(q["file"].(string))
		if err != nil {
			// very beautiful, isn't it?
			panic(err)
		}
		return database.NewQuery(
			// fmt.Sprintf("./.%s/%s", appName, q["file"].(string)),
			qFile,
			q["title"].(string),
			q["description"].(string),
			format["heading whith results"].(string),
			format["heading without results"].(string),
			format["result line"].(string),
		)
	})
	report["policies"] = map[string]interface{}{}
	for _, pol := range queries {
		res, err := dbManager.ExecuteQueryFile(pol.File)
		if err != nil {
			return fmt.Errorf("error executing query from '%s': %s", pol.File, err)
		}
		// TODO: operate on the results
		b, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			fmt.Println("error parsing query results:", err)
		}
		fmt.Printf("Violations of '%s': %s\n", pol.Title, b)
		resReport := report["policies"].(map[string]interface{})
		resReport[pol.Title] = res
	}

	// 4. Verify contract compliance
	/*
		contractDir := fmt.Sprintf("./.%s/contracts/", appName)
		contractFiles, err := database.FindQueryFiles(contractDir)
		if err != nil {
			return err
		}
		for _, con := range contractFiles {
			res, err := dbManager.ExecuteQueryFile(con)
			if err != nil {
				return err
			}
			// TODO: operate on the results
			fmt.Printf("%s\n", res)
		}
	*/
	/* // TODO uncomment
	fmt.Println("===\nContract Compliance\n===")
	contractFile := fmt.Sprintf("./.%s/contracts/contracts.yml", appName)
	yamlQueries, err = schema.ReadYAML(contractFile, polschema)
	if err != nil {
		return err
	}
	// yamlQueriesList = yamlQueries.([]map[string]interface{})
	// queries = util.Map(yamlQueriesList, func(q map[string]interface{}) database.Query {
	yamlQueriesList = yamlQueries.([]interface{})
	queries = util.Map(yamlQueriesList, func(q1 interface{}) database.Query {
		q := q1.(map[interface{}]interface{})
		format := q["format"].(map[interface{}]interface{})
		return database.NewQuery(
			// q["file"].(string),
			fmt.Sprintf("./.%s/%s", appName, q["file"].(string)),
			q["title"].(string),
			q["description"].(string),
			format["heading whith results"].(string),
			format["heading without results"].(string),
			format["result line"].(string),
		)
	})
	for _, contract := range queries {
		res, err := dbManager.ExecuteQueryFile(contract.File)
		if err != nil {
			return err
		}
		// TODO: operate on the results
		fmt.Printf("%s\n", res)
	}
	*/
	// 5. Run all attack trees
	fmt.Println("===\nAttack Trees\n===")
	// atkDir := fmt.Sprintf("./.%s/attack_trees/", appName)
	atkDir, err := fs.GetFile("attack_trees/")
	if err != nil {
		return err
	}
	files, err = os.ReadDir(atkDir)
	if err != nil {
		return err
	}
	//	atkSchema := fmt.Sprintf("./.%s/atk-tree-schema.json", appName)
	atkSchema, err := fs.GetFile("atk-tree-schema.json")
	if err != nil {
		return err
	}
	report["attack_trees"] = map[string]interface{}{}
	for _, file := range files {
		// fPath := fmt.Sprintf("./.%s/attack_trees/%s", appName, file.Name())
		fPath, err := fs.GetFile(fmt.Sprintf("attack_trees/%s", file.Name()))
		if err != nil {
			return err
		}
		tree, err := attacktree.NewAttackTreeFromYaml(fPath, atkSchema)
		if err != nil {
			// fmt.Printf("ERROR!!!! %s: %s\n", fPath, err)
			return err
		}

		// query code, failingNode, err
		_, failingNode, err := dbManager.ExecuteAttackTree(tree)
		if err != nil {
			return fmt.Errorf("error at node '%s': %s", failingNode.Description, err)
		}

		//		report["attack_trees"].(map[string]interface{})[tree.Root.Query], err = json.Marshal(tree)
		//		if err != nil {
		//			fmt.Println("error parsing attack tree:", err)
		//		}
		report["attack_trees"].(map[string]interface{})[file.Name()] = tree
	}

	// fmt.Printf("Analyzing database at endpoint: %s:%d\n", ip, port)
	jsonReport, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Println("error parsing report:", err)
	}
	fmt.Printf("==============\nReport: %s\n", jsonReport)

	if err := os.WriteFile("report.json", []byte(jsonReport), 0666); err != nil {
		return err
	}

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
		RunE:  execute,
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
