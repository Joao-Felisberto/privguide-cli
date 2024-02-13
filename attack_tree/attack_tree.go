package attacktree

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Joao-Felisberto/devprivops/schema"
)

type AttackNode struct {
	Description string
	Query       string
	Children    []AttackNode
}

type AttackTree struct {
	Root AttackNode
}

func parseNode(data interface{}) (*AttackNode, error) {
	switch node := data.(type) {
	case map[string]interface{}: // , map[interface{}]interface{}:
		description, descOk := node["description"].(string)
		query, queryOk := node["query"].(string)
		childrenData, childrenOk := node["children"].([]interface{})

		if !descOk || !queryOk || !childrenOk {
			return nil, errors.New("missing required fields in node")
		}

		children := make([]AttackNode, len(childrenData))
		for i, childData := range childrenData {
			childNode, err := parseNode(childData)
			if err != nil {
				return nil, fmt.Errorf("error parsing child node: %w", err)
			}
			children[i] = *childNode
		}

		return &AttackNode{
			Description: description,
			Query:       query,
			Children:    children,
		}, nil
	case map[interface{}]interface{}:
		description, descOk := node["description"].(string)
		query, queryOk := node["query"].(string)
		childrenData, childrenOk := node["children"].([]interface{})

		if !descOk || !queryOk || !childrenOk {
			return nil, errors.New("missing required fields in node")
		}

		children := make([]AttackNode, len(childrenData))
		for i, childData := range childrenData {
			childNode, err := parseNode(childData)
			if err != nil {
				return nil, fmt.Errorf("error parsing child node: %w", err)
			}
			children[i] = *childNode
		}

		return &AttackNode{
			Description: description,
			Query:       query,
			Children:    children,
		}, nil
	default:
		return nil, fmt.Errorf("invalid node data type: %s", reflect.TypeOf(data))
	}
}

func NewAttackTreeFromYaml(yamlFile string, atkTreeSchema string) (*AttackTree, error) {
	yamlTree, err := schema.ReadYAML(yamlFile, atkTreeSchema)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("%s\n", reflect.TypeOf(yamlTree))

	rootNode, err := parseNode(yamlTree)
	if err != nil {
		return nil, err
	}
	return &AttackTree{Root: *rootNode}, nil
}
