package main

import (
	"fmt"
	"os"
	"strconv"

	attacktree "github.com/Joao-Felisberto/devprivops/attack_tree"
	"github.com/Joao-Felisberto/devprivops/database"
	"github.com/spf13/cobra"
)

func main() {
	appName := "devprivops"

	dfdSchema := fmt.Sprintf("./.%s/schemas/dfd-schema.json", appName)
	attackTreeSchema := fmt.Sprintf("./.%s/schemas/atk-tree-schema.json", appName)

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
		RunE: func(cmd *cobra.Command, args []string) error {
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

			// 1. Load DFD into DB
			/*
				dfd, err := schema.ReadYAML(fmt.Sprintf("./.%s/dfd/dfd.yml", appName), "") // TODO: add schema
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
			*/

			// 2. Run all the reasoner rules
			reasonDir := fmt.Sprintf("./.%s/reasoner/", appName)
			files, err := os.ReadDir(reasonDir)
			if err != nil {
				return err
			}
			for _, file := range files {
				fPath := fmt.Sprintf("./.%s/reasoner/%s", appName, file.Name())
				if err := dbManager.ExecuteReasonerRule(fPath); err != nil {
					return err
				}
			}

			// 3. Run all attack trees
			atkDir := fmt.Sprintf("./.%s/attack_trees/", appName)
			files, err = os.ReadDir(atkDir)
			if err != nil {
				return err
			}
			for _, file := range files {
				fPath := fmt.Sprintf("./.%s/attack_trees/%s", appName, file.Name())
				tree, err := attacktree.NewAttackTreeFromYaml(fPath, "") // TODO schema
				if err != nil {
					fmt.Printf("ERROR!!!! %s: %s\n", fPath, err)
					return err
				}

				// query code, failingNode, err
				_, failingNode, err := dbManager.ExecuteAttackTree(tree)
				if err != nil {
					return fmt.Errorf("error at node '%s': %s", failingNode.Description, err)
				}
			}

			// fmt.Printf("Analyzing database at endpoint: %s:%d\n", ip, port)
			return nil
		},
	}

	var devCmd = &cobra.Command{
		Use:   "dev",
		Short: "Development tests only",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Running development command")
			return nil
		},
	}

	analyseCmd.Flags().StringVar(&dfdSchema, "dfd-schema", dfdSchema, "Custom DFD schema file")
	analyseCmd.Flags().StringVar(&attackTreeSchema, "attack-tree-schema", attackTreeSchema, "Custom attack tree schema file")

	rootCmd.AddCommand(analyseCmd)
	rootCmd.AddCommand(devCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
