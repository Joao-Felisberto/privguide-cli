// package for types and functions that deal with attack and harm trees.
package attacktree

import (
	"fmt"
	"reflect"

	"github.com/Joao-Felisberto/devprivops/schema"
)

// Represents the execution status of a tree node, either before or after the execution of its associated query
type ExecutionStatus int

const (
	NOT_EXECUTED ExecutionStatus = iota // The node has not yet been executed
	NOT_POSSIBLE                        // The node's condition is deemed not possible
	POSSIBLE                            // The node's condition is deemed possible
	ERROR                               // There was an error when executing the node
)

// Represents a node in the attack tree.
//
// A node is composed of a query, which is its condition, the child nodes and some metadata.
// A node is only evaluated if at least one of its pre-conditions (its children) is possible, or has no children.
type AttackNode struct {
	Description     string                    `json:"description"`      // Brief textual description of the node's condition
	Query           string                    `json:"query"`            // Path to the query that encodes the condition
	Children        []*AttackNode             `json:"children"`         // The node's pre-conditions
	ExecutionStatus ExecutionStatus           `json:"execution status"` // The current execution stats of the node, may change when the tree is executed
	ExecutionResult *[]map[string]interface{} `json:"execution result"` // The result of running the query, if it was run, else nil
}

// Represents the whole attack/harm tree.
//
// Is represented by a singular root node.
// When the root node's condition is possible, the attack/harm is deemed present in the system.
type AttackTree struct {
	Root AttackNode `json:"root"` // The root node of the attack tree
}

// Setter for the execution status.
//
// Should be called after node's condition execution to change the status accordingly
func (node *AttackNode) SetExecutionResults(status ExecutionStatus, results *[]map[string]interface{}) {
	node.ExecutionStatus = status
	node.ExecutionResult = results
}

// Construct a node from the json representation of the tree, prior to its execution.
// All nodes are initialized with ExecutionStatus `NOT_EXECUTED` and ExecutionResult `nil`.
//
// The dict structure should follow this pattern:
//
//	{
//		"description": "some text",
//		"query": "path to the query file",
//		"children": [] // more nodes like this one in the array
//	}
//
// Calls to this method should pass a root node, the children are processed recursivelly.
//
// `data`: the node represented by a go map
//
// returns: the parsed AttackNode, or an error when:
//
//   - there was an error parsing the child node
//   - required fields are missing from the dict
//   - the node is passed in an incorrect data type
func parseNode(data interface{}) (*AttackNode, error) {
	switch node := data.(type) {
	case map[interface{}]interface{}:
		description, descOk := node["description"].(string)
		query, queryOk := node["query"].(string)
		childrenData, childrenOk := node["children"].([]interface{})

		if !descOk || !queryOk || !childrenOk {
			return nil, fmt.Errorf("missing required fields in node")
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
			ExecutionResult: nil,
		}, nil
	default:
		return nil, fmt.Errorf("invalid node data type: %s", reflect.TypeOf(data))
	}
}

// Constructs a full attack/harm tree struct from the YAML description in a file.
//
// `yamlFile`: The file in which the tree is represented
//
// `atkTreeSchema`: The file with the attack tree schema
//
// returns: The parsed ATtackTree, or an error if
//   - The `schema.ReadYAML` call fails
//   - The node could not be parsed
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
