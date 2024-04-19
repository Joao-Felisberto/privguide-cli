package attacktree

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Joao-Felisberto/devprivops/schema"
)

type ExecutionStatus int

const (
	NOT_EXECUTED ExecutionStatus = iota
	NOT_POSSIBLE
	POSSIBLE
	ERROR
)

type AttackNode struct {
	Description     string
	Query           string
	Children        []*AttackNode
	ExecutionStatus ExecutionStatus
}

type AttackTree struct {
	Root AttackNode
}

func (node *AttackNode) SetExecutionStatus(status ExecutionStatus) {
	node.ExecutionStatus = status
}

func parseNode(data interface{}) (*AttackNode, error) {
	switch node := data.(type) {
	/*
		case map[string]interface{}: // , map[interface{}]interface{}:
			description, descOk := node["description"].(string)
			query, queryOk := node["query"].(string)
			childrenData, childrenOk := node["children"].([]interface{})

			if !descOk || !queryOk || !childrenOk {
				return nil, errors.New("missing required fields in node")
			}

			children := make([]*AttackNode, len(childrenData))
			for i, childData := range childrenData {
				childNode, err := parseNode(childData)
				if err != nil {
					return nil, fmt.Errorf("error parsing child node: %w", err)
				}
				children[i] = childNode
			}

			return &AttackNode{
				Description:     description,
				Query:           query,
				Children:        children,
				ExecutionStatus: NOT_EXECUTED,
			}, nil
	*/
	case map[interface{}]interface{}:
		description, descOk := node["description"].(string)
		query, queryOk := node["query"].(string)
		childrenData, childrenOk := node["children"].([]interface{})

		if !descOk || !queryOk || !childrenOk {
			return nil, errors.New("missing required fields in node")
		}

		children := make([]*AttackNode, len(childrenData))
		for i, childData := range childrenData {
			childNode, err := parseNode(childData)
			if err != nil {
				return nil, fmt.Errorf("error parsing child node: %w", err)
			}
			children[i] = childNode
		}

		return &AttackNode{
			Description:     description,
			Query:           query,
			Children:        children,
			ExecutionStatus: NOT_EXECUTED,
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

	rootNode, err := parseNode(yamlTree)
	if err != nil {
		return nil, err
	}

	return &AttackTree{Root: *rootNode}, nil
}
