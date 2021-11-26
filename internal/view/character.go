package view

// type CharacterViewData struct {
// 	Character *graphql.Character
// 	Skills []*graphql.Skill

// }

// func (s *Service) CharacterSkillboard(ctx context.Context, characterID uint64) (*CharacterViewData, error) {
// 	skillboard, gqlErr := s.graphql.Skillboard(ctx, characterID)
// 	if gqlErr != nil {
// 		fmt.Println(gqlErr)
// 		return
// 	}

// 	groupedSkills, err := s.skill.GroupedSkills(ctx, skillboard.User.Skills.Skills)

// 	viewData := &characterViewData{
// 		Character:  skillboard.User.Character,
// 		Skills:     skillboard.User.Skills,
// 		Implants:   skillboard.User.Implants,
// 		Attributes: skillboard.User.Attributes,
// 		Queue:      skillboard.User.Queue,
// 	}
// }
