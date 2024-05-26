package database

// Describes a user story and its associated requirements
type UserStory struct {
	UseCase      string        // The use case name
	IsMisuseCase bool          // Whether it should be regarded as a misuse case
	Requirements []Requirement // The list of requirements to be satisfied
}

// Describes a requirement
type Requirement struct {
	Title       string // The requirement's title
	Description string // Its description
	Query       string // The query that encodes the requirement validation
}

// Reads the user stories from a map.
//
// `yaml`: The requirement data
//
// returns: the list of user stories or an error if there is any when the map parsed
func USFromYAML(yaml []interface{}) ([]*UserStory, error) {
	userStories := []*UserStory{}
	for _, usRaw := range yaml {
		us := usRaw.(map[interface{}]interface{})
		useCase := us["use case"].(string)
		isMisuseCase := us["is misuse case"].(bool)
		requirementsRaw := us["requirements"].([]interface{})

		requirements := []Requirement{}
		for _, reqRaw := range requirementsRaw {
			req := reqRaw.(map[interface{}]interface{})
			title := req["title"].(string)
			description := req["description"].(string)
			query := req["query"].(string)

			requirements = append(requirements, Requirement{
				Title:       title,
				Description: description,
				Query:       query,
			})
		}

		userStories = append(userStories, &UserStory{
			UseCase:      useCase,
			IsMisuseCase: isMisuseCase,
			Requirements: requirements,
		})
	}

	return userStories, nil
}
