package database

import "github.com/Joao-Felisberto/devprivops/util"

// Describes a user story and its associated requirements
type UserStory struct {
	UseCase      string        // The use case name
	IsMisuseCase bool          // Whether it should be regarded as a misuse case
	Requirements []Requirement // The list of requirements to be satisfied
	ClearenceLvl int           // The minimum hierarchical level required to see this in the visualizer
	Groups       []string      // The groups allowed to see this in the visualizer
}

// Describes a requirement
type Requirement struct {
	Title        string   // The requirement's title
	Description  string   // Its description
	Query        string   // The query that encodes the requirement validation
	ClearenceLvl int      // The minimum hierarchical level required to see this in the visualizer
	Groups       []string // The groups allowed to see this in the visualizer
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
		clearenceLvl := us["clearence level"].(int)
		groupsRaw := us["groups"].([]interface{})
		groups := util.Map(groupsRaw, func(raw interface{}) string { return raw.(string) })
		requirementsRaw := us["requirements"].([]interface{})

		requirements := []Requirement{}
		for _, reqRaw := range requirementsRaw {
			req := reqRaw.(map[interface{}]interface{})
			title := req["title"].(string)
			description := req["description"].(string)
			query := req["query"].(string)
			clearenceLvl := req["clearence level"].(int)
			groupsRaw := us["groups"].([]interface{})
			groups := util.Map(groupsRaw, func(raw interface{}) string { return raw.(string) })

			requirements = append(requirements, Requirement{
				Title:        title,
				Description:  description,
				Query:        query,
				ClearenceLvl: clearenceLvl,
				Groups:       groups,
			})
		}

		userStories = append(userStories, &UserStory{
			UseCase:      useCase,
			IsMisuseCase: isMisuseCase,
			ClearenceLvl: clearenceLvl,
			Groups:       groups,
			Requirements: requirements,
		})
	}

	return userStories, nil
}
