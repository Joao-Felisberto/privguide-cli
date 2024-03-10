package database

type Query struct {
	File           string
	Title          string
	Description    string
	ResultsHeading string
	EmptyHeading   string
	ResultLine     string
}

func NewQuery(
	file,
	title,
	description,
	resultsHeading,
	emptyHeader,
	resultLine string,
) Query {
	return Query{
		File:           file,
		Title:          title,
		Description:    description,
		ResultsHeading: resultsHeading,
		EmptyHeading:   emptyHeader,
		ResultLine:     resultLine,
	}
}
