package database

type UserStory struct {
	UseCase      string
	IsMisuseCase bool
	Requirements []Requirement
}

type Requirement struct {
	Title       string
	Description string
	Query       string
}

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
