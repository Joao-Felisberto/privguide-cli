package database

type Query struct {
	File           string
	Title          string
	Description    string
	MaxViolations  int
	ResultsHeading string
	EmptyHeading   string
	ResultLine     string
}

func NewQuery(
	file string,
	title string,
	description string,
	maxViolations int,
	resultsHeading string,
	emptyHeader string,
	resultLine string,
) Query {
	return Query{
		File:           file,
		Title:          title,
		Description:    description,
		MaxViolations:  maxViolations,
		ResultsHeading: resultsHeading,
		EmptyHeading:   emptyHeader,
		ResultLine:     resultLine,
	}
}
