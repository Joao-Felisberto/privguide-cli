package database

type Query struct {
	File          string
	Title         string
	Description   string
	IsConsistency bool
	MaxViolations int
	// ResultsHeading string
	// EmptyHeading   string
	// ResultLine     string
	MappingMessage string
}

func NewQuery(
	file string,
	title string,
	description string,
	isConsistency bool,
	maxViolations int,
	// resultsHeading string,
	// emptyHeader string,
	// resultLine string,
	mappingMessage string,
) Query {
	return Query{
		File:          file,
		Title:         title,
		Description:   description,
		IsConsistency: isConsistency,
		MaxViolations: maxViolations,
		// 		ResultsHeading: resultsHeading,
		// 		EmptyHeading:   emptyHeader,
		// 		ResultLine:     resultLine,
		MappingMessage: mappingMessage,
	}
}
