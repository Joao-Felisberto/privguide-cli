package database

// Defines the data a query holds
type Query struct {
	File           string // The file where the query resides
	Title          string // The query's title
	Description    string // The query's purpose description
	IsConsistency  bool   // Whether the query concerns the consistency of the descriptions or not
	MaxViolations  int    // The maximum number of violations allowed
	MappingMessage string // The message instructing how to map the results of the query to solutions
}

// Constructs a new query
// `file`: The file where the query resides
// `title`: The query's title
// `description: The query's purpose description
// `isConsistency`: Whether the query concerns the consistency of the descriptions or not
// `maxViolations`: The maximum number of violations allowed
// `mappingMessage`: The message instructing how to map the results of the query to solutions
func NewQuery(
	file string,
	title string,
	description string,
	isConsistency bool,
	maxViolations int,
	mappingMessage string,
) Query {
	return Query{
		File:           file,
		Title:          title,
		Description:    description,
		IsConsistency:  isConsistency,
		MaxViolations:  maxViolations,
		MappingMessage: mappingMessage,
	}
}
